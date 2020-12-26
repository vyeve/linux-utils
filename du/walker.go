package du

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
)

type DU interface {
	Sum(dir string) ( /*files*/ int64 /*size*/, string)
}

type walker struct {
	semaphore chan struct{}
}

func NewWalker() DU {
	return &walker{
		semaphore: make(chan struct{}, maxRoutines),
	}
}

func (w *walker) Sum(dir string) ( /*files*/ int64 /*size*/, string) {
	resCh := make(chan int64, 1)
	var wg sync.WaitGroup
	wg.Add(1)
	go w.walkDir(dir, &wg, resCh)
	go func() {
		wg.Wait()
		close(resCh)
	}()
	var nFiles, nBytes int64
	for s := range resCh {
		nFiles++
		nBytes += s
	}
	var size string
	switch {
	case nBytes/1e9 > 0:
		size = fmt.Sprintf("%.3fGi", float64(nBytes)/math.Pow(1024, 3))
	case nBytes/1e6 > 0:
		size = fmt.Sprintf("%.3fMi", float64(nBytes)/math.Pow(1024, 2))
	case nBytes/1e3 > 0:
		size = fmt.Sprintf("%.3fKi", float64(nBytes)/math.Pow(1024, 1))
	default:
		size = fmt.Sprintf("%d", nBytes)
	}
	return nFiles, size
}

func (w *walker) walkDir(dir string, wg *sync.WaitGroup, resCh chan<- int64) {
	defer wg.Done()
	for _, f := range w.listFiles(dir) {
		resCh <- f.Size()
		if f.IsDir() {
			wg.Add(1)
			go w.walkDir(filepath.Join(dir, f.Name()), wg, resCh)
			// continue
		}
		// resCh <- f.Size()
	}
}

func (w *walker) listFiles(dir string) []os.FileInfo {
	w.semaphore <- struct{}{}
	defer func() {
		<-w.semaphore
	}()
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		if os.IsPermission(err) {
			return nil
		}
		log.Printf("Failed to read dir [%s]. err: %v", dir, err)
		return nil
	}
	return files
}

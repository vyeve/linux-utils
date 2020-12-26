package main

import (
	"fmt"
	"os"
	"time"

	"github.com/vyeve/linux-utils/du"
)

func main() {
	args := os.Args
	var dir string
	if len(args) < 2 {
		dir = "./"
	} else {

		dir = args[1]
	}
	w := du.NewWalker()
	tn := time.Now()
	fn, fs := w.Sum(dir)
	fmt.Printf("TotalTime: %s\n", time.Since(tn))
	fmt.Printf("Dir: %s\nFiles: %d\nSize: %s\n", dir, fn, fs)
}

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/devplayg/encdec"
)

func main() {
	t := time.Now()
	args := os.Args[1:]
	wg := new(sync.WaitGroup)

	for _, target := range args {
		files, err := filepath.Glob(target)
		if err != nil {
			fmt.Errorf(err.Error())
		}

		for _, f := range files {
			basename, _ := filepath.Abs(f)

			wg.Add(1)
			go func(basename string) {
				_, dur, err := encdec.Encrypt(basename)
				if err != nil {
					fmt.Println(err.Error())
				} else {
					fmt.Printf("Encrypted %s (%3.1fs)\n", filepath.Base(basename), time.Duration(dur).Seconds())
				}
				wg.Done()

			}(basename)
		}
	}
	wg.Wait()
	fmt.Printf("Complete %3.1fs\n", time.Since(t).Seconds())
}

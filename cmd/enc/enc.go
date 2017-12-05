package main

import (
	"fmt"
	"github.com/devplayg/encdec"
	"github.com/dustin/go-humanize"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func main() {

	// Check arguments
	args := os.Args[1:]
	wg := new(sync.WaitGroup)
	if len(args) < 1 {
		fmt.Println("Encrypt files")
		return
	}

	// Get password hash
	err := encdec.SetSecretKey(2)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Encrypt
	t := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	for _, target := range args {
		files, err := filepath.Glob(target)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Encrypting..")
		for _, f := range files {
			absPath, _ := filepath.Abs(f)
			wg.Add(1)
			go func(f string) {
				newFile, err := encdec.Encrypt(f)
				if err != nil {
					os.Remove(newFile.Name())
					fmt.Println(err.Error())
				} else {
					srcFile, _ := os.Stat(absPath)
					dstFile, _ := os.Stat(newFile.Name())
					fmt.Printf("%s (%s Bytes) => %s (%s Bytes)\n", filepath.Base(f), humanize.Comma(srcFile.Size()), filepath.Base(newFile.Name()), humanize.Comma(dstFile.Size()))
				}
				wg.Done()
			}(absPath)
		}
	}
	wg.Wait()
	fmt.Printf("Complete %3.1fs\n", time.Since(t).Seconds())
}

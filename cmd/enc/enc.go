package main

import (
	"github.com/devplayg/encdec"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
	"github.com/dustin/go-humanize"
)

func main() {

	// Check arguments
	args := os.Args[1:]
	wg := new(sync.WaitGroup)
	if len(args) < 1 {
		log.Println("Encrypt files")
		return
	}

	// Get password hash
	encdec.SetSecretKey(2)

	// Encrypt
	t := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())
	for _, target := range args {
		files, err := filepath.Glob(target)
		if err != nil {
			log.Println(err.Error())
		}

		log.Println("Encrypting..")
		for _, f := range files {
			absPath, _ := filepath.Abs(f)
			wg.Add(1)
			go func(f string) {
				newFile, err := encdec.Encrypt(f)
				if err != nil {
					os.Remove(newFile.Name())
					log.Println("[error]", err.Error())
				} else {
					fi, _ := os.Stat(newFile.Name())
					log.Printf("[%s] %s => %s\n", humanize.Comma(fi.Size()), filepath.Base(f), filepath.Base(newFile.Name()))
				}
				wg.Done()
			}(absPath)
		}
	}
	wg.Wait()
	log.Printf("Complete %3.1fs\n", time.Since(t).Seconds())
}

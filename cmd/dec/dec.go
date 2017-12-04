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
		log.Println("Decrypt files")
		return
	}

	// Get password hash
	encdec.SetSecretKey(1)

	// Decrypt
	t := time.Now()
	runtime.GOMAXPROCS(runtime.NumCPU())

	nameMap := encdec.NewNameMap()

	log.Println("Decrypting..")
	for _, target := range args {
		files, err := filepath.Glob(target)
		if err != nil {
			log.Println(err.Error())
		}
		for _, f := range files {
			absPath, _ := filepath.Abs(f)
			wg.Add(1)
			go func(f string) {
				decFile, originFileName, err := encdec.Decrypt(f)
				if err != nil {
					os.Remove(decFile.Name())
					log.Printf("[error] %s: %s", err.Error(), filepath.Base(f))
				} else {
					newName, err2 := encdec.Rename(decFile, originFileName, nameMap)
					if err2 != nil {
						log.Printf("[error] %s: %s => %s", err2.Error(), filepath.Base(f), newName)
					} else {
						fi, _ := os.Stat(newName)
						log.Printf("[%s] %s => %s \n", humanize.Comma(fi.Size()), filepath.Base(f), newName)
					}
				}
				wg.Done()

			}(absPath)
		}
	}
	wg.Wait()
	log.Printf("Complete %3.1fs\n", time.Since(t).Seconds())
}

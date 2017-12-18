package main

import (
	"github.com/devplayg/encdec"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
	"github.com/dustin/go-humanize"
	"sync/atomic"
	"fmt"
	"log"
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
	var count uint64 = 0
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
					fmt.Printf("%s: %s", err.Error(), filepath.Base(f))
				} else {
					newName, err2 := encdec.Rename(decFile, originFileName, nameMap)
					if err2 != nil {
						fmt.Printf("%s: %s => %s", err2.Error(), filepath.Base(f), newName)
					} else {
						srcFile, _ := os.Stat(f)
						dstFile, _ := os.Stat(newName)
						fmt.Printf("%s (%s Bytes) => %s (%s Bytes)\n", filepath.Base(f), humanize.Comma(srcFile.Size()), newName, humanize.Comma(dstFile.Size()))
					}
				}
				atomic.AddUint64(&count, 1)
				wg.Done()
			}(absPath)
		}
	}
	wg.Wait()
	fmt.Printf("Count: %d, Duration %3.1fs\n", count, time.Since(t).Seconds())
}

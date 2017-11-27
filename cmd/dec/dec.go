package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/devplayg/encdec"
)

func main() {
	args := os.Args[1:]

	for _, target := range args {
		files, err := filepath.Glob(target)
		if err != nil {
			fmt.Errorf(err.Error())
		}
		for _, f := range files {
			fp, _ := filepath.Abs(f)
			tempFile, dur, err := encdec.Decrypt(fp)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			stripedBasename := strings.TrimSuffix(f, filepath.Ext(f))
			newFileName := fp + filepath.Ext(stripedBasename)
			os.Rename(tempFile.Name(), newFileName)
			fmt.Printf("Decrypting %s (%3.1fs)\n", f, time.Duration(dur).Seconds())
		}
	}
}

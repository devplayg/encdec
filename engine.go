package encdec

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/devplayg/golibs/crypto"
	"github.com/howeyc/gopass"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

var PrivateKey []byte
var Version = []byte{1}

func SetSecretKey(count int) {
	fmt.Printf("Password: ")
	password1, err := gopass.GetPasswd()
	if err != nil {
		log.Fatal(err)
	}
	h1 := md5.New()
	h1.Write(password1)

	if count > 1 {
		fmt.Printf("Password confirm: ")
		password2, err := gopass.GetPasswd()
		if err != nil {
			log.Fatal(err)
		}
		h2 := md5.New()
		h2.Write(password2)
		if !bytes.Equal(h1.Sum(nil), h2.Sum(nil)) {
			log.Fatal("Incorrect password")
		}
	}
	PrivateKey = h1.Sum(nil)
}

func Encrypt(fp string) (*os.File, time.Duration, error) {
	t := time.Now()

	// Create temp file
	tempFile, err := ioutil.TempFile(filepath.Dir(fp), "enc_")
	if err != nil {
		return tempFile, time.Since(t), err
	}
	defer tempFile.Close()

	// Write version
	tempFile.Write(Version)

	// Encrypt filename
	encFileName, err := crypto.EncAes256(PrivateKey, []byte(filepath.Base(fp)))
	if err != nil {
		return tempFile, time.Since(t), err
	}

	// Encrypt file name
	nameLen := make([]byte, 2)
	binary.BigEndian.PutUint16(nameLen, uint16(len(encFileName)*2)) // version
	_, err = tempFile.Write(nameLen)
	if err != nil {
		return tempFile, time.Since(t), err
	}
	_, err = tempFile.Write(encFileName)
	if err != nil {
		return tempFile, time.Since(t), err
	}

	// Encrypt data
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return tempFile, time.Since(t), err
	}
	encData, err := crypto.EncAes256(PrivateKey, b)
	if err != nil {
		return tempFile, time.Since(t), err
	}
	_, err = tempFile.Write(encData)
	if err != nil {
		return tempFile, time.Since(t), err
	}

	return tempFile, time.Since(t), nil
}

func Decrypt(fp string) (*os.File, string, time.Duration, error) {
	t := time.Now()

	// Create temp file
	tempFile, err := ioutil.TempFile(filepath.Dir(fp), "dec_")
	if err != nil {
		return nil, "", 0, err
	}
	defer tempFile.Close()

	// Read decrypted file
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, "", 0, err
	}

	// Decrypt file name
	nameLen := binary.BigEndian.Uint16(b[1:3]) / 2
	encFileName := b[3 : nameLen+3]
	originFileName, err := crypto.DecAes256(PrivateKey, encFileName)
	if err != nil {
		return tempFile, "", time.Since(t), err
	}

	// Decrypt data
	decData, err := crypto.DecAes256(PrivateKey, b[nameLen+3:])
	if err != nil {
		return tempFile, "", time.Since(t), err
	}

	_, err = tempFile.Write(decData)
	if err != nil {
		return tempFile, "", time.Since(t), err
	}
	//

	return tempFile, string(originFileName), time.Since(t), nil
}

func Rename(decFile *os.File, originFileName string, nameTable map[string]bool) (string, error) {
	suffix := 0
	var newName string

	for suffix < 10 {
		if suffix > 0 {
			newName = strings.TrimSuffix(string(originFileName), filepath.Ext(string(originFileName)))
			newName = newName + "_" + strconv.Itoa(suffix) + filepath.Ext(string(originFileName))
		} else {
			newName = originFileName
		}
		newFilePath := filepath.Join(filepath.Dir(decFile.Name()), newName)
		if _, err := os.Stat(newFilePath); os.IsNotExist(err) {
			rwMutex := new(sync.RWMutex)
			rwMutex.Lock()
			_, ok := nameTable[newFilePath]

			if ok {
				rwMutex.Unlock()
				continue
			} else {
				nameTable[newFilePath] = true
			}
			rwMutex.Unlock()
			err2 := os.Rename(decFile.Name(), newFilePath)
			rwMutex.Lock()
			delete(nameTable, newFilePath)
			rwMutex.Unlock()
			if err2 == nil {
				return filepath.Base(newFilePath), nil
			}

		}

		suffix += 1
	}

	return filepath.Base(decFile.Name()), errors.New("Failed to rename file")

}

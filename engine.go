package encdec

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/devplayg/golibs/crypto"
	"github.com/howeyc/gopass"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var PrivateKey []byte
var Version = []byte{1}

func SetSecretKey(count int) error {
	fmt.Printf("Password: ")
	password1, err := gopass.GetPasswd()
	if err != nil {
		return err
	}
	if len(password1) < 1 {
		return errors.New("password is too short")
	}
	h1 := sha256.Sum256([]byte(password1))

	if count > 1 {
		fmt.Printf("Password confirm: ")
		password2, err := gopass.GetPasswd()
		if err != nil {
			return err
		}
		h2 := sha256.Sum256([]byte(password2))
		if h1 != h2 {
			return errors.New("incorrect password")
		}
	}
	PrivateKey = h1[:]

	return nil
}

func Encrypt(fp string) (*os.File, error) {
	// Create temp file
	tempFile, err := ioutil.TempFile(filepath.Dir(fp), "enc_")
	if err != nil {
		return tempFile, err
	}
	defer tempFile.Close()

	// Write version
	tempFile.Write(Version)

	// Encrypt filename
	encFileName, err := crypto.EncAes256(PrivateKey, []byte(filepath.Base(fp)))
	if err != nil {
		return tempFile, err
	}

	// Encrypt file name
	nameLen := make([]byte, 2)
	binary.BigEndian.PutUint16(nameLen, uint16(len(encFileName)*2)) // version
	_, err = tempFile.Write(nameLen)
	if err != nil {
		return tempFile, err
	}
	_, err = tempFile.Write(encFileName)
	if err != nil {
		return tempFile, err
	}

	// Encrypt data
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return tempFile, err
	}
	encData, err := crypto.EncAes256(PrivateKey, b)
	if err != nil {
		return tempFile, err
	}
	_, err = tempFile.Write(encData)
	if err != nil {
		return tempFile, err
	}

	return tempFile, nil
}

func Decrypt(fp string) (*os.File, string, error) {
	// Create temp file
	tempFile, err := ioutil.TempFile(filepath.Dir(fp), "dec_")
	if err != nil {
		return nil, "", err
	}
	defer tempFile.Close()

	// Read decrypted file
	b, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, "", err
	}

	// Decrypt file name
	nameLen := binary.BigEndian.Uint16(b[1:3]) / 2
	encFileName := b[3 : nameLen+3]
	originFileName, err := crypto.DecAes256(PrivateKey, encFileName)
	if err != nil {
		return tempFile, "", errors.New("failed to decrypt")
	}

	// Decrypt data
	decData, err := crypto.DecAes256(PrivateKey, b[nameLen+3:])
	if err != nil {
		return tempFile, "", errors.New("failed to decrypt(-2)")
	}

	_, err = tempFile.Write(decData)
	if err != nil {
		return tempFile, "", errors.New("failed to decrypt(-3)")
	}

	return tempFile, string(originFileName), nil
}

func Rename(decFile *os.File, originFileName string, nameMap *NameMap) (string, error) {
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

			_, ok := nameMap.Load(newFilePath)
			if ok {
				continue
			} else {
				nameMap.Store(newFilePath, true)
			}
			err2 := os.Rename(decFile.Name(), newFilePath)
			nameMap.Delete(newFilePath)
			if err2 == nil {
				return filepath.Base(newFilePath), nil
			}

		}

		suffix += 1
	}

	return filepath.Base(decFile.Name()), errors.New("Failed to rename file")
}

type NameMap struct {
	sync.RWMutex
	internal map[string]bool
}

func NewNameMap() *NameMap {
	return &NameMap{
		internal: make(map[string]bool),
	}
}

func (rm *NameMap) Load(key string) (value bool, ok bool) {
	rm.RLock()
	result, ok := rm.internal[key]
	rm.RUnlock()
	return result, ok
}

func (rm *NameMap) Delete(key string) {
	rm.Lock()
	delete(rm.internal, key)
	rm.Unlock()
}

func (rm *NameMap) Store(key string, value bool) {
	rm.Lock()
	rm.internal[key] = value
	rm.Unlock()
}

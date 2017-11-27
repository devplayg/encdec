package encdec

import (
	"io/ioutil"
	"os"
	"path/filepath"
	//	"strings"
	"time"

	"github.com/devplayg/golibs/crypto"
)

var PrivateKey = []byte("_AES256_ENC_KEY_")

func Encrypt(basename string) (*os.File, time.Duration, error) {

	t := time.Now()

	// Create temp file
	tempFile, err := ioutil.TempFile(filepath.Dir(basename), filepath.Base(basename)+".")
	if err != nil {
		return nil, 0, err
	}
	defer tempFile.Close()

	// Read file
	b, err := ioutil.ReadFile(basename)
	if err != nil {
		return nil, 0, err
	}

	// Encrypt data
	enc, err := crypto.EncAes256(PrivateKey, b)
	if err != nil {
		return nil, 0, err
	}

	// Write to file
	_, err = tempFile.Write(enc)
	if err != nil {
		return nil, 0, err
	}
	return tempFile, time.Since(t), nil
}

func Decrypt(basename string) (*os.File, time.Duration, error) {
	t := time.Now()

	// Create temp file
	tempFile, err := ioutil.TempFile(filepath.Dir(basename), "")
	if err != nil {
		return nil, 0, err
	}
	defer tempFile.Close()

	// Read file
	b, err := ioutil.ReadFile(basename)
	if err != nil {
		return nil, 0, err
	}

	// Encrypt data
	dec, _ := crypto.DecAes256(PrivateKey, b)
	if err != nil {
		return nil, 0, err
	}

	// Write to file
	_, err = tempFile.Write(dec)
	if err != nil {
		return nil, 0, err
	}

	return tempFile, time.Since(t), nil

	//	t := time.Now()

	//	stripedBasename := strings.TrimSuffix(basename, filepath.Ext(basename))
	//	newFile := basename + filepath.Ext(stripedBasename)
	//		 ioutil.WriteFile("/tmp/dat1", d1, 0644)

	// Read file

	// Write file

	//	fmt.Println(filepath.Ext(stripedBasename))
	//	fmt.Println(stripedFp)
	// zxcvasdfasdf.text.213412341234

	// Create temp file
	//	tempFile, err := ioutil.TempFile(basename + filepath.Ext(stripedBasename))
	//	ioutil.
	//	if err != nil {
	//		return nil, 0, err
	//	}
	//	defer tempFile.Close()

	//	// Read file
	//	b, err := ioutil.ReadFile(fp)
	//	if err != nil {
	//		return nil, 0, err
	//	}

	//	// Encrypt data
	//	enc, _ := crypto.EncAes256(PrivateKey, b)
	//	if err != nil {
	//		return nil, 0, err
	//	}

	//	// Write to file
	//	_, err = tempFile.Write(enc)
	//	if err != nil {
	//		return nil, 0, err
	//	}
	//	return nil, 0, nil
}

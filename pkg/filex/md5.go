package filex

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

//export
func CalculateReaderMd5(reader io.Reader) string {
	r := bufio.NewReader(reader)
	h := md5.New()
	_, err := io.Copy(h, r)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}

//export
func CalculateBytesMd5(data []byte) string {
	mh := md5.New()
	if _, terr := mh.Write(data); terr != nil {
		fmt.Printf("error found during calculate bytes md5 : %s", terr.Error())
		return ""
	}
	return hex.EncodeToString(mh.Sum(nil))
}

//export
func CalculateFileMd5(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer func() {
		if terr := f.Close(); terr != nil {
			fmt.Println(err)
		}
	}()
	r := bufio.NewReader(f)

	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}

	//return fmt.Sprintf("%x", h.Sum(nil))
	return hex.EncodeToString(h.Sum(nil))
}

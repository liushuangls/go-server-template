package xfile

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// GetFileMIME return file mime type: https://mimesniff.spec.whatwg.org/
func GetFileMIME(file io.ReadSeeker) (string, error) {
	mime, err := GetFileMIMEType(file)
	if err != nil {
		return "", err
	}
	return mime.String(), nil
}

func GetFileMIMEType(file io.ReadSeeker) (*mimetype.MIME, error) {
	defer file.Seek(0, 0)

	buffer := make([]byte, 1024*10)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return mimetype.Detect(buffer[:n]), nil
}

func FileIsImage(file io.ReadSeeker) (bool, error) {
	mime, err := GetFileMIME(file)
	if err != nil {
		return false, err
	}
	return MIMEIsImage(mime), nil
}

func MIMEIsImage(mime string) bool {
	return strings.HasPrefix(mime, "image/")
}

func CalcFileHexMD5(file io.ReadSeeker) (string, error) {
	defer file.Seek(0, 0)

	hash := md5.New()
	_, err := io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

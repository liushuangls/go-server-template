package xfile

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"mime/multipart"
	"strings"
)

func IsImageFile(fh *multipart.FileHeader) (bool, string, error) {
	return IsImageFileByName(fh.Filename)
}

func IsImageFileByName(name string) (bool, string, error) {
	mime := TypeByExtension(name)
	if strings.HasPrefix(mime, "image/") {
		return true, mime, nil
	}
	return false, "", nil
}

func CalcFileHexMD5(file io.ReadSeeker) (string, error) {
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", err
	}
	hash := md5.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

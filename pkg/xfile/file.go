package xfile

import (
	"io"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// GetFileMIME return file mime type: https://mimesniff.spec.whatwg.org/
func GetFileMIME(file io.ReadSeeker) (string, error) {
	// 确保从开头读取
	_, err := file.Seek(0, 0)
	if err != nil {
		return "", err
	}

	// 读取前512字节
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// 读取完成后重置数据
	defer func() {
		_, _ = file.Seek(0, 0)
	}()

	return mimetype.Detect(buffer[:n]).String(), nil
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

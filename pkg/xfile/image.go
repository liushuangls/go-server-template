package xfile

import (
	"bytes"
	"encoding/base64"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"strings"
)

func GetImageFileMIMEForFileHeader(header *multipart.FileHeader) (mime string, err error) {
	imgFile, err := header.Open()
	defer imgFile.Close()
	if err != nil {
		return ``, err
	}
	mime, err = GetImageFileMIME(imgFile)
	if err != nil {
		return ``, err
	}

	return mime, nil
}

func GetImageDimensionForFileHeader(header *multipart.FileHeader) (width int, height int, err error) {
	imgFile, err := header.Open()
	defer imgFile.Close()
	return getImageDimensionForReader(imgFile)
}

func GetImageDimensionForByteData(imgData []byte) (width int, height int, err error) {
	reader := bytes.NewReader(imgData)
	return getImageDimensionForReader(reader)
}

func GetImageDimensionForBase64(base64Data string) (width int, height int, err error) {
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64Data))
	return getImageDimensionForReader(reader)
}

func getImageDimensionForReader(reader io.Reader) (width int, height int, err error) {
	im, _, err := image.DecodeConfig(reader)
	if err != nil {
		return 0, 0, err

	}

	return im.Width, im.Height, nil
}

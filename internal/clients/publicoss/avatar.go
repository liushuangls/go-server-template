package publicoss

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/liushuangls/go-server-template/configs"
	"github.com/liushuangls/go-server-template/pkg/xoss"
)

var downloadFileFailedErr = errors.New("download file failed")
var uploadFileFailedErr = errors.New("upload file failed")

type avatarStoreInfo struct {
	storeKey  string
	storeLink string
}

type Avatar struct {
	conf *xoss.Config
	oss  xoss.Oss
}

func NewAvatar(conf *configs.Config) *Avatar {
	return &Avatar{
		oss:  xoss.NewOss(&conf.PublicOSS),
		conf: &conf.PublicOSS,
	}
}

func randomRename(oldName string, defaultExt string) string {
	ext := filepath.Ext(oldName)
	if ext == "" {
		ext = "." + defaultExt
	}
	return fmt.Sprintf("avatar/%s%s", uuid.New(), ext)
}

func (a *Avatar) genFileNameAndOssLink(name string) *avatarStoreInfo {
	splits := strings.Split(name, "?") //去掉url的query
	storeKey := randomRename(splits[0], "png")
	return &avatarStoreInfo{
		storeKey:  storeKey,
		storeLink: a.conf.GetHost() + "/" + storeKey,
	}
}

func (a *Avatar) UploadAvatar(ctx context.Context, file io.Reader, fileName, mimeType string) (string, error) {
	storeKey := randomRename(fileName, "png")
	storeLink := a.conf.GetHost() + "/" + storeKey
	if err := a.oss.PutObject(ctx, storeKey, file, mimeType); err != nil {
		return "", err
	}
	return storeLink, nil
}

func (a *Avatar) genAvatarByURL(url string, info *avatarStoreInfo) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*8)
	defer cancel()

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp == nil {
		return "", errors.Join(downloadFileFailedErr, fmt.Errorf("genAvatarByURL.geterr err:%s", err))
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.Join(downloadFileFailedErr, fmt.Errorf("genAvatarByURL.geterr statuscode:%d", resp.StatusCode))
	}
	defer resp.Body.Close()

	err = a.oss.PutObject(ctx, info.storeKey, resp.Body, mime.TypeByExtension(info.storeKey))
	if err != nil {
		return "", errors.Join(uploadFileFailedErr, fmt.Errorf("genAvatarByURL.PutObject err:%s", err.Error()))
	}
	return info.storeLink, nil
}

func (a *Avatar) genAvatarByRandom() string {
	return fmt.Sprintf("%s/avatar/default-avatar-%d.png", a.conf.GetHost(), rand.Intn(16))
}

func (a *Avatar) Handle(avatar string) string {
	if avatar != "" {
		link, err := a.genAvatarByURL(avatar, a.genFileNameAndOssLink(avatar))
		if err == nil {
			return link
		}
	}
	return a.genAvatarByRandom()
}

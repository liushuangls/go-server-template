package xoss

import (
	"context"
	"io"
	"time"
)

type Config struct {
	Platform string `yaml:"Platform"`
	Aws      S3Conf `yaml:"AWS"`
}

func (c *Config) GetBucket() string {
	switch c.Platform {
	case "aws":
		return c.Aws.Bucket
	}
	return ""
}

func (c *Config) GetHost() string {
	switch c.Platform {
	case "aws":
		return c.Aws.Host
	}
	return ""
}

type Oss interface {
	// PutObject 上传文件
	PutObject(ctx context.Context, origin string, reader io.Reader, contentType string) error
	// PreSignURL 预先签名授权访问url
	PreSignURL(key string, ttl time.Duration) (string, error)
	// PreSignDownloadURL 预先签名授权下载url
	PreSignDownloadURL(key string, ttl time.Duration, fileName string) (string, error)
	// SignPolicy 签署上传url和format
	SignPolicy(key string, ttl time.Duration, size int64, contentType string) (*PolicyToken, error)
	// DownloadObject 下载文件
	DownloadObject(remotePath, localPath string) error
	// DeleteObject 删除文件
	DeleteObject(key string) error
	// DeleteObjects 批量删除文件
	DeleteObjects(keys []string) error
	//Credentials 临时凭证
	Credentials(key string) (*Credentials, error)
	//CompleteMultipartUpload 结束分片上传
	CompleteMultipartUpload(uploadID string, key string) error
}

func NewOss(conf *Config) Oss {
	switch conf.Platform {
	case "aws":
		return newAwsS3Oss(&conf.Aws)
	default:
		return newAwsS3Oss(&conf.Aws)
	}
}

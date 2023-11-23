package xoss

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const _s3Policy = `{"Version":"2012-10-17","Statement":[{"Sid":"Stmt1","Effect":"Allow","Action":"s3:*","Resource":"arn:aws:s3:::%s/%s"}]}`

func newAwsS3Oss(conf *S3Conf) Oss {
	awsS3 := new(S3Oss)
	sess, err := newSession(conf)
	if err != nil {
		panic(err)
	}
	awsS3.session = sess
	awsS3.client = newClient(sess)
	awsS3.uploader = newUploader(awsS3.client)
	awsS3.downloader = newDownloader(awsS3.client)
	awsS3.conf = conf
	return awsS3
}

func newSession(conf *S3Conf) (*session.Session, error) {
	err := os.Setenv("AWS_ACCESS_KEY", conf.AccessKeyID)
	if err != nil {
		return nil, err
	}
	err = os.Setenv("AWS_SECRET_KEY", conf.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return session.NewSession(&aws.Config{
		Region:          aws.String(conf.Region),
		S3UseAccelerate: aws.Bool(conf.S3UseAccelerate),
	})
}
func newClient(sess *session.Session) *s3.S3 {
	return s3.New(sess, &aws.Config{
		DisableRestProtocolURICleaning: aws.Bool(true),
	})
}
func newUploader(svc s3iface.S3API) *s3manager.Uploader {
	return s3manager.NewUploaderWithClient(svc)
}
func newDownloader(svc s3iface.S3API) *s3manager.Downloader {
	return s3manager.NewDownloaderWithClient(svc, func(d *s3manager.Downloader) {
	})
}

type S3Conf struct {
	Host            string `yaml:"Host"`
	Bucket          string `yaml:"Bucket"`
	Region          string `yaml:"Region"`
	AccessKeyID     string `yaml:"AccessKeyID"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	S3UseAccelerate bool   `yaml:"S3UseAccelerate"`
}
type S3Oss struct {
	session    *session.Session
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
	conf       *S3Conf
}

func (o *S3Oss) PutObject(ctx context.Context, origin string, reader io.Reader, contentType string) error {
	input := &s3manager.UploadInput{
		Bucket:      aws.String(o.conf.Bucket),
		Key:         &origin,
		Body:        reader,
		ContentType: &contentType,
	}
	output, err := o.uploader.UploadWithContext(ctx, input)
	if err != nil {
		marshal, _ := json.Marshal(output)
		return errors.Join(err, errors.New(string(marshal)))
	}
	return nil
}

func (o *S3Oss) PreSignURL(key string, ttl time.Duration) (string, error) {
	req, _ := o.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(o.conf.Bucket),
		Key:    aws.String(key),
	})
	return req.Presign(ttl)
}
func (o *S3Oss) PreSignDownloadURL(key string, ttl time.Duration, fileName string) (string, error) {
	cd := fmt.Sprintf("attachment; filename=\"%s\"", url.PathEscape(fileName))
	req, _ := o.client.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(o.conf.Bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(cd),
	})
	return req.Presign(ttl)
}
func (o *S3Oss) SignPolicy(key string, ttl time.Duration, size int64, contentType string) (*PolicyToken, error) {
	var policyConf = map[string]interface{}{
		"expiration": time.Now().Add(ttl).Format(time.RFC3339),
		"conditions": [][]interface{}{
			{"eq", "$bucket", o.conf.Bucket},
			{"eq", "$key", key},
			{"eq", "$acl", "private"},
			{"eq", "$Content-Type", contentType},
			{"content-length-range", 0, size},
		},
	}
	marshal, err := json.Marshal(policyConf)
	if err != nil {
		return nil, err
	}
	policyConfStr := base64.StdEncoding.EncodeToString(marshal)
	h := hmac.New(sha1.New, []byte(o.conf.AccessKeySecret))
	h.Write([]byte(policyConfStr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	formData := map[string]string{
		"AWSAccessKeyId": o.conf.AccessKeyID,
		"acl":            "private",
		"Content-Type":   contentType,
		"policy":         policyConfStr,
		"signature":      signature,
		"key":            key,
		"bucket":         o.conf.Bucket,
	}
	return &PolicyToken{
		Host:     o.conf.Host,
		FormData: formData,
	}, nil
}

func (o *S3Oss) CopyObject(srcBucket, srcKey string, targetBucket, targetKey string) error {
	upload, err := o.client.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket: aws.String(targetBucket),
		Key:    aws.String(targetKey),
	})
	if err != nil {
		return err
	}
	output, err := o.client.UploadPartCopy(&s3.UploadPartCopyInput{
		Bucket:     aws.String(targetBucket),
		CopySource: aws.String(url.PathEscape(srcBucket + "/" + srcKey)),
		PartNumber: aws.Int64(1000),
		Key:        aws.String(targetKey),
		UploadId:   upload.UploadId,
	})
	if err != nil {
		marshal, _ := json.Marshal(output)
		return errors.Join(err, errors.New(string(marshal)))
	}

	var parts []*s3.CompletedPart
	result, err := o.client.ListParts(&s3.ListPartsInput{
		Bucket:   aws.String(targetBucket),
		Key:      aws.String(targetKey),
		UploadId: upload.UploadId,
	})
	if err != nil {
		return err
	}
	for _, p := range result.Parts {
		parts = append(parts, &s3.CompletedPart{
			ETag:       p.ETag,
			PartNumber: p.PartNumber,
		})
	}
	_, err = o.client.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(targetBucket),
		Key:      aws.String(targetKey),
		UploadId: upload.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: parts,
		},
	})
	return err
}

func (o *S3Oss) DownloadObject(remotePath, localPath string) error {
	fd, err := os.OpenFile(localPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0664)
	if err != nil {
		return err
	}
	defer fd.Close()
	_, err = o.downloader.Download(fd, &s3.GetObjectInput{
		Bucket: aws.String(o.conf.Bucket),
		Key:    aws.String(remotePath),
	})
	return err
}

func (o *S3Oss) DeleteObject(key string) error {
	_, err := o.client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(o.conf.Bucket),
		Key:    aws.String(key),
	})
	return err
}

func (o *S3Oss) DeleteObjects(keys []string) error {
	objs := make([]*s3.ObjectIdentifier, 0, len(keys))
	for _, key := range keys {
		objs = append(objs, &s3.ObjectIdentifier{
			Key: aws.String(key),
		})
	}
	_, err := o.client.DeleteObjects(&s3.DeleteObjectsInput{
		Bucket: aws.String(o.conf.Bucket),
		Delete: &s3.Delete{
			Objects: objs,
			Quiet:   aws.Bool(false),
		},
	})
	return err
}

func (o *S3Oss) Credentials(key string) (*Credentials, error) {
	expireAt := time.Now().Add(stscreds.DefaultDuration)
	creds := stscreds.NewCredentials(o.session, "arn:aws:iam::314459260709:role/AWS-ServiceRoleForS3", func(provider *stscreds.AssumeRoleProvider) {
		provider.Policy = aws.String(fmt.Sprintf(_s3Policy, o.conf.Bucket, key))
	})
	val, err := creds.Get()
	if err != nil {
		return nil, err
	}

	return &Credentials{
		AccessKeyID:     val.AccessKeyID,
		SecretAccessKey: val.SecretAccessKey,
		SessionToken:    val.SessionToken,
		Bucket:          o.conf.Bucket,
		ExpireAt:        expireAt.Unix(),
		Region:          o.conf.Region,
	}, nil
}

func (o *S3Oss) CompleteMultipartUpload(uploadID string, key string) error {
	var parts []*s3.CompletedPart
	result, err := o.client.ListParts(&s3.ListPartsInput{
		Bucket:   aws.String(o.conf.Bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
	})
	if err != nil {
		return err
	}
	for _, p := range result.Parts {
		parts = append(parts, &s3.CompletedPart{
			ETag:       p.ETag,
			PartNumber: p.PartNumber,
		})
	}
	_, err = o.client.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(o.conf.Bucket),
		Key:      aws.String(key),
		UploadId: aws.String(uploadID),
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: parts,
		},
	})
	return err
}

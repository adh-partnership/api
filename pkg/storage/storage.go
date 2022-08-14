package storage

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/kzdv/api/pkg/config"
)

type Client struct {
	Session *session.Session
	S3      *s3.S3
	Bucket  string
}

var client map[string]*Client

func Storage(name string) *Client {
	return client[name]
}

func Configure(c config.ConfigStorage, name string, bucket string) *Client {
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
		Region:      aws.String(c.Region),
		Endpoint:    aws.String(c.Endpoint),
	}

	cl := &Client{}
	cl.Session = session.Must(session.NewSession(s3Config))
	cl.S3 = s3.New(cl.Session)
	cl.Bucket = bucket

	client[name] = cl

	return cl
}

func (c *Client) PutObject(key string, filepath string, private bool, length int64, contenttype string) error {
	acl := aws.String("public-read")
	if private {
		acl = aws.String("private")
	}

	uploader := s3manager.NewUploader(c.Session, func(u *s3manager.Uploader) {
		u.PartSize = 10 * 1024 * 1024
	})

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = f.Close()
	}()

	_, err = uploader.Upload(&s3manager.UploadInput{
		ACL:         acl,
		Body:        f,
		Bucket:      aws.String(c.Bucket),
		ContentType: aws.String(contenttype),
		Key:         aws.String(key),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ListObjects() ([]string, error) {
	var objects []string
	err := c.S3.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(c.Bucket),
	}, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		for _, o := range page.Contents {
			objects = append(objects, *o.Key)
		}
		return true
	})

	return objects, err
}

func (c *Client) DeleteObject(key string) error {
	_, err := c.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(c.Bucket),
		Key:    aws.String(key),
	})

	return err
}

/*
 * Copyright ADH Partnership
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package storage

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/adh-partnership/api/pkg/config"
)

type Client struct {
	Session *session.Session
	S3      *s3.S3
	Bucket  string
}

var client map[string]*Client

func init() {
	client = make(map[string]*Client)
}

func Storage(name string) *Client {
	return client[name]
}

func Configure(c config.ConfigStorage, name string) (*Client, error) {
	if c.AccessKey == "" || c.SecretKey == "" || c.Bucket == "" {
		return nil, fmt.Errorf("storage configuration not configured")
	}
	s3Config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, ""),
		Region:      aws.String(c.Region),
		Endpoint:    aws.String(c.Endpoint),
	}

	cl := &Client{}
	cl.Session = session.Must(session.NewSession(s3Config))
	cl.S3 = s3.New(cl.Session)
	cl.Bucket = c.Bucket

	client[name] = cl

	return cl, nil
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

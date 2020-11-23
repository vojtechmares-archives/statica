package aws

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	amazon_aws "github.com/aws/aws-sdk-go/aws"
	amazon_credentials "github.com/aws/aws-sdk-go/aws/credentials"
	amazon_session "github.com/aws/aws-sdk-go/aws/session"
	amazon_s3 "github.com/aws/aws-sdk-go/service/s3"
)

type AWS struct {
	logger  *log.Logger
	session *amazon_session.Session
	Region  string
}

type AWSCredentials struct {
	AccessKeyID string
	SecretKey   string
}

func NewAWS(l *log.Logger, c *AWSCredentials, region string) *AWS {
	session, err := amazon_session.NewSession(
		&amazon_aws.Config{
			Region: amazon_aws.String(region),
			Credentials: amazon_credentials.NewStaticCredentials(
				c.AccessKeyID,
				c.SecretKey,
				"",
			),
		},
	)
	if err != nil {
		l.Fatalf("Unable to create AWS Session\n%v", err)
	}

	return &AWS{
		session: session,
		logger:  l,
		Region:  region,
	}
}

func (aws *AWS) UploadFilesFromDir(b *Bucket, dir string) {
	var files []string

	aws.logger.Println("Uploading files...")

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		files = append(files, path)
		return nil
	})
	if err != nil {
		aws.logger.Fatalf("Error while reading upload dir\n%v", err)
	}

	svc := aws.getS3()

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			aws.logger.Fatalf("Could not open file\n%v", err)
		}
		defer f.Close()

		key := strings.Replace(file, fmt.Sprintf("%s/", dir), "", 1)

		aws.logger.Printf("Uploading %s to %s\n", file, key)

		_, err = svc.PutObject(&amazon_s3.PutObjectInput{
			Bucket:      amazon_aws.String(b.Name),
			Key:         amazon_aws.String(key),
			Body:        f,
			ContentType: amazon_aws.String(mime.TypeByExtension(filepath.Ext(file))),
		})
		if err != nil {
			aws.logger.Fatalf("Error while uploading file\n%v", err)
		}
	}
}

func (aws *AWS) getS3() *amazon_s3.S3 {
	return amazon_s3.New(aws.session)
}

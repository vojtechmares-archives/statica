package aws

import (
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"

	amazon_aws "github.com/aws/aws-sdk-go/aws"
	amazon_credentials "github.com/aws/aws-sdk-go/aws/credentials"
	amazon_session "github.com/aws/aws-sdk-go/aws/session"
	amazon_s3 "github.com/aws/aws-sdk-go/service/s3"
	amazon_s3manager "github.com/aws/aws-sdk-go/service/s3/s3manager"
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

type directoryIterator struct {
	filePaths []string
	bucket    string
	next      struct {
		path string
		f    *os.File
	}
	err error
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
	svc := aws.getS3()
	iter := newDirectoryIterator(b.Name, dir)
	uploader := amazon_s3manager.NewUploaderWithClient(svc)

	aws.logger.Println("Uploading files...")

	if err := uploader.UploadWithIterator(amazon_aws.BackgroundContext(), iter); err != nil {
		aws.logger.Fatalf("Error while uploading files to bucket\n%v", err)
	}

	aws.logger.Println("Upload complete")
}

func (aws *AWS) getS3() *amazon_s3.S3 {
	return amazon_s3.New(aws.session)
}

func newDirectoryIterator(bucket, dir string) amazon_s3manager.BatchUploadIterator {
	paths := []string{}
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// We care only about files, not directories
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	return &directoryIterator{
		filePaths: paths,
		bucket:    bucket,
	}
}

func (iter *directoryIterator) UploadObject() amazon_s3manager.BatchUploadObject {
	f := iter.next.f
	key := iter.next.path
	contentType := mime.TypeByExtension(filepath.Ext(iter.next.path))

	if i := strings.Index(iter.next.path, "/"); i != -1 {
		key = iter.next.path[i:]
	}
	return amazon_s3manager.BatchUploadObject{
		Object: &amazon_s3manager.UploadInput{
			Bucket:      &iter.bucket,
			Key:         &key,
			Body:        f,
			ContentType: &contentType,
		},

		After: func() error {
			return f.Close()
		},
	}
}

func (iter *directoryIterator) Next() bool {
	if len(iter.filePaths) == 0 {
		iter.next.f = nil
		return false
	}

	f, err := os.Open(iter.filePaths[0])
	iter.err = err

	iter.next.f = f
	iter.next.path = iter.filePaths[0]

	iter.filePaths = iter.filePaths[1:]
	return true && iter.Err() == nil
}

func (iter *directoryIterator) Err() error {
	return iter.err
}

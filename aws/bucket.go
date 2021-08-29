package aws

import (
	"fmt"

	amazon_aws "github.com/aws/aws-sdk-go/aws"
	amazon_s3 "github.com/aws/aws-sdk-go/service/s3"
)

type Bucket struct {
	Name   string
	Region string
}

func (aws *AWS) NewBucket(bucketName, region string) *Bucket {
	svc := aws.getS3()
	aws.logger.Println("Looking for existing bucket...")
	blo, err := svc.ListBuckets(&amazon_s3.ListBucketsInput{})
	if err != nil {
		aws.logger.Fatalf("Unable to get list of buckets\n%v\n", err)
	}

	exists := false
	for _, b := range blo.Buckets {
		if *b.Name == bucketName {
			aws.logger.Println("Bucket found.")
			exists = true
			break
		}
	}

	if !exists {
		return aws.createBucket(bucketName, region)
	}

	// Tag Bucket
	aws.tagBucket(bucketName, map[string]string{
		"statica.created-by": "statica",
	})

	return &Bucket{
		Name:   bucketName,
		Region: region,
	}
}

func (b *Bucket) GetHTTPEndpoint() string {
	return fmt.Sprintf("%s.s3-website.%s.amazonaws.com", b.Name, b.Region)
}

func (aws *AWS) createBucket(bucketName, region string) *Bucket {
	svc := aws.getS3()

	var err error

	aws.logger.Println("Bucket does not exists yet. Creating...")
	_, err = svc.CreateBucket(&amazon_s3.CreateBucketInput{Bucket: amazon_aws.String(bucketName)})
	if err != nil {
		aws.logger.Fatalf("Unable to create a bucket\n%v\n", err)
	}

	err = svc.WaitUntilBucketExists(&amazon_s3.HeadBucketInput{Bucket: amazon_aws.String(bucketName)})
	if err != nil {
		aws.logger.Fatalf("Error occurred while waiting for bucket creation\n%v\n", err)
	}
	aws.logger.Println("Bucket created")

	// Bucket configuration (static website hosting, perimissions, etc.)

	aws.logger.Println("Configuring bucket...")
	_, err = svc.PutBucketWebsite(&amazon_s3.PutBucketWebsiteInput{
		Bucket: amazon_aws.String(bucketName),
		WebsiteConfiguration: &amazon_s3.WebsiteConfiguration{
			IndexDocument: &amazon_s3.IndexDocument{Suffix: amazon_aws.String("index.html")},
			ErrorDocument: &amazon_s3.ErrorDocument{Key: amazon_aws.String("error.html")},
		},
	})
	if err != nil {
		aws.logger.Fatalf("Error occurred while configuring bucket (website)\n%v\n", err)
	}

	_, err = svc.PutBucketPolicy(&amazon_s3.PutBucketPolicyInput{
		Bucket: amazon_aws.String(bucketName),
		Policy: amazon_aws.String(
			fmt.Sprintf(`{
						"Version": "2012-10-17",
						"Statement": [
							{
								"Sid": "PublicReadGetObject",
								"Effect": "Allow",
								"Principal": "*",
								"Action": [
									"s3:GetObject"
								],
								"Resource": [
									"arn:aws:s3:::%s/*"
								]
							}
						]
					}`,
				bucketName,
			),
		),
	})
	if err != nil {
		aws.logger.Fatalf("Error occurred while configuring bucket (policy)\n%v\n", err)
	}

	aws.logger.Println("Bucket configured")

	return &Bucket{
		Name:   bucketName,
		Region: region,
	}
}

func (aws *AWS) tagBucket(bucketName string, tags map[string]string) {
	svc := aws.getS3()

	var err error

	aws.logger.Println("Tagging bucket...")

	var awsTags []*amazon_s3.Tag

	for key, val := range tags {
		awsTags = append(awsTags, &amazon_s3.Tag{
			Key:   amazon_aws.String(key),
			Value: amazon_aws.String(val),
		})
	}

	putInput := &amazon_s3.PutBucketTaggingInput{
		Bucket: amazon_aws.String(bucketName),
		Tagging: &amazon_s3.Tagging{
			TagSet: awsTags,
		},
	}

	_, err = svc.PutBucketTagging(putInput)
	if err != nil {
		aws.logger.Fatalf("Error occurred while tagging the bucket\n%v\n", err)
	}

	aws.logger.Println("Bucket tagged")
}

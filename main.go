package main

import (
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
	"github.com/vojtechmares/statica/aws"
	"github.com/vojtechmares/statica/cloudflare"
)

var (
	version string = ""
	commit  string = ""
)

func buildBucketName(prefix, name, suffix string) string {
	return fmt.Sprintf("%s%s%s", prefix, name, suffix)
}

func main() {
	l := log.New(os.Stdout, "[statica] ", log.LstdFlags)

	var bucketName string
	flag.StringVar(&bucketName, "bucket-name", "", "Overrides bucket name")
	var bucketNamePrefix string
	flag.StringVar(&bucketNamePrefix, "bucket-prefix", "", "Bucket name prefix (empty by default)")
	var bucketNameSuffix string
	flag.StringVar(&bucketNameSuffix, "bucket-suffix", "", "Bucket name prefix (empty by default)")
	var v bool
	flag.BoolVar(&v, "version", false, "Statica version")
	flag.Parse()

	if flag.Lookup("version") != nil && v {
		fmt.Printf("Statica version: %s\n", version)
		os.Exit(0)
	}

	accessKeyID, exists := os.LookupEnv("STATICA_AWS_ACCESS_KEY_ID")
	if !exists {
		l.Fatalln("Missing ENV variable STATICA_AWS_ACCESS_KEY_ID")
	}
	secretKey, exists := os.LookupEnv("STATICA_AWS_SECRET_KEY")
	if !exists {
		l.Fatalln("Missing ENV variable STATICA_AWS_SECRET_KEY")
	}
	region, exists := os.LookupEnv("STATICA_AWS_REGION")
	if !exists {
		l.Fatalln("Missing ENV variable STATICA_AWS_REGION")
	}
	apiToken, exists := os.LookupEnv("STATICA_CF_API_TOKEN")
	if !exists {
		l.Fatalln("Missing ENV variable STATICA_CF_API_TOKEN")
	}

	if len(os.Args) < 2 {
		l.Fatalln("You must specify domain or use --version flag")
	}

	domainArg := os.Args[1]

	if bucketName == "" {
		bucketName = domainArg
	}

	uploadDir := "."
	if len(os.Args) == 3 {
		uploadDir = os.Args[2]
	}

	l.Println("Deploying...")

	sa := aws.NewAWS(
		l,
		&aws.AWSCredentials{
			AccessKeyID: accessKeyID,
			SecretKey:   secretKey,
		}, region,
	)

	fullBucketName := buildBucketName(bucketNamePrefix, bucketName, bucketNameSuffix)

	b := sa.NewBucket(fullBucketName, sa.Region)

	sa.UploadFilesFromDir(b, uploadDir)

	cf := cloudflare.NewCloudflareWithAPIToken(l, apiToken)

	cf.ConfigureDomain(domainArg, b.GetHTTPEndpoint())

	l.Println("Deploy completed")
}

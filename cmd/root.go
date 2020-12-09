package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/vojtechmares/statica/aws"
	"github.com/vojtechmares/statica/cloudflare"
)

var (
	accessKeyID string
	secretKey   string
	region      string
	apiToken    string

	bucketName       string
	bucketNamePrefix string
	bucketNameSuffix string
	noDns            bool

	domain    string
	uploadDir string = "."
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVar(&bucketName, "bucket-name", "", "Overrides bucket name")
	rootCmd.Flags().StringVar(&bucketNamePrefix, "bucket-prefix", "", "Bucket name prefix (empty by default)")
	rootCmd.Flags().StringVar(&bucketNameSuffix, "bucket-suffix", "", "Bucket name prefix (empty by default)")
	rootCmd.Flags().BoolVar(&noDns, "no-dns", false, "Omits creation of DNS record")
}

func initConfig() {
	viper.SetEnvPrefix("STATICA")
	viper.BindEnv("AWS_ACCESS_KEY_ID")
	viper.BindEnv("AWS_SECRET_KEY")
	viper.BindEnv("AWS_REGION")
	viper.BindEnv("CF_API_TOKEN")

	accessKeyID = viper.GetString("AWS_ACCESS_KEY_ID")
	secretKey = viper.GetString("AWS_SECRET_KEY")
	region = viper.GetString("AWS_REGION")
	apiToken = viper.GetString("CF_API_TOKEN")
}

func areConfVarsSet() bool {
	missingConfVar := false
	if !viper.IsSet("AWS_ACCESS_KEY_ID") {
		fmt.Println("Missing required environment variable STATICA_AWS_ACCESS_KEY_ID")
		missingConfVar = true
	}

	if !viper.IsSet("AWS_SECRET_KEY") {
		fmt.Println("Missing required environment variable STATICA_AWS_SECRET_KEY")
		missingConfVar = true
	}

	if !viper.IsSet("AWS_REGION") {
		fmt.Println("Missing required environment variable STATICA_AWS_REGION")
		missingConfVar = true
	}

	if !viper.IsSet("CF_API_TOKEN") {
		fmt.Println("Missing required environment variable STATICA_CF_API_TOKEN")

		// If --no-dns flag is provided, we do not need Cloudflare and therefore nor the env var
		if !noDns {
			missingConfVar = true
		}
	}

	return missingConfVar
}

var rootCmd = &cobra.Command{
	Use:   "statica <domain>",
	Short: "Deploys static content",
	Long: `Deploys static content (website) to AWS S3.
		Configures the S3 bucket for static web hosting.
		Points a Cloudflare DNS CNAME record to S3's website endpoint.
	`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(c *cobra.Command, args []string) {
		if areConfVarsSet() {
			os.Exit(1)
		}

		l := log.New(os.Stdout, "[statica] ", log.LstdFlags)

		domain = args[0]

		if len(args) == 2 {
			uploadDir = args[1]
		}

		if bucketName == "" {
			bucketName = domain
		}

		l.Println("Deploying...")

		sa := aws.NewAWS(
			l,
			&aws.AWSCredentials{
				AccessKeyID: accessKeyID,
				SecretKey:   secretKey,
			}, region,
		)

		fullBucketName := fmt.Sprintf("%s%s%s", bucketNamePrefix, bucketName, bucketNameSuffix)

		b := sa.NewBucket(fullBucketName, sa.Region)

		sa.UploadFilesFromDir(b, uploadDir)

		if !noDns {
			cf := cloudflare.NewCloudflareWithAPIToken(l, apiToken)

			cf.ConfigureDomain(domain, b.GetHTTPEndpoint())
		} else {
			l.Printf("Bucket HTTP endpoint http://%s\n", b.GetHTTPEndpoint())
		}

		l.Println("Deploy completed")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

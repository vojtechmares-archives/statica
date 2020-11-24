package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/vojtechmares/statica/aws"
	"github.com/vojtechmares/statica/cloudflare"
)

func init() {
	cobra.OnInitialize(initConfig)
	destroyCmd.Flags().StringVar(&bucketName, "bucket-name", "", "Overrides bucket name")
	destroyCmd.Flags().StringVar(&bucketNamePrefix, "bucket-prefix", "", "Bucket name prefix (empty by default)")
	destroyCmd.Flags().StringVar(&bucketNameSuffix, "bucket-suffix", "", "Bucket name prefix (empty by default)")
	rootCmd.AddCommand(destroyCmd)
}

var destroyCmd = &cobra.Command{
	Use:   "destroy <domain>",
	Short: "Destroys resources (S3 bucket and DNS record)",
	Long: `Destroys everything related for this domain.
		- AWS S3 bucket
		- Cloudflare DNS record
	`,
	Args: cobra.ExactArgs(1),
	Run: func(c *cobra.Command, args []string) {
		if areConfVarsSet() {
			os.Exit(1)
		}

		l := log.New(os.Stdout, "[statica] ", log.LstdFlags)

		domain = args[0]

		if bucketName == "" {
			bucketName = domain
		}

		l.Println("Destroying...")

		a := aws.NewAWS(
			l,
			&aws.AWSCredentials{
				AccessKeyID: accessKeyID,
				SecretKey:   secretKey,
			}, region,
		)

		fullBucketName := fmt.Sprintf("%s%s%s", bucketNamePrefix, bucketName, bucketNameSuffix)

		a.DestroyBucket(fullBucketName)

		cf := cloudflare.NewCloudflareWithAPIToken(l, apiToken)

		cf.DestroyDomain(domain)

		l.Println("Destroy completed")
	},
}

package s3client

import (
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client
var PresignClient *s3.PresignClient

func InitS3Client() *s3.Client {
	S3Client = s3.NewFromConfig(aws.Config{
		Region:      "us-east-1",
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), "")),
	})
	PresignClient = s3.NewPresignClient(S3Client)
	return S3Client
}

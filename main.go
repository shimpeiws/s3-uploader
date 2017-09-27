package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"gopkg.in/yaml.v2"
)

func main() {
	buf, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}

	var fileTransfer FileTransferToS3
	err = yaml.Unmarshal(buf, &fileTransfer)
	if err != nil {
		panic(err)
	}

	fileInfos, err := ioutil.ReadDir(fileTransfer.TargetDir)
	if err != nil {
		panic(err)
	}

	for _, f := range fileInfos {
		fileTransfer.PutToS3(fileTransfer.TargetDir, f.Name())
	}
}

type FileTransferToS3 struct {
	AwsAccessKeyID       string `yaml:"AWS_ACCESS_KEY_ID"`
	AwsAccessSecretKeyID string `yaml:"AWS_ACCESS_SECRET_KEY_ID"`
	AwsRegion            string `yaml:"AWS_REGION"`
	AwsBucketName        string `yaml:"AWS_BUCKET_NAME"`
	TargetDir            string `yaml:"TARGET_DIR"`
}

func (f *FileTransferToS3) PutToS3(path string, filename string) {
	file, err := os.Open(fmt.Sprintf("%s%s", path, filename))
	if err != nil {
		log.Println(err.Error())
	}
	defer file.Close()

	sess := session.Must(session.NewSession())
	cli := s3.New(sess, &aws.Config{
		Credentials: credentials.NewStaticCredentials(f.AwsAccessKeyID, f.AwsAccessSecretKeyID, ""),
		Region:      aws.String(f.AwsRegion),
	})

	resp, err := cli.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(f.AwsBucketName),
		Key:    aws.String(path + filename),
		Body:   file,
	})
	if err != nil {
		log.Println("upload error")
		log.Println(err.Error())
	}

	log.Println(awsutil.StringValue(resp))
}

package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
    "github.com/bitly/go-simplejson"
	"github.com/disintegration/imaging"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

/* log event
	eventJson, _ := json.MarshalIndent(event, "", "  ")
	log.Printf("EVENT: %s", eventJson)

	log.Printf("REGION: %s", os.Getenv("AWS_REGION"))
	log.Println("ALL ENV VARS:")
	for _, element := range os.Environ() {
	log.Println(element)
}
*/

const tmpDir = "/tmp"

func thumbnail(inputFileName string, width int, height int, outFilename string) {
	srcImage, err := imaging.Open(inputFileName)
	if err != nil {
		panic(err)
	}

	srcImage = imaging.Sharpen(srcImage, 1.0)

	dstImage := imaging.Thumbnail(srcImage, width, height, imaging.Lanczos)
	imaging.Save(dstImage, outFilename)
}

func getSNSCallBackBucket(json_byte []byte) (string, string) {
	// 从SNS回调通知中提取bucket_name和object_key
	mdic, err := simplejson.NewJson(json_byte)
	if err != nil {
		log.Println("err", err)
	}
	record := mdic.Get("Records").GetIndex(0)
	s3_bucket_name, _ := record.Get("s3").Get("bucket").Get("name").String()
	s3_object_key, _ := record.Get("s3").Get("object").Get("key").String()

	return s3_bucket_name, s3_object_key
}

func saveLocalFile(srcBucketName string, srcObjectKey string) (localFileName string, fileName string) {
	// 从S3获取文件内容
	svc := s3.New(session.New())
	input := &s3.GetObjectInput{
		Bucket: aws.String(srcBucketName),
		Key:    aws.String(srcObjectKey),
	}
	result, err := svc.GetObject(input)
	if err != nil {
		log.Println("err", err)
	}
	defer result.Body.Close()

	// 创建本地全路径, 确保目录存在
	sections := strings.Split(srcObjectKey, "/")
	dirName, fileName := sections[0], sections[1]

	log.Println("dirName", dirName)
	log.Println("fileName", fileName)

	localDirName := fmt.Sprintf("%s/%s", tmpDir, dirName)

	_, err = os.Stat(localDirName)
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(localDirName, 0755)
		if errDir != nil {
			log.Fatal(err)
		}
	}

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, result.Body); err != nil {
		log.Println("err", err)
	}

	localFileName = fmt.Sprintf("%s/%s", tmpDir, srcObjectKey)
	log.Println("LocalFileName", localFileName)

	f, err := os.Create(localFileName)
	defer f.Close()
	f.Write(buf.Bytes())

	return localFileName, fileName
}

func saveTargetFile(destFileName string, bucketName string, objectKey string) {
	var err error
	fileHandler, err := os.Open(destFileName)
	if err != nil {
		log.Println("os.Open - filename: %s, err: %v", destFileName, err)
	}
	defer fileHandler.Close()

	svc := s3.New(session.New())
	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   fileHandler,
	})
	log.Println("destBucketName:", bucketName)
	log.Println("destKey:", objectKey)

	if err != nil {
		log.Println("put err", err)
	}
}

func pickupResoultion(resolution string) (width int, height int) {
	sections := strings.Split(resolution, "x")
	width, _ = strconv.Atoi(sections[0])
	height, _ = strconv.Atoi(sections[1])
	return width, height
}

func Handler(ctx context.Context, event events.SNSEvent) {
	// nftavatar.ascendex.io/origin/1.png, 1.jpg

	var srcBucketName, srcObjectKey string
	var localFileName string
	var shortFileName string
	// var region string
	var destBucketName string
	var thumbSetEnv string
	var thumbSetsList []string

	// region = os.Getenv("AWS_REGION")
	destBucketName = os.Getenv("DestinationBucket")

	// 拿到SNS事件
	log.Println("records", event.Records)
	message := event.Records[0].SNS.Message

	// 从S3的Message中提取文件
	srcBucketName, srcObjectKey = getSNSCallBackBucket([]byte(message))

	// 本地存临时文件, 用于生成缩略图
	localFileName, shortFileName = saveLocalFile(srcBucketName, srcObjectKey)

	thumbSetEnv = os.Getenv("ThumbSets")
	thumbSetsList = strings.Split(thumbSetEnv, ",")

	for _, thumbSet := range thumbSetsList {

		width, height := pickupResoultion(thumbSet)
		targetFilename := fmt.Sprintf("%s/%s.jpg", tmpDir, thumbSet)

		// 生成缩略图
		thumbnail(localFileName, width, height, targetFilename)

		destKey := thumbSet + "/" + shortFileName

		// 写目标bucket
		saveTargetFile(targetFilename, destBucketName, destKey)
	}
}

func main() {
	lambda.Start(Handler)
}

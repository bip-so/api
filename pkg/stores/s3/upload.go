package s3

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/disintegration/imaging"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

func UploadObjectToBucket(imagePath string, object io.Reader, isPublic bool) (string, error) {
	client := s3.New(s3.Options{
		Region:      configs.GetAWSS3Config().Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(configs.GetAWSS3Config().AccessKeyID, configs.GetAWSS3Config().AccessSecretKey, "")),
	})

	var acl types.ObjectCannedACL = types.ObjectCannedACLPrivate
	if isPublic {
		acl = types.ObjectCannedACLPublicRead
	}

	_, err := client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: &configs.GetAWSS3Config().BucketName,
		Key:    &imagePath,
		ACL:    acl,
		Body:   object,
	})

	if err != nil {
		return "", err
	}

	url := configs.GetAWSS3Config().CloudFrontURL + "/" + imagePath
	if !isPublic {
		url = "https://" + configs.GetAWSS3Config().BucketName + ".s3." + configs.GetAWSS3Config().Region + ".amazonaws.com/" + imagePath
	}

	return url, nil
}

type ImageResponse struct {
	URL     string
	TinyURL string
}

func UploadImageToBucket(imagePath string, object io.Reader, isPublic bool, isTiny bool) (*ImageResponse, error) {

	byts, err := ioutil.ReadAll(object)
	if err != nil {
		return nil, nil
	}

	object = bytes.NewReader(byts)

	url, err := UploadObjectToBucket(imagePath, object, isPublic)
	if err != nil {
		return nil, nil
	}

	if !isTiny {
		return &ImageResponse{
			URL:     url,
			TinyURL: "",
		}, nil
	}

	object = bytes.NewReader(byts)

	image, _, err := image.Decode(object)
	if err != nil {
		fmt.Println(err)
	}

	imageTiny := imaging.Resize(image, 60, 60, imaging.Lanczos)

	tinyObjectBuff := new(bytes.Buffer)
	err = png.Encode(tinyObjectBuff, imageTiny)
	if err != nil {
		fmt.Println(err)
	}
	tinyObject := bytes.NewReader(tinyObjectBuff.Bytes())

	re, _ := regexp.Compile(".(jpg|jpeg|png|gif|JPG|JPEG|PNG|GIF)$")
	var tinyImagePath string
	if match := re.Match([]byte(imagePath)); match {
		strs := strings.Split(imagePath, ".")
		strs[len(strs)-2] = strs[len(strs)-2] + "-tiny"
		tinyImagePath = strings.Join(strs, ".")
	} else {
		tinyImagePath = imagePath + "-tiny"
	}

	tinyURL, err := UploadObjectToBucket(tinyImagePath, tinyObject, isPublic)
	if err != nil {
		return nil, nil
	}
	return &ImageResponse{
		URL:     url,
		TinyURL: tinyURL,
	}, nil
}

func UploadImageFromURLToS3(srcURL, imagePath string, isPublic bool, isTiny bool) (*ImageResponse, error) {
	response, err := http.Get(srcURL)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer response.Body.Close()

	return UploadImageToBucket(imagePath, response.Body, isPublic, isTiny)
}

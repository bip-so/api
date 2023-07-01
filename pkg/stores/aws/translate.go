package aws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/translate"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

func Translate(text, srcLanguageCode, targetLanguageCode string) (string, error) {
	accessId := configs.GetAWSS3Config().AccessKeyID
	awsSecret := configs.GetAWSS3Config().AccessSecretKey
	fmt.Println("access id and secret ", accessId, awsSecret, configs.GetAWSS3Config().Region)
	credentialsData := credentials.NewStaticCredentials(accessId, awsSecret, "")

	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentialsData,
	})

	awsTranslate := translate.New(sess)
	output, err := awsTranslate.Text(&translate.TextInput{
		Text:               &text,
		SourceLanguageCode: &srcLanguageCode,
		TargetLanguageCode: &targetLanguageCode,
	})
	if err != nil {
		return "", err
	}
	fmt.Println(output)
	return *output.TranslatedText, nil
}

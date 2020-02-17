package email

import (
	"aws-api/setting"
	"aws-api/utils"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"strings"
)

var Sender string


const (
	// The character encoding for the email.
	CharSet = "UTF-8"
)

func SetSender(sender string) {
	Sender = sender
}

func SendEmail(recipients string, subject string, body string) string {

	if Sender == "" {
		Sender = setting.AppSetting.AwsSesSenderEmail
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: aws.StringSlice(strings.Split(recipients, ",")),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(strings.ReplaceAll(body, "\n", "<br>")),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return utils.JsonError(err)
	}

	fmt.Println("Email Sent")
	fmt.Println(result)

	return result.String()
}

func SendRawEmail(recipients string, subject string, body string, file string, fileName string) string {

	raw := "From: {FROM}\nTo: {RECVS}\nSubject: {SUBJECT}\nMIME-Version: 1.0\nContent-type: Multipart/Mixed; boundary=\"NextPart\"\n\n--NextPart\nContent-Type: text/plain\nContent-Transfer-Encoding: base64\n\n{BODY}\n\n--NextPart\nContent-Type: text/plain;\nContent-Disposition: attachment; filename=\"{FILENAME}\"\nContent-Transfer-Encoding: base64\n\n{ATTACHMENT}\n\n--NextPart--"

	raw = strings.Replace(raw, "{FROM}", Sender, -1)
	raw = strings.Replace(raw, "{RECVS}", recipients, -1)
	raw = strings.Replace(raw, "{SUBJECT}", subject, -1)

	bodyBase64 := base64.StdEncoding.EncodeToString([]byte(body))

	raw = strings.Replace(raw, "{BODY}", bodyBase64, -1)

	fileBase64 := base64.StdEncoding.EncodeToString([]byte(file))

	raw = strings.Replace(raw, "{ATTACHMENT}", fileBase64, -1)

	raw = strings.Replace(raw, "{FILENAME}", fileName, -1)

	fmt.Println(raw)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create an SES session.
	svc := ses.New(sess)

	input := &ses.SendRawEmailInput{

		RawMessage: &ses.RawMessage{
			Data: []byte(raw),
		},
	}

	result, err := svc.SendRawEmail(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			case ses.ErrCodeConfigurationSetSendingPausedException:
				fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
			case ses.ErrCodeAccountSendingPausedException:
				fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
			default:
				fmt.Println("*" + aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(">>" + err.Error())
		}

		return utils.JsonError(err)
	}

	return result.String()
}

package awsops

import (
	"aws-api/setting"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"log"
	"net/http"
	"os"
)

func ConnectAWS() (sess *session.Session) {

	//var sess *session.Session

	if setting.AppSetting.AwsCredential == "SharedConfigEnable" {
		// Load session from shared config
		sess = session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	} else if setting.AppSetting.AwsCredential == "NewStaticCredentials" {

		sess, _ = session.NewSession(&aws.Config{
			Region:      aws.String(setting.AppSetting.AwsRegion),
			Credentials: credentials.NewStaticCredentials(setting.AppSetting.AwsAccessKeyId, setting.AppSetting.AwsSecretAccessKey, ""),
		})
	} else if setting.AppSetting.AwsCredential == "Environment" {

		aws_region := os.Getenv("AWS_DEFAULT_REGION")
		aws_user := os.Getenv("AWS_ACCESS_KEY_ID")
		aws_pass := os.Getenv("AWS_SECRET_ACCESS_KEY")
		sess, _ = session.NewSession(&aws.Config{
			Region:      aws.String(aws_region),
			Credentials: credentials.NewStaticCredentials(aws_user, aws_pass, ""),
		})
	} else {

		log.Fatal("Invalid INI parameter: setting.AppSetting.AwsCredential")
	}
	return
}

func ReleaseAddress(public_ip string) (*ec2.ReleaseAddressOutput, error) {
	sess := ConnectAWS()

	svc := ec2.New(sess)

	input := &ec2.ReleaseAddressInput{
		PublicIp: aws.String(public_ip),
	}

	result, err := svc.ReleaseAddress(input)

	return result, err

}

func DescribeInstance(instanceId string) (*ec2.DescribeInstancesOutput, error) {

	sess := ConnectAWS()

	svc := ec2.New(sess)

	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []*string{aws.String(instanceId)},
			},
		},
	}

	result, err := svc.DescribeInstances(input)

	return result, err
}

func TerminateInstances(instanceId string) (*ec2.TerminateInstancesOutput, error) {

	sess := ConnectAWS()

	svc := ec2.New(sess)

	input := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceId),
		},
	}

	result, err := svc.TerminateInstances(input)

	return result, err

}

func GetDisableTerminationProtection(instanceId string) (*ec2.DescribeInstanceAttributeOutput, error) {

	sess := ConnectAWS()

	svc := ec2.New(sess)
	input := &ec2.DescribeInstanceAttributeInput{
		Attribute:  aws.String("disableApiTermination"),
		InstanceId: aws.String(instanceId),
	}

	result, err := svc.DescribeInstanceAttribute(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return nil, aerr
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, err
		}
		return nil, nil
	}
	return result, nil
}

func SetDisableTerminationProtection(instanceId string) (*ec2.ModifyInstanceAttributeOutput, error) {

	sess := ConnectAWS()

	svc := ec2.New(sess)

	input := &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instanceId),
		DisableApiTermination: &ec2.AttributeBooleanValue{
			Value: aws.Bool(true),
		},
	}

	result, err := svc.ModifyInstanceAttribute(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return nil, aerr
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			return nil, err
		}
		return nil, nil
	}
	return result, nil
}

func CreateSnapshot(instanceId, tagName, tagValue, tagStateName, instanceType, publicIp string) (*string, error) {

	sess := ConnectAWS()

	svc := ec2.New(sess)

	response := describeVolumes(svc, instanceId)

	for _, volume := range response.Volumes {

		tagList := &ec2.TagSpecification{
			Tags: []*ec2.Tag{
				{Key: aws.String(tagName), Value: aws.String(tagValue)},
				{Key: aws.String("StateName"), Value: aws.String(tagStateName)},
				{Key: aws.String("Type"), Value: aws.String(instanceType)},
				{Key: aws.String("PublicIP"), Value: aws.String(publicIp)},
			},
			ResourceType: aws.String(ec2.ResourceTypeSnapshot),
		}

		input := &ec2.CreateSnapshotInput{
			Description:       aws.String("AWS-API Snapshot"),
			VolumeId:          volume.VolumeId,
			TagSpecifications: []*ec2.TagSpecification{tagList},
		}

		_, err := svc.CreateSnapshot(input)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				default:
					fmt.Println("a Error :" + aerr.Error())
				}
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println("Error:" + err.Error())
			}
			return nil, err
		}

	}

	result := "ok"

	return &result, nil

}

func describeVolumes(svc *ec2.EC2, instanceId string) *ec2.DescribeVolumesOutput {

	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("attachment.instance-id"),
				Values: []*string{
					aws.String(instanceId),
				},
			},
			/*
				{
					Name: aws.String("attachment.delete-on-termination"),
					Values: []*string{
						aws.String("true"),
					},
				},
			*/
		},
	}

	result, err := svc.DescribeVolumes(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return &ec2.DescribeVolumesOutput{}
	}

	return result

}


func DescribeAddresses(r *http.Request) (*ec2.DescribeAddressesOutput, error) {
	sess := ConnectAWS()

	svc := ec2.New(sess)

	var input *ec2.DescribeAddressesInput

	if r.Form.Get("association-id") != "" {

		input = &ec2.DescribeAddressesInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("association-id"),
					Values: []*string{aws.String(r.Form.Get("association-id"))},
				},
			},
		}
	} else {
		log.Println("DescribeAddresses: All")
		input = &ec2.DescribeAddressesInput{}
	}

	result, err := svc.DescribeAddresses(input)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return &ec2.DescribeAddressesOutput{}, err
	}

	return result, nil
}

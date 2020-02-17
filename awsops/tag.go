package awsops

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func DeleteTag(resource, tagName, tagValue string) (*ec2.DeleteTagsOutput, error) {

	sess := ConnectAWS()

	ec2svc := ec2.New(sess)

	input := &ec2.DeleteTagsInput{
		Resources: []*string{aws.String(resource)},
		Tags:      []*ec2.Tag{{Key: aws.String(tagName), Value: aws.String(tagValue)}},
	}

	result, err := ec2svc.DeleteTags(input)

	return result, err
}

func CreateTag(resource, tagName, tagValue string) (*ec2.CreateTagsOutput, error) {

	sess := ConnectAWS()

	ec2svc := ec2.New(sess)

	input := &ec2.CreateTagsInput{
		Resources: []*string{aws.String(resource)},
		Tags:      []*ec2.Tag{{Key: aws.String(tagName), Value: aws.String(tagValue)}},
	}

	result, err := ec2svc.CreateTags(input)

	return result, err
}

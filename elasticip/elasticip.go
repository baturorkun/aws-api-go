package elasticip

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"os"
	"strings"
)

func AllocateIP(number int) string {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)
	if err != nil {
		exitErrorf("Session, %v", err)
	}
	// Create an EC2 service client.
	svc := ec2.New(sess)

	ips := []string{}

	for i := 0; i < number; i++ {
		// Attempt to allocate the Elastic IP address.
		allocRes, err := svc.AllocateAddress(&ec2.AllocateAddressInput{
			Domain: aws.String("vpc"),
		})
		if err != nil {
			exitErrorf("Unable to allocate IP address, %v", err)
		}
		ips = append(ips, *allocRes.PublicIp)
	}

	res := map[string]string{"ips": strings.Join(ips, ","), "region": "us-east-1"}

	urlsJson, _ := json.Marshal(res)

	return string(urlsJson)
}

func GetUsingIPs() {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	// Create an EC2 service client.
	svc := ec2.New(sess)

	result, err := svc.DescribeAddresses(&ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("domain"),
				Values: aws.StringSlice([]string{"vpc"}),
			},
		},
	})
	if err != nil {
		exitErrorf("Unable to elastic IP address, %v", err)
	}

	// Printout the IP addresses if there are any.
	if len(result.Addresses) == 0 {
		fmt.Printf("No elastic IPs for %s region\n", *svc.Config.Region)
	} else {
		fmt.Println("Elastic IPs")
		for _, addr := range result.Addresses {
			fmt.Println("*", fmtAddress(addr))
		}
	}
}

/*
 Private Functions
*/

func fmtAddress(addr *ec2.Address) string {
	out := fmt.Sprintf("IP: %s,  allocation id: %s",
		aws.StringValue(addr.PublicIp), aws.StringValue(addr.AllocationId))
	if addr.InstanceId != nil {
		out += fmt.Sprintf(", instance-id: %s", *addr.InstanceId)
	}
	return out
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

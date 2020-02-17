package main

import (
	"aws-api/awsops"
	"aws-api/cost"
	"aws-api/elasticip"
	"aws-api/email"
	"aws-api/setting"
	"aws-api/utils"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/gorilla/mux"
	"log"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type KeyVal struct {
	Key   string
	Value string
}

func runMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here

		ip, _, _ := net.SplitHostPort(r.RemoteAddr)

		log.Printf("Request from " + ip + " called " + r.RequestURI)

		AllowIps := strings.Replace(setting.AppSetting.AllowIps, " ", "", -1)

		AllowIpsSlice := strings.Split(AllowIps, ",")

		if ip == "127.0.0.1" || ip == "::1" || ip[0:3] == "192" || ip[0:3] == "172" || utils.StringInSlice(ip, AllowIpsSlice) {
			next.ServeHTTP(w, r)
			return
		} else {
			myToken := tokenGenerator(ip, setting.AppSetting.TokenSalt)

			r.ParseForm()

			var reqToken string

			if r.Form.Get("token") != "" {
				reqToken = r.Form.Get("token")
			} else {
				reqToken = r.Header.Get("Token")
			}

			if reqToken == myToken {
				next.ServeHTTP(w, r)
				return
			}

			fmt.Fprintf(w, "Invalid Token")
			//fmt.Fprintf(w, ip + "-" + setting.AppSetting.TokenSalt + "-" + myToken + "-" + reqToken)
		}

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		//next.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome AWS-API")
	fmt.Fprintf(w, "(by Batur Orkun)")
}

func instanceSearch(w http.ResponseWriter, r *http.Request) {

	sess := awsops.ConnectAWS()

	ec2svc := ec2.New(sess)

	r.ParseForm()

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running"), aws.String("stopped")},
			},
		},
	}

	if r.Form.Get("tag-name") != "" && r.Form.Get("tag-value") != "" {

		filter := ec2.Filter{
			Name:   aws.String("tag:" + r.Form.Get("tag-name")),
			Values: []*string{aws.String(r.Form.Get("tag-value"))},
		}

		params.Filters = append(params.Filters, &filter)
	}

	resp, err := ec2svc.DescribeInstances(params)


	if err != nil {
		fmt.Println("There was an error listing instances in", err.Error())
		log.Fatal(err.Error())
	}

	type Instance struct {
		Region          string
		InstanceId      string
		PublicIpAddress string
		StateName       string
		InstanceType    string
		Tags            []KeyVal
	}

	arr := make([]Instance, len(resp.Reservations))

	for idx, res := range resp.Reservations {

		for _, inst := range res.Instances {

			arr[idx].Region = "us-east-1"
			arr[idx].InstanceId = *inst.InstanceId
			arr[idx].InstanceType = *inst.InstanceType
			if inst.NetworkInterfaces[0].PrivateIpAddresses[0].Association != nil {
				arr[idx].PublicIpAddress = *inst.NetworkInterfaces[0].PrivateIpAddresses[0].Association.PublicIp
			} else {
				arr[idx].PublicIpAddress = ""
			}
			arr[idx].StateName = *inst.State.Name
			for _, tag := range inst.Tags {
				item := KeyVal{Key: *tag.Key, Value: *tag.Value}
				arr[idx].Tags = append(arr[idx].Tags, item)
			}
		}
	}

	pagesJson, err := json.Marshal(arr)


	if err != nil {
		log.Println("Cannot encode to JSON ", err)
	}

	fmt.Fprintf(w, "%s", pagesJson)

}

func instanceStop(w http.ResponseWriter, r *http.Request) {

	sess := awsops.ConnectAWS()
	ec2svc := ec2.New(sess)

	r.ParseForm()

	input := &ec2.StopInstancesInput{
		InstanceIds: []*string{
			aws.String(r.Form.Get("instance-id")),
		},
		DryRun: aws.Bool(false),
	}
	result, err := ec2svc.StopInstances(input)

	if err != nil {
		data := map[string]string{"Error": err.Error()}
		jsonData, _ := json.Marshal(data)

		fmt.Fprintf(w, string(jsonData))
	} else {
		fmt.Fprint(w, result)
	}

}

func instanceStart(w http.ResponseWriter, r *http.Request) {

	sess := awsops.ConnectAWS()
	ec2svc := ec2.New(sess)

	r.ParseForm()

	input := &ec2.StartInstancesInput{
		InstanceIds: []*string{
			aws.String(r.Form.Get("instance-id")),
		},
		DryRun: aws.Bool(false),
	}
	result, err := ec2svc.StartInstances(input)

	if err != nil {

		data := map[string]string{"Error": err.Error()}
		jsonData, _ := json.Marshal(data)

		fmt.Fprintf(w, string(jsonData))
	} else {
		fmt.Fprint(w, result)
	}

}

func sesSendEmail(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var result string

	if r.Form.Get("file") != "" {
		log.Println("SendRawEmail")
		result = email.SendRawEmail(r.Form.Get("recipients"), r.Form.Get("subject"), r.Form.Get("body"), r.Form.Get("file"), r.Form.Get("filename"))
	} else {
		log.Println("SendEmail")
		result = email.SendEmail(r.Form.Get("recipients"), r.Form.Get("subject"), r.Form.Get("body"))
	}

	fmt.Fprint(w, result)
}

func elasticipAllocate(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	number, err := strconv.Atoi(r.Form.Get("number"))

	if err != nil {
		fmt.Fprint(w, "Error", err)
	}

	json := elasticip.AllocateIP(number)

	fmt.Fprint(w, json)
}

func lbtargetgroupSearch(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	sess := awsops.ConnectAWS()

	svc := elbv2.New(sess)

	var input *elbv2.DescribeTargetGroupsInput
	var names = r.Form.Get("names")

	namesArr := strings.Split(names, ",")

	if names == "" {
		input = &elbv2.DescribeTargetGroupsInput{}
	} else {
		input = &elbv2.DescribeTargetGroupsInput{
			Names: aws.StringSlice(namesArr),
		}
	}

	result, err := svc.DescribeTargetGroups(input)

	if err != nil {

		data := map[string]string{"Error": err.Error()}
		jsonData, _ := json.Marshal(data)

		fmt.Fprintf(w, string(jsonData))
	} else {

		data := map[string]string{}

		for _, item := range result.TargetGroups {
			rt := reflect.TypeOf(item.LoadBalancerArns)

			if rt.Kind() == reflect.Slice {
				if len(item.LoadBalancerArns) > 0 {
					arr := strings.Split(*item.LoadBalancerArns[0], "/")
					data[*item.TargetGroupName] = arr[len(arr)-2]
				}
			}
		}

		jsonData, _ := json.Marshal(data)

		fmt.Fprintf(w, string(jsonData))
	}

}

func billingDaily(w http.ResponseWriter, r *http.Request) {

	sess := awsops.ConnectAWS()

	result := cost.DailyCost(sess)

	fmt.Fprint(w, result)

}

func billingMonthly(w http.ResponseWriter, r *http.Request) {

	sess := awsops.ConnectAWS()

	result := cost.MonthlyCost(sess)

	fmt.Fprint(w, result)

}

func billingByTag(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	if r.Form.Get("cid") == "" {
		data := map[string]string{"Error": "Missing cid parameter"}
		jsonData, _ := json.Marshal(data)

		fmt.Fprintf(w, string(jsonData))
		return
	}

	sess := awsops.ConnectAWS()

	result := cost.TagFilteredMonthlyCost(sess, r.Form.Get("tagName"), r.Form.Get("tagValue"))

	fmt.Fprint(w, result)
}


func remoteCopySshKey(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	out := awsops.CopyShhKey(r.Form.Get("sshkey"), r.Form.Get("public-ip"))

	fmt.Fprint(w, out)

}


func remoteGetMessagesLog(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	out := awsops.GetMessagesLog(r.Form.Get("public-ip"), r.Form.Get("lines"))

	fmt.Fprint(w, out)
}



func instanceTerminate(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	instanceId := r.Form.Get("instance_id")
	// releaseIp := r.Form.Get("release_ip")   // not ready

	res, err := awsops.DescribeInstance(instanceId)

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	if len(res.Reservations) > 0 {
		if *res.Reservations[0].Instances[0].State.Name != "stopped" {
			fmt.Fprint(w, errors.New("Only stopped instances can be terminated"))
			return
		}
	}

	result, err := awsops.TerminateInstances(instanceId)

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				fmt.Fprint(w, aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Fprint(w, err.Error())
		}
		return
	}

	fmt.Fprint(w, result)

}

func elasticipRelease(w http.ResponseWriter, r *http.Request) {

		// Not ready
}

func instanceGetDisableApiTermination(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	instanceId := r.Form.Get("instance_id")

	result, err := awsops.GetDisableTerminationProtection(instanceId)

	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, result)
	}
}

func instanceSetDisableApiTermination(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	instanceId := r.Form.Get("instance_id")

	result, err := awsops.SetDisableTerminationProtection(instanceId)

	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, result)
	}
}

func tagDelete(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	instanceId := r.Form.Get("instance-id")
	tag := r.Form.Get("tag")
	value := r.Form.Get("value")

	result, err := awsops.DeleteTag(instanceId, tag, value)

	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, result)
	}
}

func tagCreate(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	instanceId := r.Form.Get("instance_id")
	tag := r.Form.Get("tag")
	value := r.Form.Get("value")

	result, err := awsops.CreateTag(instanceId, tag, value)

	if err != nil {
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, result)
	}
}

func snapshotCreate(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	instanceId := r.Form.Get("instance-id")
	tagName := r.Form.Get("tag-name")
	tagValue := r.Form.Get("tag-value")
	tagStateName := r.Form.Get("state-name")
	instanceType := r.Form.Get("instance-type")
	publicIp := r.Form.Get("public-ip")

	result, _ := awsops.CreateSnapshot(instanceId, tagName, tagValue, tagStateName, instanceType, publicIp)

	fmt.Fprint(w, *result)
}


func elasticipSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	result, err := awsops.DescribeAddresses(r)

	if err != nil {
		log.Println("Error")
		fmt.Fprint(w, err)
	} else {
		fmt.Fprint(w, result)
	}
}


func main() {

	setting.Setup()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", home)
	router.HandleFunc("/ses-sendemail", sesSendEmail)
	router.HandleFunc("/instance-search", instanceSearch)
	router.HandleFunc("/instance-start", instanceStart)
	router.HandleFunc("/instance-stop", instanceStop)
	router.HandleFunc("/instance-terminate", instanceTerminate)
	router.HandleFunc("/instance-setting-getdisabletermination", instanceGetDisableApiTermination)
	router.HandleFunc("/instance-setting-setdisabletermination", instanceSetDisableApiTermination)
	router.HandleFunc("/snapshot-create", snapshotCreate)
	router.HandleFunc("/elasticip-allocate", elasticipAllocate)
	router.HandleFunc("/elasticip-search", elasticipSearch)
	router.HandleFunc("/elasticip-release", elasticipRelease)
	router.HandleFunc("/lbtargetgroup-search", lbtargetgroupSearch)
	router.HandleFunc("/tag-delete", tagDelete)
	router.HandleFunc("/tag-create", tagCreate)
	router.HandleFunc("/billing-daily", billingDaily)
	router.HandleFunc("/billing-monthly", billingMonthly)
	router.HandleFunc("/billing-bytag", billingByTag)
	router.HandleFunc("/remote-copy-sshkey", remoteCopySshKey)
	router.HandleFunc("/remote-get-messsages-log", remoteGetMessagesLog)

	router.Use(runMiddleware)

	log.Println("Starting server on " + setting.ServerSetting.Port)

	log.Fatal(http.ListenAndServe(setting.ServerSetting.Port, router))

}
func tokenGenerator(ip string, salt string) string {

	h := sha256.New()
	h.Write([]byte(ip + salt))

	return hex.EncodeToString(h.Sum(nil))

}

func check(e error) {

	if e != nil {
		panic(e)
	}
}

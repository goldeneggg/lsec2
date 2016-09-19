package awsec2

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

import (
	"github.com/goldeneggg/lsec2/constants"
)

const (
	fieldSeparater = "\t"
)

type Opt struct {
	PrintHeader   bool
	OnlyPrivateIP bool
	Region        string
	Tags          []string
}

type instanceInfo struct {
	instanceID       string            `header:"INSTANCE_ID"`
	privateIPAddress string            `header:"PRIVATE_IP"`
	publicIPAddress  string            `header:"PUBLIC_IP"`
	instanceType     string            `header:"TYPE"`
	stateName        string            `header:"STATE"`
	tags             map[string]string `header:"TAGS"`
}

func (opt *Opt) Show() int {
	infos, err := opt.buildInfos()
	if err != nil {
		fmt.Errorf("error: %v", err)
		return constants.ExitStsNg
	}

	opt.printInfos(infos)

	return constants.ExitStsOk
}

func (opt *Opt) printInfos(infos []*instanceInfo) {
	if opt.PrintHeader {
		printHeader(*infos[0])
	}

	for _, info := range infos {
		if opt.OnlyPrivateIP {
			fmt.Printf("%s\n", info.privateIPAddress)
		} else {
			printRow(info)
		}
	}
}

func (opt *Opt) buildInfos() ([]*instanceInfo, error) {
	svc := ec2.New(session.New(), &aws.Config{Region: aws.String(opt.Region)})

	output, err := svc.DescribeInstances(opt.filterParams())
	if err != nil {
		return nil, fmt.Errorf("aws describe instances error: %v", err)
	}

	var infos []*instanceInfo

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			infos = append(infos, newInstanceInfo(instance))
		}
	}

	return infos, nil
}

func (opt *Opt) filterParams() *ec2.DescribeInstancesInput {
	var filters []*ec2.Filter

	// opt.tags is separated key-value pair by "=", and values are separated by ","(comma)
	// ex. "Name=Value"
	// ex. "Name=Value1,Value2"
	for _, tag := range opt.Tags {
		tagNameValue := strings.Split(tag, "=")
		name := aws.String("tag:" + tagNameValue[0])
		values := make([]*string, 0, 3)
		for _, value := range strings.Split(tagNameValue[1], ",") {
			values = append(values, aws.String(value))
		}

		tagFilter := &ec2.Filter{
			Name:   name,
			Values: values,
		}
		filters = append(filters, tagFilter)
	}

	if len(filters) == 0 {
		return nil
	}

	return &ec2.DescribeInstancesInput{Filters: filters}
}

func newInstanceInfo(instance *ec2.Instance) *instanceInfo {
	i := &instanceInfo{
		privateIPAddress: fetchItem(instance.PrivateIpAddress),
		instanceID:       fetchItem(instance.InstanceId),
		instanceType:     fetchItem(instance.InstanceType),
		stateName:        fetchItem(instance.State.Name),
		publicIPAddress:  fetchItem(instance.PublicIpAddress),
	}

	tags := make(map[string]string)
	for _, tag := range instance.Tags {
		tags[*tag.Key] = *tag.Value
	}

	if len(tags) > 0 {
		i.tags = tags
	}

	return i
}

func printHeader(info instanceInfo) {
	var headers []string

	rt := reflect.TypeOf(info)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		headers = append(headers, field.Tag.Get("header"))
	}

	fmt.Printf("%s\n", strings.Join(headers, fieldSeparater))
}

func printRow(info *instanceInfo) {
	fmt.Printf("%s\n", parseRow(info))
}

func parseRow(info *instanceInfo) string {
	var values []string

	values = append(values, info.instanceID)
	values = append(values, info.privateIPAddress)
	values = append(values, info.publicIPAddress)
	values = append(values, info.instanceType)
	values = append(values, info.stateName)

	var tagNames []string
	for name := range info.tags {
		tagNames = append(tagNames, name)
	}
	sort.Strings(tagNames)

	var tags []string
	for _, name := range tagNames {
		tags = append(tags, name+"="+info.tags[name])
	}
	values = append(values, strings.Join(tags, ","))

	return fmt.Sprintf("%s", strings.Join(values, fieldSeparater))
}

func fetchItem(instanceItem *string) string {
	if len(aws.StringValue(instanceItem)) == 0 {
		return "UNDEFINED"
	}

	return *instanceItem
}

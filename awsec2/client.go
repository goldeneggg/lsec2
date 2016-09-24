package awsec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

import (
	"github.com/goldeneggg/lsec2/constants"
)

type Client struct {
	PrintHeader   bool
	OnlyPrivateIP bool
	Region        string
	Tags          []string
}

func (client *Client) Print() error {
	infos, err := client.buildInfos()
	if err != nil {
		return fmt.Errorf("buildInfos error: %v", err)
	}

	client.printInfos(infos)

	return nil
}

func (client *Client) printInfos(infos []*instanceInfo) {
	if client.PrintHeader {
		infos[0].printHeader()
	}

	for _, info := range infos {
		if client.OnlyPrivateIP {
			fmt.Printf("%s\n", info.privateIPAddress)
		} else {
			info.printRow()
		}
	}
}

func (client *Client) buildInfos() ([]*instanceInfo, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(client.Region)})
	if err != nil {
		return nil, fmt.Errorf("aws new session error: %v", err)
	}

	svc := ec2.New(sess)

	output, err := svc.DescribeInstances(client.filterParams())
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

func (client *Client) filterParams() *ec2.DescribeInstancesInput {
	var filters []*ec2.Filter

	// client.tags is separated key-value pair by "=", and values are separated by ","(comma)
	// ex. "Name=Value"
	// ex. "Name=Value1,Value2"
	for _, tag := range client.Tags {
		tagNameValue := strings.Split(tag, constants.TagPairSeparator)
		name := aws.String(constants.TagFilterPrefix + tagNameValue[0])
		values := make([]*string, 0, 3)
		for _, value := range strings.Split(tagNameValue[1], constants.TagValueSeparator) {
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

func fetchItem(instanceItem *string) string {
	if len(aws.StringValue(instanceItem)) == 0 {
		return constants.UndefinedItem
	}

	return *instanceItem
}

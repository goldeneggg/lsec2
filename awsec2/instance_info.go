package awsec2

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type InstanceInfo struct {
	InstanceID       string            `header:"INSTANCE_ID"`
	PrivateIPAddress string            `header:"PRIVATE_IP"`
	PublicIPAddress  string            `header:"PUBLIC_IP"`
	InstanceType     string            `header:"TYPE"`
	StateName        string            `header:"STATE"`
	Tags             map[string]string `header:"TAGS"`
}

func NewInstanceInfo(instance *ec2.Instance) (*InstanceInfo, error) {
	if instance == nil {
		return nil, fmt.Errorf("ec2.Instance: %v is invalid", instance)
	}

	i := &InstanceInfo{
		InstanceID:       fetchItem(instance.InstanceId),
		PrivateIPAddress: fetchItem(instance.PrivateIpAddress),
		PublicIPAddress:  fetchItem(instance.PublicIpAddress),
		InstanceType:     fetchItem(instance.InstanceType),
		StateName:        fetchItem(instance.State.Name),
	}

	tags := make(map[string]string)
	for _, tag := range instance.Tags {
		tags[*tag.Key] = *tag.Value
	}

	if len(tags) > 0 {
		i.Tags = tags
	}

	return i, nil
}

func (info *InstanceInfo) ParseRow() string {
	var values []string

	values = append(values, info.InstanceID)
	values = append(values, info.PrivateIPAddress)
	values = append(values, info.PublicIPAddress)
	values = append(values, info.InstanceType)
	values = append(values, info.StateName)

	values = append(values, strings.Join(info.parseTags(), ParsedTagSeparator))

	return fmt.Sprintf("%s", strings.Join(values, FieldSeparater))
}

func (info *InstanceInfo) printHeader() {
	var headers []string

	rt := reflect.TypeOf(*info)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		headers = append(headers, field.Tag.Get("header"))
	}

	fmt.Printf("%s\n", strings.Join(headers, FieldSeparater))
}

func (info *InstanceInfo) printRow() {
	fmt.Printf("%s\n", info.ParseRow())
}

func (info *InstanceInfo) parseTags() []string {
	var tagNames []string
	for name := range info.Tags {
		tagNames = append(tagNames, name)
	}
	sort.Strings(tagNames)

	var tags []string
	for _, name := range tagNames {
		tags = append(tags, name+"="+info.Tags[name])
	}

	return tags
}

func fetchItem(instanceItem *string) string {
	if len(aws.StringValue(instanceItem)) == 0 {
		return UndefinedItem
	}

	return *instanceItem
}

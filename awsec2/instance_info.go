package awsec2

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fatih/color"
)

// InstanceInfo is attributes of aws ec2 instance
type InstanceInfo struct {
	InstanceID       string            `header:"INSTANCE_ID"`
	PrivateIPAddress string            `header:"PRIVATE_IP"`
	PublicIPAddress  string            `header:"PUBLIC_IP"`
	InstanceType     string            `header:"TYPE"`
	StateName        string            `header:"STATE"`
	Tags             map[string]string `header:"TAGS"`
}

// NewInstanceInfo creates a new InstanceInfo
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

// ParseRow parses from InstanceInfo to one line string
func (info *InstanceInfo) ParseRow(withColor bool) string {
	var values []string

	values = append(values, info.InstanceID)
	values = append(values, info.PrivateIPAddress)
	values = append(values, info.PublicIPAddress)
	values = append(values, info.InstanceType)
	values = append(values, info.decorateStateName(withColor))

	values = append(values, strings.Join(info.parseTags(), parsedTagSeparator))

	return fmt.Sprintf("%s", strings.Join(values, fieldSeparater))
}

func (info *InstanceInfo) printHeader() {
	var headers []string

	rt := reflect.TypeOf(*info)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		headers = append(headers, field.Tag.Get("header"))
	}

	fmt.Printf("%s\n", strings.Join(headers, fieldSeparater))
}

func (info *InstanceInfo) printRow(withColor bool) {
	fmt.Printf("%s\n", info.ParseRow(withColor))
}

func (info *InstanceInfo) decorateStateName(withColor bool) string {
	if !withColor {
		return info.StateName
	}

	var colorAttr color.Attribute
	switch info.StateName {
	case "running":
		colorAttr = color.FgGreen
	case "stopped":
		colorAttr = color.FgRed
	default:
		colorAttr = color.FgYellow
	}

	return color.New(colorAttr).SprintFunc()(info.StateName)
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
		return undefinedItem
	}

	return *instanceItem
}

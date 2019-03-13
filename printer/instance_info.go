package printer

import (
	"reflect"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/fatih/color"
)

const (
	parsedTagSeparator = ","
	undefinedItem      = "UNDEFINED"
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
func NewInstanceInfo(instance *ec2.Instance) *InstanceInfo {
	if instance == nil {
		return &InstanceInfo{}
	}

	// TODO: coldef.yml対応
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

	return i
}

// Headers gets header string array
func (info *InstanceInfo) Headers() []string {
	var headers []string

	rt := reflect.TypeOf(*info)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		headers = append(headers, field.Tag.Get("header"))
	}

	return headers
}

// Values gets value string array
func (info *InstanceInfo) Values(withColor bool) []string {
	var values []string

	values = append(values, info.InstanceID)
	values = append(values, info.PrivateIPAddress)
	values = append(values, info.PublicIPAddress)
	values = append(values, info.InstanceType)
	values = append(values, info.decorateStateName(withColor))

	values = append(values, strings.Join(info.parseTags(), parsedTagSeparator))

	return values
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

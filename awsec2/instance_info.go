package awsec2

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

const (
	fieldSeparater = "\t"
)

type instanceInfo struct {
	instanceID       string            `header:"INSTANCE_ID"`
	privateIPAddress string            `header:"PRIVATE_IP"`
	publicIPAddress  string            `header:"PUBLIC_IP"`
	instanceType     string            `header:"TYPE"`
	stateName        string            `header:"STATE"`
	tags             map[string]string `header:"TAGS"`
}

func (info *instanceInfo) printHeader() {
	var headers []string

	rt := reflect.TypeOf(*info)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		headers = append(headers, field.Tag.Get("header"))
	}

	fmt.Printf("%s\n", strings.Join(headers, fieldSeparater))
}

func (info *instanceInfo) printRow() {
	fmt.Printf("%s\n", info.parseRow())
}

func (info *instanceInfo) parseRow() string {
	var values []string

	values = append(values, info.instanceID)
	values = append(values, info.privateIPAddress)
	values = append(values, info.publicIPAddress)
	values = append(values, info.instanceType)
	values = append(values, info.stateName)

	values = append(values, strings.Join(info.parseTags(), ","))

	return fmt.Sprintf("%s", strings.Join(values, fieldSeparater))
}

func (info *instanceInfo) parseTags() []string {
	var tagNames []string
	for name := range info.tags {
		tagNames = append(tagNames, name)
	}
	sort.Strings(tagNames)

	var tags []string
	for _, name := range tagNames {
		tags = append(tags, name+"="+info.tags[name])
	}

	return tags
}

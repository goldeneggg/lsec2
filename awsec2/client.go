package awsec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

const (
	tagFilterPrefix   = "tag:"
	tagPairSeparator  = "="
	tagValueSeparator = ","

	defaultTagValuesCap = 3
)

// Client is attributes definition for filtering ec2 instances
type Client struct {
	ec2iface.EC2API
	StateName string
	Tags      []string
}

// NewClient returns a new DefaultClient
func NewClient(region string, profile string) *Client {
	client := new(Client)
	client.EC2API = defaultEC2Client(region, profile)

	return client
}

// NewClientWithEC2API returns a new Client with assigned EC2API variable
func NewClientWithEC2API(maybeEC2Client interface{}) (*Client, error) {
	client := new(Client)

	if ec2Client, ok := maybeEC2Client.(ec2iface.EC2API); ok {
		client.EC2API = ec2Client
		return client, nil
	} else {
		return nil, fmt.Errorf("maybeEC2Client %#v does not implement ec2.EC2API methods", maybeEC2Client)
	}
}

// EC2Instances gets filtered EC2 instances
func (client *Client) EC2Instances() ([]*ec2.Instance, error) {
	output, err := client.EC2API.DescribeInstances(client.buildFilter())
	if err != nil {
		return nil, fmt.Errorf("aws describe instances error: %v", err)
	}

	var res []*ec2.Instance
	for _, r := range output.Reservations {
		res = append(res, r.Instances...)
	}

	return res, nil
}

func defaultEC2Client(region string, profile string) ec2iface.EC2API {
	config := aws.NewConfig()
	if region != "" {
		config = config.WithRegion(region)
	}
	if profile != "" {
		config = config.WithCredentials(credentials.NewSharedCredentials("", profile))
	}
	sess := session.Must(session.NewSession(config))

	return ec2.New(sess)
}

func (client *Client) buildFilter() *ec2.DescribeInstancesInput {
	var filters []*ec2.Filter

	// client.tags is separated key-value pair by "=", and values are separated by ","(comma)
	// ex. "Name=Value"
	// ex. "Name=Value1,Value2"
	for _, tag := range client.Tags {
		kv := strings.Split(tag, tagPairSeparator)
		name := aws.String(tagFilterPrefix + kv[0])
		values := make([]*string, 0, defaultTagValuesCap)
		for _, v := range strings.Split(kv[1], tagValueSeparator) {
			values = append(values, aws.String(v))
		}

		filters = append(filters, &ec2.Filter{Name: name, Values: values})
	}

	if len(filters) == 0 {
		return nil
	}

	return &ec2.DescribeInstancesInput{Filters: filters}
}

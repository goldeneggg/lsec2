package awsec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
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
func NewClient(region, state, profile string, tags []string) (*Client, error) {
	return NewClientWithEC2API(region, state, profile, tags, nil)
}

// NewClientWithEC2API returns a new Client with assigned EC2API variable
func NewClientWithEC2API(region, state, profile string, tags []string, c interface{}) (*Client, error) {
	client := new(Client)

	client.StateName = state
	client.Tags = tags

	if c != nil {
		if ec2, ok := c.(ec2iface.EC2API); ok {
			client.EC2API = ec2
		} else {
			return nil, fmt.Errorf("arg %#v does not implement ec2.EC2API methods", c)
		}
	} else {
		client.EC2API = defaultEC2Client(region, profile)
	}

	return client, nil
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
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true)
	if region != "" {
		config = config.WithRegion(region)
	}

	sessOpts := session.Options{
		Config:                  *config,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		// Load ~/.aws/config regardless of whether AWS_SDK_LOAD_CONFIG env is set
		SharedConfigState: session.SharedConfigEnable,
	}
	if profile != "" {
		sessOpts.Profile = profile
	}

	return ec2.New(session.Must(session.NewSessionWithOptions(sessOpts)))
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

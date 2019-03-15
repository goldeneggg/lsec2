package awsec2_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	. "github.com/goldeneggg/lsec2/awsec2"
)

var (
	dummyInstances = []*ec2.Instance{
		&ec2.Instance{
			InstanceId:       aws.String("i-xxxxxxxxx1"),
			PrivateIpAddress: aws.String("10.0.0.1"),
			PublicIpAddress:  aws.String("1.2.3.4"),
			InstanceType:     aws.String("t1.micro"),
			State: &ec2.InstanceState{
				Name: aws.String("running"),
			},
		},
		&ec2.Instance{
			InstanceId:       aws.String("i-xxxxxxxxx2"),
			PrivateIpAddress: aws.String("10.0.0.2"),
			PublicIpAddress:  aws.String("1.2.3.5"),
			InstanceType:     aws.String("t2.micro"),
			State: &ec2.InstanceState{
				Name: aws.String("stopped"),
			},
		},
		&ec2.Instance{
			InstanceId:       aws.String("i-xxxxxxxxx9"),
			PrivateIpAddress: aws.String("10.0.0.9"),
			PublicIpAddress:  aws.String("1.2.3.9"),
			InstanceType:     aws.String("t2.large"),
			State: &ec2.InstanceState{
				Name: aws.String("stopped"),
			},
		},
	}
)

type mockEC2API struct {
	ec2iface.EC2API
}

// override for mock
func (mock *mockEC2API) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	instancesOutput := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			&ec2.Reservation{
				Instances: []*ec2.Instance{dummyInstances[0], dummyInstances[1]},
			},
			&ec2.Reservation{
				Instances: []*ec2.Instance{dummyInstances[2]},
			},
		},
	}

	return instancesOutput, nil
}

func TestNewClient(t *testing.T) {
	cases := []struct {
		region   string
		profile  string
		expected *Client
	}{
		{
			region:  "ap-northeast-1",
			profile: "",
			expected: &Client{
				StateName: "",
				Tags:      []string{},
			},
		},
		{
			region:  "us-east-1",
			profile: "test",
			expected: &Client{
				StateName: "",
				Tags:      []string{},
			},
		},
	}

	for _, c := range cases {
		client := NewClient(c.region, c.profile)

		if client.StateName != c.expected.StateName {
			t.Errorf("expected: %#v, but actual: %#v", c.expected.StateName, client.StateName)
		}
		if client.Tags != nil {
			t.Errorf("expected: nil, but actual: %#v", client.Tags)
		}
	}
}

func TestNewClientEC2API(t *testing.T) {
	mockClient := &mockEC2API{}

	cases := []struct {
		ec2Client ec2iface.EC2API
		expected  *Client
	}{
		{
			ec2Client: mockClient,
			expected: &Client{
				StateName: "",
				Tags:      []string{},
				EC2API:    mockClient,
			},
		},
	}

	for _, c := range cases {
		client, err := NewClientWithEC2API(c.ec2Client)
		if err != nil {
			t.Errorf("error occured. err: %#v", err)
		}

		if client.StateName != c.expected.StateName {
			t.Errorf("expected: %#v, but actual: %#v", c.expected.StateName, client.StateName)
		}
		if client.Tags != nil {
			t.Errorf("expected: nil, but actual: %#v", client.Tags)
		}
		if client.EC2API != c.expected.EC2API {
			t.Errorf("expected: %#v, but actual: %#v", c.expected.EC2API, client.EC2API)
		}
	}
}

func TestEC2Instances(t *testing.T) {
	cases := []struct {
		client   *Client
		expected []*ec2.Instance
	}{
		{
			client: &Client{
				EC2API:    &mockEC2API{},
				StateName: "running",
				Tags:      []string{"Role=roleA,roleB", "Service=serviceX"},
			},
			expected: dummyInstances,
		},
		{
			client: &Client{
				EC2API:    &mockEC2API{},
				StateName: "stopped",
			},
			expected: dummyInstances,
		},
	}

	for _, c := range cases {
		insts, err := c.client.EC2Instances()

		if err != nil {
			t.Errorf("error occured. err: %#v", err)
		}

		if len(c.expected) != len(insts) {
			t.Errorf("not same instances length. expected: %#v, but actual: %#v", len(c.expected), len(insts))
		}
	}
}

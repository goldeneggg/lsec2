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
		state    string
		profile  string
		tags     []string
		expected *Client
	}{
		{
			region:  "ap-northeast-1",
			state:   "running",
			profile: "",
			tags:    []string{"Tag1=Value1", "Tag2=Value2"},
			expected: &Client{
				StateName: "running",
				Tags:      []string{"Tag1=Value1", "Tag2=Value2"},
			},
		},
		{
			region:  "us-east-1",
			state:   "stopped",
			profile: "test",
			tags:    []string{},
			expected: &Client{
				StateName: "stopped",
				Tags:      []string{},
			},
		},
	}

	for _, c := range cases {
		client, err := NewClient(c.region, c.state, c.profile, c.tags)
		if err != nil {
			t.Errorf("error occured. err: %#v", err)
		}

		if client.StateName != c.expected.StateName {
			t.Errorf("expected: %#v, but actual: %#v", c.expected.StateName, client.StateName)
		}
		for i, tag := range client.Tags {
			if tag != c.expected.Tags[i] {
				t.Errorf("expected: %#v, but actual: %#v", c.expected.Tags[i], tag)
			}
		}
		if _, ok := client.EC2API.(*ec2.EC2); !ok {
			t.Errorf("expected: *ec2.EC2, but actual: %#v", client.EC2API)
		}
	}
}

func TestNewClientEC2API(t *testing.T) {
	mockClient := &mockEC2API{}

	cases := []struct {
		region    string
		state     string
		profile   string
		tags      []string
		ec2Client ec2iface.EC2API
		expected  *Client
	}{
		{
			region:    "ap-northeast-1",
			state:     "running",
			profile:   "",
			tags:      []string{"Tag1=Value1", "Tag2=Value2"},
			ec2Client: mockClient,
			expected: &Client{
				EC2API:    mockClient,
				StateName: "running",
				Tags:      []string{"Tag1=Value1", "Tag2=Value2"},
			},
		},
		{
			region:    "us-east-1",
			state:     "stopped",
			profile:   "test",
			tags:      []string{},
			ec2Client: mockClient,
			expected: &Client{
				EC2API:    mockClient,
				StateName: "stopped",
				Tags:      []string{},
			},
		},
	}

	for _, c := range cases {
		client, err := NewClientWithEC2API(c.region, c.state, c.profile, c.tags, c.ec2Client)
		if err != nil {
			t.Errorf("error occured. err: %#v", err)
		}

		if client.StateName != c.expected.StateName {
			t.Errorf("expected: %#v, but actual: %#v", c.expected.StateName, client.StateName)
		}
		for i, tag := range client.Tags {
			if tag != c.expected.Tags[i] {
				t.Errorf("expected: %#v, but actual: %#v", c.expected.Tags[i], tag)
			}
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

		for i, inst := range insts {
			if inst.InstanceId != c.expected[i].InstanceId {
				t.Errorf("expected: %#v, but actual: %#v", c.expected[i].InstanceId, inst.InstanceId)
			}
			if inst.PrivateIpAddress != c.expected[i].PrivateIpAddress {
				t.Errorf("expected: %#v, but actual: %#v", c.expected[i].PrivateIpAddress, inst.PrivateIpAddress)
			}
			if inst.PublicIpAddress != c.expected[i].PublicIpAddress {
				t.Errorf("expected: %#v, but actual: %#v", c.expected[i].PublicIpAddress, inst.PublicIpAddress)
			}
			if inst.InstanceType != c.expected[i].InstanceType {
				t.Errorf("expected: %#v, but actual: %#v", c.expected[i].InstanceType, inst.InstanceType)
			}
			if inst.State.Name != c.expected[i].State.Name {
				t.Errorf("expected: %#v, but actual: %#v", c.expected[i].State.Name, inst.State.Name)
			}
		}
	}
}

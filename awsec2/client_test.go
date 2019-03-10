package awsec2_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	. "github.com/goldeneggg/lsec2/awsec2"
)

type (
	caseNewClient struct {
		region string
		exp    *DefaultClient
	}

	caseEC2Instances struct {
		client *DefaultClient
		exp    []*ec2.Instance
	}
)

type mockEC2API struct {
	ec2iface.EC2API
}

func (mock *mockEC2API) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	// TODO
	return nil, nil
}

func TestNewClient(t *testing.T) {
	cases := []caseNewClient{
		{"ap-northeast-1", &DefaultClient{Region: "ap-northeast-1"}},
		{"us-east-1", &DefaultClient{Region: "us-east-1"}},
	}

	for _, c := range cases {
		client := NewClient(c.region)

		dClient, ok := client.(*DefaultClient)
		if !ok {
			t.Errorf("client is not *DefaultClient: %#v, but exp: %#v", c.exp.Region, dClient.Region)
		}

		if dClient.Region != c.exp.Region {
			t.Errorf("expected: %#v, but exp: %#v", c.exp.Region, dClient.Region)
		}
	}
}

func TestEC2Instances(t *testing.T) {
	cases := []caseEC2Instances{
		{
			client: &DefaultClient{EC2API: &mockEC2API{}, Region: "ap-northeast-1"},
			exp:    []*ec2.Instance{},
		},
	}

	for _, c := range cases {
		insts, err := c.client.EC2Instances()

		if err != nil {
			t.Errorf("error occured. err: %#v", err)
		}

		if len(c.exp) != len(insts) {
			t.Errorf("not same instances length. expected: %#v, but exp: %#v", len(c.exp), len(insts))
		}
	}
}

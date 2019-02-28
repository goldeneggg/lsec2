package awsec2

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

const (
	dummyInstanceID       = "i-xxxxxxxx"
	dummyPrivateIPAddress = "10.0.0.1"
	dummyPublicIPAddress  = "1.2.3.4"
	dummyInstanceType     = "t2.micro"
	dummyStateCode        = 16
	dummyStateName        = "running"
)

var (
	dummyTagKeys   = []string{"Name", "Role"}
	dummyTagValues = []string{"test-name", "test-role"}
)

type infoTest struct {
	in             *ec2.Instance
	expected       *InstanceInfo
	expectedParsed string
}

type headerTest struct {
	in       *InstanceInfo
	expected string
}

var infoTests = []infoTest{
	// normal case
	{
		in: &ec2.Instance{
			InstanceId:       aws.String(dummyInstanceID),
			PrivateIpAddress: aws.String(dummyPrivateIPAddress),
			PublicIpAddress:  aws.String(dummyPublicIPAddress),
			InstanceType:     aws.String(dummyInstanceType),
			State: &ec2.InstanceState{
				Code: aws.Int64(dummyStateCode),
				Name: aws.String(dummyStateName),
			},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(dummyTagKeys[0]),
					Value: aws.String(dummyTagValues[0]),
				},
				{
					Key:   aws.String(dummyTagKeys[1]),
					Value: aws.String(dummyTagValues[1]),
				},
			},
		},
		expected: &InstanceInfo{
			InstanceID:       dummyInstanceID,
			PrivateIPAddress: dummyPrivateIPAddress,
			PublicIPAddress:  dummyPublicIPAddress,
			InstanceType:     dummyInstanceType,
			StateName:        dummyStateName,
			Tags: map[string]string{
				dummyTagKeys[0]: dummyTagValues[0],
				dummyTagKeys[1]: dummyTagValues[1],
			},
		},
		expectedParsed: fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s=%s,%s=%s",
			dummyInstanceID,
			dummyPrivateIPAddress,
			dummyPublicIPAddress,
			dummyInstanceType,
			dummyStateName,
			dummyTagKeys[0],
			dummyTagValues[0],
			dummyTagKeys[1],
			dummyTagValues[1],
		),
	},
	{
		// public ip is nil
		in: &ec2.Instance{
			InstanceId:       aws.String(dummyInstanceID),
			PrivateIpAddress: aws.String(dummyPrivateIPAddress),
			PublicIpAddress:  nil,
			InstanceType:     aws.String(dummyInstanceType),
			State: &ec2.InstanceState{
				Code: aws.Int64(dummyStateCode),
				Name: aws.String(dummyStateName),
			},
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(dummyTagKeys[0]),
					Value: aws.String(dummyTagValues[0]),
				},
			},
		},
		expected: &InstanceInfo{
			InstanceID:       dummyInstanceID,
			PrivateIPAddress: dummyPrivateIPAddress,
			PublicIPAddress:  "UNDEFINED",
			InstanceType:     dummyInstanceType,
			StateName:        dummyStateName,
			Tags: map[string]string{
				dummyTagKeys[0]: dummyTagValues[0],
			},
		},
		expectedParsed: fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s=%s",
			dummyInstanceID,
			dummyPrivateIPAddress,
			"UNDEFINED",
			dummyInstanceType,
			dummyStateName,
			dummyTagKeys[0],
			dummyTagValues[0],
		),
	},
	{
		// tags are empty
		in: &ec2.Instance{
			InstanceId:       aws.String(dummyInstanceID),
			PrivateIpAddress: aws.String(dummyPrivateIPAddress),
			PublicIpAddress:  aws.String(dummyPublicIPAddress),
			InstanceType:     aws.String(dummyInstanceType),
			State: &ec2.InstanceState{
				Code: aws.Int64(dummyStateCode),
				Name: aws.String(dummyStateName),
			},
			Tags: []*ec2.Tag{},
		},
		expected: &InstanceInfo{
			InstanceID:       dummyInstanceID,
			PrivateIPAddress: dummyPrivateIPAddress,
			PublicIPAddress:  dummyPublicIPAddress,
			InstanceType:     dummyInstanceType,
			StateName:        dummyStateName,
			Tags:             map[string]string{},
		},
		expectedParsed: fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t",
			dummyInstanceID,
			dummyPrivateIPAddress,
			dummyPublicIPAddress,
			dummyInstanceType,
			dummyStateName,
		),
	},
}

var headerTests = []headerTest{
	{
		in: &InstanceInfo{
			InstanceID:       dummyInstanceID,
			PrivateIPAddress: dummyPrivateIPAddress,
			PublicIPAddress:  dummyPublicIPAddress,
			InstanceType:     dummyInstanceType,
			StateName:        dummyStateName,
			Tags: map[string]string{
				dummyTagKeys[0]: dummyTagValues[0],
				dummyTagKeys[1]: dummyTagValues[1],
			},
		},
		expected: "INSTANCE_ID\tPRIVATE_IP\tPUBLIC_IP\tTYPE\tSTATE\tTAGS\n",
	},
}

func TestNewInstanceInfo(t *testing.T) {
	for _, it := range infoTests {
		out, err := NewInstanceInfo(it.in)
		if err != nil {
			t.Errorf("occurred error: %v, in: %#v", err, it.in)
		}
		if !compare(out, it.expected) {
			t.Errorf("expected: %#v, but out: %#v", it.expected, out)
		}
	}
}

func compare(out *InstanceInfo, expected *InstanceInfo) bool {
	if out.InstanceID != expected.InstanceID {
		return false
	}
	if out.PrivateIPAddress != expected.PrivateIPAddress {
		return false
	}
	if out.PublicIPAddress != expected.PublicIPAddress {
		return false
	}
	if out.InstanceType != expected.InstanceType {
		return false
	}
	if out.StateName != expected.StateName {
		return false
	}

	if len(out.Tags) != len(expected.Tags) {
		return false
	}
	if len(out.Tags) > 0 {
		for k, v := range out.Tags {
			if v != expected.Tags[k] {
				return false
			}
		}
	}

	return true
}

func TestNewInstanceInfoByNil(t *testing.T) {
	out, err := NewInstanceInfo(nil)
	if err == nil {
		t.Errorf("not occurred expected error. out: %#v", out)
	}
}

func TestParseRow(t *testing.T) {
	for _, it := range infoTests {
		out := it.expected.ParseRow(false)
		if out != it.expectedParsed {
			t.Errorf("expected: [%s], but out: [%s]", it.expectedParsed, out)
		}
	}
}

func TestPrintHeader(t *testing.T) {
	var _ = (*InstanceInfo).printHeader
	for _, it := range headerTests {
		out := it.in.printHeader()
		if (out != it.expected) {
			t.Errorf("expected: %#v, but out: %#v", it.expected, out)
		}
	}
}

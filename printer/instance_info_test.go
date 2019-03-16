package printer_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	. "github.com/goldeneggg/lsec2/printer"
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
	expectedValues []string
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
		expectedValues: []string{
			dummyInstanceID,
			dummyPrivateIPAddress,
			dummyPublicIPAddress,
			dummyInstanceType,
			dummyStateName,
			dummyTagKeys[0] + "=" + dummyTagValues[0] + "," + dummyTagKeys[1] + "=" + dummyTagValues[1],
		},
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
		expectedValues: []string{
			dummyInstanceID,
			dummyPrivateIPAddress,
			dummyPublicIPAddress,
			dummyInstanceType,
			dummyStateName,
			dummyTagKeys[0] + "=" + dummyTagValues[0],
		},
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
		expectedValues: []string{
			dummyInstanceID,
			dummyPrivateIPAddress,
			dummyPublicIPAddress,
			dummyInstanceType,
			dummyStateName,
			"",
		},
	},
}

func TestNewInstanceInfo(t *testing.T) {
	for _, it := range infoTests {
		out := NewInstanceInfo(it.in)

		if !compare(out, it.expected) {
			t.Errorf("expected: %#v, but out: %#v", it.expected, out)
		}
	}
}

func TestValues(t *testing.T) {
	for _, it := range infoTests {
		out := it.expected.Values(false)
		if len(out) != len(it.expectedValues) {
			t.Errorf("expected: [%#v], but out: [%#v]", it.expectedValues, out)
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

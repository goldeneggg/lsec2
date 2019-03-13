package printer_test

import (
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"github.com/goldeneggg/lsec2/awsec2"
	. "github.com/goldeneggg/lsec2/printer"
)

var (
	dummyInstances = []*ec2.Instance{
		&ec2.Instance{
			InstanceId:       aws.String("i-xxxxxxxx91"),
			PrivateIpAddress: aws.String("192.168.56.101"),
			PublicIpAddress:  aws.String("54.0.0.1"),
			InstanceType:     aws.String("t3.large"),
			State: &ec2.InstanceState{
				Name: aws.String("running"),
			},
		},
		&ec2.Instance{
			InstanceId:       aws.String("i-xxxxxxxx99"),
			PrivateIpAddress: aws.String("192.168.56.199"),
			PublicIpAddress:  aws.String("54.0.0.9"),
			InstanceType:     aws.String("t4.large"),
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
				Instances: []*ec2.Instance{dummyInstances[0]},
			},
			&ec2.Reservation{
				Instances: []*ec2.Instance{dummyInstances[1]},
			},
		},
	}

	return instancesOutput, nil
}

func TestNewPrinter(t *testing.T) {
	dummyFile := os.NewFile(uintptr(3), "printer_test.go.out")

	cases := []struct {
		maybeWriter io.Writer
		exp         *Printer
	}{
		{
			maybeWriter: nil,
			exp: &Printer{
				Writer: os.Stdout,
			},
		},
		{
			maybeWriter: dummyFile,
			exp: &Printer{
				Writer: dummyFile,
			},
		},
	}

	for _, c := range cases {
		printer := NewPrinter(c.maybeWriter)

		if printer.Writer != c.exp.Writer {
			t.Errorf("expected: %#v, but exp: %#v", c.exp.Writer, printer.Writer)
		}
	}
}

func TestPrintAll(t *testing.T) {
	client := &awsec2.Client{
		EC2API:    &mockEC2API{},
		Region:    "ap-northeast-1",
		StateName: "running",
		Tags:      []string{"Role=roleA,roleB", "Service=serviceX"},
	}

	cases := []struct {
		printer *Printer
	}{
		{
			printer: &Printer{
				Writer: os.Stdout,
			},
		},
	}

	for _, c := range cases {
		err := c.printer.PrintAll(client)

		if err != nil {
			t.Errorf("error occured. err: %#v", err)
		}
	}
}

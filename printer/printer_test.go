package printer_test

import (
	"io"
	"io/ioutil"
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
		pr := NewPrinter(c.maybeWriter)

		if pr.Delimiter != "\t" {
			t.Errorf("expected: %#v, but exp: %#v", c.exp.Writer, pr.Writer)
		}

		if pr.Writer != c.exp.Writer {
			t.Errorf("expected: %#v, but exp: %#v", c.exp.Writer, pr.Writer)
		}
	}
}

func TestPrintAll(t *testing.T) {
	client := &awsec2.Client{
		EC2API:    &mockEC2API{},
		StateName: "running",
	}

	tmp, err := ioutil.TempFile("", "test_lsec2_printer")
	if err != nil {
		t.Errorf("failed to ioutil.TempFile. err: %#v", err)
	}
	defer os.Remove(tmp.Name())

	cases := []struct {
		pr *Printer
	}{
		{
			pr: &Printer{
				Writer:        os.Stdout,
				PrintHeader:   true,
				OnlyPrivateIP: false,
				WithColor:     true,
				Delimiter:     "\t",
			},
		},
		{
			pr: &Printer{
				Writer:        tmp,
				PrintHeader:   false,
				OnlyPrivateIP: false,
				WithColor:     false,
				Delimiter:     ",",
			},
		},
		{
			pr: &Printer{
				Writer:        tmp,
				PrintHeader:   false,
				OnlyPrivateIP: true,
				WithColor:     false,
			},
		},
	}

	for _, c := range cases {
		err := c.pr.PrintAll(client)

		if err != nil {
			t.Errorf("error occured. err: %#v", err)
		}
	}
}

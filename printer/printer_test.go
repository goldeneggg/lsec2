package printer_test

import (
	"io"
	"io/ioutil"
	"os"
	"testing"
	"text/tabwriter"

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
		delim     string
		header    bool
		onlyPvtIP bool
		withColor bool
		w         io.Writer
		expected  *Printer
	}{
		{
			delim:     "",
			header:    true,
			onlyPvtIP: false,
			withColor: true,
			w:         nil,
			expected: &Printer{
				Delimiter:     "\t",
				PrintHeader:   true,
				OnlyPrivateIP: false,
				WithColor:     true,
			},
		},
		{
			delim:     ",",
			header:    false,
			onlyPvtIP: true,
			withColor: false,
			w:         dummyFile,
			expected: &Printer{
				Writer:        dummyFile,
				Delimiter:     ",",
				PrintHeader:   false,
				OnlyPrivateIP: true,
				WithColor:     false,
			},
		},
	}

	for _, c := range cases {
		pr := NewPrinter(c.delim, c.header, c.onlyPvtIP, c.withColor, c.w)

		if pr.Delimiter != c.expected.Delimiter {
			t.Errorf("expected: %#v, but actual: %#v", c.expected.Delimiter, pr.Delimiter)
		}
		if pr.PrintHeader != c.expected.PrintHeader {
			t.Errorf("expected: %#v, but actual: %#v", c.expected.PrintHeader, pr.PrintHeader)
		}
		if pr.OnlyPrivateIP != c.expected.OnlyPrivateIP {
			t.Errorf("expected: %#v, but actual: %#v", c.expected.OnlyPrivateIP, pr.OnlyPrivateIP)
		}
		if pr.WithColor != c.expected.WithColor {
			t.Errorf("expected: %#v, but actual: %#v", c.expected.WithColor, pr.WithColor)
		}
		if c.w == nil {
			if _, ok := pr.Writer.(*tabwriter.Writer); !ok {
				t.Errorf("expected: *tabwriter.Writer, but actual: %#v", pr.Writer)
			}
		} else {
			if pr.Writer != c.expected.Writer {
				t.Errorf("expected: %#v, but actual: %#v", c.expected.Writer, pr.Writer)
			}
		}
	}
}

func TestPrintAll(t *testing.T) {
	tmp, err := ioutil.TempFile("", "test_lsec2_printer")
	if err != nil {
		t.Errorf("failed to ioutil.TempFile. err: %#v", err)
	}
	defer os.Remove(tmp.Name())

	client := &awsec2.Client{
		EC2API:    &mockEC2API{},
		StateName: "running",
	}
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
		if err := c.pr.PrintAll(client); err != nil {
			t.Errorf("error occured. err: %#v", err)
		}
	}
}

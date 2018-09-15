package awsec2

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// Client is options definition of print
// client attributes
type Client struct {
	PrintHeader   bool
	OnlyPrivateIP bool
	Region        string
	Tags          []string
	WithColor     bool
	StateName     string
}

// Print is print method for aws ec2 instances
// print information of aws ec2 instances
func (client *Client) Print() error {
	infos, err := client.buildInfos()
	if err != nil {
		return fmt.Errorf("buildInfos error: %v", err)
	}

	client.printInfos(infos)

	return nil
}

func (client *Client) buildInfos() ([]*InstanceInfo, error) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(client.Region)})
	if err != nil {
		return nil, fmt.Errorf("aws new session error: %v", err)
	}

	svc := ec2.New(sess)

	output, err := svc.DescribeInstances(client.filterParams())
	if err != nil {
		return nil, fmt.Errorf("aws describe instances error: %v", err)
	}

	var infos []*InstanceInfo

	for _, reservation := range output.Reservations {
		for _, instance := range reservation.Instances {
			info, err := NewInstanceInfo(instance)
			if err != nil {
				return nil, err
			}
			infos = append(infos, info)
		}
	}

	return infos, nil
}

func (client *Client) filterParams() *ec2.DescribeInstancesInput {
	var filters []*ec2.Filter

	// client.tags is separated key-value pair by "=", and values are separated by ","(comma)
	// ex. "Name=Value"
	// ex. "Name=Value1,Value2"
	for _, tag := range client.Tags {
		tagNameValue := strings.Split(tag, tagPairSeparator)
		name := aws.String(tagFilterPrefix + tagNameValue[0])
		values := make([]*string, 0, 3)
		for _, value := range strings.Split(tagNameValue[1], tagValueSeparator) {
			values = append(values, aws.String(value))
		}

		tagFilter := &ec2.Filter{
			Name:   name,
			Values: values,
		}
		filters = append(filters, tagFilter)
	}

	if len(filters) == 0 {
		return nil
	}

	return &ec2.DescribeInstancesInput{Filters: filters}
}

func (client *Client) printInfos(infos []*InstanceInfo) {
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	if client.PrintHeader {
		w.Write([]byte(infos[0].printHeader()))
	}

	for _, info := range infos {
		if len(client.StateName) == 0 || client.StateName == info.StateName {
			if client.OnlyPrivateIP {
				w.Write([]byte(fmt.Sprintf("%s\n", info.PrivateIPAddress)))
			} else {
				w.Write([]byte(info.printRow(client.WithColor)))
			}
		}
	}
	w.Flush()
}

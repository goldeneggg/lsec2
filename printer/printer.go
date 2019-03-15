package printer

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/goldeneggg/lsec2/awsec2"
)

const (
	defaultDelimiter = "\t"

	tabMinWidth = 0
	tabTabWidth = 4
	tabPadding  = 4
	tabPadChar  = ' '
	tabFlags    = 0
)

/*
type flushableWriter interface {
	io.Writer
	Flush() error
}
*/

// Printer is options definition of print
type Printer struct {
	io.Writer
	Delimiter     string
	PrintHeader   bool
	OnlyPrivateIP bool
	WithColor     bool
}

// NewPrinter returns a new Printer
func NewPrinter(delim string, header, onlyPvtIP, withColor bool, w interface{}) *Printer {
	pr := new(Printer)

	if delim == "" {
		delim = defaultDelimiter
	}
	pr.Delimiter = delim
	pr.PrintHeader = header
	pr.OnlyPrivateIP = onlyPvtIP
	pr.WithColor = withColor

	if writer, ok := w.(io.Writer); ok {
		pr.Writer = writer
	} else {
		pr.Writer = defaultWriter(delim)
	}

	return pr
}

func defaultWriter(delim string) io.Writer {
	if delim == defaultDelimiter {
		return tabwriter.NewWriter(
			os.Stdout,
			tabMinWidth,
			tabTabWidth,
			tabPadding,
			tabPadChar,
			tabFlags,
		)
	}

	return os.Stdout
}

// PrintAll prints information all of aws ec2 instances
func (pr *Printer) PrintAll(client *awsec2.Client) error {
	instances, err := client.EC2Instances()
	if err != nil {
		return fmt.Errorf("get EC2 Instances error: %v", err)
	}

	defer pr.flushIfFlushable()

	if pr.PrintHeader {
		if err := pr.printHeader(); err != nil {
			return fmt.Errorf("print header error: %v", err)
		}
	}

	for _, inst := range instances {
		if err := pr.printInstance(client, inst); err != nil {
			return fmt.Errorf("print instance error: %v", err)
		}
	}

	return nil
}

func (pr *Printer) flushIfFlushable() {
	//if fw, ok := pr.Writer.(flushableWriter); ok {
	if fw, ok := pr.Writer.(interface{ Flush() error }); ok {
		fw.Flush()
	}
}

func (pr *Printer) printHeader() error {
	return pr.printArray(NewInstanceInfo(nil).Headers())
}

func (pr *Printer) printInstance(client *awsec2.Client, inst *ec2.Instance) error {
	var (
		ii  *InstanceInfo
		err error
	)

	ii = NewInstanceInfo(inst)
	if len(client.StateName) == 0 || client.StateName == ii.StateName {
		if pr.OnlyPrivateIP {
			err = pr.printArray([]string{ii.PrivateIPAddress})
		} else {
			err = pr.printArray(ii.Values(pr.WithColor))
		}
	}

	return err
}

func (pr *Printer) printArray(sArr []string) error {
	if _, err := fmt.Fprintln(pr.Writer, strings.Join(sArr, pr.Delimiter)); err != nil {
		return err
	}

	return nil
}

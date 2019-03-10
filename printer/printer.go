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

	tabWriterMinWidth = 0
	tabWriterTabWidth = 4
	tabWriterPadding  = 4
	tabWriterPadChar  = ' '
	tabWriterFlags    = 0
)

var (
	defaultWriter = os.Stdout
)

type flushableWriter interface {
	io.Writer
	Flush() error
}

// Printer is options definition of print
type Printer struct {
	io.Writer
	PrintHeader   bool
	OnlyPrivateIP bool
	WithColor     bool
	Delimeter     string
	ColDef        string
}

// NewPrinter returns a new Printer
func NewPrinter(maybeWriter interface{}) *Printer {
	printer := new(Printer)

	printer.Delimeter = defaultDelimiter

	if writer, ok := maybeWriter.(io.Writer); ok {
		printer.Writer = writer
	} else {
		printer.Writer = defaultWriter
	}

	return printer
}

// Print is print method for aws ec2 instances
// print information of aws ec2 instances
func (printer *Printer) PrintAll(client *awsec2.Client) error {
	instances, err := client.EC2Instances()
	if err != nil {
		return fmt.Errorf("get EC2 Instances error: %v", err)
	}

	printer.wrapWriterIfDefault()
	defer printer.flushIfFlushable()

	if printer.PrintHeader {
		printer.printHeader()
	}

	for _, inst := range instances {
		if err := printer.printInstance(client, inst); err != nil {
			return fmt.Errorf("print instance error: %v", err)
		}
	}

	return nil
}

func (printer *Printer) wrapWriterIfDefault() {
	if printer.Writer == defaultWriter {
		if printer.Delimeter == defaultDelimiter {
			printer.Writer = tabwriter.NewWriter(
				printer.Writer,
				tabWriterMinWidth,
				tabWriterTabWidth,
				tabWriterPadding,
				tabWriterPadChar,
				tabWriterFlags,
			)
		}
	}
}

func (printer *Printer) flushIfFlushable() {
	// FIXME:
	// printer.Writer.(interface{Flush() error}) とも書けそう。どちらがパフォーマンスが良いか？
	if fw, ok := printer.Writer.(flushableWriter); ok {
		fw.Flush()
	}
}

func (printer *Printer) printHeader() {
	printer.printArray(NewInstanceInfo(nil).Headers())
}

func (printer *Printer) printInstance(client *awsec2.Client, inst *ec2.Instance) error {
	var (
		ii  *InstanceInfo
		err error
	)

	ii = NewInstanceInfo(inst)
	if len(client.StateName) == 0 || client.StateName == ii.StateName {
		if printer.OnlyPrivateIP {
			err = printer.printArray([]string{ii.PrivateIPAddress})
		} else {
			err = printer.printArray(ii.Values(printer.WithColor))
		}
	}

	return err
}

func (printer *Printer) printArray(sArr []string) error {
	if _, err := fmt.Fprintln(printer.Writer, strings.Join(sArr, printer.Delimeter)); err != nil {
		return err
	}

	return nil
}

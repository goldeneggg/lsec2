package coldef

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
)

const (
	ColDefFileName = "coldef"

	ArraySign = "ARRAY"

	DefaultYAML = `
- InstanceId
- PrivateIpAddress
- PublicIpAddress
- InstanceType
- State: Name
- Tags
`
)

var (
	DefaultDir = filepath.Join("$HOME", ".lsec2")
)

type Selector interface {
	SelectInstanceInfo(*ec2.Instance) error
}

type colDef struct {
	Columns []interface{}
}

type Selected struct {
	Name  string
	Value interface{}
}

type DefaultSelector struct {
	*colDef
	Selection []*Selected
}

// NewSelectorFromYAMLFile read a yaml file named by dir argument
func NewSelectorFromYAMLFile(dir string) (Selector, error) {
	var cd colDef

	viper.SetConfigName(ColDefFileName)
	viper.SetConfigType("yml")
	viper.AddConfigPath(dir)

	if err := viper.ReadInConfig(); err != nil {
		// if yaml not found, use default yaml definition
		viper.ReadConfig(bytes.NewBuffer([]byte(DefaultYAML)))
	}

	err := viper.Unmarshal(&cd)
	return &DefaultSelector{colDef: &cd}, err
}

func (ds *DefaultSelector) SelectInstanceInfo(instance *ec2.Instance) error {
	instValue := reflect.ValueOf(instance)
	log.Printf("--- instValue: %+v", instValue)

	for _, column := range ds.colDef.Columns {
		log.Printf("--- column: %+v", column)
		sel := &Selected{}
		err := sel.setSelected(column, instValue)
		if err != nil {
			return err
		}

		log.Printf("selected: %+v", sel)
		ds.Selection = append(ds.Selection, sel)
	}

	return nil
}

func (sel *Selected) setSelected(column interface{}, instValue reflect.Value) (err error) {
	instance := reflect.Indirect(instValue)

	switch instance.Kind() {
	case reflect.Array, reflect.Slice:
		for index := 0; index < instance.Len(); index++ {
			err = sel.setSelected(column, instance.Index(index))
			if err != nil {
				return
			}
		}
		// 何かのスライス型のValue、例えば []*ec2.Tags とか
		// return sel.parseArray(column, indirected)
	default:
		switch typedColumn := column.(type) {
		case string:
			// e.g. "- InstanceId"
			err = sel.processStringColumn(typedColumn, instance)
		case map[interface{}]interface{}:
			// e.g. "- State: Name"
			vals, isTags := typedColumn["Tags"]
			if isTags {
				aVals := vals.([]interface{})
				cols := make([]string, len(aVals))
				for i, v := range aVals {
					cols[i] = v.(string)
				}
				err = sel.processTags(cols, instance)
			} else {
				err = sel.processMapColumn(typedColumn, instance)
			}
		case []interface{}:
			// e.g.
			// "- Tags"
			// "  - Key"
			// "  - Value"
			log.Printf("@@@ array typedColumn: %+v, instance: %#v", typedColumn, instance)
			for _, next := range typedColumn {
				err = sel.setSelected(next, instance)
				if err != nil {
					return
				}
			}
		default:
			// FIXME
			log.Printf("----- default typedColumn: %#v", typedColumn)
		}
	}

	return
}

/*
func (sel *Selected) parseArray(columnIntf interface{}, value reflect.Value) error {
	log.Printf("## parseArray columnIntf: %+v, value: %+v", columnIntf, value)

	for index := 0; index < value.Len(); index++ {
		log.Printf("## parseArray call parseObject: %+v, value: %+v", columnIntf, value.Index(index))
		return sel.fromObject(columnIntf, value.Index(index))
	}

	return nil
}
*/

func (sel *Selected) processStringColumn(typedColumn string, value reflect.Value) error {
	str, err := stringFieldValue(typedColumn, value)
	if err != nil {
		return err
	}

	sel.buildName(typedColumn)
	sel.setValue(str)
	return nil
}

func (sel *Selected) processMapColumn(typedColumn map[interface{}]interface{}, value reflect.Value) error {
	// e.g. "- State: Name"
	var current string
	var next interface{}
	for c, n := range typedColumn {
		current, _ = c.(string)
		next = n
		break
	}
	log.Printf("------- map column current: %s", current)

	sel.buildName(current)

	return sel.setSelected(next, value.FieldByName(current))
}

func (sel *Selected) processTags(cols []string, value reflect.Value) error {
	tags := value.FieldByName("Tags")

	var tag reflect.Value
	var k, v string
	val := make(map[string]string)

	for index := 0; index < tags.Len(); index++ {
		tag = reflect.Indirect(tags.Index(index))
		for i := 0; i < tag.NumField(); i++ {
			switch tag.Field(i).Kind() {
			case reflect.Ptr:
				if i%2 == 0 {
					v = tag.Field(i).Elem().String()
				} else {
					k = tag.Field(i).Elem().String()
				}
				val[k] = v
			}
		}
	}

	sel.buildName("Tags")
	sel.setValue(val)

	return nil
}

func (sel *Selected) buildName(s string) {
	if sel.Name == "" {
		sel.Name = s
	} else {
		sel.Name = sel.Name + "." + s
	}
}

func (sel *Selected) setValue(v interface{}) {
	sel.Value = v
}

func stringFieldValue(fieldName string, value reflect.Value) (s string, err error) {
	v := value.FieldByName(fieldName)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.String:
		s = v.String()
		log.Printf("----- string: %s, field: %s", s, fieldName)
	case reflect.Int64:
		s = strconv.FormatInt(v.Int(), 10)
		log.Printf("----- int64: %s, field: %s", s, fieldName)
	case reflect.Invalid:
		log.Printf("SKIP INVALID. field: %s, casted: %+v", fieldName, v)
	default:
		err = fmt.Errorf("field: %s does not have string value. %s", fieldName, v.String())
	}

	return s, err
}

func Hoge() {
	//cd, err := NewColDefFromYAMLFile(filepath.Join("$HOME", ".lsec2"))
	ds, err := NewSelectorFromYAMLFile(filepath.Join(".", "tmp"))
	if err != nil {
		log.Printf("@@ err: %v", err)
	}

	log.Printf("@@ ds: %+v", ds)

	instance := &ec2.Instance{}
	instance.InstanceId = aws.String("instance-id")
	instance.PublicIpAddress = aws.String("1.2.3.4")

	instance.State = &ec2.InstanceState{
		Code: aws.Int64(1234),
		Name: aws.String("running"),
	}
	instance.CapacityReservationId = aws.String("capacity-reservation-id")
	instance.NetworkInterfaces = []*ec2.InstanceNetworkInterface{
		{
			Association: &ec2.InstanceNetworkInterfaceAssociation{PublicDnsName: aws.String("xyz-dns")},
			Groups: []*ec2.GroupIdentifier{
				{GroupId: aws.String("xxxgroup111")},
				{GroupId: aws.String("xxxgroup222")},
			},
		},
	}
	instance.Tags = []*ec2.Tag{
		{Key: aws.String("tag1"), Value: aws.String("value1")},
		{Key: aws.String("tag2"), Value: aws.String("value2")},
	}

	err = ds.SelectInstanceInfo(instance)
	if err != nil {
		log.Fatalf("!!!!! Select error: %v", err)
	}
}

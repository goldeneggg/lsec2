package coldef

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
	"gopkg.in/yaml.v2"
)

const (
	yamlFile = "coldef.yml"
)

var (
	DefaultPath = os.Getenv("HOME") + string(os.PathSeparator) + ".lsec2" + string(os.PathSeparator) + yamlFile

	DefaultYAML = `
- InstanceId
- PrivateIpAddress
- PublicIpAddress
- InstanceType
- State.Name
- Tags
`

	hogeYAML = `
- InstanceId
- PrivateIpAddress
- PublicIpAddress
- InstanceType
- State.Name
- Licenses:
  - LicenseConfigurationArn21
    LicenseConfigurationArn22
- Tags
`
)

// Read read a yaml file named by path argument
func Read(path string) (unmarshals []interface{}, err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Unmarshal([]byte(DefaultYAML))
	}

	f, err := os.Open(path)
	if err != nil {
		return
	}

	d, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}

	return Unmarshal(d)
}

func Unmarshal(data []byte) (unmarshals []interface{}, err error) {
	err = yaml.UnmarshalStrict(data, &unmarshals)

	return
}

type ColDef struct {
	Columns []interface{}
}

func NewColDef(unmarshals []interface{}) (*ColDef, error) {
	cd := &ColDef{Columns: make([]interface{}, 0, 10)}

	for _, um := range unmarshals {
		switch v := um.(type) {
		case string:
			cd.Columns = append(cd.Columns, newItem(v))
		case map[interface{}]interface{}:
			if ri, err := newRepeatedItem(v); err == nil {
				cd.Columns = append(cd.Columns, ri)
			} else {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("%v is invalid type", um)
		}
	}

	return cd, nil
}

type Item struct {
	Receivers []string
	Name      string
}

func newItem(s string) Item {
	a := strings.Split(s, ".")
	if len(a) == 0 {
		return Item{Receivers: nil, Name: s}
	}
	return Item{Receivers: a[0 : len(a)-1], Name: a[len(a)-1]}
}

func (i Item) EmptyReceivers() bool {
	return len(i.Receivers) == 0
}

type RepeatedItem struct {
	Root  string
	Items []Item
}

func newRepeatedItem(m map[interface{}]interface{}) (RepeatedItem, error) {
	var ri RepeatedItem

	if len(m) > 1 {
		return ri, fmt.Errorf("%v is not a map of size one", m)
	}

	var (
		root  string
		items []Item
		ok    bool
	)

	for key, value := range m {
		root, ok = key.(string)
		if !ok {
			return ri, fmt.Errorf("key %v is not a string", key)
		}

		arrIntf, ok := value.([]interface{})
		if !ok {
			return ri, fmt.Errorf("value %v is not a interface array", value)
		}

		for _, intf := range arrIntf {
			s, ok := intf.(string)
			if !ok {
				return ri, fmt.Errorf("intf %v is not a string", intf)
			}
			items = append(items, newItem(s))
		}

		break
	}

	return RepeatedItem{Root: root, Items: items}, nil
}

func Hoge() {
	unmarshals, err := Unmarshal([]byte(hogeYAML))
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- unmarshals: \n%#v\n\n", unmarshals)

	cd, err := NewColDef(unmarshals)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- coldef: \n%#v\n\n", cd)

	d, err := yaml.Marshal(&unmarshals)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("--- unmarshals dump:\n%s\n\n", string(d))

	for _, col := range cd.Columns {
		var ok bool

		fmt.Printf(" col: %#v\n", col)

		switch t := col.(type) {
		case Item:
			if t.EmptyReceivers() {
				fmt.Println("EmptyReceivers")
				_, ok = reflect.TypeOf(ec2.Instance{}).FieldByName(t.Name)
			} else {
				fmt.Println("Non EmptyReceivers")
				var st reflect.StructField

				for idx, r := range t.Receivers {
					if idx == 0 {
						fmt.Printf(" --- from EC2, st %#v, receiver: %#v, idx: %d\n", st, r, idx)
						st, ok = reflect.TypeOf(ec2.Instance{}).FieldByName(r)
					} else {
						fmt.Printf(" --- from st, st %#v, receiver: %#v, idx: %d\n", st, r, idx)
						st, ok = st.Type.FieldByName(r)
					}

					fmt.Printf("  FieldByName ok: %v, st: %#v\n", ok, st)
					fmt.Printf("  ! st.Type %#v\n", st.Type)
					if !ok {
						break
					}
				}

				if ok {
					fmt.Printf(" Type: %#v\n", st.Type)
					st, ok = st.Type.FieldByName(t.Name)
					fmt.Printf(" final st: %#v\n", st)
				}
			}

			if !ok {
				log.Fatalf("error: %v", fmt.Errorf("does not exist in field"))
			}

		case RepeatedItem:
			fmt.Printf(" I am RepeatedItem TODO\n")
		default:
			fmt.Printf(" t: %v\n", t)
		}
	}
	// reflect.TypeOf(ec2.Instance{}).FieldByName("xxx")
	// FieldByName(name string) (StructField, bool)
}

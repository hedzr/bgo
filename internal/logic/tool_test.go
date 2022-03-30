package logic

//nolint:goimports
import (
	"fmt"
	"github.com/hedzr/cmdr"
	"gopkg.in/yaml.v3"
	"reflect"
	"strconv"
	"testing"
)

func TestDottedKeyInYaml(t *testing.T) {
	doc := `---
apiVersion: v1
kind: ServiceAccount
001.key: 11
metadata:
  labels:
    k8s-app: kubernetes-dashboard
    addonmanager.kubernetes.io/mode: Reconcile
  name: kubernetes-dashboard
  namespace: kube-system
`

	m := make(map[string]interface{})
	err := yaml.Unmarshal([]byte(doc), m)
	if err != nil {
		t.Error(err)
	}

	t.Log(m)

	out, err := yaml.Marshal(m)
	if err != nil {
		t.Error(err)
	}
	t.Log("\n", string(out))
}

type (
	SA struct {
		A string
		B int
	}

	SB struct {
		*SA
		OS string
	}
)

func TestEmbedStruct0(t *testing.T) {
	sb := SB{
		SA: &SA{
			A: "aa",
			B: -3,
		},
		OS: "a",
	}

	tgt := new(SB)

	_ = cmdr.CloneViaGob(tgt, sb)
}

func TestEmbedStruct1(t *testing.T) {
	sb := SB{
		SA: &SA{
			A: "aa",
			B: -3,
		},
		OS: "a",
	}

	tgt := new(SB)

	_ = cmdr.Clone(tgt, sb)
}

func TestEmbedStruct2(t *testing.T) {
	oo := new(BgoSettings)

	var name string
	var to = indirect(reflect.ValueOf(oo))
	var toType = indirectType(to.Type())
	for _, field := range deepFields(toType) {
		name = field.Name
		t.Log(name, field.Type, field.Index, field.Anonymous, field.Offset)
	}
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) reflect.Type {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	var fields []reflect.StructField

	if reflectType = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}
	}

	return fields
}

func TestUnquote(t *testing.T) {
	str := "'hello!'"
	t.Log(strconv.Unquote(str))
	str = "`hello!`"
	t.Log(strconv.Unquote(str))
	str = "\"hello!\""
	t.Log(strconv.Unquote(str))

	s, err := strconv.Unquote("You can't unquote a string without quotes")
	fmt.Printf("%q, %v\n", s, err)
	s, err = strconv.Unquote("\"The string must be either double-quoted\"")
	fmt.Printf("%q, %v\n", s, err)
	s, err = strconv.Unquote("`or backquoted.`")
	fmt.Printf("%q, %v\n", s, err)
	s, err = strconv.Unquote("'\u263a'") // single character only allowed in single quotes
	fmt.Printf("%q, %v\n", s, err)
	s, err = strconv.Unquote("'\u2639\u2639'")
	fmt.Printf("%q, %v\n", s, err)
}

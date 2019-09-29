package enki

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type unmarshalTest struct {
	in  string
	out interface{}
}

var unmarshalTests = []unmarshalTest{
	// empty
	{in: `level1:`, out: map[string]interface{}{"level1": map[string]interface{}{}}},
	{in: `level1:level2:`, out: map[string]interface{}{"level1": map[string]interface{}{"level2": map[string]interface{}{}}}},
	// basic types
	{in: `level1:int=1`, out: map[string]interface{}{"level1": map[string]interface{}{"int": 1}}},
	{in: `level1:float=1.1`, out: map[string]interface{}{"level1": map[string]interface{}{"float": 1.1}}},
	{in: `level1:bool`, out: map[string]interface{}{"level1": map[string]interface{}{"bool": true}}},
	{in: `level1:bool=true`, out: map[string]interface{}{"level1": map[string]interface{}{"bool": true}}},
	{in: `level1:bool=false`, out: map[string]interface{}{"level1": map[string]interface{}{"bool": false}}},
	{in: `level1:string="a string"`, out: map[string]interface{}{"level1": map[string]interface{}{"string": "a string"}}},
	{in: `level1:string="a:string"`, out: map[string]interface{}{"level1": map[string]interface{}{"string": "a:string"}}},
	{in: `level1:string=a string`, out: map[string]interface{}{"level1": map[string]interface{}{"string": "a string"}}},
	// lists
	{in: `level1:list=[string]`, out: map[string]interface{}{"level1": map[string]interface{}{"list": []interface{}{"string"}}}},
	{in: `level1:list=[string,1]`, out: map[string]interface{}{"level1": map[string]interface{}{"list": []interface{}{"string", 1}}}},
	{in: `level1:list=[string,1,true]`, out: map[string]interface{}{"level1": map[string]interface{}{"list": []interface{}{"string", 1, true}}}},
	{in: `level1:list=[string,1,false]`, out: map[string]interface{}{"level1": map[string]interface{}{"list": []interface{}{"string", 1, false}}}},
	// multiple levels
	{in: `level1:level2:string="a string"`, out: map[string]interface{}{"level1": map[string]interface{}{"level2": map[string]interface{}{"string": "a string"}}}},
	// errors
	{in: `level1:string=a=string`, out: map[string]interface{}{"level1": map[string]interface{}{"string": "a=string"}}},
	// multiple fields
	{in: `level1:string="a string",bool=true,int=1`, out: map[string]interface{}{"level1": map[string]interface{}{"string": "a string", "bool": true, "int": 1}}},
	// nested node
	{in: `level1:string="a string",nested:bool`, out: map[string]interface{}{"level1": map[string]interface{}{"string": "a string", "nested": map[string]interface{}{"bool": true}}}},
}

func TestUnmarshal(t *testing.T) {
	for _, k := range unmarshalTests {
		g, err := Unmarshal([]byte(k.in))
		if err != nil {
			t.Error(err.Error())
		}
		if !cmp.Equal(k.out, g) {
			t.Errorf("Unmarshal([]byte(\"%s\")) = got: %v want: %v", k.in, g, k.out)
		}
	}
}

type hasSuffxRemoveTest struct {
	in     string
	suffix string
	out    string
}

var hasSuffxRemoveTests = []hasSuffxRemoveTest{
	{in: "test:", suffix: ":", out: "test"},
	{in: "test::", suffix: "::", out: "test"},
	{in: "test,", suffix: ",", out: "test"},
	{in: "test=", suffix: "=", out: "test"},
}

func TestHasSuffixRemove(t *testing.T) {
	for _, k := range hasSuffxRemoveTests {
		g := hasSuffixRemove(k.in, k.suffix)
		if g != k.out {
			t.Errorf("hasSuffixRemove(\"%s\",\"%s\")) = got: '%s' want: '%s'\n", k.in, k.suffix, g, k.out)
		}
	}
}

package enki

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type parseTest struct {
	in           []string
	out          interface{}
	token        string
	namespace    string
	includeInput bool
}

var parseTests = []parseTest{
	{
		in:           []string{`# +!!:level1:level2:note=["this,that"]`},
		token:        `+!!:`,
		namespace:    `test`,
		includeInput: false,
		out: map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"note": []interface{}{"this,that"},
				},
			},
		},
	},
	{
		in:           []string{`# +!!:level1`},
		token:        `+!!:`,
		namespace:    `test`,
		includeInput: true,
		out: map[string]interface{}{
			"_input": map[string]interface{}{
				"test": map[string]interface{}{
					"content": "IyArISE6bGV2ZWwxCg==",
					"src": map[string]interface{}{
						"0": "IyArISE6bGV2ZWwx",
					},
				},
			},
			"level1": true,
		},
	},
	{
		in:           []string{`# +!!:level1=`, `text`},
		token:        `+!!:`,
		namespace:    `test`,
		includeInput: false,
		out: map[string]interface{}{
			"level1": "dGV4dAo=",
		},
	},
	{
		in:           []string{`# +!!:level1:src=`, `text`, `# +!!:level1:test`},
		token:        `+!!:`,
		namespace:    `test`,
		includeInput: false,
		out: map[string]interface{}{
			"level1": map[string]interface{}{
				"src":  "dGV4dAo=",
				"test": true,
			},
		},
	},
	{
		in:           []string{`# +!!:level1:src=`, `text`, `# +!!:level2:test`},
		token:        `+!!:`,
		namespace:    `test`,
		includeInput: false,
		out: map[string]interface{}{
			"level1": map[string]interface{}{
				"src": "dGV4dAo=",
			},
			"level2": map[string]interface{}{
				"test": true,
			},
		},
	},
	{
		in:           []string{`# +!!:level1:level2:test="pass"`},
		token:        `+!!:`,
		namespace:    `test`,
		includeInput: false,
		out:          makeMap([]string{"level1", "level2"}, map[string]interface{}{"test": "pass"}),
	},
}

func TestParse(t *testing.T) {
	for _, k := range parseTests {
		c := ParserConfig{
			Token:        k.token,
			Namespace:    k.namespace,
			IncludeInput: k.includeInput,
		}

		p := NewParser(strings.NewReader(strings.Join(k.in, "\n")), c)

		if err := p.Parse(); err != nil {
			t.Error(err.Error())
		}

		if !cmp.Equal(k.out, p.Data) {
			t.Errorf("---\nParse:  %s\nGot:\t%v\nWant:\t%v\n", strings.Join(k.in, "\n"), p.Data, k.out)
		}
	}
}

func makeMap(s []string, e interface{}) map[string]interface{} {
	m := make([]map[string]interface{}, len(s))

	for i, k := range s {
		if i == (len(s) - 1) {
			m[i] = map[string]interface{}{
				k: e,
			}
		} else {
			m[i] = map[string]interface{}{
				k: map[string]interface{}{},
			}
		}
	}

	for i := len(s) - 2; i != -1; i-- {
		m[i][s[i]] = m[i+1]
	}

	return m[0]
}

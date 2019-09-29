# Enki

Enki is for adding structured data as "markers" into arbitrary files. It uses a syntax that's similar to JSON and YAML but is meant for single lines. This module provides the unmarshalling and parsing functions.

This project is very much a work in progress.

## Installation

The import path for the package is github.com/stevexnicholls/enki

To install it, run:

```
go get github.com/stevexnicholls/enki
```

## Getting Started

### Syntax & Data Types

Enki data types are:

- Object: an unordered collection of name-value pairs (expressions) where the names are stings
- Number: a signed decimal number. This can be an integer or float.
- String: a sequence of characters. Strings that include special characters ```,:[]``` need to be surrounded by double quotes.
- Boolean: either of the values `true` or `false`. Using only an expression without a equals `=` translates to a `true` boolean
- Array: an ordered list of zero or more values surrounded by `[]`

### Markers 

Markers begin with the token `+!!:`. What follows is a combination of the supported data types above. The token is used by the Parser to identify enki lines.

### Examples

The following is an example of the enki syntax:

```
meta:file="thisfile.yml",process=true,choices=[delete,rename,copy,move],count=42
```

The following is an example of a marker. The `#` has been added to show that you should add markers as comments in code.

```
# +!!:meta:file="thisfile.yml",process=true,choices=[delete,rename,copy,move],count=42
```

The above would translate into the following JSON:

```
{
  "meta": {
    "choices": [
      "delete",
      "rename",
      "copy",
      "move"
    ],
    "count": 42,
    "file": "thisfile.yml",
    "process": true
  }
}
```

Example using `Unmarshal` in code to decode a single line:

```
package main

import (
  "fmt"
  "log"

  "github.com/stevexnicholls/enki"
)

var data = `meta:file="thisfile.yml",process=true,choices=[delete,rename,copy,move],count=42`

func main() {       
  m := make(map[string]interface{})

  err = enki.Unmarshal([]byte(data), &m)
  if err != nil {
          log.Fatalf("error: %v", err)
  }
  fmt.Printf("--- m:\n%v", m)
}
```

This example will generate the following output:

```
--- m:
map[
meta:map[choices:[delete rename copy move] count:42
 file:thisfile.yml process:true]]
```

Example using a `Parser` to parse multiple lines (as code comments):

```
package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/stevexnicholls/enki"
)

var data = `
# +!!:meta:file="thisfile.yml",process=true,choices=[delete,rename,copy,move],count=42
# +!!:vars:name="somevar",value=42,type=int
# +!!:meta:hello="world"
`

func main() {

	m := make(map[string]interface{})

	p := enki.NewParser(strings.NewReader(data), enki.ParserConfig{
		Token:        "+!!",
		Namespace:    "root",
		IncludeInput: false,
	})

	if err := p.Parse(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("--- m:\n%v", p.Data)
}
```

This example will generate the following output:

```
--- m:
map[:map[meta:map[choices:[delete rename copy move] count:42 file:thisfile.yml hello:world process:true] vars:map[name:somevar type:int value:42]]]
```

## Contributing

Please read [CONTRIBUTING.md](/CONTRIBUTING.md) for details on this project's code of conduct, and the process for submitting pull requests.

## Versioning

This project uses [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/stevexnicholls/enki/tags). 

## License

This project is licensed under the Apache License - see the [LICENSE](LICENSE) file for details

## Acknowledgments

The idea for this project was based off the [markers](https://book.kubebuilder.io/reference/markers.html) from the [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) project.
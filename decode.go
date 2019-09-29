package enki

import (
	"log"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
)

// Unmarshal parses the encoded data and returns the result
//
func Unmarshal(d []byte) (map[string]interface{}, error) {
	return unmarshal(d)
}

func unmarshal(d []byte) (map[string]interface{}, error) {
	return decode(string(d)), nil
}

func decode(s string) map[string]interface{} {
	// map to hold return map
	root := make(map[string]interface{}, 1)

	// split string into two parts on : or ,
	// : for objects
	// , when decoding fields
	el := splitStringAnyOutsideQuotesBrackets(s, ":,", true)

	// if there was strings returned
	if len(el) != 0 {

		// check if this is an object or fields based on ending char of first element
		// ends with : then it's an object, everything else is considered fields
		if el[0][len(el[0])-1:] == ":" { // object
			// todo(steve): make this less ugly
			// remove the ending : if exists
			lhs := hasSuffixRemove(el[0], ":")

			if len(el) == 2 {
				// object has a value; decode the second element
				// this is recursive and will process until end
				root[lhs] = decode(el[1])
			} else {
				// object is empty so make it an empty object
				root[lhs] = map[string]interface{}{}
			}
		} else { // expressions
			// create expression map
			exp := make(map[string]interface{})
			// split on =, returns two parts
			e := splitStringAnyOutsideQuotesBrackets(el[0], "=", false)
			lhs := hasSuffixRemove(e[0], ",")
			var rhs string
			if len(e) == 2 {
				rhs = hasSuffixRemove(e[1], ",")
				// equals something
				// check for rhs list, if not just use value
				if len(rhs) > 2 && rhs[0] == '[' && rhs[len(rhs)-1] == ']' {
					exp[lhs] = decodeList(rhs)
				} else {
					exp[lhs] = convertStringTo(rhs)
				}
			} else {
				// no equals, make it bool
				exp[lhs] = true
			}
			// handle nested objects
			if len(el) == 2 {
				if err := mergo.Map(&exp, decode(el[1])); err != nil {
					log.Fatal(err.Error())
				}
			}
			for k, z := range exp {
				root[k] = z
			}
		}
	}
	return root
}

func hasSuffixRemove(s string, suf string) string {
	if strings.HasSuffix(s, suf) {
		return s[:len(s)-len(suf)]
	}
	return s
}

func decodeList(s string) []interface{} {
	var o []interface{}
	g := splitStringOutsideQuotesBrackets(s[1:len(s)-1], ',')
	for _, a := range g {
		o = append(o, convertStringTo(a))
	}
	return o
}

func convertStringTo(s string) interface{} {
	switch {
	case isInt(s):
		c, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err.Error())
		}
		return c
	case isFloat(s):
		c, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Fatal(err.Error())
		}
		return c
	case isBool(s):
		c, err := strconv.ParseBool(s)
		if err != nil {
			log.Fatal(err.Error())
		}
		return c
	case isQuoted(s):
		c, err := strconv.Unquote(s)
		if err != nil {
			log.Fatal(err.Error())
		}
		return c
	default:
		c := s
		return c
	}
}

func isFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func isInt(s string) bool {
	_, err := strconv.ParseInt(s, 0, 64)
	return err == nil
}

func isBool(s string) bool {
	_, err := strconv.ParseBool(s)
	return err == nil
}

func isQuoted(s string) bool {
	_, err := strconv.Unquote(s)
	return err == nil
}

func indexAnyOutsideQuotesBrackets(s string, f string) int {
	q := false
	b := false
	for i, c := range s {
		switch c {
		case '"':
			q = !q
		case '[':
			b = true
		case ']':
			b = false
		default:
			if strings.ContainsAny(string(c), f) && !b && !q {
				return i
			}
		}
	}
	return -1
}

// TODO(steve): combine this with the other one and or make more and better funcs
// splitStringAnyOutsideQuotesBrackets splits a string using multiple splitters
// returns only two strings
// keep, keeps the splitter
func splitStringAnyOutsideQuotesBrackets(s string, f string, keep bool) []string {

	i := indexAnyOutsideQuotesBrackets(s, f)
	var r []string

	if i != -1 && i != len(s)-1 {
		if keep {
			r = append(r, s[0:i+1])
		} else {
			r = append(r, s[0:i])
		}
		r = append(r, s[i+1:])
	} else {
		r = append(r, s)
	}
	return r
}

// splitStringOutsideQuotesBrackets split the string
// returns all splits
func splitStringOutsideQuotesBrackets(str string, r rune) []string {
	inQuote := false
	inBracket := false
	f := func(c rune) bool {
		switch {
		case c == '"':
			inQuote = !inQuote
			return false
		case c == '[':
			inBracket = true
			return false
		case c == ']':
			inBracket = false
			return false
		case inQuote:
			return false
		case inBracket:
			return false
		default:
			return c == r
		}
	}
	return strings.FieldsFunc(str, f)
}

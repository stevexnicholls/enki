package enki

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"github.com/imdario/mergo"
)

// Parser ...
//
type Parser struct {
	Config ParserConfig
	Reader io.Reader
	Data   map[string]interface{}
}

// ParserConfig ...
//
type ParserConfig struct {
	Token        string
	Namespace    string
	IncludeInput bool
}

type line struct {
	text     []byte
	position int
}

// capture
//
type capture struct {
	buffer []byte
	header []byte
}

// NewParser ...
//
func NewParser(r io.Reader, c ParserConfig) *Parser {
	p := new(Parser)
	p.Config = c
	p.Reader = r
	return p
}

// Parse ...
//
func (p *Parser) Parse() error {
	scanner := bufio.NewScanner(p.Reader)
	var o []map[string]interface{}
	var l []line
	var c capture
	capture := false
	i := 0

	for scanner.Scan() {

		t := scanner.Text()

		l = append(l, line{
			text:     []byte(t),
			position: i,
		})

		if strings.Contains(t, p.Config.Token) { // token line
			h := strings.TrimSpace(t[strings.Index(t, p.Config.Token)+len(p.Config.Token):])

			if capture {
				m := c.process()
				o = append(o, m)
				capture = false
			}

			if (h[len(h)-1:]) == "=" {
				capture = true
				c.header = []byte(h)
			} else {
				m, err := Unmarshal([]byte(h))
				if err != nil {
					log.Fatal(err.Error())
				}
				o = append(o, m)
			}
		} else { // no token line
			if capture {
				c.buffer = append(c.buffer, []byte(t+"\n")...)
			}
		}
		i++
	}

	if capture {
		m := c.process()
		o = append(o, m)
		capture = false
	}

	if p.Config.IncludeInput {
		// full content
		var a []byte
		for _, l := range l {
			a = append(a, l.text...)
			a = append(a, []byte("\n")...)
		}
		t := fmt.Sprintf("_input:" + p.Config.Namespace + ":content=" + base64.StdEncoding.EncodeToString(a))
		m, err := Unmarshal([]byte(t))
		if err != nil {
			log.Fatal(err.Error())
		}
		o = append(o, m)

		// lines
		d := len(strconv.Itoa(i))
		for _, l := range l {
			h := fmt.Sprintf("_input:" + p.Config.Namespace + ":src:")
			l := fmt.Sprintf("%s%0[2]*d=\"%s\"", h, d, l.position, base64.StdEncoding.EncodeToString(l.text))
			m, err := Unmarshal([]byte(l))
			if err != nil {
				log.Fatal(err.Error())
			}
			o = append(o, m)
		}
	}

	for _, e := range o {
		if err := mergo.Map(&p.Data, e); err != nil {
			log.Fatal(err.Error())
		}
	}

	return nil
}

func (c *capture) process() map[string]interface{} {
	e := []byte(base64.StdEncoding.EncodeToString(c.buffer))
	m, _ := Unmarshal(append(c.header, e...))
	c.buffer = nil
	c.header = nil
	return m
}

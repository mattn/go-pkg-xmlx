// Copyright (c) 2010, Jim Teeuwen. All rights reserved.
// This code is subject to a 1-clause BSD license.
// The contents of which can be found in the LICENSE file.

/*
 This package wraps the standard XML library and uses it to build a node tree of
 any document you load. This allows you to look up nodes forwards and backwards,
 as well as perform simple search queries.

 Nodes now simply become collections and don't require you to read them in the
 order in which the xml.Parser finds them.

 The Document currently implements 2 search functions which allow you to
 look for specific nodes.

   *xmlx.Document.SelectNode(namespace, name string) *Node;
   *xmlx.Document.SelectNodes(namespace, name string) []*Node;

 SelectNode() returns the first, single node it finds matching the given name
 and namespace. SelectNodes() returns a slice containing all the matching nodes.

 Note that these search functions can be invoked on individual nodes as well.
 This allows you to search only a subset of the entire document.
*/
package xmlx

import (
	"os"
	"io"
	"io/ioutil"
	"path"
	"strings"
	"xml"
	"fmt"
	"http"
	"iconv"
)

type Document struct {
	Version     string
	Encoding    string
	StandAlone  string
	SaveDocType bool
	Root        *Node
	Entity      map[string]string
	Verbose     bool
}

func New() *Document {
	return &Document{
		Version:     "1.0",
		Encoding:    "utf-8",
		StandAlone:  "yes",
		SaveDocType: true,
		Entity:      make(map[string]string),
		Verbose:     false,
	}
}

// This loads a rather massive table of non-conventional xml escape sequences.
// Needed to make the parser map them to characters properly. It is advised to
// set only those entities needed manually using the document.Entity map, but
// if need be, this method can be called to fill the map with the entire set
// defined on http://www.w3.org/TR/html4/sgml/entities.html
func (this *Document) LoadExtendedEntityMap() { loadNonStandardEntities(this.Entity) }

func (this *Document) String() string {
	s, _ := this.SaveString()
	return s
}

func (this *Document) SelectNode(namespace, name string) *Node {
	return this.Root.SelectNode(namespace, name)
}

func (this *Document) SelectNodes(namespace, name string) []*Node {
	return this.Root.SelectNodes(namespace, name)
}

// *****************************************************************************
// *** Satisfy ILoader interface
// *****************************************************************************
func (this *Document) LoadString(s string) (err os.Error) {
	// Ensure we are passing UTF-8 encoding content to the XML tokenizer.
	if s, err = this.correctEncoding(s); err != nil {
		return
	}

	// tokenize data
	xp := xml.NewParser(strings.NewReader(s))
	xp.Entity = this.Entity

	this.Root = NewNode(NT_ROOT)
	ct := this.Root

	var tok xml.Token
	var t *Node
	var i int
	var doctype string
	var v xml.Attr

	for {
		if tok, err = xp.Token(); err != nil {
			if err == os.EOF {
				return nil
			}

			if this.Verbose {
				fmt.Fprintf(os.Stderr, "Xml Error: %s\n", err)
			}
			return err
		}

		switch tt := tok.(type) {
		case xml.SyntaxError:
			return os.NewError(tt.String())
		case xml.CharData:
			ct.Value = strings.TrimSpace(string([]byte(tt)))
		case xml.Comment:
			t := NewNode(NT_COMMENT)
			t.Value = strings.TrimSpace(string([]byte(tt)))
			ct.AddChild(t)
		case xml.Directive:
			t = NewNode(NT_DIRECTIVE)
			t.Value = strings.TrimSpace(string([]byte(tt)))
			ct.AddChild(t)
		case xml.StartElement:
			t = NewNode(NT_ELEMENT)
			t.Name = tt.Name
			t.Attributes = make([]*Attr, len(tt.Attr))
			for i, v = range tt.Attr {
				t.Attributes[i] = new(Attr)
				t.Attributes[i].Name = v.Name
				t.Attributes[i].Value = v.Value
			}
			ct.AddChild(t)
			ct = t
		case xml.ProcInst:
			if tt.Target == "xml" { // xml doctype
				doctype = strings.TrimSpace(string(tt.Inst))
				if i = strings.Index(doctype, `standalone="`); i > -1 {
					this.StandAlone = doctype[i+len(`standalone="`) : len(doctype)]
					i = strings.Index(this.StandAlone, `"`)
					this.StandAlone = this.StandAlone[0:i]
				}
			} else {
				t = NewNode(NT_PROCINST)
				t.Target = strings.TrimSpace(tt.Target)
				t.Value = strings.TrimSpace(string(tt.Inst))
				ct.AddChild(t)
			}
		case xml.EndElement:
			if ct = ct.Parent; ct == nil {
				return
			}
		}
	}

	return
}

func (this *Document) LoadFile(filename string) (err os.Error) {
	var data []byte

	if data, err = ioutil.ReadFile(path.Clean(filename)); err != nil {
		return
	}

	return this.LoadString(string(data))
}

func (this *Document) LoadUri(uri string) (err os.Error) {
	var r *http.Response
	if r, _, err = http.Get(uri); err != nil {
		return
	}

	defer r.Body.Close()

	var b []byte
	if b, err = ioutil.ReadAll(r.Body); err != nil {
		return
	}

	return this.LoadString(string(b))
}

func (this *Document) LoadStream(r io.Reader) (err os.Error) {
	var b []byte
	if b, err = ioutil.ReadAll(r); err != nil {
		return
	}
	return this.LoadString(string(b))
}

// *****************************************************************************
// *** Satisfy ISaver interface
// *****************************************************************************
func (this *Document) SaveFile(path string) (err os.Error) {
	var data string
	if data, err = this.SaveString(); err != nil {
		return
	}

	return ioutil.WriteFile(path, []byte(data), 0600)
}

func (this *Document) SaveString() (s string, err os.Error) {
	if this.SaveDocType {
		s = fmt.Sprintf(`<?xml version="%s" encoding="%s" standalone="%s"?>`,
			this.Version, this.Encoding, this.StandAlone)
	}

	s += this.Root.String()
	return
}

func (this *Document) SaveStream(w io.Writer) (err os.Error) {
	var s string
	if s, err = this.SaveString(); err != nil {
		return
	}
	_, err = w.Write([]byte(s))
	return
}

// Use libiconv to ensure we get UTF-8 encoded data. The Go Xml tokenizer will
// throw a tantrum if we give it anything else.
func (this *Document) correctEncoding(data string) (ret string, err os.Error) {
	var cd *iconv.Iconv
	var tok xml.Token

	enc := "utf-8"
	xp := xml.NewParser(strings.NewReader(data))
	xp.Entity = this.Entity

loop:
	for {
		if tok, err = xp.Token(); err != nil {
			if err == os.EOF {
				break loop
			}
			return "", err
		}

		switch tt := tok.(type) {
		case xml.ProcInst:
			if tt.Target == "xml" { // xml doctype
				enc = strings.ToLower(string(tt.Inst))
				if i := strings.Index(enc, `encoding="`); i > -1 {
					enc = enc[i+len(`encoding="`):]
					i = strings.Index(enc, `"`)
					enc = enc[:i]
					break loop
				}
			}
		}
	}

	if enc == "utf-8" {
		return data, nil
	}

	if cd, err = iconv.Open("utf-8", enc); err != nil {
		return
	}

	defer cd.Close()
	return cd.Conv(data)
}

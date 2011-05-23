// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

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
	"strings"
	"bytes"
	"xml"
	"fmt"
	"http"
	"go-charset.googlecode.com/hg/charset"
)

// represents a single XML document.
type Document struct {
	Version     string            // XML version
	Encoding    string            // Encoding found in document. If absent, assumes UTF-8.
	StandAlone  string            // Value of XML doctype's 'standalone' attribute.
	SaveDocType bool              // Whether not to include the XML doctype in saves.
	Root        *Node             // The document's root node.
	Entity      map[string]string // Mapping of custom entity conversions.
	Verbose     bool              // [depracated] Not actually used anymore.
}

// Create a new, empty XML document instance.
func New() *Document {
	return &Document{
		Version:     "1.0",
		Encoding:    "utf-8",
		StandAlone:  "yes",
		SaveDocType: true,
		Entity:      make(map[string]string),
	}
}

// This loads a rather massive table of non-conventional xml escape sequences.
// Needed to make the parser map them to characters properly. It is advised to
// set only those entities needed manually using the document.Entity map, but
// if need be, this method can be called to fill the map with the entire set
// defined on http://www.w3.org/TR/html4/sgml/entities.html
func (this *Document) LoadExtendedEntityMap() { loadNonStandardEntities(this.Entity) }

// Select a single node with the given namespace and name. Returns nil if no
// matching node was found.
func (this *Document) SelectNode(namespace, name string) *Node {
	return this.Root.SelectNode(namespace, name)
}

// Select all nodes with the given namespace and name. Returns an empty slice
// if no matches were found.
func (this *Document) SelectNodes(namespace, name string) []*Node {
	return this.Root.SelectNodes(namespace, name)
}

// Load the contents of this document from the supplied reader.
func (this *Document) LoadStream(r io.Reader) (err os.Error) {
	xp := xml.NewParser(r)
	xp.Entity = this.Entity
	xp.CharsetReader = func(enc string, input io.Reader) (io.Reader, os.Error) {
		return charset.NewReader(enc, input)
	}

	this.Root = NewNode(NT_ROOT)
	ct := this.Root

	var tok xml.Token
	var t *Node
	var doctype string

	for {
		if tok, err = xp.Token(); err != nil {
			if err == os.EOF {
				return nil
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
			for i, v := range tt.Attr {
				t.Attributes[i] = new(Attr)
				t.Attributes[i].Name = v.Name
				t.Attributes[i].Value = v.Value
			}
			ct.AddChild(t)
			ct = t
		case xml.ProcInst:
			if tt.Target == "xml" { // xml doctype
				doctype = strings.TrimSpace(string(tt.Inst))
				if i := strings.Index(doctype, `standalone="`); i > -1 {
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

// Load the contents of this document from the supplied byte slice.
func (this *Document) LoadBytes(d []byte) (err os.Error) {
	return this.LoadStream(bytes.NewBuffer(d))
}

// Load the contents of this document from the supplied string.
func (this *Document) LoadString(s string) (err os.Error) {
	return this.LoadStream(strings.NewReader(s))
}

// Load the contents of this document from the supplied file.
func (this *Document) LoadFile(filename string) (err os.Error) {
	var fd *os.File
	if fd, err = os.Open(filename); err != nil {
		return
	}

	defer fd.Close()
	return this.LoadStream(fd)
}

// Load the contents of this document from the supplied uri.
func (this *Document) LoadUri(uri string) (err os.Error) {
	var r *http.Response
	if r, err = http.Get(uri); err != nil {
		return
	}

	defer r.Body.Close()
	return this.LoadStream(r.Body)
}

// Save the contents of this document to the supplied file.
func (this *Document) SaveFile(path string) os.Error {
	return ioutil.WriteFile(path, this.SaveBytes(), 0600)
}

// Save the contents of this document as a byte slice.
func (this *Document) SaveBytes() []byte {
	var b bytes.Buffer

	if this.SaveDocType {
		b.WriteString(fmt.Sprintf(`<?xml version="%s" encoding="%s" standalone="%s"?>`,
			this.Version, this.Encoding, this.StandAlone))
	}

	b.Write(this.Root.Bytes())
	return b.Bytes()
}

// Save the contents of this document as a string.
func (this *Document) SaveString() string { return string(this.SaveBytes()) }

// Alias for Document.SaveString(). This one is invoked by anything looking for
// the standard String() method (eg: fmt.Printf("%s\n", mydoc).
func (this *Document) String() string { return string(this.SaveBytes()) }

// Save the contents of this document to the supplied writer.
func (this *Document) SaveStream(w io.Writer) (err os.Error) {
	_, err = w.Write(this.SaveBytes())
	return
}

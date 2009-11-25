// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 Author: Jim Teeuwen <jimteeuwen@gmail.com>

 This package wraps the standard XML library and uses it to build a node tree of
 any document you load. This allows you to look up nodes forwards and backwards,
 as well as perform search queries (no xpath support yet).

 Nodes now simply become collections and don't require you to read them in the
 order in which the xml.Parser finds them.

 The Document currently implements 2 simple search functions which allow you to
 look for specific nodes.
 
   Document.SelectNode(namespace, name string) *Node;
   Document.SelectNodes(namespace, name string) []*Node;
 
 SelectNode() returns the first, single node it finds matching the given name
 and namespace. SelectNodes() returns a slice containing all the matching nodes.
 
 Note that these search functions can be invoked on individual nodes as well.
 This allows you to search only a subset of the entire document.

*/
package xmlx

import "os"
import "io"
import "strings"
import "xml"
import "fmt"
import "http"

type Document struct {
	Version		string;
	Encoding	string;
	StandAlone	string;
	SaveDocType	bool;
	Root		*Node;
	Entity		map[string]string;
	Verbose		bool;
}

func New() *Document {
	return &Document{
		Version: "1.0",
		Encoding: "utf-8",
		StandAlone: "yes",
		SaveDocType: true,
		Entity: make(map[string]string),
		Verbose: false
	}
}

// This loads a rather massive table of non-conventional xml escape sequences.
// Needed to make the parser map them to characters properly. It is advised to
// set only those entities needed manually using the document.Entity map, but
// if need be, this method can be called to fill the map with the entire set
// defined on http://www.w3.org/TR/html4/sgml/entities.html
func (this *Document) LoadExtendedEntityMap() {
	loadNonStandardEntities(&this.Entity);
}

func (this *Document) String() string {
	s, _ := this.SaveString();
	return s;
}

func (this *Document) SelectNode(namespace, name string) *Node {
	return this.Root.SelectNode(namespace, name);
}

func (this *Document) SelectNodes(namespace, name string) []*Node {
	return this.Root.SelectNodes(namespace, name);
}

// *****************************************************************************
// *** Satisfy ILoader interface
// *****************************************************************************
func (this *Document) LoadString(s string) (err os.Error) {
	xp := xml.NewParser(strings.NewReader(s));
	xp.Entity = this.Entity;

	this.Root = NewNode(NT_ROOT);
	ct := this.Root;

	for {
		tok, err := xp.Token();
		if err != nil {
			if err != os.EOF && this.Verbose {
				fmt.Fprintf(os.Stderr, "Xml Error: %s\n", err);
			}
			return
		}

		t1, ok := tok.(xml.SyntaxError);
		if ok {
			err = os.NewError(t1.String());
			return
		}

		t2, ok := tok.(xml.CharData);
		if ok && ct != nil {
			ct.Value = strings.TrimSpace(string(t2));
			continue
		}

		t3, ok := tok.(xml.Comment);
		if ok && ct != nil {
			t := NewNode(NT_COMMENT);
			t.Value = strings.TrimSpace(string(t3));
			ct.AddChild(t);
			continue
		}

		t4, ok := tok.(xml.Directive);
		if ok && ct != nil {
			t := NewNode(NT_DIRECTIVE);
			t.Value = strings.TrimSpace(string(t4));
			ct.AddChild(t);
			continue
		}

		t5, ok := tok.(xml.StartElement);
		if ok && ct != nil {
			t := NewNode(NT_ELEMENT);
			t.Name = t5.Name;
			t.Attributes = make([]Attr, len(t5.Attr));
			for i, v := range t5.Attr {
				t.Attributes[i].Name = v.Name;
				t.Attributes[i].Value = v.Value;
			}
			ct.AddChild(t);
			ct = t;
			continue
		}

		t6, ok := tok.(xml.ProcInst);
		if ok {
			if t6.Target == "xml" {	// xml doctype
				doctype := strings.TrimSpace(string(t6.Inst));
				
				/* // Not needed. There is only xml version 1.0
				pos := strings.Index(doctype, `version="`);
				if pos > -1 {
					this.Version = doctype[pos+len(`version="`) : len(doctype)];
					pos = strings.Index(this.Version, `"`);
					this.Version = this.Version[0:pos];
				}
				*/

				/* // Not needed. Any string we handle in Go is UTF8
				   // encoded. This means we will save UTF8 data as well. 
				pos = strings.Index(doctype, `encoding="`);
				if pos > -1 {
					this.Encoding = doctype[pos+len(`encoding="`) : len(doctype)];
					pos = strings.Index(this.Encoding, `"`);
					this.Encoding = this.Encoding[0:pos];
				}
				*/

				pos := strings.Index(doctype, `standalone="`);
				if pos > -1 {
					this.StandAlone = doctype[pos+len(`standalone="`) : len(doctype)];
					pos = strings.Index(this.StandAlone, `"`);
					this.StandAlone = this.StandAlone[0:pos];
				}
			} else if ct != nil {
				t := NewNode(NT_PROCINST);
				t.Target = strings.TrimSpace(t6.Target);
				t.Value = strings.TrimSpace(string(t6.Inst));
				ct.AddChild(t);
			}
			continue
		}

		_, ok = tok.(xml.EndElement);
		if ok {
			ct = ct.Parent;
			continue
		}
	}

	return;
}

func (this *Document) LoadFile(path string) (err os.Error) {
	file, err := os.Open(path, os.O_RDONLY, 0600);
	if err != nil {
		return
	}
	defer file.Close();

	content := "";
	buff := make([]byte, 256);
	for {
		_, err := file.Read(buff);
		if err != nil {
			break
		}
		content += string(buff);
	}

	err = this.LoadString(content);
	return;
}

func (this *Document) LoadUri(uri string) (err os.Error) {
	r, _, err := http.Get(uri);
	if err != nil {
		return
	}

	defer r.Body.Close();

	b, err := io.ReadAll(r.Body);
	if err != nil {
		return
	}

	err = this.LoadString(string(b));
	return;
}

func (this *Document) LoadStream(r *io.Reader) (err os.Error) {
	content := "";
	buff := make([]byte, 256);
	for {
		_, err := r.Read(buff);
		if err != nil {
			break
		}
		content += string(buff);
	}

	err = this.LoadString(content);
	return;
}

// *****************************************************************************
// *** Satisfy ISaver interface
// *****************************************************************************
func (this *Document) SaveFile(path string) (err os.Error) {
	file, err := os.Open(path, os.O_WRONLY | os.O_CREAT, 0600);
	if err != nil {
		return
	}
	defer file.Close();

	content, err  := this.SaveString();
	if err != nil {
		return
	}

	file.Write(strings.Bytes(content));
	return
}

func (this *Document) SaveString() (s string, err os.Error) {
	if this.SaveDocType {
		s = fmt.Sprintf(`<?xml version="%s" encoding="%s" standalone="%s"?>`,
			this.Version, this.Encoding, this.StandAlone)
	}

	s += this.Root.String();
	return;
}

func (this *Document) SaveStream(w *io.Writer) (err os.Error) {
	s, err := this.SaveString();
	if err != nil {
		return
	}
	w.Write(strings.Bytes(s));
	return;
}


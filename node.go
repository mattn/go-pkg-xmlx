/*
Copyright (c) 2010, Jim Teeuwen.
All rights reserved.

This code is subject to a 1-clause BSD license.
The contents of which can be found in the LICENSE file.
*/

package xmlx

import (
	"os"
	"strings"
	"xml"
	"bytes"
	"fmt"
	"strconv"
)

const (
	NT_ROOT = iota
	NT_DIRECTIVE
	NT_PROCINST
	NT_COMMENT
	NT_ELEMENT
)

type Attr struct {
	Name  xml.Name
	Value string
}

type Node struct {
	Type       byte
	Name       xml.Name
	Children   []*Node
	Attributes []*Attr
	Parent     *Node
	Value      string
	Target     string // procinst field
}

func NewNode(tid byte) *Node {
	n := new(Node)
	n.Type = tid
	n.Children = make([]*Node, 0, 10)
	n.Attributes = make([]*Attr, 0, 10)
	return n
}

// This wraps the standard xml.Unmarshal function and supplies this particular
// node as the content to be unmarshalled.
func (this *Node) Unmarshal(obj interface{}) os.Error {
	return xml.Unmarshal(strings.NewReader(this.String()), obj)
}

// Get node value as string
func (this *Node) S(namespace, name string) string {
	if node := rec_SelectNode(this, namespace, name); node != nil {
		return node.Value
	}
	return ""
}

// Deprecated - use Node.S()
func (this *Node) GetValue(namespace, name string) string { return this.S(namespace, name) }


// Get node value as int
func (this *Node) I(namespace, name string) int {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atoi(node.Value)
		return n
	}
	return 0
}

// Deprecated - use Node.I()
func (this *Node) GetValuei(namespace, name string) int { return this.I(namespace, name) }


// Get node value as int64
func (this *Node) I64(namespace, name string) int64 {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atoi64(node.Value)
		return n
	}
	return 0
}

// Deprecated - use Node.I64()
func (this *Node) GetValuei64(namespace, name string) int64 { return this.I64(namespace, name) }


// Get node value as uint
func (this *Node) U(namespace, name string) uint {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atoui(node.Value)
		return n
	}
	return 0
}

// Deprecated - use Node.U()
func (this *Node) GetValueui(namespace, name string) uint { return this.U(namespace, name) }


// Get node value as uint64
func (this *Node) U64(namespace, name string) uint64 {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atoui64(node.Value)
		return n
	}
	return 0
}

// Deprecated - use Node.U64()
func (this *Node) GetValueui64(namespace, name string) uint64 { return this.U64(namespace, name) }

// Get node value as float32
func (this *Node) F32(namespace, name string) float32 {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atof32(node.Value)
		return n
	}
	return 0
}

// Deprecated - use Node.F32()
func (this *Node) GetValuef32(namespace, name string) float32 { return this.F32(namespace, name) }


// Get node value as float64
func (this *Node) F64(namespace, name string) float64 {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atof64(node.Value)
		return n
	}
	return 0
}

// Deprecated - use Node.F64()
func (this *Node) GetValuef64(namespace, name string) float64 { return this.F64(namespace, name) }


// Get node value as bool
func (this *Node) B(namespace, name string) bool {
	if node := rec_SelectNode(this, namespace, name); node != nil && node.Value != "" {
		n, _ := strconv.Atob(node.Value)
		return n
	}
	return false
}

// Get attribute value as string
func (this *Node) As(namespace, name string) string {
	for _, v := range this.Attributes {
		if namespace == v.Name.Space && name == v.Name.Local {
			return v.Value
		}
	}
	return ""
}

// Deprecated - use Node.As()
func (this *Node) GetAttr(namespace, name string) string { return this.As(namespace, name) }


// Get attribute value as int
func (this *Node) Ai(namespace, name string) int {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atoi(s)
		return n
	}
	return 0
}

// Deprecated - use Node.Ai()
func (this *Node) GetAttri(namespace, name string) int { return this.Ai(namespace, name) }


// Get attribute value as uint
func (this *Node) Au(namespace, name string) uint {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atoui(s)
		return n
	}
	return 0
}

// Deprecated - use Node.Au()
func (this *Node) GetAttrui(namespace, name string) uint { return this.Au(namespace, name) }


// Get attribute value as uint64
func (this *Node) Au64(namespace, name string) uint64 {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atoui64(s)
		return n
	}
	return 0
}

// Deprecated - use Node.Au64()
func (this *Node) GetAttrui64(namespace, name string) uint64 { return this.Au64(namespace, name) }


// Get attribute value as int64
func (this *Node) Ai64(namespace, name string) int64 {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atoi64(s)
		return n
	}
	return 0
}

// Deprecated - use Node.Ai64()
func (this *Node) GetAttri64(namespace, name string) int64 { return this.Ai64(namespace, name) }

// Get attribute value as float32
func (this *Node) Af32(namespace, name string) float32 {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atof32(s)
		return n
	}
	return 0
}

// Deprecated - use Node.Af32()
func (this *Node) GetAttrf32(namespace, name string) float32 { return this.Af32(namespace, name) }

// Get attribute value as float64
func (this *Node) Af64(namespace, name string) float64 {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atof64(s)
		return n
	}
	return 0
}

// Deprecated - use Node.Af64()
func (this *Node) GetAttrf64(namespace, name string) float64 { return this.Af64(namespace, name) }


// Get attribute value as bool
func (this *Node) Ab(namespace, name string) bool {
	if s := this.As(namespace, name); s != "" {
		n, _ := strconv.Atob(s)
		return n
	}
	return false
}


// Returns true if this node has the specified attribute. False otherwise.
func (this *Node) HasAttr(namespace, name string) bool {
	for _, v := range this.Attributes {
		if namespace == v.Name.Space && name == v.Name.Local {
			return true
		}
	}
	return false
}

// Select single node by name
func (this *Node) SelectNode(namespace, name string) *Node {
	return rec_SelectNode(this, namespace, name)
}

func rec_SelectNode(cn *Node, namespace, name string) *Node {
	if cn.Name.Space == namespace && cn.Name.Local == name {
		return cn
	}

	var tn *Node
	for _, v := range cn.Children {
		if tn = rec_SelectNode(v, namespace, name); tn != nil {
			return tn
		}
	}
	return nil
}

// Select multiple nodes by name
func (this *Node) SelectNodes(namespace, name string) []*Node {
	list := make([]*Node, 0, 16)
	rec_SelectNodes(this, namespace, name, &list)
	return list
}

func rec_SelectNodes(cn *Node, namespace, name string, list *[]*Node) {
	if cn.Name.Space == namespace && cn.Name.Local == name {
		*list = append(*list, cn)
		return
	}

	for _, v := range cn.Children {
		rec_SelectNodes(v, namespace, name, list)
	}
}

// Convert node to appropriate string representation based on it's @Type.
// Note that NT_ROOT is a special-case empty node used as the root for a
// Document. This one has no representation by itself. It merely forwards the
// String() call to it's child nodes.
func (this *Node) String() (s string) {
	switch this.Type {
	case NT_PROCINST:
		s = this.printProcInst()
	case NT_COMMENT:
		s = this.printComment()
	case NT_DIRECTIVE:
		s = this.printDirective()
	case NT_ELEMENT:
		s = this.printElement()
	case NT_ROOT:
		s = this.printRoot()
	}
	return
}

func (this *Node) printRoot() (s string) {
	var data []byte
	buf := bytes.NewBuffer(data)
	for _, v := range this.Children {
		buf.WriteString(v.String())
	}
	return buf.String()
}

func (this *Node) printProcInst() string {
	return "<?" + this.Target + " " + this.Value + "?>"
}

func (this *Node) printComment() string {
	return "<!-- " + this.Value + " -->"
}

func (this *Node) printDirective() string {
	return "<!" + this.Value + "!>"
}

func (this *Node) printElement() string {
	var data []byte
	buf := bytes.NewBuffer(data)

	if len(this.Name.Space) > 0 {
		buf.WriteRune('<')
		buf.WriteString(this.Name.Space)
		buf.WriteRune(':')
		buf.WriteString(this.Name.Local)
	} else {
		buf.WriteRune('<')
		buf.WriteString(this.Name.Local)
	}

	for _, v := range this.Attributes {
		if len(v.Name.Space) > 0 {
			buf.WriteString(fmt.Sprintf(` %s:%s="%s"`, v.Name.Space, v.Name.Local, v.Value))
		} else {
			buf.WriteString(fmt.Sprintf(` %s="%s"`, v.Name.Local, v.Value))
		}
	}

	if len(this.Children) == 0 && len(this.Value) == 0 {
		buf.WriteString(" />")
		return buf.String()
	}

	buf.WriteRune('>')

	for _, v := range this.Children {
		buf.WriteString(v.String())
	}

	buf.WriteString(this.Value)
	if len(this.Name.Space) > 0 {
		buf.WriteString("</")
		buf.WriteString(this.Name.Space)
		buf.WriteRune(':')
		buf.WriteString(this.Name.Local)
		buf.WriteRune('>')
	} else {
		buf.WriteString("</")
		buf.WriteString(this.Name.Local)
		buf.WriteRune('>')
	}

	return buf.String()
}

// Add a child node
func (this *Node) AddChild(t *Node) {
	if t.Parent != nil {
		t.Parent.RemoveChild(t)
	}
	t.Parent = this
	this.Children = append(this.Children, t)
}

// Remove a child node
func (this *Node) RemoveChild(t *Node) {
	p := -1
	for i, v := range this.Children {
		if v == t {
			p = i
			break
		}
	}

	if p == -1 {
		return
	}

	copy(this.Children[p:], this.Children[p+1:])
	this.Children = this.Children[0 : len(this.Children)-1]

	t.Parent = nil
}

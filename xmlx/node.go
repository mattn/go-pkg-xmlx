/*
Copyright (c) 2010, Jim Teeuwen.
All rights reserved.

This code is subject to a 1-clause BSD license.
The contents of which can be found in the LICENSE file.
*/

package xmlx

import "os"
import "strings"
import "xml"
import "bytes"
import "fmt"
import "strconv"

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
func (this *Node) GetValue(namespace, name string) string {
	node := rec_SelectNode(this, namespace, name)
	if node == nil {
		return ""
	}
	return node.Value
}

// Get node value as int
func (this *Node) GetValuei(namespace, name string) int {
	node := rec_SelectNode(this, namespace, name)
	if node == nil || node.Value == "" {
		return 0
	}
	n, _ := strconv.Atoi(node.Value)
	return n
}

// Get node value as int64
func (this *Node) GetValuei64(namespace, name string) int64 {
	node := rec_SelectNode(this, namespace, name)
	if node == nil || node.Value == "" {
		return 0
	}
	n, _ := strconv.Atoi64(node.Value)
	return n
}

// Get node value as uint
func (this *Node) GetValueui(namespace, name string) uint {
	node := rec_SelectNode(this, namespace, name)
	if node == nil || node.Value == "" {
		return 0
	}
	n, _ := strconv.Atoui(node.Value)
	return n
}

// Get node value as uint64
func (this *Node) GetValueui64(namespace, name string) uint64 {
	node := rec_SelectNode(this, namespace, name)
	if node == nil || node.Value == "" {
		return 0
	}
	n, _ := strconv.Atoui64(node.Value)
	return n
}

// Get node value as float
func (this *Node) GetValuef(namespace, name string) float {
	node := rec_SelectNode(this, namespace, name)
	if node == nil || node.Value == "" {
		return 0
	}
	n, _ := strconv.Atof(node.Value)
	return n
}

// Get node value as float32
func (this *Node) GetValuef32(namespace, name string) float32 {
	node := rec_SelectNode(this, namespace, name)
	if node == nil || node.Value == "" {
		return 0
	}
	n, _ := strconv.Atof32(node.Value)
	return n
}

// Get node value as float64
func (this *Node) GetValuef64(namespace, name string) float64 {
	node := rec_SelectNode(this, namespace, name)
	if node == nil || node.Value == "" {
		return 0
	}
	n, _ := strconv.Atof64(node.Value)
	return n
}

// Get attribute value as string
func (this *Node) GetAttr(namespace, name string) string {
	for _, v := range this.Attributes {
		if namespace == v.Name.Space && name == v.Name.Local {
			return v.Value
		}
	}
	return ""
}

// Get attribute value as int
func (this *Node) GetAttri(namespace, name string) int {
	s := this.GetAttr(namespace, name)
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}

// Get attribute value as uint
func (this *Node) GetAttrui(namespace, name string) uint {
	s := this.GetAttr(namespace, name)
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoui(s)
	return n
}

// Get attribute value as uint64
func (this *Node) GetAttrui64(namespace, name string) uint64 {
	s := this.GetAttr(namespace, name)
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoui64(s)
	return n
}

// Get attribute value as int64
func (this *Node) GetAttri64(namespace, name string) int64 {
	s := this.GetAttr(namespace, name)
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi64(s)
	return n
}

// Get attribute value as float
func (this *Node) GetAttrf(namespace, name string) float {
	s := this.GetAttr(namespace, name)
	if s == "" {
		return 0
	}
	n, _ := strconv.Atof(s)
	return n
}

// Get attribute value as float32
func (this *Node) GetAttrf32(namespace, name string) float32 {
	s := this.GetAttr(namespace, name)
	if s == "" {
		return 0
	}
	n, _ := strconv.Atof32(s)
	return n
}

// Get attribute value as float64
func (this *Node) GetAttrf64(namespace, name string) float64 {
	s := this.GetAttr(namespace, name)
	if s == "" {
		return 0
	}
	n, _ := strconv.Atof64(s)
	return n
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
		l := len(*list)
		if l >= cap(*list) {
			c := make([]*Node, l, l+16)
			copy(c, *list)
			*list = c
		}
		*list = (*list)[0 : l+1]
		(*list)[l] = cn
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

	l := len(this.Children)
	if l >= cap(this.Children) {
		c := make([]*Node, l, l+10)
		copy(c, this.Children)
		this.Children = c
	}

	this.Children = this.Children[0 : l+1]
	this.Children[l] = t
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

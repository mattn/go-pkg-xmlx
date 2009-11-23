package xmlx

import "xml"
import "fmt"

const (
	NT_ROOT		= 0x00;
	NT_DIRECTIVE	= 0x01;
	NT_PROCINST	= 0x02;
	NT_COMMENT	= 0x03;
	NT_ELEMENT	= 0x04;
)

type Attr struct {
	Name	xml.Name;
	Value	string;
}

type Node struct {
	Type		byte;
	Name		xml.Name;
	Children	[]*Node;
	Attributes	[]Attr;
	Parent		*Node;
	Value		string;

	// procinst field
	Target	string;
}

func NewNode(tid byte) *Node	{ return &Node{Type: tid} }

func (this *Node) SelectNode(namespace, name string) *Node {
	return rec_SelectNode(this, namespace, name);
}

func rec_SelectNode(cn *Node, namespace, name string) *Node {
	if cn.Name.Space == namespace && cn.Name.Local == name {
		return cn;
	}

	for _, v := range cn.Children {
		tn := rec_SelectNode(v, namespace, name);
		if tn != nil { return tn }
	}
	return nil;
}

func (this *Node) SelectNodes(namespace, name string) []*Node {
	list := make([]*Node, 0);
	rec_SelectNodes(this, namespace, name, &list);
	return list;
}

func rec_SelectNodes(cn *Node, namespace, name string, list *[]*Node) {
	if cn.Name.Space == namespace && cn.Name.Local == name {
		slice := make([]*Node, len(*list) + 1);
		for i,v := range *list {
			slice[i] = v;
		}
		slice[len(slice) - 1] = cn;
		*list = slice;
		return
	}

	for _, v := range cn.Children {
		rec_SelectNodes(v, namespace, name, list);
	}
}

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
	return;
}

func (this *Node) printRoot() (s string) {
	for _, v := range this.Children {
		s += v.String()
	}
	return;
}

func (this *Node) printProcInst() (s string) {
	s = "<?" + this.Target + " " + this.Value + "?>";
	return;
}

func (this *Node) printComment() (s string) {
	s = "<!-- " + this.Value + " -->";
	return;
}

func (this *Node) printDirective() (s string) {
	s = "<!" + this.Value + "!>";
	return;
}

func (this *Node) printElement() (s string) {
	if len(this.Name.Space) > 0 {
		s = "<" + this.Name.Space + ":" + this.Name.Local
	} else {
		s = "<" + this.Name.Local
	}

	for _, v := range this.Attributes {
		if len(v.Name.Space) > 0 {
			s += fmt.Sprintf(` %s:%s="%s"`, v.Name.Space, v.Name.Local, v.Value)
		} else {
			s += fmt.Sprintf(` %s="%s"`, v.Name.Local, v.Value)
		}
	}

	if len(this.Children) == 0 && len(this.Value) == 0 {
		s += " />";
		return;
	}

	s += ">";

	for _, v := range this.Children {
		s += v.String()
	}

	s += this.Value;
	if len(this.Name.Space) > 0 {
		s += "</" + this.Name.Space + ":" + this.Name.Local + ">"
	} else {
		s += "</" + this.Name.Local + ">"
	}
	return;
}

func (this *Node) AddChild(t *Node) {
	if t.Parent != nil {
		t.Parent.RemoveChild(t)
	}
	t.Parent = this;

	slice := make([]*Node, len(this.Children)+1);
	for i, v := range this.Children {
		slice[i] = v
	}
	slice[len(slice)-1] = t;
	this.Children = slice;
}

func (this *Node) RemoveChild(t *Node) {
	pos := -1;
	for i, v := range this.Children {
		if v == t {
			pos = i;
			break;
		}
	}

	if pos == -1 {
		return
	}
	slice := make([]*Node, len(this.Children)-1);

	idx := 0;
	for i, v := range this.Children {
		if i != pos {
			slice[idx] = v;
			idx++;
		}
	}

	t.Parent = nil;
}

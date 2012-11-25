// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package xmlx

import "testing"

func TestLoadLocal(t *testing.T) {
	doc := New()

	if err := doc.LoadFile("test.xml", nil); err != nil {
		t.Error(err.Error())
		return
	}

	if len(doc.Root.Children) == 0 {
		t.Errorf("Root node has no children.")
		return
	}
}

func TestWildcard(t *testing.T) {
	doc := New()

	if err := doc.LoadFile("test2.xml", nil); err != nil {
		t.Error(err.Error())
		return
	}

	list := doc.SelectNodes("ns", "*")

	if len(list) != 1 {
		t.Errorf("Wrong number of child elements. Expected 1, got %d.", len(list))
		return
	}
}

func TestWildcardRecursive(t *testing.T) {
	doc := New()

	if err := doc.LoadFile("test2.xml", nil); err != nil {
		t.Error(err.Error())
		return
	}

	list := doc.SelectNodesRecursive("ns", "*")

	if len(list) != 7 {
		t.Errorf("Wrong number of child elements. Expected 7, got %d.", len(list))
		return
	}
}

func _TestLoadRemote(t *testing.T) {
	doc := New()

	if err := doc.LoadUri("http://blog.golang.org/feeds/posts/default", nil); err != nil {
		t.Error(err.Error())
		return
	}

	if len(doc.Root.Children) == 0 {
		t.Errorf("Root node has no children.")
		return
	}
}

func TestSave(t *testing.T) {
	doc := New()

	if err := doc.LoadFile("test.xml", nil); err != nil {
		t.Errorf("LoadFile(): %s", err)
		return
	}

	IndentPrefix = "\t"
	if err := doc.SaveFile("test1.xml"); err != nil {
		t.Errorf("SaveFile(): %s", err)
		return
	}
}

func TestNodeSearch(t *testing.T) {
	doc := New()

	if err := doc.LoadFile("test1.xml", nil); err != nil {
		t.Errorf("LoadFile(): %s", err)
		return
	}

	if node := doc.SelectNode("", "item"); node == nil {
		t.Errorf("SelectNode(): No node found.")
		return
	}

	nodes := doc.SelectNodes("", "item")
	if len(nodes) == 0 {
		t.Errorf("SelectNodes(): no nodes found.")
		return
	}
}

type Image struct {
	Title       string `xml:"title"`
	Url         string `xml:"url"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Width       int    `xml:"width"`
	Height      int    `xml:"height"`
}

func TestUnmarshal(t *testing.T) {
	doc := New()
	err := doc.LoadFile("test1.xml", nil)

	if err != nil {
		t.Errorf("LoadFile(): %s", err)
		return
	}

	node := doc.SelectNode("", "image")
	if node == nil {
		t.Errorf("SelectNode(): No node found.")
		return
	}

	var img Image
	if err = node.Unmarshal(&img); err != nil {
		t.Errorf("Unmarshal(): %s", err)
		return
	}

	if img.Title != "WriteTheWeb" {
		t.Errorf("Image.Title has incorrect value. Got '%s', expected 'WriteTheWeb'.", img.Title)
		return
	}
}

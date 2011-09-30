// This work is subject to the CC0 1.0 Universal (CC0 1.0) Public Domain Dedication
// license. Its contents can be found at:
// http://creativecommons.org/publicdomain/zero/1.0/

package xmlx

import "testing"

func TestLoadLocal(t *testing.T) {
	doc := New()

	if err := doc.LoadFile("test.xml"); err != nil {
		t.Error(err.String())
		return
	}

	if len(doc.Root.Children) == 0 {
		t.Errorf("Root node has no children.")
		return
	}
}

func TestWildcard(t *testing.T) {
	doc := New()

	if err := doc.LoadFile("test2.xml"); err != nil {
		t.Error(err.String())
		return
	}

	list := doc.SelectNodes("ns", "*")

	if len(list) != 1 {
		t.Errorf("Wrong number of child elements. Expected 4, got %d.", len(list))
		return
	}
}

func _TestLoadRemote(t *testing.T) {
	doc := New()

	if err := doc.LoadUri("http://blog.golang.org/feeds/posts/default"); err != nil {
		t.Error(err.String())
		return
	}

	if len(doc.Root.Children) == 0 {
		t.Errorf("Root node has no children.")
		return
	}
}

func TestSave(t *testing.T) {
	doc := New()

	if err := doc.LoadFile("test.xml"); err != nil {
		t.Errorf("LoadFile(): %s", err)
		return
	}

	if err := doc.SaveFile("test1.xml"); err != nil {
		t.Errorf("SaveFile(): %s", err)
		return
	}
}

func TestNodeSearch(t *testing.T) {
	doc := New()

	if err := doc.LoadFile("test1.xml"); err != nil {
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
	Title       string
	Url         string
	Link        string
	Width       string
	Height      string
	Description string
}

func TestUnmarshal(t *testing.T) {
	doc := New()
	err := doc.LoadFile("test1.xml")

	if err != nil {
		t.Errorf("LoadFile(): %s", err)
		return
	}

	node := doc.SelectNode("", "image")
	if node == nil {
		t.Errorf("SelectNode(): No node found.")
		return
	}

	img := Image{}
	if err = node.Unmarshal(&img); err != nil {
		t.Errorf("Unmarshal(): %s", err)
		return
	}

	if img.Title != "WriteTheWeb" {
		t.Errorf("Image.Title has incorrect value. Got '%s', expected 'WriteTheWeb'.", img.Title)
		return
	}
}

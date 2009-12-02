package xmlx

import "testing"

func TestLoadLocal(t *testing.T) {
	doc := New();
	err := doc.LoadFile("test.xml");

	if err != nil {
		t.Errorf("%s", err);
		return;
	}

	if len(doc.Root.Children) == 0 {
		t.Errorf("Root node has no children.");
		return;
	}
}

func TestLoadRemote(t *testing.T) {
	doc := New();
	err := doc.LoadUri("http://tldp.org/authors/template/Sample-HOWTO.xml");

	if err != nil {
		t.Errorf("%s", err);
		return;
	}

	if len(doc.Root.Children) == 0 {
		t.Errorf("Root node has no children.");
		return;
	}
}

func TestSave(t *testing.T) {
	doc := New();
	err := doc.LoadFile("test.xml");

	if err != nil {
		t.Errorf("LoadFile(): %s", err);
		return;
	}

	err = doc.SaveFile("test1.xml");
	if err != nil {
		t.Errorf("SaveFile(): %s", err);
		return;
	}
}

func TestNodeSearch(t *testing.T) {
	doc := New();
	err := doc.LoadFile("test.xml");

	if err != nil {
		t.Errorf("LoadFile(): %s", err);
		return;
	}

	node := doc.SelectNode("", "item");
	if node == nil {
		t.Errorf("SelectNode(): No node found.");
		return;
	}

	nodes := doc.SelectNodes("", "item");
	if len(nodes) == 0 {
		t.Errorf("SelectNodes(): no nodes found.");
		return;
	}
}

type Image struct {
	Title		string;
	Url		string;
	Link		string;
	Width		string;
	Height		string;
	Description	string;
}

func TestUnmarshal(t *testing.T) {
	doc := New();
	err := doc.LoadFile("test.xml");

	if err != nil {
		t.Errorf("LoadFile(): %s", err);
		return;
	}

	node := doc.SelectNode("", "image");
	if node == nil {
		t.Errorf("SelectNode(): No node found.");
		return;
	}

	img := Image{};
	err = node.Unmarshal(&img);
	if err != nil {
		t.Errorf("Unmarshal(): %s", err);
		return;
	}

	if img.Title != "WriteTheWeb" {
		t.Errorf("Image.Title has incorrect value. Got '%s', expected 'WriteTheWeb'.", img.Title);
		return;
	}
}



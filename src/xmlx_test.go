package xmlx

import "testing"

func TestDoc(t *testing.T) {
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


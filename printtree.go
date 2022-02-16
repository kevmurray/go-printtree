package printtree

type Tree struct {
	Name string
	Parent *Tree
	Sibling *Tree
	Child *Tree
	context *treeContext
}

type treeContext struct {
}

func (t *Tree)Indent() *Tree {
	if t.Child == nil {
		t.Child = 
	}
	return t.Child
}

func newTree(context *treeContext, parent *Tree) *Tree {
	return &Tree{
		context: context
		Parent: parent
	}
}


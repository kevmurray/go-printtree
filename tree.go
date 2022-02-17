package printtree

import (
	"fmt"
	"sort"
	"strings"
)

// A simple tree that prints heirarchically. Create a new tree with "NewTree()", then add
// children to it by calling Add(name), Addf(format, args...) or AddAll(name,...). Each of these
// will return the child(ren) that were added so you can continue to recursively add more nodes
// to the tree.
//
// When done, you can call Print() or PrintStyle(style) to get a string representing the tree
// that has traditional ascii-like heirarchy markup.
//
// Example:
//    root := NewTree()
//    colorTree := root.Add("Monitors")
//    monoTree := colorTree.Add("Monocrome")
//    monoTree.Add("Old School").AddAll("black", "green")
//    monoTree.Add("Contemporary").AddAll("black", "white")
//    colorTree.Add("Color").AddAll("red", "green", "blue")
//    fmt.Println(root.PrintStyle(Heavy))
//
// Which will print
//    Monitors
//    ├── Monocrome
//    │   ├── Old School
//    │   │   ├── black
//    │   │   ╰── green
//    │   ╰── Contemporary
//    │       ├── black
//    │       ╰── white
//    ╰── Color
//        ├── red
//        ├── green
//        ╰── blue
//
// You can use PrintStyle() to access other styles
//    |-- ASCII
//    ├── Box
//    ┣━━ Heavy
//    |-ASCIINarrow
//    ├ BoxNarrow
//    ┣ HeavyNarrow
//        WhiteSpace
//    ··· Dots
//    →   Arrows

// TreeStyle is the markup style of the tree
type TreeStyle int

const (
	ASCII       TreeStyle = 0
	Box         TreeStyle = 1
	Heavy       TreeStyle = 2
	ASCIINarrow TreeStyle = 3
	BoxNarrow   TreeStyle = 4
	HeavyNarrow TreeStyle = 5
	WhiteSpace  TreeStyle = 6
	Dots        TreeStyle = 7
	Arrows      TreeStyle = 8
)

// for the different branch types
const (
	childBranch     int = 0
	bypassBranch    int = 1
	lastChildBranch int = 2
	noBranch        int = 3
)

type Tree struct {
	Label    string
	Children []*Tree
}

// NewTree returns a new tree node that has no label. This is the root of a tree that you can
// work with later. When you print this tree, this root node will not be printed, only the
// children of this root node. This allows for creating multiple visual roots like
//   root := NewTree()
//   root.Add("1. First").Add("a. Alpha")
//   root.Add("2. second").Add("b. Bravo")
// will print the following
//   1. First
//   '-- a. Alpha
//   2. Second
//   '-- b. Bravo
func NewTree() *Tree {
	return &Tree{}
}

// Add creates a new child Tree (at the end of the child list) in this tree and returns that child Tree
func (tree *Tree) Add(childName string) *Tree {
	childTree := Tree{
		Label: childName,
	}
	tree.Children = append(tree.Children, &childTree)
	return &childTree
}

// AddAll creates multiple children at one time in the order they are listed
func (tree *Tree) AddAll(childrenNames ...string) []*Tree {
	children := make([]*Tree, 0, len(childrenNames))
	for _, childName := range childrenNames {
		child := tree.Add(childName)
		children = append(children, child)
	}
	return children
}

// Add creates a new child Tree (at the end of the child list) in this tree and returns that child Tree
func (tree *Tree) Addf(label string, a ...interface{}) *Tree {
	return tree.Add(fmt.Sprintf(label, a...))
}

// AddTree adds another tree as a child of this tree. If the other tree has no label, then it is
// assumed to be a root node, and all it's children will be added. If it does have a label, then
// it will be added as a child
//
// WARNING: It is up to you not to create infinite tree loops
func (tree *Tree) AddTree(other *Tree) {
	if other.Label == "" {
		// this is a root tree, copy all it's children
		tree.Children = append(tree.Children, other.Children...)
	} else {
		// tree branch, add the branch to this tree
		tree.Children = append(tree.Children, other)
	}
}

// Sort sorts the children of this tree by the labels
func (tree *Tree) Sort() {
	sort.SliceStable(tree.Children, func(i, j int) bool {
		return tree.Children[i].Label < tree.Children[j].Label
	})
}

// DeepSort sorts the children of this tree and all sub-trees by the labels
func (tree *Tree) DeepSort() {
	tree.Sort()
	for index := range tree.Children {
		tree.Children[index].DeepSort()
	}
}

// String returns a string representation of this tree indented with whitespace
func (tree *Tree) String() string {
	return tree.PrintStyle(WhiteSpace)
}

// Print returns a string which is this tree in a heirarchical format
func (tree *Tree) Print() string {
	return tree.PrintStyle(Box)
}
func (tree *Tree) PrintStyle(style TreeStyle) string {
	dict := [][]string{
		{"|-- ", "|   ", "'-- ", "    "},
		{"├── ", "│   ", "╰── ", "    "},
		{"┣━━ ", "┃   ", "┗━━ ", "    "},
		{"|-", "| ", "'-", "  "},
		{"├ ", "│ ", "╰ ", "  "},
		{"┣ ", "┃ ", "┗ ", "  "},
		{"    ", "    ", "    ", "    "},
		{"··· ", "····", "··· ", "····"},
		{"→   ", "    ", "→   ", "    "},
	}[style]
	buf := strings.Builder{}
	tree.print(&buf, true, "", dict)
	return buf.String()
}

// print is the internal, recursive hook for printing the tree
func (tree *Tree) print(buf *strings.Builder, isRoot bool, padding string, dict []string) {
	for index := range tree.Children {
		child := tree.Children[index]
		for lineIndex, line := range strings.Split(child.Label, "\n") {
			if lineIndex == 0 {
				// first line of the block, normal padding
				buf.WriteString(padding + tree.labelPadding(isRoot, index, dict) + line + "\n")
			} else {
				// subsequent lines of the block, flow padding
				buf.WriteString(padding + tree.flowPadding(isRoot, index, dict) + line + "\n")
			}
		}
		child.print(buf, false, padding+tree.flowPadding(isRoot, index, dict), dict)
	}
}

func (tree *Tree) labelPadding(isRoot bool, index int, dict []string) string {
	switch {
	case isRoot:
		return ""
	case index == len(tree.Children)-1:
		return dict[lastChildBranch]
	}
	return dict[childBranch]
}

func (tree *Tree) flowPadding(isRoot bool, index int, dict []string) string {
	switch {
	case isRoot:
		return ""
	case index == len(tree.Children)-1:
		return dict[noBranch]
	}
	return dict[bypassBranch]
}

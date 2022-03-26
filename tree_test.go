package printtree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddBranch(t *testing.T) {
	// given
	tree := NewTree()

	// when
	root := tree.AddBranch("label")

	// then
	assert.Equal(t, "label", root.Label)
}

func ExampleTree_AddBranch() {
	tree := NewTree()
	fruit := tree.AddBranch("Fruit")
	fruit.AddBranch("Lemmon")
	fruit.AddBranch("Orange").AddBranch("Mandarin")
	fruit.AddBranch("Lime")
	fmt.Print(tree.Print())
	// Output:
	// Fruit
	// ├── Lemmon
	// ├── Orange
	// │   ╰── Mandarin
	// ╰── Lime
}

func TestBranchFormatting(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranchf("label %s", "one")
	assert.Equal(t, "label one", root.Label)
}

func ExampleTree_AddBranchf() {
	tree := NewTree()
	users := tree.AddBranch("/Users (disk space)")
	users.AddBranchf("%-12s (%6dMB)", "lister", 12)
	users.AddBranchf("%-12s (%6dMB)", "kryten", 167)
	users.AddBranchf("%-12s (%6dMB)", "rimmer", 876252)
	fmt.Print(tree.Print())
	// Output:
	// /Users (disk space)
	// ├── lister       (    12MB)
	// ├── kryten       (   167MB)
	// ╰── rimmer       (876252MB)
}

func TestAddBranches(t *testing.T) {
	tree := NewTree()
	roots := tree.AddBranches("alfa", "bravo", "charlie")
	assert.Equal(t, "alfa", roots[0].Label)
	assert.Equal(t, "bravo", roots[1].Label)
	assert.Equal(t, "charlie", roots[2].Label)
}

func ExampleTree_AddBranches() {
	tree := NewTree()
	fruit := tree.AddBranch("Fruit")
	fruits := fruit.AddBranches("Lemmon", "Orange", "Lime")
	fruits[1].AddBranch("Mandarin")
	fmt.Print(tree.Print())
	// Output:
	// Fruit
	// ├── Lemmon
	// ├── Orange
	// │   ╰── Mandarin
	// ╰── Lime
}

func TestAddTreeAsBranch_Tree(t *testing.T) {
	tree := NewTree()
	tree.AddBranch("original")

	otherTree := NewTree()
	otherTree.AddBranches("extra", "crispy")

	// when
	tree.AddTreeAsBranch(otherTree)

	// then
	assert.Len(t, tree.Branches, 3)
	assert.Equal(t, "original", tree.Branches[0].Label)
	assert.Equal(t, "extra", tree.Branches[1].Label)
	assert.Equal(t, "crispy", tree.Branches[2].Label)
}

func TestAddTreeAsBranch_Branch(t *testing.T) {
	tree := NewTree()
	tree.AddBranch("original")

	otherTree := NewTree()
	branches := otherTree.AddBranches("extra", "crispy")

	// when
	tree.AddTreeAsBranch(branches[1])

	// then
	assert.Len(t, tree.Branches, 2)
	assert.Equal(t, "original", tree.Branches[0].Label)
	assert.Equal(t, "crispy", tree.Branches[1].Label)
}

func TestSort(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranch("root")
	branchBlack := root.AddBranches("sphinx", "of", "black", "quartz")[2]
	branchBlack.AddBranches("obsidian", "midnight", "dark")

	root.Sort()

	// top level should be sorted
	assert.Equal(t, "black", root.Branches[0].Label)
	assert.Equal(t, "of", root.Branches[1].Label)
	assert.Equal(t, "quartz", root.Branches[2].Label)
	assert.Equal(t, "sphinx", root.Branches[3].Label)

	// next level should not be sorted
	assert.Equal(t, "obsidian", root.Branches[0].Branches[0].Label)
	assert.Equal(t, "midnight", root.Branches[0].Branches[1].Label)
	assert.Equal(t, "dark", root.Branches[0].Branches[2].Label)
}

func TestDeepSort(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranch("root")
	branchBlack := root.AddBranches("sphinx", "of", "black", "quartz")[2]
	branchBlack.AddBranches("obsidian", "midnight", "dark")

	root.DeepSort()

	// check top level
	assert.Equal(t, "black", root.Branches[0].Label)
	assert.Equal(t, "of", root.Branches[1].Label)
	assert.Equal(t, "quartz", root.Branches[2].Label)
	assert.Equal(t, "sphinx", root.Branches[3].Label)

	// check next level
	assert.Equal(t, "dark", root.Branches[0].Branches[0].Label)
	assert.Equal(t, "midnight", root.Branches[0].Branches[1].Label)
	assert.Equal(t, "obsidian", root.Branches[0].Branches[2].Label)
}

func TestCustomSort(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranch("root")
	branchBlack := root.AddBranches("sphinx", "of", "black", "quartz")[2]
	branchBlack.AddBranches("obsidian", "midnight", "dark")

	// sort by length of label
	root.SortCustom(func(b1, b2 *Tree) bool {
		return len(b1.Label) < len(b2.Label)
	})

	// top level should be sorted
	assert.Equal(t, "of", root.Branches[0].Label)
	assert.Equal(t, "black", root.Branches[1].Label)
	assert.Equal(t, "sphinx", root.Branches[2].Label)
	assert.Equal(t, "quartz", root.Branches[3].Label)

	// next level should not be sorted
	assert.Equal(t, "obsidian", root.Branches[1].Branches[0].Label)
	assert.Equal(t, "midnight", root.Branches[1].Branches[1].Label)
	assert.Equal(t, "dark", root.Branches[1].Branches[2].Label)
}

func TestDeepCustomSort(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranch("root")
	branchBlack := root.AddBranches("sphinx", "of", "black", "quartz")[2]
	branchBlack.AddBranches("obsidian", "midnight", "dark")

	// sort by length of label
	root.DeepSortCustom(func(b1, b2 *Tree) bool {
		return len(b1.Label) < len(b2.Label)
	})

	// check top level
	assert.Equal(t, "of", root.Branches[0].Label)
	assert.Equal(t, "black", root.Branches[1].Label)
	assert.Equal(t, "sphinx", root.Branches[2].Label)
	assert.Equal(t, "quartz", root.Branches[3].Label)

	// check next level
	assert.Equal(t, "dark", root.Branches[1].Branches[0].Label)
	assert.Equal(t, "obsidian", root.Branches[1].Branches[1].Label)
	assert.Equal(t, "midnight", root.Branches[1].Branches[2].Label)
}

func TestRoot(t *testing.T) {
	root := NewTree()
	root.AddBranches("1", "2", "3")

	t.Logf("TREE:\n%s", root.String())
	t.Logf("OUT :\n%s", root.PrintStyle(BoxStyle))

	result := root.PrintStyle(ASCIIStyle)
	assert.Equal(t, "1\n2\n3\n", result)
}

func TestChildren(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranch("1")
	root.AddBranches("a", "b")

	t.Logf("TREE:\n%s", tree.String())
	t.Logf("OUT :\n%s", tree.PrintStyle(BoxStyle))

	result := tree.PrintStyle(ASCIIStyle)
	assert.Equal(t, `1
|-- a
'-- b
`, result)
}

func TestGrandChildren(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranch("1")
	branchA := root.AddBranch("a")
	branchA.AddBranches("i", "ii")

	t.Logf("TREE:\n%s", tree.String())
	t.Logf("OUT :\n%s", tree.PrintStyle(BoxStyle))

	result := tree.PrintStyle(ASCIIStyle)
	assert.Equal(t, `1
'-- a
    |-- i
    '-- ii
`, result)
}

func TestMultipleRoots(t *testing.T) {
	tree := NewTree()
	root1 := tree.AddBranch("1")
	root1.AddBranches("a", "b")
	root2 := tree.AddBranch("2")
	root2.AddBranches("a", "b")

	t.Logf("TREE:\n%s", tree.String())
	t.Logf("OUT :\n%s", tree.PrintStyle(BoxStyle))

	result := tree.PrintStyle(ASCIIStyle)
	assert.Equal(t, `1
|-- a
'-- b
2
|-- a
'-- b
`, result)
}

func TestNephews(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranch("1")
	branchA := root.AddBranch("a")
	branchA.AddBranch("i")
	branchA.AddBranch("ii")
	branchB := root.AddBranch("b")
	branchB.AddBranch("i")
	branchB.AddBranch("ii")

	t.Logf("TREE:\n%s", tree.String())
	t.Logf("OUT :\n%s", tree.PrintStyle(BoxStyle))

	result := tree.PrintStyle(ASCIIStyle)
	assert.Equal(t, `1
|-- a
|   |-- i
|   '-- ii
'-- b
    |-- i
    '-- ii
`, result)
}

func TestComplexTree(t *testing.T) {
	// replicate a complex subdirectory:
	//   vda
	//   |-- api
	//   |   |-- auth.go
	//   |   |-- engine.go
	//   |   '-- graphql.go
	//   |-- clair
	//   |   '-- api
	//   |       '-- v1
	//   |           |-- models.go
	//   |           '-- readme.md
	//   |-- engine-run-test
	//   |   '-- engine-run.go
	//   |-- errors.go
	//   '-- go-config
	//       |-- config.go
	//       |-- config_test.go
	//       '-- testdata
	//           |-- config_test.json
	//           '-- config_test.yml

	tree := NewTree()
	branchVDA := tree.AddBranch("vda")
	vdaChildren := branchVDA.AddBranches("api", "clair", "engine-run-test", "errors.go", "go-config")
	// api
	vdaChildren[0].AddBranches("auth.go", "engine.go", "graphql.go")
	// clair
	vdaChildren[1].AddBranch("api").AddBranch("v1").AddBranches("models.go", "readme.md")
	// engine-run-test
	vdaChildren[2].AddBranch("engine-run.go")
	// go-config
	vdaChildren[4].AddBranches("config.go", "config_test.go", "testdata")[2].AddBranches("config_test.json", "config_test.yaml")

	// assert structural data
	t.Logf("OUT :\n%s", tree.PrintStyle(ASCIIStyle))
	// the result was confirmed visually (from the above log statement) then copied and pasted
	result := tree.PrintStyle(ASCIIStyle)
	assert.Equal(t, `vda
|-- api
|   |-- auth.go
|   |-- engine.go
|   '-- graphql.go
|-- clair
|   '-- api
|       '-- v1
|           |-- models.go
|           '-- readme.md
|-- engine-run-test
|   '-- engine-run.go
|-- errors.go
'-- go-config
    |-- config.go
    |-- config_test.go
    '-- testdata
        |-- config_test.json
        '-- config_test.yaml
`, result)

	// assert number list
	t.Logf("OUT :\n%s", tree.PrintStyle(OrderedStyle))
	// the result was confirmed visually (from the above log statement) then copied and pasted
	result = tree.PrintStyle(OrderedStyle)
	assert.Equal(t, `vda
 1. api
     a. auth.go
     b. engine.go
     c. graphql.go
 2. clair
     a. api
         i. v1
             A. models.go
             B. readme.md
 3. engine-run-test
     a. engine-run.go
 4. errors.go
 5. go-config
     a. config.go
     b. config_test.go
     c. testdata
         i. config_test.json
        ii. config_test.yaml
`, result)
}

func TestMultilineBranches(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranch("1")
	root.AddBranch("a\nalfa\nalpha\nable")
	root.AddBranch("b")

	t.Logf("TREE:\n%s", tree.String())
	t.Logf("OUT :\n%s", tree.PrintStyle(BoxStyle))

	result := tree.PrintStyle(ASCIIStyle)
	assert.Equal(t, `1
|-- a
|   alfa
|   alpha
|   able
'-- b
`, result)
}

func TestStructuralStyles(t *testing.T) {
	// create minimal tree that contains every style
	tree := NewTree()
	root := tree.AddBranch("1")
	branchA := root.AddBranch("a")
	branchA.AddBranch("i")
	root.AddBranch("b\nB")

	cases := []struct {
		style    TreeStyle
		expected string
	}{
		{style: ASCIIStyle, expected: `1
|-- a
|   '-- i
'-- b
    B
`},
		{style: BoxStyle, expected: `1
├── a
│   ╰── i
╰── b
    B
`},
		{style: BoxBoldStyle, expected: `1
┣━━ a
┃   ┗━━ i
┗━━ b
    B
`},
		{style: ASCIINarrowStyle, expected: `1
|-a
| '-i
'-b
  B
`},
		{style: BoxNarrowStyle, expected: `1
├ a
│ ╰ i
╰ b
  B
`},
		{style: BoxBoldNarrowStyle, expected: `1
┣ a
┃ ┗ i
┗ b
  B
`},
	}

	for index, tc := range cases {
		t.Logf("Style %d:\n%s", tc.style, tree.PrintStyle(tc.style))
		assert.Equal(t, tc.expected, tree.PrintStyle(tc.style), "Test case %d failed", index)
	}
}

func TestReplaceNumberPlaceholder(t *testing.T) {
	tree := NewTree()

	cases := []struct {
		s           string
		placeholder string
		actual      string
		expected    string
	}{
		{"1", "1", "x", "x"},
		{"1", "1", "xxx", "xxx"},
		{"  1", "1", "x", "  x"},
		{"  1", "1", "xx", " xx"},
		{"  1", "1", "xxx", "xxx"},
		{"  1", "1", "xxxx", "xxxx"},
		{"(  1)", "1", "x", "(  x)"},
		{"(  1)", "1", "xx", "( xx)"},
		{"(  1)", "1", "xxx", "(xxx)"},
		{"(  1)", "1", "xxxx", "(xxxx)"},
	}

	for index, tc := range cases {
		assert.Equal(t, tc.expected, tree.replaceNumberPlaceholder(tc.s, tc.placeholder, tc.actual), "test case %d failed", index)
	}
}

func TestReplaceNumberListMarkup(t *testing.T) {
	tree := NewTree()

	cases := []struct {
		markup   string
		index    int
		expected string
	}{
		{"1", 9, "9"},
		{"1", 19, "19"},
		{" 1)", 29, "29)"},
		{"  1.", 99, " 99."},
		{" (1)", 1239, " (1239)"},

		{"a", 0, "-"},
		{"a", 1, "a"},
		{"a", 9, "i"},
		{"a", 19, "s"},
		{"a", 26, "z"},
		{"a", 27, "aa"},
		{" a)", 29, "ac)"},
		{"  a.", 99, " cu."},
		{" (a)", 1239, " (auq)"},

		{"A", 0, "-"},
		{"A", 1, "A"},
		{"A", 9, "I"},
		{"A", 19, "S"},
		{"A", 26, "Z"},
		{"A", 27, "AA"},
		{" A)", 29, "AC)"},
		{"  A.", 99, " CU."},
		{" (A)", 1239, " (AUQ)"},

		{"i", 0, "-"},
		{" i", 1, " i"},
		{" i", 2, "ii"},
		{" i", 3, "iii"},
		{" i", 4, "iv"},
		{" i", 5, " v"},
		{"i", 9, "ix"},
		{"i", 19, "xix"},
		{" i)", 29, "xxix)"},
		{" i)", 33, "xxxiii)"},
		{"  i.", 99, "xcix."},
		{" (i)", 1239, " (mccxxxix)"},

		{"I", 0, "-"},
		{" I", 1, " I"},
		{" I", 2, "II"},
		{" I", 3, "III"},
		{" I", 4, "IV"},
		{" I", 5, " V"},
		{"I", 9, "IX"},
		{"I", 19, "XIX"},
		{" I)", 29, "XXIX)"},
		{" I)", 33, "XXXIII)"},
		{"  I.", 99, "XCIX."},
		{" (I)", 1239, " (MCCXXXIX)"},
	}

	for index, tc := range cases {
		assert.Equal(t, tc.expected, tree.replaceNumberListMarkup(tc.markup, tc.index), "test case %d failed", index)
	}
}

func TestListStyles(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranch("Monitors")
	monoTree := root.AddBranch("Monochrome")
	monoTree.AddBranch("Old School").AddBranches("black", "green")
	monoTree.AddBranch("Contemporary").AddBranches("black", "white")
	root.AddBranch("Color").AddBranches("red", "green", "blue")

	cases := []struct {
		style    TreeStyle
		expected string
	}{
		{style: WhiteSpaceStyle, expected: `Monitors
    Monochrome
        Old School
            black
            green
        Contemporary
            black
            white
    Color
        red
        green
        blue
`},
		{style: ASCIIBulletStyle, expected: `Monitors
* Monochrome
  + Old School
    - black
    - green
  + Contemporary
    - black
    - white
* Color
  + red
  + green
  + blue
`},
		{style: BulletStyle, expected: `Monitors
● Monochrome
  ○ Old School
    ■ black
    ■ green
  ○ Contemporary
    ■ black
    ■ white
● Color
  ○ red
  ○ green
  ○ blue
`},
		{style: OrderedStyle, expected: `Monitors
 1. Monochrome
     a. Old School
         i. black
        ii. green
     b. Contemporary
         i. black
        ii. white
 2. Color
     a. red
     b. green
     c. blue
`},
		{style: NumberStyle, expected: `Monitors
 1. Monochrome
     1. Old School
         1. black
         2. green
     2. Contemporary
         1. black
         2. white
 2. Color
     1. red
     2. green
     3. blue
`},
		{style: AlphaStyle, expected: `Monitors
 a. Monochrome
     a. Old School
         a. black
         b. green
     b. Contemporary
         a. black
         b. white
 b. Color
     a. red
     b. green
     c. blue
`},
		{style: AlphaUCStyle, expected: `Monitors
 A. Monochrome
     A. Old School
         A. black
         B. green
     B. Contemporary
         A. black
         B. white
 B. Color
     A. red
     B. green
     C. blue
`},
		{style: RomanStyle, expected: `Monitors
   i. Monochrome
         i. Old School
               i. black
              ii. green
        ii. Contemporary
               i. black
              ii. white
  ii. Color
         i. red
        ii. green
       iii. blue
`},
		{style: RomanUCStyle, expected: `Monitors
   I. Monochrome
         I. Old School
               I. black
              II. green
        II. Contemporary
               I. black
              II. white
  II. Color
         I. red
        II. green
       III. blue
`},
	}

	for index, tc := range cases {
		t.Logf("Style %d:\n%s", tc.style, tree.PrintStyle(tc.style))
		assert.Equal(t, tc.expected, tree.PrintStyle(tc.style), "Test case %d failed", index)
	}
}

func TestStyles_ListCycleBullets(t *testing.T) {
	tree := NewTree()
	tree.AddBranch("One").
		AddBranch("Two").
		AddBranch("Three").
		AddBranch("Four").
		AddBranch("Five").
		AddBranch("Six")

	result := tree.PrintStyle(ASCIIBulletStyle)
	assert.Equal(t, `One
* Two
  + Three
    - Four
      * Five
        + Six
`, result)

}

func TestIllegalStyle(t *testing.T) {
	tree := NewTree()
	root := tree.AddBranch("1")
	branchA := root.AddBranch("a")
	branchA.AddBranch("i")
	root.AddBranch("b\nB")

	// illegal scaffold defaults
	result := tree.PrintStyle(TreeStyle(999999))
	assert.Equal(t, `1
├── a
│   ╰── i
╰── b
    B
`, result)
}

func ExampleAddStructuralStyle() {
	tree := NewTree()
	root := tree.AddBranch("Mom")
	branchA := root.AddBranch("Myself")
	branchA.AddBranch("Child")
	root.AddBranch("Sister\nBrother")

	myStyle := AddStructuralStyle(">- ", "*- ", "}  ", "...")

	fmt.Print(tree.PrintStyle(myStyle))
	// Output:
	// Mom
	// >- Myself
	// }  *- Child
	// *- Sister
	// ...Brother
}

func ExampleAddListStyle() {
	tree := NewTree()
	root := tree.AddBranch("Mom")
	branchA := root.AddBranch("Myself")
	branchA.AddBranches("Child1", "Child2")
	root.AddBranch("Sister\nBrother")

	myStyle := AddListStyle("   ", "(1)", "(•)", "(i)")

	fmt.Print(tree.PrintStyle(myStyle))
	// Output:
	// Mom
	// (1)Myself
	//    (•)Child1
	//    (•)Child2
	// (2)Sister
	//    Brother
}

func TestDepth(t *testing.T) {
	tree := NewTree()
	assert.Equal(t, 0, tree.Depth())

	root := tree.AddBranch("root")
	assert.Equal(t, 1, tree.Depth())
	assert.Equal(t, 0, root.Depth())

	branchA := root.AddBranch("a")
	branchB := root.AddBranch("b")
	branchB.AddBranches("i", "ii", "iii")
	assert.Equal(t, 3, tree.Depth())
	assert.Equal(t, 2, root.Depth())
	assert.Equal(t, 0, branchA.Depth())
	assert.Equal(t, 1, branchB.Depth())
}

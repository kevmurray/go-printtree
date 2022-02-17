package printtree

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type treeTestSuite struct {
	suite.Suite
}

func TestTree(t *testing.T) {
	suite.Run(t, new(treeTestSuite))
}

func (t *treeTestSuite) SetupSuite() {
}

func (t *treeTestSuite) SetupTest() {
}

func (t *treeTestSuite) AfterTest(suiteName, testName string) {
}

func TestAdd(t *testing.T) {
	// given
	tree := NewTree()

	// when
	root := tree.Add("label")

	// then
	if root.Label != "label" {
		t.Errorf("TestAdd: Expected '%s', got '%s'","label", root.Label)
	}
}

func (t *treeTestSuite) TestAddf() {
	tree := NewTree()
	root := tree.Addf("label %s", "one")
	t.Equal("label one", root.Label)
}

func (t *treeTestSuite) TestAddAll() {
	tree := NewTree()
	roots := tree.AddAll("alfa", "bravo", "charlie")
	t.Equal("alfa", roots[0].Label)
	t.Equal("bravo", roots[1].Label)
	t.Equal("charlie", roots[2].Label)
}

func (t *treeTestSuite) TestAddTree_Tree() {
	tree := NewTree()
	tree.Add("original")

	otherTree := NewTree()
	otherTree.AddAll("extra", "crispy")

	// when
	tree.AddTree(otherTree)

	// then
	t.Len(tree.Children, 3)
	t.Equal("original", tree.Children[0].Label)
	t.Equal("extra", tree.Children[1].Label)
	t.Equal("crispy", tree.Children[2].Label)
}

func (t *treeTestSuite) TestAddTree_Branch() {
	tree := NewTree()
	tree.Add("original")

	otherTree := NewTree()
	branches := otherTree.AddAll("extra", "crispy")

	// when
	tree.AddTree(branches[1])

	// then
	t.Len(tree.Children, 2)
	t.Equal("original", tree.Children[0].Label)
	t.Equal("crispy", tree.Children[1].Label)
}

func (t *treeTestSuite) TestSort() {
	tree := NewTree()
	tRoot := tree.Add("root")
	tRoot.AddAll("sphinx", "of", "black", "quartz")
	tRoot.Sort()

	t.Equal("black", tRoot.Children[0].Label)
	t.Equal("of", tRoot.Children[1].Label)
	t.Equal("quartz", tRoot.Children[2].Label)
	t.Equal("sphinx", tRoot.Children[3].Label)
}

func (t *treeTestSuite) TestSortAll() {
	tree := NewTree()
	tRoot := tree.Add("root")
	tRoot.AddAll("sphinx", "of", "black", "quartz")
	tRoot.Children[2].AddAll("obsidian", "midnight", "dark")

	tRoot.DeepSort()

	// check top level
	t.Equal("black", tRoot.Children[0].Label)
	t.Equal("of", tRoot.Children[1].Label)
	t.Equal("quartz", tRoot.Children[2].Label)
	t.Equal("sphinx", tRoot.Children[3].Label)

	// check next level
	t.Equal("dark", tRoot.Children[0].Children[0].Label)
	t.Equal("midnight", tRoot.Children[0].Children[1].Label)
	t.Equal("obsidian", tRoot.Children[0].Children[2].Label)
}

func (t *treeTestSuite) TestRoot() {
	rootTree := NewTree()
	rootTree.Add("1")
	rootTree.Add("2")
	rootTree.Add("3")

	t.T().Logf("TREE:\n%s", rootTree.String())
	t.T().Logf("OUT :\n%s", rootTree.PrintStyle(Box))

	result := rootTree.PrintStyle(ASCII)
	t.Equal("1\n2\n3\n", result)
}

func (t *treeTestSuite) TestChildren() {
	rootTree := NewTree()
	oneTree := rootTree.Add("1")
	oneTree.Add("a")
	oneTree.Add("b")

	t.T().Logf("TREE:\n%s", rootTree.String())
	t.T().Logf("OUT :\n%s", rootTree.PrintStyle(Box))

	result := rootTree.PrintStyle(ASCII)
	t.Equal(`1
|-- a
'-- b
`, result)
}

func (t *treeTestSuite) TestGrandChildren() {
	rootTree := NewTree()
	oneTree := rootTree.Add("1")
	aTree := oneTree.Add("a")
	aTree.Add("i")
	aTree.Add("ii")

	t.T().Logf("TREE:\n%s", rootTree.String())
	t.T().Logf("OUT :\n%s", rootTree.PrintStyle(Box))

	result := rootTree.PrintStyle(ASCII)
	t.Equal(`1
'-- a
    |-- i
    '-- ii
`, result)
}

func (t *treeTestSuite) TestMultipleRoots() {
	rootTree := NewTree()
	oneTree := rootTree.Add("1")
	oneTree.Add("a")
	oneTree.Add("b")
	twoTree := rootTree.Add("2")
	twoTree.Add("a")
	twoTree.Add("b")

	t.T().Logf("TREE:\n%s", rootTree.String())
	t.T().Logf("OUT :\n%s", rootTree.PrintStyle(Box))

	result := rootTree.PrintStyle(ASCII)
	t.Equal(`1
|-- a
'-- b
2
|-- a
'-- b
`, result)
}

func (t *treeTestSuite) TestNephews() {
	tRoot := NewTree()
	t1 := tRoot.Add("1")
	tA := t1.Add("a")
	tA.Add("i")
	tA.Add("ii")
	tB := t1.Add("b")
	tB.Add("i")
	tB.Add("ii")

	t.T().Logf("TREE:\n%s", tRoot.String())
	t.T().Logf("OUT :\n%s", tRoot.PrintStyle(Box))

	result := tRoot.PrintStyle(ASCII)
	t.Equal(`1
|-- a
|   |-- i
|   '-- ii
'-- b
    |-- i
    '-- ii
`, result)
}

func (t *treeTestSuite) TestComplexTree() {
	// replicate a complex subdirectory:
	//   |-- vda
	//   |   |-- api
	//   |   |   |-- auth.go
	//   |   |   |-- engine.go
	//   |   |   `-- graphql.go
	//   |   |-- clair
	//   |   |   `-- api
	//   |   |       `-- v1
	//   |   |           |-- models.go
	//   |   |           `-- readme.md
	//   |   |-- engine-run-test
	//   |   |   `-- engine-run.go
	//   |   |-- errors.go
	//   |   |-- go-config
	//   |   |   |-- config.go
	//   |   |   |-- config_test.go
	//   |   |   `-- testdata
	//   |   |       |-- config_test.json
	//   |   |       |-- config_test.yml

	tRoot := NewTree()
	tVda := tRoot.Add("vda")
	vdaChildren := tVda.AddAll("api", "clair", "engine-run-test", "errors.go", "go-config")
	// api
	vdaChildren[0].AddAll("auth.go", "engine.go", "graphql.go")
	// clair
	vdaChildren[1].Add("api").Add("v1").AddAll("models.go", "readme.md")
	// engine-run-test
	vdaChildren[2].Add("engine-run.go")
	// go-config
	vdaChildren[4].AddAll("config.go", "config_test.go", "testdata")[2].AddAll("config_test.json", "config_test.yaml")

	t.T().Logf("OUT :\n%s", tRoot.PrintStyle(ASCII))

	// the result was confirmed visually (from the above log statement) then copied and pasted
	result := tRoot.PrintStyle(ASCII)
	t.Equal("vda\n|-- api\n|   |-- auth.go\n|   |-- engine.go\n|   '-- graphql.go\n|-- clair\n|   '-- api\n|       '-- v1\n|           |-- models.go\n|           '-- readme.md\n|-- engine-run-test\n|   '-- engine-run.go\n|-- errors.go\n'-- go-config\n    |-- config.go\n    |-- config_test.go\n    '-- testdata\n        |-- config_test.json\n        '-- config_test.yaml\n", result)
}

func (t *treeTestSuite) TestMultilineBranches() {
	rootTree := NewTree()
	oneTree := rootTree.Add("1")
	oneTree.Add("a\nalfa\nalpha\nable")
	oneTree.Add("b")

	t.T().Logf("TREE:\n%s", rootTree.String())
	t.T().Logf("OUT :\n%s", rootTree.PrintStyle(Box))

	result := rootTree.PrintStyle(ASCII)
	t.Equal(`1
|-- a
|   alfa
|   alpha
|   able
'-- b
`, result)
}

func (t *treeTestSuite) TestStyles() {
	tRoot := NewTree()
	t1 := tRoot.Add("1")
	tA := t1.Add("a")
	tA.Add("i")
	tA.Add("ii")
	tB := t1.Add("b")
	tB.Add("i")
	tB.Add("ii")

	cases := []struct {
		style    TreeStyle
		expected string
	}{
		{style: ASCII, expected: `1
|-- a
|   |-- i
|   '-- ii
'-- b
    |-- i
    '-- ii
`},
		{style: Box, expected: `1
├── a
│   ├── i
│   ╰── ii
╰── b
    ├── i
    ╰── ii
`},
		{style: Heavy, expected: `1
┣━━ a
┃   ┣━━ i
┃   ┗━━ ii
┗━━ b
    ┣━━ i
    ┗━━ ii
`},
		{style: ASCIINarrow, expected: `1
|-a
| |-i
| '-ii
'-b
  |-i
  '-ii
`},
		{style: BoxNarrow, expected: `1
├ a
│ ├ i
│ ╰ ii
╰ b
  ├ i
  ╰ ii
`},
		{style: HeavyNarrow, expected: `1
┣ a
┃ ┣ i
┃ ┗ ii
┗ b
  ┣ i
  ┗ ii
`},
		{style: WhiteSpace, expected: `1
    a
        i
        ii
    b
        i
        ii
`},
		{style: Dots, expected: `1
··· a
······· i
······· ii
··· b
······· i
······· ii
`},
		{style: Arrows, expected: `1
→   a
    →   i
    →   ii
→   b
    →   i
    →   ii
`},
	}

	for index, tc := range cases {
		t.T().Logf("Style %d:\n%s", tc.style, tRoot.PrintStyle(tc.style))
		t.Equal(tc.expected, tRoot.PrintStyle(tc.style), "Test case %d failed", index)
	}
}

func (*treeTestSuite) TestForDocumentationExample() {
	root := NewTree()
	colorTree := root.Add("Monitors")
	monoTree := colorTree.Add("Monocrome")
	monoTree.Add("Old School").AddAll("black", "green")
	monoTree.Add("Contemporary").AddAll("black", "white")
	colorTree.Add("Color").AddAll("red", "green", "blue")

	fmt.Println(root.Print())
}

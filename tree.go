package printtree

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Tree struct {
	Label    string // branch name. will be "" in the root node
	Branches []*Tree
}

// BranchLess accepts two branches and returns true if the first branch is less than (comes
// before) the second branch
type BranchLess func(branch1, branch2 *Tree) bool

// TreeStyle is the markup style of the tree. It may be one of the `...Style` constants or a
// higher number if custom styles have been added
type TreeStyle int

const (
	ASCIIStyle = iota
	BoxStyle
	BoxBoldStyle
	ASCIINarrowStyle
	BoxNarrowStyle
	BoxBoldNarrowStyle
	WhiteSpaceStyle
	ASCIIBulletStyle
	BulletStyle
	OrderedStyle
	NumberStyle
	AlphaStyle
	AlphaUCStyle
	RomanStyle
	RomanUCStyle
)

// for the different structural scaffolding types
const (
	midBranchScaffold = iota
	lastBranchScaffold
	bypassBranchScaffold
	noBranchScaffold
)

// for the different list scaffolding types
const (
	indentList = iota
	levelList
)

type scaffolding struct {
	isList bool     // true if this is a bullet style list
	markup []string // the markup for different types/levels of branches
}

var scaffoldingDict = []scaffolding{
	{false, []string{"|-- ", "'-- ", "|   ", "    "}},
	{false, []string{"├── ", "╰── ", "│   ", "    "}},
	{false, []string{"┣━━ ", "┗━━ ", "┃   ", "    "}},
	{false, []string{"|-", "'-", "| ", "  "}},
	{false, []string{"├ ", "╰ ", "│ ", "  "}},
	{false, []string{"┣ ", "┗ ", "┃ ", "  "}},
	{true, []string{"    ", "    "}},
	{true, []string{"  ", "* ", "+ ", "- "}},
	{true, []string{"  ", "● ", "○ ", "■ ", "□ "}},
	{true, []string{"    ", " 1. ", " a. ", " i. ", " A. ", " I. "}},
	{true, []string{"    ", " 1. "}},
	{true, []string{"    ", " a. "}},
	{true, []string{"    ", " A. "}},
	{true, []string{"      ", "   i. "}},
	{true, []string{"      ", "   I. "}},
}

// NewTree returns a new tree node that has no label. This is the root of a tree that you can
// work with later. When you print this tree, this root node will not be printed, only the
// branches of this root node. This allows for creating multiple visual roots like
//   root := NewTree()
//   root.AddBranch("1. First").AddBranch("a. Alpha")
//   root.AddBranch("2. second").AddBranch("b. Bravo")
// will print the following
//   1. First
//   '-- a. Alpha
//   2. Second
//   '-- b. Bravo
func NewTree() *Tree {
	return &Tree{}
}

// AddBranch creates a new branch Tree (at the end of the branch list) in this tree and returns that
// branch
func (tree *Tree) AddBranch(branchName string) *Tree {
	childTree := Tree{
		Label: branchName,
	}
	tree.Branches = append(tree.Branches, &childTree)
	return &childTree
}

// AddBranches creates multiple branches in one call in the order they are listed. Returns each
// new branch in the same order they were created
func (tree *Tree) AddBranches(branchNames ...string) []*Tree {
	branches := make([]*Tree, 0, len(branchNames))
	for _, childName := range branchNames {
		child := tree.AddBranch(childName)
		branches = append(branches, child)
	}
	return branches
}

// AddBranchf creates a new branch (at the end of the branch list) in this tree and returns that
// branch. This is a convenience function for adding a branch with a formatted name to the tree
func (tree *Tree) AddBranchf(label string, a ...interface{}) *Tree {
	return tree.AddBranch(fmt.Sprintf(label, a...))
}

// AddTreeAsBranch grafts in a tree as a branch of this tree. If the other tree has no label,
// then it is assumed to be a root node, and all it's branches will be added. If it does have a
// label, then it will be added as a branch
//
// WARNING: It is up to the caller not to create infinite tree loops
func (tree *Tree) AddTreeAsBranch(other *Tree) {
	if other.Label == "" {
		// this is a root tree, copy all it's children
		tree.Branches = append(tree.Branches, other.Branches...)
	} else {
		// tree branch, add the branch to this tree
		tree.Branches = append(tree.Branches, other)
	}
}

// Depth returns the maximum depth of the tree. A root tree with no branches is depth 0, a tree
// with one level of branches has depth 1
func (tree *Tree) Depth() int {
	depth := 0

	for index := range tree.Branches {
		branchDepth := 1 + tree.Branches[index].Depth()
		if branchDepth > depth {
			depth = branchDepth
		}
	}

	return depth
}

// Sort sorts the children of this tree by the labels
func (tree *Tree) Sort() {
	sort.SliceStable(tree.Branches, func(i, j int) bool {
		return tree.Branches[i].Label < tree.Branches[j].Label
	})
}

// DeepSort sorts the children of this tree and all sub-trees by the labels
func (tree *Tree) DeepSort() {
	tree.Sort()
	for index := range tree.Branches {
		tree.Branches[index].DeepSort()
	}
}

// SortCustom sorts the children of this tree by calling a custom function. The less function
// must return true if child1 comes before child2 in the list
func (tree *Tree) SortCustom(less BranchLess) {
	sort.SliceStable(tree.Branches, func(i, j int) bool {
		return less(tree.Branches[i], tree.Branches[j])
	})
}

// DeepSortCustom sorts the children of this tree and all sub-trees by calling a custom function. The
// less function must return true if child1 comes before child2 in the list
func (tree *Tree) DeepSortCustom(less BranchLess) {
	tree.SortCustom(less)
	for index := range tree.Branches {
		tree.Branches[index].DeepSortCustom(less)
	}
}

// AddStructuralStyle adds a new, custom style to the dictionary of structural styles. Pass in the
// strucutre that you want to use for different types of branches. Best results are obtained if
// all the branch structures are the same length.
//
// For example
//   middle branch   "|>- "
//   last branch     " `- "
//   bypass branch   "|   "
//   no branch       "O   "
// Produces a structure like
//   Parent
//   |>- Child1
//   |>- Child2
//   |   |>- Grandchild1
//   |    `- Grandchild2
//    `- Child3
//   O    `- Grandchild3
// The return value will be the value you can pass to `PrintStyle()` to use this style
func AddStructuralStyle(middleBranch, lastBranch, bypassBranch, noBranch string) TreeStyle {
	scaffoldingDict = append(scaffoldingDict, scaffolding{
		isList: false,
		markup: []string{middleBranch, lastBranch, bypassBranch, noBranch},
	})
	return TreeStyle(len(scaffoldingDict) - 1)
}

// AddListStyle adds a new, custom style to the dictionary of bulleted or ordered list styles.
// The indent is the string to use to indent each level of nesting, and the remaining values are
// the bullets or numbers to use for the branches of the tree
//
// Unordered List
//
// An unordered list uses the same "bullet" for each branch. You can use any character (or
// string of characters) as a bullet except for 1, a, A, i, or I (see below for the meanings of
// these characters). You may include as many bullets as you want and each bullet will be used
// for the corresponding indent levels of the tree. If there are more indent levels in the tree
// than bullet types, then the bullets will be recycled.
//
// For example:
//    myStyle = treeprint.AddListStyle("  ", "* ", "+ ", "- ")
// Would produce a tree like
//    Grandpappy
//    * Mom
//      + Me
//        - Son
//          * Grand-daughter
//            + Great-grand-son
//        - Daughter
//      + Sister
//
// Ordered List
//
// An ordered list is indicated by special "bullets" that contain 1, a, A, i, or I. When these
// characters are encountered, they will be replaced with a sequence that increments with every
// branch as follows:
//    1 -> 1, 2, 3, ... 9, 10, 11, ... 99, 100, 101 ...
//    a -> a, b, c, ... z, aa, ab, ... az, ba, bb, ... zz, aaa, aab, ...
//    A -> A, B, C, ... Z, AA, AB, ... AZ, BA, BB, ... ZZ, AAA, AAB, ...
//    i -> i, ii, iii, ... ix, x, xi, ... xlix, l, li, ... cmxcix, m, mi, ...
//    I -> I, II, III, ... IX, X, XI, ... XLIX, L, LI, ... CMXCIX, M, MI, ...
// When expanding the "bullet" to numbers, spaces to the left are used before expanding the
// width of the string. For example, if the "bullet" is "( 1)", "( a)" or "( i)" then the
// sequence will have 1 character to expand into on the left before expanding to the right. For
// example, the columns below show the expansion of numbers, alphabetic and roman numerals
//    ( 9)       ( x)       ( i)
//    (98)       (xx)       (ii)
//    (987)      (xxx)      (iii)
//    (9876)     (xxxx)     (xxix)
//
// You can mix ordered and unordered lists with no problem:
//    myStyle = treeprint.AddListStyle("  ", "* ", "1 ", "- ")
// Would produce a tree like
//    Grandpappy
//    * Mom
//      1 Me
//        - Son
//          * Grand-daughter
//            1 Great-grand-son
//        - Daughter
//      2 Sister
//
// The return value will be the value you can pass to `PrintStyle()` to use this style
func AddListStyle(indent string, bullets ...string) TreeStyle {
	scaffoldingDict = append(scaffoldingDict, scaffolding{
		isList: true,
		markup: []string{indent},
	})
	styleIndex := len(scaffoldingDict) - 1
	scaffoldingDict[styleIndex].markup = append(scaffoldingDict[styleIndex].markup, bullets...)
	return TreeStyle(styleIndex)
}

// String returns a string representation of this tree indented with whitespace
func (tree *Tree) String() string {
	return tree.PrintStyle(WhiteSpaceStyle)
}

// Print returns a string which is this tree in the default style (BoxStyle).
func (tree *Tree) Print() string {
	return tree.PrintStyle(BoxStyle)
}

// PrintStyle returns a string which is this tree printed with custom settings. The TreeStyle
// indicates what style of markup should be used on the left side of the tree.
func (tree *Tree) PrintStyle(style TreeStyle) string {
	// sanity checks
	if style < 0 || int(style) >= len(scaffoldingDict) {
		style = BoxStyle
	}

	scaffold := scaffoldingDict[style]
	buf := strings.Builder{}
	tree.print(&buf, 0, "", scaffold)
	return buf.String()
}

// print is the internal, recursive hook for printing the tree
func (tree *Tree) print(buf *strings.Builder, depth int, padding string, scaffold scaffolding) {
	var prefix string // prefix of each line

	for index := range tree.Branches {
		branch := tree.Branches[index]

		// handle each line of a block of text separately
		for lineIndex, line := range strings.Split(branch.Label, "\n") {
			if lineIndex == 0 {
				// first (or only) line of a block of text.
				prefix = padding + tree.labelPadding(depth, index, scaffold)
			} else {
				// subsequent lines of a block of text. the scaffold is one that indicates that
				// indicates we are flowing some text
				prefix = padding + tree.flowPadding(depth, index, scaffold)
			}
			buf.WriteString(prefix + line + "\n")
		}

		prefix = padding + tree.flowPadding(depth, index, scaffold)
		branch.print(buf, depth+1, prefix, scaffold)
	}
}

func (tree *Tree) labelPadding(depth int, index int, scaffold scaffolding) string {
	if depth == 0 {
		return ""
	}

	if scaffold.isList {
		// scaffold is a bulleted or numbered list
		offset := (depth - 1) % (len(scaffold.markup) - 1)
		return tree.replaceNumberListMarkup(scaffold.markup[levelList+offset], index+1)
	}

	// scaffold is structural
	switch {
	case index == len(tree.Branches)-1:
		return scaffold.markup[lastBranchScaffold]
	}
	return scaffold.markup[midBranchScaffold]
}

func (tree *Tree) flowPadding(depth int, index int, scaffold scaffolding) string {
	if depth == 0 {
		return ""
	}

	if scaffold.isList {
		// scaffold is a bulleted or numbered list
		return scaffold.markup[indentList]
	}

	// scaffold is structural
	switch {
	case index == len(tree.Branches)-1:
		return scaffold.markup[noBranchScaffold]
	}
	return scaffold.markup[bypassBranchScaffold]
}

// replaceNumberListMarkup replaces number markup (1, a, i) with a version of the number in the appropriate
// format. That is, replaces "1" with 1, 2, 3, etc; replaces "a" with a, b, c, etc; replaces "i"
// with i, ii, iii. Uses uppercase in the case of A and I. In order to keep alignment, attempts
// to consume spaces to the left of the markup character first
func (tree *Tree) replaceNumberListMarkup(markup string, index int) string {
	switch {
	case strings.Contains(markup, "1"):
		return tree.replaceNumberPlaceholder(markup, "1", strconv.Itoa(index))
	case strings.Contains(markup, "a"):
		return tree.replaceNumberPlaceholder(markup, "a", tree.convertToAlpha(index))
	case strings.Contains(markup, "A"):
		return tree.replaceNumberPlaceholder(markup, "A", strings.ToUpper(tree.convertToAlpha(index)))
	case strings.Contains(markup, "i"):
		return tree.replaceNumberPlaceholder(markup, "i", tree.convertToRoman(index))
	case strings.Contains(markup, "I"):
		return tree.replaceNumberPlaceholder(markup, "I", strings.ToUpper(tree.convertToRoman(index)))
	}

	// this markup did not contain a number field
	return markup
}

// replaceNumberPlaceholder replaces a number-placeholder in a markup string with the specific
// numeric value. This will replace the placeholder with the actualValue, while consuming
// whitespace to the left of the placeholder before expanding to the right.
func (tree *Tree) replaceNumberPlaceholder(s string, placeholder string, actualValue string) string {
	// find the placeholder and all whitespace to the left
	pattern := regexp.MustCompile(" *" + placeholder)
	loc := pattern.FindIndex([]byte(s))
	if loc == nil {
		return s
	}

	// replace with the actual value, padded to the same length as being replaced. if the actual
	// value is longer, that is fine and it will just flow to the right
	actualValue = fmt.Sprintf("%*s", loc[1]-loc[0], actualValue)
	return s[:loc[0]] + actualValue + s[loc[1]:]
}

// convertToAlpha converts a base-1 integer to a base-26 number using lower case alphabetic values
func (tree *Tree) convertToAlpha(n int) string {
	if n <= 0 {
		return "-"
	}
	s := ""
	for n > 0 {
		s = string(rune(((n-1)%26)+int('a'))) + s
		n = (n - 1) / 26
	}
	return s
}

var romanMap = []struct {
	value int
	key   string
}{
	{value: 1000, key: "m"},
	{value: 900, key: "cm"},
	{value: 500, key: "d"},
	{value: 400, key: "cd"},
	{value: 100, key: "c"},
	{value: 90, key: "xc"},
	{value: 50, key: "l"},
	{value: 40, key: "xl"},
	{value: 10, key: "x"},
	{value: 9, key: "ix"},
	{value: 5, key: "v"},
	{value: 4, key: "iv"},
	{value: 1, key: "i"},
}

// convertToAlpha converts a base-1 integer to roman numerals using lower case alphabetic values
func (tree *Tree) convertToRoman(n int) string {
	if n <= 0 {
		return "-"
	}

	s := ""
	for n > 0 {
		for _, rn := range romanMap {
			if n >= rn.value {
				s += rn.key
				n -= rn.value
				break
			}
		}
	}
	return s
}

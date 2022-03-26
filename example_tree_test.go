package printtree

import (
	"fmt"
)

func Example() {
	// initialize tree
	tree := NewTree()
	// add root element
	root := tree.AddBranch("Monitors")
	// add branch under the root
	monochromeBranch := root.AddBranch("Monocrome")
	// add nested branch with another level of deeper values
	monochromeBranch.AddBranch("Old School").AddBranches("black", "green")
	monochromeBranch.AddBranch("Contemporary").AddBranches("black", "white")
	root.AddBranch("Color").AddBranches("red", "green", "blue")

	fmt.Println(tree.Print())

	// Output:
	// Monitors
	// ├── Monocrome
	// │   ├── Old School
	// │   │   ├── black
	// │   │   ╰── green
	// │   ╰── Contemporary
	// │       ├── black
	// │       ╰── white
	// ╰── Color
	//     ├── red
	//     ├── green
	//     ╰── blue
}

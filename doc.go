// Package printtree is a simple tree that prints hierarchical lists with either bullets or
// markup that displays the list hierarchy. Create a new tree with "NewTree()", then add
// branches to it by calling AddBranch(name), AddBranchf(format, args...) or
// AddBranches(name...). Each of these will return the branch(es) that were added so you can
// continue to recursively add more nodes to the tree.
//
// When done, you can call Print() or PrintStyle(style) to get a string representing the tree
// that has traditional ascii-like heirarchy markup or bullets
//
// Example:
//    root := NewTree()
//    colorTree := root.Add("Monitors")
//    monoTree := colorTree.Add("Monocrome")
//    monoTree.Add("Old School").AddAll("black", "green")
//    monoTree.Add("Contemporary").AddAll("black", "white")
//    colorTree.Add("Color").AddAll("red", "green", "blue")
//    fmt.Println(root.Print())
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
// You can use PrintStyle() (instead of Print()) to access other tree-branch styles that show
// the structure of the tree
//    tree.PrintStyle(ASCIIStyle)         = |-- ASCII
//    tree.PrintStyle(BoxStyle)           = ├── Box
//    tree.PrintStyle(BoxBoldStyle)       = ┣━━ Box Bold
//    tree.PrintStyle(ASCIINarrowStyle)   = |-ASCII Narrow
//    tree.PrintStyle(BoxNarrowStyle)     = ├ Box Narrow
//    tree.PrintStyle(BoxBoldNarrowStyle) = ┗ Box Bold Narrow
//
// You can also print the tree as a list with bullets or numbering
//    tree.PrintStyle(WhiteSpaceStyle)     =     White Space
//    tree.PrintStyle(ASCIIBulletStyle)    = * ASCII Bullet (*, +, -)
//    tree.PrintStyle(BulletStyle)         = ● Bullet (●, ○, ■, □ )
//    tree.PrintStyle(OrderedStyle)        =  1. Ordered (1., a., i., A., I.)
//    tree.PrintStyle(NumberStyle)         =  1. Numbers
//    tree.PrintStyle(AlphaStyle)          =  a. Alphabetic
//    tree.PrintStyle(AlphaUCStyle)        =  A. ALPHABETIC
//    tree.PrintStyle(RomanStyle)          =    i. Roman Numerals
//    tree.PrintStyle(RomanUCStyle)        =    I. ROMAN NUMERALS
//
// Or you can define your own style
//    myStyle := printtree.AddListStyle("    ", "(1) ", "(A) ", "(I) ")
//    tree.PrintStyle(myStyle)
// Prints
//    Monitors
//    (1) Monocrome
//        (A) Old School
//            (I) black
//            (II) green
//        (B) Contemporary
//            (I) black
//            (II) white
//    (2) Color
//        (A) red
//        (B) green
//        (C) blue
package printtree

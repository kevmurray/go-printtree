# go-printtree

Printtree is a Go module that allows you to define a tree of data then print that data as a
tree, ordered, or unordered list in a variety of styles.

For example, given a tree constructed with
```
	tree := NewTree()
	root := tree.AddBranch("Monitors")
	monoTree := root.AddBranch("Monochrome")
	monoTree.AddBranch("Old School").AddBranches("black", "green")
	monoTree.AddBranch("Contemporary").AddBranches("black", "white")
	root.AddBranch("Color").AddBranches("red", "green", "blue")
```

This can be printed as a structured tree:

**ASCII tree**
```
ASCII tree                    Box tree
----------                    --------
Monitors                      Monitors
|-- Monochrome                ├── Monochrome
|   |-- Old School            │   ├── Old School
|   |   |-- black             │   │   ├── black
|   |   '-- green             │   │   ╰── green
|   '-- Contemporary          │   ╰── Contemporary
|       |-- black             │       ├── black
|       '-- white             │       ╰── white
'-- Color                     ╰── Color
    |-- red                       ├── red
    |-- green                     ├── green
    '-- blue                      ╰── blue
```

Or as lists
```
Unordered list                Ordered list
--------------                ------------
Monitors                      Monitors
● Monochrome                    1. Monochrome
  ○ Old School                      a. Old School
    ■ black                             i. black
    ■ green                            ii. green
  ○ Contemporary                    b. Contemporary
    ■ black                             i. black
    ■ white                            ii. white
● Color                         2. Color
  ○ red                             a. red
  ○ green                           b. green
  ○ blue                            c. blue
```

Or you can define your own structural tree, ordered, or unordered list style.

## Features
- Simple tree building
- Trees can be sorted
- Many pre-defined tree and list styles
- Customizable tree and list styles

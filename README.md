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

This can be printed in a variety of ways:

**ASCII tree**
```
Monitors
|-- Monochrome
|   |-- Old School
|   |   |-- black
|   |   '-- green
|   '-- Contemporary
|       |-- black
|       '-- white
'-- Color
    |-- red
    |-- green
    '-- blue
```

**Box characters**
```
Monitors
├── Monochrome
│   ├── Old School
│   │   ├── black
│   │   ╰── green
│   ╰── Contemporary
│       ├── black
│       ╰── white
╰── Color
    ├── red
    ├── green
    ╰── blue
```

**Unordered list**
```
Monitors
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
```

**Ordered list**
```
Monitors
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
```

Or you can define your own structural tree, ordered, or unordered list style.

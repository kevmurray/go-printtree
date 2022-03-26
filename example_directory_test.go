package printtree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Example_directoryTree demonstrates recursively walking a directory tree and creating a
// printable tree from it
func Example_directoryTree() {
	// create tree and root
	rootDir := "testdata"
	tree := NewTree()
	root := tree.AddBranch(rootDir)

	// build the tree
	if err := addFiles(root, rootDir); err != nil {
		panic("unable to walk test directory")
	}

	// print tree
	fmt.Println(tree.Print())
}

// addFiles recursively adds all the files and directories to the tree
func addFiles(tree *Tree, path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return fmt.Errorf("unable to read directory %s: %w", path, err)
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			// ignore hidden files
			continue
		}

		// add the file/dir as a child branch of the tree
		child := tree.AddBranch(file.Name())
		if file.IsDir() {
			// recur into directories
			if err := addFiles(child, filepath.Join(path, file.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

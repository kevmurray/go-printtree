
test:
	go test ./...

dir-tree:
	go test -run "^TestDirTree$$" -v github.com/kevmurray/printtree

size-tree:
	go test -run "^TestDirTreeWithSize$$" -v github.com/kevmurray/printtree

size-aligned-tree:
	go test -run "^TestDirTreeWithCustomStyle$$" -v github.com/kevmurray/printtree

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type fileEntry struct {
	Name string
	isFile bool
	size int64
}

func (entry fileEntry)getSize() string {
	if entry.size == 0 {
		return "empty"
	}
	return fmt.Sprintf("%db", entry.size)
}

type FileList []*fileEntry
func (list FileList) Len() int {
	return len(list)
}
func (list FileList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}
func (list FileList) Less(i, j int) bool {
	return list[i].Name < list[j].Name
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}

func dirTree(writer io.Writer, path string, printFiles bool) error {
	fileInfo, err := ioutil.ReadDir(path)
		if err != nil {
		return err
	}
	list,_ := processFilesList(fileInfo, printFiles)
	for i := 0; i < len(list); i++ {
		entry := list[i]
		isLast := i + 1 == len(list)
		if !entry.isFile {
			subdir := filepath.Join(path, entry.Name)
			fmt.Fprintf(writer, "%s%s\n", getPrefix(isLast), entry.Name)
			ident := "│\t"
			if (isLast) {
				ident = "\t"
			}
			writeDir(writer, subdir, printFiles, ident)
		}
		if entry.isFile {
			fmt.Fprintf(writer,"%s%s (%s)\n",getPrefix(isLast), entry.Name, entry.getSize())
		}
	}
	return nil
}

func processFilesList(filesList []os.FileInfo, printFiles bool) (FileList, error) {
	listDir := make(FileList, 0, len(filesList))
	listFile := make(FileList, 0, len(filesList))
	for _, entry := range filesList {
		if entry.IsDir() {
			listDir = append(listDir, &fileEntry{entry.Name(), false, 0})
		} else {
			listFile = append(listFile, &fileEntry{entry.Name(), true, entry.Size()})
		}
	}
	if printFiles {
		listDir = append(listDir, listFile...)
	}
	sort.Sort(listDir)
	return listDir, nil
}

func getPrefix(isLast bool) string {
	if isLast {
		return "└───"
	}
	return "├───"
}

func writeDir(writer io.Writer, path string, printFiles bool, ident string) error {
	fileInfo, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	list,_ := processFilesList(fileInfo, printFiles)
	for i := 0; i < len(list); i++ {
		entry := list[i]
		isLast := i + 1 == len(list)
		if !entry.isFile {
			subdir := filepath.Join(path, entry.Name)
			fmt.Fprintf(writer, "%s%s%s\n",ident, getPrefix(isLast), entry.Name)
			
			newIdent := ident
			if (isLast) {
				newIdent += "\t"
			} else {
				newIdent += ident
			}
			writeDir(writer, subdir, printFiles, newIdent)
		}
		if entry.isFile {
			fmt.Fprintf(writer,"%s%s%s (%s)\n",ident, getPrefix(isLast), entry.Name, entry.getSize())
		}
	}
	return nil
}

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

var (
	sizeInMB float64 = 999 // This is in megabytes
	suffixes [5]string
)

type Folder struct {
	Name string
	Size int64
}

type ByInt []Folder

func (a ByInt) Len() int           { return len(a) }
func (a ByInt) Less(i, j int) bool { return a[i].Size < a[j].Size }
func (a ByInt) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func main() {
	var (
		root    string
		files   []string
		folders []Folder
		err     error
	)

	fmt.Println("Enter Directory: ")
	fmt.Scanln(&root)
	fmt.Println("\nStarting Scan")

	files, err = IOReadDir(root)
	if err != nil {
		panic(err)
	}

	folders = make([]Folder, len(files))

	for _, name := range files {
		var (
			num int64
			_   error
		)

		num, _ = DirSize(root + "/" + name)
		folders = append(folders, Folder{name, num})
	}

	sort.Sort(ByInt(folders))

	for _, data := range folders {
		fmt.Println(data.Name + " - " + ByteCountSI(data.Size))
	}
}

func IOReadDir(root string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(root)
	if err != nil {
		return files, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}
	return files, nil
}

func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func ByteCountSI(b int64) string {
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}

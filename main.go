package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nigh/transliterate/pkg/transliterate"
)

var (
	help      bool
	inputPath string
	isDry     bool
)
var trans func(string) string

func init() {
	flag.BoolVar(&help, "help", false, "帮助")
	flag.StringVar(&inputPath, "in", "", "需要重命名的路径")
	flag.BoolVar(&isDry, "dry", false, "仅查看重命名结果，而不改变文件名")

	trans = transliterate.Sugar("-", "")
}

func main() {
	flag.Parse()
	if help || len(inputPath) == 0 {
		flag.Usage()
		return
	}
	if len(inputPath) == 0 {
		flag.Usage()
		return
	}

	filepath.Walk(inputPath, walker)
}

func walker(realPath string, f os.FileInfo, err error) error {
	ext := filepath.Ext(f.Name())
	oldName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
	newName := trans(oldName) + ext

	if f.IsDir() {
		fmt.Printf("[D] %s", oldName)
	} else {
		fmt.Printf("\t[F] %s", oldName)
	}
	fmt.Printf(" ---> %s", newName)
	fmt.Print("\n")
	return nil
}

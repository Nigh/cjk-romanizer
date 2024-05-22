package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Nigh/transliterate/pkg/transliterate"
	"github.com/logrusorgru/aurora/v4"
)

var (
	help        bool
	inputPath   string
	isDry       bool
	skipComfirm bool
	isSilent    bool
	separator   string
)
var trans func(string) string

type FilePath string
type FileRanames struct {
	path    FilePath
	oldName string
	newName string
	isDir   bool
}
type FilePaths []FileRanames
type FilePathSlice []string

func (a FilePath) Depth() (depth int) {
	slice := strings.Split(string(a), string(filepath.Separator))
	for _, v := range slice {
		if len(v) > 0 {
			depth++
		}
	}
	return
}

func (a FilePaths) Len() int      { return len(a) }
func (a FilePaths) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a FilePaths) Less(i, j int) bool {
	di := a[i].path.Depth()
	dj := a[j].path.Depth()
	return di > dj
}

var file2Rename FilePaths
var colorize *aurora.Aurora

func init() {
	flag.BoolVar(&help, "help", false, "help")
	flag.BoolVar(&isDry, "dry", false, "run without actually rename files")
	flag.BoolVar(&skipComfirm, "y", false, "skip comfirm")
	flag.BoolVar(&isSilent, "s", false, "silence output")

	flag.StringVar(&separator, "sp", "-", "separator between characters")
	flag.StringVar(&inputPath, "in", "", "input path")

	flag.Parse()
	if help || len(inputPath) == 0 {
		flag.Usage()
		return
	}
	if len(inputPath) == 0 {
		flag.Usage()
		return
	}
}

func askForAnswer() string {
	var ans string
	fmt.Scanln(&ans)
	return strings.ToLower(ans)
}
func askForContinue() bool {
	var yes string
	fmt.Scanln(&yes)
	return strings.ToLower(yes) == "y"
}

func main() {
	file2Rename = make(FilePaths, 0)
	trans = transliterate.Sugar(separator, "")
	colorize = aurora.New()
	inputPath = filepath.Clean(inputPath)
	_, err := os.Stat(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	filepath.Walk(inputPath, walker)
	if !isSilent {
		fmt.Print("There are total ", colorize.Green(len(file2Rename)), " files to rename\n\n")
		if !skipComfirm {
			fmt.Println(colorize.Cyan("Confirm to rename? [y/n]"))
			if !askForContinue() {
				return
			}
		}
	}

	if len(file2Rename) > 0 {
		ignoreError := false
		sort.Sort(file2Rename)
		fmtStr := fmt.Sprintf("[%%s] %%0%dd. %%s %%s%%s -> %%s\n", int(math.Ceil(math.Log10(float64(len(file2Rename)+1)))))
		for i, v := range file2Rename {
			var fType string
			if v.isDir {
				fType = "D"
			} else {
				fType = "F"
			}
			var err error
			if !isDry {
				old := filepath.Clean(string(v.path) + string(filepath.Separator) + v.oldName)
				new := filepath.Clean(string(v.path) + string(filepath.Separator) + v.newName)
				err = os.Rename(old, new)
			}
			if !isSilent {
				var result string
				if isDry {
					result = "↩️"
				} else {
					if err == nil {
						result = "✅"
					} else {
						result = "❌"
					}
				}
				fmt.Printf(fmtStr, result, i+1, colorize.BgWhite(fType).Black(), strings.TrimSuffix(string(v.path), string(filepath.Separator))+string(filepath.Separator), colorize.Yellow(v.oldName).Faint(), colorize.Green(v.newName).Bold())
			}

			if err != nil {
				fmt.Println(err)
				if !ignoreError {
					fmt.Println(colorize.Red("Continue with error? [Y/N/All]"))
					switch askForAnswer()[0] {
					case 'a':
						ignoreError = true
					case 'n':
						return
					}
				}
			}
		}
	}
}

func walker(realPath string, f os.FileInfo, err error) error {
	ext := filepath.Ext(f.Name())
	oldName := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
	newName := trans(oldName)

	if !isSilent {
		if f.IsDir() {
			fmt.Printf("[D] %s%c%s", strings.TrimSuffix(filepath.Dir(realPath), string(filepath.Separator)), filepath.Separator, colorize.BrightBlue(oldName+ext))
			if oldName != newName {
				fmt.Printf(" -> %s", colorize.BrightCyan(newName+ext))
			}
			fmt.Printf("\n")
		} else {
			if oldName != newName {
				fmt.Printf("\t[F] %s\n\t--> %s\n", colorize.BrightBlue(oldName+ext), colorize.BrightCyan(newName+ext))
			}
		}
	}

	if oldName != newName {
		file2Rename = append(file2Rename, FileRanames{path: FilePath(filepath.Dir(realPath)), oldName: oldName + ext, newName: newName + ext, isDir: f.IsDir()})
	}
	return nil
}

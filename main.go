package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Nigh/transliterate/pkg/transliterate"
	"github.com/integrii/flaggy"
	"github.com/logrusorgru/aurora/v4"
)

var (
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
var version string

func init() {
	separator = "-"
	flaggy.SetName("cjk-romanizer")
	flaggy.SetDescription("Rename CJK characters to roman characters")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/Nigh/cjk-romanizer"
	flaggy.AddPositionalValue(&inputPath, "path", 1, true, "the path to start rename")
	flaggy.Bool(&isDry, "d", "dry", "run without actually rename files")
	flaggy.Bool(&skipComfirm, "y", "comfirm", "skip comfirm")
	flaggy.Bool(&isSilent, "s", "silent", "silence output")
	flaggy.String(&separator, "sp", "separator", "separator between characters")
	flaggy.SetVersion(version)
	flaggy.Parse()
}

func askForAnswer() (ans string) {
	fmt.Scanln(&ans)
	ans = strings.ToLower(ans)
	if len(ans) == 0 {
		ans = " "
	}
	return
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

	inputPath, _ = filepath.Abs(inputPath)
	_, err := os.Stat(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	filepath.Walk(inputPath, walker)
	if !isSilent {
		fmt.Print("There are total ", colorize.Green(len(file2Rename)), " files to rename\n\n")
		if !skipComfirm {
			fmt.Println(colorize.Cyan("Confirm to rename? [Yes/No]"))
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
				var result aurora.Value
				if isDry {
					result = colorize.Yellow("DRY")
				} else {
					if err == nil {
						result = colorize.Green("✅")
					} else {
						result = colorize.Red("❌")
					}
				}
				fmt.Printf(fmtStr, result, i+1, colorize.BgWhite(fType).Black(), strings.TrimSuffix(string(v.path), string(filepath.Separator))+string(filepath.Separator), colorize.Yellow(v.oldName).Faint(), colorize.Green(v.newName).Bold())
			}

			if err != nil {
				fmt.Println(err)
				if !ignoreError {
					fmt.Println(colorize.Red("Continue with error? [Yes/No/All]"))
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
	if len(f.Name()) == 0 {
		return nil
	}
	if f.Name()[0] == '.' {
		return filepath.SkipDir
	}
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

<div align="center">
<img src="https://socialify.git.ci/Nigh/cjk-romanizer/image?description=1&font=KoHo&logo=https%3A%2F%2Fraw.githubusercontent.com%2FNigh%2Fcjk-romanizer%2Fmain%2Fassets%2Flogo.png&name=1&owner=1&pattern=Floating%20Cogs&stargazers=1&theme=Auto" alt="cjk-romanizer" width="640" height="320" />
</div>

# cjk-romanizer

Rename CJK characters in file and directory names to roman letters (中日韩文件名罗马化重命名).

Uses [`github.com/Nigh/transliterate`](https://github.com/Nigh/transliterate). Only names that change after transliteration are renamed; pure ASCII names are left alone. Extension is preserved. Hidden entries (names starting with `.`) are skipped.

![](assets/usage.gif)

## Install

- Download binaries from [Releases](https://github.com/Nigh/cjk-romanizer/releases) (`linux` / `darwin` / `windows`, amd64), or
- Build from source (see below).

## Usage

```shell
cjk-romanizer ~/download/test --dry -sp "_"
```

Safety defaults:

1. Walk the tree and collect renames
2. Print a count and ask for confirmation (`y`)
3. Rename deepest paths first (so parent dirs move after children)

| Flag | Meaning |
|------|---------|
| `-d` / `--dry` | Print planned renames only, do not rename |
| `-y` / `--comfirm` | Skip confirmation (flag name is `comfirm` in the CLI, matching the code) |
| `-s` / `--silent` | Quiet output |
| `-sp` / `--separator` | Separator between romanized segments (default: `-`) |

```text
Usage:
	cjk-romanizer [path]

Positional Variables:
	path   the path to start rename (Required)
Flags:
	   --version     Displays the program version string.
	-h --help        Displays help with available flag, subcommand, and positional value parameters.
	-d --dry         run without actually rename files
	-y --comfirm     skip comfirm
	-s --silent      silence output
	-sp --separator  separator between characters (default: -)
```

## Build

Requires Go 1.22+.

```shell
make          # build/cjk-romanizer-{linux,darwin,win}-amd64
# or one platform:
make linux
make darwin
make win64
```

Or:

```shell
go build -ldflags="-s -w -X main.version=1.0.0" -o cjk-romanizer .
```

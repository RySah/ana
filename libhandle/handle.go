package libhandle

import (
	"ana/compiler"
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func LibException(errorMsg string) {
	fmt.Printf("Lib Exception: %s\n", errorMsg)
	os.Exit(1)
}

const (
	DEFAULT_BYTE        = "0"
	DEFAULT_WORD        = "0"
	DEFAULT_DOUBLE_WORD = "0"
	DEFAULT_QUAD_WORD   = "0"
	DEFAULT_TEN_BYTE    = "0"
)

func tokenizeLine(line string) (tokens []string) {
	tokens = make([]string, 0)
	line = strings.Trim(line, " ")
	preTokens := strings.Split(line, " ")
	for _, t := range preTokens {
		if t != "" && t != " " {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

func inStrArray(tokens *[]string, token string) bool {
	for _, v := range *tokens {
		if v == token {
			return true
		}
	}
	return false
}

const SEP = string(filepath.Separator)

func GetLibPaths(exeDirPath string) (lib32DirPath string, lib64DirPath string) {
	lib32DirPath = exeDirPath + fmt.Sprintf("lib%s32bit", SEP)
	lib64DirPath = exeDirPath + fmt.Sprintf("lib%s64bit", SEP)
	return lib32DirPath, lib64DirPath
}
func GetLibFiles(libDirPath string) map[string][]fs.DirEntry {
	dirs, err := os.ReadDir(libDirPath)
	if err != nil {
		log.Fatal("Lib Exception:", err)
	}

	libMap := make(map[string][]fs.DirEntry)

	for _, dir := range dirs {
		if !dir.IsDir() {
			LibException(fmt.Sprintf("Lib Exception: Expected \"%s\" to be a directory", dir.Name()))
		}

		files, err := os.ReadDir(libDirPath + SEP + dir.Name())
		if err != nil {
			LibException(fmt.Sprintf("Unable to read from from directory \"%s\"", libDirPath+SEP+dir.Name()))
		}

		libMap[dir.Name()] = files
	}

	return libMap
}

func HandlePreProcessors(libPkgPath string, nasmInstance compiler.Nasm) compiler.Nasm {
	fmt.Printf("Installing package \"%s\"\n", libPkgPath)
	file, err := os.Open(libPkgPath)
	if err != nil {
		LibException("Unable to locate lib directory")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	definedAliases := make([]string, 0)
	reservedAliases := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) <= 2 {
			continue
		}
		if !(line[0] == ';' && line[1] == ':') {
			continue
		}
		line = line[2:]
		lineTokens := tokenizeLine(line)

		if len(lineTokens) > 0 {
			switch lineTokens[0] {
			case "define":
				if inStrArray(&definedAliases, lineTokens[2]) {
					LibException("This variable is defined elsewhere")
				} else {
					switch lineTokens[1] {
					case "b":
						fallthrough
					case "byte":
						nasmInstance.DataSection.DefineByte(lineTokens[2], DEFAULT_BYTE)

					case "w":
						fallthrough
					case "word":
						nasmInstance.DataSection.DefineWord(lineTokens[2], DEFAULT_WORD)

					case "dw":
						fallthrough
					case "dword":
						nasmInstance.DataSection.DefineDoubleWord(lineTokens[2], DEFAULT_DOUBLE_WORD)

					case "qw":
						fallthrough
					case "qword":
						nasmInstance.DataSection.DefineQuadWord(lineTokens[2], DEFAULT_QUAD_WORD)

					case "t":
						fallthrough
					case "tbytes":
						nasmInstance.DataSection.DefineTenBytes(lineTokens[2], DEFAULT_TEN_BYTE)

					case "const":
						nasmInstance.DataSection.DefineConstant(lineTokens[2], lineTokens[3])
					}
					definedAliases = append(definedAliases, lineTokens[2])
				}
			case "reserve":
				if inStrArray(&reservedAliases, lineTokens[2]) {
					LibException("This buffer is defined elsewhere")
				} else {
					switch lineTokens[1] {
					case "b":
						fallthrough
					case "byte":
						nasmInstance.BufferSection.ReserveByte(lineTokens[2], lineTokens[3])

					case "w":
						fallthrough
					case "word":
						nasmInstance.BufferSection.ReserveWord(lineTokens[2], lineTokens[3])

					case "dw":
						fallthrough
					case "dword":
						nasmInstance.BufferSection.ReserveDoubleWord(lineTokens[2], lineTokens[3])

					case "qw":
						fallthrough
					case "qword":
						nasmInstance.BufferSection.ReserveQuadWord(lineTokens[2], lineTokens[3])

					case "t":
						fallthrough
					case "tbytes":
						nasmInstance.BufferSection.ReserveTenBytes(lineTokens[2], lineTokens[3])
					}
					reservedAliases = append(reservedAliases, lineTokens[2])
				}
			}
		}
	}

	return nasmInstance
}

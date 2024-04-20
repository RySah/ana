package main

import (
	"ana/args"
	"ana/compiler"
	"ana/libhandle"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	WINDOWS_OS  = "windows"
	LINUX_OS    = "linux"
	ARCH_64_BIT = "amd64"
)

func main() {
	// comp := *compiler.NewCompiler("hello")
	arguments := args.HandleArgs()
	build := args.Exists(&arguments, args.BUILD)
	run := args.Exists(&arguments, args.RUN)
	if build || run {

		gen := ""
		op := "out.asm"

		{
			fp := ""
			if !run {
				v := args.GetValueOf(&arguments, args.BUILD)
				if v != nil {
					fp = *v
				} else {
					args.Exception(nil, "Could not locate file path argument in arguments")
				}
			} else {
				v := args.GetValueOf(&arguments, args.RUN)
				if v != nil {
					fp = *v
				} else {
					args.Exception(nil, "Could not locate file path argument in arguments")
				}
			}
			if args.Exists(&arguments, args.OUTPUT) {
				v := args.GetValueOf(&arguments, args.OUTPUT)
				if v != nil {
					op = *v
				} else {
					args.Exception(nil, "Could not locate file path argument in arguments")
				}
			}

			bytestream, err := os.ReadFile(fp)
			if err != nil {
				args.Exception(&fp, "Unable to read byte file")
			}

			arch := runtime.GOARCH
			if args.Exists(&arguments, args.ARCH) {
				v := args.GetValueOf(&arguments, args.ARCH)
				if v != nil {
					arch = *v
				} else {
					args.Exception(nil, "Could not parse the architecture in arguments")
				}
			}
			libPath := getLibPath(arch)
			libFiles := libhandle.GetLibFiles(libPath)
			lib := *libhandle.NewLibary(libPath, libFiles)

			cmp := compiler.NewCompiler("_start")
			srcLines := make([]string, 0)
			{
				built := ""
				for _, b := range bytestream {
					if b == byte('\n') {
						srcLines = append(srcLines, built)
						built = ""
					} else {
						built += string(rune(b))
					}
				}
			}
			cmp.HandleLines(srcLines)
			reqPkgs := cmp.Compile()

			const pkgDataSep = "/"

			for i, pkgData := range reqPkgs {
				tokens := strings.Split(pkgData, pkgDataSep)
				if !(len(tokens) > 1) {
					compiler.CompilerException(fmt.Sprintf("Unable to import package \"%s\"", reqPkgs[i]))
				}
				pkgSrc := tokens[0]
				pkgName := strings.Join(tokens[1:], pkgDataSep)

				pkg := lib.GetPackage(pkgSrc)
				srcPath, lines := pkg.Export(pkgName)
				if srcPath == "" {
					continue
				}

				cmp.NasmInstance = libhandle.HandlePreProcessors(srcPath, cmp.NasmInstance)

				conc := make([]string, len(cmp.ImportedLines)+len(lines))
				copy(conc, cmp.ImportedLines)
				copy(conc[len(cmp.ImportedLines):], lines)
				cmp.ImportedLines = conc
			}
			gen = cmp.NasmInstance.GeneratedFileContent() + strings.Join(cmp.ImportedLines, "\n")
		}

		{
			of, err := os.OpenFile(op, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				compiler.CompilerException("Unable to output compiled code")
			}
			defer of.Close()

			_, err = of.WriteString(gen)
			if err != nil {
				compiler.CompilerException("Unable to output compiled code")
			}
		}

		if run {
			on := ""
			{
				fileNameWithExt := filepath.Base(op)
				on = fileNameWithExt[:len(fileNameWithExt)-len(filepath.Ext(fileNameWithExt))]
			}

			toObjCmd := ""
			toObjArgs := make([]string, 0)

			linkCmd := ""
			linkArgs := make([]string, 0)

			switch runtime.GOOS {
			case WINDOWS_OS:
				if compiler.BitProcessIs64 {
					toObjCmd = "nasm"
					toObjArgs = []string{"-f", "win64", op, "-o", fmt.Sprintf("%s.obj", on)}
					linkCmd = "link"
					linkArgs = []string{"/subsystem:console", "/ENTRY:_start", fmt.Sprintf("%s.obj", on), "user32.lib", "kernel32.lib"}
				} else {
					toObjCmd = "nasm"
					toObjArgs = []string{"-f", "win32", op, "-o", fmt.Sprintf("%s.obj", on)}
					linkCmd = "link"
					linkArgs = []string{"/subsystem:console", "/ENTRY:_start", fmt.Sprintf("%s.obj", on), "user32.lib", "kernel32.lib"}
				}
			case LINUX_OS:
				if compiler.BitProcessIs64 {
					toObjCmd = "nasm"
					toObjArgs = []string{"-f", "elf64", op, "-o", fmt.Sprintf("%s.o", on)}
					linkCmd = "ld"
					linkArgs = []string{fmt.Sprintf("%s.o", on), "-o", on}
				} else {
					toObjCmd = "nasm"
					toObjArgs = []string{"-f", "elf32", op, "-o", fmt.Sprintf("%s.o", on)}
					linkCmd = "ld"
					linkArgs = []string{"-m", "elf_i386", fmt.Sprintf("%s.o", on), "-o", on}
				}
			}

			{
				_, err := exec.Command(toObjCmd, toObjArgs...).Output()
				if err != nil {
					compiler.CompilerException("Unable to convert \"" + op + "\" to object file")
				}
			}
			{
				_, err := exec.Command(linkCmd, linkArgs...).Output()
				if err != nil {
					compiler.CompilerException("Unable to link object file")
				}
			}

			{
				switch runtime.GOOS {
				case WINDOWS_OS:
					_, err := exec.Command(fmt.Sprintf("%s.exe", on)).Output()
					if err != nil {
						compiler.CompilerException("Unable to run generated executable")
					}
				case LINUX_OS:
					_, err := exec.Command(fmt.Sprintf("./%s", on)).Output()
					if err != nil {
						compiler.CompilerException("Unable to run generated executable")
					}
				}
			}
		}
	}
	os.Exit(0)
}

func getExeDir() (exeDir string) {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal("Error: ", err)
	}

	exeDir = filepath.Dir(exePath)
	if exeDir[len(exeDir)-1] != '\\' {
		exeDir += string(filepath.Separator)
	}
	return exeDir
}

func getLibPath(arch string) (evalLibPath string) {
	os := runtime.GOOS
	lib32Path, lib64Path := libhandle.GetLibPaths(getExeDir())
	compiler.BitProcessIs64 = arch == ARCH_64_BIT
	switch os {
	case WINDOWS_OS:
		if compiler.BitProcessIs64 {
			lib64Path += libhandle.SEP + "windows"
		} else {
			lib32Path += libhandle.SEP + "windows"
		}
	case LINUX_OS:
		if compiler.BitProcessIs64 {
			lib64Path += libhandle.SEP + "linux"
		} else {
			lib32Path += libhandle.SEP + "linux"
		}
	default:
		libhandle.LibException("Unsupported architecture")
	}

	if compiler.BitProcessIs64 {
		evalLibPath = lib64Path
	} else {
		evalLibPath = lib32Path
	}

	return evalLibPath
}

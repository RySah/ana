package libhandle

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Package struct {
	Name    string
	libPath string
	files   []fs.DirEntry
}

func NewPackage(name string, libPath string, files []fs.DirEntry) *Package {
	return &Package{
		Name:    name,
		libPath: libPath,
		files:   files,
	}
}

const (
	STD_PACKAGE_EXT = ".txt"
)

func (p *Package) Exists(n string) bool {
	for _, file := range p.files {
		if file.Name() == n+STD_PACKAGE_EXT {
			return true
		}
	}
	return false
}

func (p *Package) GetAllNames() (names []string) {
	names = make([]string, len(p.files))
	for i, file := range p.files {
		names[i] = file.Name()[:len(file.Name())-len(STD_PACKAGE_EXT)-1]
	}
	return names
}

func (p *Package) Export(n string) (filePath string, lines []string) {
	lines = make([]string, 0)
	for _, file := range p.files {
		if file.Name() == n+STD_PACKAGE_EXT {
			filePath = filepath.Join(p.libPath, p.Name, file.Name())
			content, err := os.Open(filePath)
			if err != nil {
				PkgException(fmt.Sprintf("Unable to read from package \"%s\"", n))
			}
			defer content.Close()

			scanner := bufio.NewScanner(content)

			for scanner.Scan() {
				text := scanner.Text()
				lines = append(lines, text)
			}
		}
	}

	return filePath, lines
}

package libhandle

import (
	"fmt"
	"io/fs"
	"os"
)

type Libary struct {
	path   string
	libMap map[string][]fs.DirEntry
}

func NewLibary(path string, libMap map[string][]fs.DirEntry) *Libary {
	fmt.Println(path)
	return &Libary{
		path:   path,
		libMap: libMap,
	}
}

func PkgException(errMsg string) {
	fmt.Printf("Import Exception: %s", errMsg)
	os.Exit(1)
}

func (l *Libary) Exists(pkgName string) bool {
	for name := range l.libMap {
		if pkgName == name {
			return true
		}
	}
	return false
}
func (l *Libary) GetPackage(pkgName string) *Package {
	for name := range l.libMap {
		if pkgName == name {
			return NewPackage(name, l.path, l.libMap[name])
		}
	}
	PkgException(fmt.Sprintf("Unable to find package \"%s\"", pkgName))
	return nil
}

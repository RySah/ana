package compiler

import "fmt"

type DataSection struct {
	lines []string
}

func NewDataSection() DataSection {
	return DataSection{
		lines: make([]string, 0),
	}
}

func (section *DataSection) GeneratedFileContent() string {
	total := "section .data\n"
	for _, line := range section.lines {
		total += fmt.Sprintf("\t%s\n", line)
	}
	return total
}

func (section *DataSection) GetLines() []string {
	return section.lines
}

func (section *DataSection) DefineConstant(alias string, data string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\tequ\t%s", alias, data))
}
func (section *DataSection) DefineByte(alias string, data string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\tdb\t%s", alias, data))
}
func (section *DataSection) DefineWord(alias string, data string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\tdw\t%s", alias, data))
}
func (section *DataSection) DefineDoubleWord(alias string, data string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\tdd\t%s", alias, data))
}
func (section *DataSection) DefineQuadWord(alias string, data string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\tdq\t%s", alias, data))
}
func (section *DataSection) DefineTenBytes(alias string, data string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\tdt\t%s", alias, data))
}

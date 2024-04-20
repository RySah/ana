package compiler

import "fmt"

type TextSection struct {
	lines []string
}

func NewTextSection() TextSection {
	return TextSection{
		lines: make([]string, 0),
	}
}

func (section *TextSection) GeneratedFileContent() string {
	total := ""
	for _, line := range section.lines {
		total += fmt.Sprintf("\t%s\n", line)
	}
	return total
}

func (section *TextSection) GetLines() []string {
	return section.lines
}

type TextSectionCmd func(a ...any) string

func (section *TextSection) Add(d string) {
	section.lines = append(section.lines, d)
}
func (section *TextSection) Addf(format string, a ...any) {
	section.Add(fmt.Sprintf(format, a...))
}

func (section *TextSection) CallAddr(addr string) {
	section.Addf("call %s", addr)
}
func (section *TextSection) MovData(dst string, d string) {
	section.Addf("mov %s, %s", dst, d)
}
func (section *TextSection) PushData(d string) {
	section.Addf("push %s", d)
}
func (section *TextSection) PopInto(dst string) {
	section.Addf("pop %s", dst)
}

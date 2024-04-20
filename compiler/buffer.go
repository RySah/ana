package compiler

import "fmt"

type BufferSection struct {
	lines []string
}

func NewBufferSection() BufferSection {
	return BufferSection{
		lines: make([]string, 0),
	}
}

func (section *BufferSection) GeneratedFileContent() string {
	total := "section .bss\n"
	for _, line := range section.lines {
		total += fmt.Sprintf("\t%s\n", line)
	}
	return total
}

func (section *BufferSection) GetLines() []string {
	return section.lines
}

func (section *BufferSection) ReserveByte(alias string, size string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\tresb\t%s", alias, size))
}
func (section *BufferSection) ReserveWord(alias string, size string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\tresw\t%s", alias, size))
}
func (section *BufferSection) ReserveDoubleWord(alias string, size string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\tresd\t%s", alias, size))
}
func (section *BufferSection) ReserveQuadWord(alias string, size string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\tresq\t%s", alias, size))
}
func (section *BufferSection) ReserveTenBytes(alias string, size string) {
	section.lines = append(section.lines, fmt.Sprintf("%s\trest\t%s", alias, size))
}

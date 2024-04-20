package compiler

import "fmt"

type Nasm struct {
	entry         string
	lines         []string
	DataSection   DataSection
	BufferSection BufferSection
	TextSection   TextSection
}

func NewNASM(entry string) *Nasm {
	return &Nasm{
		entry:         entry,
		lines:         make([]string, 0),
		DataSection:   NewDataSection(),
		BufferSection: NewBufferSection(),
		TextSection:   NewTextSection(),
	}
}

func (n *Nasm) GeneratedFileContent() string {
	total := fmt.Sprintf("global %s\n\n", n.entry)
	if len(n.DataSection.lines) > 0 {
		total += n.DataSection.GeneratedFileContent() + "\n"
	}
	if len(n.BufferSection.lines) > 0 {
		total += n.BufferSection.GeneratedFileContent() + "\n"
	}
	if len(n.TextSection.lines) > 0 {
		total += "section .text\n\n" + fmt.Sprintf("%s:\n", n.entry) + n.TextSection.GeneratedFileContent() + "\n"
	}
	return total
}

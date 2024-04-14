package compiler

type Nasm struct {
	name        string
	DataSection DataSection
}

func NewNASM(name string) *Nasm {
	return &Nasm{
		name:        name,
		DataSection: NewDataSection(),
	}
}

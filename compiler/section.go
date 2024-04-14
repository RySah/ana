package compiler

type Section interface {
	GenerateFileContent() string
}

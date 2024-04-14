package compiler

type Compiler struct {
	NasmInstance Nasm
	instrStack   Stack
	paramStack   Stack
}

func NewCompiler(name string) *Compiler {
	return &Compiler{
		NasmInstance: *NewNASM(name),
		instrStack:   *NewStack(),
		paramStack:   *NewStack(),
	}
}

func (c *Compiler) HandleBytes(byteStream []byte) {
	handler := NewBytecodeHandler(byteStream)
	c.instrStack, c.paramStack = handler.GetStacks()
}
func (c *Compiler) Compile() {

}

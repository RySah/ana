package compiler

type Stack struct {
	memory []byte
}

func NewStack() *Stack {
	return &Stack{
		memory: make([]byte, 0),
	}
}

func (s *Stack) Push(v byte) {
	s.memory = append(s.memory, v)
}
func (s *Stack) Pop() (popped byte) {
	popped = s.Peek()
	s.memory = s.memory[:len(s.memory)-1]
	return popped
}
func (s *Stack) Peek() (v byte) {
	v = s.memory[len(s.memory)-1]
	return v
}
func (s *Stack) Size() int {
	return len(s.memory)
}

package common

type Node struct {
	value interface{}
	prev  *Node
}

type Stack struct {
	top    *Node
	length int
}

func NewStack() *Stack {
	return &Stack{top: nil, length: 0}
}

func (s *Stack) Len() int {
	return s.length
}

func (s *Stack) Peek() interface{} {
	if s.length == 0 {
		return nil
	}

	return s.top.value
}

func (s *Stack) Push(x interface{}) {
	n := &Node{value: x, prev: s.top}
	s.top = n
	s.length++
}

func (s *Stack) Pop() interface{} {
	if s.length == 0 {
		return nil
	}

	n := s.top
	s.top = n.prev
	s.length--
	return n.value
}

func (s *Stack) Top() *Node {
	return s.top
}

func (n *Node) Next() *Node {
	return n.prev
}

func (n *Node) Value() interface{} {
	return n.value
}

package main

type Stack []Coord

func NewStack() Stack {
	return make([]Coord, 0)
}
func (stack *Stack) IsEmpty() bool {
	return len(*stack) == 0
}
func (stack *Stack) Push(c Coord) {
	*stack = append(*stack, c)
}
func (stack *Stack) Pop() (Coord, bool) {
	if stack.IsEmpty() {
		return Coord{}, false
	}
	index := len(*stack) - 1
	element := (*stack)[index]
	*stack = (*stack)[:index]

	return element, true
}

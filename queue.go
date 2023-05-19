package main

type Queue struct {
	Elements []Coord
}

func (q *Queue) IsEmpty() bool {
	return len(q.Elements) == 0
}
func (q *Queue) Length() int {
	return len(q.Elements)
}
func (q *Queue) Enqueue(c Coord) {
	q.Elements = append(q.Elements, c)
}
func (q *Queue) Dequeue() (Coord, bool) {
	if q.IsEmpty() {
		return Coord{}, false
	}
	c := q.Elements[0]
	if q.Length() == 1 {
		q.Elements = nil
		return c, true
	}

	q.Elements = q.Elements[1:]

	return c, true
}

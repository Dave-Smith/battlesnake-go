package main

type HeadZone struct {
	SnakeHead   Coord
	SnakeLength int
	SnakeName   string
	Zone        []Coord
}

func MakeHeadZones(snakes []Battlesnake, board Board, depthLimit int) []HeadZone {
	zones := make([]HeadZone, 0)
	for _, s := range snakes {
		zones = append(zones, MakeHeadZone(s, board, depthLimit))
	}

	return zones
}

func MakeHeadZone(snake Battlesnake, board Board, depthLimit int) HeadZone {
	head := snake.Head
	length := snake.Length
	name := snake.Name
	body := snake.Body[1:]

	seen := make([]Coord, 0)
	q := Queue{}
	q.Enqueue(head)

	depthTrack := 1
	countdownToNextDepth := 1

	for !q.IsEmpty() {
		curr, _ := q.Dequeue()

		// depth bookkeeping
		countdownToNextDepth--
		if countdownToNextDepth == 0 {
			depthTrack++
		}

		if depthTrack > depthLimit {
			continue
		}

		if hasCoord(curr, seen) {
			continue
		}
		if isOffBoard(curr, board) {
			continue
		}

		// don't traverse over snake body
		if hasCoord(curr, body) {
			continue
		}

		seen = append(seen, curr)

		for _, d := range directionalMoves {
			q.Enqueue(Coord{head.X + d.X, head.Y + d.Y})
			countdownToNextDepth++
		}
	}

	return HeadZone{SnakeHead: head, SnakeLength: length, SnakeName: name, Zone: seen}
}

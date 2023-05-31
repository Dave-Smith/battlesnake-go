package main

import (
	"math"
)

type Coords []Coord
type GameMap [][]CellOccupant
type CellOccupant int

const (
	None CellOccupant = iota
	Snake
	Hazard
	Food
	VulnerableSnake
)

type Movement int

const (
	Up Movement = iota
	Right
	Down
	Left
)

var directionalMoves = [4]Coord{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

func direction(curr Coord, next Coord) string {
	if curr.X == next.X {
		if curr.Y > next.Y {
			return "down"
		}
		return "up"
	}
	if curr.Y == next.Y {
		if curr.X > next.X {
			return "right"
		}
		return "left"
	}
	return "left"
}

func pathToCoord(curr, dest Coord) []Coord {
	path := make([]Coord, 0)
	moveY := 1
	moveX := 1
	if dest.Y-curr.Y < 0 {
		moveY = -1
	}
	if dest.X-curr.X < 0 {
		moveX = -1
	}

	var next Coord = curr
	for i := 0; i < int(math.Abs(float64(dest.Y-curr.Y))); i++ {
		next = Coord{next.X, next.Y + moveY}
		path = append(path, next)
	}
	for i := 0; i < int(math.Abs(float64(dest.X-curr.X))); i++ {
		next = Coord{next.X + moveX, next.Y}
		path = append(path, next)
	}
	return path
}

func fillMap(board Board, me Battlesnake) GameMap {
	var gameMap GameMap = make([][]CellOccupant, board.Width)
	for i := range gameMap {
		gameMap[i] = make([]CellOccupant, board.Width)
	}

	for x := 0; x < board.Width; x++ {
		for y := 0; y < board.Height; y++ {
			curr := Coord{x, y}
			if hasCoord(curr, board.Food) {
				gameMap[x][y] = Food
			}
			if hasCoord(curr, board.Hazards) {
				gameMap[x][y] = Hazard
			}
			for s := 0; s < len(board.Snakes); s++ {
				if hasCoord(curr, board.Snakes[s].Body) {
					gameMap[x][y] = Snake
				}
				if curr == board.Snakes[s].Head && board.Snakes[s].Length < me.Length {
					gameMap[x][y] = VulnerableSnake
				}
			}
		}
	}

	return gameMap
}

func hasCoord(curr Coord, coords []Coord) bool {
	for i := 0; i < len(coords); i++ {
		if coords[i] == curr {
			return true
		}
	}
	return false
}

func nearestCorner(curr Coord, board Board) Coord {
	x := 0
	y := 0
	middle := board.Height / 2

	if x > middle {
		x = board.Width - 1
	}
	if y > middle {
		y = board.Height - 1
	}
	return Coord{x, y}
}

func safeMoves(curr Coord, gameMap GameMap) []Coord {
	moves := make([]Coord, 0)
	for i := 0; i < len(directionalMoves); i++ {
		m := directionalMoves[i]
		if exists, occupant := cell(curr.X+m.X, curr.Y+m.Y, gameMap); exists && isSafe(*occupant) {
			moves = append(moves, Coord{curr.X + m.X, curr.Y + m.Y})
		}
	}
	return moves
}

func saferMoves(curr Coord, gameMap GameMap, seen []Coord, pathLength int, opponents []Battlesnake) []Coord {
	moves := make([]Coord, 0)
	if hasCoord(curr, seen) {
		return moves
	}
	if pathLength == 0 && isSafe(gameMap[curr.X][curr.Y]) {
		return []Coord{curr}
	}
	for i := 0; i < len(directionalMoves); i++ {
		m := directionalMoves[i]
		seen = append(seen, curr)
		if exists, occupant := cell(curr.X+m.X, curr.Y+m.Y, gameMap); exists && isSafe(*occupant) {
			next := Coord{curr.X + m.X, curr.Y + m.Y}
			if len(seen) <= 0 && adjecentCellHasSnakeHead(next, opponents) {
				continue
			}
			if len(saferMoves(next, gameMap, seen, pathLength-1, opponents)) > 0 {
				moves = append(moves, Coord{curr.X + m.X, curr.Y + m.Y})
			}
		}
	}
	return moves
}

func adjecentCellHasSnakeHead(curr Coord, snakes []Battlesnake) bool {
	for i := 0; i < len(snakes); i++ {
		snake := snakes[i]
		for d := 0; d < len(directionalMoves); d++ {
			adj := directionalMoves[d]
			next := Coord{curr.X + adj.X, curr.Y + adj.Y}
			if next == snake.Head {
				return true
			}
		}
	}
	return false
}

func isSafe(cell CellOccupant) bool {
	if cell == Hazard || cell == Snake {
		return false
	}
	return true
}

func cell(x, y int, gameMap GameMap) (bool, *CellOccupant) {
	if x >= len(gameMap) {
		return false, nil
	}
	if x < 0 {
		return false, nil
	}
	if y >= len(gameMap[0]) {
		return false, nil
	}
	if y < 0 {
		return false, nil
	}
	return true, &gameMap[x][y]
}

func dir(curr, next Coord) string {
	if curr.X == next.X {
		if curr.Y > next.Y {
			return "down"
		}
		if curr.Y < next.Y {
			return "up"
		}
	}
	if curr.Y == next.Y {
		if curr.X < next.X {
			return "right"
		}
		if curr.X > next.X {
			return "left"
		}
	}
	return ""
}

func makeNextMoves(curr Coord) []Coord {
	moves := make([]Coord, 0)
	for i := 0; i < len(directionalMoves); i++ {
		m := directionalMoves[i]
		moves = append(moves, Coord{curr.X + m.X, curr.Y + m.Y})
	}
	return moves
}
func nearest(curr Coord, food []Coord) Coord {
	if len(food) == 0 {
		return Coord{}
	}
	if len(food) == 1 {
		return food[0]
	}
	nearest := food[0]

	for i := 1; i < len(food); i++ {
		if distanceTo(curr, food[i]) < distanceTo(curr, nearest) {
			nearest = food[i]
		}
	}
	return nearest
}
func distanceTo(from, to Coord) int {
	return abs(from.X-to.X) + abs(from.Y-to.Y)
}

func abs(num int) int {
	if num < 0 {
		return num * -1
	}
	return num
}

func (coords Coords) has(c Coord) bool {
	for i := 0; i > len(coords); i++ {
		if coords[i] == c {
			return true
		}
	}
	return false
}

func floodFill(curr, target Coord, board Board, snakes []Battlesnake) Coords {
	seen := make([]Coord, 0)
	q := Queue{}

	q.Enqueue(curr)
	for !q.IsEmpty() {
		c, _ := q.Dequeue()
		seen = append(seen, c)
		for i := 0; i < len(makeNextMoves(c)); i++ {
			// if is safe and not seen
			// add to queue
		}
	}

	return seen
}

func hasSnakeCollision(curr Coord, snakes []Battlesnake) bool {
	for _, snake := range snakes {
		if hasCoord(curr, snake.Body) {
			return true
		}
	}
	return false
}

func isOffBoard(curr Coord, board Board) bool {
	if curr.X < 0 || curr.Y < 0 {
		return true
	}
	if curr.X >= board.Width || curr.Y >= board.Height {
		return true
	}
	return false
}

func (m Movement) asString() string {
	if m == Up {
		return "up"
	}
	if m == Right {
		return "right"
	}
	if m == Down {
		return "down"
	}
	return "left"
}

func UseMovement(source, target Coord) Movement {
	if source.X < target.X {
		return Right
	}
	if source.X > target.X {
		return Left
	}
	if source.Y < target.Y {
		return Up
	}
	return Down
}

package main

import (
	"log"
	"math"
)

type TailChasing struct {
	path         []Coord
	pathToCorner []Coord
	pathIndex    int
	corner       Coord
	fromWallY    int
	fromWallX    int
}

type GameMap [][]CellOccupant
type CellOccupant int

const (
	None CellOccupant = iota
	Snake
	Hazard
	Food
	VulnerableSnake
)

var directionalMoves = [4]Coord{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}

func NewTailChasingPlan() *TailChasing {
	return &TailChasing{}
}

func (plan *TailChasing) start(state GameState) {
	head := state.You.Head
	middle := state.Board.Height / 2
	plan.fromWallX = 1
	plan.fromWallY = 1

	plan.path = make([]Coord, 0)
	if head.X > state.Board.Width/2 {
		plan.fromWallX = -1
	}
	if head.Y > state.Board.Height/2 {
		plan.fromWallY = -1
	}

	plan.corner = nearestCorner(state.You.Head, state.Board)
	for i := 0; i < middle; i++ {
		plan.path = append(plan.path, Coord{plan.corner.X + i*plan.fromWallX, plan.corner.Y})
	}
	for i := middle; i >= 0; i-- {
		plan.path = append(plan.path, Coord{plan.corner.X + i*plan.fromWallX, plan.corner.Y + plan.fromWallY})
	}
	log.Printf("Board: %d wide, %d high", state.Board.Width, state.Board.Height)
	log.Printf("Snake body: %v", state.You.Body)
	log.Printf("Path created: %v", plan.path)
	plan.pathIndex = pathIndex(state.You.Head, plan.path)
	plan.pathToCorner = pathToCoord(state.You.Head, plan.corner)
	log.Printf("Path to corner: %v", plan.pathToCorner)
}

func (plan *TailChasing) move(state GameState) BattlesnakeMoveResponse {
	curr := state.You.Head
	plan.pathIndex = pathIndex(state.You.Head, plan.path)
	if plan.pathIndex != -1 {
		next := plan.path[plan.pathIndex+1%len(plan.path)]
		return BattlesnakeMoveResponse{Move: direction(curr, next)}
	}
	log.Printf("PathToCorner %v", plan.pathToCorner)
	next := plan.pathToCorner[0]
	plan.pathToCorner = plan.pathToCorner[1:]
	return BattlesnakeMoveResponse{Move: direction(curr, next)}
}

func pathIndex(head Coord, path []Coord) int {
	for i := 0; i < len(path); i++ {
		if path[i].X == head.X && path[i].Y == head.Y {
			return i
		}
	}
	return -1
}

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

func (plan *TailChasing) end(state GameState) {

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
	cell := func(x int, y int, gameMap GameMap) (bool, *CellOccupant) {
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
	isSafe := func(cell CellOccupant) bool {
		if cell == Hazard || cell == Snake {
			return false
		}
		return true
	}
	for i := 0; i < len(directionalMoves); i++ {
		m := directionalMoves[i]
		if exists, occupant := cell(curr.X+m.X, curr.Y+m.Y, gameMap); exists && isSafe(*occupant) {
			moves = append(moves, Coord{curr.X + m.X, curr.Y + m.Y})
		}
	}
	return moves
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

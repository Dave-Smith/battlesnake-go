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

func NewTailChasingPlan() *TailChasing {
	return &TailChasing{}
}

func (plan *TailChasing) start(state GameState) {
	head := state.You.Head
	x := 0
	y := 0
	plan.fromWallX = 1
	plan.fromWallY = 1

	plan.path = make([]Coord, 0)
	if head.X > state.Board.Width/2 {
		x = state.Board.Width
		plan.fromWallX = -1
	}
	if head.Y > state.Board.Height/2 {
		y = state.Board.Height
		plan.fromWallY = -1
	}

	plan.corner = Coord{x, y}
	for i := 0; i < 5; i++ {
		plan.path = append(plan.path, Coord{plan.corner.X + i*plan.fromWallX, plan.corner.Y})
	}
	for i := 5; i >= 0; i-- {
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

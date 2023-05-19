package main

import "log"

type WeightedMovement struct {
	movement     Movement
	root         Coord
	obstacles    int
	open         []Coord
	heads        int
	food         int
	deadEnd      bool
	certainDeath bool
}
type WeightedMovementSet []WeightedMovement

func makeOpeningMoves(c Coord) []WeightedMovement {
	return []WeightedMovement{
		{movement: Up, root: Coord{c.X, c.Y + 1}, open: make([]Coord, 0)},
		{movement: Right, root: Coord{c.X + 1, c.Y}, open: make([]Coord, 0)},
		{movement: Down, root: Coord{c.X, c.Y - 1}, open: make([]Coord, 0)},
		{movement: Left, root: Coord{c.X - 1, c.Y}, open: make([]Coord, 0)}}
}
func (w *WeightedMovement) addOpenSpot(c Coord) {
	w.open = append(w.open, c)
}
func fillToDepth(start Coord, depthLimit int, board Board) WeightedMovementSet {
	movements := makeOpeningMoves(start)

	for i := 0; i < len(movements); i++ {
		depth := 0
		countdownToNextDepth := 1
		move := movements[i]
		seen := make([]Coord, 0)
		q := Queue{}
		q.Enqueue(move.root)

		if isOffBoard(move.root, board) || hasSnakeCollision(move.root, board.Snakes) {
			log.Printf("Not moving %s to %v because of certain death", move.movement.asString(), move.root)
			move.certainDeath = true
		}

		for !q.IsEmpty() {
			curr, _ := q.Dequeue()

			// depth bookkeeping
			countdownToNextDepth--
			if countdownToNextDepth == 0 {
				depth++
			}

			if depth > depthLimit {
				continue
			}

			// if move is safe
			if hasCoord(curr, seen) {
				continue
			}
			if isOffBoard(curr, board) {
				continue
			}
			if hasSnakeCollision(curr, board.Snakes) {
				move.obstacles++
				continue
			}
			for _, food := range board.Food {
				if curr == food {
					move.food++
				}
			}
			for _, snake := range board.Snakes {
				if curr == snake.Head {
					move.heads++
				}
			}

			move.addOpenSpot(curr)
			seen = append(seen, curr)
			nextMoves := makeNextMoves(curr)
			for _, next := range nextMoves {
				countdownToNextDepth++
				q.Enqueue(next)
			}
		}
	}
	return movements
}

func (moves WeightedMovementSet) bestMoveForFood(you Battlesnake) WeightedMovement {
	best := moves[0]
	for i := 1; i < len(moves); i++ {
		move := moves[i]
		if best.certainDeath || move.food > best.food {
			best = move
		}
	}
	return best
}

func (moves WeightedMovementSet) bestMoveToAvoidFood(you Battlesnake) WeightedMovement {
	best := moves[0]
	for i := 1; i < len(moves); i++ {
		move := moves[i]
		if best.certainDeath || move.food > best.food {
			best = move
		}
	}
	return best
}

func (moves WeightedMovementSet) bestMoveForRoaming(you Battlesnake) WeightedMovement {
	best := moves[0]
	for i := 1; i < len(moves); i++ {
		move := moves[i]
		if best.certainDeath || len(move.open) > len(best.open) {
			best = move
		}
	}
	return best
}

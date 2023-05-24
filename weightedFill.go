package main

import (
	"log"
	"math/rand"
	"time"
)

type Opponent struct {
	distance  int
	length    int
	headCoord Coord
}
type WeightedMovement struct {
	movement        Movement
	root            Coord
	obstacles       int
	open            []Coord
	heads           int
	food            int
	deadEnd         bool
	certainDeath    bool
	nearestOpponent Opponent
}
type WeightedMovementSet []WeightedMovement

func makeOpeningMoves(c Coord) WeightedMovementSet {
	return []WeightedMovement{
		{movement: Up, root: Coord{c.X, c.Y + 1}, open: make([]Coord, 0, 5), nearestOpponent: Opponent{}},
		{movement: Right, root: Coord{c.X + 1, c.Y}, open: make([]Coord, 0, 5), nearestOpponent: Opponent{}},
		{movement: Down, root: Coord{c.X, c.Y - 1}, open: make([]Coord, 0, 5), nearestOpponent: Opponent{}},
		{movement: Left, root: Coord{c.X - 1, c.Y}, open: make([]Coord, 0, 5), nearestOpponent: Opponent{}}}
}

func (w *WeightedMovement) addOpenSpot(c Coord) {
	if w.open == nil {
		w.open = make([]Coord, 5)
	}
	w.open = append(w.open, c)
}

func fillToDepth(start Coord, depthLimit int, board Board) WeightedMovementSet {
	movements := makeOpeningMoves(start)

	for i := 0; i < len(movements); i++ {
		depth := 0
		countdownToNextDepth := 1
		seen := make([]Coord, 0)
		q := Queue{}
		q.Enqueue(movements[i].root)

		if isOffBoard(movements[i].root, board) || hasSnakeCollision(movements[i].root, board.Snakes) {
			movements[i].certainDeath = true
			log.Printf("Not moving %s to %v because of certain death, move deets %v", movements[i].movement.asString(), movements[i].root, movements[i])
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
				movements[i].obstacles++
				continue
			}
			for _, food := range board.Food {
				if curr == food {
					movements[i].food++
				}
			}
			for _, snake := range board.Snakes {
				if curr == snake.Head {
					movements[i].heads++
					if movements[i].nearestOpponent.distance >= depth {
						movements[i].nearestOpponent = Opponent{
							distance:  depth,
							length:    snake.Length,
							headCoord: snake.Head,
						}
					}
				}
				// log.Printf("Nearest snake %v", movements[i].nearestOpponent)
			}

			//log.Printf("Adding open spot %v to %v", curr, movements[i].root)
			movements[i].addOpenSpot(curr)
			seen = append(seen, curr)
			nextMoves := makeNextMoves(curr)
			for _, next := range nextMoves {
				countdownToNextDepth++
				q.Enqueue(next)
			}
		}
	}
	//log.Printf("After flood fill %v", movements)
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
	//log.Printf("Avoiding Food: Possible movements %v", moves)
	best := moves[0]
	for i := 1; i < len(moves); i++ {
		move := moves[i]
		if best.certainDeath || move.food > best.food {
			best = move
		}
	}
	return best
}

func (moves WeightedMovementSet) bestMoveForDefense(you Battlesnake) WeightedMovement {
	best := moves[0]
	for i := 1; i < len(moves); i++ {
		move := moves[i]
		if best.certainDeath || (move.nearestOpponent.distance > best.nearestOpponent.distance && move.nearestOpponent.length >= you.Length) {
			best = move
		}
	}
	return best
}

func (moves WeightedMovementSet) bestMoveForOffense(you Battlesnake) WeightedMovement {
	var best WeightedMovement
	for i := 0; i < len(moves); i++ {
		move := moves[i]
		if !move.certainDeath && (move.nearestOpponent.distance < best.nearestOpponent.distance && move.nearestOpponent.length < you.Length) {
			best = move
		}
	}
	return best
}

func (moves WeightedMovementSet) avoidCertainDeath() WeightedMovementSet {
	n := 0
	for _, val := range moves {
		if !val.certainDeath {
			moves[n] = val
			n++
		}
	}
	return moves[:n]
}

func (moves WeightedMovementSet) bestMoveForRoaming(you Battlesnake) WeightedMovement {
	//log.Printf("Roaming: Possible movements %v", moves)
	safest := make([]WeightedMovement, 0)
	safer := make([]WeightedMovement, 0)
	for i := 0; i >= len(moves); i++ {
		move := moves[i]
		if !move.certainDeath {
			safest = append(safest, move)
		}
		if !move.certainDeath && (move.nearestOpponent.distance == 0 || move.nearestOpponent.distance > 1) {
			safer = append(safer, move)
		}
	}
	if len(safest) == 0 {
		return moves[0]
	}

	rand.Seed(time.Now().Unix())
	if len(safer) == 0 {
		return safest[rand.Intn(len(safest))]
	}
	return safest[rand.Intn(len(safer))]
}

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
	opponentInDmz   bool
	movingToCorner  bool
	distanceToFood  int
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
	otherSnakes := make([]Battlesnake, 0)
	for i := 0; i < len(board.Snakes); i++ {
		if start != board.Snakes[i].Head {
			otherSnakes = append(otherSnakes, board.Snakes[i])
		}
	}

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

		// todo look out for corners.  might get trapped
		corners := []Coord{{0, 0}, {0, board.Height - 1}, {board.Width - 1, 0}, {board.Width - 1, board.Width - 1}}
		if hasCoord(movements[i].root, corners) {
			movements[i].movingToCorner = true
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
			if isOffBoard(curr, board) {
				continue
			}

			if hasCoord(curr, seen) {
				continue
			}
			seen = append(seen, curr)

			for _, snake := range otherSnakes {
				if hasCoord(curr, snake.Body) {
					if curr == snake.Head {
						movements[i].heads++
						if movements[i].nearestOpponent.distance == 0 || movements[i].nearestOpponent.distance >= depth {
							movements[i].nearestOpponent = Opponent{
								distance:  depth,
								length:    snake.Length,
								headCoord: snake.Head,
							}
						}
						if depth == 2 {
							movements[i].opponentInDmz = true
							log.Printf("Opponent located in DMZ")
						}
					}
					movements[i].obstacles++
					continue
				}
				// log.Printf("Nearest snake %v", movements[i].nearestOpponent)
			}

			for _, food := range board.Food {
				if curr == food {
					movements[i].food++
					if movements[i].distanceToFood == 0 && movements[i].distanceToFood > depth {
						movements[i].distanceToFood = depth
					}
				}
			}

			//log.Printf("Adding open spot %v to %v", curr, movements[i].root)
			movements[i].addOpenSpot(curr)
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
		if move.food > 0 && move.distanceToFood < best.distanceToFood {
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
		if move.food == 0 || move.distanceToFood > best.distanceToFood {
			best = move
		}
	}
	return best
}

func (moves WeightedMovementSet) bestMoveForDefense(you Battlesnake) WeightedMovement {
	best := moves[0]
	for i := 1; i < len(moves); i++ {
		move := moves[i]
		if move.nearestOpponent.distance > best.nearestOpponent.distance {
			best = move
		}
	}
	return best
}

func (moves WeightedMovementSet) bestMoveForOffense(you Battlesnake) WeightedMovement {
	var best WeightedMovement
	for i := 0; i < len(moves); i++ {
		move := moves[i]
		if move.nearestOpponent.distance < best.nearestOpponent.distance && move.nearestOpponent.length < you.Length {
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
	for i := 0; i < len(moves); i++ {
		move := moves[i]
		if (move.nearestOpponent.distance == 0 || move.nearestOpponent.distance > 2) && !move.opponentInDmz && !move.movingToCorner {
			safest = append(safest, move)
		}
		if move.nearestOpponent.distance == 0 || move.nearestOpponent.distance > 1 {
			safer = append(safer, move)
		}
	}

	rand.Seed(time.Now().Unix())
	if len(safest) > 0 {
		return safest[rand.Intn(len(safest))]
	} else {
		log.Printf("[%s] No safest moves available", you.Name)
	}

	if len(safer) > 0 {
		return safer[rand.Intn(len(safer))]
	} else {
		log.Printf("[%s] No safer moves available", you.Name)
	}

	return moves[rand.Intn(len(moves))]
}

func isCorner(c Coord, board Board) bool {
	corners := []Coord{{0, 0}, {0, board.Height - 1}, {board.Width - 1, 0}, {board.Width - 1, board.Width - 1}}
	if hasCoord(c, corners) {
		return true
	}
	return false
}

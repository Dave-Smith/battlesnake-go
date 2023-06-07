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
		w.open = make([]Coord, 0)
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

		// look out for corners.  might get trapped
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
								distance:  distanceTo(start, curr),
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
		opponent := move.nearestOpponent

		if (opponent.distance == 0 || opponent.distance > 2) && !move.opponentInDmz && !move.movingToCorner && len(move.open) >= you.Length {
			safest = append(safest, move)
		}
		if opponent.distance == 0 || opponent.distance > 1 && len(move.open) >= you.Length && (opponent.headCoord.X != you.Head.X || opponent.headCoord.Y != you.Head.Y) {
			safer = append(safer, move)
		}
	}

	rand.Seed(time.Now().Unix())
	if len(safest) > 0 {
		return mostOpenMoves(safest)
	} else {
		log.Printf("[%s] No safest moves available", you.Name)
	}

	if len(safer) > 0 {
		return mostOpenMoves(safer)
	} else {
		log.Printf("[%s] No safer moves available", you.Name)
	}

	return mostOpenMoves(moves)
}

func isCorner(c Coord, board Board) bool {
	corners := []Coord{{0, 0}, {0, board.Height - 1}, {board.Width - 1, 0}, {board.Width - 1, board.Width - 1}}
	if hasCoord(c, corners) {
		return true
	}
	return false
}

func isOnBorder(c Coord, board Board) bool {
	if c.X == 0 || c.X == board.Width-1 {
		return true
	}
	if c.Y == 0 || c.Y == board.Height-1 {
		return true
	}
	return false
}

func isOccupied(c Coord, allSnakes []Battlesnake) bool {
	for _, s := range allSnakes {
		if hasCoord(c, s.Body) {
			return true
		}
	}
	return false
}

func isOccupiedBySmallerSnake(c Coord, you Battlesnake, otherSnakes []Battlesnake) bool {
	for _, s := range otherSnakes {
		if c == s.Head && you.Length > s.Length {
			return true
		}
	}
	return false
}

func isNearSnakeHead(c Coord, snakeZone []HeadZone) bool {
	for _, s := range snakeZone {
		if hasCoord(c, s.Zone) {
			return true
		}
	}
	return false
}

func isNearSmallerSnakeHead(c Coord, you Battlesnake, snakeZone []HeadZone) bool {
	for _, s := range snakeZone {
		if hasCoord(c, s.Zone) && you.Length > s.SnakeLength {
			return true
		}
	}
	return false
}

func hasFood(c Coord, food []Coord) bool {
	for _, f := range food {
		if c == f {
			return true
		}
	}
	return false
}

func nearestSmallerSnake(c Coord, you Battlesnake, snakes []Battlesnake) Coord {
	var head Coord
	for _, s := range snakes {
		if you.Length > s.Length {
			head = c
		}
	}

	return head
}

func mostOpenMoves(possible WeightedMovementSet) WeightedMovement {
	move := possible[0]
	for _, m := range possible {
		if len(m.open) > len(move.open) {
			move = m
		}
	}
	return move
}

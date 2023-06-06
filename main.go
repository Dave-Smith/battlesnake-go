package main

import (
	"log"
	"math/rand"
	"time"
)

var priorMoves = make(map[string]string)

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func info() BattlesnakeInfoResponse {
	log.Println("Creating new battlesnake")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Dave-Smith",
		Color:      "#7ABF36",
		Head:       "all-seeing",
		Tail:       "do-sammy",
		Version:    "0.0.1-beta",
	}
}
func infoSalazar() BattlesnakeInfoResponse {
	log.Println("Creating new battlesnake salazar")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Dave-Smith",
		Color:      "#7ABF36",
		Head:       "all-seeing",
		Tail:       "do-sammy",
		Version:    "0.0.1-beta",
	}
}
func infoCoward() BattlesnakeInfoResponse {
	log.Println("Creating new battlesnake coward")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Dave-Smith",
		Color:      "#e6e600",
		Head:       "all-seeing",
		Tail:       "do-sammy",
		Version:    "0.0.1-beta",
	}
}
func infoVNext() BattlesnakeInfoResponse {
	log.Println("Creating new battlesnake vNext")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Dave-Smith",
		Color:      "#9af5b2",
		Head:       "silly",
		Tail:       "bolt",
		Version:    "0.0.1-beta",
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Printf("[%s] GAME START", state.You.Name)
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("[%s] GAME OVER\n\n", state.You.Name)
	log.Printf("[%s] Ending position: [%d,%d], Body: %v, ending health %d, ending length %d", state.You.Name, state.You.Head.X, state.You.Head.Y, state.You.Body, state.You.Health, state.You.Length)
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// We've included code to prevent your Battlesnake from moving backwards
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of your "neck"

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false

	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false

	} else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false

	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	// TODO: Step 1 - Prevent your Battlesnake from moving out of bounds
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height
	if myHead.X == 0 {
		isMoveSafe["left"] = false
	}
	if myHead.Y == 0 {
		isMoveSafe["down"] = false
	}
	if myHead.X == boardWidth-1 {
		isMoveSafe["right"] = false
	}
	if myHead.Y == boardHeight-1 {
		isMoveSafe["up"] = false
	}

	// TODO: Step 2 - Prevent your Battlesnake from colliding with itself
	mybody := state.You.Body
	for i := 0; i < len(mybody); i++ {
		part := mybody[i]
		if myHead.X+1 == part.X && myHead.Y == part.Y {
			isMoveSafe["right"] = false
		}
		if myHead.X-1 == part.X && myHead.Y == part.Y {
			isMoveSafe["left"] = false
		}
		if myHead.Y+1 == part.Y && myHead.X == part.X {
			isMoveSafe["up"] = false
		}
		if myHead.Y-1 == part.Y && myHead.X == part.X {
			isMoveSafe["left"] = false
		}
	}

	// TODO: Step 3 - Prevent your Battlesnake from colliding with other Battlesnakes
	// opponents := state.Board.Snakes

	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		log.Printf("[%s] MOVE %d: No safe moves detected! Moving down\n", state.You.Name, state.Turn)
		return BattlesnakeMoveResponse{Move: "down"}
	}

	// Choose a random move from the safe ones
	// log.Printf("Current position (%d,%d). Available safe moves %s", state.You.Head.X, state.You.Head.Y, strings.Join(safeMoves, ","))

	var nextMove string
	if val, ok := priorMoves[state.Game.ID]; len(safeMoves) > 1 && ok && isMoveSafe[val] {
		nextMove = val
	} else {
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
	}

	priorMoves[state.Game.ID] = nextMove

	//nextMove := safeMoves[rand.Intn(len(safeMoves))]

	// TODO: Step 4 - Move towards food instead of random, to regain health and survive longer
	// food := state.Board.Food

	log.Printf("MOVE %d: %s\n", state.Turn, nextMove)
	return BattlesnakeMoveResponse{Move: nextMove}
}

func moveSemiBlindWandering(state GameState) BattlesnakeMoveResponse {
	curr := state.You.Head
	gameMap := fillMap(state.Board, state.You)
	safe := safeMoves(curr, gameMap)
	log.Printf("[%s] Safe coordinates for next move %v", state.You.Name, safe)
	rand.Seed(time.Now().Unix())
	next := safe[rand.Intn(len(safe))]
	return BattlesnakeMoveResponse{Move: dir(curr, next)}
}

func moveLessBlindWandering(state GameState) BattlesnakeMoveResponse {
	curr := state.You.Head
	gameMap := fillMap(state.Board, state.You)
	opponents := make([]Battlesnake, 0)
	for i := 0; i < len(state.Board.Snakes); i++ {
		if state.You.ID != state.Board.Snakes[i].ID {
			opponents = append(opponents, state.Board.Snakes[i])
		}
	}
	safe := saferMoves(curr, gameMap, make([]Coord, 0), 4, opponents)
	log.Printf("[%s] Safe coordinates for next move %v", state.You.Name, safe)
	rand.Seed(time.Now().Unix())
	next := safe[rand.Intn(len(safe))]
	return BattlesnakeMoveResponse{Move: dir(curr, next)}
}

func moveSmart(state GameState) BattlesnakeMoveResponse {
	// scan the board for a possible moves
	//myLength := state.You.Length
	log.Printf("[%s] Starting Turn %d", state.You.Name, state.Turn)
	possible := fillToDepth(state.You.Head, state.You.Length, state.Board)
	possible = possible.avoidCertainDeath()

	var bestMove WeightedMovement

	otherSnakes := make([]Battlesnake, 0)
	for _, s := range state.Board.Snakes {
		if s.Name != state.You.Name {
			otherSnakes = append(otherSnakes, s)
		}
	}

	dangerishZones := MakeHeadZones(otherSnakes, state.Board, 2)
	log.Printf("Other snakes bubbles %v", dangerishZones)

	// possible offensive attack
	for _, p := range possible {
		for _, danger := range dangerishZones {
			if hasCoord(p.root, danger.Zone) && danger.SnakeLength < state.You.Length {
				return BattlesnakeMoveResponse{Move: p.movement.asString()}
			}
		}
	}

	// start game hunting for food
	if state.Turn < 70 || state.You.Health < 50 {
		food := NearestFoods(state.You, state.Board)
		if len(food) > 0 {
			for _, f := range food {
				if !f.Collision {
					isSafe := false
					movement := UseMovement(state.You.Head, f.Coords[0])
					for _, p := range possible {
						if p.movement == movement && !isCorner(p.root, state.Board) {
							isSafe = true
						}
					}
					for _, danger := range dangerishZones {
						if !isSafe {
							break
						}
						if hasCoord(f.Coords[0], danger.Zone) {
							break
						}
					}
					return BattlesnakeMoveResponse{Move: movement.asString()}
				}
			}
		}
	}

	// if board has more than 3 snakes, stay small
	if len(state.Board.Snakes) > 7 {
		if state.You.Health < 30 {
			bestMove = possible.bestMoveForFood(state.You)
		}
		if state.You.Health > 30 {
			bestMove = possible.bestMoveToAvoidFood(state.You)
		}
		bestMove = possible.bestMoveForRoaming(state.You)
		return BattlesnakeMoveResponse{Move: bestMove.movement.asString()}
	}

	// if 3 snakes, find food
	// if len(state.Board.Snakes) == 3 {
	// 	bestMove = possible.bestMoveForFood(state.You)
	// 	return BattlesnakeMoveResponse{Move: bestMove.movement.asString()}
	// }

	// head to head, roam
	starvingMove := possible.bestMoveToAvoidFood(state.You)
	defensiveMove := possible.bestMoveForRoaming(state.You)
	//offensiveMove := possible.bestMoveForOffense(state.You)
	if starvingMove.root == defensiveMove.root {
		bestMove = defensiveMove
	}

	if state.You.Health < 45 {
		bestMove = possible.bestMoveForFood(state.You)
		log.Printf("[%s] looking for food, %d moves away", state.You.Name, bestMove.distanceToFood)
	} else {
		bestMove = defensiveMove
	}

	return BattlesnakeMoveResponse{Move: bestMove.movement.asString(), Shout: "Avoiding food"}
}

func moveAggressive(state GameState) BattlesnakeMoveResponse {
	return BattlesnakeMoveResponse{Move: "up", Shout: "I'm coming after you"}
}

func movePassive(state GameState) BattlesnakeMoveResponse {
	return BattlesnakeMoveResponse{Move: "down", Shout: "Run away"}
}

func main() {
	RunServer()
}

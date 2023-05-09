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
	log.Println("INFO")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "Dave-Smith",
		Color:      "#7ABF36",
		Head:       "all-seeing",
		Tail:       "do-sammy",
		Version:    "0.0.1-beta",
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Println("GAME START")
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("GAME OVER\n\n")
	log.Printf("Ending position: %d,%d, Body: %v, ending health %d, ending length %d", state.You.Head.X, state.You.Head.Y, state.You.Body, state.You.Health, state.You.Length)
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
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
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
	log.Printf("Safe coordinates for next move %v", safe)
	rand.Seed(time.Now().Unix())
	next := safe[rand.Intn(len(safe))]
	return BattlesnakeMoveResponse{Move: dir(curr, next)}
}

func main() {
	RunServer()
}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type SnakeMoverFunc func(state GameState) BattlesnakeMoveResponse
type SnakeStartFunc func(state GameState)
type SnakeInfoFunc func() BattlesnakeInfoResponse
type SnakeEndFunc func(state GameState)

func HandleStart(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode start json, %s", err)
		return
	}

	start(state)

	// Nothing to respond with here
}
func HandleMove(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode move json, %s", err)
		return
	}
	log.Printf("Head position: (%d,%d), Body: %v, Health: %d, Length: %d", state.You.Head.X, state.You.Head.Y, state.You.Body, state.You.Health, state.You.Length)

	// response := move(state)
	// response := moveSemiBlindWandering(state)
	response := moveLessBlindWandering(state)

	log.Printf("Moving %s", response.Move)

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode move response, %s", err)
		return
	}
}

func HandleEnd(w http.ResponseWriter, r *http.Request) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		log.Printf("ERROR: Failed to decode end json, %s", err)
		return
	}

	end(state)

	// Nothing to respond with here
}

// Middleware

const ServerID = "battlesnake/dave-smith/salazar"
const ServerIdSal = "battlesnake/dave-smith/salazar"
const ServerIdVNext = "battlesnake/dave-smith/vNext"
const ServerIdCoward = "battlesnake/dave-smith/coward"
const ServerIdAgg = "battlesnake/dave-smith/aggressive"

func withServerID(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", ServerID)
		if next != nil {
			next(w, r)
		}
	}
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	response := info()

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("ERROR: Failed to encode info response, %s", err)
	}
}

func SnakeHandlerMove(mover SnakeMoverFunc, serverId string, next http.HandlerFunc) http.HandlerFunc {
	log.Printf("Move")
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", serverId)
		if next != nil {
			next(w, r)
		}

		state, err := unmarshalState(r)
		if err != nil {
			log.Printf("ERROR: Failed to decode move json, %s", err)
			return
		}
		log.Printf("[%s] Head position: (%d,%d), Body: %v, Health: %d, Length: %d", state.You.Name, state.You.Head.X, state.You.Head.Y, state.You.Body, state.You.Health, state.You.Length)

		response := mover(state)

		log.Printf("[%s] Moving %s", state.You.Name, response.Move)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("ERROR: Failed to encode move response, %s", err)
			return
		}
	}
}

func SnakeHandlerStart(starter SnakeStartFunc, serverId string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", serverId)
		if next != nil {
			next(w, r)
		}
		state, err := unmarshalState(r)
		log.Printf("[%s] Starting new game", state.You.Name)
		if err != nil {
			log.Printf("ERROR: Failed to decode move json, %s", err)
			return
		}
		log.Printf("[%s] Head position: (%d,%d), Body: %v, Health: %d, Length: %d", state.You.Name, state.You.Head.X, state.You.Head.Y, state.You.Body, state.You.Health, state.You.Length)

		starter(state)

		return
	}
}

func SnakeHandlerInfo(snakeInfo SnakeInfoFunc, serverId string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", serverId)

		response := snakeInfo()

		if next != nil {
			next(w, r)
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Printf("ERROR: Failed to encode info response, %s", err)
		}
	}
}

func SnakeHandlerEnd(gameEnd SnakeEndFunc, serverId string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		state, err := unmarshalState(r)
		if err != nil {
			log.Printf("ERROR: Failed to decode move json, %s", err)
			return
		}
		gameEnd(state)
	}
}

func unmarshalState(r *http.Request) (GameState, error) {
	state := GameState{}
	err := json.NewDecoder(r.Body).Decode(&state)
	if err != nil {
		return GameState{}, err
	}
	return state, nil
}

// Start Battlesnake Server

func RunServer() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	http.HandleFunc("/", withServerID(HandleIndex))
	http.HandleFunc("/start", withServerID(HandleStart))
	http.HandleFunc("/move", withServerID(HandleMove))
	http.HandleFunc("/end", withServerID(HandleEnd))

	http.HandleFunc("/agg", withServerID(HandleIndex))
	http.HandleFunc("/agg/start", withServerID(HandleStart))
	http.HandleFunc("/agg/move", withServerID(HandleMove))
	http.HandleFunc("/agg/end", withServerID(HandleEnd))

	// http.HandleFunc("/coward", withServerID(HandleIndex))
	http.HandleFunc("/coward", SnakeHandlerInfo(infoCoward, ServerIdCoward, nil))
	http.HandleFunc("/coward/start", SnakeHandlerStart(start, ServerIdCoward, nil))
	http.HandleFunc("/coward/move", SnakeHandlerMove(moveLessBlindWandering, ServerIdCoward, nil))
	http.HandleFunc("/coward/end", SnakeHandlerEnd(end, ServerIdCoward, nil))

	http.HandleFunc("/vnext", SnakeHandlerInfo(infoVNext, ServerIdVNext, nil))
	http.HandleFunc("/vnext/start", SnakeHandlerStart(start, ServerIdVNext, nil))
	http.HandleFunc("/vnext/move", SnakeHandlerMove(moveSmart, ServerIdVNext, nil))
	http.HandleFunc("/vnext/end", SnakeHandlerEnd(end, ServerIdVNext, nil))

	http.HandleFunc("/salazar", SnakeHandlerInfo(infoSalazar, ServerIdSal, nil))
	http.HandleFunc("/salazar/start", SnakeHandlerStart(start, ServerIdSal, nil))
	http.HandleFunc("/salazar/move", SnakeHandlerMove(moveLessBlindWandering, ServerIdSal, nil))
	http.HandleFunc("/salazar/end", SnakeHandlerEnd(end, ServerIdSal, nil))

	log.Printf("Running Battlesnake at http://0.0.0.0:%s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

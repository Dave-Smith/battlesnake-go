package main

type PossibleMove struct {
	Dir                      Movement
	Curr                     Coord
	Next                     Coord
	Prev                     Coord
	IsOffBoard               bool
	IsBackwards              bool
	HasFood                  bool
	IsOccupied               bool
	IsOccupiedBySmallerSnake bool
	IsCorner                 bool
	IsOnBorder               bool
	//IsOnObstacle             bool
	IsNearAnySnake     bool
	IsNearSmallerSnake bool
	//IsAdjacentToSmallerSnake bool
	//NearestFood              int
	//NearestSnake             int
	NearestSmallerSnake Coord

	//NearestSnakeBody         int
	//Paths                    []Path
}

type OtherSnake struct {
	Head         Coord
	Prev         Coord
	Tail         Coord
	Body         []Coord
	HeadlessBody []Coord
	PrevMove     Movement
	Length       int
	Health       int
	NextMoves    []Coord
	IsOnObstacle bool
	IsOnBorder   bool
	IsInCorner   bool
}

type GameFood struct {
	Location     Coord
	NearestSnake Coord
	DistToSnake  int
	DistToMe     int
	IsOnObstacle bool
}

func FindFood(you Battlesnake, other []Battlesnake, food []Coord) []GameFood {
	f := make([]GameFood, 0)

	return f
}

func FindOtherSnakes(snakes []Battlesnake, board Board) []OtherSnake {
	s := make([]OtherSnake, 0)

	return s
}

func FindNextMoves(you Battlesnake, other []Battlesnake, food []Coord, board Game) []PossibleMove {
	var m [4]PossibleMove
	return m[:]
}

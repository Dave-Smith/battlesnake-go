package main

import (
	"sort"
)

type Path struct {
	Source    Coord
	Target    Coord
	Coords    []Coord
	Collision bool
}

type byDistance []Path

func (p Path) Len() int {
	return len(p.Coords)
}

func AllFood(you Battlesnake, board Board) []Path {
	var paths = make([]Path, 0)

	for _, v := range board.Food {
		toFood := MakePaths(you.Head, v, board)
		paths = append(paths, toFood...)
	}
	return paths
}

func NearestFoods(you Battlesnake, board Board) []Path {
	var foods = AllFood(you, board)

	sort.Sort(byDistance(foods))
	return foods
}

func NearestFood(you Battlesnake, board Board) Path {
	nearest := NearestFoods(you, board)
	if len(nearest) > 0 {
		return nearest[0]
	}
	return Path{}
}

func (v byDistance) Len() int {
	return len(v)
}

func (v byDistance) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func (v byDistance) Less(i, j int) bool {
	return v[i].Len() < v[j].Len()
}

func MakePaths(source, target Coord, board Board) []Path {
	paths := make([]Path, 0)
	h := target.X - source.X
	v := target.Y - source.Y

	stepH := 1
	stepV := 1

	if h < 0 {
		stepH = -1
	}

	if v < 0 {
		stepV = -1
	}

	// horizontal first
	var coords = make([]Coord, 0)
	for i := 0; i < h*stepH; i++ {
		coords = append(coords, Coord{source.X + i + 1*stepH, source.Y})
	}
	for i := 0; i < v*stepV; i++ {
		coords = append(coords, Coord{source.X + h*stepH, source.Y + i + 1*stepV})
	}
	collision := false
	for _, c := range coords {
		for _, s := range board.Snakes {
			if hasCoord(c, s.Body) {
				collision = true
			}
		}
	}
	paths = append(paths, Path{Source: source, Target: target, Coords: coords, Collision: collision})

	coords = make([]Coord, 0)
	// vert first
	for i := 0; i < v*stepV; i++ {
		coords = append(coords, Coord{source.X, source.Y + i + 1*stepV})
	}
	for i := 0; i < h*stepH; i++ {
		coords = append(coords, Coord{source.X + i + 1*stepH, source.Y + v*stepV})
	}
	collision = false
	for _, c := range coords {
		for _, s := range board.Snakes {
			if hasCoord(c, s.Body) {
				collision = true
			}
		}
	}
	paths = append(paths, Path{Source: source, Target: target, Coords: coords, Collision: collision})

	return paths
}

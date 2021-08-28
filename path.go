package MarsRover

import (
	"errors"
	"fmt"
	"math"
)

type Point struct {
	X, Y     int
	Obstacle bool
}

// NewPoint creates a point with the given X and Y coordinates.
func NewPoint(x, y int) Point {
	return Point{
		X: x,
		Y: y,
	}
}

//Map is a 2D map.
type Map struct {
	cache     map[Point]bool
	obstacles []Point
	X, Y      int
}

// NewMap creates a map with the given width and height.
func NewMap(x, y int) *Map {
	return &Map{
		cache: make(map[Point]bool),
		X:     x,
		Y:     y,
	}
}

// AddObstacle adds as an obstacle to the map.
func (m *Map) AddObstacle(o Point) {
	m.obstacles = append(m.obstacles, o)
}

// IsObstacle returns whether the given point is an obstacle.
func (m *Map) IsObstacle(p Point) bool {
	if is, ok := m.cache[p]; ok {
		return is
	}

	var blocked bool
	for _, o := range m.obstacles {
		if o.Equal(p) {
			blocked = true
			break
		}
	}
	m.cache[p] = blocked
	return blocked
}

// Neighbours returns all neighbours within a distance of 1
func (m *Map) Neighbours(p Point) []Point {
	var result []Point
	var ne = [][2]int{
		{0, -1},
		{-1, 0},
		{0, 1},
		{1, 0},
	}
	for _, n := range ne {
		newPoint := NewPoint(p.X+n[0], p.Y+n[1])
		if m.IsObstacle(newPoint) {
			continue
		}
		result = append(result, newPoint)
	}
	return result
}

// Path returns all the points that need to be visited in order to get from start
// to goal in the least amount of steps
func (m *Map) Path(start, goal Point, direction string) (string, []Point, error) {
	closedList := make(map[Point]struct{})
	openList := map[Point]struct{}{start: {}}
	fScore := map[Point]float64{start: 0}
	gScore := map[Point]float64{start: 0}
	cameFrom := make(map[Point]Point)

	for len(openList) > 0 {
		var current Point
		var currentScore = math.MaxFloat64

		for node := range openList {
			score := fScore[node]
			if score < currentScore {
				current = node
				currentScore = score
			}
		}

		delete(openList, current)
		closedList[current] = struct{}{}

		if current.Equal(goal) {
			var path []Point

			for {
				path = append(path, current)
				var ok bool
				current, ok = cameFrom[current]
				if !ok {
					break
				}
			}

			var revPath = make([]Point, len(path))
			for i, node := range path {
				revPath[len(path)-1-i] = node
			}

			return pathToCommand(direction, revPath), revPath, nil
		}

	NeighborsLoop:
		for _, n := range m.Neighbours(current) {
			if _, ok := closedList[n]; ok {
				continue
			}

			//https://en.wikipedia.org/wiki/A*_search_algorithm#Pseudocode
			gScore[n] = gScore[current] + 1
			fScore[n] = gScore[n] + n.Distance(goal)

			for on := range openList {
				if n.Equal(on) {
					if g, ok := gScore[on]; ok && gScore[n] > g {
						continue NeighborsLoop
					}
				}
			}

			openList[n] = struct{}{}
			cameFrom[n] = current
		}
	}

	return "", nil, errors.New("could not reach goal")
}

func pathToCommand(direction string, points []Point) string {
	var cmd string
	currentDirection := direction
	for index := range points {
		if index <= len(points)-2 {
			currtDirection, command := TwoNeighborPath(currentDirection, points[index], points[index+1])
			cmd += command
			currentDirection = currtDirection
		}
	}
	return cmd
}

func TwoNeighborPath(direction string, start, end Point) (string, string) {
	var cmd, currentDirection string
	rover := Rover{
		X:         0,
		Y:         0,
		Direction: direction,
	}
	if start.X < end.X && start.Y == end.Y {
		currentDirection = "EAST"
		cmd += rover.CorrectDirection("EAST")
	} else if start.X > end.X && start.Y == end.Y {
		currentDirection = "WEST"
		cmd += rover.CorrectDirection("WEST")
	} else if start.X == end.X && start.Y < end.Y {
		currentDirection = "NORTH"
		cmd += rover.CorrectDirection("NORTH")
	} else if start.X == end.X && start.Y > end.Y {
		currentDirection = "SOUTH"
		cmd += rover.CorrectDirection("SOUTH")
	} else {
		currentDirection = direction
		fmt.Println("these two points are not neighbor")
		return currentDirection, cmd
	}
	return currentDirection, cmd + "F"
}

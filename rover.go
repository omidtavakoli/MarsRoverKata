package MarsRover

import (
	"errors"
	"fmt"
	"strings"
)

type Rover struct {
	X, Y      int
	Direction string
}

func CreateRover(position [2]int, direction string) *Rover {
	rover := Rover{
		X:         position[0],
		Y:         position[1],
		Direction: direction,
	}
	return &rover
}

func (r *Rover) CommandReceiver(command string, obstacles map[string]bool) string {
	for i := 0; i < len(command); i++ {
		cmd := strings.ToLower(string(command[i]))
		err := r.Move(cmd, obstacles)
		if err != nil {
			msg := fmt.Sprintf("(%d, %d) %s %s", r.X, r.Y, r.Direction, err.Error())
			return msg
		}
	}
	msg := fmt.Sprintf("(%d, %d) %s", r.X, r.Y, r.Direction)
	return msg
}

func (r *Rover) Move(move string, obstacles map[string]bool) error {
	x, y := r.X, r.Y
	if move == "l" {
		r.Turn("LEFT")
	} else if move == "r" {
		r.Turn("RIGHT")
	} else if (r.Direction == "NORTH" && move == "f") || (r.Direction == "SOUTH" && move == "b") {
		x, y = r.X, r.Y+1
	} else if (r.Direction == "SOUTH" && move == "f") || (r.Direction == "NORTH" && move == "b") {
		x, y = r.X, r.Y-1
	} else if (r.Direction == "WEST" && move == "f") || (r.Direction == "EAST" && move == "b") {
		x, y = r.X-1, r.Y
	} else if (r.Direction == "EAST" && move == "f") || (r.Direction == "WEST" && move == "b") {
		x, y = r.X+1, r.Y
	}

	if checkObstacle(x, y, obstacles) {
		return errors.New("STOPPED")
	}
	r.X, r.Y = x, y
	return nil
}

//Turn turns the rover on the rotation commands
func (r *Rover) Turn(direction string) {
	directions := map[string][]string{"WEST": {"SOUTH", "NORTH"}, "NORTH": {"WEST", "EAST"}, "EAST": {"NORTH", "SOUTH"}, "SOUTH": {"EAST", "WEST"}}
	if direction == "RIGHT" || direction == "r" || direction == "R" {
		r.Direction = directions[r.Direction][1]
	} else {
		r.Direction = directions[r.Direction][0]
	}
}

//checkObstacle checks if obstacle happened
func checkObstacle(x, y int, obstacles map[string]bool) bool {
	point := fmt.Sprintf("%d,%d", x, y)
	return obstacles[point]
}

//obstacle makes a hashmap to access obstacle faster
func createObstaclesHashMap(obstacles [][2]int) map[string]bool {
	obstacleMap := make(map[string]bool, len(obstacles))
	for _, obst := range obstacles {
		position := fmt.Sprintf("%d,%d", obst[0], obst[1])
		obstacleMap[position] = true
	}
	return obstacleMap
}

func CalculateDirectionCost(now, want string) [2]int {
	//0: left rotation & 1:right rotation
	if now == want {
		return [2]int{4, 4}
	} else {
		cost := [2]int{}
		directions := map[string][]string{"WEST": {"SOUTH", "NORTH"}, "NORTH": {"WEST", "EAST"}, "EAST": {"NORTH", "SOUTH"}, "SOUTH": {"EAST", "WEST"}}

		for key := range cost {
			cursor := now
			for {
				cursor = directions[cursor][key]
				cost[key]++
				if cursor == want {
					break
				}
			}
		}

		return cost
	}
}

func (r *Rover) CorrectDirection(direction string) string {
	var command string
	rotation := "L"
	rotationCost := CalculateDirectionCost(r.Direction, direction)
	if rotationCost[1] < rotationCost[0] {
		rotation = "R"
	}

	if r.Direction != direction {
		for {
			if r.Direction == direction {
				break
			}
			r.Turn(rotation)
			command += rotation
		}
	}

	return command
}

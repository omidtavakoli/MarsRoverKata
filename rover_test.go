package MarsRover

import (
	"errors"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestCreateRover(t *testing.T) {
	rover := CreateRover([2]int{3, 4}, "NORTH")
	assert.Equal(t, 3, rover.X)
	assert.Equal(t, 4, rover.Y)
	assert.Equal(t, "NORTH", rover.Direction)
}

func TestMoverError(t *testing.T) {
	tests := []struct {
		description string
		now         [2]int
		obstacles   [][2]int
		direction   string
		command     string
		want        error
	}{
		{
			description: "obstacle check",
			now:         [2]int{4, 0},
			obstacles:   [][2]int{{5, 0}},
			direction:   "EAST",
			command:     "f",
			want:        errors.New("STOPPED"),
		},
		{
			description: "obstacle check",
			now:         [2]int{3, 2},
			obstacles:   [][2]int{{3, 3}},
			direction:   "NORTH",
			command:     "f",
			want:        errors.New("STOPPED"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			rover := CreateRover(tt.now, tt.direction)
			obstaclesHashMap := createObstaclesHashMap(tt.obstacles)
			msg := rover.Move(tt.command, obstaclesHashMap)
			assert.Equal(t, msg, tt.want)
		})
	}

}

func TestMover(t *testing.T) {
	tests := []struct {
		description string
		now         [2]int
		obstacles   [][2]int
		direction   string
		command     string
		want        [2]int
	}{
		{
			description: "mover check",
			now:         [2]int{4, 0},
			obstacles:   [][2]int{},
			direction:   "WEST",
			command:     "f",
			want:        [2]int{3, 0},
		},
		{
			description: "mover check",
			now:         [2]int{3, 2},
			obstacles:   [][2]int{},
			direction:   "NORTH",
			command:     "f",
			want:        [2]int{3, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			rover := CreateRover(tt.now, tt.direction)
			obstaclesHashMap := createObstaclesHashMap(tt.obstacles)
			rover.Move(tt.command, obstaclesHashMap)
			assert.Equal(t, tt.want[0], rover.X)
			assert.Equal(t, tt.want[1], rover.Y)
		})
	}

}

func TestCommandReceiver(t *testing.T) {
	tests := []struct {
		description string
		now         [2]int
		obstacles   [][2]int
		direction   string
		command     string
		want        [2]int
		message     string
	}{
		{
			description: "cmd receiver check",
			now:         [2]int{4, 0},
			obstacles:   [][2]int{{3, 1}},
			direction:   "EAST",
			command:     "LFFFFLF",
			want:        [2]int{3, 4},
			message:     "(3, 4) WEST",
		},
		{
			description: "cmd receiver check",
			now:         [2]int{4, 2},
			obstacles:   [][2]int{},
			direction:   "EAST",
			command:     "FLFFFRFLB",
			want:        [2]int{6, 4},
			message:     "(6, 4) NORTH",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			rover := CreateRover(tt.now, tt.direction)
			obstaclesHashMap := createObstaclesHashMap(tt.obstacles)
			msg := rover.CommandReceiver(tt.command, obstaclesHashMap)
			assert.Equal(t, tt.want[0], rover.X)
			assert.Equal(t, tt.want[1], rover.Y)
			assert.Equal(t, tt.message, msg)
		})
	}

}

func TestTurn(t *testing.T) {
	rover := CreateRover([2]int{4, 2}, "WEST")

	tests := []struct {
		description string
		command     string
		want        string
	}{
		{
			description: "turn left",
			command:     "LEFT",
			want:        "SOUTH",
		},
		{
			description: "turn right",
			command:     "RIGHT",
			want:        "WEST",
		},
		{
			description: "turn right",
			command:     "RIGHT",
			want:        "NORTH",
		},
		{
			description: "turn right",
			command:     "RIGHT",
			want:        "EAST",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			rover.Turn(tt.command)
			assert.Equal(t, tt.want, rover.Direction)
		})
	}
}

func TestObstacles(t *testing.T) {
	obstacles := [][2]int{{1, 4}, {3, 5}, {7, 4}}
	rover := CreateRover([2]int{3, 4}, "WEST")
	obstaclesHashMap := createObstaclesHashMap(obstacles)
	msg := rover.CommandReceiver("FF", obstaclesHashMap)

	assert.Equal(t, "(2, 4) WEST STOPPED", msg)
}

func TestCalculateDirectionCost(t *testing.T) {
	tests := []struct {
		description string
		now         string
		want        string
		cost        [2]int
	}{
		{
			description: "both 4",
			now:         "SOUTH",
			want:        "SOUTH",
			cost:        [2]int{4, 4},
		},
		{
			description: "left 1 and right 3",
			now:         "NORTH",
			want:        "WEST",
			cost:        [2]int{1, 3},
		},
		{
			description: "both 2",
			now:         "EAST",
			want:        "WEST",
			cost:        [2]int{2, 2},
		},
		{
			description: "both",
			now:         "EAST",
			want:        "WEST",
			cost:        [2]int{2, 2},
		},
		{
			description: "",
			now:         "EAST",
			want:        "NORTH",
			cost:        [2]int{1, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cost := CalculateDirectionCost(tt.now, tt.want)
			assert.Equal(t, tt.cost, cost)
		})
	}
}

func TestCorrectDirection(t *testing.T) {
	tests := []struct {
		description string
		now         string
		want        string
		command     string
	}{
		{
			description: "both 4",
			now:         "SOUTH",
			want:        "SOUTH",
			command:     "",
		},
		{
			description: "just left",
			now:         "NORTH",
			want:        "WEST",
			command:     "L",
		},
		{
			description: "just right",
			now:         "EAST",
			want:        "SOUTH",
			command:     "R",
		},
		{
			description: "just right",
			now:         "EAST",
			want:        "NORTH",
			command:     "L",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			rover := CreateRover([2]int{10, 50}, tt.now)
			cmd := rover.CorrectDirection(tt.want)
			assert.Equal(t, tt.command, cmd)
		})
	}
}

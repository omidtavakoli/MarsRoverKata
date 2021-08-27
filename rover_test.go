package MarsRover

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestCreateRover(t *testing.T) {
	rover := CreateRover([]int{3, 4}, "NORTH")
	assert.Equal(t, 3, rover.X)
	assert.Equal(t, 4, rover.Y)
	assert.Equal(t, "n", rover.Direction)
}

func TestMover(t *testing.T) {
	rover := CreateRover([]int{5, 8}, "SOUTH")
	obstacles := make(map[string]bool, 0)
	_ = rover.Move("f", obstacles)
	assert.Equal(t, 7, rover.Y)
}

func TestCommandReceiver(t *testing.T) {
	mars := Init([][]int{})
	rover := CreateRover([]int{8, 10}, "WEST")
	msg := rover.CommandReceiver("fff", mars)
	assert.Equal(t, 5, rover.X)
	assert.Equal(t, "(5, 10) WEST", msg)
}

func TestCreateRoverSecond(t *testing.T) {
	mars := Init([][]int{})
	rover := CreateRover([]int{4, 2}, "EAST")
	msg := rover.CommandReceiver("FLFFFRFLB", mars)
	assert.Equal(t, "(6, 4) NORTH", msg)
}

func TestTurn(t *testing.T) {
	rover := CreateRover([]int{4, 2}, "WEST")

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
			fmt.Println(rover)
			assert.Equal(t, tt.want, rover.Direction)
		})
	}
}

func TestObstacles(t *testing.T) {
	mars := Init([][]int{{1, 4}, {3, 5}, {7, 4}})
	rover := CreateRover([]int{3, 4}, "WEST")
	msg := rover.CommandReceiver("FF", mars)

	assert.Equal(t, "(2, 4) WEST STOPPED", msg)
}

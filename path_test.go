package MarsRover

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestPath(t *testing.T) {
	m := NewMap(5, 5)
	obstacles := [][2]int{{3, 3}, {1, 4}, {2, 4}, {3, 1}}
	for _, obst := range obstacles {
		point := Point{
			X:        obst[0],
			Y:        obst[1],
			Obstacle: true,
		}
		m.AddObstacle(point)
		m.cache[point] = true
	}
	start := NewPoint(2, 1)
	end := NewPoint(3, 4)
	command, _, err := m.Path(start, end, "SOUTH")
	if err != nil {
		t.Error(err)
	}
	if command != "LLFRFFLFFLF" {
		t.Error("command is wrong")
	}
}

func TestTwoNeighborPath(t *testing.T) {
	tests := []struct {
		description string
		direction   string
		point1      Point
		point2      Point
		want        string
	}{
		{
			description: "distance func test",
			direction:   "SOUTH",
			point1: Point{
				X: 2,
				Y: 1,
			},
			point2: Point{
				X: 3,
				Y: 1,
			},
			want: "LF",
		},
		{
			description: "distance func test",
			direction:   "WEST",
			point1: Point{
				X: 2,
				Y: 1,
			},
			point2: Point{
				X: 2,
				Y: 2,
			},
			want: "RF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			_, msg := TwoNeighborPath(tt.direction, tt.point1, tt.point2)
			assert.Equal(t, tt.want, msg)
		})
	}

}

func TestPathToCommand(t *testing.T) {
	tests := []struct {
		description string
		direction   string
		points      []Point
		want        string
	}{
		{
			description: "path func test",
			direction:   "SOUTH",
			points: []Point{{
				X: 2,
				Y: 1,
			},
				{
					X: 3,
					Y: 1,
				}},
			want: "LF",
		},
		{
			description: "path func test",
			direction:   "WEST",
			points: []Point{{
				X: 2,
				Y: 1,
			},
				{
					X: 1,
					Y: 1,
				}},
			want: "F",
		},
		{
			description: "path func test",
			direction:   "WEST",
			points: []Point{{
				X: 2,
				Y: 1,
			},
				{
					X: 2,
					Y: 2,
				}, {
					X: 3,
					Y: 2,
				},
				{
					X: 4,
					Y: 2,
				}, {
					X: 4,
					Y: 3,
				},
				{
					X: 4,
					Y: 4,
				}, {
					X: 3,
					Y: 4,
				},
			},
			want: "RFRFFLFFLF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cmd := pathToCommand(tt.direction, tt.points)
			assert.Equal(t, tt.want, cmd)
		})
	}

}

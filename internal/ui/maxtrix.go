package ui

import (
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type matrixTickMsg struct{}

func matrixTicker() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return matrixTickMsg{}
	})
}

type Stream struct {
	col        int
	heads      []Head
	column     []rune
	bright     []float64
	lastUpdate time.Time
}
type Head struct {
	row   int
	speed float64
}

func (m *model) initStreams(width, height int) {
	m.streams = make([]Stream, width)
	for i := 0; i < width; i++ {
		col := make([]rune, height)
		bright := make([]float64, height)
		for j := range col {
			col[j] = randomMatrixChar()
			bright[j] = rand.Float64() * 0.3
		}
		heads := []Head{
			{row: -rand.Intn(height), speed: 5 + rand.Float64()*10},
			{row: -rand.Intn(height), speed: 5 + rand.Float64()*10},
		}
		m.streams[i] = Stream{
			col:        i,
			heads:      heads,
			column:     col,
			bright:     bright,
			lastUpdate: time.Now(),
		}
	}
}

// func newStream(col, h int) *Stream {
// 	s := &Stream{
// 		col:        col,
// 		lastUpdate: time.Now(),
// 		speed:      rand.Float64()*15 + 10,
// 	}
// 	s.heads = append(s.heads, Head{row: -rand.Intn(h)})
// 	return s
// }

func (s *Stream) step(height int) {
	now := time.Now()
	dt := now.Sub(s.lastUpdate).Seconds()
	s.lastUpdate = now

	for i := range s.bright {
		s.bright[i] *= 0.8 // decay factor, adjust to taste
		if s.bright[i] < 0.1 {
			s.bright[i] = 0
			s.column[i] = randomMatrixChar()
		}
	}

	//update heads
	for i := range s.heads {
		head := &s.heads[i]
		head.row += int(head.speed * dt)
		if head.row >= height {
			head.row = -rand.Intn(height)
			head.speed = 5 + rand.Float64()*10
		}

		if head.row >= 0 && head.row < height {
			s.column[head.row] = randomMatrixChar()
			s.bright[head.row] = 1.0
		}
	}
}

func randomMatrixChar() rune {
	r := rune(rand.Intn(94) + 33) // 33–126 are printable ASCII
	return r
}

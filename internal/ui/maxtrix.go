package ui

import (
	"math/rand"
	"strings"
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
	heads  []Head
	column []rune
	bright []float64
	height int
}
type Head struct {
	row         int
	speed       int
	trailLength int
}

// Initialize streams
func (m *model) initStreams(width, height int) {
	m.streams = make([]Stream, width)
	for i := 0; i < width; i++ {
		col := make([]rune, height)
		bright := make([]float64, height)

		// fill with random chars and base faint brightness
		for j := range col {
			col[j] = randomMatrixChar()
			bright[j] = 0.05 + rand.Float64()*0.1
		}

		// multiple heads per column
		numHeads := 4 + rand.Intn(3) // 4–6 heads per column
		heads := make([]Head, numHeads)
		for h := range heads {
			heads[h] = Head{
				row:   rand.Intn(height), // random start anywhere
				speed: 1 + rand.Intn(2),  // speed 1–2
			}
		}

		m.streams[i] = Stream{
			heads:  heads,
			column: col,
			bright: bright,
		}
	}
}

// Step the stream
func (s *Stream) step(height int) {
	// fade previous brightness
	for i := range s.bright {
		s.bright[i] *= 0.8
		if s.bright[i] < 0.02 {
			s.bright[i] = 0.05 // faint base glow
			s.column[i] = randomMatrixChar()
		}
	}

	// move heads and draw trails
	for i := range s.heads {
		head := &s.heads[i]
		head.row += head.speed
		if head.row >= height {
			head.row = -rand.Intn(height / 2)
			head.speed = 1 + rand.Intn(2)
		}

		trailLength := 5 + rand.Intn(6) // 5–10 characters
		for t := 0; t < trailLength; t++ {
			pos := head.row - t
			if pos >= 0 && pos < height {
				s.column[pos] = randomMatrixChar()
				s.bright[pos] = 0.1 + 0.9*(float64(t)/float64(trailLength))
			}
		}
	}
}

func randomMatrixChar() rune {
	r := rune(rand.Intn(94) + 33) // 33–126 are printable ASCII
	return r
}

func unicodeCorrupt(r rune) rune {
	lookalikes := []rune{
		'@', '#', '$', '%', '&', '*', '!', '?', '+', '-', '=',
		'~', '^', ':', ';', '.', ',', '|', '/', '\\',
	}
	if rand.Float64() < 0.5 {
		return lookalikes[rand.Intn(len(lookalikes))]
	}

	return r
}

func (m *model) generateMatrixFrame() []string {
	lines := make([]string, m.modalHeight)

	for row := 0; row < m.modalHeight; row++ {
		var b strings.Builder
		for col := 0; col < m.modalWidth-4; col++ {
			s := &m.streams[col]

			char := s.column[row]
			brightness := s.bright[row]

			// overlay log content if exists
			logIndex := len(m.extractionLogs) - m.modalHeight + row
			if logIndex >= 0 && logIndex < len(m.extractionLogs) {
				line := padLogLine(m.extractionLogs[logIndex], m.modalWidth-4)
				if col < len(line) {
					orig := rune(line[col])
					r := rand.Float64()
					switch {
					case r < 0.5:
						char = orig
					case r < 0.8:
						char = unicodeCorrupt(orig)
					case r < 0.9:
						char = ' '
					default:
						char = orig
					}
				}
			}

			b.WriteString(renderChar(char, brightness))
		}
		lines[row] = b.String()
	}

	return lines
}

func padLogLine(line string, width int) string {
	if len(line) < width {
		// Fill empty space with spaces or light dots for consistent width perception
		line += strings.Repeat(string(randomMatrixChar()), width-len(line))
	} else if len(line) > width {
		line = line[:width]
	}
	return line
}

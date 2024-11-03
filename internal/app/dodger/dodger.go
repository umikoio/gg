package dodger

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type spawnBlockMsg struct{}
type moveBlockMsg struct{}

type vector struct {
	X int
	Y int
}

type model struct {
	size   vector
	player vector
	blocks []vector
	score  int

	blockStyle  lipgloss.Style
	playerStyle lipgloss.Style
}

func initialModel() tea.Model {
	size := vector{30, 20}
	return model{
		size:   size,
		player: vector{int(size.X / 2), size.Y - 1},
		blocks: []vector{}, score: 0,
		blockStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("#cccccc")),
		playerStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("#aaaaff")),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "left", "a":
			m.player.X--
			if m.player.X < 0 {
				m.player.X = m.size.X - 1
			}
		case "right", "d":
			m.player.X++
			if m.player.X >= m.size.X {
				m.player.X = 0
			}
		}
	case spawnBlockMsg:
		m.blocks = append(m.blocks, vector{rand.Intn(m.size.X), 0})
	case moveBlockMsg:
		m.MoveBlocks()
	}

	for _, b := range m.blocks {
		if b.X == m.player.X && b.Y == m.player.Y {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	for y := 0; y < m.size.Y; y++ {
		for x := 0; x < m.size.X; x++ {
			drew := false
			for _, b := range m.blocks {
				if b.X == x && b.Y == y {
					s += m.blockStyle.Render(string(rune(0x2022)))
					drew = true
				}
			}
			if !drew {
				if x == m.player.X && y == m.player.Y {
					s += m.playerStyle.Render(string(rune(0x2205)))
				} else {
					s += " "
				}
			}
		}
		s += "\n"
	}

	s += fmt.Sprintf("\nScore: %d\n", m.score)

	return s
}

func (m *model) MoveBlocks() {
	for i, _ := range m.blocks {
		m.blocks[i].Y++
	}

	for i, _ := range m.blocks {
		if m.blocks[i].Y > m.size.Y {
			m.blocks = append(m.blocks[:i], m.blocks[i+1:]...)
			m.score++
			// There can only be one block to be removed every time.
			break
		}
	}
}

func Run() {
	p := tea.NewProgram(initialModel())

	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			p.Send(spawnBlockMsg{})
			p.Send(moveBlockMsg{})
		}
	}()

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

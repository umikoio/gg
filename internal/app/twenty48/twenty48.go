package twenty48

import (
	"math/rand"
	"time"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	score int
	colors map[int]lipgloss.Style
	grid [4][4]int
}

func initialModel() tea.Model {
	defaultStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff"))

	m := model {
		colors: map[int]lipgloss.Style{
			0: defaultStyle.Background(lipgloss.Color("#3c3a32")), 
			2: defaultStyle.Background(lipgloss.Color("#eee4da")), 
			4: defaultStyle.Background(lipgloss.Color("#ede0c8")), 
			8: defaultStyle.Background(lipgloss.Color("#f2b179")), 
			16:defaultStyle.Background(lipgloss.Color("#f59563")), 
			32:defaultStyle.Background(lipgloss.Color("#f67c5f")), 
			64:defaultStyle.Background(lipgloss.Color("#f65e3b")), 
			128:defaultStyle.Background(lipgloss.Color("#edcf72")), 
			256:defaultStyle.Background(lipgloss.Color("#edcc61")), 
			512:defaultStyle.Background(lipgloss.Color("#edc850")), 
			1024:defaultStyle.Background(lipgloss.Color("#edc53f")), 
			2048:defaultStyle.Background(lipgloss.Color("#edc22e")), 
		},
		grid: [4][4]int{},
	}

	m.AddTile()
	m.AddTile()
	return m
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
			m.MergeTilesLeft()
			if !m.AddTile() {
				return m, tea.Quit
			}
		case "down", "s":
			m.Rotate90(false)
			m.MergeTilesLeft()
			m.Rotate90(true)
			if !m.AddTile() {
				return m, tea.Quit
			}
		case "up", "w":
			m.Rotate90(true)
			m.MergeTilesLeft()
			m.Rotate90(false)
			if !m.AddTile() {
				return m, tea.Quit
			}
		case "right", "d":
			m.Rotate90(false)
			m.Rotate90(false)
			m.MergeTilesLeft()
			m.Rotate90(true)
			m.Rotate90(true)
			if !m.AddTile() {
				return m, tea.Quit
			}}
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			stringifiedNum := strconv.Itoa(m.grid[y][x])

			for i := 0; i < 4 - len(stringifiedNum); i++ {
				s += m.colors[m.grid[y][x]].Render(" ")
			}
			s += m.colors[m.grid[y][x]].Render(stringifiedNum)
		}
		s += "\n"
	}

	return s
}

func (m *model) MergeTilesLeft() {
	for i, _ := range m.grid {
		stopMerge := 0
		for j := 1; j < len(m.grid[i]); j++ {
			if m.grid[i][j] != 0 {
				for k := j; k > stopMerge; k-- {
					if m.grid[i][k-1] == 0 {
						m.grid[i][k-1] = m.grid[i][k]
						m.grid[i][k] = 0
					} else if m.grid[i][k-1] == m.grid[i][k] {
						m.grid[i][k-1] += m.grid[i][k]
						m.grid[i][k] = 0
						stopMerge = k
						break
					} else {
						break
					}
				}
			}
		}
	}
}

func (m *model) AddTile() bool {
	empty := []int{}
	for y, row := range m.grid {
		for x, cell := range row {
			if cell == 0 {
				empty = append(empty, y * 4 + x)
			}
		}
	}

	if len(empty) == 0 {
		return false
	}
	
	rndSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSrc)

	cell := empty[rnd.Intn(len(empty))]

	m.grid[cell/len(m.grid)][cell % len(m.grid)] = 2
	
	return true
}

func (m *model) Rotate90(counterClockWise bool) {
	rotatedGrid := [4][4]int{}
	for i, row := range m.grid {
		rotatedGrid[i] = [4]int{}
		for j := range row {
			if counterClockWise {
				rotatedGrid[i][j] = m.grid[j][len(m.grid)-i-1]
			} else {
				rotatedGrid[i][j] = m.grid[len(m.grid)-j-1][i]
			}
		}
	}
	m.grid = rotatedGrid
}

func Run() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

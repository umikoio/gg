package twenty48

import (
	"math/rand"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	score  int
	colors map[int]lipgloss.Style
	grid   [4][4]int
}

func initialModel() tea.Model {
	defaultStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#f9f6f2"))
	c := func(s string) lipgloss.Color {
		return lipgloss.Color(s)
	}

	m := model{
		colors: map[int]lipgloss.Style{
			0:    defaultStyle.Background(c("#3c3a32")),
			2:    defaultStyle.Background(c("#eee4da")).Foreground(c("#000000")),
			4:    defaultStyle.Background(c("#ede0c8")).Foreground(c("#000000")),
			8:    defaultStyle.Background(c("#f2b179")).Foreground(c("#f9f6f2")),
			16:   defaultStyle.Background(c("#f59563")).Foreground(c("#f9f6f2")),
			32:   defaultStyle.Background(c("#f67c5f")).Foreground(c("#f9f6f2")),
			64:   defaultStyle.Background(c("#f65e3b")).Foreground(c("#f9f6f2")),
			128:  defaultStyle.Background(c("#edcf72")).Foreground(c("#f9f6f2")),
			256:  defaultStyle.Background(c("#edcc61")).Foreground(c("#f9f6f2")),
			512:  defaultStyle.Background(c("#edc850")).Foreground(c("#f9f6f2")),
			1024: defaultStyle.Background(c("#edc53f")).Foreground(c("#f9f6f2")),
			2048: defaultStyle.Background(c("#edc22e")).Foreground(c("#f9f6f2")),
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
			}
		}
	}

	if m.CheckForWin() {
		return m, tea.Quit
	}

	return m, nil
}

func (m model) View() string {
	s := ""

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			s += m.colors[m.grid[y][x]].Render("      ")
		}
		s += "\n"
		for x := 0; x < 4; x++ {
			stringifiedNum := strconv.Itoa(m.grid[y][x])
			if stringifiedNum == "0" {
				stringifiedNum = "."
			}

			for i := 0; i < 5-len(stringifiedNum); i++ {
				s += m.colors[m.grid[y][x]].Render(" ")
			}
			s += m.colors[m.grid[y][x]].Render(stringifiedNum + " ")
		}
		s += "\n"
		for x := 0; x < 4; x++ {
			s += m.colors[m.grid[y][x]].Render("      ")
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
				empty = append(empty, y*4+x)
			}
		}
	}

	if len(empty) == 0 {
		return false
	}

	rndSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSrc)

	cell := empty[rnd.Intn(len(empty))]

	if rnd.Intn(10) < 9 {
		m.grid[cell/len(m.grid)][cell%len(m.grid)] = 2
	} else {
		m.grid[cell/len(m.grid)][cell%len(m.grid)] = 4
	}

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

func (m model) CheckForWin() bool {
	for _, row := range m.grid {
		for x, _ := range row {
			if row[x] == 2048 {
				return true
			}
		}
	}

	return false
}

func Run() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

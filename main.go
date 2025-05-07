package main

import (
	"fmt"
	"os"
	"os/user"
	"syscall"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc     string
	callingFunction func() string
}

func (i item) Title() string { return i.title }

// func (i item) Description() string { return i.desc }
func (i item) Description() string {
	if i.callingFunction != nil {
		return i.callingFunction()
	}
	return ""
}

func (i item) FilterValue() string { return i.title }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func main() {
	items := []list.Item{

		item{title: "Get available memory.", desc: "",
			callingFunction: func() string {
				return fmt.Sprintf("%d GB", ram()/GB)
			},
		},
		item{title: "Get the current logged in user.", desc: "",
			callingFunction: func() string {
				u, _ := user.Current()
				return u.Username
			},
		},
		item{title: "Quit", desc: ""},
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func ram() uint64 {
	totalRamBit := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(totalRamBit); err != nil {
		return 0
	}
	return uint64(totalRamBit.Totalram) * uint64(totalRamBit.Unit)
}

package main

import (
    "fmt"
    "os"
    "os/exec"
    "log"
    "strings"

    "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
    title string
}

func (i item) Title() string { return i.title }
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

// Uses exec to run AWS CLI command to list profiles
func GetProfiles() []string {
    cmd := exec.Command("aws", "configure", "list-profiles")
    cmd.Stderr = os.Stderr
    data, err := cmd.Output()

    if err != nil {
        log.Fatalf("Failed to call cmd.Output(): %v", err)
    }

   profiles := strings.Split(string(data), "\n")

   return profiles
}

func main() {
    items := []list.Item{}

    for _, profile := range GetProfiles() {
        items = append(items, profile)
    }

    fmt.Sprintf("Profiles: %v", items)
}

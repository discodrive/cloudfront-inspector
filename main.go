package main

import (
	"cf-check/profiles"
	"context"
	"fmt"
	"io"
	"log"
	"os"
    "strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF5194"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		distribution := GetDistributions(m.choice)

		for _, dist := range distribution.DistributionList.Items {
            fmt.Println(dist.Status)
		}
		//return quitTextStyle.Render(fmt.Sprintf("%s", GetDistributions()))
	}
	if m.quitting {
		return quitTextStyle.Render("Quit without making a selection.")
	}
	return "\n" + m.list.View()
}

func ProfilesList() tea.Model {
	items := []list.Item{}

	for _, profile := range profiles.GetProfiles() {
		items = append(items, item(profile))
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "AWS Profiles"

	m := model{list: l}

	return m
}

func GetDistributions(profile string) *cloudfront.ListDistributionsOutput {
	// Load config based on a selected profile
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile))

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	client := cloudfront.NewFromConfig(cfg)

	res, err := client.ListDistributions(context.TODO(), &cloudfront.ListDistributionsInput{})

	return res
}

func main() {

	if err := tea.NewProgram(ProfilesList()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

}

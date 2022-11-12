package main

import (
	"cf-check/profiles"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
)

const (
	listHeight    = 14
	listBatchSize = 100
	defaultWidth  = 20
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#FF5194"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string
type itemDelegate struct{}
type model struct {
	profileList   list.Model
	profileChoice string
	distChoice    string
	quitting      bool
}

type Distribution struct {
	// "InProgress": distribution is updating
	// "Deployed": distribution is active and working
	distributionId string
	Domain         string
	Comment        string
	Status         string
}

func (i item) FilterValue() string { return "" }

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

func main() {
	//initialModel := model{list.Model{}, "", "", false}
	if err := tea.NewProgram(ProfilesList()).Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

// Primary update function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "q" || k == "esc" || k == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	if m.profileChoice == "" {
		return updateProfileChoices(msg, m)
	}
	//return updateDistChoices(msg, m)
	var cmd tea.Cmd
	m.profileList, cmd = m.profileList.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.profileChoice != "" {
		GetDistributions(m.profileChoice)
	}
	if m.quitting {
		return quitTextStyle.Render("Quit without making a selection.")
	}
	return "\n" + m.profileList.View()
}

// Sub-update functions

// Update loop for the first view where you're choosing a profile.
func updateProfileChoices(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "enter":
			i, ok := m.profileList.SelectedItem().(item)
			if ok {
				m.profileChoice = string(i)
			}
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.profileList, cmd = m.profileList.Update(msg)
	return m, cmd
}

// Utilities

// Loop through the profiles and create the profileList
func ProfilesList() tea.Model {
	items := []list.Item{}

	for _, profile := range profiles.GetProfiles() {
		items = append(items, item(profile))
	}

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "AWS Profiles"

	m := model{profileList: l}

	return m
}

// Return the list of distributions for the selected profile
func GetDistributions(profile string) ([]*Distribution, error) {
	// Load config based on a selected profile
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile(profile))

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	client := cloudfront.NewFromConfig(cfg)
	res, err := client.ListDistributions(context.TODO(), &cloudfront.ListDistributionsInput{})
	ret := make([]*Distribution, 0, listBatchSize)
	nitems := int(*res.DistributionList.Quantity)

	for i := 0; i < nitems; i++ {
		cfrDist := res.DistributionList.Items[i]
		dist := Distribution{
			Status:         *cfrDist.Status,
			Comment:        *cfrDist.Comment,
			Domain:         *cfrDist.DomainName,
			distributionId: *cfrDist.Id,
		}
		ret = append(ret, &dist)
		fmt.Println(dist)
	}

	return ret, nil
}

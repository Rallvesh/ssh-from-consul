package table

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// NodeData представляет структуру данных из JSON.
type NodeData struct {
	Node            string `json:"Node"`
	TaggedAddresses struct {
		LanIPV4 string `json:"lan_ipv4"`
		WanIPV4 string `json:"wan_ipv4"`
	} `json:"TaggedAddresses"`
}

// Model представляет собой основную модель приложения.
type Model struct {
	table table.Model
}

// Init инициализирует модель. В данном случае не выполняет никаких действий.
func (m Model) Init() tea.Cmd { return nil }

// Update обрабатывает сообщения и обновляет состояние модели.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Selected Node: %s", m.table.SelectedRow()[0]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

// View возвращает строку, которая представляет текущее состояние интерфейса.
func (m Model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

// NewModel создает и настраивает модель с таблицей.
func NewModel(data []NodeData) Model {
	columns := []table.Column{
		{Title: "Node", Width: 20},
		{Title: "LAN IPv4", Width: 15},
		{Title: "WAN IPv4", Width: 15},
	}

	// Преобразуем данные в строки для таблицы.
	rows := make([]table.Row, 0, len(data))
	for _, node := range data {
		rows = append(rows, table.Row{
			node.Node,
			node.TaggedAddresses.LanIPV4,
			node.TaggedAddresses.WanIPV4,
		})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return Model{table: t}
}

// LoadData загружает данные из JSON-файла.
func LoadData(filename string) ([]NodeData, error) {
	filePath, err := filepath.Abs(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var nodes []NodeData
	if err := json.Unmarshal(fileData, &nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nodes, nil
}

// Run запускает приложение.
func Run() {
	nodes, err := LoadData("nodes.json")
	if err != nil {
		fmt.Println("Error loading data:", err)
		os.Exit(1)
	}

	m := NewModel(nodes)
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

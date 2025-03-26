package table

import (
	"encoding/json"
	"fmt"

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
	table        table.Model
	selectedNode string // Поле для хранения выбранного узла.
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
			// Сохраняем выбранный узел.
			m.selectedNode = m.table.SelectedRow()[0]
			return m, tea.Quit // Завершаем программу.
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

// ParseData парсит JSON-строку в структуру данных.
func ParseData(jsonData string) ([]NodeData, error) {
	var nodes []NodeData
	if err := json.Unmarshal([]byte(jsonData), &nodes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return nodes, nil
}

// GetSelectedNode возвращает выбранный узел.
func (m Model) GetSelectedNode() string {
	return m.selectedNode
}

// Run запускает приложение, принимая JSON-строку, и возвращает выбранный узел.
func Run(jsonData string) (string, error) {
	nodes, err := ParseData(jsonData)
	if err != nil {
		return "", fmt.Errorf("error parsing JSON data: %w", err)
	}

	m := NewModel(nodes)
	// Запускаем программу и получаем финальную модель.
	finalModel, err := tea.NewProgram(m).Run()
	if err != nil {
		return "", fmt.Errorf("error running program: %w", err)
	}

	// Приводим finalModel к типу Model и получаем выбранный узел.
	if model, ok := finalModel.(Model); ok {
		return model.GetSelectedNode(), nil
	}

	return "", fmt.Errorf("failed to get selected node")
}

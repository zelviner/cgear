package ui

import (
	"errors"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

// 泛型 item 封装
type genericItem[T any] struct {
	Data  T
	Label string
}

func (i genericItem[T]) Title() string       { return i.Label }
func (i genericItem[T]) Description() string { return "" }
func (i genericItem[T]) FilterValue() string { return i.Label }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i := listItem.(genericItem[any])
	str := fmt.Sprintf("%d. %s", index+1, i.Title())
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(_ ...string) string {
			return selectedItemStyle.Render("> " + str)
		}
	}
	fmt.Fprint(w, fn(str))
}

type model[T any] struct {
	list     list.Model
	choice   *genericItem[T]
	quitting bool
}

func (m model[T]) Init() tea.Cmd {
	return nil
}

func (m model[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			if i, ok := m.list.SelectedItem().(genericItem[any]); ok {
				if data, ok := i.Data.(genericItem[T]); ok {
					m.choice = &data
				}
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model[T]) View() string {
	if m.choice != nil {
		return quitTextStyle.Render(fmt.Sprintf("已选择: %s", m.choice.Title()))
	}
	if m.quitting {
		return quitTextStyle.Render("已取消选择。")
	}
	return "\n" + m.list.View()
}

// ListOption[T] 展示交互式泛型列表，返回所选项、是否取消、错误
func ListOption[T any](title string, options []T, render func(T) string) (T, bool, error) {
	if len(options) == 0 {
		var zero T
		return zero, false, errors.New("无选项可选")
	}

	// 构造 genericItem[any]
	items := make([]list.Item, len(options)+1) // 加一个“退出”
	for i, opt := range options {
		items[i] = genericItem[any]{Data: genericItem[T]{Data: opt, Label: render(opt)}, Label: render(opt)}
	}
	exitLabel := "退出"
	items[len(options)] = genericItem[any]{Data: genericItem[T]{}, Label: exitLabel}

	// 初始化模型
	l := list.New(items, itemDelegate{}, 0, listHeight)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model[T]{list: l}
	prog := tea.NewProgram(m)
	finalModel, err := prog.Run()
	if err != nil {
		var zero T
		return zero, false, err
	}

	result := finalModel.(model[T])
	if result.quitting || result.choice == nil || result.choice.Label == exitLabel {
		var zero T
		return zero, true, nil
	}
	return result.choice.Data, false, nil
}

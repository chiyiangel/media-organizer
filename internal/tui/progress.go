package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ProgressModel struct {
	progress    progress.Model
	total       int
	current     int
	currentFile string
	done        bool
	err         error
}

type progressMsg struct {
	current     int
	total       int
	currentFile string
}

type doneMsg struct{}
type errMsg struct{ err error }

func NewProgress(total int) *ProgressModel {
	p := progress.New(
		progress.WithDefaultGradient(),
		progress.WithWidth(40),
	)
	return &ProgressModel{
		progress: p,
		total:    total,
	}
}

func (m *ProgressModel) Init() tea.Cmd {
	return nil
}

func (m *ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}

	case progressMsg:
		m.current = msg.current
		m.currentFile = msg.currentFile
		if m.current >= m.total {
			m.done = true
			return m, tea.Quit
		}

	case errMsg:
		m.err = msg.err
		return m, tea.Quit

	case doneMsg:
		m.done = true
		return m, tea.Quit
	}

	return m, nil
}

func (m *ProgressModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("\n处理出错: %v\n", m.err)
	}

	if m.done {
		return fmt.Sprintf("\n✨ 处理完成! 共处理 %d 个文件\n", m.total)
	}

	str := strings.Builder{}
	str.WriteString("\n�� 媒体文件整理进度\n\n")

	// 进度条
	percent := float64(m.current) / float64(m.total)
	prog := m.progress.ViewAs(percent)
	str.WriteString(prog)
	str.WriteString("\n\n")

	// 当前文件
	if m.currentFile != "" {
		style := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
		str.WriteString(style.Render(fmt.Sprintf(" 正在处理: %s", m.currentFile)))
		str.WriteString("\n")
	}

	// 进度数字
	str.WriteString(fmt.Sprintf(" %d/%d 文件", m.current, m.total))
	str.WriteString("\n\n")
	str.WriteString(" Ctrl+C 退出\n")

	return str.String()
}

// UpdateProgress 更新进度
func (m *ProgressModel) UpdateProgress(current int, currentFile string) tea.Msg {
	return progressMsg{
		current:     current,
		total:       m.total,
		currentFile: currentFile,
	}
}

// Done 标记完成
func (m *ProgressModel) Done() tea.Msg {
	return doneMsg{}
}

// Error 标记错误
func (m *ProgressModel) Error(err error) tea.Msg {
	return errMsg{err: err}
}

// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package menus

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/khgreav/chronosplit/common"
	"github.com/khgreav/chronosplit/repos"
	"github.com/khgreav/chronosplit/services"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
)

var subjectOptions = []common.MenuItem{
	{ID: "list", Label: "List subjects"},
	{ID: "create", Label: "Create subject"},
	{ID: "delete", Label: "Delete subject"},
	{ID: "back", Label: "Back to main menu"},
	{ID: "exit", Label: "Exit"},
}

var subjectResultOptions = []common.MenuItem{
	{ID: "subjects", Label: "Back to subjects"},
}

type SubjectMenu struct {
	common.BaseMenu
}

func NewSubjectMenu(db *sql.DB) *SubjectMenu {
	return &SubjectMenu{
		BaseMenu: common.BaseMenu{
			Header:  "Manage subjects\n\n",
			Options: subjectOptions,
			Index:   0,
			DB:      db,
		},
	}
}

func (m SubjectMenu) Init() tea.Cmd {
	return nil
}

func (m SubjectMenu) View() tea.View {
	var v tea.View

	var sb strings.Builder
	sb.WriteString(m.Header)

	for i, option := range m.Options {
		cursor := m.GetCursor(i)
		if m.Index == i {
			line := fmt.Sprintf("%s %s", cursor, option.Label)
			fmt.Fprintf(&sb, "%s\n", common.SelectedStyle.Render(line))
		} else {
			fmt.Fprintf(&sb, "%s %s\n", cursor, option.Label)
		}
	}

	sb.WriteString(
		common.RenderHints([]common.MenuHint{
			common.ConfirmHint,
			common.SelectHint,
			common.BackToMainHint,
			common.ExitHint,
		}),
	)

	v.SetContent(sb.String())
	return v
}

func (m *SubjectMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			if m.Index == 0 {
				m.Index = len(m.Options) - 1
			} else {
				m.Index--
			}
		case "down":
			if m.Index == len(m.Options)-1 {
				m.Index = 0
			} else {
				m.Index++
			}
		case "enter":
			switch m.Options[m.Index].ID {
			case "list":
				service := services.NewSubjectService(repos.NewSubjectRepo(m.DB))
				subjects, err := service.ListSubjects()
				resultMenu := NewResultMenu(
					m.DB,
					"Subject operation result",
					subjectResultOptions,
					true,
					"",
				)
				if err != nil {
					resultMenu.Success = false
					resultMenu.Message = err.Error()
				} else {
					t := table.New().
						Border(lipgloss.ASCIIBorder()).
						Headers("ID", "Name")

					for _, s := range subjects {
						t.Row(
							fmt.Sprintf("%d", s.ID),
							s.Name,
						)
					}
					t.StyleFunc(func(row, col int) lipgloss.Style {
						if row == table.HeaderRow {
							return common.HeaderStyle
						}
						return lipgloss.NewStyle()
					})
					resultMenu.Message = t.String()
				}
				return resultMenu, resultMenu.Init()
			case "create":
				subjectCreateMenu := NewSubjectCreateMenu(m.DB)
				return subjectCreateMenu, subjectCreateMenu.Init()
			case "delete":
				service := services.NewSubjectService(repos.NewSubjectRepo(m.DB))
				subjects, err := service.ListSubjects()
				if err != nil {
					resultMenu := NewResultMenu(
						m.DB,
						"Subject operation result",
						subjectResultOptions,
						false,
						err.Error(),
					)
					return resultMenu, resultMenu.Init()
				}
				subjectDeleteMenu := NewSubjectDeleteMenu(
					m.DB,
					subjects,
				)
				return subjectDeleteMenu, subjectDeleteMenu.Init()
			case "back":
				mainMenu := NewMainMenu(m.DB)
				return mainMenu, mainMenu.Init()
			case "exit":
				return m, tea.Quit
			}
		case "ctrl+c":
			mainMenu := NewMainMenu(m.DB)
			return mainMenu, mainMenu.Init()
		case "ctrl+q":
			return m, tea.Quit
		}
	}
	return m, nil
}

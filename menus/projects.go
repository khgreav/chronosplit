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

var projectOptions = []common.MenuItem{
	{ID: "list", Label: "List projects"},
	{ID: "create", Label: "Create project"},
	{ID: "delete", Label: "Delete project"},
	{ID: "back", Label: "Back to main menu"},
	{ID: "exit", Label: "Exit"},
}

var projectResultOptions = []common.MenuItem{
	{ID: "projects", Label: "Back to projects"},
}

type ProjectMenu struct {
	common.BaseMenu
}

func NewProjectMenu(db *sql.DB) *ProjectMenu {
	return &ProjectMenu{
		BaseMenu: common.BaseMenu{
			Header:  "Manage projects\n\n",
			Options: projectOptions,
			Index:   0,
			DB:      db,
		},
	}
}

func (m ProjectMenu) Init() tea.Cmd {
	return nil
}

func (m ProjectMenu) View() tea.View {
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

func (m *ProjectMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.Up()
		case "down":
			m.Down()
		case "enter":
			switch m.Options[m.Index].ID {
			case "list":
				service := services.NewProjectService(repos.NewProjectRepo(m.DB))
				projects, err := service.ListProjects()
				resultMenu := NewResultMenu(
					m.DB,
					"Project operation result",
					projectResultOptions,
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

					for _, p := range projects {
						t.Row(
							fmt.Sprintf("%d", p.ID),
							p.Name,
						)
					}
					t.StyleFunc(func(row, col int) lipgloss.Style {
						if row == table.HeaderRow {
							return common.TableHeaderStyle
						}
						return lipgloss.NewStyle()
					})
					resultMenu.Message = t.String()
				}
				return resultMenu, resultMenu.Init()
			case "create":
				projectCreateMenu := NewProjectCreateMenu(m.DB)
				return projectCreateMenu, projectCreateMenu.Init()
			case "delete":
				service := services.NewProjectService(repos.NewProjectRepo(m.DB))
				projects, err := service.ListProjects()
				if err != nil {
					resultMenu := NewResultMenu(
						m.DB,
						"Project operation result",
						projectResultOptions,
						false,
						err.Error(),
					)
					return resultMenu, resultMenu.Init()
				}
				projectDeleteMenu := NewProjectDeleteMenu(
					m.DB,
					projects,
				)
				return projectDeleteMenu, projectDeleteMenu.Init()
			case "back":
				mainMenu := NewMainMenu(m.DB)
				return mainMenu, mainMenu.Init()
			case "exit":
				return NewQuitMenu(), tea.Quit
			}
		case "ctrl+c":
			mainMenu := NewMainMenu(m.DB)
			return mainMenu, mainMenu.Init()
		case "ctrl+q":
			return NewQuitMenu(), tea.Quit
		}
	}
	return m, nil
}

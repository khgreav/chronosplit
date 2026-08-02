// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package menus

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/khgreav/chronosplit/common"
	"github.com/khgreav/chronosplit/models"
	"github.com/khgreav/chronosplit/repos"
	"github.com/khgreav/chronosplit/services"

	tea "charm.land/bubbletea/v2"
)

type ProjectDeleteMenu struct {
	common.BaseMenu
	Projects []models.Project
}

func NewProjectDeleteMenu(db *sql.DB, projects []models.Project) *ProjectDeleteMenu {
	return &ProjectDeleteMenu{
		BaseMenu: common.BaseMenu{
			Header:  "Delete a project\n\n",
			Options: []common.MenuItem{},
			Index:   0,
			DB:      db,
		},
		Projects: projects,
	}
}

func (m ProjectDeleteMenu) Init() tea.Cmd {
	return nil
}

func (m ProjectDeleteMenu) View() tea.View {
	var v tea.View

	var sb strings.Builder
	sb.WriteString(m.Header)

	for i, option := range m.Projects {
		cursor := m.GetCursor(i)
		if m.Index == i {
			line := fmt.Sprintf("%s %s", cursor, option.Name)
			fmt.Fprintf(&sb, "%s\n", common.SelectedStyle.Render(line))
		} else {
			fmt.Fprintf(&sb, "%s %s\n", cursor, option.Name)
		}
	}

	sb.WriteString("\n\n[Ctrl-C] Back to projects menu")

	v.SetContent(sb.String())
	return v
}

func (m *ProjectDeleteMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			if len(m.Projects) == 1 {
				break
			}
			if m.Index == 0 {
				m.Index = len(m.Projects) - 1
			} else {
				m.Index--
			}
		case "down":
			if len(m.Projects) == 1 {
				break
			}
			if m.Index == len(m.Projects)-1 {
				m.Index = 0
			} else {
				m.Index++
			}
		case "enter":
			id := m.Projects[m.Index].ID
			service := services.NewProjectService(repos.NewProjectRepo(m.DB))
			err := service.DeleteProject(id)
			resultMenu := NewResultMenu(
				m.DB,
				projectResultOptions,
				true,
				"",
			)
			if err != nil {
				resultMenu.Success = false
				resultMenu.Message = err.Error()
			} else {
				resultMenu.Message = fmt.Sprintf(
					"Project %s has been deleted.",
					m.Projects[m.Index].Name,
				)
			}
			return resultMenu, resultMenu.Init()
		case "ctrl+c":
			projectMenu := NewProjectMenu(m.DB)
			return projectMenu, projectMenu.Init()
		}
	}
	return m, nil
}

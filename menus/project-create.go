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

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type ProjectCreateMenu struct {
	common.BaseMenu
	Input textinput.Model
	Retry bool
}

func NewProjectCreateMenu(db *sql.DB) *ProjectCreateMenu {
	input := textinput.New()
	input.Focus()
	return &ProjectCreateMenu{
		BaseMenu: common.BaseMenu{
			Header:  "Create a project",
			Options: []common.MenuItem{},
			Index:   0,
			DB:      db,
		},
		Input: input,
		Retry: false,
	}
}

func (m ProjectCreateMenu) Init() tea.Cmd {
	return nil
}

func (m *ProjectCreateMenu) View() tea.View {
	var v tea.View

	var sb strings.Builder
	sb.WriteString(m.RenderHeader())

	if m.Retry {
		fmt.Fprintf(&sb, "%s\n", common.ErrorStyle.Render("Empty project name is not allowed\n"))
	}
	sb.WriteString("Please enter project name:\n")
	sb.WriteString(m.Input.View())
	sb.WriteString("\n")

	sb.WriteString(
		common.RenderHints([]common.MenuHint{
			common.ConfirmHint,
			common.BackToProjectsHint,
			common.BackToMainHint,
			common.ExitHint,
		}),
	)

	v.SetContent(sb.String())
	return v
}

func (m *ProjectCreateMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			name := m.Input.Value()
			if len(name) == 0 {
				m.Input.Reset()
				m.Retry = true
				return m, nil
			}
			// success
			m.Input.Blur()

			service := services.NewProjectService(repos.NewProjectRepo(m.DB))

			project, err := service.CreateProject(name)

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
				resultMenu.Message = fmt.Sprintf(
					"New project %s created with ID %d.",
					project.Name,
					project.ID,
				)
			}
			return resultMenu, resultMenu.Init()
		case "ctrl+d":
			projectMenu := NewProjectMenu(m.DB)
			return projectMenu, projectMenu.Init()
		case "ctrl+c":
			mainMenu := NewMainMenu(m.DB)
			return mainMenu, mainMenu.Init()
		case "ctrl+q":
			return NewQuitMenu(), tea.Quit
		}
	}
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

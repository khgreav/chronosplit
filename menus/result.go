// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package menus

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/khgreav/chronosplit/common"

	tea "charm.land/bubbletea/v2"
)

var resultOptions = []common.MenuItem{
	{ID: "back", Label: "Back to main menu"},
	{ID: "exit", Label: "Exit"},
}

type ResultMenu struct {
	common.BaseMenu
	Success bool
	Message string
}

func NewResultMenu(db *sql.DB, header string, customOptions []common.MenuItem, success bool, msg string) *ResultMenu {
	return &ResultMenu{
		BaseMenu: common.BaseMenu{
			Header:  header,
			Options: append(customOptions, resultOptions...),
			Index:   0,
			DB:      db,
		},
		Success: success,
		Message: msg,
	}
}

func (m ResultMenu) Init() tea.Cmd {
	return nil
}

func (m ResultMenu) View() tea.View {
	var v tea.View

	var sb strings.Builder
	sb.WriteString(m.Header)
	sb.WriteString("\n\n")
	style := common.GetStyle(m.Success)
	sb.WriteString(style.Render(m.Message))
	sb.WriteString("\n\n")

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

func (m *ResultMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			case "projects":
				projectMenu := NewProjectMenu(m.DB)
				return projectMenu, projectMenu.Init()
			case "subjects":
				subjectMenu := NewSubjectMenu(m.DB)
				return subjectMenu, subjectMenu.Init()
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

// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package menus

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/khgreav/chronosplit/common"
	"github.com/khgreav/chronosplit/models"
	"github.com/khgreav/chronosplit/repos"
	"github.com/khgreav/chronosplit/services"
)

type ReportProjectMenu struct {
	common.BaseMenu
	Projects []models.Project
	From     *time.Time
	To       *time.Time
}

func NewReportProjectMenu(db *sql.DB, projects []models.Project, from, to *time.Time) *ReportProjectMenu {
	return &ReportProjectMenu{
		BaseMenu: common.BaseMenu{
			Header:  "Select a project",
			Options: []common.MenuItem{},
			Index:   0,
			DB:      db,
		},
		Projects: projects,
		From:     from,
		To:       to,
	}
}

func (m ReportProjectMenu) Init() tea.Cmd {
	return nil
}

func (m ReportProjectMenu) View() tea.View {
	var v tea.View

	var sb strings.Builder
	sb.WriteString(m.RenderHeader())

	for i, option := range m.Projects {
		cursor := m.GetCursor(i)
		if m.Index == i {
			line := fmt.Sprintf("%s %s", cursor, option.Name)
			fmt.Fprintf(&sb, "%s\n", common.SelectedStyle.Render(line))
		} else {
			fmt.Fprintf(&sb, "%s %s\n", cursor, option.Name)
		}
	}

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

func (m *ReportProjectMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.Up(len(m.Projects))
		case "down":
			m.Down(len(m.Projects))
		case "enter":
			id := m.Projects[m.Index].ID
			service := services.NewReportService(repos.NewReportRepo(m.DB))
			report, err := service.GetProjectReport(
				id,
				m.From,
				m.To,
			)
			resultMenu := NewResultMenu(
				m.DB,
				"Report operation result",
				[]common.MenuItem{},
				true,
				"",
			)
			if err != nil {
				resultMenu.Success = false
				resultMenu.Message = err.Error()
			} else {
				resultMenu.Message = report.Render()
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
	return m, nil
}

// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package menus

import (
	"database/sql"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/khgreav/chronosplit/common"
	"github.com/khgreav/chronosplit/repos"
	"github.com/khgreav/chronosplit/services"
)

var reportOptions = []common.MenuItem{
	{ID: "all", Label: "All time"},
	{ID: "week", Label: "This week"},
	{ID: "month", Label: "This month"},
	common.BackToMainItem,
	common.ExitItem,
}

type ReportMenu struct {
	common.BaseMenu
}

func NewReportMenu(db *sql.DB) *ReportMenu {
	return &ReportMenu{
		BaseMenu: common.BaseMenu{
			Header:  "Report data",
			Options: reportOptions,
			Index:   0,
			DB:      db,
		},
	}
}

func (m ReportMenu) Init() tea.Cmd {
	return nil
}

func (m ReportMenu) View() tea.View {
	var v tea.View

	var sb strings.Builder
	sb.WriteString(m.RenderHeader())
	sb.WriteString(m.RenderOptions())
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

func (m *ReportMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.Up()
		case "down":
			m.Down()
		case "enter":
			option := m.Options[m.Index].ID
			switch option {
			case "all", "week", "month":
				service := services.NewProjectService(repos.NewProjectRepo(m.DB))
				projects, err := service.ListProjects()
				if err != nil {
					resultMenu := NewResultMenu(
						m.DB,
						"Report operation result",
						[]common.MenuItem{},
						false,
						err.Error(),
					)
					return resultMenu, resultMenu.Init()
				}
				if len(projects) == 0 {
					resultMenu := NewResultMenu(
						m.DB,
						"Report operation result",
						[]common.MenuItem{},
						false,
						"There are no available projects to generate report from.",
					)
					return resultMenu, resultMenu.Init()
				}
				var from, to *time.Time
				switch option {
				case "week":
					f, t := common.GetThisWeek()
					from, to = &f, &t
				case "month":
					f, t := common.GetThisMonth()
					from, to = &f, &t
				}
				reportProjectMenu := NewReportProjectMenu(
					m.DB,
					projects,
					from,
					to,
				)
				return reportProjectMenu, reportProjectMenu.Init()
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

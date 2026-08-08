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
	"github.com/khgreav/chronosplit/version"

	tea "charm.land/bubbletea/v2"
)

var mainMenuOptions = []common.MenuItem{
	{ID: "report", Label: "Reports"},
	{ID: "start", Label: "Start block"},
	{ID: "stop", Label: "Stop block"},
	{ID: "checkpoint", Label: "Create a checkpoint"},
	{ID: "projects", Label: "Manage projects"},
	{ID: "subjects", Label: "Manage subjects"},
	{ID: "exit", Label: "Exit"},
}

type MainMenu struct {
	common.BaseMenu
}

func NewMainMenu(db *sql.DB) *MainMenu {
	return &MainMenu{
		BaseMenu: common.BaseMenu{
			Header:  "Make your life slightly less hellish with a tailored time tracking solution.\n\n",
			Options: mainMenuOptions,
			Index:   0,
			DB:      db,
		},
	}
}

func (m MainMenu) Init() tea.Cmd {
	return nil
}

func (m MainMenu) View() tea.View {
	var v tea.View

	var sb strings.Builder
	sb.WriteString(" ██████╗██╗  ██╗██████╗  ██████╗ ███╗  ██╗ ██████╗ ███████╗██████╗ ██╗     ██╗████████╗\n")
	sb.WriteString("██╔════╝██║  ██║██╔══██╗██╔═══██╗████╗ ██║██╔═══██╗██╔════╝██╔══██╗██║     ██║╚══██╔══╝\n")
	sb.WriteString("██║     ███████║██████╔╝██║   ██║██║██╗██║██║   ██║███████╗██████╔╝██║     ██║   ██║\n")
	sb.WriteString("██║     ██╔══██║██╔══██╗██║   ██║██║╚████║██║   ██║╚════██║██╔═══╝ ██║     ██║   ██║\n")
	sb.WriteString("╚██████╗██║  ██║██║  ██║╚██████╔╝██║ ╚███║╚██████╔╝███████║██║     ███████║██║   ██║\n")
	sb.WriteString(" ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚══╝ ╚═════╝ ╚══════╝╚═╝     ╚══════╝╚═╝   ╚═╝\n")
	sb.WriteString(version.String())
	sb.WriteString("\n\n")
	sb.WriteString(m.Header)
	sb.WriteString(m.RenderOptions())
	sb.WriteString(
		common.RenderHints([]common.MenuHint{
			common.ConfirmHint,
			common.SelectHint,
			common.ExitHint,
		}),
	)

	v.SetContent(sb.String())
	return v
}

func (m *MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			m.Up()
		case "down":
			m.Down()
		case "enter":
			switch m.Options[m.Index].ID {
			case "report":
				reportMenu := NewReportMenu(m.DB)
				return reportMenu, reportMenu.Init()
			case "start":
				service := services.NewBlockService(
					repos.NewBlockRepo(m.DB),
				)
				exists, err := service.ActiveBlockExists()
				resultMenu := NewResultMenu(
					m.DB,
					"Block operation result",
					[]common.MenuItem{},
					true,
					"",
				)
				if err != nil {
					resultMenu.Success = false
					resultMenu.Message = err.Error()
					return resultMenu, resultMenu.Init()
				}
				if exists {
					resultMenu.Success = false
					resultMenu.Message = "An active block already exists."
					return resultMenu, resultMenu.Init()
				}
				id, err := service.StartBlock()
				if err != nil {
					resultMenu.Success = false
					resultMenu.Message = err.Error()
					return resultMenu, resultMenu.Init()
				}
				resultMenu.Message = fmt.Sprintf("New block with ID %d started.", *id)
				return resultMenu, resultMenu.Init()
			case "stop", "checkpoint":
				header := "Checkpoint operation result"
				stopBlock := false
				if m.Options[m.Index].ID == "stop" {
					stopBlock = true
					header = "Block operation result"
				}
				service := services.NewBlockService(
					repos.NewBlockRepo(m.DB),
				)
				exists, err := service.ActiveBlockExists()
				if !exists {
					resultMenu := NewResultMenu(
						m.DB,
						header,
						[]common.MenuItem{},
						false,
						"",
					)
					if err != nil {
						resultMenu.Message = err.Error()
					} else {
						resultMenu.Message = "There is no active block."
					}
					return resultMenu, resultMenu.Init()
				}
				block, err := service.GetActiveBlock()
				if err != nil {
					resultMenu := NewResultMenu(
						m.DB,
						header,
						[]common.MenuItem{},
						false,
						err.Error(),
					)
					return resultMenu, resultMenu.Init()
				}
				projectService := services.NewProjectService(repos.NewProjectRepo(m.DB))
				subjectService := services.NewSubjectService(repos.NewSubjectRepo(m.DB))
				projects, _ := projectService.ListProjects()
				if len(projects) == 0 {
					resultMenu := NewResultMenu(
						m.DB,
						header,
						[]common.MenuItem{},
						false,
						"There are no defined projects.",
					)
					return resultMenu, resultMenu.Init()
				}
				subjects, _ := subjectService.ListSubjects()
				if len(subjects) == 0 {
					resultMenu := NewResultMenu(
						m.DB,
						header,
						[]common.MenuItem{},
						false,
						"There are no defined subjects.",
					)
					return resultMenu, resultMenu.Init()
				}
				checkpointMenu := NewCheckpointMenu(
					m.DB,
					projects,
					subjects,
					block,
					stopBlock,
				)
				return checkpointMenu, checkpointMenu.Init()
			case "projects":
				projectMenu := NewProjectMenu(m.DB)
				return projectMenu, projectMenu.Init()
			case "subjects":
				subjectMenu := NewSubjectMenu(m.DB)
				return subjectMenu, subjectMenu.Init()
			case "exit":
				return NewQuitMenu(), tea.Quit
			}
		case "ctrl+q":
			return NewQuitMenu(), tea.Quit
		}
	}
	return m, nil
}

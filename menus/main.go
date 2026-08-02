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
	{ID: "show", Label: "Show work blocks"},
	{ID: "start", Label: "Start work block"},
	{ID: "stop", Label: "Stop work block"},
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
			Db:      db,
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

	for i, option := range m.Options {
		cursor := m.GetCursor(i)
		if m.Index == i {
			line := fmt.Sprintf("%s %s", cursor, option.Label)
			fmt.Fprintf(&sb, "%s\n", common.SelectedStyle.Render(line))
		} else {
			fmt.Fprintf(&sb, "%s %s\n", cursor, option.Label)
		}
	}

	v.SetContent(sb.String())
	return v
}

func (m *MainMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			case "start":
				service := services.NewBlockService(
					repos.NewBlockRepo(m.Db),
				)
				exists, err := service.ActiveBlockExists()
				resultMenu := NewResultMenu(
					m.Db,
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
					resultMenu.Message = "An active work block already exists."
					return resultMenu, resultMenu.Init()
				}
				id, err := service.StartBlock()
				if err != nil {
					resultMenu.Success = false
					resultMenu.Message = err.Error()
					return resultMenu, resultMenu.Init()
				}
				resultMenu.Message = fmt.Sprintf("New work block with ID %d started.", *id)
				return resultMenu, resultMenu.Init()
			case "stop", "checkpoint":
				stopBlock := false
				if m.Options[m.Index].ID == "stop" {
					stopBlock = true
				}
				service := services.NewBlockService(
					repos.NewBlockRepo(m.Db),
				)
				exists, err := service.ActiveBlockExists()
				if !exists {
					resultMenu := NewResultMenu(
						m.Db,
						[]common.MenuItem{},
						false,
						"",
					)
					if err != nil {
						resultMenu.Message = err.Error()
					} else {
						resultMenu.Message = "There is no active work block."
					}
					return resultMenu, resultMenu.Init()
				}
				block, err := service.GetActiveBlock()
				if err != nil {
					resultMenu := NewResultMenu(
						m.Db,
						[]common.MenuItem{},
						false,
						err.Error(),
					)
					return resultMenu, resultMenu.Init()
				}
				projectService := services.NewProjectService(repos.NewProjectRepo(m.Db))
				subjectService := services.NewSubjectService(repos.NewSubjectRepo(m.Db))
				projects, _ := projectService.ListProjects()
				subjects, _ := subjectService.ListSubjects()
				checkpointMenu := NewCheckpointMenu(
					m.Db,
					projects,
					subjects,
					block,
					stopBlock,
				)
				return checkpointMenu, checkpointMenu.Init()
			case "projects":
				projectMenu := NewProjectMenu(m.Db)
				return projectMenu, projectMenu.Init()
			case "subjects":
				subjectMenu := NewSubjectMenu(m.Db)
				return subjectMenu, subjectMenu.Init()
			case "exit":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

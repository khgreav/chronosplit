package menus

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/khgreav/chronosplit/common"
	"github.com/khgreav/chronosplit/repos"

	tea "charm.land/bubbletea/v2"
)

var mainMenuOptions = []common.MenuItem{
	{ID: "show", Label: "Show work blocks"},
	{ID: "start", Label: "Start work block"},
	{ID: "stop", Label: "Stop work block"},
	{ID: "checkpoint", Label: "Make a checkpoint"},
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
			Header:  "Make your life slightly less hellish with tailored time tracking solution.\n\n",
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
				repo := repos.NewBlockRepo(m.Db)
				id, err := repo.StartBlock()
				if err != nil {
					resultMenu := NewResultMenu(
						m.Db,
						[]common.MenuItem{},
						false,
						err.Error(),
					)
					return resultMenu, resultMenu.Init()
				}
				resultMenu := NewResultMenu(
					m.Db,
					[]common.MenuItem{},
					true,
					fmt.Sprintf("New work block with ID %d started.", *id),
				)
				return resultMenu, resultMenu.Init()
			case "stop":
				blockRepo := repos.NewBlockRepo(m.Db)
				id, err := blockRepo.GetActiveBlockId()
				if err != nil {
					resultMenu := NewResultMenu(
						m.Db,
						[]common.MenuItem{},
						false,
						err.Error(),
					)
					return resultMenu, resultMenu.Init()
				}
				now := time.Now()
				err = blockRepo.StopBlock(*id, now)
				if err != nil {
					resultMenu := NewResultMenu(
						m.Db,
						[]common.MenuItem{},
						false,
						err.Error(),
					)
					return resultMenu, resultMenu.Init()
				}
				resultMenu := NewResultMenu(
					m.Db,
					[]common.MenuItem{},
					true,
					fmt.Sprintf("Work block ID %d stopped.", *id),
				)
				return resultMenu, resultMenu.Init()
			case "projects":
				projectMenu := NewProjectMenu(m.Db)
				return projectMenu, projectMenu.Init()
			case "exit":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

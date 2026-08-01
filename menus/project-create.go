package menus

import (
	"github.com/khgreav/chronosplit/common"
	"github.com/khgreav/chronosplit/repos"
	"github.com/khgreav/chronosplit/services"
	"database/sql"
	"fmt"
	"strings"

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
			Header:  "Creating a project\n\n",
			Options: []common.MenuItem{},
			Index:   0,
			Db:      db,
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

	if m.Retry {
		sb.WriteString("Empty project name is not allowed\n\n")
	}
	sb.WriteString("Please enter project name:\n")
	sb.WriteString(m.Input.View())

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

			repo := repos.NewProjectRepo(m.Db)
			service := services.NewProjectService(repo)

			project, err := service.CreateProject(name)

			resultMenu := NewResultMenu(
				m.Db,
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
		}
	}
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

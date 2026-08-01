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

type SubjectCreateMenu struct {
	common.BaseMenu
	Input textinput.Model
	Retry bool
}

func NewSubjectCreateMenu(db *sql.DB) *SubjectCreateMenu {
	input := textinput.New()
	input.Focus()
	return &SubjectCreateMenu{
		BaseMenu: common.BaseMenu{
			Header:  "Creating a subject\n\n",
			Options: []common.MenuItem{},
			Index:   0,
			Db:      db,
		},
		Input: input,
		Retry: false,
	}
}

func (m SubjectCreateMenu) Init() tea.Cmd {
	return nil
}

func (m *SubjectCreateMenu) View() tea.View {
	var v tea.View

	var sb strings.Builder

	if m.Retry {
		sb.WriteString("Empty subject name is not allowed\n\n")
	}
	sb.WriteString("Please enter subject name:\n")
	sb.WriteString(m.Input.View())

	v.SetContent(sb.String())
	return v
}

func (m *SubjectCreateMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

			repo := repos.NewSubjectRepo(m.Db)
			service := services.NewSubjectService(repo)

			subject, err := service.CreateSubject(name)

			resultMenu := NewResultMenu(
				m.Db,
				subjectResultOptions,
				true,
				"",
			)
			if err != nil {
				resultMenu.Success = false
				resultMenu.Message = err.Error()
			} else {
				resultMenu.Message = fmt.Sprintf(
					"New subject %s created with ID %d.",
					subject.Name,
					subject.ID,
				)
			}
			return resultMenu, resultMenu.Init()
		}
	}
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

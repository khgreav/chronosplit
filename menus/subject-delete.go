package menus

import (
	"github.com/khgreav/chronosplit/common"
	"github.com/khgreav/chronosplit/models"
	"github.com/khgreav/chronosplit/repos"
	"github.com/khgreav/chronosplit/services"
	"database/sql"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type SubjectDeleteMenu struct {
	common.BaseMenu
	Subjects []models.Subject
}

func NewSubjectDeleteMenu(db *sql.DB, subjects []models.Subject) *SubjectDeleteMenu {
	return &SubjectDeleteMenu{
		BaseMenu: common.BaseMenu{
			Header:  "Delete a subject\n\n",
			Options: []common.MenuItem{},
			Index:   0,
			Db:      db,
		},
		Subjects: subjects,
	}
}

func (m SubjectDeleteMenu) Init() tea.Cmd {
	return nil
}

func (m SubjectDeleteMenu) View() tea.View {
	var v tea.View

	var sb strings.Builder
	sb.WriteString(m.Header)

	for i, option := range m.Subjects {
		cursor := m.GetCursor(i)
		if m.Index == i {
			line := fmt.Sprintf("%s %s", cursor, option.Name)
			fmt.Fprintf(&sb, "%s\n", common.SelectedStyle.Render(line))
		} else {
			fmt.Fprintf(&sb, "%s %s\n", cursor, option.Name)
		}
	}

	sb.WriteString("\n\n[Ctrl-C] Back to subjects menu")

	v.SetContent(sb.String())
	return v
}

func (m *SubjectDeleteMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			if len(m.Subjects) == 1 {
				break
			}
			if m.Index == 0 {
				m.Index = len(m.Options) - 1
			} else {
				m.Index--
			}
		case "down":
			if len(m.Subjects) == 1 {
				break
			}
			if m.Index == len(m.Options)-1 {
				m.Index = 0
			} else {
				m.Index++
			}
		case "enter":
			id := m.Subjects[m.Index].ID
			repo := repos.NewSubjectRepo(m.Db)
			service := services.NewSubjectService(repo)
			err := service.DeleteSubject(id)
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
					"Subject %s has been deleted.",
					m.Subjects[m.Index].Name,
				)
			}
			return resultMenu, resultMenu.Init()
		case "ctrl+c":
			subjectMenu := NewSubjectMenu(m.Db)
			return subjectMenu, subjectMenu.Init()
		}
	}
	return m, nil
}

// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package menus

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"github.com/khgreav/chronosplit/common"
	"github.com/khgreav/chronosplit/models"
	"github.com/khgreav/chronosplit/repos"
	"github.com/khgreav/chronosplit/services"
)

type CheckpointState int

const (
	ProjectState CheckpointState = iota
	SubjectState
	DescriptionState
)

type CheckpointMenu struct {
	common.BaseMenu
	CurrentState CheckpointState
	Projects     []models.Project
	Subjects     []models.Subject
	Block        *models.Block
	ProjectID    int64
	SubjectID    int64
	StopBlock    bool
	Input        textinput.Model
	Retry        bool
}

func NewCheckpointMenu(
	db *sql.DB,
	projects []models.Project,
	subjects []models.Subject,
	block *models.Block,
	stopBlock bool,
) *CheckpointMenu {
	input := textinput.New()
	return &CheckpointMenu{
		BaseMenu: common.BaseMenu{
			Header:  "Create a checkpoint",
			Options: []common.MenuItem{},
			Index:   0,
			DB:      db,
		},
		CurrentState: ProjectState,
		Projects:     projects,
		Subjects:     subjects,
		Block:        block,
		ProjectID:    -1,
		SubjectID:    -1,
		StopBlock:    stopBlock,
		Input:        input,
		Retry:        false,
	}
}

func (m CheckpointMenu) Init() tea.Cmd {
	return nil
}

func (m *CheckpointMenu) View() tea.View {
	var v tea.View

	var sb strings.Builder
	sb.WriteString(m.RenderHeader())

	switch m.CurrentState {
	case ProjectState:
		sb.WriteString("Select a project\n\n")
		for i, project := range m.Projects {
			cursor := m.GetCursor(i)
			if m.Index == i {
				line := fmt.Sprintf("%s %s", cursor, project.Name)
				fmt.Fprintf(&sb, "%s\n", common.SelectedStyle.Render(line))
			} else {
				fmt.Fprintf(&sb, "%s %s\n", cursor, project.Name)
			}
		}
	case SubjectState:
		sb.WriteString("Select a subject\n\n")
		for i, subject := range m.Subjects {
			cursor := m.GetCursor(i)
			if m.Index == i {
				line := fmt.Sprintf("%s %s", cursor, subject.Name)
				fmt.Fprintf(&sb, "%s\n", common.SelectedStyle.Render(line))
			} else {
				fmt.Fprintf(&sb, "%s %s\n", cursor, subject.Name)
			}
		}
	case DescriptionState:
		if m.Retry {
			fmt.Fprintf(&sb, "%s\n", common.ErrorStyle.Render("Empty description is not allowed\n"))
		}
		sb.WriteString("Enter a description:\n")
		sb.WriteString(m.Input.View())
		sb.WriteString("\n")
	}

	sb.WriteString(
		common.RenderHints([]common.MenuHint{
			common.ConfirmHint,
			common.BackToMainHint,
			common.ExitHint,
		}),
	)

	v.SetContent(sb.String())
	return v
}

func (m *CheckpointMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "up":
			if m.CurrentState != DescriptionState {
				optLen := 0
				if m.CurrentState == ProjectState {
					optLen = len(m.Projects)
				} else {
					optLen = len(m.Subjects)
				}
				m.Up(optLen)
			}
		case "down":
			if m.CurrentState != DescriptionState {
				optLen := 0
				if m.CurrentState == ProjectState {
					optLen = len(m.Projects)
				} else {
					optLen = len(m.Subjects)
				}
				m.Down(optLen)
			}
		case "enter":
			switch m.CurrentState {
			case ProjectState:
				m.ProjectID = m.Projects[m.Index].ID
				m.CurrentState = SubjectState
			case SubjectState:
				m.SubjectID = m.Subjects[m.Index].ID
				m.CurrentState = DescriptionState
				m.Input.Focus()
			case DescriptionState:
				desc := m.Input.Value()
				if len(desc) == 0 {
					m.Input.Reset()
					m.Retry = true
					return m, nil
				}
				m.Input.Blur()

				resultMenu := NewResultMenu(
					m.DB,
					"Checkpoint operation result",
					[]common.MenuItem{},
					true,
					"",
				)

				service := services.NewCheckpointService(repos.NewCheckpointRepo(m.DB))
				now := time.Now().UTC()
				var previousCheckpointID *int64

				tx, err := m.DB.Begin()
				if err != nil {
					resultMenu.Success = false
					resultMenu.Message = fmt.Sprintf("Failed to start transaction: %v", err)
					_ = tx.Rollback()
					return resultMenu, resultMenu.Init()
				}

				// GET POSSIBLE LAST CHECKPOINT
				lastCheckpoint, err := service.GetLastCheckpoint(m.Block.ID)
				if err != nil {
					resultMenu.Success = false
					resultMenu.Message = err.Error()
					_ = tx.Rollback()
					return resultMenu, resultMenu.Init()
				}
				// IF PREVIOUS CHECKPOINT EXISTS, END CONCLUDE IT
				if lastCheckpoint != nil {
					previousCheckpointID = &lastCheckpoint.ID
					err := service.UpdateCheckpoint(lastCheckpoint)
					if err != nil {
						resultMenu.Success = false
						resultMenu.Message = err.Error()
						_ = tx.Rollback()
						return resultMenu, resultMenu.Init()
					}
				}

				if m.StopBlock {
					blockService := services.NewBlockService(repos.NewBlockRepo(m.DB))
					err := blockService.StopBlock(m.Block.ID, now)
					if err != nil {
						resultMenu.Success = false
						resultMenu.Message = err.Error()
						_ = tx.Rollback()
						return resultMenu, resultMenu.Init()
					}
				}

				newCheckpoint := &models.Checkpoint{
					ID:                   0,
					BlockID:              m.Block.ID,
					ProjectID:            m.ProjectID,
					SubjectID:            m.SubjectID,
					PreviousCheckpointID: previousCheckpointID,
					StartTime:            m.Block.CreatedAt,
					Description:          desc,
				}
				if lastCheckpoint != nil {
					newCheckpoint.StartTime = lastCheckpoint.EndTime
				}
				newCheckpoint.EndTime = now
				newCheckpoint, err = service.CreateCheckpoint(newCheckpoint)

				if err != nil {
					resultMenu.Success = false
					resultMenu.Message = err.Error()
					_ = tx.Rollback()
				} else {
					err := tx.Commit()
					if err != nil {
						resultMenu.Success = false
						resultMenu.Message = fmt.Sprintf("Failed to commit transaction: %v", err)
						_ = tx.Rollback()
						return resultMenu, resultMenu.Init()
					}
					line := fmt.Sprintf(
						"Created new checkpoint ID %d for block ID %d.\n",
						newCheckpoint.ID,
						m.Block.ID,
					)
					if m.StopBlock {
						line += fmt.Sprintf(
							"Work block ID %d stopped.\n",
							m.Block.ID,
						)
					}
					resultMenu.Message = line
				}
				return resultMenu, resultMenu.Init()
			}
		case "ctrl+c":
			mainMenu := NewMainMenu(m.DB)
			return mainMenu, mainMenu.Init()
		}
	}
	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package common

import (
	"database/sql"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type Menu interface {
	Init() tea.Cmd
	Update(tea.Msg) (tea.Model, tea.Cmd)
	View() tea.View
}

type MenuItem struct {
	ID    string
	Label string
}

type MenuHint struct {
	Key   string
	Label string
}

func (h MenuHint) String() string {
	return fmt.Sprintf("[%s] %s", h.Key, h.Label)
}

type BaseMenu struct {
	Header  string
	Options []MenuItem
	Index   int
	DB      *sql.DB
}

func (m BaseMenu) RenderHeader() string {
	return fmt.Sprintf("%s\n\n", HeaderStyle.Render(m.Header))
}

func (m BaseMenu) RenderOptions() string {
	var sb strings.Builder

	for i, option := range m.Options {
		cursor := m.GetCursor(i)
		if m.Index == i {
			line := fmt.Sprintf("%s %s", cursor, option.Label)
			fmt.Fprintf(&sb, "%s\n", SelectedStyle.Render(line))
		} else {
			fmt.Fprintf(&sb, "%s %s\n", cursor, option.Label)
		}
	}
	return sb.String()
}

func (m BaseMenu) GetCursor(i int) string {
	if m.Index == i {
		return ">"
	}
	return " "
}

func (m *BaseMenu) Up(itemCount ...int) {
	var cnt int
	if len(itemCount) == 0 {
		cnt = len(m.Options)
	} else {
		cnt = itemCount[0]
	}
	if cnt == 1 {
		return
	}
	if m.Index == 0 {
		m.Index = cnt - 1
	} else {
		m.Index--
	}
}

func (m *BaseMenu) Down(itemCount ...int) {
	var cnt int
	if len(itemCount) == 0 {
		cnt = len(m.Options)
	} else {
		cnt = itemCount[0]
	}
	if cnt == 1 {
		return
	}
	if m.Index == cnt-1 {
		m.Index = 0
	} else {
		m.Index++
	}
}

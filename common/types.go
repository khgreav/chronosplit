// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package common

import (
	"database/sql"
	"fmt"

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

func (m BaseMenu) GetCursor(i int) string {
	if m.Index == i {
		return ">"
	}
	return " "
}

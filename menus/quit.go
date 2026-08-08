// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package menus

import tea "charm.land/bubbletea/v2"

type QuitMenu struct{}

func NewQuitMenu() *QuitMenu {
	return &QuitMenu{}
}

func (m QuitMenu) Init() tea.Cmd {
	return nil
}

func (m QuitMenu) View() tea.View {
	return tea.View{}
}

func (m QuitMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

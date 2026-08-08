// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package common

import "charm.land/lipgloss/v2"

var (
	HeaderStyle      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f0de1f"))
	TableHeaderStyle = lipgloss.NewStyle().Bold(true).Align(lipgloss.Center)
	SelectedStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#f0de1f"))
	SuccessStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#20e618"))
	ErrorStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("#e61818"))
)

func GetStyle(success bool) lipgloss.Style {
	if success {
		return SuccessStyle
	}
	return ErrorStyle
}

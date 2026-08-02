package common

import "strings"

var (
	ConfirmHint        = MenuHint{Key: "enter", Label: "Confirm"}
	SelectHint         = MenuHint{Key: "up/down", Label: "Select option"}
	BackToProjectsHint = MenuHint{Key: "ctrl-d", Label: "Back to projects menu"}
	BackToSubjectsHint = MenuHint{Key: "ctrl-d", Label: "Back to subjects menu"}
	BackToMainHint     = MenuHint{Key: "ctrl-c", Label: "Back to main menu"}
	ExitHint           = MenuHint{Key: "ctrl-q", Label: "Quit"}
)

func RenderHints(hints []MenuHint) string {
	var sb strings.Builder
	sb.WriteString("\n")
	for i, hint := range hints {
		sb.WriteString(hint.String())
		if i != len(hints)-1 {
			sb.WriteString(" - ")
		}
	}
	return sb.String()
}

// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package common_test

import (
	"testing"

	"github.com/khgreav/chronosplit/common"
)

func TestRenderHints(t *testing.T) {
	tests := []struct {
		name     string
		hints    []common.MenuHint
		expected string
	}{
		{
			name: "single hint",
			hints: []common.MenuHint{
				{Key: "enter", Label: "Confirm"},
			},
			expected: "\n[enter] Confirm",
		},
		{
			name: "multiple hints",
			hints: []common.MenuHint{
				{Key: "enter", Label: "Confirm"},
				{Key: "ctrl-c", Label: "Back"},
			},
			expected: "\n[enter] Confirm - [ctrl-c] Back",
		},
		{
			name:     "no hints",
			hints:    []common.MenuHint{},
			expected: "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := common.RenderHints(tt.hints); got != tt.expected {
				t.Errorf("RenderHints() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGlobalHints(t *testing.T) {
	if common.ConfirmHint.Key != "enter" || common.ConfirmHint.Label != "Confirm" {
		t.Errorf("ConfirmHint mismatch: got %v", common.ConfirmHint)
	}
	if common.SelectHint.Key != "up/down" || common.SelectHint.Label != "Select option" {
		t.Errorf("SelectHint mismatch: got %v", common.SelectHint)
	}
	if common.BackToProjectsHint.Key != "ctrl-d" || common.BackToProjectsHint.Label != "Back to projects menu" {
		t.Errorf("BackToProjectsHint mismatch: got %v", common.BackToProjectsHint)
	}
	if common.BackToSubjectsHint.Key != "ctrl-d" || common.BackToSubjectsHint.Label != "Back to subjects menu" {
		t.Errorf("BackToSubjectsHint mismatch: got %v", common.BackToSubjectsHint)
	}
	if common.BackToMainHint.Key != "ctrl-c" || common.BackToMainHint.Label != "Back to main menu" {
		t.Errorf("BackToMainHint mismatch: got %v", common.BackToMainHint)
	}
	if common.ExitHint.Key != "ctrl-q" || common.ExitHint.Label != "Quit" {
		t.Errorf("ExitHint mismatch: got %v", common.ExitHint)
	}
}

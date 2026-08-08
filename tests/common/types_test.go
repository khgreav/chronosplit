// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package common_test

import (
	"testing"

	"github.com/khgreav/chronosplit/common"
)

func TestMenuHint_String(t *testing.T) {
	tests := []struct {
		name     string
		hint     common.MenuHint
		expected string
	}{
		{
			name:     "standard hint",
			hint:     common.MenuHint{Key: "k", Label: "key"},
			expected: "[k] key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hint.String(); got != tt.expected {
				t.Errorf("MenuHint.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBaseMenu_GetCursor(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		menuIdx  int
		expected string
	}{
		{
			name:     "cursor at index",
			index:    1,
			menuIdx:  1,
			expected: ">",
		},
		{
			name:     "cursor not at index",
			index:    0,
			menuIdx:  1,
			expected: " ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := common.BaseMenu{Index: tt.menuIdx}
			if got := menu.GetCursor(tt.index); got != tt.expected {
				t.Errorf("BaseMenu.GetCursor() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestBaseMenu_Up(t *testing.T) {
	tests := []struct {
		name          string
		initialIndex  int
		itemCount     int
		optionsCount  int
		expectedIndex int
	}{
		{
			name:          "move up from middle",
			initialIndex:  1,
			itemCount:     3,
			optionsCount:  3,
			expectedIndex: 0,
		},
		{
			name:          "wrap around to end from zero",
			initialIndex:  0,
			itemCount:     3,
			optionsCount:  3,
			expectedIndex: 2,
		},
		{
			name:          "no change if itemCount is 1",
			initialIndex:  0,
			itemCount:     1,
			optionsCount:  1,
			expectedIndex: 0,
		},
		{
			name:          "use default itemCount from Options",
			initialIndex:  1,
			itemCount:     0,
			optionsCount:  3,
			expectedIndex: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := common.BaseMenu{
				Index:   tt.initialIndex,
				Options: make([]common.MenuItem, tt.optionsCount),
			}

			if tt.itemCount == 0 {
				menu.Up()
			} else {
				menu.Up(tt.itemCount)
			}

			if menu.Index != tt.expectedIndex {
				t.Errorf("BaseMenu.Up() index = %v, want %v", menu.Index, tt.expectedIndex)
			}
		})
	}
}

func TestBaseMenu_Down(t *testing.T) {
	tests := []struct {
		name          string
		initialIndex  int
		itemCount     int
		optionsCount  int
		expectedIndex int
	}{
		{
			name:          "move down from middle",
			initialIndex:  1,
			itemCount:     3,
			optionsCount:  3,
			expectedIndex: 2,
		},
		{
			name:          "wrap around to start from end",
			initialIndex:  2,
			itemCount:     3,
			optionsCount:  3,
			expectedIndex: 0,
		},
		{
			name:          "no change if itemCount is 1",
			initialIndex:  0,
			itemCount:     1,
			optionsCount:  1,
			expectedIndex: 0,
		},
		{
			name:          "use default itemCount from Options",
			initialIndex:  1,
			itemCount:     0,
			optionsCount:  3,
			expectedIndex: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			menu := common.BaseMenu{
				Index:   tt.initialIndex,
				Options: make([]common.MenuItem, tt.optionsCount),
			}

			if tt.itemCount == 0 {
				menu.Down()
			} else {
				menu.Down(tt.itemCount)
			}

			if menu.Index != tt.expectedIndex {
				t.Errorf("BaseMenu.Down() index = %v, want %v", menu.Index, tt.expectedIndex)
			}
		})
	}
}

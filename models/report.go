// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package models

import (
	"fmt"
	"strings"
	"time"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/table"
	"github.com/khgreav/chronosplit/common"
)

type ProjectReport struct {
	Name          string
	TotalDuration time.Duration
	Subjects      []SubjectReport
}

type SubjectReport struct {
	Name        string
	Checkpoints []CheckpointReport
}

type CheckpointReport struct {
	Description string
	From        time.Time
	To          time.Time
}

func (pr ProjectReport) Render() string {
	var sb strings.Builder

	fmt.Fprintf(&sb, "Project name: %s\n", pr.Name)
	fmt.Fprintf(&sb, "Total time: %s\n", common.FormatDuration(pr.TotalDuration))

	for _, subject := range pr.Subjects {
		fmt.Fprintf(&sb, "\nSubject name: %s\n", subject.Name)
		t := table.New().
			Border(lipgloss.ASCIIBorder()).
			Headers("Start", "End", "Activity")

		for _, c := range subject.Checkpoints {
			t.Row(
				c.From.Format(time.RFC3339),
				c.To.Format(time.RFC3339),
				c.Description,
			)
		}

		t.StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow || (row == table.HeaderRow && col == 0) {
				return common.TableHeaderStyle
			}
			return lipgloss.NewStyle()
		})

		sb.WriteString(t.String())
	}

	return sb.String()
}

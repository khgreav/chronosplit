// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package repos

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/khgreav/chronosplit/models"
)

type ReportRepo struct {
	db *sql.DB
}

func NewReportRepo(db *sql.DB) *ReportRepo {
	return &ReportRepo{db: db}
}

func (r *ReportRepo) GetProjectReport(projectID int64, start, end *time.Time) (*models.ProjectReport, error) {
	query := `
		SELECT
			p.name AS project_name,
			s.name AS subject_name,
			c.description,
			c.start_time,
			c.end_time
		FROM projects p
		JOIN checkpoints c ON p.id = c.project_id
		JOIN subjects s ON s.id = c.subject_id
		WHERE p.id = ?
	`
	args := []interface{}{projectID}

	if start != nil {
		query += " AND c.start_time >= ?"
		args = append(args, start.UTC())
	}
	if end != nil {
		query += " AND c.end_time <= ?"
		args = append(args, end.UTC())
	}

	query += " ORDER BY s.name, c.start_time;"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("Failed to query project report: %w", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	var report *models.ProjectReport
	subjectMap := make(map[string]*models.SubjectReport)
	var subjectNames []string // To maintain order of subjects

	for rows.Next() {
		var projectName, subjectName, description string
		var start, end time.Time

		err := rows.Scan(
			&projectName,
			&subjectName,
			&description,
			&start,
			&end,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan report row: %w", err)
		}

		if report == nil {
			report = &models.ProjectReport{
				Name: projectName,
			}
		}

		report.TotalDuration += end.Sub(start)

		if _, ok := subjectMap[subjectName]; !ok {
			newSubject := &models.SubjectReport{
				Name:        subjectName,
				Checkpoints: []models.CheckpointReport{},
			}
			subjectMap[subjectName] = newSubject
			subjectNames = append(subjectNames, subjectName)
		}
		subjectMap[subjectName].Checkpoints = append(
			subjectMap[subjectName].Checkpoints,
			models.CheckpointReport{
				Description: description,
				From:        start,
				To:          end,
			},
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error occurred while iterating rows: %w", err)
	}

	if report == nil {
		return nil, fmt.Errorf("No records found for specified time frame.")
	}

	// Convert map to ordered slice
	for _, name := range subjectNames {
		report.Subjects = append(report.Subjects, *subjectMap[name])
	}

	return report, nil
}

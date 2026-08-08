// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package services

import (
	"time"

	"github.com/khgreav/chronosplit/models"
	"github.com/khgreav/chronosplit/repos"
)

type ReportService struct {
	repo *repos.ReportRepo
}

func NewReportService(repo *repos.ReportRepo) *ReportService {
	return &ReportService{repo: repo}
}

func (s *ReportService) GetProjectReport(projectID int64, start, end *time.Time) (*models.ProjectReport, error) {
	return s.repo.GetProjectReport(projectID, start, end)
}

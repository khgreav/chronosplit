// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package services

import (
	"github.com/khgreav/chronosplit/models"
	"github.com/khgreav/chronosplit/repos"
)

type SubjectService struct {
	repo *repos.SubjectRepo
}

func NewSubjectService(repo *repos.SubjectRepo) *SubjectService {
	return &SubjectService{repo: repo}
}

func (s *SubjectService) ListSubjects() ([]models.Subject, error) {
	return s.repo.List()
}

func (s *SubjectService) GetSubject(id int64) (*models.Subject, error) {
	return s.repo.Get(id)
}

func (s *SubjectService) CreateSubject(name string) (*models.Subject, error) {
	return s.repo.Create(name)
}

func (s *SubjectService) DeleteSubject(id int64) error {
	return s.repo.Delete(id)
}

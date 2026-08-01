package services

import (
	"chronosplit/models"
	"chronosplit/repos"
)

type ProjectService struct {
	repo *repos.ProjectRepo
}

func NewProjectService(repo *repos.ProjectRepo) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) ListProjects() ([]models.Project, error) {
	return s.repo.List()
}

func (s *ProjectService) GetProject(id int64) (*models.Project, error) {
	return s.repo.Get(id)
}

func (s *ProjectService) CreateProject(name string) (*models.Project, error) {
	return s.repo.Create(name)
}

func (s *ProjectService) DeleteProject(id int64) error {
	return s.repo.Delete(id)
}

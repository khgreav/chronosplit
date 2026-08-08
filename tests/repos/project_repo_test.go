// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package repos_test

import (
	"database/sql"
	"testing"

	"github.com/khgreav/chronosplit/repos"
)

func setupProjectDB(t *testing.T) *sql.DB {
	return newTestDB(t, "")
}

func TestProjectRepo_Create(t *testing.T) {
	db := setupProjectDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewProjectRepo(db)
	name := "Test Project"

	project, err := repo.Create(name)
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	if project.Name != name {
		t.Errorf("expected name %q, got %q", name, project.Name)
	}

	if project.ID <= 0 {
		t.Errorf("expected non-zero ID, got %d", project.ID)
	}
}

func TestProjectRepo_Get(t *testing.T) {
	db := setupProjectDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewProjectRepo(db)
	name := "Get Project"
	created, err := repo.Create(name)
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	project, err := repo.Get(created.ID)
	if err != nil {
		t.Fatalf("failed to get project: %v", err)
	}

	if project.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, project.ID)
	}

	if project.Name != name {
		t.Errorf("expected name %q, got %q", name, project.Name)
	}
}

func TestProjectRepo_List(t *testing.T) {
	db := setupProjectDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewProjectRepo(db)
	names := []string{"P1", "P2", "P3"}

	for _, name := range names {
		_, err := repo.Create(name)
		if err != nil {
			t.Fatalf("failed to create project %s: %v", name, err)
		}
	}

	projects, err := repo.List()
	if err != nil {
		t.Fatalf("failed to list projects: %v", err)
	}

	if len(projects) != len(names) {
		t.Errorf("expected %d projects, got %d", len(names), len(projects))
	}

	for i, name := range names {
		if projects[i].Name != name {
			t.Errorf("expected project %d name %q, got %q", i, name, projects[i].Name)
		}
	}
}

func TestProjectRepo_Delete(t *testing.T) {
	db := setupProjectDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewProjectRepo(db)
	created, err := repo.Create("Delete Me")
	if err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	err = repo.Delete(created.ID)
	if err != nil {
		t.Fatalf("failed to delete project: %v", err)
	}

	_, err = repo.Get(created.ID)
	if err == nil {
		t.Error("expected error when getting deleted project, got nil")
	}

	err = repo.Delete(999)
	if err == nil {
		t.Error("expected error when deleting non-existent project, got nil")
	}
}

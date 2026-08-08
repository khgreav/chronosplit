// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package repos_test

import (
	"database/sql"
	"testing"

	"github.com/khgreav/chronosplit/repos"

	_ "modernc.org/sqlite"
)

func setupSubjectDB(t *testing.T) *sql.DB {
	return newTestDB(t, "")
}

func TestSubjectRepo_Create(t *testing.T) {
	db := setupSubjectDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewSubjectRepo(db)
	name := "Test Subject"

	subject, err := repo.Create(name)
	if err != nil {
		t.Fatalf("failed to create subject: %v", err)
	}

	if subject.Name != name {
		t.Errorf("expected name %q, got %q", name, subject.Name)
	}

	if subject.ID <= 0 {
		t.Errorf("expected non-zero ID, got %d", subject.ID)
	}
}

func TestSubjectRepo_Get(t *testing.T) {
	db := setupSubjectDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewSubjectRepo(db)
	name := "Get Subject"
	created, err := repo.Create(name)
	if err != nil {
		t.Fatalf("failed to create subject: %v", err)
	}

	subject, err := repo.Get(created.ID)
	if err != nil {
		t.Fatalf("failed to get subject: %v", err)
	}

	if subject.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, subject.ID)
	}

	if subject.Name != name {
		t.Errorf("expected name %q, got %q", name, subject.Name)
	}
}

func TestSubjectRepo_List(t *testing.T) {
	db := setupSubjectDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewSubjectRepo(db)
	names := []string{"S1", "S2", "S3"}

	for _, name := range names {
		_, err := repo.Create(name)
		if err != nil {
			t.Fatalf("failed to create subject %s: %v", name, err)
		}
	}

	subjects, err := repo.List()
	if err != nil {
		t.Fatalf("failed to list subjects: %v", err)
	}

	if len(subjects) != len(names) {
		t.Errorf("expected %d subjects, got %d", len(names), len(subjects))
	}

	for i, name := range names {
		if subjects[i].Name != name {
			t.Errorf("expected subject %d name %q, got %q", i, name, subjects[i].Name)
		}
	}
}

func TestSubjectRepo_Delete(t *testing.T) {
	db := setupSubjectDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewSubjectRepo(db)
	created, err := repo.Create("Delete Me")
	if err != nil {
		t.Fatalf("failed to create subject: %v", err)
	}

	err = repo.Delete(created.ID)
	if err != nil {
		t.Fatalf("failed to delete subject: %v", err)
	}

	_, err = repo.Get(created.ID)
	if err == nil {
		t.Error("expected error when getting deleted subject, got nil")
	}

	err = repo.Delete(999)
	if err == nil {
		t.Error("expected error when deleting non-existent subject, got nil")
	}
}

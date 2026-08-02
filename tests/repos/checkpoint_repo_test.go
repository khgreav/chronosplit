// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/khgreav/chronosplit/models"
	"github.com/khgreav/chronosplit/repos"
)

func setupCheckpointDB(t *testing.T) *sql.DB {
	return newTestDB(t, "")
}


func TestCheckpointRepo_Create(t *testing.T) {
	db := setupCheckpointDB(t)
	defer db.Close()

	repo := repos.NewCheckpointRepo(db)
	now := time.Now().UTC()
	c := &models.Checkpoint{
		BlockID:       1,
		ProjectID:     1,
		SubjectID:     1,
		StartTime:     now,
		Description:   "Test Checkpoint",
	}

	created, err := repo.Create(c)
	if err != nil {
		t.Fatalf("failed to create checkpoint: %v", err)
	}

	if created.ID == 0 {
		t.Error("expected non-zero ID")
	}

	if created.Description != "Test Checkpoint" {
		t.Errorf("expected description 'Test Checkpoint', got %q", created.Description)
	}
}

func TestCheckpointRepo_Get(t *testing.T) {
	db := setupCheckpointDB(t)
	defer db.Close()

	repo := repos.NewCheckpointRepo(db)
	now := time.Now().UTC()
	c := &models.Checkpoint{
		Description: "Get Me",
		StartTime:   now,
	}
	created, _ := repo.Create(c)

	got, err := repo.Get(created.ID)
	if err != nil {
		t.Fatalf("failed to get checkpoint: %v", err)
	}

	if got.Description != "Get Me" {
		t.Errorf("expected description 'Get Me', got %q", got.Description)
	}
}

func TestCheckpointRepo_Update(t *testing.T) {
	db := setupCheckpointDB(t)
	defer db.Close()

	repo := repos.NewCheckpointRepo(db)
	c := &models.Checkpoint{Description: "Old"}
	created, _ := repo.Create(c)

	created.Description = "New"
	err := repo.Update(created)
	if err != nil {
		t.Fatalf("failed to update checkpoint: %v", err)
	}

	updated, _ := repo.Get(created.ID)
	if updated.Description != "New" {
		t.Errorf("expected description 'New', got %q", updated.Description)
	}
}

func TestCheckpointRepo_Delete(t *testing.T) {
	db := setupCheckpointDB(t)
	defer db.Close()

	repo := repos.NewCheckpointRepo(db)
	c := &models.Checkpoint{Description: "To Delete"}
	created, _ := repo.Create(c)

	err := repo.Delete(created.ID)
	if err != nil {
		t.Fatalf("failed to delete checkpoint: %v", err)
	}

	_, err = repo.Get(created.ID)
	if err == nil {
		t.Error("expected error getting deleted checkpoint, got nil")
	}
}

func TestCheckpointRepo_ListByBlock(t *testing.T) {
	db := setupCheckpointDB(t)
	defer db.Close()

	repo := repos.NewCheckpointRepo(db)
	c1 := &models.Checkpoint{BlockID: 1, Description: "B1"}
	c2 := &models.Checkpoint{BlockID: 1, Description: "B2"}
	c3 := &models.Checkpoint{BlockID: 2, Description: "B3"}

	repo.Create(c1)
	repo.Create(c2)
	repo.Create(c3)

	checkpoints, err := repo.ListByBlock(1)
	if err != nil {
		t.Fatalf("failed to list by block: %v", err)
	}

	if len(checkpoints) != 2 {
		t.Errorf("expected 2 checkpoints, got %d", len(checkpoints))
	}
}

func TestCheckpointRepo_HasCheckpoints(t *testing.T) {
	db := setupCheckpointDB(t)
	defer db.Close()

	repo := repos.NewCheckpointRepo(db)
	
	exists, _ := repo.HasCheckpoints(1)
	if exists {
		t.Error("expected no checkpoints for block 1")
	}

	repo.Create(&models.Checkpoint{BlockID: 1})

	exists, _ = repo.HasCheckpoints(1)
	if !exists {
		t.Error("expected checkpoints for block 1")
	}
}

func TestCheckpointRepo_GetLastCheckpoint(t *testing.T) {
	db := setupCheckpointDB(t)
	defer db.Close()

	repo := repos.NewCheckpointRepo(db)
	
	// No checkpoints
	last, err := repo.GetLastCheckpoint(1)
	if err != nil {
		t.Fatalf("failed to get last checkpoint: %v", err)
	}
	if last != nil {
		t.Error("expected nil for no checkpoints")
	}

	// One active checkpoint
	repo.Create(&models.Checkpoint{BlockID: 1, Description: "Active"})

	last, err = repo.GetLastCheckpoint(1)
	if err != nil {
		t.Fatalf("failed to get last checkpoint: %v", err)
	}
	if last == nil || last.Description != "Active" {
		t.Errorf("expected 'Active' checkpoint, got %v", last)
	}

	// One finished checkpoint
	now := time.Now().UTC()
	repo.Create(&models.Checkpoint{BlockID: 1, Description: "Finished", EndTime: &now})

	last, err = repo.GetLastCheckpoint(1)
	if err != nil {
		t.Fatalf("failed to get last checkpoint: %v", err)
	}
	if last != nil && last.Description == "Finished" {
		t.Error("expected nil for finished checkpoint")
	}
}

func TestCheckpointRepo_List(t *testing.T) {
	db := setupCheckpointDB(t)
	defer db.Close()

	repo := repos.NewCheckpointRepo(db)
	repo.Create(&models.Checkpoint{Description: "C1"})
	repo.Create(&models.Checkpoint{Description: "C2"})

	checkpoints, err := repo.List()
	if err != nil {
		t.Fatalf("failed to list checkpoints: %v", err)
	}

	if len(checkpoints) != 2 {
		t.Errorf("expected 2 checkpoints, got %d", len(checkpoints))
	}
}

func TestCheckpointRepo_ListByProject(t *testing.T) {
	db := setupCheckpointDB(t)
	defer db.Close()

	repo := repos.NewCheckpointRepo(db)
	repo.Create(&models.Checkpoint{ProjectID: 1, Description: "P1"})
	repo.Create(&models.Checkpoint{ProjectID: 2, Description: "P2"})

	checkpoints, err := repo.ListByProject(1)
	if err != nil {
		t.Fatalf("failed to list by project: %v", err)
	}

	if len(checkpoints) != 1 {
		t.Errorf("expected 1 checkpoint, got %d", len(checkpoints))
	}
	if checkpoints[0].Description != "P1" {
		t.Errorf("expected description 'P1', got %q", checkpoints[0].Description)
	}
}

func TestCheckpointRepo_ListBySubject(t *testing.T) {
	db := setupCheckpointDB(t)
	defer db.Close()

	repo := repos.NewCheckpointRepo(db)
	repo.Create(&models.Checkpoint{SubjectID: 1, Description: "S1"})
	repo.Create(&models.Checkpoint{SubjectID: 2, Description: "S2"})

	checkpoints, err := repo.ListBySubject(1)
	if err != nil {
		t.Fatalf("failed to list by subject: %v", err)
	}

	if len(checkpoints) != 1 {
		t.Errorf("expected 1 checkpoint, got %d", len(checkpoints))
	}
	if checkpoints[0].Description != "S1" {
		t.Errorf("expected description 'S1', got %q", checkpoints[0].Description)
	}
}

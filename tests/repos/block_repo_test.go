// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT

package repos_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/khgreav/chronosplit/repos"
	_ "modernc.org/sqlite"
)

func setupBlockDB(t *testing.T) *sql.DB {
	return newTestDB(t, "")
}

func TestBlockRepo_StartBlock(t *testing.T) {
	db := setupBlockDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewBlockRepo(db)

	id, err := repo.StartBlock()
	if err != nil {
		t.Fatalf("failed to start block: %v", err)
	}

	if id == nil {
		t.Fatal("expected non-nil ID")
	}

	if *id <= 0 {
		t.Errorf("expected non-zero ID, got %d", *id)
	}
}

func TestBlockRepo_ActiveBlockExists(t *testing.T) {
	db := setupBlockDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewBlockRepo(db)

	exists, err := repo.ActiveBlockExists()
	if err != nil {
		t.Fatalf("failed to check active block existence: %v", err)
	}
	if exists {
		t.Error("expected no active block, but found one")
	}

	_, err = repo.StartBlock()
	if err != nil {
		t.Fatalf("failed to start block: %v", err)
	}

	exists, err = repo.ActiveBlockExists()
	if err != nil {
		t.Fatalf("failed to check active block existence: %v", err)
	}
	if !exists {
		t.Error("expected an active block, but found none")
	}
}

func TestBlockRepo_GetActiveBlock(t *testing.T) {
	db := setupBlockDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewBlockRepo(db)

	_, err := repo.GetActiveBlock()
	if err == nil {
		t.Error("expected error when getting active block when none exists, got nil")
	}

	id, err := repo.StartBlock()
	if err != nil {
		t.Fatalf("failed to start block: %v", err)
	}

	block, err := repo.GetActiveBlock()
	if err != nil {
		t.Fatalf("failed to get active block: %v", err)
	}

	if block.ID != *id {
		t.Errorf("expected ID %d, got %d", *id, block.ID)
	}
}

func TestBlockRepo_GetActiveBlockID(t *testing.T) {
	db := setupBlockDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewBlockRepo(db)

	_, err := repo.GetActiveBlockID()
	if err == nil {
		t.Error("expected error when getting active block ID when none exists, got nil")
	}

	id, err := repo.StartBlock()
	if err != nil {
		t.Fatalf("failed to start block: %v", err)
	}

	gotID, err := repo.GetActiveBlockID()
	if err != nil {
		t.Fatalf("failed to get active block ID: %v", err)
	}

	if gotID == nil || *gotID != *id {
		t.Errorf("expected ID %d, got %v", *id, gotID)
	}
}

func TestBlockRepo_StopBlock(t *testing.T) {
	db := setupBlockDB(t)
	defer func() {
		_ = db.Close()
	}()

	repo := repos.NewBlockRepo(db)

	id, err := repo.StartBlock()
	if err != nil {
		t.Fatalf("failed to start block: %v", err)
	}

	now := time.Now().UTC()
	err = repo.StopBlock(*id, now)
	if err != nil {
		t.Fatalf("failed to stop block: %v", err)
	}

	exists, err := repo.ActiveBlockExists()
	if err != nil {
		t.Fatalf("failed to check active block existence: %v", err)
	}
	if exists {
		t.Error("expected no active block after stopping, but found one")
	}

	err = repo.StopBlock(*id, now)
	if err == nil {
		t.Error("expected error when stopping an already stopped block, got nil")
	}
}

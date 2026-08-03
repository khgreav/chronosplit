// Copyright (c) 2026 Karel Hanák
// SPDX-License-Identifier: MIT
package migrations

var (
	Migrations = []string{
		`
			CREATE TABLE IF NOT EXISTS blocks (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				created_at DATETIME NOT NULL,
				ended_at DATETIME DEFAULT NULL
			);

			CREATE TABLE IF NOT EXISTS projects (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL
			);

			CREATE TABLE IF NOT EXISTS subjects (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL
			);

			CREATE TABLE IF NOT EXISTS checkpoints (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				block_id INTEGER NOT NULL REFERENCES blocks(id) ON DELETE CASCADE,
				project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
				subject_id INTEGER NOT NULL REFERENCES subjects(id) ON DELETE CASCADE,
				previous_checkpoint_id INTEGER REFERENCES checkpoints(id),
				start_time DATETIME NOT NULL,
				end_time DATETIME NOT NULL,
				description TEXT NOT NULL
			);
		`,
	}
)

/*
* Copyright 2022, The Cozo Project Authors.
*
* This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
* If a copy of the MPL was not distributed with this file,
* You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package cozo

import (
	"os"
	"testing"
)

func TestDb(t *testing.T) {
	db, err := New("mem", "", nil)
	if err != nil {
		t.Errorf("Failed to create: %v", err)
	}
	_, err = db.Run("?[a,b,c] <- [[1,2,3]] :create s{a, b, c}", nil, false)
	if err != nil {
		t.Errorf("Failed to run query: %v", err)
	}
	res, _ := db.Run("?[a,b,c] := *s[a,b,c]", nil, true)
	if len(res.Rows) != 1 || len(res.Rows[0]) != 3 {
		t.Error("Bad number of rows")
	}
	_, err = db.Run("?[x] <- [[1,2,3]]", nil, true)
	if err == nil {
		t.Error("Expected error wasn't thrown")
	}

	stat, err := os.Stat("test.db")
	if !os.IsNotExist(err) && !stat.IsDir() {
		err = os.Remove("test.db")
		if err != nil {
			t.Errorf("Failed to remove 'test.db' file: %v", err)
		}
	}

	err = db.Backup("test.db")
	if err != nil {
		t.Errorf("Failed to backup: %v", err)
	}

	db2, err := New("mem", "", nil)
	if err != nil {
		t.Errorf("Failed to create: %v", err)

	}
	err = db2.Restore("test.db")
	if err != nil {
		t.Errorf("Failed to restore: %v", err)
	}

	res, err = db2.Run("?[a,b,c] := *s[a,b,c]", nil, true)
	if err != nil {
		t.Errorf("Failed to run query: %v", err)
	}
	if len(res.Rows) != 1 || len(res.Rows[0]) != 3 {
		t.Error("Bad number of rows")
	}

	data, err := db.ExportRelations([]string{"s"})
	if err != nil {
		t.Errorf("Failed to export: %v", err)
	}
	db3, err := New("mem", "", nil)
	if err != nil {
		t.Errorf("Failed to create: %v", err)
	}
	_, err = db3.Run(":create s {a, b, c}", nil, false)
	if err != nil {
		t.Errorf("Failed to run query: %v", err)
	}

	res, err = db3.Run("?[a,b,c] := *s[a,b,c]", nil, true)
	if err != nil {
		t.Errorf("Failed to run query: %v", err)
	}
	if len(res.Rows) != 0 {
		t.Error("Bad number of rows")
	}
	err = db3.ImportRelations(data)
	if err != nil {
		t.Errorf("Failed to import: %v", err)
	}

	res, err = db3.Run("?[a,b,c] := *s[a,b,c]", nil, true)
	if err != nil {
		t.Errorf("Failed to run query: %v", err)
	}

	if len(res.Rows) != 1 || len(res.Rows[0]) != 3 {
		t.Error("Bad number of rows")
	}

	db4, err := New("mem", "", nil)
	if err != nil {
		t.Errorf("Failed to create: %v", err)
	}
	_, err = db4.Run(":create s {a, b, c}", nil, false)
	if err != nil {
		t.Errorf("Failed to run query: %v", err)
	}

	res, err = db4.Run("?[a,b,c] := *s[a,b,c]", nil, true)
	if err != nil {
		t.Errorf("Failed to run query: %v", err)
	}
	if len(res.Rows) != 0 {
		t.Error("Bad number of rows")
	}
	err = db4.ImportRelationsFromBackup("test.db", []string{"s"})
	if err != nil {
		t.Errorf("Failed to import: %v", err)
	}

	res, err = db4.Run("?[a,b,c] := *s[a,b,c]", nil, false)
	if err != nil {
		t.Errorf("Failed to run query: %v", err)
	}

	if len(res.Rows) != 1 || len(res.Rows[0]) != 3 {
		t.Error("Bad number of rows")
	}

	db.Close()
}

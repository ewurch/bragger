package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ewurch/resume-tracker/internal/models"
)

func setupTestStorage(t *testing.T) (*Storage, string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "app-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	filePath := filepath.Join(tmpDir, "test-applications.jsonl")
	store := New(filePath)

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return store, filePath, cleanup
}

func TestStorageNew(t *testing.T) {
	t.Run("default path", func(t *testing.T) {
		store := New("")
		if store.filePath != DefaultFilePath {
			t.Errorf("expected default path %q, got %q", DefaultFilePath, store.filePath)
		}
	})

	t.Run("custom path", func(t *testing.T) {
		customPath := "/tmp/custom.jsonl"
		store := New(customPath)
		if store.filePath != customPath {
			t.Errorf("expected path %q, got %q", customPath, store.filePath)
		}
	})
}

func TestStorageAddAndLoad(t *testing.T) {
	store, _, cleanup := setupTestStorage(t)
	defer cleanup()

	app := models.NewApplication("TestCorp", "Engineer")
	app.JDURL = "https://example.com/job"
	app.JDContent = "Test job description"
	app.Notes = "Test notes"

	// Add application
	if err := store.Add(app); err != nil {
		t.Fatalf("failed to add application: %v", err)
	}

	// Load and verify
	apps, err := store.Load()
	if err != nil {
		t.Fatalf("failed to load applications: %v", err)
	}

	if len(apps) != 1 {
		t.Fatalf("expected 1 application, got %d", len(apps))
	}

	loaded := apps[0]
	if loaded.Company != "TestCorp" {
		t.Errorf("expected company TestCorp, got %q", loaded.Company)
	}
	if loaded.Role != "Engineer" {
		t.Errorf("expected role Engineer, got %q", loaded.Role)
	}
	if loaded.JDURL != "https://example.com/job" {
		t.Errorf("expected JDURL, got %q", loaded.JDURL)
	}
	if loaded.JDContent != "Test job description" {
		t.Errorf("expected JDContent, got %q", loaded.JDContent)
	}
}

func TestStorageLoadEmpty(t *testing.T) {
	store, _, cleanup := setupTestStorage(t)
	defer cleanup()

	apps, err := store.Load()
	if err != nil {
		t.Fatalf("failed to load from non-existent file: %v", err)
	}

	if len(apps) != 0 {
		t.Errorf("expected 0 applications, got %d", len(apps))
	}
}

func TestStorageGet(t *testing.T) {
	store, _, cleanup := setupTestStorage(t)
	defer cleanup()

	app := models.NewApplication("TestCorp", "Engineer")
	store.Add(app)

	t.Run("existing ID", func(t *testing.T) {
		found, err := store.Get(app.ID)
		if err != nil {
			t.Fatalf("failed to get application: %v", err)
		}
		if found.Company != "TestCorp" {
			t.Errorf("expected company TestCorp, got %q", found.Company)
		}
	})

	t.Run("non-existing ID", func(t *testing.T) {
		_, err := store.Get("app-nonexistent")
		if err == nil {
			t.Error("expected error for non-existent ID")
		}
	})
}

func TestStorageUpdate(t *testing.T) {
	store, _, cleanup := setupTestStorage(t)
	defer cleanup()

	app := models.NewApplication("TestCorp", "Engineer")
	store.Add(app)

	t.Run("update status", func(t *testing.T) {
		err := store.Update(app.ID, func(a *models.Application) {
			a.Status = models.StatusInterviewing
		})
		if err != nil {
			t.Fatalf("failed to update: %v", err)
		}

		updated, _ := store.Get(app.ID)
		if updated.Status != models.StatusInterviewing {
			t.Errorf("expected status interviewing, got %q", updated.Status)
		}
	})

	t.Run("update notes", func(t *testing.T) {
		err := store.Update(app.ID, func(a *models.Application) {
			a.Notes = "Updated notes"
		})
		if err != nil {
			t.Fatalf("failed to update: %v", err)
		}

		updated, _ := store.Get(app.ID)
		if updated.Notes != "Updated notes" {
			t.Errorf("expected notes 'Updated notes', got %q", updated.Notes)
		}
	})

	t.Run("update non-existing ID", func(t *testing.T) {
		err := store.Update("app-nonexistent", func(a *models.Application) {
			a.Notes = "test"
		})
		if err == nil {
			t.Error("expected error for non-existent ID")
		}
	})
}

func TestStorageRemove(t *testing.T) {
	store, _, cleanup := setupTestStorage(t)
	defer cleanup()

	app1 := models.NewApplication("Company1", "Role1")
	app2 := models.NewApplication("Company2", "Role2")
	store.Add(app1)
	store.Add(app2)

	t.Run("remove existing", func(t *testing.T) {
		err := store.Remove(app1.ID)
		if err != nil {
			t.Fatalf("failed to remove: %v", err)
		}

		apps, _ := store.Load()
		if len(apps) != 1 {
			t.Errorf("expected 1 application, got %d", len(apps))
		}
		if apps[0].ID != app2.ID {
			t.Errorf("expected remaining app to be app2")
		}
	})

	t.Run("remove non-existing", func(t *testing.T) {
		err := store.Remove("app-nonexistent")
		if err == nil {
			t.Error("expected error for non-existent ID")
		}
	})
}

func TestStorageMultilineJDContent(t *testing.T) {
	store, _, cleanup := setupTestStorage(t)
	defer cleanup()

	multilineJD := `Senior Software Engineer

About the Role:
We are looking for an experienced engineer.

Requirements:
- 5+ years experience
- Strong Go skills
- Experience with distributed systems`

	app := models.NewApplication("TestCorp", "Senior Engineer")
	app.JDContent = multilineJD

	if err := store.Add(app); err != nil {
		t.Fatalf("failed to add: %v", err)
	}

	loaded, err := store.Get(app.ID)
	if err != nil {
		t.Fatalf("failed to get: %v", err)
	}

	if loaded.JDContent != multilineJD {
		t.Errorf("multiline JD content not preserved.\nExpected:\n%s\n\nGot:\n%s", multilineJD, loaded.JDContent)
	}
}

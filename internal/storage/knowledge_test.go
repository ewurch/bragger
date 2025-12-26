package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ewurch/brag/internal/models"
)

func setupKBTestStorage(t *testing.T) (*KBStorage, string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "kb-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	filePath := filepath.Join(tmpDir, "test-kb.jsonl")
	store := NewKBStorage(filePath)

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return store, filePath, cleanup
}

func TestKBStorageNew(t *testing.T) {
	t.Run("default path", func(t *testing.T) {
		store := NewKBStorage("")
		if store.filePath != DefaultKBFilePath {
			t.Errorf("expected default path %q, got %q", DefaultKBFilePath, store.filePath)
		}
	})

	t.Run("custom path", func(t *testing.T) {
		customPath := "/tmp/custom-kb.jsonl"
		store := NewKBStorage(customPath)
		if store.filePath != customPath {
			t.Errorf("expected path %q, got %q", customPath, store.filePath)
		}
	})
}

func TestKBStorageAddAndLoad(t *testing.T) {
	store, _, cleanup := setupKBTestStorage(t)
	defer cleanup()

	entry := models.NewProfileEntry(models.CategoryContact, models.ContactData{
		Name:  "John Doe",
		Email: "john@example.com",
	}, "test")

	// Add entry
	if err := store.Add(entry); err != nil {
		t.Fatalf("failed to add entry: %v", err)
	}

	// Load and verify
	entries, err := store.Load()
	if err != nil {
		t.Fatalf("failed to load entries: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}

	loaded := entries[0]
	if loaded.ID != entry.ID {
		t.Errorf("expected ID %q, got %q", entry.ID, loaded.ID)
	}
	if loaded.Type != models.KBTypeProfile {
		t.Errorf("expected type profile, got %q", loaded.Type)
	}
	if loaded.Category != string(models.CategoryContact) {
		t.Errorf("expected category contact, got %q", loaded.Category)
	}
}

func TestKBStorageLoadEmpty(t *testing.T) {
	store, _, cleanup := setupKBTestStorage(t)
	defer cleanup()

	entries, err := store.Load()
	if err != nil {
		t.Fatalf("failed to load from non-existent file: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(entries))
	}
}

func TestKBStorageGet(t *testing.T) {
	store, _, cleanup := setupKBTestStorage(t)
	defer cleanup()

	entry := models.NewContextEntry("achievement", "Led a team", "user")
	store.Add(entry)

	t.Run("existing ID", func(t *testing.T) {
		found, err := store.Get(entry.ID)
		if err != nil {
			t.Fatalf("failed to get entry: %v", err)
		}
		if found.Content != "Led a team" {
			t.Errorf("expected content 'Led a team', got %q", found.Content)
		}
	})

	t.Run("non-existing ID", func(t *testing.T) {
		_, err := store.Get("kb-nonexistent")
		if err == nil {
			t.Error("expected error for non-existent ID")
		}
	})
}

func TestKBStorageUpdate(t *testing.T) {
	store, _, cleanup := setupKBTestStorage(t)
	defer cleanup()

	entry := models.NewContextEntry("achievement", "Original content", "user")
	store.Add(entry)

	t.Run("update content", func(t *testing.T) {
		err := store.Update(entry.ID, func(e *models.KBEntry) {
			e.Content = "Updated content"
		})
		if err != nil {
			t.Fatalf("failed to update: %v", err)
		}

		updated, _ := store.Get(entry.ID)
		if updated.Content != "Updated content" {
			t.Errorf("expected content 'Updated content', got %q", updated.Content)
		}
	})

	t.Run("update non-existing ID", func(t *testing.T) {
		err := store.Update("kb-nonexistent", func(e *models.KBEntry) {
			e.Content = "test"
		})
		if err == nil {
			t.Error("expected error for non-existent ID")
		}
	})
}

func TestKBStorageRemove(t *testing.T) {
	store, _, cleanup := setupKBTestStorage(t)
	defer cleanup()

	entry1 := models.NewContextEntry("achievement", "First", "user")
	entry2 := models.NewContextEntry("achievement", "Second", "user")
	store.Add(entry1)
	store.Add(entry2)

	t.Run("remove existing", func(t *testing.T) {
		err := store.Remove(entry1.ID)
		if err != nil {
			t.Fatalf("failed to remove: %v", err)
		}

		entries, _ := store.Load()
		if len(entries) != 1 {
			t.Errorf("expected 1 entry, got %d", len(entries))
		}
		if entries[0].ID != entry2.ID {
			t.Errorf("expected remaining entry to be entry2")
		}
	})

	t.Run("remove non-existing", func(t *testing.T) {
		err := store.Remove("kb-nonexistent")
		if err == nil {
			t.Error("expected error for non-existent ID")
		}
	})
}

func TestKBStorageGetByType(t *testing.T) {
	store, _, cleanup := setupKBTestStorage(t)
	defer cleanup()

	// Add profile and context entries
	profile1 := models.NewProfileEntry(models.CategoryContact, models.ContactData{Name: "John", Email: "j@e.com"}, "test")
	profile2 := models.NewProfileEntry(models.CategoryExperience, models.ExperienceEntry{Company: "Acme", Role: "Dev", StartDate: "2020-01"}, "test")
	context1 := models.NewContextEntry("achievement", "Achievement 1", "user")

	store.Add(profile1)
	store.Add(profile2)
	store.Add(context1)

	t.Run("get profile entries", func(t *testing.T) {
		entries, err := store.GetByType(models.KBTypeProfile)
		if err != nil {
			t.Fatalf("failed to get by type: %v", err)
		}
		if len(entries) != 2 {
			t.Errorf("expected 2 profile entries, got %d", len(entries))
		}
	})

	t.Run("get context entries", func(t *testing.T) {
		entries, err := store.GetByType(models.KBTypeContext)
		if err != nil {
			t.Fatalf("failed to get by type: %v", err)
		}
		if len(entries) != 1 {
			t.Errorf("expected 1 context entry, got %d", len(entries))
		}
	})
}

func TestKBStorageGetByCategory(t *testing.T) {
	store, _, cleanup := setupKBTestStorage(t)
	defer cleanup()

	exp1 := models.NewProfileEntry(models.CategoryExperience, models.ExperienceEntry{Company: "Acme", Role: "Dev", StartDate: "2020-01"}, "test")
	exp2 := models.NewProfileEntry(models.CategoryExperience, models.ExperienceEntry{Company: "Beta", Role: "Senior", StartDate: "2022-01"}, "test")
	edu := models.NewProfileEntry(models.CategoryEducation, models.EducationEntry{Institution: "MIT", Degree: "BS"}, "test")

	store.Add(exp1)
	store.Add(exp2)
	store.Add(edu)

	t.Run("get experience entries", func(t *testing.T) {
		entries, err := store.GetByCategory(string(models.CategoryExperience))
		if err != nil {
			t.Fatalf("failed to get by category: %v", err)
		}
		if len(entries) != 2 {
			t.Errorf("expected 2 experience entries, got %d", len(entries))
		}
	})

	t.Run("get education entries", func(t *testing.T) {
		entries, err := store.GetByCategory(string(models.CategoryEducation))
		if err != nil {
			t.Fatalf("failed to get by category: %v", err)
		}
		if len(entries) != 1 {
			t.Errorf("expected 1 education entry, got %d", len(entries))
		}
	})
}

func TestKBStorageHelperMethods(t *testing.T) {
	store, _, cleanup := setupKBTestStorage(t)
	defer cleanup()

	contact := models.NewProfileEntry(models.CategoryContact, models.ContactData{Name: "John", Email: "j@e.com"}, "test")
	exp := models.NewProfileEntry(models.CategoryExperience, models.ExperienceEntry{Company: "Acme", Role: "Dev", StartDate: "2020-01"}, "test")
	edu := models.NewProfileEntry(models.CategoryEducation, models.EducationEntry{Institution: "MIT", Degree: "BS"}, "test")
	skills := models.NewProfileEntry(models.CategorySkills, models.SkillsData{Languages: []string{"Go", "Python"}}, "test")

	store.Add(contact)
	store.Add(exp)
	store.Add(edu)
	store.Add(skills)

	t.Run("GetProfile", func(t *testing.T) {
		entries, err := store.GetProfile()
		if err != nil {
			t.Fatalf("GetProfile failed: %v", err)
		}
		if len(entries) != 4 {
			t.Errorf("expected 4 profile entries, got %d", len(entries))
		}
	})

	t.Run("GetContact", func(t *testing.T) {
		entry, err := store.GetContact()
		if err != nil {
			t.Fatalf("GetContact failed: %v", err)
		}
		if entry.ID != contact.ID {
			t.Errorf("expected contact ID %q, got %q", contact.ID, entry.ID)
		}
	})

	t.Run("GetExperience", func(t *testing.T) {
		entries, err := store.GetExperience()
		if err != nil {
			t.Fatalf("GetExperience failed: %v", err)
		}
		if len(entries) != 1 {
			t.Errorf("expected 1 experience entry, got %d", len(entries))
		}
	})

	t.Run("GetEducation", func(t *testing.T) {
		entries, err := store.GetEducation()
		if err != nil {
			t.Fatalf("GetEducation failed: %v", err)
		}
		if len(entries) != 1 {
			t.Errorf("expected 1 education entry, got %d", len(entries))
		}
	})

	t.Run("GetSkills", func(t *testing.T) {
		entry, err := store.GetSkills()
		if err != nil {
			t.Fatalf("GetSkills failed: %v", err)
		}
		if entry.ID != skills.ID {
			t.Errorf("expected skills ID %q, got %q", skills.ID, entry.ID)
		}
	})
}

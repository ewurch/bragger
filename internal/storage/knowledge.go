package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/ewurch/bragger/internal/models"
)

const DefaultKBFilePath = "candidate-kb.jsonl"

type KBStorage struct {
	filePath string
}

func NewKBStorage(filePath string) *KBStorage {
	if filePath == "" {
		filePath = DefaultKBFilePath
	}
	return &KBStorage{filePath: filePath}
}

func (s *KBStorage) Load() ([]*models.KBEntry, error) {
	file, err := os.Open(s.filePath)
	if os.IsNotExist(err) {
		return []*models.KBEntry{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var entries []*models.KBEntry
	scanner := bufio.NewScanner(file)

	// Increase buffer size for large entries
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		var entry models.KBEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue // Skip malformed lines
		}
		entries = append(entries, &entry)
	}
	return entries, scanner.Err()
}

func (s *KBStorage) Save(entries []*models.KBEntry) error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, entry := range entries {
		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		file.Write(data)
		file.WriteString("\n")
	}
	return nil
}

func (s *KBStorage) Add(entry *models.KBEntry) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}
	entries = append(entries, entry)
	return s.Save(entries)
}

func (s *KBStorage) Get(id string) (*models.KBEntry, error) {
	entries, err := s.Load()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.ID == id {
			return entry, nil
		}
	}
	return nil, os.ErrNotExist
}

func (s *KBStorage) Update(id string, updateFn func(*models.KBEntry)) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}

	found := false
	for _, entry := range entries {
		if entry.ID == id {
			updateFn(entry)
			entry.UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return os.ErrNotExist
	}

	return s.Save(entries)
}

func (s *KBStorage) Remove(id string) error {
	entries, err := s.Load()
	if err != nil {
		return err
	}

	var filtered []*models.KBEntry
	found := false
	for _, entry := range entries {
		if entry.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, entry)
	}

	if !found {
		return os.ErrNotExist
	}

	return s.Save(filtered)
}

// GetByType returns all entries of a specific type
func (s *KBStorage) GetByType(entryType models.KBEntryType) ([]*models.KBEntry, error) {
	entries, err := s.Load()
	if err != nil {
		return nil, err
	}

	var filtered []*models.KBEntry
	for _, entry := range entries {
		if entry.Type == entryType {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}

// GetByCategory returns all entries of a specific category
func (s *KBStorage) GetByCategory(category string) ([]*models.KBEntry, error) {
	entries, err := s.Load()
	if err != nil {
		return nil, err
	}

	var filtered []*models.KBEntry
	for _, entry := range entries {
		if entry.Category == category {
			filtered = append(filtered, entry)
		}
	}
	return filtered, nil
}

// GetProfile returns all profile entries
func (s *KBStorage) GetProfile() ([]*models.KBEntry, error) {
	return s.GetByType(models.KBTypeProfile)
}

// GetContext returns all contextual entries
func (s *KBStorage) GetContext() ([]*models.KBEntry, error) {
	return s.GetByType(models.KBTypeContext)
}

// GetContact returns the contact profile entry (there should be at most one)
func (s *KBStorage) GetContact() (*models.KBEntry, error) {
	entries, err := s.GetByCategory(string(models.CategoryContact))
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, os.ErrNotExist
	}
	return entries[0], nil
}

// GetExperience returns all experience entries
func (s *KBStorage) GetExperience() ([]*models.KBEntry, error) {
	return s.GetByCategory(string(models.CategoryExperience))
}

// GetEducation returns all education entries
func (s *KBStorage) GetEducation() ([]*models.KBEntry, error) {
	return s.GetByCategory(string(models.CategoryEducation))
}

// GetSkills returns the skills profile entry (there should be at most one)
func (s *KBStorage) GetSkills() (*models.KBEntry, error) {
	entries, err := s.GetByCategory(string(models.CategorySkills))
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, os.ErrNotExist
	}
	return entries[0], nil
}

package storage

import (
	"bufio"
	"encoding/json"
	"os"
	"time"

	"github.com/ewurch/resume-tracker/internal/models"
)

const DefaultFilePath = "applications.jsonl"

type Storage struct {
	filePath string
}

func New(filePath string) *Storage {
	if filePath == "" {
		filePath = DefaultFilePath
	}
	return &Storage{filePath: filePath}
}

func (s *Storage) Load() ([]*models.Application, error) {
	file, err := os.Open(s.filePath)
	if os.IsNotExist(err) {
		return []*models.Application{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var apps []*models.Application
	scanner := bufio.NewScanner(file)

	// Increase buffer size for large JD content
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	for scanner.Scan() {
		var app models.Application
		if err := json.Unmarshal(scanner.Bytes(), &app); err != nil {
			continue // Skip malformed lines
		}
		apps = append(apps, &app)
	}
	return apps, scanner.Err()
}

func (s *Storage) Save(apps []*models.Application) error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, app := range apps {
		data, err := json.Marshal(app)
		if err != nil {
			return err
		}
		file.Write(data)
		file.WriteString("\n")
	}
	return nil
}

func (s *Storage) Add(app *models.Application) error {
	apps, err := s.Load()
	if err != nil {
		return err
	}
	apps = append(apps, app)
	return s.Save(apps)
}

func (s *Storage) Update(id string, updateFn func(*models.Application)) error {
	apps, err := s.Load()
	if err != nil {
		return err
	}

	found := false
	for _, app := range apps {
		if app.ID == id {
			updateFn(app)
			app.UpdatedAt = time.Now()
			found = true
			break
		}
	}

	if !found {
		return os.ErrNotExist
	}

	return s.Save(apps)
}

func (s *Storage) Remove(id string) error {
	apps, err := s.Load()
	if err != nil {
		return err
	}

	var filtered []*models.Application
	found := false
	for _, app := range apps {
		if app.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, app)
	}

	if !found {
		return os.ErrNotExist
	}

	return s.Save(filtered)
}

func (s *Storage) Get(id string) (*models.Application, error) {
	apps, err := s.Load()
	if err != nil {
		return nil, err
	}

	for _, app := range apps {
		if app.ID == id {
			return app, nil
		}
	}
	return nil, os.ErrNotExist
}

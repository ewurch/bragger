package models

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Status string

const (
	StatusApplied      Status = "applied"
	StatusInterviewing Status = "interviewing"
	StatusRejected     Status = "rejected"
	StatusOffer        Status = "offer"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusApplied, StatusInterviewing, StatusRejected, StatusOffer:
		return true
	}
	return false
}

type Application struct {
	ID          string    `json:"id"`
	Company     string    `json:"company"`
	Role        string    `json:"role"`
	Status      Status    `json:"status"`
	DateApplied string    `json:"date_applied"`
	JDURL       string    `json:"jd_url,omitempty"`
	JDContent   string    `json:"jd_content,omitempty"`
	ResumePath  string    `json:"resume_path,omitempty"`
	CompanyURL  string    `json:"company_url,omitempty"`
	Notes       string    `json:"notes,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func GenerateID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return "app-" + hex.EncodeToString(bytes)
}

func NewApplication(company, role string) *Application {
	now := time.Now()
	return &Application{
		ID:          GenerateID(),
		Company:     company,
		Role:        role,
		Status:      StatusApplied,
		DateApplied: now.Format("2006-01-02"),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

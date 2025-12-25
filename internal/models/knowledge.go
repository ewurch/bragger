package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// KBEntryType distinguishes between profile and contextual entries
type KBEntryType string

const (
	KBTypeProfile KBEntryType = "profile"
	KBTypeContext KBEntryType = "context"
)

func (t KBEntryType) IsValid() bool {
	return t == KBTypeProfile || t == KBTypeContext
}

// ProfileCategory defines the categories for structured profile data
type ProfileCategory string

const (
	CategoryContact        ProfileCategory = "contact"
	CategoryExperience     ProfileCategory = "experience"
	CategoryEducation      ProfileCategory = "education"
	CategorySkills         ProfileCategory = "skills"
	CategoryCertifications ProfileCategory = "certifications"
	CategoryLanguages      ProfileCategory = "languages"
)

var validProfileCategories = []ProfileCategory{
	CategoryContact, CategoryExperience, CategoryEducation,
	CategorySkills, CategoryCertifications, CategoryLanguages,
}

func (c ProfileCategory) IsValid() bool {
	for _, valid := range validProfileCategories {
		if c == valid {
			return true
		}
	}
	return false
}

// ContactData holds candidate contact information
type ContactData struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone,omitempty"`
	Location string `json:"location,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`
	GitHub   string `json:"github,omitempty"`
	Website  string `json:"website,omitempty"`
}

func (c *ContactData) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

// ExperienceEntry holds a single work experience entry
type ExperienceEntry struct {
	Company     string   `json:"company"`
	Role        string   `json:"role"`
	StartDate   string   `json:"start_date"`         // YYYY-MM format
	EndDate     string   `json:"end_date,omitempty"` // YYYY-MM or "present"
	Location    string   `json:"location,omitempty"`
	Description string   `json:"description,omitempty"`
	Highlights  []string `json:"highlights,omitempty"`
}

func (e *ExperienceEntry) Validate() error {
	if e.Company == "" {
		return fmt.Errorf("company is required")
	}
	if e.Role == "" {
		return fmt.Errorf("role is required")
	}
	if e.StartDate == "" {
		return fmt.Errorf("start_date is required")
	}
	return nil
}

// EducationEntry holds a single education entry
type EducationEntry struct {
	Institution string `json:"institution"`
	Degree      string `json:"degree"`
	Field       string `json:"field,omitempty"`
	StartDate   string `json:"start_date,omitempty"`
	EndDate     string `json:"end_date,omitempty"`
	GPA         string `json:"gpa,omitempty"`
}

func (e *EducationEntry) Validate() error {
	if e.Institution == "" {
		return fmt.Errorf("institution is required")
	}
	if e.Degree == "" {
		return fmt.Errorf("degree is required")
	}
	return nil
}

// SkillsData holds categorized skills
type SkillsData struct {
	Languages  []string `json:"languages,omitempty"` // Programming languages
	Frameworks []string `json:"frameworks,omitempty"`
	Tools      []string `json:"tools,omitempty"`
	Databases  []string `json:"databases,omitempty"`
	Cloud      []string `json:"cloud,omitempty"`
	Other      []string `json:"other,omitempty"`
}

func (s *SkillsData) Validate() error {
	// Skills can be partially filled, no required fields
	return nil
}

func (s *SkillsData) IsEmpty() bool {
	return len(s.Languages) == 0 && len(s.Frameworks) == 0 &&
		len(s.Tools) == 0 && len(s.Databases) == 0 &&
		len(s.Cloud) == 0 && len(s.Other) == 0
}

// CertificationEntry holds a single certification
type CertificationEntry struct {
	Name         string `json:"name"`
	Issuer       string `json:"issuer,omitempty"`
	Date         string `json:"date,omitempty"`        // YYYY-MM format
	ExpiryDate   string `json:"expiry_date,omitempty"` // YYYY-MM format
	CredentialID string `json:"credential_id,omitempty"`
}

func (c *CertificationEntry) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("certification name is required")
	}
	return nil
}

// LanguageEntry holds spoken/written language proficiency
type LanguageEntry struct {
	Language    string `json:"language"`
	Proficiency string `json:"proficiency,omitempty"` // e.g., "native", "fluent", "intermediate", "basic"
}

func (l *LanguageEntry) Validate() error {
	if l.Language == "" {
		return fmt.Errorf("language is required")
	}
	return nil
}

// KBEntry is the unified wrapper for all knowledge base entries
type KBEntry struct {
	ID        string      `json:"id"`
	Type      KBEntryType `json:"type"`
	Category  string      `json:"category"`
	Data      any         `json:"data,omitempty"`    // For profile entries (typed struct)
	Content   string      `json:"content,omitempty"` // For context entries (free text)
	Source    string      `json:"source,omitempty"`  // "cv-import", "user", "app-xxx"
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

// GenerateKBID generates a unique ID for KB entries
func GenerateKBID() string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return "kb-" + hex.EncodeToString(bytes)
}

// NewProfileEntry creates a new profile KB entry
func NewProfileEntry(category ProfileCategory, data any, source string) *KBEntry {
	now := time.Now()
	return &KBEntry{
		ID:        GenerateKBID(),
		Type:      KBTypeProfile,
		Category:  string(category),
		Data:      data,
		Source:    source,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// NewContextEntry creates a new contextual KB entry
func NewContextEntry(category string, content string, source string) *KBEntry {
	now := time.Now()
	return &KBEntry{
		ID:        GenerateKBID(),
		Type:      KBTypeContext,
		Category:  category,
		Content:   content,
		Source:    source,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Validate checks if the KB entry is valid
func (e *KBEntry) Validate() error {
	if !e.Type.IsValid() {
		return fmt.Errorf("invalid entry type: %s", e.Type)
	}

	if e.Category == "" {
		return fmt.Errorf("category is required")
	}

	if e.Type == KBTypeProfile {
		cat := ProfileCategory(e.Category)
		if !cat.IsValid() {
			return fmt.Errorf("invalid profile category: %s", e.Category)
		}
		if e.Data == nil {
			return fmt.Errorf("data is required for profile entries")
		}
	}

	if e.Type == KBTypeContext {
		if e.Content == "" {
			return fmt.Errorf("content is required for context entries")
		}
	}

	return nil
}

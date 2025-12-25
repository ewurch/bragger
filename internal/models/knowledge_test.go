package models

import (
	"strings"
	"testing"
)

func TestKBEntryTypeIsValid(t *testing.T) {
	tests := []struct {
		entryType KBEntryType
		valid     bool
	}{
		{KBTypeProfile, true},
		{KBTypeContext, true},
		{KBEntryType("invalid"), false},
		{KBEntryType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.entryType), func(t *testing.T) {
			if got := tt.entryType.IsValid(); got != tt.valid {
				t.Errorf("KBEntryType(%q).IsValid() = %v, want %v", tt.entryType, got, tt.valid)
			}
		})
	}
}

func TestProfileCategoryIsValid(t *testing.T) {
	tests := []struct {
		category ProfileCategory
		valid    bool
	}{
		{CategoryContact, true},
		{CategoryExperience, true},
		{CategoryEducation, true},
		{CategorySkills, true},
		{CategoryCertifications, true},
		{CategoryLanguages, true},
		{ProfileCategory("invalid"), false},
		{ProfileCategory(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.category), func(t *testing.T) {
			if got := tt.category.IsValid(); got != tt.valid {
				t.Errorf("ProfileCategory(%q).IsValid() = %v, want %v", tt.category, got, tt.valid)
			}
		})
	}
}

func TestContactDataValidate(t *testing.T) {
	tests := []struct {
		name    string
		data    ContactData
		wantErr bool
	}{
		{
			name:    "valid contact",
			data:    ContactData{Name: "John Doe", Email: "john@example.com"},
			wantErr: false,
		},
		{
			name:    "valid with all fields",
			data:    ContactData{Name: "John", Email: "j@e.com", Phone: "123", Location: "NYC", LinkedIn: "linkedin.com/in/john"},
			wantErr: false,
		},
		{
			name:    "missing name",
			data:    ContactData{Email: "john@example.com"},
			wantErr: true,
		},
		{
			name:    "missing email",
			data:    ContactData{Name: "John Doe"},
			wantErr: true,
		},
		{
			name:    "empty",
			data:    ContactData{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ContactData.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExperienceEntryValidate(t *testing.T) {
	tests := []struct {
		name    string
		data    ExperienceEntry
		wantErr bool
	}{
		{
			name:    "valid experience",
			data:    ExperienceEntry{Company: "Acme", Role: "Engineer", StartDate: "2020-01"},
			wantErr: false,
		},
		{
			name:    "valid with all fields",
			data:    ExperienceEntry{Company: "Acme", Role: "Engineer", StartDate: "2020-01", EndDate: "present", Location: "NYC", Description: "Built things", Highlights: []string{"Led team"}},
			wantErr: false,
		},
		{
			name:    "missing company",
			data:    ExperienceEntry{Role: "Engineer", StartDate: "2020-01"},
			wantErr: true,
		},
		{
			name:    "missing role",
			data:    ExperienceEntry{Company: "Acme", StartDate: "2020-01"},
			wantErr: true,
		},
		{
			name:    "missing start_date",
			data:    ExperienceEntry{Company: "Acme", Role: "Engineer"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("ExperienceEntry.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEducationEntryValidate(t *testing.T) {
	tests := []struct {
		name    string
		data    EducationEntry
		wantErr bool
	}{
		{
			name:    "valid education",
			data:    EducationEntry{Institution: "MIT", Degree: "BS Computer Science"},
			wantErr: false,
		},
		{
			name:    "missing institution",
			data:    EducationEntry{Degree: "BS"},
			wantErr: true,
		},
		{
			name:    "missing degree",
			data:    EducationEntry{Institution: "MIT"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("EducationEntry.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCertificationEntryValidate(t *testing.T) {
	tests := []struct {
		name    string
		data    CertificationEntry
		wantErr bool
	}{
		{
			name:    "valid certification",
			data:    CertificationEntry{Name: "AWS Solutions Architect"},
			wantErr: false,
		},
		{
			name:    "valid with all fields",
			data:    CertificationEntry{Name: "AWS SA", Issuer: "Amazon", Date: "2023-01", CredentialID: "ABC123"},
			wantErr: false,
		},
		{
			name:    "missing name",
			data:    CertificationEntry{Issuer: "Amazon"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CertificationEntry.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLanguageEntryValidate(t *testing.T) {
	tests := []struct {
		name    string
		data    LanguageEntry
		wantErr bool
	}{
		{
			name:    "valid language",
			data:    LanguageEntry{Language: "English"},
			wantErr: false,
		},
		{
			name:    "valid with proficiency",
			data:    LanguageEntry{Language: "Spanish", Proficiency: "fluent"},
			wantErr: false,
		},
		{
			name:    "missing language",
			data:    LanguageEntry{Proficiency: "native"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.data.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("LanguageEntry.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSkillsDataIsEmpty(t *testing.T) {
	tests := []struct {
		name    string
		data    SkillsData
		isEmpty bool
	}{
		{
			name:    "empty",
			data:    SkillsData{},
			isEmpty: true,
		},
		{
			name:    "with languages",
			data:    SkillsData{Languages: []string{"Go"}},
			isEmpty: false,
		},
		{
			name:    "with tools",
			data:    SkillsData{Tools: []string{"Git"}},
			isEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.data.IsEmpty(); got != tt.isEmpty {
				t.Errorf("SkillsData.IsEmpty() = %v, want %v", got, tt.isEmpty)
			}
		})
	}
}

func TestGenerateKBID(t *testing.T) {
	id1 := GenerateKBID()
	id2 := GenerateKBID()

	if id1 == id2 {
		t.Error("expected unique IDs, got duplicates")
	}
	if !strings.HasPrefix(id1, "kb-") {
		t.Errorf("expected ID to start with 'kb-', got %q", id1)
	}
	if len(id1) != 11 { // "kb-" + 8 hex chars
		t.Errorf("expected ID length 11, got %d", len(id1))
	}
}

func TestNewProfileEntry(t *testing.T) {
	data := ContactData{Name: "John", Email: "j@e.com"}
	entry := NewProfileEntry(CategoryContact, data, "test-source")

	if entry.Type != KBTypeProfile {
		t.Errorf("expected type %q, got %q", KBTypeProfile, entry.Type)
	}
	if entry.Category != string(CategoryContact) {
		t.Errorf("expected category %q, got %q", CategoryContact, entry.Category)
	}
	if entry.Source != "test-source" {
		t.Errorf("expected source 'test-source', got %q", entry.Source)
	}
	if !strings.HasPrefix(entry.ID, "kb-") {
		t.Errorf("expected ID to start with 'kb-', got %q", entry.ID)
	}
}

func TestNewContextEntry(t *testing.T) {
	entry := NewContextEntry("achievement", "Led a team", "user")

	if entry.Type != KBTypeContext {
		t.Errorf("expected type %q, got %q", KBTypeContext, entry.Type)
	}
	if entry.Category != "achievement" {
		t.Errorf("expected category 'achievement', got %q", entry.Category)
	}
	if entry.Content != "Led a team" {
		t.Errorf("expected content 'Led a team', got %q", entry.Content)
	}
}

func TestKBEntryValidate(t *testing.T) {
	tests := []struct {
		name    string
		entry   KBEntry
		wantErr bool
	}{
		{
			name: "valid profile entry",
			entry: KBEntry{
				ID:       "kb-test",
				Type:     KBTypeProfile,
				Category: string(CategoryContact),
				Data:     ContactData{Name: "John", Email: "j@e.com"},
			},
			wantErr: false,
		},
		{
			name: "valid context entry",
			entry: KBEntry{
				ID:       "kb-test",
				Type:     KBTypeContext,
				Category: "achievement",
				Content:  "Some achievement",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			entry: KBEntry{
				ID:       "kb-test",
				Type:     KBEntryType("invalid"),
				Category: "test",
			},
			wantErr: true,
		},
		{
			name: "missing category",
			entry: KBEntry{
				ID:   "kb-test",
				Type: KBTypeProfile,
			},
			wantErr: true,
		},
		{
			name: "profile with invalid category",
			entry: KBEntry{
				ID:       "kb-test",
				Type:     KBTypeProfile,
				Category: "invalid-category",
				Data:     "test",
			},
			wantErr: true,
		},
		{
			name: "profile without data",
			entry: KBEntry{
				ID:       "kb-test",
				Type:     KBTypeProfile,
				Category: string(CategoryContact),
			},
			wantErr: true,
		},
		{
			name: "context without content",
			entry: KBEntry{
				ID:       "kb-test",
				Type:     KBTypeContext,
				Category: "achievement",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entry.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("KBEntry.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

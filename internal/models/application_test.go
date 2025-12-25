package models

import (
	"strings"
	"testing"
	"time"
)

func TestNewApplication(t *testing.T) {
	company := "Test Corp"
	role := "Software Engineer"

	app := NewApplication(company, role)

	if app.Company != company {
		t.Errorf("expected company %q, got %q", company, app.Company)
	}
	if app.Role != role {
		t.Errorf("expected role %q, got %q", role, app.Role)
	}
	if app.Status != StatusApplied {
		t.Errorf("expected status %q, got %q", StatusApplied, app.Status)
	}
	if app.DateApplied != time.Now().Format("2006-01-02") {
		t.Errorf("expected date_applied to be today, got %q", app.DateApplied)
	}
	if !strings.HasPrefix(app.ID, "app-") {
		t.Errorf("expected ID to start with 'app-', got %q", app.ID)
	}
	if len(app.ID) != 12 { // "app-" + 8 hex chars
		t.Errorf("expected ID length 12, got %d", len(app.ID))
	}
}

func TestGenerateID(t *testing.T) {
	id1 := GenerateID()
	id2 := GenerateID()

	if id1 == id2 {
		t.Error("expected unique IDs, got duplicates")
	}
	if !strings.HasPrefix(id1, "app-") {
		t.Errorf("expected ID to start with 'app-', got %q", id1)
	}
}

func TestStatusIsValid(t *testing.T) {
	tests := []struct {
		status Status
		valid  bool
	}{
		{StatusApplied, true},
		{StatusInterviewing, true},
		{StatusRejected, true},
		{StatusOffer, true},
		{Status("invalid"), false},
		{Status(""), false},
		{Status("APPLIED"), false}, // case sensitive
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.valid {
				t.Errorf("Status(%q).IsValid() = %v, want %v", tt.status, got, tt.valid)
			}
		})
	}
}

func TestApplicationFields(t *testing.T) {
	app := NewApplication("Company", "Role")

	// Test setting optional fields
	app.JDURL = "https://example.com/job"
	app.JDContent = "This is a job description"
	app.CompanyURL = "https://example.com"
	app.ResumePath = "outputs/resume.html"
	app.Notes = "Applied via referral"

	if app.JDURL != "https://example.com/job" {
		t.Errorf("JDURL not set correctly")
	}
	if app.JDContent != "This is a job description" {
		t.Errorf("JDContent not set correctly")
	}
	if app.CompanyURL != "https://example.com" {
		t.Errorf("CompanyURL not set correctly")
	}
	if app.ResumePath != "outputs/resume.html" {
		t.Errorf("ResumePath not set correctly")
	}
	if app.Notes != "Applied via referral" {
		t.Errorf("Notes not set correctly")
	}
}

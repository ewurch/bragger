package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	// Build the binary for integration tests
	tmpDir, err := os.MkdirTemp("", "app-integration-test-*")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath = filepath.Join(tmpDir, "app")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = "."
	if err := cmd.Run(); err != nil {
		panic("failed to build binary: " + err.Error())
	}

	os.Exit(m.Run())
}

func setupIntegrationTest(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "app-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	// Create empty applications.jsonl
	jsonlPath := filepath.Join(tmpDir, "applications.jsonl")
	os.WriteFile(jsonlPath, []byte{}, 0644)

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func runApp(t *testing.T, workDir string, args ...string) (string, error) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)
	cmd.Dir = workDir
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func TestCLIAddWithFlags(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Run("minimal flags", func(t *testing.T) {
		output, err := runApp(t, workDir, "add", "--company", "TestCorp", "--role", "Engineer")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Application added successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
		if !strings.Contains(output, "TestCorp") {
			t.Errorf("expected company in output, got: %s", output)
		}
	})

	t.Run("with jd-url", func(t *testing.T) {
		output, err := runApp(t, workDir, "add",
			"--company", "URLCorp",
			"--role", "Developer",
			"--jd-url", "https://example.com/job")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Application added successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
	})

	t.Run("with jd-content inline", func(t *testing.T) {
		output, err := runApp(t, workDir, "add",
			"--company", "InlineCorp",
			"--role", "Designer",
			"--jd-content", "Looking for a creative designer")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Application added successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
	})

	t.Run("with jd-file", func(t *testing.T) {
		// Create temp JD file
		jdPath := filepath.Join(workDir, "jd.txt")
		jdContent := "Senior Engineer\n\nRequirements:\n- 5+ years\n- Go expertise"
		os.WriteFile(jdPath, []byte(jdContent), 0644)

		output, err := runApp(t, workDir, "add",
			"--company", "FileCorp",
			"--role", "Senior Engineer",
			"--jd-file", jdPath)
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Application added successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
	})

	t.Run("with all optional flags", func(t *testing.T) {
		output, err := runApp(t, workDir, "add",
			"--company", "FullCorp",
			"--role", "Staff Engineer",
			"--jd-url", "https://example.com/job",
			"--jd-content", "Full stack role",
			"--company-url", "https://fullcorp.com",
			"--resume-path", "outputs/resume.html",
			"--notes", "Referral from John")
		// This should fail because both jd-content and jd-file logic
		// Wait, jd-content and jd-url can coexist, only jd-content and jd-file conflict
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
	})

	t.Run("with custom status", func(t *testing.T) {
		output, err := runApp(t, workDir, "add",
			"--company", "StatusCorp",
			"--role", "Engineer",
			"--status", "interviewing")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Application added successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
		if !strings.Contains(output, "interviewing") {
			t.Errorf("expected status in output, got: %s", output)
		}
	})

	t.Run("with custom date", func(t *testing.T) {
		output, err := runApp(t, workDir, "add",
			"--company", "DateCorp",
			"--role", "Developer",
			"--date", "2025-01-15")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Application added successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
		if !strings.Contains(output, "2025-01-15") {
			t.Errorf("expected date in output, got: %s", output)
		}
	})

	t.Run("with status and date for backfill", func(t *testing.T) {
		output, err := runApp(t, workDir, "add",
			"--company", "OldCorp",
			"--role", "Senior Dev",
			"--status", "offer",
			"--date", "2024-12-01")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "offer") {
			t.Errorf("expected status in output, got: %s", output)
		}
		if !strings.Contains(output, "2024-12-01") {
			t.Errorf("expected date in output, got: %s", output)
		}
	})
}

func TestCLIAddErrors(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Run("missing role", func(t *testing.T) {
		output, err := runApp(t, workDir, "add", "--company", "TestCorp")
		if err == nil {
			t.Error("expected error when role is missing")
		}
		if !strings.Contains(output, "--company and --role are required") {
			t.Errorf("expected error message about required flags, got: %s", output)
		}
	})

	t.Run("missing company", func(t *testing.T) {
		output, err := runApp(t, workDir, "add", "--role", "Engineer")
		if err == nil {
			t.Error("expected error when company is missing")
		}
		if !strings.Contains(output, "--company and --role are required") {
			t.Errorf("expected error message about required flags, got: %s", output)
		}
	})

	t.Run("both jd-content and jd-file", func(t *testing.T) {
		jdPath := filepath.Join(workDir, "jd.txt")
		os.WriteFile(jdPath, []byte("test"), 0644)

		output, err := runApp(t, workDir, "add",
			"--company", "TestCorp",
			"--role", "Engineer",
			"--jd-content", "inline content",
			"--jd-file", jdPath)
		if err == nil {
			t.Error("expected error when both jd-content and jd-file provided")
		}
		if !strings.Contains(output, "--jd-content and --jd-file cannot be used together") {
			t.Errorf("expected conflict error message, got: %s", output)
		}
	})

	t.Run("jd-file not found", func(t *testing.T) {
		output, err := runApp(t, workDir, "add",
			"--company", "TestCorp",
			"--role", "Engineer",
			"--jd-file", "/nonexistent/path/jd.txt")
		if err == nil {
			t.Error("expected error when jd-file not found")
		}
		if !strings.Contains(output, "Error reading JD file") {
			t.Errorf("expected file read error message, got: %s", output)
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		output, err := runApp(t, workDir, "add",
			"--company", "TestCorp",
			"--role", "Engineer",
			"--status", "invalid-status")
		if err == nil {
			t.Error("expected error for invalid status")
		}
		if !strings.Contains(output, "--status must be one of") {
			t.Errorf("expected status error message, got: %s", output)
		}
	})

	t.Run("invalid date format", func(t *testing.T) {
		output, err := runApp(t, workDir, "add",
			"--company", "TestCorp",
			"--role", "Engineer",
			"--date", "15-01-2025")
		if err == nil {
			t.Error("expected error for invalid date format")
		}
		if !strings.Contains(output, "--date must be in YYYY-MM-DD format") {
			t.Errorf("expected date format error message, got: %s", output)
		}
	})

	t.Run("invalid date value", func(t *testing.T) {
		output, err := runApp(t, workDir, "add",
			"--company", "TestCorp",
			"--role", "Engineer",
			"--date", "not-a-date")
		if err == nil {
			t.Error("expected error for invalid date")
		}
		if !strings.Contains(output, "--date must be in YYYY-MM-DD format") {
			t.Errorf("expected date format error message, got: %s", output)
		}
	})
}

func TestCLIList(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	t.Run("empty list", func(t *testing.T) {
		output, err := runApp(t, workDir, "list")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "No applications found") {
			t.Errorf("expected empty message, got: %s", output)
		}
	})

	t.Run("list after adding", func(t *testing.T) {
		runApp(t, workDir, "add", "--company", "Corp1", "--role", "Role1")
		runApp(t, workDir, "add", "--company", "Corp2", "--role", "Role2")

		output, err := runApp(t, workDir, "list")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Corp1") || !strings.Contains(output, "Corp2") {
			t.Errorf("expected both companies in list, got: %s", output)
		}
		if !strings.Contains(output, "Total: 2 applications") {
			t.Errorf("expected total count, got: %s", output)
		}
	})
}

// Helper to extract app ID from add command output
func extractAppID(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "ID:"))
		}
	}
	return ""
}

func TestCLIUpdateWithFlags(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create an application to update
	addOutput, _ := runApp(t, workDir, "add", "--company", "UpdateCorp", "--role", "Engineer")
	appID := extractAppID(addOutput)
	if appID == "" {
		t.Fatal("could not extract app ID")
	}

	t.Run("update status only", func(t *testing.T) {
		output, err := runApp(t, workDir, "update", appID, "--status", "interviewing")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Application updated successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
		if !strings.Contains(output, "interviewing") {
			t.Errorf("expected new status in output, got: %s", output)
		}

		// Verify with show
		showOutput, _ := runApp(t, workDir, "show", appID)
		if !strings.Contains(showOutput, "interviewing") {
			t.Errorf("status not persisted, got: %s", showOutput)
		}
	})

	t.Run("update notes only", func(t *testing.T) {
		output, err := runApp(t, workDir, "update", appID, "--notes", "Had first interview")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Application updated successfully") {
			t.Errorf("expected success message, got: %s", output)
		}

		// Verify with show
		showOutput, _ := runApp(t, workDir, "show", appID)
		if !strings.Contains(showOutput, "Had first interview") {
			t.Errorf("notes not persisted, got: %s", showOutput)
		}
	})

	t.Run("update multiple fields", func(t *testing.T) {
		output, err := runApp(t, workDir, "update", appID,
			"--status", "offer",
			"--notes", "Got the offer!",
			"--jd-url", "https://example.com/job")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "offer") {
			t.Errorf("expected new status in output, got: %s", output)
		}

		// Verify with show
		showOutput, _ := runApp(t, workDir, "show", appID)
		if !strings.Contains(showOutput, "offer") {
			t.Errorf("status not persisted, got: %s", showOutput)
		}
		if !strings.Contains(showOutput, "Got the offer!") {
			t.Errorf("notes not persisted, got: %s", showOutput)
		}
		if !strings.Contains(showOutput, "https://example.com/job") {
			t.Errorf("jd-url not persisted, got: %s", showOutput)
		}
	})

	t.Run("update with jd-file", func(t *testing.T) {
		jdPath := filepath.Join(workDir, "updated-jd.txt")
		os.WriteFile(jdPath, []byte("Updated job description content"), 0644)

		output, err := runApp(t, workDir, "update", appID, "--jd-file", jdPath)
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}

		// Verify with show
		showOutput, _ := runApp(t, workDir, "show", appID)
		if !strings.Contains(showOutput, "Updated job description content") {
			t.Errorf("jd content not persisted, got: %s", showOutput)
		}
	})

	t.Run("update company and role", func(t *testing.T) {
		output, err := runApp(t, workDir, "update", appID,
			"--company", "NewCompany",
			"--role", "Senior Engineer")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "NewCompany") {
			t.Errorf("expected new company in output, got: %s", output)
		}
		if !strings.Contains(output, "Senior Engineer") {
			t.Errorf("expected new role in output, got: %s", output)
		}
	})

	t.Run("update date", func(t *testing.T) {
		output, err := runApp(t, workDir, "update", appID, "--date", "2025-01-01")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "2025-01-01") {
			t.Errorf("expected new date in output, got: %s", output)
		}
	})
}

func TestCLIUpdateErrors(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create an application
	addOutput, _ := runApp(t, workDir, "add", "--company", "ErrorCorp", "--role", "Dev")
	appID := extractAppID(addOutput)

	t.Run("update non-existent app", func(t *testing.T) {
		output, err := runApp(t, workDir, "update", "app-nonexistent", "--status", "offer")
		if err == nil {
			t.Error("expected error for non-existent app")
		}
		if !strings.Contains(output, "Application not found") {
			t.Errorf("expected not found message, got: %s", output)
		}
	})

	t.Run("update with invalid status", func(t *testing.T) {
		output, err := runApp(t, workDir, "update", appID, "--status", "invalid")
		if err == nil {
			t.Error("expected error for invalid status")
		}
		if !strings.Contains(output, "--status must be one of") {
			t.Errorf("expected status error, got: %s", output)
		}
	})

	t.Run("update with invalid date", func(t *testing.T) {
		output, err := runApp(t, workDir, "update", appID, "--date", "not-a-date")
		if err == nil {
			t.Error("expected error for invalid date")
		}
		if !strings.Contains(output, "--date must be in YYYY-MM-DD format") {
			t.Errorf("expected date error, got: %s", output)
		}
	})

	t.Run("update with both jd-content and jd-file", func(t *testing.T) {
		jdPath := filepath.Join(workDir, "jd.txt")
		os.WriteFile(jdPath, []byte("test"), 0644)

		output, err := runApp(t, workDir, "update", appID,
			"--jd-content", "inline",
			"--jd-file", jdPath)
		if err == nil {
			t.Error("expected error for conflicting jd flags")
		}
		if !strings.Contains(output, "--jd-content and --jd-file cannot be used together") {
			t.Errorf("expected conflict error, got: %s", output)
		}
	})
}

func TestCLIShowAndVerifyJDContent(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Add with JD content from file
	jdContent := "Senior Software Engineer\n\nRequirements:\n- 5+ years experience\n- Go expertise"
	jdPath := filepath.Join(workDir, "jd.txt")
	os.WriteFile(jdPath, []byte(jdContent), 0644)

	output, _ := runApp(t, workDir, "add",
		"--company", "TestCorp",
		"--role", "Senior Engineer",
		"--jd-file", jdPath)

	appID := extractAppID(output)
	if appID == "" {
		t.Fatal("could not extract app ID from output")
	}

	// Show and verify JD content is preserved
	showOutput, err := runApp(t, workDir, "show", appID)
	if err != nil {
		t.Fatalf("show command failed: %v\nOutput: %s", err, showOutput)
	}

	if !strings.Contains(showOutput, "Senior Software Engineer") {
		t.Errorf("expected JD title in show output, got: %s", showOutput)
	}
	if !strings.Contains(showOutput, "5+ years experience") {
		t.Errorf("expected JD requirements in show output, got: %s", showOutput)
	}
	if !strings.Contains(showOutput, "Go expertise") {
		t.Errorf("expected JD skills in show output, got: %s", showOutput)
	}
}

func TestCLIHelp(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	output, err := runApp(t, workDir, "help")
	if err != nil {
		t.Fatalf("command failed: %v\nOutput: %s", err, output)
	}

	// Check that new flags are documented
	expectedFlags := []string{
		"--company",
		"--role",
		"--status",
		"--date",
		"--jd-url",
		"--jd-content",
		"--jd-file",
		"--company-url",
		"--resume-path",
		"--notes",
	}

	for _, flag := range expectedFlags {
		if !strings.Contains(output, flag) {
			t.Errorf("expected help to contain %q, got: %s", flag, output)
		}
	}

	// Check that update command examples are documented
	if !strings.Contains(output, "update app-a1b2c3d4 --status") {
		t.Errorf("expected update flag example in help, got: %s", output)
	}
}

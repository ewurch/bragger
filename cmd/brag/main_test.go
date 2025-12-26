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

	// Check that kb command is documented
	if !strings.Contains(output, "kb <subcommand>") {
		t.Errorf("expected kb command in help, got: %s", output)
	}
}

// Helper to extract kb ID from add command output
func extractKBID(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ID:") {
			return strings.TrimSpace(strings.TrimPrefix(line, "ID:"))
		}
	}
	return ""
}

func TestKBHelp(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	output, _ := runApp(t, workDir, "kb")
	// Should show help when no subcommand

	expectedStrings := []string{
		"brag kb <subcommand>",
		"--type",
		"--category",
		"--data",
		"--content",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(output, s) {
			t.Errorf("expected kb help to contain %q", s)
		}
	}
}

func TestKBAddProfile(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	// Create KB file
	kbPath := filepath.Join(workDir, "candidate-kb.jsonl")
	os.WriteFile(kbPath, []byte{}, 0644)

	t.Run("add contact", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "profile",
			"--category", "contact",
			"--source", "test",
			"--data", `{"name":"John Doe","email":"john@example.com"}`)
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Entry added successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
		if !strings.Contains(output, "kb-") {
			t.Errorf("expected kb ID in output, got: %s", output)
		}
	})

	t.Run("add experience", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "profile",
			"--category", "experience",
			"--source", "cv-import",
			"--data", `{"company":"Acme Corp","role":"Senior Engineer","start_date":"2020-01"}`)
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Entry added successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
	})

	t.Run("add skills", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "profile",
			"--category", "skills",
			"--data", `{"languages":["Go","Python"],"frameworks":["React"]}`)
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Entry added successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
	})
}

func TestKBAddContext(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	kbPath := filepath.Join(workDir, "candidate-kb.jsonl")
	os.WriteFile(kbPath, []byte{}, 0644)

	t.Run("add achievement", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "context",
			"--category", "achievement",
			"--source", "user",
			"--content", "Led migration to microservices, reducing latency by 40%")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Entry added successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
		if !strings.Contains(output, "context") {
			t.Errorf("expected type context in output, got: %s", output)
		}
	})
}

func TestKBAddErrors(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	kbPath := filepath.Join(workDir, "candidate-kb.jsonl")
	os.WriteFile(kbPath, []byte{}, 0644)

	t.Run("missing type", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--category", "contact",
			"--data", `{"name":"John","email":"j@e.com"}`)
		if err == nil {
			t.Error("expected error when type is missing")
		}
		if !strings.Contains(output, "--type is required") {
			t.Errorf("expected type error, got: %s", output)
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "invalid",
			"--category", "contact",
			"--data", `{"name":"John"}`)
		if err == nil {
			t.Error("expected error for invalid type")
		}
		if !strings.Contains(output, "--type must be") {
			t.Errorf("expected type error, got: %s", output)
		}
	})

	t.Run("missing category", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "profile",
			"--data", `{"name":"John"}`)
		if err == nil {
			t.Error("expected error when category is missing")
		}
		if !strings.Contains(output, "--category is required") {
			t.Errorf("expected category error, got: %s", output)
		}
	})

	t.Run("invalid profile category", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "profile",
			"--category", "invalid-category",
			"--data", `{"name":"John"}`)
		if err == nil {
			t.Error("expected error for invalid category")
		}
		if !strings.Contains(output, "invalid profile category") {
			t.Errorf("expected category error, got: %s", output)
		}
	})

	t.Run("missing data for profile", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "profile",
			"--category", "contact")
		if err == nil {
			t.Error("expected error when data is missing")
		}
		if !strings.Contains(output, "--data is required") {
			t.Errorf("expected data error, got: %s", output)
		}
	})

	t.Run("missing content for context", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "context",
			"--category", "achievement")
		if err == nil {
			t.Error("expected error when content is missing")
		}
		if !strings.Contains(output, "--content is required") {
			t.Errorf("expected content error, got: %s", output)
		}
	})

	t.Run("invalid JSON data", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "profile",
			"--category", "contact",
			"--data", `{invalid json}`)
		if err == nil {
			t.Error("expected error for invalid JSON")
		}
		if !strings.Contains(output, "invalid JSON") {
			t.Errorf("expected JSON error, got: %s", output)
		}
	})

	t.Run("contact missing required field", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "add",
			"--type", "profile",
			"--category", "contact",
			"--data", `{"name":"John"}`) // missing email
		if err == nil {
			t.Error("expected error for missing email")
		}
		if !strings.Contains(output, "email is required") {
			t.Errorf("expected email required error, got: %s", output)
		}
	})
}

func TestKBShow(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	kbPath := filepath.Join(workDir, "candidate-kb.jsonl")
	os.WriteFile(kbPath, []byte{}, 0644)

	// Add some entries
	runApp(t, workDir, "kb", "add", "--type", "profile", "--category", "contact",
		"--data", `{"name":"John","email":"j@e.com"}`)
	runApp(t, workDir, "kb", "add", "--type", "profile", "--category", "experience",
		"--data", `{"company":"Acme","role":"Dev","start_date":"2020-01"}`)
	runApp(t, workDir, "kb", "add", "--type", "context", "--category", "achievement",
		"--content", "Built something great")

	t.Run("show all", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "show")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Profile Entries") {
			t.Errorf("expected profile section, got: %s", output)
		}
		if !strings.Contains(output, "Context Entries") {
			t.Errorf("expected context section, got: %s", output)
		}
		if !strings.Contains(output, "John") {
			t.Errorf("expected contact name, got: %s", output)
		}
	})

	t.Run("show profile only", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "show", "profile")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Profile Entries") {
			t.Errorf("expected profile section, got: %s", output)
		}
		if strings.Contains(output, "Context Entries") {
			t.Errorf("should not contain context section, got: %s", output)
		}
	})

	t.Run("show context only", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "show", "context")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if strings.Contains(output, "Profile Entries") {
			t.Errorf("should not contain profile section, got: %s", output)
		}
		if !strings.Contains(output, "Context Entries") {
			t.Errorf("expected context section, got: %s", output)
		}
	})
}

func TestKBUpdate(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	kbPath := filepath.Join(workDir, "candidate-kb.jsonl")
	os.WriteFile(kbPath, []byte{}, 0644)

	// Add an entry to update
	addOutput, _ := runApp(t, workDir, "kb", "add",
		"--type", "context",
		"--category", "achievement",
		"--content", "Original content")
	kbID := extractKBID(addOutput)

	t.Run("update content", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "update", kbID,
			"--content", "Updated content")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "Entry updated successfully") {
			t.Errorf("expected success message, got: %s", output)
		}
	})

	t.Run("update non-existent", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "update", "kb-nonexistent",
			"--content", "test")
		if err == nil {
			t.Error("expected error for non-existent entry")
		}
		if !strings.Contains(output, "Entry not found") {
			t.Errorf("expected not found error, got: %s", output)
		}
	})
}

func TestKBContext(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	kbPath := filepath.Join(workDir, "candidate-kb.jsonl")
	os.WriteFile(kbPath, []byte{}, 0644)

	t.Run("empty KB", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "context")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}
		if !strings.Contains(output, "# Candidate Knowledge Base") {
			t.Errorf("expected markdown header, got: %s", output)
		}
		if !strings.Contains(output, "No entries found") {
			t.Errorf("expected empty message, got: %s", output)
		}
	})

	// Add sample data
	runApp(t, workDir, "kb", "add", "--type", "profile", "--category", "contact",
		"--data", `{"name":"Jane Smith","email":"jane@example.com","location":"Berlin","linkedin":"linkedin.com/in/jane"}`)
	runApp(t, workDir, "kb", "add", "--type", "profile", "--category", "experience",
		"--data", `{"company":"TechCorp","role":"Senior Engineer","start_date":"2020-01","end_date":"present","highlights":["Led team of 5","Reduced latency by 40%"]}`)
	runApp(t, workDir, "kb", "add", "--type", "profile", "--category", "education",
		"--data", `{"institution":"MIT","degree":"BS Computer Science","field":"Computer Science","end_date":"2019"}`)
	runApp(t, workDir, "kb", "add", "--type", "profile", "--category", "skills",
		"--data", `{"languages":["Go","Python","TypeScript"],"frameworks":["React","FastAPI"],"databases":["PostgreSQL","Redis"]}`)
	runApp(t, workDir, "kb", "add", "--type", "profile", "--category", "certifications",
		"--data", `{"name":"AWS Solutions Architect","issuer":"Amazon","date":"2023-06"}`)
	runApp(t, workDir, "kb", "add", "--type", "profile", "--category", "languages",
		"--data", `{"language":"German","proficiency":"fluent"}`)
	runApp(t, workDir, "kb", "add", "--type", "context", "--category", "achievement",
		"--source", "user", "--content", "Built real-time notification system handling 10k msgs/sec")
	runApp(t, workDir, "kb", "add", "--type", "context", "--category", "preference",
		"--source", "user", "--content", "Prefer remote-first companies")

	t.Run("full context output", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "context")
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, output)
		}

		// Check markdown structure
		if !strings.Contains(output, "# Candidate Knowledge Base") {
			t.Errorf("expected main header, got: %s", output)
		}

		// Check contact section
		if !strings.Contains(output, "## Contact") {
			t.Errorf("expected Contact section, got: %s", output)
		}
		if !strings.Contains(output, "**Name:** Jane Smith") {
			t.Errorf("expected name in contact, got: %s", output)
		}
		if !strings.Contains(output, "**Email:** jane@example.com") {
			t.Errorf("expected email in contact, got: %s", output)
		}
		if !strings.Contains(output, "**LinkedIn:** linkedin.com/in/jane") {
			t.Errorf("expected linkedin in contact, got: %s", output)
		}

		// Check experience section
		if !strings.Contains(output, "## Experience") {
			t.Errorf("expected Experience section, got: %s", output)
		}
		if !strings.Contains(output, "### Senior Engineer @ TechCorp") {
			t.Errorf("expected experience header, got: %s", output)
		}
		if !strings.Contains(output, "Led team of 5") {
			t.Errorf("expected highlight, got: %s", output)
		}

		// Check education section
		if !strings.Contains(output, "## Education") {
			t.Errorf("expected Education section, got: %s", output)
		}
		if !strings.Contains(output, "BS Computer Science") {
			t.Errorf("expected degree, got: %s", output)
		}

		// Check skills section
		if !strings.Contains(output, "## Skills") {
			t.Errorf("expected Skills section, got: %s", output)
		}
		if !strings.Contains(output, "**Programming Languages:** Go, Python, TypeScript") {
			t.Errorf("expected languages in skills, got: %s", output)
		}

		// Check certifications section
		if !strings.Contains(output, "## Certifications") {
			t.Errorf("expected Certifications section, got: %s", output)
		}
		if !strings.Contains(output, "AWS Solutions Architect") {
			t.Errorf("expected certification, got: %s", output)
		}

		// Check languages section
		if !strings.Contains(output, "## Languages") {
			t.Errorf("expected Languages section, got: %s", output)
		}
		if !strings.Contains(output, "**German** - fluent") {
			t.Errorf("expected language entry, got: %s", output)
		}

		// Check context entries section
		if !strings.Contains(output, "## Context Entries") {
			t.Errorf("expected Context Entries section, got: %s", output)
		}
		if !strings.Contains(output, "### Achievement") {
			t.Errorf("expected Achievement category, got: %s", output)
		}
		if !strings.Contains(output, "10k msgs/sec") {
			t.Errorf("expected achievement content, got: %s", output)
		}

		// Check entry IDs are included
		if !strings.Contains(output, "Entry ID: kb-") {
			t.Errorf("expected entry IDs in output, got: %s", output)
		}
		if !strings.Contains(output, "ID: kb-") {
			t.Errorf("expected IDs for context entries, got: %s", output)
		}
	})
}

func TestKBRemove(t *testing.T) {
	workDir, cleanup := setupIntegrationTest(t)
	defer cleanup()

	kbPath := filepath.Join(workDir, "candidate-kb.jsonl")
	os.WriteFile(kbPath, []byte{}, 0644)

	// Add an entry to remove
	addOutput, _ := runApp(t, workDir, "kb", "add",
		"--type", "context",
		"--category", "test",
		"--content", "To be removed")
	kbID := extractKBID(addOutput)

	t.Run("remove with confirmation", func(t *testing.T) {
		// Use stdin to provide "y" confirmation
		cmd := exec.Command(binaryPath, "kb", "remove", kbID)
		cmd.Dir = workDir
		cmd.Stdin = strings.NewReader("y\n")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("command failed: %v\nOutput: %s", err, string(output))
		}
		if !strings.Contains(string(output), "Entry removed") {
			t.Errorf("expected removed message, got: %s", string(output))
		}
	})

	t.Run("remove non-existent", func(t *testing.T) {
		output, err := runApp(t, workDir, "kb", "remove", "kb-nonexistent")
		if err == nil {
			t.Error("expected error for non-existent entry")
		}
		if !strings.Contains(output, "Entry not found") {
			t.Errorf("expected not found error, got: %s", output)
		}
	})
}

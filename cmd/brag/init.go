package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ewurch/brag/templates"
)

func cmdInit() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Check if already initialized
	bragDir := filepath.Join(cwd, ".brag")
	if _, err := os.Stat(bragDir); err == nil {
		fmt.Println("Workspace already initialized. Use 'brag upgrade' to update to the latest version.")
		return
	}

	fmt.Printf("Initializing Brag workspace in %s\n\n", cwd)

	var created []string
	var skipped []string

	// Create .brag directory with version file
	if err := os.MkdirAll(bragDir, 0755); err != nil {
		fmt.Printf("Error creating .brag directory: %v\n", err)
		os.Exit(1)
	}
	versionFile := filepath.Join(bragDir, "version")
	if err := os.WriteFile(versionFile, []byte(templates.Version), 0644); err != nil {
		fmt.Printf("Error creating version file: %v\n", err)
		os.Exit(1)
	}
	created = append(created, ".brag/version")

	// Create .claude/skills directories and copy skill files
	skillDirs := []string{
		"resume-builder",
		"cover-letter",
		"candidate-kb",
		"applications",
	}

	for _, skill := range skillDirs {
		skillDir := filepath.Join(cwd, ".claude", "skills", skill)
		if err := os.MkdirAll(skillDir, 0755); err != nil {
			fmt.Printf("Error creating skill directory %s: %v\n", skill, err)
			os.Exit(1)
		}

		skillFile := filepath.Join(skillDir, "SKILL.md")
		if _, err := os.Stat(skillFile); err == nil {
			skipped = append(skipped, fmt.Sprintf(".claude/skills/%s/SKILL.md", skill))
			continue
		}

		content, err := templates.Files.ReadFile(fmt.Sprintf("skills/%s/SKILL.md", skill))
		if err != nil {
			fmt.Printf("Error reading skill template %s: %v\n", skill, err)
			os.Exit(1)
		}

		if err := os.WriteFile(skillFile, content, 0644); err != nil {
			fmt.Printf("Error writing skill file %s: %v\n", skill, err)
			os.Exit(1)
		}
		created = append(created, fmt.Sprintf(".claude/skills/%s/SKILL.md", skill))
	}

	// Copy AGENTS.md
	if err := copyTemplateFile("AGENTS.md", filepath.Join(cwd, "AGENTS.md"), &created, &skipped); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Copy package.json
	if err := copyTemplateFile("package.json", filepath.Join(cwd, "package.json"), &created, &skipped); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Create scripts directory and copy html-to-pdf.js
	scriptsDir := filepath.Join(cwd, "scripts")
	if err := os.MkdirAll(scriptsDir, 0755); err != nil {
		fmt.Printf("Error creating scripts directory: %v\n", err)
		os.Exit(1)
	}

	scriptDest := filepath.Join(cwd, "scripts", "html-to-pdf.js")
	if _, err := os.Stat(scriptDest); os.IsNotExist(err) {
		content, err := templates.Files.ReadFile("scripts/html-to-pdf.js")
		if err != nil {
			fmt.Printf("Error reading script template: %v\n", err)
			os.Exit(1)
		}
		if err := os.WriteFile(scriptDest, content, 0644); err != nil {
			fmt.Printf("Error writing script: %v\n", err)
			os.Exit(1)
		}
		created = append(created, "scripts/html-to-pdf.js")
	} else {
		skipped = append(skipped, "scripts/html-to-pdf.js")
	}

	// Create outputs directory
	outputsDir := filepath.Join(cwd, "outputs")
	if _, err := os.Stat(outputsDir); os.IsNotExist(err) {
		if err := os.MkdirAll(outputsDir, 0755); err != nil {
			fmt.Printf("Error creating outputs directory: %v\n", err)
			os.Exit(1)
		}
		// Create .gitkeep to ensure directory is tracked
		gitkeep := filepath.Join(outputsDir, ".gitkeep")
		os.WriteFile(gitkeep, []byte{}, 0644)
		created = append(created, "outputs/")
	} else {
		skipped = append(skipped, "outputs/")
	}

	// Create empty data files
	dataFiles := []string{"applications.jsonl", "candidate-kb.jsonl"}
	for _, file := range dataFiles {
		filePath := filepath.Join(cwd, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
				fmt.Printf("Error creating %s: %v\n", file, err)
				os.Exit(1)
			}
			created = append(created, file)
		} else {
			skipped = append(skipped, file)
		}
	}

	// Print results
	if len(created) > 0 {
		fmt.Println("Created:")
		for _, f := range created {
			fmt.Printf("  %s\n", f)
		}
	}

	if len(skipped) > 0 {
		fmt.Println("\nSkipped (already exist):")
		for _, f := range skipped {
			fmt.Printf("  %s\n", f)
		}
	}

	fmt.Println("\nNext steps:")
	fmt.Println("  1. Run 'npm install' to enable PDF generation")
	fmt.Println("  2. Run 'brag kb add' to add your profile information")
	fmt.Println("  3. Run 'brag add' to track your first application")
}

func copyTemplateFile(templatePath, destPath string, created, skipped *[]string) error {
	if _, err := os.Stat(destPath); err == nil {
		*skipped = append(*skipped, filepath.Base(destPath))
		return nil
	}

	content, err := templates.Files.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("reading template %s: %v", templatePath, err)
	}

	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return fmt.Errorf("writing %s: %v", destPath, err)
	}

	*created = append(*created, filepath.Base(destPath))
	return nil
}

func cmdUpgrade() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Check if workspace is initialized
	bragDir := filepath.Join(cwd, ".brag")
	versionFile := filepath.Join(bragDir, "version")

	if _, err := os.Stat(bragDir); os.IsNotExist(err) {
		fmt.Println("Not a Brag workspace. Run 'brag init' first.")
		os.Exit(1)
	}

	// Read current version
	currentVersion := "unknown"
	if versionBytes, err := os.ReadFile(versionFile); err == nil {
		currentVersion = string(versionBytes)
	}

	if currentVersion == templates.Version {
		fmt.Printf("Workspace is already at version %s (latest)\n", templates.Version)
		return
	}

	fmt.Printf("Upgrading workspace from v%s to v%s\n\n", currentVersion, templates.Version)

	// Create backup directory
	backupDir := filepath.Join(bragDir, "backup", currentVersion)
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		fmt.Printf("Error creating backup directory: %v\n", err)
		os.Exit(1)
	}

	var updated []string

	// Backup and update skill files
	skillDirs := []string{
		"resume-builder",
		"cover-letter",
		"candidate-kb",
		"applications",
	}

	for _, skill := range skillDirs {
		skillFile := filepath.Join(cwd, ".claude", "skills", skill, "SKILL.md")
		backupFile := filepath.Join(backupDir, fmt.Sprintf("%s-SKILL.md", skill))

		// Backup existing file if it exists
		if _, err := os.Stat(skillFile); err == nil {
			if content, err := os.ReadFile(skillFile); err == nil {
				os.WriteFile(backupFile, content, 0644)
			}
		}

		// Write new version
		content, err := templates.Files.ReadFile(fmt.Sprintf("skills/%s/SKILL.md", skill))
		if err != nil {
			fmt.Printf("Error reading skill template %s: %v\n", skill, err)
			continue
		}

		// Ensure directory exists
		os.MkdirAll(filepath.Dir(skillFile), 0755)

		if err := os.WriteFile(skillFile, content, 0644); err != nil {
			fmt.Printf("Error writing skill file %s: %v\n", skill, err)
			continue
		}
		updated = append(updated, fmt.Sprintf(".claude/skills/%s/SKILL.md", skill))
	}

	// Backup and update AGENTS.md
	agentsFile := filepath.Join(cwd, "AGENTS.md")
	if _, err := os.Stat(agentsFile); err == nil {
		if content, err := os.ReadFile(agentsFile); err == nil {
			os.WriteFile(filepath.Join(backupDir, "AGENTS.md"), content, 0644)
		}
	}

	if content, err := templates.Files.ReadFile("AGENTS.md"); err == nil {
		if err := os.WriteFile(agentsFile, content, 0644); err == nil {
			updated = append(updated, "AGENTS.md")
		}
	}

	// Update version file
	os.WriteFile(versionFile, []byte(templates.Version), 0644)

	// Print results
	fmt.Printf("Backed up to: %s\n\n", backupDir)

	if len(updated) > 0 {
		fmt.Println("Updated:")
		for _, f := range updated {
			fmt.Printf("  %s\n", f)
		}
	}

	fmt.Printf("\nUpgrade to v%s complete!\n", templates.Version)
}

// isWorkspaceInitialized checks if the current directory is a Brag workspace
func isWorkspaceInitialized() bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}

	// Check for .brag directory or applications.jsonl
	bragDir := filepath.Join(cwd, ".brag")
	if _, err := os.Stat(bragDir); err == nil {
		return true
	}

	// Also check for applications.jsonl for backwards compatibility
	appsFile := filepath.Join(cwd, "applications.jsonl")
	if _, err := os.Stat(appsFile); err == nil {
		return true
	}

	return false
}

// Unused variable fix - this function uses fs package
var _ fs.FS = templates.Files

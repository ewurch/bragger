package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/ewurch/brag/internal/models"
	"github.com/ewurch/brag/internal/storage"
)

func printKBUsage() {
	fmt.Println(`Candidate Knowledge Base - Manage your profile and contextual information

Usage:
  brag kb <subcommand> [arguments]

Subcommands:
  show [profile|context]   Show knowledge base entries (all, profile only, or context only)
  context                  Export full KB in LLM-friendly markdown format
  add                      Add a new KB entry
  update <id>              Update an existing KB entry
  remove <id>              Remove a KB entry

Flags for add/update:
  --type         Entry type: "profile" or "context" (required for add)
  --category     Category (required). Profile: contact, experience, education, skills, certifications, languages
  --data         JSON data for profile entries (e.g., '{"name":"John","email":"john@example.com"}')
  --content      Text content for context entries
  --source       Source of information (e.g., "cv-import", "user", "app-xxx")

Profile Categories & Required Fields:
  contact:        name, email (optional: phone, location, linkedin, github, website)
  experience:     company, role, start_date (optional: end_date, location, description, highlights)
  education:      institution, degree (optional: field, start_date, end_date, gpa)
  skills:         (optional: languages, frameworks, tools, databases, cloud, other)
  certifications: name (optional: issuer, date, expiry_date, credential_id)
  languages:      language (optional: proficiency)

Examples:
  brag kb show                                    # Show all KB entries
  brag kb show profile                            # Show profile entries only
  brag kb show context                            # Show context entries only
  brag kb context                                 # Export full KB in markdown (for LLM context)

  # Add contact info
  brag kb add --type profile --category contact --source "cv-import" \
    --data '{"name":"John Doe","email":"john@example.com","location":"New York"}'

  # Add work experience
  brag kb add --type profile --category experience --source "cv-import" \
    --data '{"company":"Acme Corp","role":"Senior Engineer","start_date":"2020-01","end_date":"present"}'

  # Add contextual information
  brag kb add --type context --category achievement --source "user" \
    --content "Led migration to microservices, reducing latency by 40%"

  brag kb update kb-a1b2c3d4 --content "Updated achievement description"
  brag kb remove kb-a1b2c3d4`)
}

// kbFlags holds flags for kb add/update commands
type kbFlags struct {
	entryType string
	category  string
	data      string
	content   string
	source    string
}

func registerKBFlags(fs *flag.FlagSet) *kbFlags {
	f := &kbFlags{}
	fs.StringVar(&f.entryType, "type", "", "Entry type: profile or context")
	fs.StringVar(&f.category, "category", "", "Category (e.g., contact, experience, skills, achievement)")
	fs.StringVar(&f.data, "data", "", "JSON data for profile entries")
	fs.StringVar(&f.content, "content", "", "Text content for context entries")
	fs.StringVar(&f.source, "source", "", "Source of information")
	return f
}

func (f *kbFlags) hasAnyFlag() bool {
	return f.entryType != "" || f.category != "" || f.data != "" || f.content != "" || f.source != ""
}

func cmdKB(store *storage.KBStorage, subcommand string, args []string) {
	switch subcommand {
	case "show":
		cmdKBShow(store, args)
	case "context":
		cmdKBContext(store)
	case "add":
		cmdKBAdd(store, args)
	case "update":
		if len(args) < 1 {
			fmt.Println("Usage: brag kb update <id> [flags]")
			os.Exit(1)
		}
		cmdKBUpdate(store, args[0], args[1:])
	case "remove":
		if len(args) < 1 {
			fmt.Println("Usage: brag kb remove <id>")
			os.Exit(1)
		}
		cmdKBRemove(store, args[0])
	default:
		fmt.Printf("Unknown kb subcommand: %s\n", subcommand)
		printKBUsage()
		os.Exit(1)
	}
}

func cmdKBShow(store *storage.KBStorage, args []string) {
	filter := ""
	if len(args) > 0 {
		filter = args[0]
	}

	var entries []*models.KBEntry
	var err error

	switch filter {
	case "profile":
		entries, err = store.GetProfile()
	case "context":
		entries, err = store.GetContext()
	case "":
		entries, err = store.Load()
	default:
		// Treat as category filter
		entries, err = store.GetByCategory(filter)
	}

	if err != nil {
		fmt.Printf("Error loading knowledge base: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("No entries found.")
		return
	}

	// Group entries by type for display
	var profileEntries, contextEntries []*models.KBEntry
	for _, e := range entries {
		if e.Type == models.KBTypeProfile {
			profileEntries = append(profileEntries, e)
		} else {
			contextEntries = append(contextEntries, e)
		}
	}

	if len(profileEntries) > 0 && (filter == "" || filter == "profile") {
		fmt.Println("=== Profile Entries ===")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tCATEGORY\tSOURCE\tSUMMARY")
		fmt.Fprintln(w, "--\t--------\t------\t-------")
		for _, e := range profileEntries {
			summary := summarizeProfileData(e)
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.ID, e.Category, e.Source, truncate(summary, 50))
		}
		w.Flush()
		fmt.Printf("\nProfile entries: %d\n", len(profileEntries))
	}

	if len(contextEntries) > 0 && (filter == "" || filter == "context") {
		if len(profileEntries) > 0 {
			fmt.Println()
		}
		fmt.Println("=== Context Entries ===")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tCATEGORY\tSOURCE\tCONTENT")
		fmt.Fprintln(w, "--\t--------\t------\t-------")
		for _, e := range contextEntries {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", e.ID, e.Category, e.Source, truncate(e.Content, 50))
		}
		w.Flush()
		fmt.Printf("\nContext entries: %d\n", len(contextEntries))
	}
}

func cmdKBContext(store *storage.KBStorage) {
	entries, err := store.Load()
	if err != nil {
		fmt.Printf("Error loading knowledge base: %v\n", err)
		os.Exit(1)
	}

	if len(entries) == 0 {
		fmt.Println("# Candidate Knowledge Base")
		fmt.Println("\n*No entries found. Use `brag kb add` to populate the knowledge base.*")
		return
	}

	// Group entries by type and category
	var contact *models.KBEntry
	var experiences []*models.KBEntry
	var education []*models.KBEntry
	var skills *models.KBEntry
	var certifications []*models.KBEntry
	var languages []*models.KBEntry
	var contextEntries []*models.KBEntry

	for _, e := range entries {
		if e.Type == models.KBTypeProfile {
			switch models.ProfileCategory(e.Category) {
			case models.CategoryContact:
				contact = e
			case models.CategoryExperience:
				experiences = append(experiences, e)
			case models.CategoryEducation:
				education = append(education, e)
			case models.CategorySkills:
				skills = e
			case models.CategoryCertifications:
				certifications = append(certifications, e)
			case models.CategoryLanguages:
				languages = append(languages, e)
			}
		} else {
			contextEntries = append(contextEntries, e)
		}
	}

	// Output in LLM-friendly markdown format
	fmt.Println("# Candidate Knowledge Base")
	fmt.Println()

	// Contact
	if contact != nil {
		fmt.Println("## Contact")
		fmt.Printf("*Entry ID: %s*\n\n", contact.ID)
		printContactMarkdown(contact)
		fmt.Println()
	}

	// Experience
	if len(experiences) > 0 {
		fmt.Println("## Experience")
		fmt.Println()
		for _, e := range experiences {
			printExperienceMarkdown(e)
			fmt.Println()
		}
	}

	// Education
	if len(education) > 0 {
		fmt.Println("## Education")
		fmt.Println()
		for _, e := range education {
			printEducationMarkdown(e)
			fmt.Println()
		}
	}

	// Skills
	if skills != nil {
		fmt.Println("## Skills")
		fmt.Printf("*Entry ID: %s*\n\n", skills.ID)
		printSkillsMarkdown(skills)
		fmt.Println()
	}

	// Certifications
	if len(certifications) > 0 {
		fmt.Println("## Certifications")
		fmt.Println()
		for _, e := range certifications {
			printCertificationMarkdown(e)
		}
		fmt.Println()
	}

	// Languages
	if len(languages) > 0 {
		fmt.Println("## Languages")
		fmt.Println()
		for _, e := range languages {
			printLanguageMarkdown(e)
		}
		fmt.Println()
	}

	// Context entries grouped by category
	if len(contextEntries) > 0 {
		fmt.Println("## Context Entries")
		fmt.Println()

		// Group by category
		categories := make(map[string][]*models.KBEntry)
		for _, e := range contextEntries {
			categories[e.Category] = append(categories[e.Category], e)
		}

		for category, entries := range categories {
			fmt.Printf("### %s\n\n", capitalizeFirst(category))
			for _, e := range entries {
				fmt.Printf("- %s *(ID: %s, source: %s)*\n", e.Content, e.ID, e.Source)
			}
			fmt.Println()
		}
	}
}

func printContactMarkdown(e *models.KBEntry) {
	dataBytes, err := json.Marshal(e.Data)
	if err != nil {
		return
	}
	var c models.ContactData
	if json.Unmarshal(dataBytes, &c) != nil {
		return
	}

	fmt.Printf("- **Name:** %s\n", c.Name)
	fmt.Printf("- **Email:** %s\n", c.Email)
	if c.Phone != "" {
		fmt.Printf("- **Phone:** %s\n", c.Phone)
	}
	if c.Location != "" {
		fmt.Printf("- **Location:** %s\n", c.Location)
	}
	if c.LinkedIn != "" {
		fmt.Printf("- **LinkedIn:** %s\n", c.LinkedIn)
	}
	if c.GitHub != "" {
		fmt.Printf("- **GitHub:** %s\n", c.GitHub)
	}
	if c.Website != "" {
		fmt.Printf("- **Website:** %s\n", c.Website)
	}
}

func printExperienceMarkdown(e *models.KBEntry) {
	dataBytes, err := json.Marshal(e.Data)
	if err != nil {
		return
	}
	var exp models.ExperienceEntry
	if json.Unmarshal(dataBytes, &exp) != nil {
		return
	}

	endDate := exp.EndDate
	if endDate == "" {
		endDate = "present"
	}

	fmt.Printf("### %s @ %s\n", exp.Role, exp.Company)
	fmt.Printf("*Entry ID: %s*\n\n", e.ID)
	fmt.Printf("**Period:** %s - %s\n", exp.StartDate, endDate)
	if exp.Location != "" {
		fmt.Printf("**Location:** %s\n", exp.Location)
	}
	if exp.Description != "" {
		fmt.Printf("\n%s\n", exp.Description)
	}
	if len(exp.Highlights) > 0 {
		fmt.Println("\n**Highlights:**")
		for _, h := range exp.Highlights {
			fmt.Printf("- %s\n", h)
		}
	}
}

func printEducationMarkdown(e *models.KBEntry) {
	dataBytes, err := json.Marshal(e.Data)
	if err != nil {
		return
	}
	var edu models.EducationEntry
	if json.Unmarshal(dataBytes, &edu) != nil {
		return
	}

	fmt.Printf("### %s", edu.Degree)
	if edu.Field != "" {
		fmt.Printf(" in %s", edu.Field)
	}
	fmt.Printf(" - %s\n", edu.Institution)
	fmt.Printf("*Entry ID: %s*\n", e.ID)

	if edu.StartDate != "" || edu.EndDate != "" {
		fmt.Printf("**Period:** %s - %s\n", edu.StartDate, edu.EndDate)
	}
	if edu.GPA != "" {
		fmt.Printf("**GPA:** %s\n", edu.GPA)
	}
}

func printSkillsMarkdown(e *models.KBEntry) {
	dataBytes, err := json.Marshal(e.Data)
	if err != nil {
		return
	}
	var s models.SkillsData
	if json.Unmarshal(dataBytes, &s) != nil {
		return
	}

	if len(s.Languages) > 0 {
		fmt.Printf("- **Programming Languages:** %s\n", joinStrings(s.Languages))
	}
	if len(s.Frameworks) > 0 {
		fmt.Printf("- **Frameworks:** %s\n", joinStrings(s.Frameworks))
	}
	if len(s.Tools) > 0 {
		fmt.Printf("- **Tools:** %s\n", joinStrings(s.Tools))
	}
	if len(s.Databases) > 0 {
		fmt.Printf("- **Databases:** %s\n", joinStrings(s.Databases))
	}
	if len(s.Cloud) > 0 {
		fmt.Printf("- **Cloud:** %s\n", joinStrings(s.Cloud))
	}
	if len(s.Other) > 0 {
		fmt.Printf("- **Other:** %s\n", joinStrings(s.Other))
	}
}

func printCertificationMarkdown(e *models.KBEntry) {
	dataBytes, err := json.Marshal(e.Data)
	if err != nil {
		return
	}
	var c models.CertificationEntry
	if json.Unmarshal(dataBytes, &c) != nil {
		return
	}

	fmt.Printf("- **%s**", c.Name)
	if c.Issuer != "" {
		fmt.Printf(" by %s", c.Issuer)
	}
	if c.Date != "" {
		fmt.Printf(" (%s)", c.Date)
	}
	fmt.Printf(" *(ID: %s)*\n", e.ID)
}

func printLanguageMarkdown(e *models.KBEntry) {
	dataBytes, err := json.Marshal(e.Data)
	if err != nil {
		return
	}
	var l models.LanguageEntry
	if json.Unmarshal(dataBytes, &l) != nil {
		return
	}

	fmt.Printf("- **%s**", l.Language)
	if l.Proficiency != "" {
		fmt.Printf(" - %s", l.Proficiency)
	}
	fmt.Printf(" *(ID: %s)*\n", e.ID)
}

func joinStrings(s []string) string {
	result := ""
	for i, str := range s {
		if i > 0 {
			result += ", "
		}
		result += str
	}
	return result
}

func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return string(s[0]-32) + s[1:]
}

func summarizeProfileData(e *models.KBEntry) string {
	if e.Data == nil {
		return ""
	}

	// Convert data back to JSON for display
	dataBytes, err := json.Marshal(e.Data)
	if err != nil {
		return ""
	}

	// Parse based on category to create summary
	switch models.ProfileCategory(e.Category) {
	case models.CategoryContact:
		var c models.ContactData
		if json.Unmarshal(dataBytes, &c) == nil {
			return fmt.Sprintf("%s <%s>", c.Name, c.Email)
		}
	case models.CategoryExperience:
		var exp models.ExperienceEntry
		if json.Unmarshal(dataBytes, &exp) == nil {
			return fmt.Sprintf("%s @ %s (%s - %s)", exp.Role, exp.Company, exp.StartDate, exp.EndDate)
		}
	case models.CategoryEducation:
		var edu models.EducationEntry
		if json.Unmarshal(dataBytes, &edu) == nil {
			return fmt.Sprintf("%s, %s", edu.Degree, edu.Institution)
		}
	case models.CategorySkills:
		var s models.SkillsData
		if json.Unmarshal(dataBytes, &s) == nil {
			count := len(s.Languages) + len(s.Frameworks) + len(s.Tools) + len(s.Databases) + len(s.Cloud) + len(s.Other)
			return fmt.Sprintf("%d skills across categories", count)
		}
	case models.CategoryCertifications:
		var c models.CertificationEntry
		if json.Unmarshal(dataBytes, &c) == nil {
			return c.Name
		}
	case models.CategoryLanguages:
		var l models.LanguageEntry
		if json.Unmarshal(dataBytes, &l) == nil {
			if l.Proficiency != "" {
				return fmt.Sprintf("%s (%s)", l.Language, l.Proficiency)
			}
			return l.Language
		}
	}

	return string(dataBytes)
}

func cmdKBAdd(store *storage.KBStorage, args []string) {
	fs := flag.NewFlagSet("kb add", flag.ExitOnError)
	flags := registerKBFlags(fs)
	fs.Parse(args)

	if !flags.hasAnyFlag() {
		fmt.Println("Error: flags required. Use --type, --category, and --data or --content")
		printKBUsage()
		os.Exit(1)
	}

	// Validate type
	if flags.entryType == "" {
		fmt.Println("Error: --type is required (profile or context)")
		os.Exit(1)
	}

	entryType := models.KBEntryType(flags.entryType)
	if !entryType.IsValid() {
		fmt.Println("Error: --type must be 'profile' or 'context'")
		os.Exit(1)
	}

	// Validate category
	if flags.category == "" {
		fmt.Println("Error: --category is required")
		os.Exit(1)
	}

	var entry *models.KBEntry

	if entryType == models.KBTypeProfile {
		// Validate profile category
		cat := models.ProfileCategory(flags.category)
		if !cat.IsValid() {
			fmt.Println("Error: invalid profile category. Must be one of: contact, experience, education, skills, certifications, languages")
			os.Exit(1)
		}

		if flags.data == "" {
			fmt.Println("Error: --data is required for profile entries")
			os.Exit(1)
		}

		// Parse and validate data based on category
		data, err := parseAndValidateProfileData(cat, flags.data)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		entry = models.NewProfileEntry(cat, data, flags.source)
	} else {
		// Context entry
		if flags.content == "" {
			fmt.Println("Error: --content is required for context entries")
			os.Exit(1)
		}

		entry = models.NewContextEntry(flags.category, flags.content, flags.source)
	}

	if err := store.Add(entry); err != nil {
		fmt.Printf("Error adding entry: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Entry added successfully!")
	fmt.Printf("ID: %s\n", entry.ID)
	fmt.Printf("Type: %s\n", entry.Type)
	fmt.Printf("Category: %s\n", entry.Category)
}

func parseAndValidateProfileData(category models.ProfileCategory, jsonData string) (any, error) {
	switch category {
	case models.CategoryContact:
		var data models.ContactData
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			return nil, fmt.Errorf("invalid JSON for contact: %v", err)
		}
		if err := data.Validate(); err != nil {
			return nil, err
		}
		return data, nil

	case models.CategoryExperience:
		var data models.ExperienceEntry
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			return nil, fmt.Errorf("invalid JSON for experience: %v", err)
		}
		if err := data.Validate(); err != nil {
			return nil, err
		}
		return data, nil

	case models.CategoryEducation:
		var data models.EducationEntry
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			return nil, fmt.Errorf("invalid JSON for education: %v", err)
		}
		if err := data.Validate(); err != nil {
			return nil, err
		}
		return data, nil

	case models.CategorySkills:
		var data models.SkillsData
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			return nil, fmt.Errorf("invalid JSON for skills: %v", err)
		}
		if err := data.Validate(); err != nil {
			return nil, err
		}
		return data, nil

	case models.CategoryCertifications:
		var data models.CertificationEntry
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			return nil, fmt.Errorf("invalid JSON for certification: %v", err)
		}
		if err := data.Validate(); err != nil {
			return nil, err
		}
		return data, nil

	case models.CategoryLanguages:
		var data models.LanguageEntry
		if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
			return nil, fmt.Errorf("invalid JSON for language: %v", err)
		}
		if err := data.Validate(); err != nil {
			return nil, err
		}
		return data, nil

	default:
		return nil, fmt.Errorf("unknown category: %s", category)
	}
}

func cmdKBUpdate(store *storage.KBStorage, id string, args []string) {
	entry, err := store.Get(id)
	if err != nil {
		fmt.Printf("Entry not found: %s\n", id)
		os.Exit(1)
	}

	fs := flag.NewFlagSet("kb update", flag.ExitOnError)
	flags := registerKBFlags(fs)
	fs.Parse(args)

	if !flags.hasAnyFlag() {
		fmt.Println("Error: at least one flag required for update")
		os.Exit(1)
	}

	// Update fields based on entry type
	if entry.Type == models.KBTypeProfile {
		if flags.data != "" {
			cat := models.ProfileCategory(entry.Category)
			data, err := parseAndValidateProfileData(cat, flags.data)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			entry.Data = data
		}
	} else {
		if flags.content != "" {
			entry.Content = flags.content
		}
	}

	if flags.category != "" {
		entry.Category = flags.category
	}
	if flags.source != "" {
		entry.Source = flags.source
	}

	err = store.Update(id, func(e *models.KBEntry) {
		e.Category = entry.Category
		e.Data = entry.Data
		e.Content = entry.Content
		e.Source = entry.Source
	})
	if err != nil {
		fmt.Printf("Error updating entry: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Entry updated successfully!")
	fmt.Printf("ID: %s\n", entry.ID)
	fmt.Printf("Type: %s\n", entry.Type)
	fmt.Printf("Category: %s\n", entry.Category)
}

func cmdKBRemove(store *storage.KBStorage, id string) {
	entry, err := store.Get(id)
	if err != nil {
		fmt.Printf("Entry not found: %s\n", id)
		os.Exit(1)
	}

	// Show what will be deleted
	fmt.Printf("Remove entry %s (%s/%s)? (y/N): ", entry.ID, entry.Type, entry.Category)
	var confirm string
	fmt.Scanln(&confirm)

	if confirm != "y" && confirm != "Y" && confirm != "yes" {
		fmt.Println("Cancelled.")
		return
	}

	if err := store.Remove(id); err != nil {
		fmt.Printf("Error removing entry: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Entry removed.")
}

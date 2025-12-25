package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ewurch/resume-tracker/internal/models"
	"github.com/ewurch/resume-tracker/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	store := storage.New("")

	switch os.Args[1] {
	case "init":
		cmdInit()
	case "add":
		cmdAdd(store, os.Args[2:])
	case "list":
		cmdList(store)
	case "update":
		if len(os.Args) < 3 {
			fmt.Println("Usage: app update <id>")
			os.Exit(1)
		}
		cmdUpdate(store, os.Args[2])
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Usage: app remove <id>")
			os.Exit(1)
		}
		cmdRemove(store, os.Args[2])
	case "show":
		if len(os.Args) < 3 {
			fmt.Println("Usage: app show <id>")
			os.Exit(1)
		}
		cmdShow(store, os.Args[2])
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Application Tracker - Track your job applications

Usage:
  app <command> [arguments]

Commands:
  init             Initialize application tracker (create applications.jsonl)
  add              Add a new application
  list             List all applications
  show <id>        Show details of an application
  update <id>      Update an application (interactive)
  remove <id>      Remove an application
  help             Show this help message

Add command flags (optional - without flags, runs interactively):
  --company        Company name (required with flags)
  --role           Job title (required with flags)
  --jd-url         Job description URL
  --jd-content     Job description text (inline)
  --jd-file        Path to file containing job description
  --company-url    Company website
  --resume-path    Path to resume file
  --notes          Notes about application

Examples:
  app init
  app add                                            # Interactive mode
  app add --company "Acme" --role "Engineer"         # Flag mode (quick)
  app add --company "Acme" --role "Engineer" --jd-url "https://..."
  app add --company "Acme" --role "Engineer" --jd-file ./jd.txt
  app list
  app show app-a1b2c3d4
  app update app-a1b2c3d4
  app remove app-a1b2c3d4`)
}

func cmdInit() {
	filePath := storage.DefaultFilePath

	// Check if file already exists
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("Application tracker already initialized (%s exists)\n", filePath)
		return
	}

	// Create empty file
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating %s: %v\n", filePath, err)
		os.Exit(1)
	}
	file.Close()

	fmt.Printf("Initialized application tracker: %s\n", filePath)
	fmt.Println("\nNext steps:")
	fmt.Println("  app add     - Add your first application")
	fmt.Println("  app list    - View all applications")
}

func cmdAdd(store *storage.Storage, args []string) {
	// Define flags
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	companyFlag := fs.String("company", "", "Company name")
	roleFlag := fs.String("role", "", "Job title")
	jdURLFlag := fs.String("jd-url", "", "Job description URL")
	jdContentFlag := fs.String("jd-content", "", "Job description text (inline)")
	jdFileFlag := fs.String("jd-file", "", "Path to file containing job description")
	companyURLFlag := fs.String("company-url", "", "Company website")
	resumePathFlag := fs.String("resume-path", "", "Path to resume file")
	notesFlag := fs.String("notes", "", "Notes about application")

	fs.Parse(args)

	// Check if any flag is provided (flag mode vs interactive mode)
	flagMode := *companyFlag != "" || *roleFlag != "" || *jdURLFlag != "" ||
		*jdContentFlag != "" || *jdFileFlag != "" ||
		*companyURLFlag != "" || *resumePathFlag != "" || *notesFlag != ""

	var app *models.Application

	if flagMode {
		// Flag mode: validate required fields
		if *companyFlag == "" || *roleFlag == "" {
			fmt.Println("Error: --company and --role are required when using flags")
			os.Exit(1)
		}

		// Validate JD content flags are not both provided
		if *jdContentFlag != "" && *jdFileFlag != "" {
			fmt.Println("Error: --jd-content and --jd-file cannot be used together")
			os.Exit(1)
		}

		app = models.NewApplication(*companyFlag, *roleFlag)
		app.JDURL = *jdURLFlag
		app.CompanyURL = *companyURLFlag
		app.ResumePath = *resumePathFlag
		app.Notes = *notesFlag

		// Handle JD content
		if *jdFileFlag != "" {
			content, err := os.ReadFile(*jdFileFlag)
			if err != nil {
				fmt.Printf("Error reading JD file: %v\n", err)
				os.Exit(1)
			}
			app.JDContent = strings.TrimSpace(string(content))
		} else if *jdContentFlag != "" {
			app.JDContent = *jdContentFlag
		}
	} else {
		// Interactive mode (original behavior)
		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Company: ")
		company, _ := reader.ReadString('\n')
		company = strings.TrimSpace(company)
		if company == "" {
			fmt.Println("Company is required")
			os.Exit(1)
		}

		fmt.Print("Role: ")
		role, _ := reader.ReadString('\n')
		role = strings.TrimSpace(role)
		if role == "" {
			fmt.Println("Role is required")
			os.Exit(1)
		}

		app = models.NewApplication(company, role)

		fmt.Print("JD URL (optional): ")
		jdURL, _ := reader.ReadString('\n')
		app.JDURL = strings.TrimSpace(jdURL)

		fmt.Print("Company URL (optional): ")
		companyURL, _ := reader.ReadString('\n')
		app.CompanyURL = strings.TrimSpace(companyURL)

		fmt.Print("Resume path (optional, e.g., outputs/company_role/resume.html): ")
		resumePath, _ := reader.ReadString('\n')
		app.ResumePath = strings.TrimSpace(resumePath)

		fmt.Println("\nPaste the job description (enter a blank line when done):")
		var jdLines []string
		for {
			line, _ := reader.ReadString('\n')
			if strings.TrimSpace(line) == "" {
				break
			}
			jdLines = append(jdLines, strings.TrimRight(line, "\n"))
		}
		app.JDContent = strings.Join(jdLines, "\n")

		fmt.Print("Notes (optional): ")
		notes, _ := reader.ReadString('\n')
		app.Notes = strings.TrimSpace(notes)
	}

	if err := store.Add(app); err != nil {
		fmt.Printf("Error adding application: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nApplication added successfully!\n")
	fmt.Printf("ID: %s\n", app.ID)
	fmt.Printf("Company: %s\n", app.Company)
	fmt.Printf("Role: %s\n", app.Role)
	fmt.Printf("Status: %s\n", app.Status)
	fmt.Printf("Date Applied: %s\n", app.DateApplied)
}

func cmdList(store *storage.Storage) {
	apps, err := store.Load()
	if err != nil {
		fmt.Printf("Error loading applications: %v\n", err)
		os.Exit(1)
	}

	if len(apps) == 0 {
		fmt.Println("No applications found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tCOMPANY\tROLE\tSTATUS\tDATE APPLIED")
	fmt.Fprintln(w, "--\t-------\t----\t------\t------------")

	for _, app := range apps {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			app.ID,
			truncate(app.Company, 20),
			truncate(app.Role, 25),
			app.Status,
			app.DateApplied,
		)
	}
	w.Flush()

	fmt.Printf("\nTotal: %d applications\n", len(apps))
}

func cmdShow(store *storage.Storage, id string) {
	app, err := store.Get(id)
	if err != nil {
		fmt.Printf("Application not found: %s\n", id)
		os.Exit(1)
	}

	fmt.Printf("ID:           %s\n", app.ID)
	fmt.Printf("Company:      %s\n", app.Company)
	fmt.Printf("Role:         %s\n", app.Role)
	fmt.Printf("Status:       %s\n", app.Status)
	fmt.Printf("Date Applied: %s\n", app.DateApplied)

	if app.JDURL != "" {
		fmt.Printf("JD URL:       %s\n", app.JDURL)
	}
	if app.CompanyURL != "" {
		fmt.Printf("Company URL:  %s\n", app.CompanyURL)
	}
	if app.ResumePath != "" {
		fmt.Printf("Resume Path:  %s\n", app.ResumePath)
	}
	if app.Notes != "" {
		fmt.Printf("Notes:        %s\n", app.Notes)
	}

	fmt.Printf("Created:      %s\n", app.CreatedAt.Format(time.RFC3339))
	fmt.Printf("Updated:      %s\n", app.UpdatedAt.Format(time.RFC3339))

	if app.JDContent != "" {
		fmt.Printf("\n--- Job Description ---\n%s\n", app.JDContent)
	}
}

func cmdUpdate(store *storage.Storage, id string) {
	app, err := store.Get(id)
	if err != nil {
		fmt.Printf("Application not found: %s\n", id)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Current status: %s\n", app.Status)
	fmt.Print("New status (applied/interviewing/rejected/offer) [enter to skip]: ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSpace(status)

	if status != "" {
		newStatus := models.Status(status)
		if !newStatus.IsValid() {
			fmt.Println("Invalid status. Must be: applied, interviewing, rejected, or offer")
			os.Exit(1)
		}

		err = store.Update(id, func(a *models.Application) {
			a.Status = newStatus
		})
		if err != nil {
			fmt.Printf("Error updating application: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Status updated to: %s\n", newStatus)
	}

	fmt.Printf("Current notes: %s\n", app.Notes)
	fmt.Print("New notes [enter to skip]: ")
	notes, _ := reader.ReadString('\n')
	notes = strings.TrimSpace(notes)

	if notes != "" {
		err = store.Update(id, func(a *models.Application) {
			a.Notes = notes
		})
		if err != nil {
			fmt.Printf("Error updating notes: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Notes updated.")
	}

	fmt.Println("Update complete.")
}

func cmdRemove(store *storage.Storage, id string) {
	app, err := store.Get(id)
	if err != nil {
		fmt.Printf("Application not found: %s\n", id)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Remove application for %s at %s? (y/N): ", app.Role, app.Company)
	confirm, _ := reader.ReadString('\n')
	confirm = strings.TrimSpace(strings.ToLower(confirm))

	if confirm != "y" && confirm != "yes" {
		fmt.Println("Cancelled.")
		return
	}

	if err := store.Remove(id); err != nil {
		fmt.Printf("Error removing application: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Application removed.")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

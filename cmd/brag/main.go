package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/ewurch/brag/internal/models"
	"github.com/ewurch/brag/internal/storage"
)

// appFlags holds all the optional flags for add/update commands
type appFlags struct {
	company    string
	role       string
	status     string
	date       string
	jdURL      string
	jdContent  string
	jdFile     string
	companyURL string
	resumePath string
	notes      string
}

// registerAppFlags registers all application flags on a FlagSet
func registerAppFlags(fs *flag.FlagSet) *appFlags {
	f := &appFlags{}
	fs.StringVar(&f.company, "company", "", "Company name")
	fs.StringVar(&f.role, "role", "", "Job title")
	fs.StringVar(&f.status, "status", "", "Application status (applied/interviewing/rejected/offer)")
	fs.StringVar(&f.date, "date", "", "Date applied (YYYY-MM-DD format)")
	fs.StringVar(&f.jdURL, "jd-url", "", "Job description URL")
	fs.StringVar(&f.jdContent, "jd-content", "", "Job description text (inline)")
	fs.StringVar(&f.jdFile, "jd-file", "", "Path to file containing job description")
	fs.StringVar(&f.companyURL, "company-url", "", "Company website")
	fs.StringVar(&f.resumePath, "resume-path", "", "Path to resume file")
	fs.StringVar(&f.notes, "notes", "", "Notes about application")
	return f
}

// hasAnyFlag returns true if any flag was provided
func (f *appFlags) hasAnyFlag() bool {
	return f.company != "" || f.role != "" || f.status != "" || f.date != "" ||
		f.jdURL != "" || f.jdContent != "" || f.jdFile != "" ||
		f.companyURL != "" || f.resumePath != "" || f.notes != ""
}

// validate checks flag values and returns an error message if invalid
func (f *appFlags) validate() string {
	// Validate JD content flags are not both provided
	if f.jdContent != "" && f.jdFile != "" {
		return "Error: --jd-content and --jd-file cannot be used together"
	}

	// Validate status if provided
	if f.status != "" {
		status := models.Status(f.status)
		if !status.IsValid() {
			return "Error: --status must be one of: applied, interviewing, rejected, offer"
		}
	}

	// Validate date format if provided
	if f.date != "" {
		if _, err := time.Parse("2006-01-02", f.date); err != nil {
			return "Error: --date must be in YYYY-MM-DD format"
		}
	}

	return ""
}

// loadJDContent reads JD content from file if jdFile is set
func (f *appFlags) loadJDContent() (string, error) {
	if f.jdFile != "" {
		content, err := os.ReadFile(f.jdFile)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(content)), nil
	}
	return f.jdContent, nil
}

// applyToApplication applies non-empty flag values to an application
func (f *appFlags) applyToApplication(app *models.Application) error {
	if f.company != "" {
		app.Company = f.company
	}
	if f.role != "" {
		app.Role = f.role
	}
	if f.status != "" {
		app.Status = models.Status(f.status)
	}
	if f.date != "" {
		app.DateApplied = f.date
	}
	if f.jdURL != "" {
		app.JDURL = f.jdURL
	}
	if f.companyURL != "" {
		app.CompanyURL = f.companyURL
	}
	if f.resumePath != "" {
		app.ResumePath = f.resumePath
	}
	if f.notes != "" {
		app.Notes = f.notes
	}

	// Handle JD content (from file or inline)
	jdContent, err := f.loadJDContent()
	if err != nil {
		return fmt.Errorf("Error reading JD file: %v", err)
	}
	if jdContent != "" {
		app.JDContent = jdContent
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	store := storage.New("")

	kbStore := storage.NewKBStorage("")

	switch os.Args[1] {
	case "init":
		cmdInit()
	case "add":
		cmdAdd(store, os.Args[2:])
	case "list":
		cmdList(store)
	case "update":
		if len(os.Args) < 3 {
			fmt.Println("Usage: brag update <id> [flags]")
			os.Exit(1)
		}
		cmdUpdate(store, os.Args[2], os.Args[3:])
	case "remove":
		if len(os.Args) < 3 {
			fmt.Println("Usage: brag remove <id>")
			os.Exit(1)
		}
		cmdRemove(store, os.Args[2])
	case "show":
		if len(os.Args) < 3 {
			fmt.Println("Usage: brag show <id>")
			os.Exit(1)
		}
		cmdShow(store, os.Args[2])
	case "upgrade":
		cmdUpgrade()
	case "kb":
		if len(os.Args) < 3 {
			printKBUsage()
			os.Exit(1)
		}
		cmdKB(kbStore, os.Args[2], os.Args[3:])
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`Brag - Job application tracker and AI-powered resume generator

Usage:
  brag <command> [arguments]

Commands:
  init             Initialize a Brag workspace in the current directory
  add              Add a new application
  list             List all applications
  show <id>        Show details of an application
  update <id>      Update an application
  remove <id>      Remove an application
  kb <subcommand>  Manage candidate knowledge base (run 'brag kb' for details)
  upgrade          Upgrade workspace to latest version
  help             Show this help message

Flags for add command (optional - without flags, runs interactively):
  --company        Company name (required with flags)
  --role           Job title (required with flags)
  --status         Status: applied, interviewing, rejected, offer (default: applied)
  --date           Date applied in YYYY-MM-DD format (default: today)
  --jd-url         Job description URL
  --jd-content     Job description text (inline)
  --jd-file        Path to file containing job description
  --company-url    Company website
  --resume-path    Path to resume file
  --notes          Notes about application

Flags for update command (optional - without flags, runs interactively):
  All flags from add command are supported. Only provided flags will be updated.

Examples:
  brag init
  brag add                                            # Interactive mode
  brag add --company "Acme" --role "Engineer"         # Flag mode (quick)
  brag add --company "Acme" --role "Engineer" --jd-url "https://..."
  brag add --company "Acme" --role "Engineer" --jd-file ./jd.txt
  brag add --company "Old" --role "Dev" --date "2025-01-15" --status "interviewing"
  brag list
  brag show app-a1b2c3d4
  brag update app-a1b2c3d4                            # Interactive mode
  brag update app-a1b2c3d4 --status "interviewing"    # Flag mode (quick)
  brag update app-a1b2c3d4 --status "offer" --notes "Accepted!"
  brag remove app-a1b2c3d4
  brag kb show                                        # Show knowledge base
  brag kb add --type profile --category contact --data '{"name":"John"}'`)
}

func cmdAdd(store *storage.Storage, args []string) {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	flags := registerAppFlags(fs)
	fs.Parse(args)

	var app *models.Application

	if flags.hasAnyFlag() {
		// Flag mode: validate required fields
		if flags.company == "" || flags.role == "" {
			fmt.Println("Error: --company and --role are required when using flags")
			os.Exit(1)
		}

		// Validate flags
		if errMsg := flags.validate(); errMsg != "" {
			fmt.Println(errMsg)
			os.Exit(1)
		}

		app = models.NewApplication(flags.company, flags.role)
		if err := flags.applyToApplication(app); err != nil {
			fmt.Println(err)
			os.Exit(1)
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

func cmdUpdate(store *storage.Storage, id string, args []string) {
	app, err := store.Get(id)
	if err != nil {
		fmt.Printf("Application not found: %s\n", id)
		os.Exit(1)
	}

	fs := flag.NewFlagSet("update", flag.ExitOnError)
	flags := registerAppFlags(fs)
	fs.Parse(args)

	if flags.hasAnyFlag() {
		// Flag mode: validate and apply
		if errMsg := flags.validate(); errMsg != "" {
			fmt.Println(errMsg)
			os.Exit(1)
		}

		if err := flags.applyToApplication(app); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Save the updated application
		err = store.Update(id, func(a *models.Application) {
			a.Company = app.Company
			a.Role = app.Role
			a.Status = app.Status
			a.DateApplied = app.DateApplied
			a.JDURL = app.JDURL
			a.JDContent = app.JDContent
			a.CompanyURL = app.CompanyURL
			a.ResumePath = app.ResumePath
			a.Notes = app.Notes
		})
		if err != nil {
			fmt.Printf("Error updating application: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Application updated successfully!")
		fmt.Printf("ID: %s\n", app.ID)
		fmt.Printf("Company: %s\n", app.Company)
		fmt.Printf("Role: %s\n", app.Role)
		fmt.Printf("Status: %s\n", app.Status)
		fmt.Printf("Date Applied: %s\n", app.DateApplied)
	} else {
		// Interactive mode
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

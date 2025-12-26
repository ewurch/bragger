# Agent Instructions

This workspace uses **Bragger** for job application tracking and AI-powered resume/cover letter generation.

## Candidate Knowledge Base

A knowledge base (`candidate-kb.jsonl`) stores verified candidate information for factual resume/cover letter generation. **All resume and cover letter content MUST be sourced from this KB.**

**Quick reference:**
```bash
bragger kb show              # View all entries (tabular)
bragger kb context           # Export full KB in markdown (for LLM context)
bragger kb add               # Add new entries
bragger kb update <id>       # Update existing entry
bragger kb remove <id>       # Remove entry
```

**When generating resumes or cover letters, agents MUST:**
1. Load KB context via `bragger kb context`
2. Load the application/JD via `bragger show <app-id>`
3. Analyze JD requirements against KB entries (gap analysis)
4. Identify gaps and collaborate with user to fill them
5. **Get explicit user approval before proceeding with any gaps**
6. Add any new reusable information to KB for future applications
7. Generate content - **EVERY claim must be traceable to a KB entry**

See `/resume-builder` and `/cover-letter` skills for detailed pre-generation protocols.

## Application Tracking

```bash
bragger add                  # Add a new job application
bragger list                 # List all applications
bragger show <id>            # Show application details with JD
bragger update <id>          # Update application status
bragger remove <id>          # Remove an application
```

## PDF Generation

After generating HTML resumes or cover letters:
```bash
npm run pdf outputs/company_role/resume.html
npm run pdf outputs/company_role/cover_letter.html
```

This will generate PDFs and:
- Copy the file path to clipboard (for Cmd+Shift+G in file dialogs)
- Open Finder with the PDF selected (for drag-and-drop)

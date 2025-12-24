# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a resume builder project that generates tailored, ATS-optimized resumes from LinkedIn PDF exports. It uses the `/resume-builder` skill to create HTML resumes that can be printed to PDF.

## Usage

### Creating a Resume

1. Provide a LinkedIn PDF export (e.g., `Profile-4.pdf`)
2. Provide the target job description
3. Optionally specify: location, citizenship, company region
4. Run the `/resume-builder` skill

The skill generates an HTML file in the `output/` directory with the naming convention `resume_[company]_[role].html`.

### Converting to PDF

Open the generated HTML in Chrome, then:
- Print (Cmd+P)
- Destination: "Save as PDF"
- Margins: None
- Enable: "Background graphics"
- Disable: "Headers and footers"

## Issue Tracking

This project uses **bd (beads)** for issue tracking.
Run `bd prime` for workflow context, or install hooks (`bd hooks install`) for auto-injection.

**Quick reference:**
- `bd ready` - Find unblocked work
- `bd create "Title" --type task --priority 2` - Create issue
- `bd close <id>` - Complete work
- `bd sync` - Sync with git (run at session end)

## Content Rules

**NEVER use em-dashes (—) in any generated content.** Use alternatives instead:
- Use a regular hyphen with spaces: ` - `
- Use a colon: `:`
- Restructure the sentence

This applies to all skills, resumes, cover letters, and any generated text.

## Directory Structure

- `outputs/[company]_[role]/` - Generated files per job application
  - `resume.html` / `resume.pdf`
  - `cover_letter.html` / `cover_letter.pdf`
- `output/` - Legacy files (pre-restructure)
- `scripts/` - CLI tools (html-to-pdf.js)
- `.claude/skills/` - Skill definitions

## PDF Generation

After generating HTML files, convert to PDF:
```bash
npm install                                    # First time only
npm run pdf outputs/company_role/resume.html   # Convert single file
```

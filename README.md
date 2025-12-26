# Bragger

Job application tracker and AI-powered resume/cover letter generator.

Bragger helps you:
- Track job applications with status, notes, and job descriptions
- Maintain a knowledge base of your professional profile
- Generate tailored resumes and cover letters using AI (with Claude/OpenCode skills)
- Ensure factual consistency by sourcing all content from your verified KB

## Installation

### Homebrew (macOS)

```bash
brew tap ewurch/tap
brew install bragger
```

### Go Install

```bash
go install github.com/ewurch/braggerger/cmd/bragger@latest
```

### From Source

```bash
git clone https://github.com/ewurch/braggerger.git
cd bragger
go build -o bragger ./cmd/bragger
```

## Quick Start

### 1. Initialize a Workspace

```bash
mkdir my-job-search
cd my-job-search
bragger init
```

This creates:
- `.bragger/` - Version tracking
- `.claude/skills/` - AI agent skills for resume/cover letter generation
- `AGENTS.md` - Instructions for AI agents
- `applications.jsonl` - Application tracking data
- `candidate-kb.jsonl` - Your professional knowledge base
- `outputs/` - Generated resumes and cover letters
- `package.json` - For PDF generation

### 2. Install PDF Dependencies

```bash
npm install
```

### 3. Add Your Profile to the Knowledge Base

```bash
# Add contact info
bragger kb add --type profile --category contact \
  --data '{"name":"Your Name","email":"you@example.com","location":"City, Country"}'

# Add work experience
bragger kb add --type profile --category experience \
  --data '{"company":"Acme Corp","role":"Senior Engineer","start_date":"2020-01","end_date":"present","highlights":["Led team of 5","Reduced latency by 40%"]}'

# Add skills
bragger kb add --type profile --category skills \
  --data '{"languages":["Go","Python"],"frameworks":["React","FastAPI"],"cloud":["AWS","GCP"]}'
```

### 4. Track Applications

```bash
# Add a new application
bragger add --company "Dream Company" --role "Staff Engineer" --jd-file job-description.txt

# List all applications
bragger list

# Show details
bragger show app-a1b2c3d4

# Update status
bragger update app-a1b2c3d4 --status "interviewing" --notes "Phone screen scheduled"
```

### 5. Generate Resumes (with AI)

When using Claude or OpenCode, invoke the resume-builder skill:

```
/resume-builder

Generate a resume for application app-a1b2c3d4
```

The AI will:
1. Load your KB context
2. Analyze the job description
3. Perform gap analysis
4. Ask for any missing information
5. Generate a tailored HTML resume
6. Save to `outputs/company_role/resume.html`

### 6. Convert to PDF

```bash
npm run pdf outputs/company_role/resume.html
```

The PDF will be created and Finder will open with the file selected for easy upload.

## Commands

| Command | Description |
|---------|-------------|
| `bragger init` | Initialize a new workspace |
| `bragger add` | Add a new application |
| `bragger list` | List all applications |
| `bragger show <id>` | Show application details |
| `bragger update <id>` | Update an application |
| `bragger remove <id>` | Remove an application |
| `bragger kb show` | Show knowledge base entries |
| `bragger kb context` | Export KB in markdown (for AI) |
| `bragger kb add` | Add a KB entry |
| `bragger kb update <id>` | Update a KB entry |
| `bragger kb remove <id>` | Remove a KB entry |
| `bragger upgrade` | Upgrade workspace to latest version |
| `bragger help` | Show help |

## Knowledge Base Categories

### Profile Entries (structured data)

| Category | Required Fields | Optional Fields |
|----------|-----------------|-----------------|
| `contact` | name, email | phone, location, linkedin, github, website |
| `experience` | company, role, start_date | end_date, location, description, highlights |
| `education` | institution, degree | field, start_date, end_date, gpa |
| `skills` | (none) | languages, frameworks, tools, databases, cloud, other |
| `certifications` | name | issuer, date, expiry_date, credential_id |
| `languages` | language | proficiency |

### Context Entries (flexible text)

Context entries store additional information that doesn't fit structured categories:

```bash
bragger kb add --type context --category achievement \
  --content "Led migration to microservices, reducing latency by 40%"

bragger kb add --type context --category preference \
  --content "Prefer remote-first companies with async culture"
```

## AI Integration

Bragger includes Claude/OpenCode skills that enforce factual consistency:

1. **Pre-generation protocol**: AI must load KB and perform gap analysis
2. **Strict factuality**: Every claim must trace to a KB entry
3. **User approval**: Gaps require explicit approval before proceeding
4. **KB enrichment**: New information is added for future applications

See `.claude/skills/` for detailed skill documentation.

## Development

```bash
# Run tests
go test ./...

# Build
go build -o bin/bragger ./cmd/bragger

# Test in a new workspace
mkdir /tmp/test-workspace && cd /tmp/test-workspace
/path/to/bin/bragger init
```

## License

MIT License - see [LICENSE](LICENSE) for details.

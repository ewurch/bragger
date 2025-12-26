# Brag

Job application tracker and AI-powered resume/cover letter generator.

Brag helps you:
- Track job applications with status, notes, and job descriptions
- Maintain a knowledge base of your professional profile
- Generate tailored resumes and cover letters using AI (with Claude/OpenCode skills)
- Ensure factual consistency by sourcing all content from your verified KB

## Installation

### Homebrew (macOS)

```bash
brew tap ewurch/tap
brew install brag
```

### Go Install

```bash
go install github.com/ewurch/brag/cmd/brag@latest
```

### From Source

```bash
git clone https://github.com/ewurch/brag.git
cd brag
go build -o brag ./cmd/brag
```

## Quick Start

### 1. Initialize a Workspace

```bash
mkdir my-job-search
cd my-job-search
brag init
```

This creates:
- `.brag/` - Version tracking
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
brag kb add --type profile --category contact \
  --data '{"name":"Your Name","email":"you@example.com","location":"City, Country"}'

# Add work experience
brag kb add --type profile --category experience \
  --data '{"company":"Acme Corp","role":"Senior Engineer","start_date":"2020-01","end_date":"present","highlights":["Led team of 5","Reduced latency by 40%"]}'

# Add skills
brag kb add --type profile --category skills \
  --data '{"languages":["Go","Python"],"frameworks":["React","FastAPI"],"cloud":["AWS","GCP"]}'
```

### 4. Track Applications

```bash
# Add a new application
brag add --company "Dream Company" --role "Staff Engineer" --jd-file job-description.txt

# List all applications
brag list

# Show details
brag show app-a1b2c3d4

# Update status
brag update app-a1b2c3d4 --status "interviewing" --notes "Phone screen scheduled"
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
| `brag init` | Initialize a new workspace |
| `brag add` | Add a new application |
| `brag list` | List all applications |
| `brag show <id>` | Show application details |
| `brag update <id>` | Update an application |
| `brag remove <id>` | Remove an application |
| `brag kb show` | Show knowledge base entries |
| `brag kb context` | Export KB in markdown (for AI) |
| `brag kb add` | Add a KB entry |
| `brag kb update <id>` | Update a KB entry |
| `brag kb remove <id>` | Remove a KB entry |
| `brag upgrade` | Upgrade workspace to latest version |
| `brag help` | Show help |

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
brag kb add --type context --category achievement \
  --content "Led migration to microservices, reducing latency by 40%"

brag kb add --type context --category preference \
  --content "Prefer remote-first companies with async culture"
```

## AI Integration

Brag includes Claude/OpenCode skills that enforce factual consistency:

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
go build -o bin/brag ./cmd/brag

# Test in a new workspace
mkdir /tmp/test-workspace && cd /tmp/test-workspace
/path/to/bin/brag init
```

## License

MIT License - see [LICENSE](LICENSE) for details.

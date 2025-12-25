---
name: candidate-kb
description: Manage the candidate knowledge base. Use when you need to store, query, or update factual information about the candidate for resume/cover letter generation.
---

# Candidate Knowledge Base Skill

You manage the candidate's knowledge base - a structured repository of factual information used to generate accurate, consistent resumes and cover letters.

## Purpose

The knowledge base serves three critical functions:

1. **Factuality** - All claims in resumes/cover letters MUST be sourced from here
2. **Consistency** - Same facts presented consistently across applications  
3. **Progressive Learning** - New information gathered during applications is stored for future use

## Storage Location

The knowledge base is stored in `candidate-kb.jsonl` in the project root. Each line is a JSON object representing one knowledge entry.

## Data Model

Each entry has a `type` field indicating the category:

### Entry Types

#### 1. `profile` - Core candidate information
```json
{
  "type": "profile",
  "id": "profile-001",
  "field": "citizenship",
  "value": "Italian",
  "context": "EU work authorization",
  "created_at": "...",
  "updated_at": "..."
}
```

Fields: `name`, `email`, `phone`, `location`, `linkedin`, `github`, `website`, `citizenship`, `languages`, `headline`

#### 2. `experience` - Work experience details
```json
{
  "type": "experience",
  "id": "exp-001",
  "company": "Company Name",
  "role": "Job Title",
  "start_date": "2020-01",
  "end_date": "2023-06",
  "location": "City, Country",
  "description": "Brief role description",
  "achievements": [
    "Led team of 5 engineers to deliver ML pipeline, reducing processing time by 40%",
    "Implemented CI/CD pipeline serving 100+ deployments/week"
  ],
  "technologies": ["Python", "AWS", "Kubernetes"],
  "created_at": "...",
  "updated_at": "..."
}
```

#### 3. `skill` - Skills with context
```json
{
  "type": "skill",
  "id": "skill-001",
  "name": "Python",
  "category": "programming_language",
  "proficiency": "expert",
  "years_experience": 8,
  "context": "Used daily for ML pipelines, data processing, and API development",
  "related_projects": ["exp-001", "proj-002"],
  "created_at": "...",
  "updated_at": "..."
}
```

Categories: `programming_language`, `framework`, `tool`, `platform`, `methodology`, `soft_skill`, `domain`
Proficiency: `beginner`, `intermediate`, `advanced`, `expert`

#### 4. `project` - Notable projects
```json
{
  "type": "project",
  "id": "proj-001",
  "name": "Project Name",
  "description": "What the project does",
  "role": "Lead Developer",
  "outcomes": [
    "Reduced latency by 60%",
    "Serves 1M+ requests/day"
  ],
  "technologies": ["Python", "FastAPI", "PostgreSQL"],
  "url": "https://github.com/...",
  "start_date": "2022-01",
  "end_date": "2022-06",
  "created_at": "...",
  "updated_at": "..."
}
```

#### 5. `education` - Education and certifications
```json
{
  "type": "education",
  "id": "edu-001",
  "institution": "University Name",
  "degree": "Master of Science",
  "field": "Computer Science",
  "graduation_year": "2018",
  "gpa": "3.8",
  "honors": "Cum Laude",
  "relevant_coursework": ["Machine Learning", "Distributed Systems"],
  "created_at": "...",
  "updated_at": "..."
}
```

#### 6. `certification` - Professional certifications
```json
{
  "type": "certification",
  "id": "cert-001",
  "name": "AWS Solutions Architect Professional",
  "issuer": "Amazon Web Services",
  "date_obtained": "2023-01",
  "expiry_date": "2026-01",
  "credential_id": "ABC123",
  "created_at": "...",
  "updated_at": "..."
}
```

#### 7. `achievement` - Standalone achievements
```json
{
  "type": "achievement",
  "id": "ach-001",
  "description": "Published paper on LLM optimization in NeurIPS 2023",
  "category": "publication",
  "date": "2023-12",
  "metrics": "Cited 50+ times",
  "url": "https://...",
  "created_at": "...",
  "updated_at": "..."
}
```

Categories: `publication`, `award`, `speaking`, `open_source`, `patent`, `other`

#### 8. `fact` - Contextual facts learned during applications
```json
{
  "type": "fact",
  "id": "fact-001",
  "category": "preference",
  "content": "Prefers remote-first companies",
  "source": "User stated during ACME Corp application",
  "created_at": "...",
  "updated_at": "..."
}
```

Categories: `preference`, `constraint`, `goal`, `strength`, `weakness`, `other`

## Capabilities

### 1. Query Knowledge Base

Read and search the knowledge base for relevant information.

**User requests:**
- "What skills does the candidate have?"
- "Show me the candidate's work experience"
- "What achievements can we use for this role?"

**Process:**
1. Read `candidate-kb.jsonl`
2. Filter by requested type or search terms
3. Return formatted results

### 2. Add New Information

Store new information gathered during resume/cover letter generation.

**User requests:**
- "Add that I have 3 years of Kubernetes experience"
- "Store this achievement: Led migration to microservices"

**Process:**
1. Read existing entries
2. Check for duplicates or conflicts
3. Generate unique ID (`type-` + 8 hex chars)
4. Append new entry with timestamps
5. Confirm addition

### 3. Update Existing Information

Modify existing entries when information changes or is refined.

**User requests:**
- "Update my Python experience to 10 years"
- "Add a new achievement to my ACME Corp experience"

**Process:**
1. Find matching entry
2. Update specified fields
3. Update `updated_at` timestamp
4. Save changes

### 4. Detect Gaps

When analyzing a job description, identify missing information that would strengthen the application.

**Process:**
1. Extract requirements from JD
2. Query KB for matching skills/experience
3. Identify gaps (required but not in KB)
4. Prompt user for missing information
5. Store responses in KB

### 5. Import from LinkedIn PDF

Parse a LinkedIn PDF export and populate the knowledge base.

**Process:**
1. Read and parse LinkedIn PDF
2. Extract profile, experience, education, skills
3. Convert to KB entry format
4. Prompt user to confirm/refine each section
5. Store confirmed entries

## Integration with Resume Builder

When the resume-builder skill needs candidate information:

1. **Query First**: Always query KB before generating content
2. **Detect Gaps**: If JD requires something not in KB, prompt user
3. **Store Learnings**: Any new information provided gets stored
4. **Cite Sources**: Generated content should be traceable to KB entries

### Gap Detection Workflow

```
JD Analysis → Required Skills/Experience
                    ↓
              Query KB
                    ↓
         ┌─────────┴─────────┐
         ↓                   ↓
    Found in KB         Not in KB
         ↓                   ↓
    Use for resume      Prompt user:
                        "The JD mentions X.
                         Do you have experience with this?
                         Please describe..."
                              ↓
                        Store response in KB
                              ↓
                        Use for resume
```

## Factuality Enforcement

**CRITICAL RULE**: The resume-builder and cover-letter skills MUST NOT fabricate information.

- All claims must be sourced from KB or explicitly provided by user in the current session
- If information isn't in KB and user doesn't provide it, omit it from resume
- Never assume, infer, or make up achievements, metrics, or experiences
- When in doubt, ask the user

## Example Interactions

**User**: "Add my Kubernetes experience - I've been using it for 3 years in production"

**Response**:
1. Create skill entry for Kubernetes
2. Set proficiency based on years
3. Ask for additional context: "Can you share a specific achievement or project using Kubernetes?"
4. Store the full entry
5. Confirm: "Added Kubernetes (3 years, production experience) to your knowledge base"

---

**User**: "Import my LinkedIn PDF"

**Response**:
1. Read the provided LinkedIn PDF
2. Extract sections (experience, education, skills)
3. For each section, show extracted info and ask:
   - "I found these 5 positions. Would you like to add details or achievements for any?"
4. Store confirmed entries
5. Summarize: "Added X experience entries, Y skills, Z education records"

---

**Skill (resume-builder)**: "Need candidate's cloud experience for AWS role"

**Response**:
1. Query KB for cloud-related skills and experience
2. Return: AWS (expert, 5 years), GCP (intermediate, 2 years)
3. Related achievements from experience entries
4. If no cloud experience found, indicate gap for user prompting

## Implementation Notes

### Reading candidate-kb.jsonl

Use the Read tool to read the file, then parse each line as JSON.

### Writing candidate-kb.jsonl

When adding or updating, read all entries, modify as needed, then write back the entire file.

### ID Generation

Generate IDs in format: `type-` + 8 random hex characters (e.g., `skill-a1b2c3d4`)

### Deduplication

Before adding new entries, check for existing entries with similar content to avoid duplicates.

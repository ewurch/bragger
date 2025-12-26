# Agent Instructions

This project uses **bd** (beads) for issue tracking.
Run `bd prime` for workflow context, or install hooks (`bd hooks install`) for auto-injection.

## Candidate Knowledge Base

A knowledge base (`candidate-kb.jsonl`) stores verified candidate information for factual resume/cover letter generation. **All resume and cover letter content MUST be sourced from this KB.**

**Quick reference:**
```bash
app kb show              # View all entries (tabular)
app kb context           # Export full KB in markdown (for LLM context)
app kb add               # Add new entries
app kb update <id>       # Update existing entry
app kb remove <id>       # Remove entry
```

**When generating resumes or cover letters, agents MUST:**
1. Load KB context via `app kb context`
2. Load the application/JD via `app show <app-id>`
3. Analyze JD requirements against KB entries (gap analysis)
4. Identify gaps and collaborate with user to fill them
5. **Get explicit user approval before proceeding with any gaps**
6. Add any new reusable information to KB for future applications
7. Generate content - **EVERY claim must be traceable to a KB entry**

See `/resume-builder` and `/cover-letter` skills for detailed pre-generation protocols.

## Quick Reference

```bash
bd ready              # Find available work
bd create "Title" --type task --priority 2  # Create issue
bd show <id>          # View issue details
bd update <id> --status in_progress  # Claim work
bd close <id>         # Complete work
bd sync               # Sync with git
```

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds


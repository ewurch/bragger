package templates

import "embed"

//go:embed AGENTS.md package.json scripts/* skills/*/SKILL.md
var Files embed.FS

// Version is the current version of Bragger
const Version = "0.4.0"

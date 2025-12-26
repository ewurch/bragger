package templates

import "embed"

//go:embed AGENTS.md package.json scripts/* skills/*/SKILL.md
var Files embed.FS

// Version is the current version of Brag
const Version = "0.1.0"

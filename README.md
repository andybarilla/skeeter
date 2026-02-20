# Skeeter

File-based project management for coding agents.

Skeeter stores tasks as markdown files with YAML frontmatter directly in your git repository. Agents read the files, pick up work, and update status — no API integrations, no platform lock-in. Git gives you history, branching, and collaboration for free.

## Install

```bash
go install github.com/andybarilla/skeeter/cmd/skeeter@latest
```

## Quick Start

```bash
# Initialize in your repo
skeeter init

# Create tasks
skeeter create "Add user authentication" -p high -t auth,security
skeeter create "Fix login redirect bug" -p critical -T bug

# Manage tasks
skeeter list
skeeter list --status backlog --priority high
skeeter show US-001
skeeter status US-001 ready-for-development
skeeter assign US-001 claude
skeeter edit US-001

# Agent workflow
skeeter next                    # Show highest-priority available task
skeeter next --assign claude    # Claim it and move to in-progress
skeeter next --quiet            # Output just the ID (for scripting)
```

## How It Works

Skeeter creates a `.skeeter/` directory in your repo:

```
.skeeter/
  config.yaml          # Project settings
  SKEETER.md           # Auto-generated agent instructions
  tasks/
    US-001.md          # One file per task
    US-002.md
  templates/
    default.md         # Task body templates
    bug.md
```

Each task is a markdown file with YAML frontmatter:

```markdown
---
id: US-001
title: Add user authentication
status: ready-for-development
priority: high
assignee:
tags: [auth, security]
created: "2026-02-20"
updated: "2026-02-20"
---

## Acceptance Criteria

- [ ] Users can sign up with email/password
- [ ] JWT tokens issued on login

## Context

- `src/routes/auth.ts` - existing stub
```

## Agent Integration

Skeeter generates a `SKEETER.md` file that acts as a natural language API for any coding agent. Agents that explore the repo will find it and immediately understand the task protocol:

1. Find tasks where `status: ready-for-development` and `assignee:` is empty
2. Set `assignee` and `status: in-progress` before starting
3. Use `Acceptance Criteria` as the definition of done
4. Set `status: done` when complete

No special integration needed — any agent that can read and write files works with Skeeter.

## Remote Access

Manage tasks via the GitHub API without a local clone:

```bash
skeeter --remote owner/repo list
skeeter --remote owner/repo create "Fix deployment" -p critical
skeeter --remote owner/repo status US-003 done
```

Authenticates via `gh auth token` or `GITHUB_TOKEN` environment variable.

## Configuration

```bash
skeeter config                                          # View settings
skeeter config set prefix TASK                          # Change ID prefix
skeeter config set statuses "backlog,todo,doing,done"   # Custom workflow
skeeter config set priorities "p0,p1,p2,p3"             # Custom priorities
skeeter config set auto_commit true                     # Auto-commit changes
```

## Templates

Tasks are created from templates stored in `.skeeter/templates/`. A `default.md` template is generated on init.

```bash
skeeter create "New feature" -T default     # Use default (this is the default)
skeeter create "Bug report" -T bug          # Use bug template
skeeter create "Quick note" --no-template   # Empty body
```

Add your own templates by dropping markdown files in `.skeeter/templates/`.

## Configurable Directory

The `.skeeter/` directory can be overridden for testing or when using Skeeter on itself:

```bash
skeeter --dir .project list              # Use .project/ instead of .skeeter/
SKEETER_DIR=.dev skeeter list            # Via environment variable
```

Resolution order: `--dir` flag > `SKEETER_DIR` env > auto-detect `.skeeter/` walking up from cwd.

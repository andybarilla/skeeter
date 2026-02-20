# Skeeter — Project Tasks

Tasks are markdown files with YAML frontmatter in the `tasks/` subdirectory.

## For Agents: Finding Work

1. Look for tasks where `status: ready-for-development` and `assignee:` is empty
2. Set `assignee: <your-name>` and `status: in-progress` before starting
3. Use `Acceptance Criteria` as your definition of done
4. Set `status: done` when complete

## Frontmatter Fields

| Field      | Description                                              |
|------------|----------------------------------------------------------|
| id         | Task identifier (e.g., US-001)                           |
| title      | Short task title                                         |
| status     | One of: backlog, ready-for-development, in-progress, done   |
| priority   | One of: critical, high, medium, low  |
| assignee   | Who is working on this (empty = available)               |
| tags       | Array of labels                                          |
| links      | Related URLs                                             |
| created    | Creation date                                            |
| updated    | Last modified date                                       |

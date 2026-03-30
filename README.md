# seifenkehrer

A modular cleanup tool to keep your macOS clean and free of clutter.

Define cleanup tasks as simple YAML files, review what will be deleted, and choose what to remove.

## Installation

```bash
go install github.com/seifenkehrer/seifenkehrer@latest
```

Or build from source:

```bash
git clone https://github.com/seifenkehrer/seifenkehrer.git
cd seifenkehrer
go build -o sk .
```

## Usage

```bash
# List installed cleanup tasks
sk tasks

# Run cleanup (interactive)
sk clean

# Use a custom tasks directory
sk clean --tasks-dir ./my-tasks
```

## Task Configuration

Tasks are YAML files placed in `~/.sk/tasks/`. See the `examples/` directory for ready-to-use tasks.

```yaml
name: Chrome Old Versions
description: Removes outdated Google Chrome framework versions
globs:
  - /Applications/Google Chrome.app/Contents/Frameworks/Google Chrome Framework.framework/Versions/*
exclude:
  - Current
keep_newest: 1
interval: 168h
```

### Fields

| Field | Required | Description |
|-------|----------|-------------|
| `name` | No | Display name (defaults to filename) |
| `description` | Yes | What this task cleans up |
| `globs` | Yes | File patterns to match |
| `exclude` | No | Basenames to skip |
| `keep_newest` | No | Keep the N most recent matches |
| `interval` | No | Minimum time between runs (e.g. `24h`, `168h`) |

## Task Management

```bash
# Enable/disable a task
sk config disable <task-name>
sk config enable <task-name>

# Set run interval
sk config interval <task-name> <duration>
```

## How It Works

1. Loads all `.yaml`/`.yml` files from the tasks directory
2. Resolves globs, applies exclusions and `keep_newest`
3. Skips tasks whose interval hasn't elapsed
4. Shows matched paths grouped by task with sizes
5. Prompts per task: accept all, skip, or choose individually

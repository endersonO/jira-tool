# jt — Jira Tool

A fast, single-binary CLI for Jira Cloud, inspired by `gh`. Built for developers who want to manage Jira without leaving the terminal — and designed to work seamlessly with AI assistants like Claude Code.

## Why jt?

- **No runtime dependencies** — single binary, works anywhere
- **`jt init`** — interactive setup, no manual config file editing
- **Rich text support** — write descriptions in Markdown, renders correctly in Jira
- **LLM-friendly output** — clean tables and JSON mode for AI analysis
- **Familiar UX** — if you've used `gh`, you already know `jt`

---

## Installation

### Option 1: Quick Install Script (macOS / Linux — Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/endersonO/jira-tool/main/install.sh | sh
```

This auto-detects your OS and architecture, downloads the correct binary, and installs it to `/usr/local/bin`.

### Option 2: Download Binary Manually

Go to the [Releases page](https://github.com/endersonO/jira-tool/releases/latest) and download the file for your platform:

| Platform | File |
|----------|------|
| macOS Apple Silicon (M1/M2/M3/M4) | `jt_<version>_darwin_arm64.tar.gz` |
| macOS Intel | `jt_<version>_darwin_amd64.tar.gz` |
| Linux AMD64 | `jt_<version>_linux_amd64.tar.gz` |
| Linux ARM64 | `jt_<version>_linux_arm64.tar.gz` |
| Windows | `jt_<version>_windows_amd64.zip` |

**macOS / Linux — after downloading:**
```bash
tar -xzf jt_*.tar.gz
sudo mv jt /usr/local/bin/
jt --version
```

**Windows — after downloading:**
1. Extract the ZIP file
2. Move `jt.exe` to a folder in your PATH (e.g., `C:\Users\<you>\bin\`)
3. Open a new terminal and run `jt --version`

### Option 3: Install from Source (Requires Go 1.26+)

```bash
go install github.com/endersonO/jt/cmd/jt@latest
```

This places `jt` in `~/go/bin/`. Make sure it's in your PATH:

```bash
# Add to ~/.zshrc or ~/.bashrc
export PATH="$HOME/go/bin:$PATH"
```

### First-time setup

```bash
jt init
```

You'll be prompted for:

```
Jira server URL [https://your-org.atlassian.net]:
Email: you@example.com
API Token (press Enter to keep existing):
Default project key (optional): SCRUM
```

> Generate your API token at: https://id.atlassian.com/manage-profile/security/api-tokens

Config is saved to your OS config directory:
- **macOS:** `~/Library/Application Support/jt/config.yml`
- **Linux:** `~/.config/jt/config.yml`
- **Windows:** `%APPDATA%\jt\config.yml`

---

## Commands

### Issues

```bash
# List issues in your default project
jt issue list

# Filter by status, assignee, or type
jt issue list --status "In Progress"
jt issue list --assignee me
jt issue list --assignee unassigned
jt issue list --type Bug
jt issue list --status "To Do" --assignee me --type Story

# Show more columns (type, priority)
jt issue list -v

# View a single issue
jt issue view SCRUM-123

# Create an issue
jt issue create --summary "Fix login redirect" --type Bug --priority High

# Write description in Markdown (renders as rich text in Jira)
jt issue create --summary "Auth refactor" \
  --type Story \
  --description "## Goals\n- Replace session tokens\n- Add OAuth2 support" \
  --assignee me \
  --priority High \
  --labels backend,security

# Open $EDITOR to write description
jt issue create --summary "New feature" --edit

# Edit an existing issue
jt issue edit SCRUM-123 --summary "Updated title"
jt issue edit SCRUM-123 --status "In Progress"
jt issue edit SCRUM-123 --assignee dev@example.com --priority Highest

# Transition status
jt issue transition SCRUM-123              # list available statuses
jt issue transition SCRUM-123 "In Review"
```

### Search (JQL)

```bash
jt search "project=SCRUM AND assignee=currentUser()"
jt search "project=SCRUM AND status='In Progress'"
jt search "project=SCRUM AND issuetype=Bug AND priority=High"
jt search "created >= -7d ORDER BY updated DESC"
```

### Projects

```bash
jt project list
```

### Configuration

```bash
jt init       # interactive setup / reconfigure
jt config     # show current config status
```

---

## Flags

| Flag | Description |
|------|-------------|
| `--max N` | Limit results (default: 30) |
| `-v, --verbose` | Show extra columns (type, priority) |
| `--json` | Output raw JSON for scripting or AI analysis |
| `--assignee me` | Filter by current user |
| `--project KEY` | Override default project |
| `--edit` | Open `$EDITOR` to write Markdown description |

---

## Environment Variables

For CI/CD or shared environments where you don't want a config file:

```bash
export JT_SERVER="https://your-org.atlassian.net"
export JT_EMAIL="you@example.com"
export JT_TOKEN="your-api-token"
export JT_PROJECT="SCRUM"
```

---

## Use with AI Assistants

`jt` is designed to feed clean data into AI tools like Claude Code.

```bash
# Analyze your backlog
jt issue list --status "To Do" -v

# Export as JSON for deeper analysis
jt issue list --json
jt search "project=SCRUM AND sprint in openSprints()" --json

# View full issue details including description
jt issue view SCRUM-123
```

**Example workflow with Claude Code:**

```
> jt issue list --assignee me -v
# Share output with Claude
"Based on these issues, help me figure out what to work on next"
```

---

## Issue Types

| Type | Description |
|------|-------------|
| `Task` | Standard work item |
| `Story` | User story |
| `Bug` | Bug report |
| `Epic` | Large feature or initiative |
| `Subtask` | Child of another issue |

## Priorities

`Highest` · `High` · `Medium` · `Low` · `Lowest`

---

## Troubleshooting

**`not configured — run jt init to get started`**
Run `jt init` to set up your credentials.

**`401 Unauthorized`**
- Verify your API token is valid (not your Jira password)
- Run `jt init` to update your credentials

**`jt: command not found`**
Make sure the binary is in your PATH:
```bash
# If installed via install.sh or manual download:
which jt  # should show /usr/local/bin/jt

# If installed via go install:
export PATH="$HOME/go/bin:$PATH"
```

**macOS: "jt can't be opened because it is from an unidentified developer"**
```bash
# Remove the quarantine attribute
xattr -d com.apple.quarantine /usr/local/bin/jt
```

**Old version still running**
Check which binary is being used: `which jt`

---

## License

MIT

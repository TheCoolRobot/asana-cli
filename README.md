# Asana CLI

A feature-rich command-line interface for Asana with an interactive TUI (Text User Interface), JSON output for automation, and a background sync daemon for offline support.

## ✨ Features

- **Interactive TUI** - Navigate, filter, and manage tasks with Bubble Tea
- **JSON Output** - Easy to parse format for scripts and automation
- **Sync Daemon** - Background service that caches your Asana data locally
- **Full CRUD** - Create, read, update, and delete tasks
- **Advanced Filtering** - Filter by assignee, tags, status, priority
- **Search** - Find tasks across your workspace
- **Task Management** - Complete, assign, set priorities and due dates
- **Configuration** - Save API tokens and default workspaces/projects

## 🚀 Quick Start

### Installation

```bash
# Clone the repo
git clone https://github.com/TheCoolRobot/asana-cli.git
cd asana-cli

# Build
make build

# Or install directly
go install
```

I like to add a function to make accessing it as easy as possible. To do this, I use a zshrc

```bash
function asana-cli(){
  cd /Users/henry/Developer/cmdln_dev/asana_cli-copilot&&./asana-cli "$@"
}
# Replace /Users/henry/Developer/cmdln_dev with the PATH to YOUR folder
```

### Authentication

Set your Asana API token:

```bash
export ASANA_TOKEN=your-token-here
# or save to config
asana-cli config set --token your-token-here
```

### Basic Usage

```bash
# List tasks in a project (interactive TUI)
asana-cli list <project-id>

# List tasks as JSON
asana-cli list <project-id> --json

# Create a task
asana-cli create <project-id> --name "My Task" --priority high

# Update a task
asana-cli update <task-id> --name "Updated Task"

# Complete a task
asana-cli complete <task-id>

# Search for tasks
asana-cli search <workspace-id> "bug fix"

# Start sync daemon
asana-cli sync --projects 12345,67890
```

## 📖 Documentation

- [Installation Guide](docs/INSTALLATION.md)
- [Usage Guide](docs/USAGE.md)
- [API Reference](docs/API.md)
- [Development](docs/DEVELOPMENT.md)

## 🔄 Sync Daemon

The sync daemon runs in the background and automatically caches your Asana data locally every 5 minutes. This enables:

- **Fast TUI loading** from local cache
- **Offline browsing** of cached tasks
- **Batch operations** without hitting API rate limits
- **History tracking** of task changes

Start the daemon:

```bash
asana-cli sync --projects project-id-1,project-id-2
```

Or run as a service (see [Development](docs/DEVELOPMENT.md) for systemd setup).

## 📝 JSON Output Examples

```bash
# Get tasks as JSON
$ asana-cli list proj-123 --json
{
  "success": true,
  "data": [
    {
      "id": "task-1",
      "name": "Build feature",
      "completed": false,
      "priority": "high",
      "due_date": "2026-03-01T00:00:00Z"
    }
  ],
  "meta": {
    "count": 1,
    "project_id": "proj-123"
  }
}

# Create a task and parse JSON
$ asana-cli create proj-123 --name "Review PR" --json | jq '.data.id'
"task-456"
```

## 🎮 Interactive TUI Controls

```
[↑↓] - Navigate tasks
[space] - Select/deselect
[c] - Mark complete
[f] - Toggle show completed
[s] - Change sort (name/due/priority)
[/] - Search
[q] - Quit
```

## 🔐 Security

- API tokens are stored in `~/.asana-cli/config.json` with restricted permissions (0600)
- Never commit `.env` files or tokens to version control
- Use environment variables for CI/CD

## 📊 Commands

### Task Management
- `list` - List tasks in a project
- `view` - View task details
- `create` - Create a new task
- `update` - Update a task
- `complete` - Mark task as complete
- `delete` - Delete a task
- `search` - Search for tasks

### System
- `config` - Manage configuration
- `sync` - Start sync daemon
- `me` - Show current user info

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -am 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing`)
5. Open a Pull Request

## 📄 License

MIT License - see LICENSE file for details

## 🐛 Issues & Support

Found a bug? Have a feature request? [Open an issue](https://github.com/TheCoolRobot/asana-cli/issues)

---

Made with ❤️ by TheCoolRobot
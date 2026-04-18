# stashctl

> CLI for managing and searching local code snippets with tag-based filtering

---

## Installation

```bash
go install github.com/yourusername/stashctl@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/stashctl.git
cd stashctl && go build -o stashctl .
```

---

## Usage

```bash
# Add a new snippet
stashctl add --title "reverse a string" --tags go,strings --file snippet.go

# List all snippets
stashctl list

# Search by tag
stashctl search --tags go,strings

# View a snippet
stashctl get <id>

# Delete a snippet
stashctl delete <id>
```

Snippets are stored locally in `~/.stashctl/snippets.json`.

---

## Commands

| Command  | Description                        |
|----------|------------------------------------|
| `add`    | Add a new snippet                  |
| `list`   | List all stored snippets           |
| `search` | Filter snippets by tags or keyword |
| `get`    | Display a specific snippet         |
| `delete` | Remove a snippet by ID             |

---

## License

MIT © [yourusername](https://github.com/yourusername)
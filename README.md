# stacksnap

> Captures and exports a snapshot of your local dev stack config for reproducible onboarding.

---

## Installation

```bash
go install github.com/yourusername/stacksnap@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/stacksnap.git && cd stacksnap && go build -o stacksnap .
```

---

## Usage

Capture a snapshot of your current dev stack and export it to a file:

```bash
stacksnap capture --output stack.json
```

Share the snapshot file with your team. A new team member can reproduce your environment by running:

```bash
stacksnap apply --input stack.json
```

### Additional Commands

```bash
stacksnap capture            # Snapshot current stack to stdout
stacksnap capture -o snap.json   # Save snapshot to a file
stacksnap apply -i snap.json     # Apply a snapshot to local environment
stacksnap diff snap.json         # Compare snapshot against current stack
stacksnap version                # Print version info
```

---

## What Gets Captured

- Language runtimes (Go, Node, Python, Ruby, etc.)
- Package manager versions
- Environment variables (non-sensitive)
- Common CLI tool versions
- Docker and container runtime info

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

---

## License

[MIT](LICENSE)
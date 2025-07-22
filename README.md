# Type CLI

A terminal-based typing speed test built with Go and Bubble Tea.

![Demo](assets/demo.gif)

## Quick Try

You can try it out immediately without installation:
```bash
ssh type.fulsiram.dev
```

## Statistics
The app currently tracks:
- **WPM** (Words Per Minute) - Raw typing speed
- **Accuracy** - Percentage of correctly typed characters

## Installation

### Prerequisites

- Go 1.24.4 or later

## Run locally

```bash
git clone https://github.com/fulsiram/type-cli.git
cd type-cli
go install ./cmd/type-cli/
type-cli
```

## Launch SSH Server

To run your own SSH server:

```bash
git clone https://github.com/fulsiram/type-cli.git
cd type-cli
go install ./cmd/type-cli-ssh-server/
type-cli-ssh-server
```

### Configuration

SSH server can be configured using environment variables:

- `SSH_HOST` - Server host (default: localhost)
- `SSH_PORT` - Server port (default: 31337)
- `SSH_HOST_KEY` - Path to server SSH host key (default: .ssh/id_ed25519)

## Contributing

Contributions are welcome! Please feel free to submit your Pull Request.


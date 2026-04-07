# NBP Exchange Rate CLI

> **Note**: This is my first ever Go project, so please don't judge! If you see something that could be improved, I'd appreciate it if you could correct me instead. Thank you!

A simple and efficient Go-based command-line tool to fetch currency exchange rates from the National Bank of Poland (NBP).

## Features

- **Normal Mode**: Get the exchange rate for a specific date or the latest available rate.
- **Tax Mode**: Automatically find the exchange rate from the last business day preceding a given date (useful for Polish tax calculations).

## Installation

### Homebrew (macOS/Linux)

The recommended way to install on macOS and Linux:

```bash
brew tap MarcusSlover/tap
brew install nbp
```

*This will automatically configure shell autocompletion for you.*

### Other Options

- **Go Install**: `go install github.com/MarcusSlover/nbp-cli@latest`
- **Releases**: Download pre-built binaries from the [Releases](https://github.com/MarcusSlover/nbp-cli/releases) page.
- **From Source**: Clone the repo and run `make install`.

## Usage

```bash
nbp [flags] <currency_code> <date|today>
```

### Examples

| Task | Command |
| :--- | :--- |
| **Get current USD rate** | `nbp USD today` |
| **Get EUR rate for date** | `nbp EUR 2024-03-15` |
| **Get tax-compliant rate** | `nbp --tax EUR 2024-03-15` |

## Shell Autocompletion

If you didn't install via Homebrew, you can enable completion for your current session:

- **Bash**: `source <(nbp completion bash)`
- **Zsh**: `source <(nbp completion zsh)`
- **Fish**: `nbp completion fish | source`

## License

MIT License. See [LICENSE.txt](LICENSE.txt) for details.

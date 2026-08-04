# notse

Notse is a terminal-based notes app. It runs as a TUI in your terminal.

<p align="center">
  <img src="images/demo.gif" width="600" alt="demo">
</p>

## Installation

### Requirements

- Go >= 1.24

### Pre-built

Download the latest binary from the [releases
page](https://github.com/RahulSandhu/notse/releases).

### Build from source

```sh
git clone https://github.com/RahulSandhu/notse.git
cd notse
make install
```

## Usage

Run `notse` from your terminal.

Notes are saved to `~/.config/notse/notes_history.json`.

## How it works

Notse is a headless TUI built with [Bubble
Tea](https://github.com/charmbracelet/bubbletea). It loads your notes on
startup, lets you create, edit, pin, and delete them, and writes everything
back to a local JSON file.

| Key            | Action              |
| -------------- | ------------------- |
| `j/k` or `↑/↓` | Navigate            |
| `enter`        | View note           |
| `n`            | New note            |
| `e`            | Edit note           |
| `s`            | Cycle status (view) |
| `d`            | Delete note         |
| `p`            | Pin/unpin           |
| `tab`          | Switch field        |
| `ctrl+s`       | Save                |
| `esc`          | Back                |
| `q` / `ctrl+c` | Quit                |

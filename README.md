# notse

notse is a terminal-based notes app. It runs as a TUI inside your terminal and
stores notes in a local JSON file.

<p align="center">
  <img src="images/demo.gif" width="600" alt="demo">
</p>

## Setup

### Requirements

- Go >= 1.24

### Pre-built

Download the latest binary from the
[releases page](https://github.com/RahulSandhu/notse/releases).

### Build from source

```sh
git clone https://github.com/RahulSandhu/notse.git
cd notse
make install
```

## Configuration

Run `notse` from your terminal:

```sh
notse
```

Notes are saved to `~/.config/notse/notes_history.json`.

### Keybindings

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

<h1 align="center"> Clerk Gopher </h1>

A command-line launcher for Toontown Rewritten

## Features
- Updates your files
- Logs you into the game
- Save logins with OS-specific keyring

## Works On
- macOS (Intel*/Apple Silicon)
- Windows (32/64-bit)*
- Linux (64-bit)*

*Not tested!

## Usage

Download the file for your operating system and architecture in the bin folder.

```
# e.g. the macOS ARM64 file

$ clerk-gopher-arm64 --help
A simple command line launcher written in Go to allow simple and fast login with login saving functionality

Usage:
  clerk-gopher [command]

Available Commands:
  help        Help about any command
  login       Logs you into Toontown Rewritten with the information given
  save        Saves the given username and password to your system's keyring

Flags:
  -h, --help      help for clerk-gopher
  -v, --verbose   sets logging to verbose

Use "clerk-gopher [command] --help" for more information about a command.
```

# gemini-cli
[![Build](https://github.com/reugn/gemini-cli/actions/workflows/build.yml/badge.svg)](https://github.com/reugn/gemini-cli/actions/workflows/build.yml)

A command-line interface (CLI) for [Google Gemini](https://deepmind.google/technologies/gemini/).

Google Gemini is a family of multimodal artificial intelligence (AI) large language models that have
capabilities in language, audio, code and video understanding.

The current version only supports multi-turn conversations (chat), using the `gemini-pro` model.

## Installation
Choose a binary from the [releases](https://github.com/reugn/gemini-cli/releases).

### Build from Source
Download and [install Go](https://golang.org/doc/install).

Install the application:

```sh
go install github.com/reugn/gemini-cli/cmd/gemini@latest
```

See the [go install](https://go.dev/ref/mod#go-install) instructions for more information about the command.

## Usage

### API key
To use `gemini-cli`, you'll need an API key set in the `GEMINI_API_KEY` environment variable. If you don't already have one, create a key in [Google AI Studio](https://makersuite.google.com/app/apikey).

### System commands
The system chat message must begin with an exclamation mark and is used for internal operations.
A short list of supported system commands:

| Command | Descripton             |
| ---     | ---                    |
| !q      | Quit the application   |
| !p      | Purge the chat history |

### CLI help
```
$ ./gemini -h
Gemini CLI Tool

Usage:
   [flags]

Flags:
  -f, --format         render markdown-formatted response (default true)
  -h, --help           help for this command
  -s, --style string   markdown format style (ascii, dark, light, pink, notty, dracula) (default "auto")
  -v, --version        version for this command
```

## License
MIT

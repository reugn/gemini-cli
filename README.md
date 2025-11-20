# gemini-cli
[![Build](https://github.com/reugn/gemini-cli/actions/workflows/build.yml/badge.svg)](https://github.com/reugn/gemini-cli/actions/workflows/build.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/reugn/gemini-cli)](https://pkg.go.dev/github.com/reugn/gemini-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/reugn/gemini-cli)](https://goreportcard.com/report/github.com/reugn/gemini-cli)

A command-line interface (CLI) for [Google Gemini](https://deepmind.google/technologies/gemini/).

Google Gemini is a family of multimodal artificial intelligence (AI) large language models that have
capabilities in language, audio, code and video understanding.

This application offers a command-line interface for interacting with various generative models through
multi-turn chat. Model selection is controlled via [system command](#system-commands) inputs.

## Installation
Choose a binary from [releases](https://github.com/reugn/gemini-cli/releases).

### Build from Source
Download and [install Go](https://golang.org/doc/install).

Install the application:
```sh
go install github.com/reugn/gemini-cli/cmd/gemini@latest
```

See the [go install](https://go.dev/ref/mod#go-install) instructions for more information about the command.

## Usage
> [!NOTE]
> For information on the available regions for the Gemini API and Google AI Studio,
> see [here](https://ai.google.dev/available_regions#available_regions).

### API key
To use `gemini-cli`, you'll need an API key set in the `GEMINI_API_KEY` environment variable.
If you don't already have one, create a key in [Google AI Studio](https://makersuite.google.com/app/apikey).

Set the environment variable in the terminal:
```sh
export GEMINI_API_KEY=<your_api_key>
```

### System commands
The system chat message must begin with an exclamation mark and is used for internal operations.
A short list of supported system commands:

| Command | Description                                                    |
|---------|----------------------------------------------------------------|
| !p      | Select the generative model system prompt <sup>1</sup>         |
| !m      | Select from a list of generative model operations <sup>2</sup> |
| !h      | Select from a list of chat history operations <sup>3</sup>     |
| !i      | Toggle the input mode (single-line <-> multi-line)             |
| !q      | Exit the application                                           |
| !help   | Show system command instructions                               |

<sup>1</sup> System instruction (also known as "system prompt") is a more forceful prompt to the model.
The model will follow instructions more closely than with a standard prompt.
The user must specify system instructions in the [configuration file](#configuration-file).
Note that not all generative models support them.

<sup>2</sup> Model operations:
* Select a generative model from the list of available models
* Show the selected model information

<sup>3</sup> History operations:
* Clear the chat history
* Store the chat history to the configuration file
* Load a chat history record from the configuration file
* Delete all history records from the configuration file

### Configuration file
The application uses a configuration file to store generative model settings and chat history. This file is optional.
If it doesn't exist, the application will attempt to create it using default values. You can use the
[config flag](#cli-help) to specify the location of the configuration file.

An example of basic configuration:
```json
{
  "system_prompts": {
    "Software Engineer": "You are an experienced software engineer.",
    "Technical Writer": "Act as a tech writer. I will provide you with the basic steps of an app functionality, and you will come up with an engaging article on how to do those steps."
  },
  "safety_settings": [
    {
      "category": "HARM_CATEGORY_HARASSMENT",
      "threshold": "LOW"
    },
    {
      "category": "HARM_CATEGORY_HATE_SPEECH",
      "threshold": "LOW"
    },
    {
      "category": "HARM_CATEGORY_SEXUALLY_EXPLICIT",
      "threshold": "LOW"
    },
    {
      "category": "HARM_CATEGORY_DANGEROUS_CONTENT",
      "threshold": "LOW"
    }
  ],
  "tools": [
    {
      "name": "GOOGLE_SEARCH",
      "enabled": true
    },
    {
      "name": "URL_CONTEXT",
      "enabled": true
    }
  ],
  "history": {
  }
}
```
<sup>1</sup> Valid safety settings threshold values include LOW (block more), MEDIUM, HIGH (block less), and OFF.

<sup>2</sup> Upon user request, the `history` map will be populated with records. Note that the chat history is stored
in plain text format. See [history operations](#system-commands) for details.

### CLI help
```console
$ ./gemini -h
Gemini CLI Tool

Usage:
   [flags]

Flags:
  -c, --config string   path to configuration file in JSON format (default "gemini_cli_config.json")
  -h, --help            help for this command
  -m, --model string    generative model name (default "gemini-2.5-flash")
      --multiline       read input as a multi-line string
  -s, --style string    markdown format style (ascii, dark, light, pink, notty, dracula, tokyo-night) (default "auto")
  -t, --term string     multi-line input terminator (default "$")
  -v, --version         version for this command
  -w, --wrap int        line length for response word wrapping (default 80)
```

## License
MIT

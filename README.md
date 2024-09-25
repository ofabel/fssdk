# Flipperzero Script SDK

Looking for a solution to upload your scripts to Flipper Zero without using [qFlipper](https://flipperzero.one/update)?

![Demo](./docs/demo.gif)

## Installation

Just download the [latest release](https://github.com/ofabel/fssdk/releases/latest) and place it somewhere on your computer:
This can be your project's root folder or a more general location like `/usr/local/bin`.

## Usage

This application is a simple CLI program. The program consists of several sub commands:

* **cli** - Opens a terminal session on the Flipper Zero (use <kbd>Ctrl</kbd> + <kbd>C</kbd> to close).
* **run** - Executes only the commands from the `run` section in the config file.
* **sync** - Do a file sync according to the settings in the config file.

If no sub command is provided, the program will perform as follows:

1. **sync** files and folders.
2. **run** the defined commands.

It is not strictly necessary to use a config file.
Small tasks, like uploading a few files and folders, can also be done with CLI arguments.
However, the [config file](#configuration) should assist you on repetitive tasks.

```plain
Usage: fssdk [--config CONFIG] [--quiet] [--port PORT] <command> [<args>]

Options:
  --config CONFIG, -c CONFIG
                         Path to the config file. [default: flipper.json]
  --quiet, -q            Don't print any output. [default: false]
  --port PORT, -p PORT   The port where your Flipper is connected.
  --help, -h             Display this help and exit
  --version              Display version and exit

Commands:
  cli
  run
  sync
```

### Configuration

By default, the application checks for a `flipper.json` file in the current working directory.

```json
{
    "source": "src",
    "target": "/ext/apps/Scripts",
    "orphans": "ignore",
    "include": [
        "*.js",
        "*.py"
    ],
    "exclude": [
        "**.git**",
        "**__pycache__**",
        "*.json"
    ],
    "run": [
        "loader close",
        "js /ext/apps/Scripts/program.js"
    ]
}
```

* **source** - Defines the source folder of your scripts. The path must be relative to the config file location.
* **target** - The target folder of your scripts on the Flipper's SD card. The path must be absolute.
* **orphans** - _optional_ - How to handle orphaned files in the target folder:
    * **ignore** - _default_ - Ignore the files.
    * **download** - Download the files.
    * **delete** - Delete the files.
* **include** - Glob patterns to match included files.
* **exclude** - Glob patterns to match excluded files.
* **run** - Commands to execute. Use `<CTRL+C>` to abort a running command.

> [!TIP]
> The config file should assist you on repetitive tasks.
> It defines all necessary settings, so you only have to run `fssdk` without any arguments.

### Synchronization

File synchronization is one of the key features of this application.
The program will only handle files, that match the `include` patterns and don't match the `exclude` patterns from the config file.
You can overwrite the `source` and/or `target` settings from the config file with the corresponding CLI arguments.

```plain
Usage: fssdk sync [--dry-run] [--force] [--list] [--local] [--source SOURCE] [--target TARGET]

Options:
  --dry-run, -d          Do a dry run, don't upload, download or delete any files. [default: false]
  --force, -f            Upload without any similarity checks. [default: false]
  --list, -l             List matching files. [default: false]
  --local, -o            List matching files from local source only. [default: false]
  --source SOURCE, -s SOURCE
                         Sync all from source to target. If source is a folder, target is also treated as a folder.
  --target TARGET, -t TARGET
                         Sync all from source to target.
```

> [!TIP]
> Use the `--list` argument to check your `include` and `exclude` patterns against the filesystem.

### Run

One of the use cases for this tool is to assist script development with the Flipper Zero:
You write your JS or Python script on your computer and have to upload it to the device in order to test it.
The `run` task executes a list of predefined commands on your Flipper.

```plain
Usage: fssdk run [--dry-run]

Options:
  --dry-run, -d          Do a dry run, don't execute any commands. [default: false]
```

### CLI

Start a terminal session without any additional program like PuTTY or Minicom.
You can also send a single command and receive its output.

```plain
Usage: fssdk cli [--command COMMAND]

Options:
  --command COMMAND, -C COMMAND
                         Execute a single command.
```

## Development

This section only applies to developers or contributors of this repositorys.

### Requirements

* [protoc](https://github.com/protocolbuffers/protobuf/releases)
* [protoc-gen-go](https://protobuf.dev/reference/go/go-generated/)

### Setup

```bash
git clone --recurse-submodules git@github.com:ofabel/flipperzero-script-sdk.git
```

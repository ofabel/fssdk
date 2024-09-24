# Flipperzero Script SDK

Looking for a solution to upload your scripts to Flipper Zero without using [qFlipper](https://flipperzero.one/update)?

![Demo](./docs/demo.gif)

## Usage

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
* **orphans** - How to handle orphaned files in the target folder:
    * **ignore** - Ignore the files.
    * **download** - Download the files.
    * **delete** - Delete the files.
* **include** - Glob patterns to match included files.
* **exclude** - Glob patterns to match excluded files.
* **run** - Commands to execute. Use `<CTRL+C>` to abort a running command.

### Synchronization

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

### Run

```plain
Usage: fssdk run [--dry-run]

Options:
  --dry-run, -d          Do a dry run, don't execute any commands. [default: false]
```

### CLI

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

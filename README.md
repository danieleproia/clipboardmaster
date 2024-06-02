# ClipboardMaster

ClipboardMaster is a clipboard management utility written in Go. It monitors the system clipboard for changes and performs transformations on the clipboard content based on a set of plugins.

## Features

- Monitors the system clipboard for changes
- Performs transformations on the clipboard content based on a set of plugins
- Each plugin can define multiple find-replace operations
- Can be enabled or disabled from the system tray
- Can be set to start at system boot

## Usage

1. Clone the repository
2. Build the project using `go build` (or use the provided powershell script)
3. Run the executable

## Plugins

Plugins are defined in YAML files in the `plugins` directory. Each plugin file should define a name and a list of transformations. Each transformation should define a `find` string and a `replace` string.

Here's an example of a plugin file:

```yaml
name: Example Plugin
replacements:
  - find: "example"
    replace: "test"
  - find: "foo"
    replace: "bar"
```

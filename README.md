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
2. use the provided powershell script, since it copies the required folders in the dist folder too
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

# Roadmap

### Short-term
- Add localization support
- Add settings window

### Long-term
- Add support for more transformation types
- Add support for more clipboard types (currently only text is supported)
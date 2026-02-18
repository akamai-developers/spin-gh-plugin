# Spin GitHub Plugin

[![Release](https://github.com/fermyon/spin-gh-plugin/actions/workflows/release.yaml/badge.svg)](https://github.com/fermyon/spin-gh-plugin/actions/workflows/release.yaml)

This is a plugin that generates GitHub Actions for your Spin Apps.

## Installation

### Install the latest version of the plugin

The latest stable release of the `gh` plugin can be installed like so:

```sh
spin plugins update
spin plugin install gh
```

### Install the canary version of the plugin

The `canary` release of the `gh` represents the most recent commits on `main` and may not be stable, with some features still in progress.

```sh
spin plugins install --url https://github.com/fermyon/spin-gh-plugin/releases/download/canary/gh.json
```

### Install from a local build

Alternatively, use the `spin pluginify` plugin to install from a fresh build. This will use the pluginify manifest (`spin-pluginify.toml`) to package the plugin and proceed to install it:

```sh
spin plugins install pluginify
go build -o gh main.go
spin pluginify --install
```

## Usage

### The `create-action` command

To create a new GitHub Action for your Spin App(s), execute the following command:

```bash
spin gh create-action
```

#### Arguments & Flags

The `create-action` command accepts a bunch of commands that you can use to customize its behavior:

##### GitHub Action Triggers

| Argument | Alias | Description | Default |
| -------- | ----- | ----------- | ------- |
| `ci` | | Run GitHub Action for every push on the specified branch (CI) | `main` |
| `cron` | | Run GitHub Action on a Cron Schedule | |
| `pr` | | Run GitHub Action for every PR targeting the specified branch | |
| `manual` | | Run GitHub Action on manual dispatch | `false` |

##### Tool Versions

You can use the following arguments to customize versions installed as part of the GitHub Action

| Argument | Alias | Description | Default |
| -------- | ----- | ----------- | ------- |
| `spin-version` | | Pin the Spin Version | latest stable release |
| `rust-version` | | Pin the Rust Version | `1.92.0` |
| `go-version` | | Pin the Go Version | `1.25.7` |
| `tinygo-version` | | Pin the TinyGo Version | `0.39.0` |
| `node-version` | | Pin the Node.js Version | `24` |
| `python-version` | | Pin the Python Version | `3.13.0` |

##### General GitHub Action customization

| Argument | Alias | Description | Default |
| -------- | ----- | ----------- | ------- |
| `env` | | Specify Environment Variables (format key=value) | |
| `name` | `n` | Specify the name of the GitHub Action | `CI` |
| `plugin` | `p` | Specify Spin Plugins that should be installed | |
| `os` | | Specify the operating system | `ubuntu-latest` |

#### Render Options

| Argument | Alias | Description | Default |
| -------- | ----- | ----------- | ------- |
| `output` | `o` | Output path | `./github/workflows/ci.yaml` |
| `dry-run` | | Print GitHub Action yaml to `stdout` | `false` |
| `overwrite` | | Overwrite the output if it exists | `false` |
| `template` | `t` | Provide a custom template | |

## The `eject` command

You can export the default template using the `eject` command. Without specifying additional arguments or flags, the template will be written to `stdout`:

```bash
spin gh eject
```

Alternatively, you can write it to a file using the `output` (short `o`) and the `overwrite` flags:

```bash
spin gh eject -o ci.yaml --overwrite
```

## Template Data

The following table lists the data passed to the template as part of the `create-action` command:

| Field | DataType | Description | Sample Value |
| ----- | -------- | ----------- | ------------ |
| `ActionName` | `string` | Name of the workflow | `CI` |
| `OperatingSystem` | `string` | Desired Operating System | `ubuntu-latest` |
| `ActionTriggers` | `ActionTriggers` | Desired triggers | See ActionTriggers section below |
| `EnvironmentVariables` | `[]EnvVar` | Slice with environment variables | (See EnvVar section below) |
| `Tools` | `Tools` | Desired tool versions | see Tools section below |
| `SpinPlugins` | `string` | A comma separated list of Spin plugins | `js2wasm,kube` |
| `Rust` | `bool` | Indicates if any Spin App or Component is built with Rust | `true` |
| `Go` | `bool` | Indicates if any Spin App or Component is built with Go | `true` |
| `JavaScript` | `bool` | Indicates if any Spin App or Component is built with JavaScript | `true` |
| `Python` | `bool` | Indicates if any Spin App or Component is built with Python | `true` |
| `SpinApps` | `[]spinAppTemplateData` | Information for every Spin App discovered | See SpinAppTemplateData section below |

### ActionTriggers

| Field | DataType | Description | Sample Value |
| ----- | -------- | ----------- | ------------ |
| `ManualDispatch` | `bool` | Is `workflow_dispatch` enabled | `true` |
| `Cron` | `string` | trigger cron expression | `0 2 * * *` |
| `PullRequest` | `string` | trigger for PRs targeting the specified branch | `main` |
| `Push` | `string` | trigger for every push on the specified branch | `main` |

### EnvVar

| Field | DataType | Description | Sample Value |
| ----- | -------- | ----------- | ------------ |
| `Key` | `string` | Name of the environment variable | `FOO` |
| `Value` | `string` | Name of the environment variable | `bar` |

### Tools

| Field | DataType | Description | Sample Value |
| ----- | -------- | ----------- | ------------ |
| `Rust` | `string` | Desired Rust Version | `1.92.0` |
| `Go` | `string` | Desired Rust Version | `1.25.7` |
| `TinyGo` | `string` | Desired Rust Version | `0.39.0` |
| `Python` | `string` | Desired Rust Version | `3.13.0` |
| `Node` | `string` | Desired Rust Version | `24` |
| `Spin` | `string` | Desired Rust Version | `3.6.0` |

### SpinAppTemplateData

| Field | DataType | Description | Sample Value |
| ----- | -------- | ----------- | ------------ |
| `Name` | `string` | Name of the Spin App | `spin-app-1` |
| `Path` | `string` | Name of the Spin App | `./src/app1` |
| `Setup` | `string` | Per Spin App setup scripts | `python3 -m venv venv && source venv/bin/activate` |
| `Teardown` | `string` | Name of the Spin App | `deactivate` |
| `Components` | `[]componentTemplateData` | Information for every Component of the App | See ComponentTemplateData section below |

### ComponentTemplateData

| Field | DataType | Description | Sample Value |
| ----- | -------- | ----------- | ------------ |
| `Language` | `string` | Name of the Language used for this component (Rust, Go, JavaScript, Python) App | `Rust` |
| `Path` | `string` | Path of the component | `./src/app1/api` |
| `InstallDependenciesCommand` | `string` | Command used to install component dependencies | `foo` |

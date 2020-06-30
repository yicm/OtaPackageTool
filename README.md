## OtaPackageTool

Archive of the diff files and changelog using git on Linux system.

The ota_packer tool provided in `bin` directory can build two types of package: full and incremental. The tool takes the target-files `.tar` and `.zip` files produced by your project installation files as input, and the project installation file must use the `Git` tool to manage.

## Installation

### Binary Installation

The executable file of `linux-x86-64` version is provided by default and placed in `bin` directory.Therefore, just download to run:

```bash
$ git clone https://github.com/yicm/OtaPackageTool.git
```

### Installing Tool from Source

```bash
$ git clone https://github.com/yicm/OtaPackageTool.git
$ cd OtaPackageTool
$ go build -o bin ./...
```

## Usage

### Preparation

1. Add executable to $PATH
2. Enter the installation file version repository
```bash
$ cd your_installation_file_version_repository
```
3. Run the `ota_packer` command

### Examples

```bash
# --------------------------------------------------
$ ota_packer version
ota_packer version 0.0.1

# --------------------------------------------------
$ ota_packer -h
Archive of the diff files using git on Linux system.

Usage:
  ota_packer [command]

Available Commands:
  gen         Generate package file
  help        Help about any command
  version     Get version of ota_packer

Flags:
  -c, --config string         Config file (default is $HOME/.ota_packer.yaml)
  -h, --help                  help for ota_packer
  -n, --project-name string   Your project name (default "OTA")

Use "ota_packer [command] --help" for more information about a command.

# --------------------------------------------------
$ ota_packer gen -h
Generate a specific version package by entering different configuration parameters.

Usage:
  ota_packer gen [flags]

Flags:
  -F, --diff-filter string       git diff --diff-filter and a similar designation (default "ACMRT")
  -e, --end-commit-id string     End revision (default "HEAD")
  -f, --format string            The format of the archive, supporting zip and tar (default "tar")
  -h, --help                     help for gen
  -o, --output string            Output destination path of the archive
  -p, --prefix string            Prefixed to the filename in the archive while project name is not set. (default "ota_packer")
  -s, --start-commit-id string   Start revision (default "HEAD~1")
  -v, --verbose                  Show packaging process statistics

Global Flags:
  -c, --config string         Config file (default is $HOME/.ota_packer.yaml)
  -n, --project-name string   Your project name (default "OTA")
```

#### Full Updates

If you set `start-commit-id` and `end-commit-id` to the same value, a full upgrade package for current commit id is generated.

```bash
$ ota_packer gen -s HEAD -e HEAD
$ ota_packer gen -s HEAD~1 -e HEAD~1
$ ota_packer gen -s HEAD~3 -e HEAD~1
$ ota_packer gen -s 6bc76a1f -e 6bc76a1f
```

#### Incremental updates

Set `start-commit-id` and `end-commit-id`. If `end-commit-id` datetime is greater than `start-commit-id` datetime, an incremental `upgrade` package is generated, otherwise incremental `downgrade` package is generated.

```bash
$ ota_packer gen -s HEAD~2 -e HEAD~0
$ ota_packer gen -s 6bc76a1f -e 9d31d032

# Set output path as 'tmp' directory, and set project name as 'Test'
$ ./output/ota_packer gen -s HEAD~1 -e HEAD~0 -o tmp -n "Test"
-----------------------------------------------------------------
    Project Name  |  Test
------------------+----------------------------------------------
     Output Path  |  tmp/
------------------+----------------------------------------------
          Output  |  Test-20200630145419-6bc76a1-to-9d31d03.tar
------------------+----------------------------------------------
       Changelog  |  ota_info.json
------------------+----------------------------------------------
{
    "project_name": "Test",
    "last_ota_version": "6bc76a1",
    "ota_version": "9d31d03",
    "is_full_update": false,
    "changes": [
        {
            "type": "D",
            "old_path": "models/y.model",
            "new_path": "models/y.model"
        }
    ]
}
------------------+----------------------------------------------
```

#### About OTA Changelog Status

```txt
A: addition of a file
C: copy of a file into a new one
D: deletion of a file
M: modification of the contents or mode of a file
R: renaming of a file
T: change in the type of the file
```

## Requirements

- Git v2.27.0
- UNIX and UNIX-like operating systems
- Go1.13+ (optional)

## License

Released under the [MIT](https://opensource.org/licenses/MIT) Licence.


# Grove

Grove is a wrapper around the `git worktree` command.

## Features

- [Customizable Hooks](#hooks)
- [Automatic Branch Name Resolution](#branch-name-resolution)
- [Worktree Seeding](#worktree-seeding)

## Installation

```sh
go install github.com/jacobdrury/grove@latest
```

## Usage

```sh
# Initialize grove in your repository
grove init

# Checkout a worktree
grove checkout <branch-name>
```

All other commands are automatically forwarded to `git worktree`.
```sh
grove prune # gets run as 'git worktree prune'
```

## Configuration

The Grove configuration file is located in `.grove/config.yaml` within your repository root.

```yaml
# The directory in which worktrees will be stored.
worktrees-directory: ./worktrees

# Used to resolve branch names
branch-resolver:
    prefix-aliases: {}
    branch-delimiter: /

# Commands to run during different events.
hooks:
    shell: C:\WINDOWS\system32\cmd.exe
    after-checkout: []
```

## Hooks

Grove supports a variety of hooks that will run the listed commands when the corresponding event is triggered. All commands will be run with the configured shell.

```yaml
hooks:
    shell: C:\WINDOWS\system32\cmd.exe
    after-checkout:
        - quick-build
```

The above config will run the `quick-build` within the new worktree directory after it's been checked out.

## Worktree Seeding

In the `.grove` directory you will find a `seed` directory. This directory contains files that you wish to seed new worktrees with when they are created. The directory structure found within the `seed` directory will be maintained when the worktree is seeded.

## Branch Name Resolution

Branch names can be resolved using custom 'prefix aliases' configured in `.grove/config.yaml`.

```yaml
branch-resolver:
    prefix-aliases:
        f: feature
        c: chore
    branch-delimiter: /
```

With the above configuration, the following would be true:
```sh
grove checkout f/some-feature # would checkout `feature/some-feature`
```

This will also perform slug resolution.

```sh
git checkout -b feature/1234-example # create a new branch
grove checkout f/1234                # resolves to `feature/1234-example`
```
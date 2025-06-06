# Grove

Grove is a wrapper around the `git worktree` command.

## Features

- After checkout hook
- Branch name resolution
- Worktree seeding

## Usage

```sh
# Initialize grove in your repository
$ grove init

# Checkout a worktree
$ grove checkout <branch-name>
```

All other commands are automatically forwarded to `git worktree`.
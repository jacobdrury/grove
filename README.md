# git worktree manager

```
go install .../wt@latest
```

```bash
wt init # Run from repo root
    # Create .wt directory
        # - config
        # - seed_files/
    # Create worktrees directory
    # Adds both to local .gitignore
```

`.wt/config`
```yaml
worktree:
    main: ./worktrees/main
    branch_prefix_aliases:
        jolondirenko-smith: j
    branch_name_format: ^{prefix}/{ticket}-.*$
```

```
wt checkout j/fm-3112 [fuckit]
```

```sh
wt rebase # called from within a worktree
# pushd to main directory
# checkout main
# pull
# popd to worktree directory
# rebase
```

```
cd ./worktrees/some-worktree
wt rebase
`.wt`, cd .., `.wt`
```

```
wt checkout j/fm-3331

slug_prefix_pattern = (fm|FM)\-[0-9]+\-

jacob/fm-3331-fuck-idk <-- matches
jacob/fm-3321-some-ther
```
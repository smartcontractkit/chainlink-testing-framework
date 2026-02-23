# Compatibility Testing

We have a simple tool to check compatibility for CL nodes. The example command will filter and sort the available tags, rollback and install the oldest version, and then begin performing automatic upgrades to verify that each subsequent version remains compatible with the previous one.

```bash
ctf compat backward \
--buildcmd "just cli" \
--envcmd "cl r" \
--testcmd "cl test ocr2 TestSmoke/rounds" \
--include_tags +compat \
--nodes 3 \
--versions_back 3
```

Use `ctf compat restore` to rollback to current branch (default is `develop`)

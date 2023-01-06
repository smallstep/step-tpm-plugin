# step-tpm-plugin

🔑 `step` plugin for interacting with TPMs. 

> ⚠️ This tool is currently in beta mode and its usage might change without
> announcements.

## TODO

Incomplete lists of things still to do:

- Cleanup
- Code docs
- Tests
- Nicer outputs (e.g. JSON)
- Extract the `tpm` package into `smallstep/crypto`
- Add `--verbose` for debugging output (incl. interactions with TPM)
- Handle/transform TPM level errors
    - Example: `NCryptCreatePersistedKey returned 8009000F: The operation completed successfully.`
- More functionalities
# Roadmap

This file is the worklist for the automated implementation loop. Each run picks the
**first unchecked item**, implements it, checks the box, commits, and pushes.

## How the loop processes this file
1. Find the first item with an unchecked box (`- [ ]`).
2. Implement it on branch `claude/happy-shannon-eo28r9`.
3. Verify it (build/run/tests as appropriate).
4. Tick the box (`- [x]`), commit with a descriptive message, push, update the draft PR.
5. Stop. The next run continues with the new first unchecked item.

If an item is ambiguous or needs a decision, leave it unchecked, add a note under it,
and ask the user instead of guessing.

## Items
- [x] Refactor `report.cbl` to use a structured `PERFORM` loop instead of `GO TO`
- [ ] Add a Go test (`report_test.go`) that runs the sample CSV and asserts the fixed-width output
- [ ] Add a Rust test for the same sample-input/expected-output behavior
- [ ] Update `README.md` to document the Rust port (currently only mentions COBOL and Go)
- [ ] Add a GitHub Actions CI workflow that builds Go + Rust and runs both test suites
- [ ] Add input validation + clear error messages for malformed CSV rows (Go and Rust)

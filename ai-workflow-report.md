# AI Workflow Report

## Tools and models

**Claude Code** with **Claude Sonnet 4.6** — used throughout: spec analysis, architecture planning, implementation, test generation, and tooling.

---

## How the work was planned

The session started by feeding both PDFs to Claude Code and asking it to discuss the architecture and show a filesystem shape before writing any code. This produced the layered package structure (`statemachine → store → engine → handler → cmd`) and surfaced the three ambiguities that needed decisions before implementation: the retry count contradiction, the SUBSCRIBER_STATE open decision, and whether SYSTEM_ERROR consumes a retry slot. Those decisions were written to `DECISIONS.md` before a single line of production code was written.

Implementation then proceeded in planning loops: generate a plan for the next chunk, implement it, review what was done and what remained, repeat.

---

## Two prompts that actually mattered

**1. "Let's build a plan according to the specs, let's discuss architecture, show me the shape of the filesystem"**

This forced Claude Code to think before coding. The resulting file layout — with `statemachine` as a pure generic package, `engine` holding all business logic, and `handler` knowing nothing about billing — came out of that conversation. The separation held cleanly through the entire project.

**2. The coverage tooling loop**

Mid-project, after noticing coverage holes with no line-level feedback in the terminal (Go's standard tooling only shows line numbers in HTML reports), a prompt asking Claude Code to implement a terminal coverage tool triggered an unexpected chain. Building the tool required a strict, predictable test structure — which led to porting the `supertape` assertion style (a battle-tested JS library) to Go as `go-tape`. While implementing `go-tape`'s test runner, it became clear that Go's test event stream — *current state + event → next state* — is structurally identical to the subscription lifecycle machine. The generic `statemachine` package was extracted at that point and published separately so both the test runner and the subscriber could use the same implementation. The subscriber's `internal/statemachine` is that package.

This was not planned. It came from following a real friction point to its root.

---

## Where the AI got it wrong

Claude Code did not write tests on the first implementation pass. The pattern was: implement a package, move to the next, implement that. Coverage holes only became visible after stepping back and reviewing. This required explicit re-planning loops — "now add tests for engine", "now add tests for handler" — that would not have been necessary under a strict TDD discipline.

The fix was structural: adopt `go-tape`'s one-assertion-per-subtest style, which makes missing coverage obvious by making each case a named, isolated test. Once that pattern was in place, gaps were easy to spot and fill.

A secondary observation: Claude Code generated `Hook()` and `Validate()` methods on the state machine that the subscriber does not call. In this case the output was evaluated and kept — they belong in a general-purpose library that other projects will pull in, and they are tested within `go-tape`. The point is that the output was reviewed and a conscious decision was made, rather than accepted or discarded automatically.

---

## On AI use generally

The expectation in this assignment was to use AI heavily and be evaluated on how it was directed, verified, and constrained. The verification piece mattered most here. Growing an application with AI requires that every generated chunk be checkable — which is why the tooling investment (coverage reporting, strict assertion style) came before the full implementation rather than after. A state machine with 9 rows and a typed transition table is also easier to verify than imperative if/else chains: the whole lifecycle fits on one screen and matches the spec's own table row-for-row.

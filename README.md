markdown-report
===============

[![Actions Status](https://github.com/julias-shaw/gauge-markdown-report/workflows/test/badge.svg)](https://github.com/julias-shaw/gauge-markdown-report/actions)

A [Gauge](https://gauge.org) execution-reporting plugin that emits
GitHub-Flavored Markdown. Forked from the `html-report` plugin, but
the output is plain `.md` so reports render in GitHub, GitLab, IDE
previewers, `glow`, and any other Markdown viewer.

Features
--------

- One `index.md` summary at the report root, one `.md` per spec under
  `specs/`, screenshots copied into `images/`.
- Status glyphs (✅ / ❌ / ⏭️) plus text status words for accessible
  rendering across viewers.
- Stack traces and large error blocks wrapped in collapsible
  `<details>` sections — keeps files scannable while preserving the
  full diagnostic context.
- Table-driven specs and scenario-table-driven scenarios render as
  GFM tables with per-row pass/fail status.
- Honors `use_nested_specs` to emit a per-directory `index.md`.

Installation
------------

This fork is not registered in the [gauge-repository](https://github.com/getgauge/gauge-repository), so the bare `gauge install markdown-report` command will not resolve. Install from a release asset instead:

Download the appropriate zip for your platform from [Releases](https://github.com/julias-shaw/gauge-markdown-report/releases) and run:

```
gauge install markdown-report --file markdown-report-0.1.0-linux.x86_64.zip
```

Substitute `linux.x86_64` with `linux.arm64`, `darwin.x86_64`, `darwin.arm64`, `windows.x86_64`, or `windows.arm64` as appropriate.

#### Build from source

Requirements:
- [Go](https://golang.org/) (matching `go.mod`)

```
go run build/make.go                 # compile for current host
go run build/make.go --all-platforms # cross-compile
go run build/make.go --install       # install locally
go run build/make.go --distro        # build distributable
```

`go run build/make.go --install --plugin-prefix CUSTOM_LOCATION` installs into a custom prefix.

Configuration
-------------

The plugin reads its configuration from `env/default/default.properties`
in the project, which Gauge surfaces as environment variables. The
configurable properties are:

| Property                  | Type    | Default     | Purpose |
| ------------------------- | ------- | ----------- | ------- |
| `gauge_reports_dir`       | path    | `reports`   | Where reports are written. Relative paths are resolved against the project root. |
| `overwrite_reports`       | boolean | `true`      | If `false`, every run lands in a new timestamped directory. |
| `use_nested_specs`        | boolean | `false`     | When `true`, emits an `index.md` for every spec subdirectory in addition to the suite-level one. |
| `save_execution_result`   | boolean | `true`      | If `true`, the plugin places a symlink (or `.bat` on Windows) to its executable inside the report directory so the report can be regenerated offline. |
| `gauge_screenshots_dir`   | path    | unset       | Where Gauge stages step / hook screenshots. Set by Gauge; the plugin reads the variable. |

Report regeneration
-------------------

The plugin saves the raw `last_run_result` proto under `.gauge/`
during normal runs. To regenerate a report from it:

1. Navigate to the report directory.
2. Move the `markdown-report` symlink to the `.gauge/` directory.
3. Navigate to `.gauge/`.
4. Run:
   ```
   ./markdown-report --input=last_run_result --output=/some/path
   ```

The output directory is created if it doesn't exist; do not point this
at a directory you don't want overwritten. Regeneration only works if
`save_execution_result` was `true` when the original run captured the
result.

Sample output
-------------

`index.md` (excerpt):

```markdown
# Gauge Report — demo

_Generated Jan 2, 2026 at 3:04pm · Environment: default · Tags: regression_

## Summary

| Scope     | Total | ✅ Passed | ❌ Failed | ⏭️ Skipped | Success rate |
| ---       | ---   | ---       | ---       | ---         | ---          |
| Specs     | 1     | 0         | 1         | 0           | 0%           |
| Scenarios | 2     | 1         | 1         | 0           | 50%          |

**Total time:** 6.00s

## Specs

| Status | Spec                            | Time  | Tags       |
| ---    | ---                             | ---   | ---        |
| ❌     | [Checkout](specs/checkout.md)   | 6.00s | regression |
```

A failing-step block in a per-spec page:

```markdown
- ❌ checkout breaks _(00:00:00)_

  **Error:** `expected 200 got 500`

  <details><summary>Stack trace</summary>

  ```
  at handler.go:42
  at router.go:11
  ```

  </details>

  ![Failure screenshot](../images/scn-fail.png)
```

The full set of representative shapes (mixed pass/fail, all-skipped,
before-suite-hook-failure, data-table-driven, with-screenshots) is
checked in under [`mdgen/_testdata/`](mdgen/_testdata/) — those files
are the format reference.

License
-------

[Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0.txt)

Copyright
---------

Copyright 2026 Julias Shaw

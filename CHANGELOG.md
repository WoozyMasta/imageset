# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog][],
and this project adheres to [Semantic Versioning][].

<!--
## Unreleased

### Added
### Changed
### Removed
-->

## [0.3.0][] - 2026-03-29

### Changed

* lint flow is now lintkit-first and uses `lint.Diagnostic` end-to-end
* `ValidationError.Diagnostics` now stores `[]lint.Diagnostic`
* `AttachLintDiagnostics` now accepts `[]lint.Diagnostic`

### Removed

* internal imageset diagnostic adapter layer (`diagnostic_lint.go`)

## [0.2.0][] - 2026-03-29

### Added

* `ValidateWithOptions(...)` and `ValidateOptions`
  for configurable validation behavior
  (including optional image padding checks)
* `lintkit` integration surface
  for rule registration and diagnostics catalog export.

### Changed

* Parser is strict for unknown/context-invalid fields

[0.2.0]: https://github.com/WoozyMasta/imageset/compare/v0.1.0...v0.2.0
[0.3.0]: https://github.com/WoozyMasta/imageset/compare/v0.2.0...v0.3.0

## [0.1.0][] - 2026-xx-xx

### Added

* First public release

[0.1.0]: https://github.com/WoozyMasta/imageset/tree/v0.1.0

<!--links-->
[Keep a Changelog]: https://keepachangelog.com/en/1.1.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

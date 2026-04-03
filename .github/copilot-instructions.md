# Copilot Skill: README-vs-Code Documentation Audit

Cross-check a repository's README (or other user-facing documentation) against the actual source code, tests, and examples. Identify every inconsistency and produce a structured issue report for each finding.

## Methodology

1. **Parse the README section by section.** For each code example or behavioral claim:
   a. Locate the corresponding source file(s) and function(s).
   b. Run or trace the code path to verify the documented behavior.
   c. Check existing tests for confirmation or contradiction.
   d. Check the examples directory (if present) for working counterexamples.

2. **Categories to check:**
   - **Code examples** — Do they compile? Do they produce the output shown?
   - **Error messages** — Does the documented error string match what the code actually returns?
   - **Default values** — Do documented defaults match the values in constructors/factory functions?
   - **Output formatting** — Do documented outputs (tables, text, YAML, JSON) match the actual rendering code?
   - **Feature descriptions** — Does the prose accurately describe what the code does, including edge cases?
   - **Struct tags / annotations** — Are tag formats and their effects described correctly?
   - **API surface** — Are all public functions, methods, and types documented? Are any deprecated items still shown?

3. **For each inconsistency found, create a report entry with:**
   - **Title:** A concise summary of the mismatch
   - **Labels:** `documentation` (always), plus `bug` if the code behavior itself is arguably wrong, or `enhancement` if the fix would require a new feature
   - **Description** containing:
     - What the README says (with line numbers or section references)
     - What the code actually does (with file and function references)
     - A concrete example showing the discrepancy
     - **Impact** — how this affects users (misleading, broken copy-paste, runtime errors, cosmetic)

4. **Severity guidance:**
   - **High:** Users copying the example will get wrong behavior or runtime errors
   - **Medium:** Misleading documentation that could cause confusion
   - **Low:** Cosmetic or minor context issues

5. **Do not:**
   - Report stylistic preferences or subjective improvements
   - Flag issues that are clearly documented as known limitations

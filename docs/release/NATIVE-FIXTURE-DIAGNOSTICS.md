# Bounded diagnostics for the fixed native test fixture

The harness emits a structured record immediately before its existing Linux
terminal failure. This changes diagnostic output only: production process
control, limits, cleanup, fixture arguments and rejection behavior are unchanged.

Schema `pscan-native-fixture-diagnostic-v2` deliberately replaces the insufficient
v1 stderr classifier with an actual escaped excerpt from this fixed synthetic
fixture's captured stderr. It is not a general logger for scanner inputs or
arbitrary commands. Stdout content and environment variables are never included.

- Entire JSON record: at most 2,048 UTF-8 bytes.
- Stderr excerpt: at most the first 256 source bytes.
- Total stderr bytes, retained source-byte count and truncation are explicit.
- ASCII bytes 32–126 remain literal, except backslash becomes two backslashes.
  All other bytes become lowercase `\xhh`. The JSON serializer then applies JSON
  escaping. Interpret the excerpt after JSON decoding using these byte-escape
  rules. Literal `\x` text remains distinguishable from an escaped byte.
- No UTF-8 decoding or replacement occurs in the excerpt. Invalid UTF-8 and a
  truncated multibyte sequence remain unambiguous original bytes.
- Missing/null/invalid values are marked unavailable, with null count/truncation
  where their values are unknown. An actual empty byte array is distinct.
- Case names and terminal reasons remain closed sets. Unknown terminal text is
  not emitted. Typed exit and containment metadata are retained without defaulting
  missing values to zero or success.

`test-native-fixture-diagnostics.ps1` validates formatting using inert result
records, including escaping, exact 256-byte boundaries, longer streams, invalid
and split UTF-8, unavailable data, inappropriate types and property getters.
These local tests neither execute a native fixture nor establish its root cause.

The actual v1 diagnostic run 35237014717 failed with 66 unclassified stderr bytes.
The v2 excerpt is intended to expose that controlled launcher error on a distinct
reviewed diagnostic invocation. Prior runs, frozen source and the R6 tag are
preserved; this documentation does not declare the runtime issue fixed.

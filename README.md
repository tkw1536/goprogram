# goprogram

![CI Status](https://github.com/tkw1536/goprogram/workflows/CI/badge.svg)

A go >= 1.19.2 package to create programs, originally designed for [ggman](https://github.com/tkw1536/ggman).

## Changelog

# 0.2.1 (Upcoming)

- add `stream.NonInteractive` method
- add `WriterGroup` to status
- `status` optimizations when there is no progress to be written

# 0.2.0 (Released [Nov 27 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.2.0))

- extend context handling
- add additional type parameters to `collection/slice`

# 0.1.1 (Released [Oct 7 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.1.1))

- remove memory leak from `slices.Filter` got types which need garbage collection

# 0.1.0 (Released [Oct 6 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.1.0))

- add `collection` utility package
- add stream package to `status`
- improve `docfmt` package
- various internal improvements

# 0.0.17 (Released [Sep 30 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.17))

- add compatibility mode to `status`

# 0.0.16 (Released [Sep 30 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.16))

- promote `status` to top-level package
- add more utility functions to `stream`

# 0.0.15 (Released [Sep 29 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.15))

- update and document `status` package

# 0.0.14 (Released [Sep 22 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.14))

- add `status` package

# 0.0.13 (Released [Sep 22 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.13))

- add `Print` and `EPrint` methods to `stream`
- add `FromDebug` method to stream
- minor internal changes

# 0.0.12 (Released [Sep 15 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.12))

- add `Streams` and `NonInteractive` utility methods to `stream`

# 0.0.11 (Released [Sep 7 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.11))

- add `program.Exec` method to execute a command from within a command

# 0.0.10 (Released [Sep 5 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.10))

- add error wrapping

# 0.0.9 (Released [Aug 26 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.9))

- add `BeforeKeyword`, `BeforeAlias`, `BeforeCommand` hooks
- add a method `StdinIsATerminal` to check if stdin is a terminal

# 0.0.8 (Released [Aug 17 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.8))

- add `ReadLine`, `ReadPassword` and `ReadPasswordStrict` methods
- minor `go1.19` formatting

# 0.0.7 (Released [May 2 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.7))

- remove `BeforeRegister` method and pass program in context
- copy commands before executing and make sure they become pointers

# 0.0.6 (Released [Apr 28 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.6))

- extend doccheck package into docfmt package

# 0.0.5 (Released [Apr 18 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.5))

- add doccheck package

# 0.0.4 (Released [Apr 15 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.4))

- refactor argument parsing

# 0.0.3 (Released [Mar 29 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.3))

- add `EmptyRequirement` struct

# 0.0.2 (Released [Mar 16 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.2))

- use `golang.org/x/exp/slices` package
- name type parameters consistently

# 0.0.1 (Released [Mar 9 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.0.1))

- initial release
# goprogram

![CI Status](https://github.com/tkw1536/goprogram/workflows/CI/badge.svg)

A go 1.18 package to create programs, originally designed for [ggman](https://github.com/tkw1536/ggman).
Documentation is a work in progress.

## Changelog

# 0.0.13 (Upcoming)

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
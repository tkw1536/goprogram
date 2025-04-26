# goprogram

![CI Status](https://github.com/tkw1536/goprogram/workflows/CI/badge.svg)

A go >= 1.24.2 package to create programs, originally designed for [ggman](https://github.com/tkw1536/ggman).

## Changelog

# 0.9.0 (Released [Apr 25 2025](https://github.com/tkw1536/goprogram/releases/tag/v0.9.0))

- rework 'exit.Error' functionality

# 0.8.4 (Released [Apr 23 2025](https://github.com/tkw1536/goprogram/releases/tag/v0.8.4))

- add 'exit.Code' function

# 0.8.3 (Released [Apr 18 2025](https://github.com/tkw1536/goprogram/releases/tag/v0.8.3))

- update dependencies

# 0.8.2 (Released [Apr 11 2025](https://github.com/tkw1536/goprogram/releases/tag/v0.8.2))

- fix broken linter dependencies

# 0.8.1 (Released [Apr 11 2025](https://github.com/tkw1536/goprogram/releases/tag/v0.8.1))

- run modernize analyzer

# 0.8.0 (Released [Apr 9 2025](https://github.com/tkw1536/goprogram/releases/tag/v0.8.0))

- update to go1.24.1
- make exit package work well with wrapped functions
- so much linting
- remove a couple of deprecated functions

# 0.7.1 (Released [Apr 3 2025](https://github.com/tkw1536/goprogram/releases/tag/v0.7.1))

- update `pkglib` version

# 0.7.0 (Released [Dec 3 2024](https://github.com/tkw1536/goprogram/releases/tag/v0.7.0))

- use builtin package instead of `golang.org/x/exp/slices`
- update `shellescape` dependency with new import path

# 0.6.0 (Released [Nov 25 2024](https://github.com/tkw1536/goprogram/releases/tag/v0.6.0))

- run a linter and fix issues
- update to go 1.23.3

# 0.5.1 (Released [May 3 2024](https://github.com/tkw1536/goprogram/releases/tag/v0.5.1))

- fix typos to make the spellchecker happy 
- upgrade dependencies

# 0.5.0 (Released [Oct 1 2023](https://github.com/tkw1536/goprogram/releases/tag/v0.5.0))

- remove automatic wrapping 

# 0.4.1 (Released [Jul 19 2023](https://github.com/tkw1536/goprogram/releases/tag/v0.4.1))

- update to new pkglib

# 0.4.0 (Released [May 10 2023](https://github.com/tkw1536/goprogram/releases/tag/v0.4.0))

- introduce `WrapError` function
- update dependencies
- fix a lot of documentation typos

# 0.3.5 (Released [Mar 16 2023](https://github.com/tkw1536/goprogram/releases/tag/v0.3.5))

- update dependencies

# 0.3.4 (Released [Mar 15 2023](https://github.com/tkw1536/goprogram/releases/tag/v0.3.4))

- move `stream` and `status` packages to `pkglib`

# 0.3.3 (Released [Mar 9 2023](https://github.com/tkw1536/goprogram/releases/tag/v0.3.3))

- update dependencies

# 0.3.2 (Released [Mar 9 2023](https://github.com/tkw1536/goprogram/releases/tag/v0.3.2))

- add `GOPROGRAM_ERRLINT_EXCEPTIONS` to `cmd/errlint`

# 0.3.1 (Released [Mar 9 2023](https://github.com/tkw1536/goprogram/releases/tag/v0.3.1))

- add `Error.DeferWrap` function

# 0.3.0 (Released [Feb 24 2023](https://github.com/tkw1536/goprogram/releases/tag/v0.3.0))

- move utility packages to pkglib
- updated errlint command

# 0.2.4 (Released [Dec 7 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.2.3))

- CI: Run `errlint` automatically
- quoted word validation bugfix

# 0.2.3 (Released [Dec 7 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.2.3))

- add `cmd/errlint` static checker
- add `IsNullWriter` function
- add choices of options to help page

# 0.2.2 (Released [Dec 2 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.2.2))

- allow accessing full context object from simple context

# 0.2.1 (Released [Nov 30 2022](https://github.com/tkw1536/goprogram/releases/tag/v0.2.1))

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
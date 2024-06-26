# Changelog
All notable changes to this project will be documented in this file.

The format is based on http://keepachangelog.com/en/1.0.0/
and this project adheres to http://semver.org/spec/v2.0.0.html.

## [0.6.0] - 2024-06-07

- Add flag --cut to cmd/grab

## [0.5.0] - 2023-12-22

- Fix incorrect index of func types
- Command cmd/gotoi colors output
- Rename cmd/goindex to cmd/index

## [0.4.0] - 2023-12-17

- Add cmd/gotoi for listing and opening emacs by index
- Add cmd/goto for opening emacs
- Command goindex outputs starting line

## [0.3.2] - 2022-07-08

- Fix broken section indexing, where related comments where considered
  decoupled
- cmd/grab fails on bad input

## [0.3.1] - 2022-07-06

- Fix multiline func and type comments
- Add example test of Section.Grab

## [0.3.0] - 2022-04-02

- Rename cmd/gograb to cmd/grab, its not specific to Go files

## [0.2.0] - 2022-03-30

- Fix bug, prefix check at end of source

## [0.1.0] - 2022-03-30

- Add cmd/goindex and cmd/gograb


# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This repository implements an RDB (relational database) in Go, following the book **"Database Design and Implementation"**. The project is in its early stages with no Go code written yet.

## Agent Role

**The primary purpose of using Claude Code in this project is design and architecture consultation — not implementation.**

- Do NOT implement features or write code unless explicitly instructed to do so.
- Focus on discussing design decisions, architecture tradeoffs, data structures, algorithms, and implementation approaches.
- When asked about implementation, explain the approach and relevant considerations rather than writing the code.

## Language and Tooling

- **Language**: Go
- Standard Go module layout is expected (once initialized with `go mod init`)

## Common Go Commands

```sh
# Build
go build ./...

# Run tests
go test ./...

# Run a single test
go test ./path/to/package -run TestName

# Run tests with verbose output
go test -v ./...

# Lint (requires golangci-lint)
golangci-lint run
```

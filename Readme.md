# Golang Integration Test

## Introduction

Here is an example for how to run Integration tests with Golang


## Contents

- [Dependencies](#dependencies)
- [Build](#build)
- [Test](#test)

## Dependencies

- Docker
- Docker Compose

```bash
# Install psql, golang
brew install postgres go
```

## Build

```bash
./build.sh
```

## Test

Run tests in `test` dir:

```bash
./test.sh
```

This will run tests via `curl` and `go test ./test`

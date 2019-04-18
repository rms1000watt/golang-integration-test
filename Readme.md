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

Run tests in `test` dir managed by `go`:

```bash
./test.sh
```

Run tests via curl commands:

```bash
./test.sh curl
```

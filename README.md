<img align="right" width="300" src="https://github.com/adrianbrad/psql-docker-tests-example/blob/image-data/logo.png?raw=true" alt="adrianbrad psqldocker">

# 📊 psql-docker-tests-example

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/adrianbrad/psql-docker-tests-example)](https://github.com/adrianbrad/psql-docker-tests-example)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/adrianbrad/psql-docker-tests-example)

[![CodeFactor](https://www.codefactor.io/repository/github/adrianbrad/psql-docker-tests-example/badge)](https://www.codefactor.io/repository/github/adrianbrad/psql-docker-tests-example)
[![Go Report Card](https://goreportcard.com/badge/github.com/adrianbrad/psql-docker-tests-example)](https://goreportcard.com/report/github.com/adrianbrad/psql-docker-tests-example)
[![lint-test](https://github.com/adrianbrad/psql-docker-tests-example/workflows/lint-test/badge.svg)](https://github.com/adrianbrad/psql-docker-tests-example/actions?query=workflow%3Alint-test)
[![codecov](https://codecov.io/gh/adrianbrad/psql-docker-tests-example/branch/main/graph/badge.svg)](https://codecov.io/gh/adrianbrad/psql-docker-tests-example)
---
### Parallel black box PostgreSQL unit tests run against a real database.

Consider reading the [Medium Story](https://adrianbrad.medium.com/parallel-postgresql-tests-go-docker-6fb51c016796) first.

This package provides examples on how to run PostgreSQL units tests against a real database
with every tests running in a separate SQL transaction. You can find the tests in [this](https://github.com/adrianbrad/psql-docker-tests-example/tree/main/internal/psql) package.

The PostgreSQL database is started using the https://github.com/adrianbrad/psqldocker package.

The sql connections are opened in an isolated SQL transaction using the https://github.com/adrianbrad/psqltest package.
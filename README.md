# goadesignupgrader

A tool to upgrade a design definition for [Goa](https://github.com/goadesign/goa) from v1 to v3

[![GoDoc](https://godoc.org/github.com/goadesign/goadesignupgrader?status.svg)](https://godoc.org/github.com/goadesign/goadesignupgrader) [![GitHub Actions](https://github.com/goadesign/goadesignupgrader/workflows/Go/badge.svg)](https://github.com/goadesign/goadesignupgrader/actions) ![GitHub](https://img.shields.io/github/license/goadesign/goadesignupgrader)

## Installation

```sh
$ go get github.com/goadesign/goadesignupgrader/...
```

## Usage

```sh
$ goadesignupgrader [design package]
```

You can use `-fix` flag to apply all suggested fixes.

```sh
$ goadesignupgrader -fix [design package]
```

It's recommended to use together with gormt.

```sh
$ goadesignupgrader -fix [design package]; gofmt -s -w .
```

## Supported diagnostics

* Import declarations
* HTTP status constants

### Supported DataTypes

* `DateTime`
* `Integer`

### Supported DSLs

* `Action`
* `BasePath`
* `CONNECT`
* `CanonicalActionName`
* `Consumes`
* `DELETE`
* `DefaultMedia`
* `GET`
* `HEAD`
* `HashOf`
* `Headers`
* `Media`
* `MediaType`
* `Metadata`
* `OPTIONS`
* `PATCH`
* `POST`
* `PUT`
* `Params`
* `Parent`
* `Produces`
* `Resource`
* `Response`
* `Routing`
* `Status`
* `TRACE`

## License

[MIT License](LICENSE)

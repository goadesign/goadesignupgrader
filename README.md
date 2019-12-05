# goadesignupgrader

A tool to upgrade a design definition for [Goa](https://github.com/goadesign/goa) from v1 to v3

[![GoDoc](https://godoc.org/github.com/tchssk/goadesignupgrader?status.svg)](https://godoc.org/github.com/tchssk/goadesignupgrader) [![CircleCI](https://circleci.com/gh/tchssk/goadesignupgrader.svg?style=shield&circle-token=736c8b4099ed93ee5f3ad19330c0751df6b86ad4)](https://circleci.com/gh/tchssk/goadesignupgrader) ![GitHub](https://img.shields.io/github/license/tchssk/goadesignupgrader)

## Installation

```sh
$ go get github.com/tchssk/goadesignupgrader/...
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
$ goadesignupgrader -fix [design package] | gofmt -s -w
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
* `Produces`
* `Resource`
* `Response`
* `Routing`
* `Status`
* `TRACE`

## License

[MIT License](LICENSE)

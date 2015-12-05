# Go URL Valuer

A generic way of creating a url.Values object from any struct. Uses the `go generate` command from Go 1.4. See [this post for details](http://blog.golang.org/generate).

## Installation

Install this package:

```
go install github.com/parkr/urlvaluer
```

## Usage

Add a `go:generate` comment to the top of any of your go files:

```go
//go:generate urlvaluer $GOFILE
```

Then run

```go
go generate ./path/to/that/file
```

Don't add a space after the `//` and before the `go:generate` or it'll break. It's a Go thing.

## What It Does

`urlvaluer` generates a set of files in the same directory with the generated code – one for each input file which contains the `//go:generate` comment. If your types already have a `.UrlValues()` function in the same file then it won't auto-generate the `.UrlValues()` method.

```go
func (t MyData) UrlValues() *url.Values {
    vals := &url.Values{}
    if t.MyField != nil {
    }
    return vals
}
```

To use the urlvaluer, just do:

```go
myData.UrlValues()
```

## About

By [Parker Moore](https://byparker.com), based on [`joiner`](https://github.com/bslatkin/joiner) by [Brett Slatkin](http://www.onebigfluke.com).

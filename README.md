[![Build Status](https://img.shields.io/travis/walle/wiki.svg?style=flat)](https://travis-ci.org/walle/wiki)
[![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/walle/wiki)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/walle/wiki/master/LICENSE)
[![Go Report Card](http://goreportcard.com/badge/walle/wiki?t=3)](http:/goreportcard.com/report/walle/wiki)

# wiki

Command line tool to get Wikipedia summaries.

The tool can fetch summaries from any MediaWiki wiki with the API active, but
defaults to the English Wikipedia.

## Installation

To be able to install with `go get` requires you to have your `$GOPATH` setup
and your `$GOPATH/bin` added to path as described here
http://golang.org/doc/code.html#GOPATH.

If you don't want the man file you can just install it with `go get`.

```shell
$ go get github.com/walle/wiki/cmd/wiki
```

If you want to install the man file you can install with `go get`,
but then use the `make install` command from the source directory.

```shell
$ go get github.com/walle/wiki
$ cd $GOPATH/src/github.com/walle/wiki
$ make install
```

or just copy the man file in `_doc/wiki.1` to `/usr/local/share/man/man1` or
where you keep your man files.

### Dependencies

* go-colorable https://github.com/mattn/go-colorable

## Usage

To get a summary from Wikipedia in English just invoke the tool with a query.

```shell
$ wiki golang
Redirected from Golang to Go (programming language)

Go, also commonly referred to as golang, is a programming language developed at Google in 2007 by Robert Griesemer, Rob Pike, and Ken Thompson. It is a statically typed language with syntax loosely derived from that of C, adding garbage collection, type safety, some structural typing capabilities, additional built-in types such as variable-length arrays & key-value maps, and a large standard library.
The language was announced in November 2009 and is now used in some of Google's production systems. Go's "gc" compiler targets the Linux, OS X, FreeBSD, NetBSD, OpenBSD, Plan 9, DragonFly BSD, Solaris, and Windows operating systems and the i386, Amd64, ARM and IBM POWER processor architectures. A second compiler, gccgo, is a GCC frontend.
Android support was added in version 1.4, which has since been ported to also run on iOS.

Read more: https://en.wikipedia.org/wiki/Go_(programming_language)
```

To get a localized result, e.g. in Swedish use the -l flag.

```shell
$ wiki -l sv ruby
```

Use the -h flag to see all options (or `man wiki` if you have it installed)

```shell
$wiki -h
wiki is a tool used to fetch exerpts from wikipedia
Usage: wiki [options...] query
Options:

  -h    Print help information and exit.
  -l string
        The language to use (default "en")
  -n    If the output should not be colorized
  -no-check-certificate
        Skip verification of certificates
  -s    If simple output should be used
  -short
        If short output should be used
  -u string
        The api url (default "https://%s.wikipedia.org/w/api.php")
  -version
        Print version information and exit.
```

### Use another wiki

To get excerpts from another wiki use the -u flag to give another url to the
API to use.

```shell
$ wiki -u https://en.wikiversity.org/w/api.php physics
```

This gives the excerpt from the wiki at wikiversity.org instead of Wikipedia.

#### Advice

If you frequently use the tool to fetch data from a custom url, add an alias.
E.g. for bash. Add an alias to your `.bash_profile` or `.bashrc` file.

```bash
alias uwiki='wiki -u https://en.wikiversity.org/w/api.php '
```

And call it using

```shell
$ uwiki physics
```

## Testing

Run the tests using the make target test `make test`, this runs both the unit and 
the integration tests. For running only one type of tests use `go test -cover` 
and `./integration-tests.sh` respectively.

```shell
$ make test
```

## Contributing

All contributions are welcome! See [CONTRIBUTING](CONTRIBUTING.md) for more
info.

## License

The code is under the MIT license. See [LICENSE](LICENSE) for more
information.

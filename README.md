# tabulate
Tabulate is a command-line tool to convert record type data (e.g. CSV files)
into a pretty printed table.

There have been multiple versions of this code over time. The original version
was `tabulate.sh`, which was written because I needed such a tool. Next came
`tabulate.py` as a Python starter project. The most recent version is
`tabulate.go`, which is a Go starter project.

Although all three versions are provided, the only version currently being
improved is the Go version. If you find this project interesting, and would
like to contribute to one or another version, feel free to do so.

## Project links
* Build Status:  [![Build Status][CIStatus]][CIProject]
* Documentation: [![GoDoc][GoDocStatus]][GoDoc]
* Views:         [![views][SGViews]][SGProject] [![views_24h][SGViews24h]][SGProject]
* Users:         [![library users][SGUsers]][SGProject] [![dependents][SGDependents]][SGProject]

## Usage
An easy example is to tabulate your `/etc/passwd` file.

```sh
$ go run tabulate.go -I : /etc/passwd
```

You can of course build tabulate into a binary, and place it into your favorite
binary location.

```sh
$ go build tabulate.go
$ mv tabulate ${HOME}/bin
```

To get a full list of options, request `--help`.

```
$ go run tabulate.go --help
Usage of /var/folders/00/0525h000h01000cxqpysvccm000m8p/T/go-build655823536/command-line-arguments/_obj/exe/tabulate:
  -I=" ": Input field separator.
  -O=" ": Output field separator.
  -cols=0: Number of columns; 0=all.
  -comment_prefix="#": Comment prefix.
  -comments=true: Ignore comments.
  -r="plain": Output renderer. (shorthand)
  -render="plain": Output renderer.
Supported renderers:
  csv
  markdown
  mysql
  plain
  sqlite3
```

<!--- Links -->
[CIProject]: https://travis-ci.org/kward/tabulate
[CIStatus]: https://travis-ci.org/kward/tabulate.png?branch=master

[GoDoc]: https://godoc.org/github.com/kward/tabulate
[GoDocStatus]: https://godoc.org/github.com/kward/tabulate?status.svg

[SGProject]: https://sourcegraph.com/github.com/kward/tabulate
[SGDependents]: https://sourcegraph.com/api/repos/github.com/kward/tabulate/.badges/dependents.svg
[SGUsers]: https://sourcegraph.com/api/repos/github.com/kward/tabulate/.badges/library-users.svg
[SGViews]: https://sourcegraph.com/api/repos/github.com/kward/tabulate/.counters/views.svg
[SGViews24h]: https://sourcegraph.com/api/repos/github.com/kward/tabulate/.counters/views-24h.svg?no-count=1


# tabulate
Tabulate is a command-line tool to convert record type data (e.g. CSV files)
into a pretty printed table.

## Project links
[![Travis CI](https://img.shields.io/travis/kward/tabulate.svg)](https://travis-ci.org/kward/tabulate)
[![GoDoc](https://godoc.org/github.com/kward/tabulate?status.svg)](https://godoc.org/github.com/kward/tabulate)

## Usage
An easy example is to tabulate your `/etc/passwd` file. Here are the first 10
non-comment lines from my `/etc/passwd` file.

```console
$ grep -v ^# /etc/passwd |head
nobody:*:-2:-2:Unprivileged User:/var/empty:/usr/bin/false
root:*:0:0:System Administrator:/var/root:/bin/sh
daemon:*:1:1:System Services:/var/root:/usr/bin/false
_uucp:*:4:4:Unix to Unix Copy Protocol:/var/spool/uucp:/usr/sbin/uucico
_taskgated:*:13:13:Task Gate Daemon:/var/empty:/usr/bin/false
_networkd:*:24:24:Network Services:/var/networkd:/usr/bin/false
_installassistant:*:25:25:Install Assistant:/var/empty:/usr/bin/false
_lp:*:26:26:Printing Services:/var/spool/cups:/usr/bin/false
_postfix:*:27:27:Postfix Mail Server:/var/spool/postfix:/usr/bin/false
_scsd:*:31:31:Service Configuration Service:/var/empty:/usr/bin/false
```

And again run through tabulate.

```console
$ grep -v ^# /etc/passwd |head |go run tabulate.go -I :
nobody            * -2 -2 Unprivileged User             /var/empty         /usr/bin/false
root              * 0  0  System Administrator          /var/root          /bin/sh
daemon            * 1  1  System Services               /var/root          /usr/bin/false
_uucp             * 4  4  Unix to Unix Copy Protocol    /var/spool/uucp    /usr/sbin/uucico
_taskgated        * 13 13 Task Gate Daemon              /var/empty         /usr/bin/false
_networkd         * 24 24 Network Services              /var/networkd      /usr/bin/false
_installassistant * 25 25 Install Assistant             /var/empty         /usr/bin/false
_lp               * 26 26 Printing Services             /var/spool/cups    /usr/bin/false
_postfix          * 27 27 Postfix Mail Server           /var/spool/postfix /usr/bin/false
_scsd             * 31 31 Service Configuration Service /var/empty         /usr/bin/false
```

And again rendered using the `mysql` renderer.

```
$ grep -v ^# /etc/passwd |head |go run tabulate.go -I : -r mysql
+-------------------+---+----+----+-------------------------------+--------------------+------------------+
| nobody            | * | -2 | -2 | Unprivileged User             | /var/empty         | /usr/bin/false   |
| root              | * | 0  | 0  | System Administrator          | /var/root          | /bin/sh          |
| daemon            | * | 1  | 1  | System Services               | /var/root          | /usr/bin/false   |
| _uucp             | * | 4  | 4  | Unix to Unix Copy Protocol    | /var/spool/uucp    | /usr/sbin/uucico |
| _taskgated        | * | 13 | 13 | Task Gate Daemon              | /var/empty         | /usr/bin/false   |
| _networkd         | * | 24 | 24 | Network Services              | /var/networkd      | /usr/bin/false   |
| _installassistant | * | 25 | 25 | Install Assistant             | /var/empty         | /usr/bin/false   |
| _lp               | * | 26 | 26 | Printing Services             | /var/spool/cups    | /usr/bin/false   |
| _postfix          | * | 27 | 27 | Postfix Mail Server           | /var/spool/postfix | /usr/bin/false   |
| _scsd             | * | 31 | 31 | Service Configuration Service | /var/empty         | /usr/bin/false   |
+-------------------+---+----+----+-------------------------------+--------------------+------------------+
```

You can of course build tabulate into a binary, and place it into your favorite
binary location.

```console
$ go build tabulate.go
$ mv tabulate ${HOME}/bin
```

To get a full list of options, request `--help`.

```console
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

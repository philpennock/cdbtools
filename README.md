cdbtools
========

Trivial wrapper programs around <https://github.com/colinmarc/cdb> to
implement `cdbdump`, `cdbmake` and `cdbget` as documented at
<https://cr.yp.to/cdb.html>.

Diagnostic messages do not attempt to match, but exit codes from `cdbget`
should conform.

There's no attempt to implement `cdbtest` or `cdbstats`.

## Installation

`go get github.com/philpennock/cdbtools/...`

## Usage

Use as command-line tools as drop-in replacements for the original native C
programs.

If any variances in output cause problems, please open an Issue on GitHub.
Pull Requests likely accepted.

For anything more than "getting a shell tool using the same lib as our Golang
code is using", please don't use these tools.  The underlying Go library
should be directly used.

## Issues

No attempts made to handle byte sequences which don't make cleanly to a string
in Golang.

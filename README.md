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

## Ambiguity

For `cdbget`, [the documentation](https://cr.yp.to/cdb/cdbget.html) of the
skip-count optional second parameter is:
> Given a numeric `s` argument, cdbget skips past the first `s` records with
> key `k`, and prints the data in the next record.

Taken literally, that's "look linearly for the first `s` instances of the key
`k` and then print the entire key/value record of whatever follows after it in
the CDB file, no matter what the key is".

I've chosen to interpret it a little less literally but more reasonably as
"look for a record with key `k`, skipping past the first `s` instances, and on
the next such record found, print the stored value".

## Issues

No attempts made to handle byte sequences which don't make cleanly to a string
in Golang.

// cdbget is a Golang implementation of the CDB cdbget command
package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/colinmarc/cdb"
)

// Exit values as per documentation at https://cr.yp.to/cdb/cdbget.html

func die(exitValue int, spec string, args ...interface{}) {
	if exitValue != 1 {
		fmt.Fprintf(os.Stderr, "%s: ", filepath.Base(os.Args[0]))
	}
	fmt.Fprintf(os.Stderr, spec, args)
	fmt.Fprintln(os.Stderr)
	os.Exit(exitValue)
}

func usage() {
	die(1, "Usage: %s <key> [<skip-count>]", filepath.Base(os.Args[0]))
}

func main() {
	if len(os.Args) != 2 && len(os.Args) != 3 {
		usage()
	}

	needleKey := []byte(os.Args[1])
	skipCount := 0
	if len(os.Args) >= 3 {
		sc, err := strconv.ParseUint(os.Args[2], 10, 32)
		if err != nil {
			usage()
		}
		skipCount = int(sc)
	}

	c, err := cdb.New(os.Stdin, nil)
	if err != nil {
		die(111, "unable to convert stdin to CDB: %s", err)
	}

	for iter := c.Iter(); iter.Next(); {
		if err := iter.Err(); err != nil {
			die(111, "iterator error: %s", err)
		}
		k := iter.Key()
		if !bytes.Equal(k, needleKey) {
			continue
		}
		if skipCount > 0 {
			skipCount -= 1
			continue
		}
		v := iter.Value()
		fmt.Printf("%s\n", v)
		os.Exit(0)
	}
	os.Exit(100)
}

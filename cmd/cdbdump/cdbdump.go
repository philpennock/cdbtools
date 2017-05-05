// cdbdump is a Golang implementation of the CDB cdbdump command
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/colinmarc/cdb"
)

func main() {
	c, err := cdb.New(os.Stdin, nil)
	if err != nil {
		log.Fatalf("unable to convert stdin to CDB: %s", err)
	}

	for iter := c.Iter(); iter.Next(); {
		if err := iter.Err(); err != nil {
			log.Fatalf("iterator error: %s", err)
		}
		k := iter.Key()
		v := iter.Value()
		fmt.Printf("+%d,%d:%s->%s\n", len(k), len(v), k, v)
	}
	fmt.Println()
}

// cdbmake is a Golang implementation of the CDB cdbmake command
package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/colinmarc/cdb"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <target> <tmp>", filepath.Base(os.Args[0]))
	}

	wr, err := cdb.Create(os.Args[2])
	if err != nil {
		log.Fatalf("unable to create CDB file at %q: %s", os.Args[2], err)
	}

	defer func() {
		if wr != nil {
			_ = wr.Close()
			_ = os.Remove(os.Args[2])
		}
	}()

	rd := bufio.NewReader(os.Stdin)
	lineCount := 0

	for {
		b, err := rd.ReadByte()
		if err != nil {
			log.Fatalf("Unable to read next line of CDB input after %d lines: %s", lineCount, err)
		}
		lineCount += 1
		if b == '\n' {
			break
		}
		if b != '+' {
			log.Fatalf("malformed line %d, does not start '+'", lineCount)
		}

		keyLenBytes, err := rd.ReadSlice(',')
		if err != nil {
			log.Fatalf("reading key-len line %d: %s", lineCount, err)
		}
		// 32-bit size is fine; a CDB file does not support a DB larger than 4GB.
		// Yes, we'll try to read an entire record into memory, let the OS kill
		// us if that's bad.
		keyLen, err := strconv.ParseUint(string(keyLenBytes[:len(keyLenBytes)-1]), 10, 32)
		if err != nil {
			log.Fatalf("converting key-len to int on line %d: %s", lineCount, err)
		}

		valueLenBytes, err := rd.ReadSlice(':')
		if err != nil {
			log.Fatalf("reading value len line %d: %s", lineCount, err)
		}
		valueLen, err := strconv.ParseUint(string(valueLenBytes[:len(valueLenBytes)-1]), 10, 32)
		if err != nil {
			log.Fatalf("converting value-len to int on line %d: %s", lineCount, err)
		}

		keyBuf := make([]byte, keyLen)
		valueBuf := make([]byte, valueLen)

		for offset := uint32(0); keyLen > 0; {
			n, err := rd.Read(keyBuf[offset:])
			if err != nil {
				log.Fatalf("reading key on line %d: %s", lineCount, err)
			}
			offset += uint32(n)
			keyLen -= uint64(n)
		}

		if b, err = rd.ReadByte(); err != nil || b != '-' {
			if err != nil {
				log.Fatalf("reading key/value separator line %d: %s", lineCount, err)
			}
			log.Fatalf("reading key/value separator line %d: does not start '-'")
		}
		if b, err = rd.ReadByte(); err != nil || b != '>' {
			if err != nil {
				log.Fatalf("reading key/value separator line %d: %s", lineCount, err)
			}
			log.Fatalf("reading key/value separator line %d: does not end '>'")
		}

		for offset := uint32(0); valueLen > 0; {
			n, err := rd.Read(valueBuf[offset:])
			if err != nil {
				log.Fatalf("reading value on line %d: %s", lineCount, err)
			}
			offset += uint32(n)
			valueLen -= uint64(n)
		}

		b, err = rd.ReadByte()
		if err != nil {
			log.Fatalf("Unable to read final newline on line %d: %s", lineCount, err)
		}
		if b != '\n' {
			log.Fatalf("malformed line %d, does not end '\\n'", lineCount)
		}

		err = wr.Put(keyBuf, valueBuf)
		if err != nil {
			log.Fatalf("unable to store key/value from line %d: %s", lineCount, err)
		}
	}

	if err := wr.Close(); err != nil {
		log.Fatalf("CDB(%q).Close failed: %s", os.Args[2], err)
	}

	wr = nil

	if err := os.Rename(os.Args[2], os.Args[1]); err != nil {
		log.Fatalf("renaming %q → %q failed: %s", os.Args[2], os.Args[1], err)
	}
}

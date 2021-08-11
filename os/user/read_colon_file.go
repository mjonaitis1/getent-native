// +build !windows

package user

import (
	"bufio"
	"bytes"
	"io"
)

// lineFunc returns a value, an error, or (nil, nil) to skip the row.
type lineFunc func(line []byte) (v interface{}, err error)

// readColonFile parses r as an /etc/group or /etc/passwd style file, running
// fn for each row. readColonFile returns a value, an error, or (nil, nil) if
// the end of the file is reached without a match.
//
// readColonFile can also be used to read /adm/users on plan9 as the
func readColonFile(r io.Reader, fn lineFunc) (v interface{}, err error) {
	bs := bufio.NewScanner(r)
	for bs.Scan() {
		line := bs.Bytes()
		// There's no spec for /etc/passwd or /etc/group, but we try to follow
		// the same rules as the glibc parser, which allows comments and blank
		// space at the beginning of a line.
		line = bytes.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		v, err = fn(line)
		if v != nil || err != nil {
			return
		}
	}
	return nil, bs.Err()
}

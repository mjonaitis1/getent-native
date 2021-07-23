// +build aix darwin dragonfly freebsd js,wasm !android,linux netbsd openbsd solaris
// +build !cgo osusergo

package user

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func allUsers() ([]*User, error){
	f, err := os.Open(userFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ch := readEtcFile(f, parseUser)
	users := make([]*User, 0)
	for efr := range ch{
		v, err := efr.v, efr.err
		if err == nil{
			if u, ok := v.(*User); ok{
				users = append(users, u)
			}
		}
	}
	return users, nil
}

func allGroups() ([]*Group, error){
	f, err := os.Open(groupFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ch := readEtcFile(f, parseGroup)
	groups := make([]*Group, 0)
	for efr := range ch{
		v, err := efr.v, efr.err
		if err == nil{
			if g, ok := v.(*Group); ok{
				groups = append(groups, g)
			}
		}
	}
	return groups, nil
}

// etcFileResult holds result from single line read
type etcFileResult struct {
	v interface{}
	err error
}

// readColonFile parses r as an /etc/group or /etc/passwd style file, running
// fn for each row. readColonFile returns a value, an error, or (nil, nil) if
// the end of the file is reached without a match.
func readEtcFile(r io.Reader, fn lineFunc) chan etcFileResult {
	ch := make(chan etcFileResult)
	isChClosed := false
	sig := make(chan os.Signal)
	closeCh := func(){
		if !isChClosed {
			close(ch)
			isChClosed = true
		}
	}
	// TODO ask motiejus if we really need this (Actually we - don't - but still not sure why this works)
	go func() {
		signal.Notify(sig, syscall.SIGPIPE)
		for{
			select {
				case <-sig:
					closeCh()
					return
			}
		}
	}()
	// TODO ask motiejus if it's better not to use channels here
	go func(){
		defer closeCh()
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

			v, err := fn(line)
			if v != nil || err != nil {
				ch <- etcFileResult{
					v: v,
					err: err,
				}
			}
		}
	}()
	return ch
}

// returns a *User for a row if that row's has the given value at the
// given index.
func parseUser (line []byte) (v interface{}, err error) {
	if bytes.Count(line, colon) < 6 {
		return
	}
	// kevin:x:1005:1006::/home/kevin:/usr/bin/zsh
	parts := strings.SplitN(string(line), ":", 7)
	if _, err := strconv.Atoi(parts[2]); err != nil {
		return nil, nil
	}
	if _, err := strconv.Atoi(parts[3]); err != nil {
		return nil, nil
	}
	u := &User{
		Username: parts[0],
		Uid:      parts[2],
		Gid:      parts[3],
		Name:     parts[4],
		HomeDir:  parts[5],
	}
	// The pw_gecos field isn't quite standardized. Some docs
	// say: "It is expected to be a comma separated list of
	// personal data where the first item is the full name of the
	// user."
	if i := strings.Index(u.Name, ","); i >= 0 {
		u.Name = u.Name[:i]
	}
	return u, nil
}


func parseGroup (line []byte) (v interface{}, err error) {
	if bytes.Count(line, colon) < 3 {
		return
	}
	// wheel:*:0:root
	parts := strings.SplitN(string(line), ":", 4)
	if len(parts) < 4 || parts[0] == "" ||
		// If the file contains +foo and you search for "foo", glibc
		// returns an "invalid argument" error. Similarly, if you search
		// for a gid for a row where the group name starts with "+" or "-",
		// glibc fails to find the record.
		parts[0][0] == '+' || parts[0][0] == '-' {
		return
	}
	if _, err := strconv.Atoi(parts[2]); err != nil {
		return nil, nil
	}
	return &Group{Name: parts[0], Gid: parts[2]}, nil
}

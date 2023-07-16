package pocketbase

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lainio/err2/try"
)

func Start() (*exec.Cmd, string) {
	pwd := getGitRoot()
	addr := getAddr()

	{
		f := try.To1(os.OpenFile(filepath.Join(pwd, "pb_try"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm))
		defer f.Close()
		fmt.Fprintln(f, "w")
	}

	build := exec.Command("go", "build", "-o", "pb_server", "./internal/pocketbase/pb")
	build.Dir = pwd
	try.To(build.Run())

	r, w := io.Pipe()
	rr := io.TeeReader(r, os.Stdout)
	// rr := r

	tmpdir := try.To1(os.MkdirTemp(os.TempDir(), "pocketbase-go-sdk"))
	cmd := exec.Command("./pb_server", "serve", "--dir", tmpdir, "--http", addr)
	cmd.Dir = pwd
	cmd.Stdout = w
	cmd.Stderr = os.Stderr

	try.To(cmd.Start())

	var buf = make([]byte, 1)
	try.To1(rr.Read(buf))
	go func() {
		var buf = make([]byte, 512)
		for {
			if _, err := rr.Read(buf); err != nil {
				break
			}
		}
	}()

	return cmd, addr
}

func getGitRoot() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	var b bytes.Buffer
	cmd.Stdout = &b
	try.To(cmd.Run())
	return strings.TrimSuffix(b.String(), "\n")
}

func getAddr() string {
	l := try.To1(net.Listen("tcp", "127.0.0.1:0"))
	defer l.Close()
	return l.Addr().String()
}

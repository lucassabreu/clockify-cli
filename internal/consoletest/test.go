package consoletest

import (
	"bytes"
	"io"
	"testing"
	"time"

	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"

	pseudotty "github.com/creack/pty"
)

// ExpectConsole is a helper to interact if the pseudo terminal on tests
type ExpectConsole interface {
	ExpectEOF()
	ExpectString(string)
	Send(string)
	SendLine(string)
}

type console struct {
	t *testing.T
	c *expect.Console
}

func (c *console) ExpectEOF() {
	if _, err := c.c.ExpectEOF(); err != nil {
		c.t.Errorf("failed to ExpectEOF %v", err)
	}
}

func (c *console) ExpectString(s string) {
	if _, err := c.c.ExpectString(s); err != nil {
		c.t.Errorf("failed to ExpectString %v", err)
	}
}

func (c *console) Send(s string) {
	if _, err := c.c.Send(s); err != nil {
		c.t.Errorf("failed to Send %v", err)
	}
}

func (c *console) SendLine(s string) {
	if _, err := c.c.SendLine(s); err != nil {
		c.t.Errorf("failed to SendLine %v", err)
	}
}

// FileWriter is a simplification of the io.Stdout struct
type FileWriter interface {
	io.Writer
	Fd() uintptr
}

// FileReader is a simplification of the io.Stdin struct
type FileReader interface {
	io.Reader
	Fd() uintptr
}

// RunTestConsole simulates a terminal for interactive tests
// This is mostly a adaptation of the RunTest function at
// [survey_test.go](https://github.com/AlecAivazis/survey/blob/e47352f914346a910cc7e1ca9f65a7ac0674449a/survey_posix_test.go#L15),
// but with interfaces exported to easy re-use on other packages.
func RunTestConsole(
	t *testing.T,
	setup func(out FileWriter, in FileReader) error,
	procedure func(c ExpectConsole),
) {
	pty, tty, err := pseudotty.Open()
	if err != nil {
		t.Fatalf("failed to open pseudotty: %v", err)
	}

	b := bytes.NewBufferString("")
	term := vt10x.New(vt10x.WithWriter(tty))
	c, err := expect.NewConsole(
		expect.WithStdin(pty),
		expect.WithStdout(term),
		expect.WithStdout(b),
		expect.WithCloser(pty, tty),
	)

	if err != nil {
		t.Fatalf("failed to create console: %v", err)
	}
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)

		go procedure(&console{c: c, t: t})
		if err = setup(c.Tty(), c.Tty()); err != nil {
			t.Error(err)
		}

		if err := c.Tty().Close(); err != nil {
			t.Errorf("error closing Tty: %v", err)
		}
	}()

	select {
	case <-time.After(time.Second * 10):
		t.Error(
			"console test timeout exceeded\n" +
				"current output:\n" +
				b.String())
	case <-donec:
	}
}

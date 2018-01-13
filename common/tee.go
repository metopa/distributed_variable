package common

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type Tee struct {
	wg         *sync.WaitGroup
	OrigStdout *os.File
	OrigStderr *os.File
	Stdout     *os.File
	Stderr     *os.File
	TeeFile    *os.File
}

func NewTee(filename string) (*Tee, error) {
	// copy stdio so we can print to dup to the
	// original files
	origStdout := os.Stdout
	origStderr := os.Stderr

	// create stdio pipe so we can read from stdout
	// and write it to our tee file
	stdoutRead, stdoutWrite, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	stderrRead, stderrWrite, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	// create the tee file to clone stdio to
	teefile, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	os.Stdout = stdoutWrite
	os.Stderr = stderrWrite

	wg := &sync.WaitGroup{}

	// function to be called as a goroutine that will
	// read from the read-end of our pipe and copy
	// to the dupd stdio plus the teefile
	f := func(reader, writer *os.File) {
		buf := make([]byte, 4096)
		for true {
			read, err := reader.Read(buf)

			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Fprintf(os.Stderr, "got err from bufReader: %s", err)
				break
			}
			if read == 0 {
				break
			}

			_, err = writer.Write(buf[:read])
			if err != nil {
				teefile.Write([]byte(fmt.Sprintf("write err: %s\n", err)))
				break
			}

			_, err = teefile.Write(buf[:read])
			if err != nil {
				break
			}
		}
		wg.Done()
	}
	wg.Add(2)

	go f(stdoutRead, origStdout)
	go f(stderrRead, origStderr)
	return &Tee{wg, origStdout, origStderr, stdoutWrite, stderrWrite, teefile}, nil
}

func (t *Tee) Sync() {
	t.Stdout.Sync()
	t.Stderr.Sync()
	t.TeeFile.Sync()
	t.OrigStdout.Sync()
	t.OrigStderr.Sync()
}

func (t *Tee) Close() {
	// close write end of pipes so
	// readers will finish up
	t.Stdout.Close()
	t.Stderr.Close()

	// wait for go routines to close
	t.wg.Wait()

	// close out tee file
	t.TeeFile.Close()

	// restore stdio to the original
	os.Stdout = t.OrigStdout
	os.Stderr = t.OrigStderr
}

package screen

import (
	"io"
	"os"
)

var (
	StdinReader  *io.PipeReader
	StdinWriter  *io.PipeWriter
	StdoutReader *io.PipeReader
	StdoutWriter *io.PipeWriter
	StderrReader *io.PipeReader
	StderrWriter *io.PipeWriter
)

func init() {
	StdinReader, StdinWriter = io.Pipe()
	StdoutReader, StdoutWriter = io.Pipe()
	StderrReader, StderrWriter = io.Pipe()
	go io.Copy(StdinWriter, os.Stdin)
	go io.Copy(os.Stdout, StdoutReader)
	go io.Copy(os.Stderr, StderrReader)
}

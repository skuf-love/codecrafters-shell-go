package main

import(
	"os/exec"
	"io"
)

type ExecCmdWraper struct {
	cmd *exec.Cmd
}

func (w *ExecCmdWraper) Run() error {
	return	w.cmd.Run()
}

func (w *ExecCmdWraper) SetStdin(stdin io.Reader) {
	w.cmd.Stdin = stdin
}

func (w *ExecCmdWraper) SetStdout(stdout io.Writer) {
	w.cmd.Stdout = stdout
}

func (w *ExecCmdWraper) SetStderr(stderr io.Writer) {
	w.cmd.Stderr = stderr
}

func (w *ExecCmdWraper) StdoutPipe() (io.ReadCloser, error) {
	return w.cmd.StdoutPipe()
}

func (w *ExecCmdWraper) Start() error {
	return	w.cmd.Start()
}

func (w *ExecCmdWraper) Wait() error {
	return	w.cmd.Wait()
}

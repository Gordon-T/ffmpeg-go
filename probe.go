package ffmpeg_go

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

// Probe Run ffprobe on the specified file and return a JSON representation of the output.
func Probe(ffprobeDir string, fileName string, kwargs ...KwArgs) (string, error) {
	return ProbeWithTimeout(ffprobeDir, fileName, 0, MergeKwArgs(kwargs))
}

func ProbeWithTimeout(ffprobeDir string, fileName string, timeOut time.Duration, kwargs KwArgs) (string, error) {
	args := KwArgs{
		"show_format":  "",
		"show_streams": "",
		"of":           "json",
	}

	return ProbeWithTimeoutExec(ffprobeDir, fileName, timeOut, MergeKwArgs([]KwArgs{args, kwargs}))
}

func ProbeWithTimeoutExec(ffprobeDir string, fileName string, timeOut time.Duration, kwargs KwArgs) (string, error) {
	if ffprobeDir == "" {
		ffprobeDir = "ffprobe.exe"
	}
	args := ConvertKwargsToCmdLineArgs(kwargs)
	args = append(args, fileName)
	ctx := context.Background()
	if timeOut > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(context.Background(), timeOut)
		defer cancel()
	}
	cmd := exec.CommandContext(ctx, ffprobeDir, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	buf := bytes.NewBuffer(nil)
	stdErrBuf := bytes.NewBuffer(nil)
	cmd.Stdout = buf
	cmd.Stderr = stdErrBuf
	for _, option := range GlobalCommandOptions {
		option(cmd)
	}
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("[%s] %w", string(stdErrBuf.Bytes()), err)
	}
	return string(buf.Bytes()), nil
}

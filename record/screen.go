package record

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/hellojukay/gors/output"
	cterminal "github.com/hellojukay/gors/terminal"

	"github.com/kr/pty"
	"golang.org/x/crypto/ssh/terminal"
)

type Screener struct {
	recorder *Recorder
	pty      *os.File

	sync.Mutex

	width  int64
	height int64
}

func NewScreener(r *Recorder) *Screener {
	s := &Screener{
		recorder: r,
	}
	return s
}

func (s *Screener) setSize() error {
	s.Lock()
	defer s.Unlock()

	ff, _ := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	defer ff.Close()

	w, h, err := terminal.GetSize(int(ff.Fd()))
	if err != nil {
		return err
	}

	s.width = int64(w)
	s.height = int64(h)

	window := [4]uint16{uint16(h), uint16(w), uint16(w * 8), uint16(h * w)}
	if _, _, err := syscall.Syscall(
		syscall.SYS_IOCTL,
		s.pty.Fd(),
		uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&window)),
	); err != 0 {
		return err
	}

	return nil
}

func (s *Screener) screen(r *Recorder) error {

	if s.pty != nil {
		panic("Screener has already running")
	}

	c := exec.Command("sh", "-c", r.Command)
	pty, err := pty.Start(c)
	if err != nil {
		panic(err)
	}

	file, err := os.Create(r.Filename)
	if err != nil {
		panic("can not create filename: " + r.Filename + " ,err is: " + err.Error())
	}
	defer file.Close()

	// 根据回车识别出敲打的命令
	ct, err := cterminal.NewCmdTermial()
	if err != nil {
		panic(err)
	}

	fileSavedCmd, err := ioutil.TempFile("/tmp", "savecmd")
	if err != nil {
		panic(err)
	}
	defer fileSavedCmd.Close()

	// do you self callback
	ct.OnCmdCallback(func(cmd []byte) error {
		_, err := fileSavedCmd.WriteString(string(cmd) + "\n")
		return err
	})

	ctErrors := make(chan error)
	go ct.IOSelect(ctErrors)

	oldState, err := terminal.MakeRaw(0)
	if err == nil {
		defer terminal.Restore(0, oldState)
	}

	s.pty = pty

	if err := s.setSize(); err != nil {
		panic(err)
	}

	now := time.Now()
	// 存放终端录屏的输出
	bufferOutput := output.NewOutput(now, file, 2, s.width, s.height, r.Command, r.Title, os.Getenv("TERM"), os.Getenv("SHELL"))

	closed := make(chan struct{}, 4)
	exit := make(chan struct{}, 1)

	pipeInR, pipeInW := io.Pipe()
	pipeOutR, pipeOutW := io.Pipe()

	go func() {
		defer func() { closed <- struct{}{} }()
		io.Copy(pipeInW, os.Stdin)
	}()

	// 输入方向的IO重定向
	go func() {
		defer func() { closed <- struct{}{} }()

		buf := make([]byte, 1024)
		for {
			size, err := pipeInR.Read(buf)
			if err != nil {
				return
			}
			// 在这里自己实现写入文件
			s.pty.Write(buf[:size])
		}
	}()

	go func() {
		defer func() { closed <- struct{}{} }()
		io.Copy(pipeOutW, s.pty)
	}()

	// 输出方向的IO重定向
	go func() {
		defer func() { closed <- struct{}{} }()

		buf := make([]byte, 1024)
		for {
			size, err := pipeOutR.Read(buf)
			if err != nil {
				return
			}
			step := float32(time.Now().UnixNano()-now.UnixNano()) / 1e9
			jsonline := []interface{}{step, "o", string(buf[:size])}
			bytes, _ := json.Marshal(jsonline)
			bufferOutput.Write(bytes)
			bufferOutput.Write([]byte("\n"))

			os.Stdout.Write(buf[:size])
		}
	}()

	go func() {
		<-closed
		c.Process.Signal(syscall.Signal(syscall.SIGHUP))
		s.pty.Close()

		exit <- struct{}{}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGWINCH)
	go func() {
		for sig := range sigs {
			if sig == os.Interrupt || sig == os.Kill || sig == syscall.SIGTERM {
				closed <- struct{}{}
			}
			if sig == syscall.SIGWINCH {
				if err := s.setSize(); err != nil {
					panic(err)
				}
			}
		}
	}()

	<-exit
	return nil
}

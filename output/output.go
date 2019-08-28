package output

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

type Output struct {
	buffer  []byte
	MaxWait int64
	sync.Mutex
	TimeStamp int64
	Writer    io.Writer
	Header    Header
}

type Header struct {
	Version   string `json:"version"`
	Width     int64  `json:"width"`
	Height    int64  `json:"height"`
	Command   string `json:"command"`
	Title     string `json:"title"`
	TimeStamp int64  `json:"timestamp"`
	Env       Env    `json:"env"`
}

type Env struct {
	TERM  string `json:"term"`
	SHELL string `json:"shell"`
}

func (d *Header) Save(file io.ReadWriter) error {
	bytes, err := json.Marshal(d)
	if err != nil {
		return err
	}
	_, err = file.Write(bytes)
	return err
}

func NewHeader(output *Output, version string, width, height int64, command, title, term, shell string) *Header {
	return &Header{
		Version:   version,
		Width:     width,
		Height:    height,
		TimeStamp: time.Now().Unix(),
		Command:   command,
		Title:     title,
		Env: Env{
			SHELL: shell,
			TERM:  term,
		},
	}
}

func (o *Output) Write(data []byte) (int, error) {
	// TODO
	// append to file
	return len(data), nil
}

func (o *Output) writeHeader() {

}
func NewOutput(w io.Writer) *Output {
	o := &Output{
		Writer: w,
	}
	o.writeHeader()
	return o
}

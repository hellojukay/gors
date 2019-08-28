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
	Version   int   `json:"version"`
	Width     int64 `json:"width"`
	Height    int64 `json:"height"`
	TimeStamp int64 `json:"timestamp"`
	Env       Env   `json:"env"`
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

func (o *Output) Write(data []byte) (int, error) {
	return o.Writer.Write(data)
}

func (o *Output) writeHeader() {
	headerJson, _ := json.Marshal(o.Header)
	o.Write(headerJson)
}
func NewOutput(w io.Writer, version int, width, height int64, command, title, term, shell string) *Output {
	o := &Output{
		Writer: w,
		Header: Header{
			Version:   version,
			Width:     width,
			Height:    height,
			TimeStamp: time.Now().Unix(),
			Env: Env{
				SHELL: shell,
				TERM:  term,
			},
		},
	}
	o.writeHeader()
	return o
}

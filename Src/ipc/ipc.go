package ipc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/Microsoft/go-winio"
)

var (
	Ipc IPC = IPC{}
)

type JsonEventData struct {
	Text string `json:"text,omitempty"`
}

type JsonEvent struct {
	Op   int           `json:"op"`
	Data JsonEventData `json:"data"`
}

type OpCodes struct {
	EXIT          int
	CLS           int
	NOTIFICATION  int
	INVALID_TOKEN int
	ERROR         int
	CONFIG_EDIT   int
	CONSOLE       int
	RESTART       int
	NEW_TOKEN     int
}

type IPC struct {
	readPipe  net.Conn
	writeFile *os.File
	OnRead    func(string)
	OpCodes   OpCodes
}

func (i *IPC) Init() {
	readPipe, err := winio.DialPipe(`\\.\pipe\LightningBotU2S`, nil)
	if err != nil {
		log.Fatalf("Error opening pipe: %v", err)
	}
	i.readPipe = readPipe

	writeFile, err := os.OpenFile(`\\.\pipe\LightningBotS2U`, os.O_WRONLY, 0)
	if err != nil {
		log.Fatalf("Error opening pipe: %v", err)
	}
	i.writeFile = writeFile

	i.OpCodes.EXIT = 100
	i.OpCodes.CLS = 101
	i.OpCodes.NOTIFICATION = 102
	i.OpCodes.INVALID_TOKEN = 103
	i.OpCodes.ERROR = 104
	i.OpCodes.CONFIG_EDIT = 105
	i.OpCodes.CONSOLE = 106
	i.OpCodes.RESTART = 107
	i.OpCodes.NEW_TOKEN = 108

	go func() {
		for {
			lengthbuf := make([]byte, 4)
			_, err := i.readPipe.Read(lengthbuf)
			if err != nil {
				if err != io.EOF {
					log.Printf("read error: %v\n", err)
				}
				break
			}
			buf := make([]byte, binary.LittleEndian.Uint32(lengthbuf))
			n, err := i.readPipe.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Printf("read error: %v\n", err)
				}
				break
			}
			str := string(buf[:n])
			i.OnRead(str)
		}
	}()
}

func (i *IPC) Write(s string) {
	sb := []byte(s)
	lb := make([]byte, 4)
	binary.LittleEndian.PutUint32(lb, uint32(len(s)))
	fb := bytes.Join([][]byte{lb, sb}, []byte(""))
	_, err := i.writeFile.Write(fb)
	if err != nil {
		fmt.Println(err.Error())
	}
}

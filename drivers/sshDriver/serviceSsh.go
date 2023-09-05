package sshDriver

import (
	"bytes"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/Archie1978/regate/drivers"
	"github.com/mitchellh/mapstructure"

	"github.com/golang/glog"

	"golang.org/x/crypto/ssh"

	_ "embed"
)

type MsgConnect struct {
	Ip                 string
	Port               int
	Username, Password string
}

type MessageTerminal struct {
	Session int
	Command string

	Msg interface{}
}

type ProcessSsh struct {
	Stdin io.WriteCloser

	client       *ssh.Client
	sessionShell *ssh.Session

	chanelWebSocket chan interface{}
	numeroSession   int
	msgConnect      MsgConnect

	Screen MsgScreen
}

func (processSSh *ProcessSsh) New() drivers.DriverRP {
	return &ProcessSsh{}
}

func (processSSh *ProcessSsh) Start(chanelWebSocket chan interface{}, numeroSession int, urlString string) {
	processSSh.chanelWebSocket = chanelWebSocket
	processSSh.numeroSession = numeroSession

	u, err := url.Parse(urlString)
	if err == nil {
		// Decode URL
		processSSh.msgConnect.Ip = u.Hostname()
		processSSh.msgConnect.Port, _ = strconv.Atoi(u.Port())
		if processSSh.msgConnect.Port == 0 {
			processSSh.msgConnect.Port = 22
		}
		if u.User != nil {
			processSSh.msgConnect.Username = u.User.Username()
			processSSh.msgConnect.Password, _ = u.User.Password()
		}
	}

}

func (processSSh *ProcessSsh) startSSh() {

	config := &ssh.ClientConfig{
		User: processSSh.msgConnect.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(processSSh.msgConnect.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if processSSh.msgConnect.Port == 0 {
		processSSh.msgConnect.Port = 22
	}

	// connect ot ssh server
	var err error
	host := fmt.Sprintf("%v:%v", processSSh.msgConnect.Ip, processSSh.msgConnect.Port)
	processSSh.client, err = ssh.Dial("tcp", host, config)
	if err != nil {
		processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "Error", Msg: err}
		glog.Error("processDial to ", host, ":", err)

		processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "End"}
		return
	}

	// Create session for shell
	processSSh.sessionShell, err = processSSh.client.NewSession()
	if err != nil {
		processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "Error", Msg: err}
		glog.Error("processSession:", err)

		processSSh.client.Close()
		processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "End"}
		return
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 1500000,
		ssh.TTY_OP_OSPEED: 1500000,
	}

	// Create Pty
	err = processSSh.sessionShell.RequestPty("linux", processSSh.Screen.Rows, processSSh.Screen.Cols, modes)
	if err != nil {
		processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "Error", Msg: err}
		glog.Error("processSession:", err)
		processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "End"}
		return
	}

	processSSh.Stdin, _ = processSSh.sessionShell.StdinPipe()
	go func() {
		defer processSSh.sessionShell.Close()
		defer processSSh.client.Close()

		r, err := processSSh.sessionShell.StdoutPipe()
		if err != nil {
			processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "Error", Msg: err}
			glog.Error("processsessionShell.StdoutPipe:", err)
			processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "End"}
			return
		}
		for {
			b1 := make([]byte, 50)
			n1, err := r.Read(b1)
			if err != nil {
				processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "End", Msg: MsgOut{}}
				glog.Error("processRead:", err)
				return
			}
			processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "Out", Msg: MsgOut{Content: b1[:n1]}}
		}
	}()

	if err := processSSh.sessionShell.Shell(); err != nil {
		processSSh.chanelWebSocket <- MessageTerminal{Session: processSSh.numeroSession, Command: "Error", Msg: err}
		glog.Error("processSSh.sessionShell:", processSSh.sessionShell.Stderr.(*bytes.Buffer).String())
		return
	}
}

func (processSSh *ProcessSsh) Close() {
	processSSh.Stdin.Close()
	processSSh.sessionShell.Close()
	processSSh.client.Close()
}

type MsgOut struct {
	Content []byte
}

type MsgScreen struct {
	Cols int
	Rows int
}

type MsgKey struct {
	Keys []byte
}

type Message struct {
	Type   string
	Key    MsgKey
	Screen MsgScreen
}

func (processSSh *ProcessSsh) Process(msg interface{}) {

	msgMsg := msg.(map[string]interface{})
	var message Message
	mapstructure.Decode(msg, &message)

	switch message.Type {

	case "size":
		processSSh.Screen = message.Screen

		glog.Info("connect", processSSh.msgConnect)
		processSSh.startSSh()

	case "key":
		if processSSh.Stdin == nil {
			glog.Info("processSSh.Stdin not init:", msg)
			return
		}
		a := msgMsg["Key"].(map[string]interface{})
		b := a["Keys"].([]interface{})
		processSSh.Stdin.Write([]byte(b[0].(string)))

	default:
		glog.Info("Message unknown:", message.Type)
	}
}

var (
	//go:embed javascriptPlugin.js
	dataplugin []byte
)

func (processSSh *ProcessSsh) GetCodeJavascript() string {
	return string(dataplugin) + "\n"
}

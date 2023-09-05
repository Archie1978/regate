package rdpDriver

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Archie1978/regate/drivers"
	"github.com/mitchellh/mapstructure"

	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/protocol/pdu"
)

type Screen struct {
	Height int `json:"height"`
	Width  int `json:"width"`
}

type Info struct {
	Domain   string `json:"domain"`
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Passwd   string `json:"password"`
	Screen   `json:"screen"`
}

type Bitmap struct {
	DestLeft     int    `json:"destLeft"`
	DestTop      int    `json:"destTop"`
	DestRight    int    `json:"destRight"`
	DestBottom   int    `json:"destBottom"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	BitsPerPixel int    `json:"bitsPerPixel"`
	IsCompress   bool   `json:"isCompress"`
	Data         []byte `json:"data"`
}

func Bpp(BitsPerPixel uint16) (pixel int) {
	switch BitsPerPixel {
	case 15:
		pixel = 1

	case 16:
		pixel = 2

	case 24:
		pixel = 3

	case 32:
		pixel = 4

	default:
		glog.Error("invalid bitmap data format")
	}
	return
}

type MessageRTM struct {
	Session int
	Command string

	Msg interface{}
}

type MsgScanCode struct {
	X, Y, Button uint16
	IsPressed    bool
}
type MsgScreen struct {
	Width  int
	Height int
}

type MsgWheel struct {
	X, Y  uint16
	Step  uint16
	IsNeg bool
	IsH   bool
}

type MsgConnect struct {
	Ip                         string
	Port                       int
	Domain, Username, Password string

	Width  int
	Height int
}

func (processRDP *ProcessRdp) New() drivers.DriverRP {
	return &ProcessRdp{}
}
func (processRDP *ProcessRdp) startRDP(chanelWebSocket chan interface{}) {
	msgConnect := &processRDP.msgConnect
	if msgConnect.Width == 0 {
		msgConnect.Width = 800
	}
	if msgConnect.Height == 0 {
		msgConnect.Height = 600
	}
	if msgConnect.Username == "" {
		msgConnect.Domain = processRDP.msgConnect.Domain
		msgConnect.Ip = processRDP.msgConnect.Ip
		msgConnect.Username = "Administrator"
		msgConnect.Password = processRDP.msgConnect.Password
	}
	if msgConnect.Port == 0 {
		msgConnect.Port = 3389
	}
	glog.Info(" Display:", msgConnect.Width, "x", msgConnect.Height, " to ", msgConnect.Ip, ' ', msgConnect.Domain, "\\", msgConnect.Username)

	processRDP.g = NewRdpClient(fmt.Sprintf("%v:%v", msgConnect.Ip, msgConnect.Port), msgConnect.Width, msgConnect.Height, glog.INFO)
	var info Info
	info.Domain = msgConnect.Domain
	info.Ip = msgConnect.Ip
	info.Port = fmt.Sprintf("%v", msgConnect.Port)
	info.Username = msgConnect.Username
	info.Passwd = msgConnect.Password
	info.Screen.Height = msgConnect.Height
	info.Screen.Width = msgConnect.Width
	processRDP.g.info = &info

	err := processRDP.g.Login()
	if err != nil {
		jsonText, _ := json.Marshal(MessageRTM{Session: processRDP.numeroSession, Command: "Error", Msg: err.Error()})
		chanelWebSocket <- jsonText
		jsonText, _ = json.Marshal(MessageRTM{Session: processRDP.numeroSession, Command: "End"})
		chanelWebSocket <- jsonText
		glog.Error("connect Error :", err)
	} else {
		// Login passe
		processRDP.g.pdu.On("error", func(e error) {
			glog.Error("on error:", e)
			jsonText, _ := json.Marshal(MessageRTM{Session: processRDP.numeroSession, Command: "Error", Msg: e.Error()})
			chanelWebSocket <- jsonText
			jsonText, _ = json.Marshal(MessageRTM{Session: processRDP.numeroSession, Command: "End"})
			chanelWebSocket <- jsonText
			//so.Emit("rdp-error", "{\"code\":1,\"message\":\""+e.Error()+"\"}")
			//wg.Done()
		}).On("close", func() {
			glog.Info("RDP on close")
			jsonText, _ := json.Marshal(MessageRTM{Session: processRDP.numeroSession, Command: "Error", Msg: "Connection closed"})
			chanelWebSocket <- jsonText
			jsonText, _ = json.Marshal(MessageRTM{Session: processRDP.numeroSession, Command: "Close"})
			chanelWebSocket <- jsonText
			err = errors.New("close")

		}).On("success", func() {
			glog.Info("RDP on success")
			jsonText, _ := json.Marshal(MessageRTM{Session: processRDP.numeroSession, Command: "Success"})
			chanelWebSocket <- jsonText
		}).On("ready", func() {
			glog.Info("RDP on ready:", processRDP.g, "connect")
			jsonText, _ := json.Marshal(MessageRTM{Session: processRDP.numeroSession, Command: "Ready"})
			chanelWebSocket <- jsonText
		}).On("update", func(rectangles []pdu.BitmapData) {
			glog.Info(time.Now(), "on update Bitmap:", len(rectangles))
			go func() {
				bs := make([]Bitmap, 0, len(rectangles))
				for _, v := range rectangles {
					IsCompress := v.IsCompress()
					data := v.BitmapDataStream

					//glog.Debug("data:", data)
					if IsCompress {
						//data = BitmapDecompress(&v)
						//IsCompress = false
					}

					glog.Debug("Bitmap", "L", int(v.DestLeft), "T:", int(v.DestTop), "R:", int(v.DestRight), "B:", int(v.DestBottom),
						"Width:", int(v.Width), "H", int(v.Height), "BitParPixel", int(v.BitsPerPixel), IsCompress)
					b := Bitmap{int(v.DestLeft), int(v.DestTop), int(v.DestRight), int(v.DestBottom),
						int(v.Width), int(v.Height), int(v.BitsPerPixel), IsCompress, data}
					bs = append(bs, b)
				}

				// Update Cavas by websocket
				chanelWebSocket <- MessageRTM{Session: processRDP.numeroSession, Command: "Update", Msg: bs}
			}()
		})
	}
}

type ProcessRdp struct {
	g             *RdpClient
	numeroSession int

	chanelWebSocket chan interface{}
	msgConnect      MsgConnect
}

func (processRDP *ProcessRdp) Start(chanelWebSocket chan interface{}, numeroSession int, urlString string) {
	processRDP.numeroSession = numeroSession
	processRDP.chanelWebSocket = chanelWebSocket

	u, err := url.Parse(urlString)
	if err == nil {
		domainUser := strings.Split(u.User.Username(), "/")
		processRDP.msgConnect.Domain = domainUser[0]
		if len(domainUser) > 1 {
			processRDP.msgConnect.Username = domainUser[1]
		}
		processRDP.msgConnect.Password, _ = u.User.Password()
		processRDP.msgConnect.Ip = u.Host
	}

	glog.Info("start driver RDP:", processRDP.msgConnect)
}
func (processRDP *ProcessRdp) Close() {
	processRDP.g.Close()
	close(processRDP.chanelWebSocket)
}

type Message struct {
	Type     string
	ScanCode MsgScanCode
	Wheel    MsgWheel
	Screen   MsgScreen
}

func (processRDP *ProcessRdp) Process(msg interface{}) {

	var message Message
	mapstructure.Decode(msg, &message)
	switch message.Type {
	case "size":
		processRDP.msgConnect.Width = message.Screen.Width
		processRDP.msgConnect.Height = message.Screen.Height

		glog.Info("connect", processRDP.msgConnect)
		processRDP.startRDP(processRDP.chanelWebSocket)

	case "mouse":
		msgMouse := message.ScanCode
		glog.Info("mouse", msgMouse.X, ":", msgMouse.Y, ":", msgMouse.Button, ":", msgMouse.IsPressed)

		p := &pdu.PointerEvent{}
		if msgMouse.IsPressed {
			p.PointerFlags |= pdu.PTRFLAGS_DOWN
		}

		switch msgMouse.Button {
		case 1:
			p.PointerFlags |= pdu.PTRFLAGS_BUTTON1
		case 2:
			p.PointerFlags |= pdu.PTRFLAGS_BUTTON2
		case 3:
			p.PointerFlags |= pdu.PTRFLAGS_BUTTON3
		default:
			p.PointerFlags |= pdu.PTRFLAGS_MOVE
		}

		p.XPos = msgMouse.X
		p.YPos = msgMouse.Y
		if processRDP.g != nil {
			if processRDP.g.pdu != nil {
				processRDP.g.pdu.SendInputEvents(pdu.INPUT_EVENT_MOUSE, []pdu.InputEventsInterface{p})
			}
		}

	case "scancode":
		msgScanCode := message.ScanCode
		glog.Info("scancode", msgScanCode.Button)

		p := &pdu.ScancodeKeyEvent{}
		p.KeyCode = msgScanCode.Button
		if !msgScanCode.IsPressed {
			p.KeyboardFlags |= pdu.KBDFLAGS_RELEASE
		}
		processRDP.g.pdu.SendInputEvents(pdu.INPUT_EVENT_SCANCODE, []pdu.InputEventsInterface{p})

	case "wheel":
		msgWheel := message.Wheel
		mapstructure.Decode(msg, &msgWheel)

		glog.Info("wheel", msgWheel.X, ":", msgWheel.Y, ":", msgWheel.Step, ":", msgWheel.IsNeg, ":", msgWheel.IsH)
		var p = &pdu.PointerEvent{}
		if msgWheel.IsH {
			p.PointerFlags |= pdu.PTRFLAGS_HWHEEL
		} else {
			p.PointerFlags |= pdu.PTRFLAGS_WHEEL
		}

		if msgWheel.IsNeg {
			p.PointerFlags |= pdu.PTRFLAGS_WHEEL_NEGATIVE
		}

		p.PointerFlags |= (msgWheel.Step & pdu.WheelRotationMask)
		p.XPos = msgWheel.X
		p.YPos = msgWheel.Y

		processRDP.g.pdu.SendInputEvents(pdu.INPUT_EVENT_SCANCODE, []pdu.InputEventsInterface{p})

	default:
		glog.Info("Message unknown:", message.Type)
	}

}

var (
	//go:embed javascriptPlugin.js
	dataplugin []byte

	//go:embed canvas.js
	datapluginCanvas []byte

	//go:embed client.js
	datapluginClient []byte

	//go:embed mstsc.js
	datapluginMstsc []byte

	//go:embed rle.js
	datapluginRle []byte

	//go:embed keyboard.js
	datapluginKeyboard []byte
)

func (processRDP *ProcessRdp) GetCodeJavascript() string {
	return string(dataplugin) + "\n" + string(datapluginMstsc) + "\n" + string(datapluginCanvas) + "\n" + string(datapluginClient) + "\n" + string(datapluginRle) + "\n" + string(datapluginKeyboard) + "\n"
}

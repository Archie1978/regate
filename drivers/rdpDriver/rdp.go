package rdpDriver

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Archie1978/regate/crypto"
	"github.com/Archie1978/regate/database"
	"github.com/tomatome/grdp/plugin"

	"github.com/tomatome/grdp/core"
	"github.com/tomatome/grdp/glog"
	"github.com/tomatome/grdp/protocol/nla"
	"github.com/tomatome/grdp/protocol/pdu"
	"github.com/tomatome/grdp/protocol/sec"
	"github.com/tomatome/grdp/protocol/t125"
	"github.com/tomatome/grdp/protocol/tpkt"
	"github.com/tomatome/grdp/protocol/x224"
)

const (
	PROTOCOL_RDP       = x224.PROTOCOL_RDP
	PROTOCOL_SSL       = x224.PROTOCOL_SSL
	PROTOCOL_HYBRID    = x224.PROTOCOL_HYBRID
	PROTOCOL_HYBRID_EX = x224.PROTOCOL_HYBRID_EX
)

type RdpClient struct {
	Host string // ip:port

	Width    int
	Height   int
	info     *Info
	tpkt     *tpkt.TPKT
	x224     *x224.X224
	mcs      *t125.MCSClient
	sec      *sec.Client
	pdu      *pdu.Client
	channels *plugin.Channels

	// connection rdp secure or not ecure
	conn net.Conn

	// Connection TCP if is secure
	connRaw net.Conn
}

func NewRdpClient(host string, width, height int, logLevel glog.LEVEL) *RdpClient {
	return &RdpClient{
		Host:   host,
		Width:  width,
		Height: height,
	}
}
func (g *RdpClient) SetRequestedProtocol(p uint32) {
	g.x224.SetRequestedProtocol(p)
}

func BitmapDecompress(bitmap *pdu.BitmapData) []byte {
	return core.Decompress(bitmap.BitmapDataStream, int(bitmap.Width), int(bitmap.Height), Bpp(bitmap.BitsPerPixel))
}

func (g *RdpClient) Login() error {

	// Get security element
	security, err := database.GetSettingSecurity()

	domain, user, pwd := g.info.Domain, g.info.Username, g.info.Passwd

	glog.Info("Connect:", g.Host, "with", domain+"\\"+user)
	g.conn, err = net.DialTimeout("tcp", g.Host, 3*time.Second)
	if err != nil {
		return fmt.Errorf("[dial err] %v", err)
	}

	// Check certificate
	if security.Cert_activate {

		// Configuration TLS de base (peut nécessiter des ajustements en fonction du serveur)
		tlsConfig := &tls.Config{
			VerifyPeerCertificate: crypto.CheckCertificate(security.Cert_list, security.Cert_list),
			InsecureSkipVerify:    true,
		}

		tlsConn := tls.Client(g.conn, tlsConfig)

		// Handshake TLS
		err = tlsConn.Handshake()
		if err != nil {
			return fmt.Errorf("Erreur lors du handshake TLS: %v", err)

		}

		// Switch connect to secure for application
		g.connRaw = g.conn
		g.conn = tlsConn
	}

	g.tpkt = tpkt.New(core.NewSocketLayer(g.conn), nla.NewNTLMv2(domain, user, pwd))
	g.x224 = x224.New(g.tpkt)
	g.mcs = t125.NewMCSClient(g.x224)
	g.sec = sec.NewClient(g.mcs)
	g.pdu = pdu.NewClient(g.sec)
	g.channels = plugin.NewChannels(g.sec)

	log.Println("connect", uint16(g.Width), uint16(g.Height))
	g.mcs.SetClientCoreData(uint16(g.Width), uint16(g.Height))

	g.sec.SetUser(user)
	g.sec.SetPwd(pwd)
	g.sec.SetDomain(domain)

	g.tpkt.SetFastPathListener(g.sec)
	g.sec.SetFastPathListener(g.pdu)
	g.sec.SetChannelSender(g.mcs)
	g.channels.SetChannelSender(g.sec)
	//g.pdu.SetFastPathSender(g.tpkt)

	//g.x224.SetRequestedProtocol(x224.PROTOCOL_RDP)
	//g.x224.SetRequestedProtocol(x224.PROTOCOL_SSL)

	err = g.x224.Connect()
	if err != nil {
		return fmt.Errorf("[x224 connect err] %v", err)
	}
	return nil
}

func (g *RdpClient) KeyUp(sc int, name string) {
	glog.Debug("KeyUp:", sc, "name:", name)

	p := &pdu.ScancodeKeyEvent{}
	p.KeyCode = uint16(sc)
	p.KeyboardFlags |= pdu.KBDFLAGS_RELEASE
	g.pdu.SendInputEvents(pdu.INPUT_EVENT_SCANCODE, []pdu.InputEventsInterface{p})
}
func (g *RdpClient) KeyDown(sc int, name string) {
	glog.Debug("KeyDown:", sc, "name:", name)

	p := &pdu.ScancodeKeyEvent{}
	p.KeyCode = uint16(sc)
	g.pdu.SendInputEvents(pdu.INPUT_EVENT_SCANCODE, []pdu.InputEventsInterface{p})
}

func (g *RdpClient) MouseMove(x, y int) {
	glog.Debug("MouseMove", x, ":", y)
	p := &pdu.PointerEvent{}
	p.PointerFlags |= pdu.PTRFLAGS_MOVE
	p.XPos = uint16(x)
	p.YPos = uint16(y)
	g.pdu.SendInputEvents(pdu.INPUT_EVENT_MOUSE, []pdu.InputEventsInterface{p})
}

func (g *RdpClient) MouseWheel(scroll, x, y int) {
	glog.Info("MouseWheel", x, ":", y)
	p := &pdu.PointerEvent{}
	p.PointerFlags |= pdu.PTRFLAGS_WHEEL
	p.XPos = uint16(x)
	p.YPos = uint16(y)
	g.pdu.SendInputEvents(pdu.INPUT_EVENT_SCANCODE, []pdu.InputEventsInterface{p})
}

func (g *RdpClient) MouseUp(button int, x, y int) {
	glog.Debug("MouseUp", x, ":", y, ":", button)
	p := &pdu.PointerEvent{}

	switch button {
	case 0:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON1
	case 2:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON2
	case 1:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON3
	default:
		p.PointerFlags |= pdu.PTRFLAGS_MOVE
	}

	p.XPos = uint16(x)
	p.YPos = uint16(y)
	g.pdu.SendInputEvents(pdu.INPUT_EVENT_MOUSE, []pdu.InputEventsInterface{p})
}
func (g *RdpClient) MouseDown(button int, x, y int) {
	glog.Info("MouseDown:", x, ":", y, ":", button)
	p := &pdu.PointerEvent{}

	p.PointerFlags |= pdu.PTRFLAGS_DOWN

	switch button {
	case 0:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON1
	case 2:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON2
	case 1:
		p.PointerFlags |= pdu.PTRFLAGS_BUTTON3
	default:
		p.PointerFlags |= pdu.PTRFLAGS_MOVE
	}

	p.XPos = uint16(x)
	p.YPos = uint16(y)
	g.pdu.SendInputEvents(pdu.INPUT_EVENT_MOUSE, []pdu.InputEventsInterface{p})
}
func (g *RdpClient) Close() {
	if g != nil && g.tpkt != nil {
		g.tpkt.Close()
	}
	if g.conn != nil {
		g.conn.Close()
	}

	if g.connRaw != nil {
		g.connRaw.Close()
	}
}

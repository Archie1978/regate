package sshDriver

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Archie1978/regate/configuration"
	"github.com/Archie1978/regate/crypto"
	"github.com/gin-gonic/gin"

	"github.com/Archie1978/regate/drivers"
	"github.com/golang/glog"
	"github.com/mitchellh/mapstructure"

	"github.com/pkg/sftp"
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
	fmt.Println("ProcessSsh) startSSh")
	config := &ssh.ClientConfig{
		User: processSSh.msgConnect.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(crypto.DecryptPasswordString(processSSh.msgConnect.Password, configuration.ConfigurationGlobal.KeyCrypt)),
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

// Downfile
func (processSSh *ProcessSsh) DownloadFile(context *gin.Context) error {

	filePath := context.Query("path")
	filePath = strings.Trim(filePath, " \n\t")

	// Ouverture d'une session SFTP
	sftpClient, err := sftp.NewClient(processSSh.client)
	if err != nil {
		err = fmt.Errorf("Erreur en établissant la connexion SFTP: %s", err)
		fmt.Println(err)
		return err
	}
	defer sftpClient.Close()

	stat, err := sftpClient.Stat(filePath)
	if err != nil {
		err = fmt.Errorf("Stat la connexion SFTP: %s: %s", filePath, err)
		fmt.Println(err)
		return err
	}

	if stat.IsDir() {
		context.Header("Content-Type", "application/tar")
		fileName := stat.Name() + ".tar"
		context.Header("Content-Disposition", "attachment; filename="+fileName)

		createTar(sftpClient, filePath, context.Writer)
		return nil
	}

	// Lecture du contenu du fichier
	file, err := sftpClient.Open(filePath)
	if err != nil {
		err = fmt.Errorf("Erreur en ouvrant le fichier sur le serveur SFTP : %s", err)
		fmt.Println(err)
		return err
	}
	defer file.Close()

	context.Header("Content-Type", "text/plain")
	fileName := stat.Name()
	context.Header("Content-Disposition", "attachment; filename="+fileName)

	// Readfile
	s, err := io.Copy(context.Writer, file)
	if err != nil {
		err = fmt.Errorf("Stat la connexion SFTP: %s", err)
		fmt.Println(err)
		return err
	}
	if s != stat.Size() {
		err = fmt.Errorf("Stat la connexion SFTP: %s", err)
		fmt.Println(err)
		return err
	}
	return nil
}

func (processSSh *ProcessSsh) UploadFile(c *gin.Context) error {
	pathDir := c.Query("pathDir")

	// Get multi path
	multipartReader, err := c.Request.MultipartReader()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return err
	}

	// Itérer sur les parties du lecteur multipart
	for {

		// Get Part
		part, err := multipartReader.NextPart()
		if err == io.EOF {
			break
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return err
		}

		// Nom du fichier
		fileName := path.Join(pathDir, part.FileName())

		// Utiliser le nom du fichier original avec un prefixe aléatoire pour éviter les conflits
		sftpClient, err := sftp.NewClient(processSSh.client)
		if err != nil {
			err = fmt.Errorf("Erreur en établissant la connexion SFTP: %s", err)
			fmt.Println(err)
			return err
		}
		defer sftpClient.Close()

		dir := filepath.Dir(fileName)
		_, err = sftpClient.Stat(dir)
		if !os.IsExist(err) {
			if err = sftpClient.MkdirAll(dir); err != nil {
				return err
			}
		}

		_, err = sftpClient.Stat(fileName)
		if os.IsExist(err) {
			sftpClient.Remove(fileName + ".old")
			err = sftpClient.Rename(fileName, fileName+".old")
			if err != nil {
				return err
			}
		}

		out, err := sftpClient.Create(fileName)
		if err != nil {
			return err
		}
		defer out.Close()

		_, err = io.Copy(out, part)

		if err != nil {

		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
	return err
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

func createTar(sftpClient *sftp.Client, sourceDir string, writer io.Writer) error {

	// Renomme tarfile
	tarFile := writer

	// Initialiser le writer tar
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	// function addtoTar generic
	addToTar := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Créer le lien symbolique pour les répertoires
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			header, err := tar.FileInfoHeader(info, link)
			if err != nil {
				return err
			}
			header.Name = path
			if err := tarWriter.WriteHeader(header); err != nil {
				return err
			}
			return nil
		}

		// Créer l'en-tête tar pour le fichier ou le dossier
		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}

		// Utiliser le chemin relatif pour les fichiers dans l'archive
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		header.Name = relPath

		// Écrire l'en-tête dans l'archive
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// Si le fichier est un répertoire, ne rien écrire, car l'en-tête suffit
		if info.Mode().IsDir() {
			return nil
		}

		// Ouvrir le fichier source
		file, err := sftpClient.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		// Copier le contenu du fichier dans l'archive
		_, err = io.Copy(tarWriter, file)
		return err
	}

	// Parcourir le répertoire source et ajouter les fichiers/dossiers à l'archive
	walkerSsh := sftpClient.Walk(sourceDir)
	for walkerSsh.Step() {
		if walkerSsh.Err() != nil {
			continue
		}

		addToTar(walkerSsh.Path(), walkerSsh.Stat(), nil)
	}
	return nil
}

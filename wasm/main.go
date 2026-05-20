package main

import (
	"context"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"nhooyr.io/websocket"
)

var sshClient *ssh.Client
var sftpClient *sftp.Client

func connectSSH(this js.Value, args []js.Value) any {
	if len(args) < 3 {
		return "Missing arguments (host, user, pass)"
	}
	host := args[0].String()
	user := args[1].String()
	pass := args[2].String()
	key := ""
	if len(args) >= 4 {
		key = args[3].String()
	}

	proxyURL := fmt.Sprintf("ws://localhost:8080/ssh?host=%s", host)
	ctx := context.Background()
	c, _, err := websocket.Dial(ctx, proxyURL, nil)
	if err != nil {
		return fmt.Sprintf("WebSocket dial error: %v", err)
	}

	netConn := websocket.NetConn(context.Background(), c, websocket.MessageBinary)

	config := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if key != "" {
		signer, err := ssh.ParsePrivateKey([]byte(key))
		if err != nil {
			return fmt.Sprintf("Private key parse error: %v", err)
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		config.Auth = []ssh.AuthMethod{ssh.Password(pass)}
	}

	sshConn, chans, reqs, err := ssh.NewClientConn(netConn, host, config)
	if err != nil {
		return fmt.Sprintf("SSH handshake error: %v", err)
	}

	sshClient = ssh.NewClient(sshConn, chans, reqs)
	
	sc, err := sftp.NewClient(sshClient)
	if err == nil {
		sftpClient = sc
	}

	return "Connected"
}

func executeCommand(this js.Value, args []js.Value) any {
	if sshClient == nil {
		return "Not connected"
	}
	if len(args) < 1 {
		return "Missing command"
	}
	cmd := args[0].String()

	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Sprintf("Failed to create session: %v", err)
	}
	defer session.Close()

	out, err := session.CombinedOutput(cmd)
	if err != nil {
		return fmt.Sprintf("Command error: %v\nOutput: %s", err, string(out))
	}
	return string(out)
}

type FileNode struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
	Size  string `json:"size,omitempty"`
}

func listSFTPDirectory(this js.Value, args []js.Value) any {
	if sftpClient == nil {
		return `{"error": "SFTP not connected"}`
	}
	path := "."
	if len(args) > 0 {
		path = args[0].String()
	}

	files, err := sftpClient.ReadDir(path)
	if err != nil {
		return fmt.Sprintf(`{"error": "%v"}`, err)
	}

	var nodes []FileNode
	for _, f := range files {
		sizeStr := ""
		if !f.IsDir() {
			sizeStr = fmt.Sprintf("%d B", f.Size())
		}
		nodes = append(nodes, FileNode{
			Name:  f.Name(),
			IsDir: f.IsDir(),
			Size:  sizeStr,
		})
	}

	b, err := json.Marshal(nodes)
	if err != nil {
		return fmt.Sprintf(`{"error": "%v"}`, err)
	}
	return string(b)
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("connectSSH", js.FuncOf(connectSSH))
	js.Global().Set("executeCommand", js.FuncOf(executeCommand))
	js.Global().Set("listSFTPDirectory", js.FuncOf(listSFTPDirectory))
	fmt.Println("Beyond-SSH WASM module loaded")
	<-c
}

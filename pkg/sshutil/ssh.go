package sshutil

import (
	"context"
	"errors"
	"fmt"
	"github.com/funnyang/jump/pkg/fileutil"
	"github.com/funnyang/jump/pkg/screen"
	"io"
	"log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/containerd/console"
	"github.com/funnyang/jump/model"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	*ssh.Client
}

func Ssh(host model.Host) {

	client, err := NewSSHClient(host)
	if err != nil {
		log.Println("ssh client: ", err)
		return
	}
	defer client.Close()

	client.Ssh()
}

func NewSSHClient(host model.Host) (*SSHClient, error) {
	//创建sshp登陆配置
	auth, err := getAuth(host)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            host.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以, 但是不够安全
		Auth:            auth,
	}

	//dial 获取ssh client
	addr := fmt.Sprintf("%s:%d", host.IP, host.Port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}

	return &SSHClient{
		Client: sshClient,
	}, err
}

// getAuth 获取认证方式
func getAuth(host model.Host) (auths []ssh.AuthMethod, err error) {
	auths = []ssh.AuthMethod{}
	// 使用密码登陆
	if host.Password != "" {
		auths = append(auths, ssh.Password(host.Password))
		return
	}

	// 密钥登陆
	var signer ssh.Signer
	var key []byte
	if host.PrivateKey != "" {
		var keyPath string
		if strings.HasPrefix(host.PrivateKey, "~/") {
			dirname, _ := os.UserHomeDir()
			keyPath = filepath.Join(dirname, host.PrivateKey[2:])
		}
		if fileutil.ExistPath(keyPath) {
			key, err = os.ReadFile(keyPath)
			if err != nil {
				return nil, err
			}
		} else {
			// 使用保存的私钥登陆
			key = []byte(host.PrivateKey)
		}
	} else {
		// 使用默认的私钥登陆
		homePath, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		key, err = os.ReadFile(path.Join(homePath, ".ssh", "id_rsa"))
		if err != nil {
			return nil, err
		}
	}

	if host.PrivateKeyPhrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(host.PrivateKeyPhrase))
	} else {
		signer, err = ssh.ParsePrivateKey(key)
	}
	if err != nil {
		return nil, err
	}

	auths = append(auths, ssh.PublicKeys(signer))
	return auths, nil
}

func (sc *SSHClient) Close() {
	sc.Client.Close()
}

func (sc *SSHClient) Ssh() {
	// 获取ssh会话
	session, err := sc.NewSession()
	if err != nil {
		log.Println("创建ssh session 失败", err)
	}
	defer session.Close()

	// Set up terminal modes
	// 使用VT100终端来实现tab键提示，上下键查看历史命令，clear键清屏等操作
	c := console.Current()
	if err := c.SetRaw(); err != nil {
		log.Println(err)
		return
	}
	defer c.Reset()

	// 监听窗口变化
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go listenWindowChange(ctx, session, c)

	// 启动一个异步的管道式复制
	var pipStdIn io.WriteCloser
	if pipStdIn, err = session.StdinPipe(); err != nil {
		log.Println("ssh: ", err)
	}

	session.Stdout = screen.StdoutWriter
	session.Stderr = screen.StderrWriter
	go func() {
		buf := make([]byte, 128)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := screen.StdinReader.Read(buf)
				if err != nil {
					log.Println("stdin read: ", err)
					return
				}
				if n > 0 {
					_, err = pipStdIn.Write(buf[:n])
					if err != nil {
						if errors.Is(err, io.EOF) {
							screen.StdinWriter.Write(buf[:n])
						}
						return
					}
				}
			}
		}
	}()

	winSize, err := c.Size()
	if err != nil {
		log.Println("ssh: ", err)
		return
	}

	modes := ssh.TerminalModes{}
	if err := session.RequestPty("xterm-256color", int(winSize.Height), int(winSize.Width), modes); err != nil {
		return
	}

	if err := session.Shell(); err != nil {
		log.Println("ssh shell: ", err)
	}
	session.Wait()
}

// 监听窗口变化，并更新
func listenWindowChange(ctx context.Context, session *ssh.Session, c console.Console) {

	rawWinSize, err := c.Size()
	if err != nil {
		return
	}

	// SIGWINCH is sent to the process when the window size of the terminal has
	// changed.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)

	for {
		select {
		case <-ctx.Done():
			return
		// The client updated the size of the local PTY. This change needs to occur
		// on the server side PTY as well.
		case sigwinch := <-ch:
			if sigwinch == nil {
				return
			}

			// 获取窗口大小
			winSize, err := c.Size()
			if err != nil {
				log.Println("get console size: ", err)
				continue
			}

			if winSize.Width == rawWinSize.Width && winSize.Height == rawWinSize.Height {
				continue
			}

			// 修改远端窗口大小
			if err := session.WindowChange(int(winSize.Height), int(winSize.Width)); err != nil {
				log.Printf("Unable to send window-change reqest: %s.", err)
				continue
			}

			rawWinSize = winSize
		}
	}

}

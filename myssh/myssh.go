package myssh

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// Cli ...
type Cli struct {
	IP       string      //IP地址
	Username string      //用户名
	Password string      //密码
	Port     int         //端口号
	client   *ssh.Client //ssh客户端
}

// New 创建命令行对象
// ip IP地址
// username 用户名
// password 密码
// port 端口号,默认22
func New(ip string, username string, password string, port ...int) (*Cli, error) {
	cli := new(Cli)
	cli.IP = ip
	cli.Username = username
	cli.Password = password
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}

	return cli, cli.connect()
}

// Run 执行 shell脚本命令
func (c Cli) Run(shell string) (string, error) {
	session, err := c.newSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	buf, err := session.CombinedOutput(shell)
	return string(buf), err
}

// RunTerminal 执行带交互的命令
func (c *Cli) RunTerminal(shell string) error {
	session, err := c.newSession()
	if err != nil {
		return err
	}
	defer session.Close()

	return c.runTerminalSession(session, shell)
}

// runTerminalSession 执行带交互的命令
func (c *Cli) runTerminalSession(session *ssh.Session, shell string) error {
	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stdin
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}
	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		return err
	}

	session.Run(shell)
	return nil
}

// EnterTerminal 完全进入终端
func (c Cli) EnterTerminal() error {
	session, err := c.newSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, oldState)

	session.Stdout = os.Stdout
	session.Stderr = os.Stdin
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	return session.Wait()
}

// Enter 完全进入终端
func (c Cli) Enter(w io.Writer, r io.Reader) error {
	session, err := c.newSession()
	if err != nil {
		return err
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	// oldState, err := terminal.MakeRaw(fd)
	// if err != nil {
	// 	return err
	// }
	// defer terminal.Restore(fd, oldState)

	session.Stdout = w
	session.Stderr = os.Stdin
	session.Stdin = r

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	err = session.RequestPty("xterm-256color", termHeight, termWidth, modes)
	if err != nil {
		return err
	}

	err = session.Shell()
	if err != nil {
		return err
	}

	return session.Wait()
}

// 连接
func (c *Cli) connect() error {
	config := ssh.ClientConfig{
		User: c.Username,
		Auth: []ssh.AuthMethod{ssh.Password(c.Password)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	c.client = sshClient
	return nil
}

// newSession new session
func (c Cli) newSession() (*ssh.Session, error) {
	if c.client == nil {
		if err := c.connect(); err != nil {
			return nil, err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}

package core

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/containerd/console"
	"github.com/fatih/color"
	"github.com/funnyang/jump/conf"
	"github.com/funnyang/jump/model"
	"github.com/funnyang/jump/service"
	"github.com/funnyang/jump/sshutil"
	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/term"
	"gorm.io/gorm"
)

const prompt = "Opt> "

var srv *service.Service

type CmdHandler func(cmd string)

// InteractiveTerminal 交互终端
type InteractiveTerminal struct {
	myTerm    *term.Terminal
	c         console.Console
	editState bool // 编辑状态
}

// 交互操作入口
func Entry() {

	it := &InteractiveTerminal{}
	it.init()
	defer it.close()

	srv = service.NewService(&conf.Config{})
	defer srv.Close()

	it.forLoop()
}

func (it *InteractiveTerminal) init() (err error) {
	displayBanner(os.Stdout)

	it.c = console.Current()
	if err = it.c.SetRaw(); err != nil {
		fmt.Println(err)
		return
	}

	it.myTerm = term.NewTerminal(os.Stdin, prompt)
	return
}

func (it *InteractiveTerminal) close() {
	it.c.Reset()
	fmt.Println("")
}

func (it *InteractiveTerminal) forLoop() {
	for {

		cmd, err := it.myTerm.ReadLine()
		if err != nil {
			break
		}

		cmd = strings.TrimSpace(cmd)

		it.c.Reset()
		it.handleCmd(cmd)
		it.c.SetRaw()
	}
}

func (it *InteractiveTerminal) handleCmd(cmd string) {
	if it.isEditState() {
		it.editMode(cmd)
	} else {
		it.cmdMode(cmd)
	}
}

// cmdMode 命令模式
func (it *InteractiveTerminal) cmdMode(cmd string) {
	if strings.HasPrefix(cmd, "/") && cmd[1:] != "" {
		ListHost(cmd[1:])
	}

	switch cmd {
	case "", "p":
		ListHost("")
	case "m":
		it.enterEditState()
	case "h":
		displayMenu(os.Stdin)
	case "q":
		os.Exit(0)
	case "c":
		// 清屏
		io.WriteString(os.Stdin, CharClear)
	default:
		it.JumpHost(cmd)
	}
}

func (it *InteractiveTerminal) editMode(cmd string) {
	switch cmd {
	case "a":
		var host model.Host
		scanHost(&host)
		srv.InsertHost(host)
	case "d":
		var id int
		fmt.Printf("请输入主机ID: ")
		fmt.Scanln(&id)
		if err := srv.DeleteHost(id); err != nil {
			fmt.Println("删除失败: ", err)
		}
	case "u":
		var id int
		fmt.Printf("请输入主机ID: ")
		fmt.Scanln(&id)
		host, err := srv.GetHost(id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				fmt.Println("输入的ID有误", err)
				return
			}
			fmt.Println("查找失败: ", err)
		}
		scanHost(&host)
		if err := srv.UpdateHost(host); err != nil {
			fmt.Println("更新失败: ", err)
		}
	case "q":
		it.exitEditState()
	default:
		displayEditMenu(os.Stdin)
	}
}

func scanHost(host *model.Host) {
	// ip
	fmt.Printf("请输入主机Host: ")
	fmt.Scanln(&host.Host)
	host.Host = strings.TrimSpace(host.Host)
	for host.Host == "" {
		fmt.Printf("请输入主机Host: ")
		fmt.Scanln(&host.Host)
	}

	// 端口
	fmt.Printf("请输入主机端口(22): ")
	fmt.Scanln(&host.Port)
	if host.Port == 0 {
		host.Port = 22
	}

	// 用户名
	fmt.Printf("请输入登录用户名(root): ")
	fmt.Scanln(&host.User)
	host.User = strings.TrimSpace(host.User)
	if host.User == "" {
		host.User = "root"
	}

	// 密码
	fmt.Printf("请输入登录密码(不填则使用私钥登录): ")
	fmt.Scanln(&host.Password)
	host.Password = strings.TrimSpace(host.Password)
	if host.Password == "" {
		// 私钥
		fmt.Printf("请输入私钥(不填则使用默认私钥): ")
		fmt.Scanln(&host.PrivateKey)
		host.PrivateKey = strings.TrimSpace(host.PrivateKey)
		if host.PrivateKey != "" {
			// 私钥密码
			fmt.Printf("请输入私钥密码: ")
			fmt.Scanln(&host.PrivateKeyPhrase)
			host.PrivateKeyPhrase = strings.TrimSpace(host.PrivateKeyPhrase)
		}
	}

	// 私钥密码
	fmt.Printf("请输入主机描述: ")
	fmt.Scanln(&host.Desc)
	host.Desc = strings.TrimSpace(host.Desc)
}

// enterEditState 进入编辑状态
func (it *InteractiveTerminal) enterEditState() {
	displayEditMenu(os.Stdin)
	it.editState = true
}

// exitEditState 退出编辑状态
func (it *InteractiveTerminal) exitEditState() {
	displayMenu(os.Stdin)
	it.editState = false
}

// isEditState 是否编辑状态
func (it *InteractiveTerminal) isEditState() bool {
	return it.editState
}

func (it *InteractiveTerminal) JumpHost(keyword string) {
	// 精确匹配
	hosts, err := srv.ExactMatchHost(keyword)
	if err != nil {
		color.Red(err.Error())
		return
	}
	if len(hosts) == 1 {
		// ssh
		sshutil.Ssh(hosts[0])
		return
	}

	// 模糊匹配
	hosts, err = srv.ListHostByKeyword(keyword)
	if err != nil {
		color.Red(err.Error())
		return
	}

	if len(hosts) == 1 {
		// ssh
		sshutil.Ssh(hosts[0])
	} else {
		PrintHosts(hosts)
	}
}

// ListHost 列出所有主机
func ListHost(keyword string) {
	hosts, err := srv.ListHostByKeyword(keyword)
	if err != nil {
		color.Red(err.Error())
		return
	}

	PrintHosts(hosts)
}

// PrintHosts 打印主机列表
func PrintHosts(hosts []model.Host) {
	if len(hosts) == 0 {
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"序号", "主机", "端口", "说明"})

	rows := make([]table.Row, 0, len(hosts))
	for _, host := range hosts {
		rows = append(rows, table.Row{
			color.YellowString(strconv.Itoa(host.ID)),
			host.Host,
			host.Port,
			host.Desc,
		})
	}

	t.AppendRows(rows)
	t.Render()
}

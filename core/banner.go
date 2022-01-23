package core

import (
	"fmt"
	"io"

	"github.com/fatih/color"
)

const (
	// 清屏
	CharClear = "\x1b[H\x1b[2J"
)

var (
	menu     []MenuItem
	editMenu []MenuItem
	green    = color.New(color.FgGreen, color.Bold).SprintfFunc()
	yellow   = color.New(color.FgYellow, color.Bold).SprintfFunc()
)

type MenuItem struct {
	id          int    // id
	instruction string // 指令
	helpText    string // 帮助信息
	showText    string
}

func (mi *MenuItem) Text() string {
	return mi.text(green)
}

func (mi *MenuItem) YellowText() string {
	return mi.text(yellow)
}

func (mi *MenuItem) text(colorSprint func(format string, a ...interface{}) string) string {
	if mi.showText != "" {
		return mi.showText
	}

	mi.showText = fmt.Sprintf("\t%d) 输入 %s %s.\r\n", mi.id, colorSprint(mi.instruction), mi.helpText)
	return mi.showText
}

func init() {
	menu = []MenuItem{
		{id: 1, instruction: "ID", helpText: "直接登陆"},
		{id: 2, instruction: "部分IP、主机名、备注", helpText: "进行搜索登录(如果唯一)"},
		{id: 3, instruction: "/ + IP，主机名 or 备注", helpText: "进行搜索，如：/192.168"},
		{id: 4, instruction: "p", helpText: "显示主机列表"},
		{id: 5, instruction: "m", helpText: "进行主机管理"},
		{id: 6, instruction: "h", helpText: "显示帮助"},
		{id: 7, instruction: "q", helpText: "退出登录"},
		{id: 8, instruction: "c", helpText: "清屏"},
	}

	editMenu = []MenuItem{
		{id: 1, instruction: "a", helpText: "添加主机"},
		{id: 2, instruction: "d", helpText: "删除主机"},
		{id: 3, instruction: "u", helpText: "更新主机"},
		{id: 4, instruction: "q", helpText: "返回上级"},
	}
}

func displayBanner(out io.ReadWriter) {
	displayWelcome(out)
	displayMenu(out)
}

func displayWelcome(out io.ReadWriter) {
	title := green("你好，欢迎使用jump ")
	welcomeMsg := fmt.Sprintf("\t\t%s\r\n\r\n", title)

	io.WriteString(out, welcomeMsg)
}

func displayMenu(out io.ReadWriter) {
	for _, v := range menu {
		io.WriteString(out, v.Text())
	}
}

// displayEditMenu 显示编辑菜单
func displayEditMenu(out io.ReadWriter) {
	for _, v := range editMenu {
		io.WriteString(out, v.YellowText())
	}
}

package commands

import (
	"flag"
	"io"
	"os"
	"strings"

	"github.com/ZEL-30/zel/logger/colors"
	"github.com/ZEL-30/zel/utils"
)

// command 是一个执行单位
type Command struct {
	// Run 运行命令, args 是命令名后面的参数
	Run func(cmd *Command, args []string) int

	// 在运行命令之前执行一个操作
	PreRun func(cmd *Command, args []string)

	// 一行Usage消息, 行中的第一个字被认为是命令名
	UsageLine string

	// 'go help' 输出中显示的简短描述
	Short string

	// 在 'go help  <this-command>' 输出中显示的长消息<this-command>
	Long string

	// 一组特定于此命令的标志
	Flag flag.FlagSet

	// 命令将执行自己的标志解析
	CustomFlags bool

	// 如果在SetOutput（w）中设置，则输出写入器
	output *io.Writer
}

// 可用的命令
var AvailableCommands = []*Command{}

// 命令使用说明
var cmdUsage = `Use {{printf "zel help %s" .Name | bold}} for more information.{{endline}}`

// 返回命令的名称：Usage行中的第一个单词
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")

	if i >= 0 {
		name = name[:i]
	}

	return name
}

// SetOutput 设置Usage和错误消息的目的地
// 如果输出为nil，则使用os.Stderr
func (c *Command) SetOutput(output io.Writer) {
	c.output = &output
}

// Out 返回当前命令的输出写入器
// 如果cmd.output为nil，则使用os.Stderr
func (c *Command) Out() io.Writer {
	if c.output != nil {
		return *c.output
	}

	return colors.NewColorWriter(os.Stderr)
}

// Usage 输出命令的使用说明
func (c *Command) Usage() {
	utils.Tmpl(cmdUsage, c)
	os.Exit(2)
}

// Runnable 报告命令是否可以运行;否则它是一个文档伪命令，如import path
func (c *Command) Runnable() bool {
	return c.Run != nil
}

func (c *Command) Options() map[string]string {
	options := make(map[string]string)
	c.Flag.VisitAll(func(f *flag.Flag) {
		defaultVal := f.DefValue
		if len(defaultVal) > 0 {
			options[f.Name+"="+defaultVal] = f.Usage
		} else {
			options[f.Name] = f.Usage
		}
	})

	return options
}

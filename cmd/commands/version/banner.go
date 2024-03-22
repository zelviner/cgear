package version

import (
	"io"
	"runtime"
	"text/template"
	"time"

	"zel/logger"
	"zel/utils"
)

// RuntimeInfo 保存有关当前运行时的信息
type RuntimeInfo struct {
	OS         string
	NumCPU     int
	Compiler   string
	ZelVersion string
	Published  string
}

// InitBanner 加载横幅并打印到输出
// 所有错误都被忽略，应用程序不会在错误的情况下打印横幅
func InitBanner(out io.Writer, in io.Reader) {
	if in == nil {
		logger.Log.Fatal("The input is nil")
	}

	banner, err := io.ReadAll(in)
	if err != nil {
		logger.Log.Fatalf("Error while trying to read the banner: %s", err)
	}

	show(out, string(banner))

}

func show(out io.Writer, content string) {
	t, err := template.New("banner").Funcs(template.FuncMap{"Now": Now}).Parse(content)
	if err != nil {
		logger.Log.Fatalf("Cannot parse the banner template: %s", err)
	}

	runtimeInfo := RuntimeInfo{
		OS:         runtime.GOOS,
		NumCPU:     runtime.NumCPU(),
		Compiler:   runtime.Compiler,
		ZelVersion: version,
		Published:  utils.GetLastPublishedTime(),
	}

	err = t.Execute(out, runtimeInfo)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}

// Now 返回指定布局中的当前本地时间
func Now(layout string) string {
	return time.Now().Format(layout)
}

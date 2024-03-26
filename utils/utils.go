package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"text/template"
	"time"

	"zel/internal/pkg/system"
	"zel/logger"
	"zel/logger/colors"
)

type Repos struct {
	UpdatedAt time.Time `json:"updatad_at"`
	PushedAt  time.Time `json:"pushed_at"`
}

type Releases struct {
	PublishedAt time.Time `json:"published_at"`
	TagName     time.Time `json:"tag_name"`
}

func GetZelWorkPath() string {
	curPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return curPath
}

// IsExist返回文件或目录是否存在
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// 检查当前路径是否为 Zel tool 生成的 C++ 项目
func IsZelProject(thePath string) bool {
	cmakeListsFiles := []string{
		thePath + `\CMakeLists.txt`,
		thePath + `\src\CMakeLists.txt`,
		thePath + `\tests\CMakeLists.txt`,
	}

	files := []string{}

	c := make(chan error)

	go func() {
		filepath.Walk(thePath, func(fpath string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if f.Name() == "CMakeLists.txt" {
				files = append(files, fpath)
			}

			return nil
		})
		close(c)
	}()

	if path := <-c; path == nil {
		if reflect.DeepEqual(cmakeListsFiles, files) {
			return true
		}
	}

	return false
}

// askForConfirmation 使用Scanln解析用户输入。
// 用户必须输入“yes”或“no”，然后按回车键。它具有模糊匹配，因此“y”、“Y”、“yes”、“YES”和“Yes”都算作确认。
// 如果输入没有被识别，它会再次询问。 该函数在得到用户的有效响应之前不会返回。
// 通常，在调用askForConfirmation之前，你应该使用fmt打印出一个问题。例如：fmt.Println（“你确定吗？(yes/无）”）
func AskForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		logger.Log.Fatalf("%s", err)
	}

	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return AskForConfirmation()
	}
}

func containsString(slice []string, element string) bool {
	for _, elem := range slice {
		if elem == element {
			return true
		}
	}

	return false
}

// FuncMap 返回不同模板中使用的函数的FuncMap。
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"trim":       strings.TrimSpace,
		"bold":       colors.Bold,
		"headline":   colors.MagentaBold,
		"foldername": colors.RedBold,
		"endline":    logger.EndLine,
		"tmpltostr":  TmplToString,
	}
}

// TmplToString 解析文本模板并将结果作为字符串返回。
func TmplToString(tmpl string, data interface{}) string {
	t := template.New("tmpl").Funcs(FuncMap())
	template.Must(t.Parse(tmpl))

	var doc bytes.Buffer
	err := t.Execute(&doc, data)
	MustCheck(err)

	return doc.String()
}

func Tmpl(text string, data interface{}) {
	output := colors.NewColorWriter(os.Stderr)

	t := template.New("Usage").Funcs(FuncMap())
	template.Must(t.Parse(text))

	err := t.Execute(output, data)
	if err != nil {
		logger.Log.Error(err.Error())
	}
}

func CheckEnv(appname string) (apppath string, packpath string, err error) {
	return
}

// 当错误不为nil时，MustCheck会出现异常
func MustCheck(err error) {
	if err != nil {
		panic(err)
	}
}

// 创建文件并向其中写入内容
func WriteToFile(filename string, content string) {
	f, err := os.Create(filename)
	MustCheck(err)
	defer CloseFile(f)

	_, err = f.WriteString(content)
	MustCheck(err)
}

// 尝试关闭传递的文件, 如果出错 panic
func CloseFile(f *os.File) {
	err := f.Close()
	MustCheck(err)
}

func PrintErrorAndExit(message string, errorTemplate string) {
	Tmpl(fmt.Sprintf(errorTemplate, message), nil)
	os.Exit(2)
}

func ZelReleasesInfo() (repos []Releases) {
	var url = "https://api.github.com/repos/beego/bee/releases"
	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Warnf("Get Zel releases from github error : %s", err)
		return
	}

	defer resp.Body.Close()
	bodyContent, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyContent, &repos); err != nil {
		logger.Log.Warnf("Unmarshal releases body error: %s", err)
		return
	}

	return
}

func UpdateLastPublishedTime() {
	info := ZelReleasesInfo()
	if len(info) == 0 {
		logger.Log.Warn("Has no releases")
		return
	}
	createdAt := info[0].PublishedAt.Format("2006-01-02")
	zelHome := system.ZelHome
	if !IsExist(zelHome) {
		if err := os.MkdirAll(zelHome, 0755); err != nil {
			logger.Log.Fatalf("Could not create the directory: %s", err)
			return
		}
	}

	fp := zelHome + "/.lastPublishedAt"
	w, err := os.OpenFile(fp, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		logger.Log.Warnf("Open .lastPublishedAt file err: %s", err)
		return
	}
	defer w.Close()

	if _, err := w.WriteString(createdAt); err != nil {
		logger.Log.Warnf("Update .lastPublishedAt file err: %s", err)
		return
	}

}

func GetLastPublishedTime() string {
	fp := system.ZelHome + "/.lastPublishedAt"
	if !IsExist(fp) {
		UpdateLastPublishedTime()
	}

	w, err := os.OpenFile(fp, os.O_RDONLY, 0644)
	if err != nil {
		logger.Log.Warnf("Open .lastPublishedAt file err: %s", err)
		return "unknown"
	}
	defer w.Close()

	t := make([]byte, 1024)
	read, err := w.Read(t)
	if err != nil {
		logger.Log.Warnf("read .lastPulishedAt file err: %s", err)
		return "unknown"
	}

	return string(t[:read])
}

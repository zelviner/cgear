package test

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"
)

type Student struct {
	Sid        int
	Name       string
	Age        int
	Gender     string
	CreateTime string
}

var studentTemplate = `
|学号|姓名|年龄|性别|创建日期|
|---|---|---|---|---|
{{- range .}}
|{{.Sid}} |{{.Name | workStatus}} |{{.Age}} |{{.Gender}} |{{.CreateTime}}|
{{- end -}}`

func TestTemplate(t *testing.T) {
	s := []Student{
		{1, "张恩乐", 25, "男", "2024.3.14"},
		{2, "王仔怡", 22, "女", "2024.3.14"},
	}

	funcMap := template.FuncMap{
		"workStatus": func(name string) string {
			if name == "张恩乐" {
				return "张恩乐 (已离职)"
			}
			return name
		},
	}

	tmpl, err := template.New("Student").Funcs(funcMap).Parse(studentTemplate)
	if err != nil {
		t.Error(err)
	}
	text := new(bytes.Buffer)
	err = tmpl.Execute(text, s)
	if err != nil {
		t.Error(err)
	}

	fmt.Print(text)
}

func TestColor(t *testing.T) {
	message := "勇士总冠军"
	str := fmt.Sprintf("\x1b[31m%s\x1b[0m", message)
	fmt.Println(str)

}

package tests

import (
	"fmt"
	"regexp"
	"testing"
)

func TestRegexp(t *testing.T) {
	// 待匹配的文件列表
	files := []string{"file.h", "file.cpp", "file.hpp", "file.txt", "file"}

	// 编译正则表达式
	re := regexp.MustCompile(`\.(h|cpp|hpp)$`)

	// 遍历文件列表，匹配正则表达式
	for _, file := range files {
		if re.MatchString(file) {
			fmt.Printf("%s 匹配\n", file)
		} else {
			fmt.Printf("%s 不匹配\n", file)
		}
	}
}

func TestNmae(t *testing.T) {

	str := `Running main() from C:/Users/ZEL/Zel/pkg/googletest/googletest/src/gtest_main.cc
filesystem.
  file
  directory
Running main() from C:/Users/ZEL/Zel/pkg/googletest/googletest/src/gtest_main.cc
ftp.
  connect
Running main() from C:/Users/ZEL/Zel/pkg/googletest/googletest/src/gtest_main.cc
ini.
  configfile
Running main() from C:/Users/ZEL/Zel/pkg/googletest/googletest/src/gtest_main.cc
json.
  base_type
  array
  object
  parser
Running main() from C:/Users/ZEL/Zel/pkg/googletest/googletest/src/gtest_main.cc
logger.
  demo
data: zel size: 9
Running main() from C:/Users/ZEL/Zel/pkg/googletest/googletest/src/gtest_main.cc
thread.
  demo
Running main() from C:/Users/ZEL/Zel/pkg/googletest/googletest/src/gtest_main.cc
xml.
  class
  parser`

	// 定义一个匹配每个类别和其条目的正则表达式
	re := regexp.MustCompile(`(?m)(?P<Category>\w+)\.\n((?:\s{2}\w+\n)*)`)

	// // 查找所有匹配的子字符串及其位置
	// matches := re.FindAllStringSubmatchIndex(str, -1)
	// categories := re.SubexpNames()

	// fmt.Println(matches)
	// fmt.Println(categories)

	arr := re.FindAllStringSubmatch(str, -1)
	for _, name := range arr {
		fmt.Println(name)
	}

	// // 输出结果
	// for _, match := range matches {
	// 	for i, name := range categories {
	// 		if i != 0 && match[2*i] != -1 {
	// 			category := str[match[2*i]:match[2*i+1]]
	// 			if name == "Category" {
	// 				fmt.Printf("Category: %s\n", category)
	// 			} else {
	// 				itemsRe := regexp.MustCompile(`\s{2}(\w+)\n`)
	// 				items := itemsRe.FindAllStringSubmatch(category, -1)
	// 				fmt.Println("Items:")
	// 				for _, item := range items {
	// 					fmt.Printf("  %s\n", item[1])
	// 				}
	// 			}
	// 		}
	// 	}
	// }
}

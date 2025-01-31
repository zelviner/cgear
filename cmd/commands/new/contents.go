package new

var gitignore = `.cache
build
bin
lib
zel.json
`

var clangFormat = `# Run manually to reformat a file:
# clang-format -i --style=file <file>

Language: Cpp                                    # 语言
BasedOnStyle: LLVM                               # 基于那个配置文件
PointerAlignment: Right                          # 指针的*的挨着哪边
IndentWidth: 4                                   # 缩进宽度
ColumnLimit: 160                                 # 行最大长度
SortIncludes: false                              # 不允许排序#include
MaxEmptyLinesToKeep: 1                           # 连续的空行保留几行
ObjCSpaceAfterProperty: true                     # 在 @property 后面添加空格, \@property (readonly) 而不是 \@property(readonly).
ObjCBlockIndentWidth: 4                          # OC block后面的缩进
AllowShortFunctionsOnASingleLine: true           # 是否允许短方法单行
AllowShortIfStatementsOnASingleLine: true        # 是否允许短if单行 If true, if (a) return; 可以放到同一行
AlignTrailingComments: true                      # 注释后面是否要对齐
AlignOperands: false                             # 换行的时候对齐操作符
SpacesInSquareBrackets: false                    # 中括号两边空格 [] 
SpacesInParentheses : false                      # 小括号两边添加空格
AlignConsecutiveDeclarations: true               # 多行声明语句按照=对齐
AlignConsecutiveAssignments: true                # 连续的赋值语句以 = 为中心对齐
SpaceBeforeAssignmentOperators: true             # 等号两边的空格
SpacesInContainerLiterals: true                  # 容器类的空格 例如 OC的字典
IndentWrappedFunctionNames: true                 # 函数名后面是否要缩进
KeepEmptyLinesAtTheStartOfBlocks: true           # 在block从空行开始
BreakConstructorInitializersBeforeComma: true    # 在构造函数初始化时按逗号断行，并以冒号对齐
AllowAllParametersOfDeclarationOnNextLine: true  # 函数参数换行
SpaceAfterCStyleCast: true                       # 类型转换后面是否要空格
TabWidth: 4                                      # 制表符宽度
`

var readme = `

# {{ .ProjectName }}

ProjectName and Description

<!-- PROJECT SHIELDS -->

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

<!-- PROJECT LOGO -->
<br />

<p align="center">
  <a href="https://github.com/shaojintian/Best_README_template/">
    <img src="https://github.com/shaojintian/Best_README_template/blob/master/images/logo.png" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">"完美的"README模板</h3>
  <p align="center">
    一个"完美的"README模板去快速开始你的项目！
    <br />
    <a href="https://github.com/shaojintian/Best_README_template"><strong>探索本项目的文档 »</strong></a>
    <br />
    <br />
    <a href="https://github.com/shaojintian/Best_README_template">查看Demo</a>
    ·
    <a href="https://github.com/shaojintian/Best_README_template/issues">报告Bug</a>
    ·
    <a href="https://github.com/shaojintian/Best_README_template/issues">提出新特性</a>
  </p>

</p>


 本篇README.md面向开发者
 
## 目录

- [上手指南](#上手指南)
  - [开发前的配置要求](#开发前的配置要求)
  - [安装步骤](#安装步骤)
- [文件目录说明](#文件目录说明)
- [开发的架构](#开发的架构)
- [部署](#部署)
- [使用到的框架](#使用到的框架)
- [贡献者](#贡献者)
  - [如何参与开源项目](#如何参与开源项目)
- [版本控制](#版本控制)
- [作者](#作者)
- [鸣谢](#鸣谢)

### 上手指南

请将所有链接中的“shaojintian/Best_README_template”改为“your_github_name/your_repository”



###### 开发前的配置要求

1. xxxxx x.x.x
2. xxxxx x.x.x

###### **安装步骤**

1. Get a free API Key at [https://example.com](https://example.com)
2. Clone the repo

` + "```sh\ngit clone https://github.com/shaojintian/Best_README_template.git\n```" + `



### 开发的架构 

请阅读[ARCHITECTURE.md](https://github.com/shaojintian/Best_README_template/blob/master/ARCHITECTURE.md) 查阅为该项目的架构。

### 部署

暂无

### 使用到的框架

- [xxxxxxx](https://getbootstrap.com)
- [xxxxxxx](https://jquery.com)
- [xxxxxxx](https://laravel.com)

### 贡献者

请阅读**CONTRIBUTING.md** 查阅为该项目做出贡献的开发者。

#### 如何参与开源项目

贡献使开源社区成为一个学习、激励和创造的绝佳场所。你所作的任何贡献都是**非常感谢**的。




### 版本控制

该项目使用Git进行版本管理。您可以在repository参看当前可用版本。

### 作者

xxx@xxxx

知乎:xxxx  &ensp; qq:xxxxxx    

 *您也可以在贡献者名单中参看所有参与该项目的开发者。*

### 版权说明

该项目签署了MIT 授权许可，详情请参阅 [LICENSE.txt](https://github.com/shaojintian/Best_README_template/blob/master/LICENSE.txt)

### 鸣谢


- [GitHub Emoji Cheat Sheet](https://www.webpagefx.com/tools/emoji-cheat-sheet)
- [Img Shields](https://shields.io)
- [Choose an Open Source License](https://choosealicense.com)
- [GitHub Pages](https://pages.github.com)
- [Animate.css](https://daneden.github.io/animate.css)
- [xxxxxxxxxxxxxx](https://connoratherton.com/loaders)

<!-- links -->
[your-project-path]:shaojintian/Best_README_template
[contributors-shield]: https://img.shields.io/github/contributors/shaojintian/Best_README_template.svg?style=flat-square
[contributors-url]: https://github.com/shaojintian/Best_README_template/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/shaojintian/Best_README_template.svg?style=flat-square
[forks-url]: https://github.com/shaojintian/Best_README_template/network/members
[stars-shield]: https://img.shields.io/github/stars/shaojintian/Best_README_template.svg?style=flat-square
[stars-url]: https://github.com/shaojintian/Best_README_template/stargazers
[issues-shield]: https://img.shields.io/github/issues/shaojintian/Best_README_template.svg?style=flat-square
[issues-url]: https://img.shields.io/github/issues/shaojintian/Best_README_template.svg
[license-shield]: https://img.shields.io/github/license/shaojintian/Best_README_template.svg?style=flat-square
[license-url]: https://github.com/shaojintian/Best_README_template/blob/master/LICENSE.txt
[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=flat-square&logo=linkedin&colorB=555
[linkedin-url]: https://linkedin.com/in/shaojintian
`

var license = `MIT License 

Copyright (c) 2018 Othneil Drew

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
`

var projectCMakeLists = `# [1] 项目设置 ------------------------------------------------------
cmake_minimum_required(VERSION 3.14)
project({{ .ProjectName }} VERSION 0.1.0 LANGUAGES CXX)

# [2] vcpkg工具链配置 ------------------------------------------------
if(DEFINED ENV{ZEL_HOME})
  if(CMAKE_SIZEOF_VOID_P EQUAL 8)
    set(VCPKG_TARGET_TRIPLET "x64-windows")
  else()
    set(VCPKG_TARGET_TRIPLET "x86-windows")
  endif()

  # 设置工具链文件
  set(CMAKE_TOOLCHAIN_FILE "$ENV{ZEL_HOME}/scripts/buildsystems/vcpkg.cmake"
    CACHE STRING "Vcpkg toolchain file" FORCE)
  
  # 自动搜索vcpkg包路径
  list(APPEND CMAKE_PREFIX_PATH "$ENV{ZEL_HOME}/installed/${VCPKG_TARGET_TRIPLET}")
  message(STATUS "[1] VCPKG_TARGET_TRIPLET=${VCPKG_TARGET_TRIPLET}")
endif()

# [3] 全局安装路径配置 ----------------------------------------------
include(GNUInstallDirs)
if(DEFINED ENV{ZEL_HOME})
  set(CMAKE_INSTALL_PREFIX "$ENV{ZEL_HOME}/installed/${VCPKG_TARGET_TRIPLET}"
    CACHE PATH "Install path" FORCE)
else()
  set(CMAKE_INSTALL_PREFIX "${CMAKE_BINARY_DIR}/installed")
endif()
message(STATUS "[2] Install prefix: ${CMAKE_INSTALL_PREFIX}")

# [4] 全局编译选项 --------------------------------------------------
if(MSVC)
  add_compile_options(/utf-8 /W4 /WX)
else()
  add_compile_options(-Wall -Wextra -Werror)
endif()

# [5] 子目录添加 ----------------------------------------------------
add_subdirectory(src)
add_subdirectory(test)`

var libSrcCMakeLists = `# [1] 库基础配置 ----------------------------------------------------
set(LIB_NAME {{ .ProjectName }})
string(TOUPPER ${LIB_NAME} UPPER_LIB_NAME)

# [2] 收集源码 ------------------------------------------------------
file(GLOB_RECURSE SOURCES CONFIGURE_DEPENDS "*.cpp")

# [3] 查找依赖 ------------------------------------------------------
find_package(fmt CONFIG REQUIRED)

# [4] 创建库目标 ----------------------------------------------------
add_library(${LIB_NAME} SHARED ${SOURCES})
add_library(${PROJECT_NAME}::${LIB_NAME} ALIAS ${LIB_NAME})

# [5] 生成导出头文件（确保安装后路径正确）---------------------------
include(GenerateExportHeader)
generate_export_header(${LIB_NAME}
  BASE_NAME ${UPPER_LIB_NAME}
  EXPORT_FILE_NAME "${CMAKE_BINARY_DIR}/include/${PROJECT_NAME}/export.h"
)
target_include_directories(${LIB_NAME}
  PUBLIC 
    "$<BUILD_INTERFACE:${CMAKE_BINARY_DIR}/include>"
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>"
)

# [6] 目标属性配置 --------------------------------------------------
target_compile_definitions(${LIB_NAME} PRIVATE ${UPPER_LIB_NAME}_EXPORTS)
target_link_libraries(${LIB_NAME} PUBLIC fmt::fmt)

# [7] 安装规则 ------------------------------------------------------
install(TARGETS ${LIB_NAME} EXPORT ${LIB_NAME}Targets
    RUNTIME DESTINATION "$<$<CONFIG:Debug>:debug/>bin"
    LIBRARY DESTINATION "$<$<CONFIG:Debug>:debug/>lib"
    ARCHIVE DESTINATION "$<$<CONFIG:Debug>:debug/>lib"
)

install(DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}/"
  DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}/${LIB_NAME}
  FILES_MATCHING PATTERN "*.h" PATTERN "*.hpp"
)

install(FILES "${CMAKE_BINARY_DIR}/include/${PROJECT_NAME}/export.h"
  DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}/${PROJECT_NAME}
)

# [8] 导出配置 ------------------------------------------------------
include(CMakePackageConfigHelpers)
configure_package_config_file(
  ${CMAKE_SOURCE_DIR}/cmake/${LIB_NAME}Config.cmake.in
  ${CMAKE_CURRENT_BINARY_DIR}/${LIB_NAME}Config.cmake
  INSTALL_DESTINATION ${CMAKE_INSTALL_DATADIR}/${LIB_NAME}
)

write_basic_package_version_file(
  ${LIB_NAME}ConfigVersion.cmake
  VERSION ${PROJECT_VERSION}
  COMPATIBILITY SameMajorVersion
)

install(EXPORT ${LIB_NAME}Targets
  FILE ${LIB_NAME}Targets.cmake
  NAMESPACE ${PROJECT_NAME}::
  DESTINATION ${CMAKE_INSTALL_DATADIR}/${LIB_NAME}
)

install(FILES
  ${CMAKE_CURRENT_BINARY_DIR}/${LIB_NAME}Config.cmake
  ${CMAKE_CURRENT_BINARY_DIR}/${LIB_NAME}ConfigVersion.cmake
  DESTINATION ${CMAKE_INSTALL_DATADIR}/${LIB_NAME}
)`

var configCMakeIn = `@PACKAGE_INIT@

include("${CMAKE_CURRENT_LIST_DIR}/{{ .ProjectName }}Targets.cmake")
check_required_components({{ .ProjectName }}) 
`

var appSrcCMakeLists = `# 设置应用程序名
set(APP_NAME {{ .ProjectName }})

# 设置二进制文件输出路径
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)

# 查找头文件
file(GLOB_RECURSE HEADERS ${CMAKE_CURRENT_LIST_DIR}/*.h ${CMAKE_CURRENT_LIST_DIR}/*.hpp)

# 查找源文件
file(GLOB_RECURSE SOURCES ${CMAKE_CURRENT_LIST_DIR}/*.cpp)

# 编译应用程序
add_executable(${APP_NAME} "")

# 链接源码
target_sources(${APP_NAME}
    PRIVATE
        ${SOURCES}
    PUBLIC
        ${HEADERS}
)

# 添加编译时的宏定义
target_compile_definitions(${APP_NAME}
PUBLIC
    NOLFS  # 可能用于禁用某些与LFS（Large File Storage）相关的功能
    _CRT_SECURE_NO_WARNINGS  # 禁用对不安全函数的警告
    _WINSOCK_DEPRECATED_NO_WARNINGS  # 禁用对已弃用Winsock功能的警告
)

# 为 target 添加头文件
target_include_directories(${APP_NAME}
    PUBLIC
        ${CMAKE_CURRENT_LIST_DIR}
)

# # 为 target 添加库文件目录 (如果有需要，可以填入库文件目录路径)
# target_link_directories(${APP_NAME}
#     PUBLIC
#         path/to/libraries
# )

# # 为 target 添加需要链接的共享库 （如果有需要，可以填入共享库名字)
# TARGET_LINK_LIBRARIES(${APP_NAME}
#     PUBLIC
#         zel
# )
`

var qtSrcCMakeLists = `# 设置应用程序名
set(APP_NAME {{ .ProjectName }})

# 设置二进制文件输出路径
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)

# 设置Qt5
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)
set(CMAKE_INCLUDE_CURRENT_DIR ON)

# 寻找Qt5库
find_package(Qt5 COMPONENTS Core Gui Widgets REQUIRED)

# 设置UI文件搜索路径
list(APPEND CMAKE_AUTOUIC_SEARCH_PATHS ${CMAKE_SOURCE_DIR}/res/ui)

# 查找源文件
file(GLOB_RECURSE SOURCES ${CMAKE_CURRENT_LIST_DIR}/*.cpp ${CMAKE_CURRENT_LIST_DIR}/*.hpp)

# 添加资源文件
file(GLOB RESOURCES ${CMAKE_SOURCE_DIR}/res/rc/*)

# 添加可执行文件
add_executable(${APP_NAME}  ${SOURCES} ${RESOURCES}) #debug
# add_executable(${APP_NAME} WIN32 ${SOURCES} ${RESOURCES}  ${MY_VERSIONINFO_RC}) # release

include_directories(${Qt5Gui_PRIVATE_INCLUDE_DIRS})

# 添加编译时的宏定义
target_compile_definitions(${APP_NAME}
PUBLIC
    NOLFS  # 可能用于禁用某些与LFS（Large File Storage）相关的功能
    _CRT_SECURE_NO_WARNINGS  # 禁用对不安全函数的警告
    _WINSOCK_DEPRECATED_NO_WARNINGS  # 禁用对已弃用Winsock功能的警告
)

# # 为target添加头文件
# target_include_directories(${APP_NAME}
#     PUBLIC
# )

# # 为 target 添加库文件目录 (如果有需要，可以填入库文件目录路径)
# target_link_directories(${APP_NAME}
#     PUBLIC
#         path/to/libraries
# )

# 为target添加需要链接的共享库
TARGET_LINK_LIBRARIES(${APP_NAME}
    PRIVATE
        Qt5::Core
        Qt5::Gui
        Qt5::Widgets
)`

var testCMakeLists = `# [1] 测试配置 -----------------------------------------------------
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin/test)
enable_testing()

# [2] 查找依赖 -----------------------------------------------------
find_package(GTest REQUIRED)

# [3] 添加测试目标 --------------------------------------------------
function(add_integration_test name)
  set(TEST_NAME "${name}-test")
  file(GLOB_RECURSE files ${name}/*.cpp)
  add_executable(${TEST_NAME} ${files})
  target_link_libraries(${TEST_NAME}
    PRIVATE 
      GTest::gtest_main
       ${ARGN}
  )
  add_test(NAME ${TEST_NAME} COMMAND ${TEST_NAME})
endfunction()

# [4] 添加具体测试 --------------------------------------------------
`

var appTestCMakeLists = `# 设置测试程序的输出目录
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin/test)

# 查找 GTest 库
find_package(GTest REQUIRED)

# 启用测试
enable_testing()

# 定义添加测试执行文件的函数
function(add_test_executable name)
    file(GLOB_RECURSE files ${name}/*.cpp)
    add_executable(${name}-test ${files})
    target_include_directories(${name}-test 
        PUBLIC
    )
    target_link_libraries(${name}-test
        PUBLIC
            GTest::gtest_main
            ${ARGN}
    )
endfunction(add_test_executable name)

# 添加测试
`

var projectHeader = `#pragma once

#include "utils/utils.h"
`

var utilsHeader = `#pragma once

#include "{{ .ProjectName }}/export.h"

namespace {{ .ProjectName }} {

    {{ .ProjectNameUpper }}_EXPORT void print_hello();

}
`

var utilsCPP = `#include "utils.h"
#include <iostream>

namespace {{ .ProjectName }} {

void print_hello() { std::cout << "Hello, world!" << std::endl; }

} // namespace {{ .ProjectName }}
`

var mainCPP = `#include <iostream>

int main(int argc, char *argv[]) {

    std::cout << "Welcome to zel!" << std::endl;

    return 0;
}
`

var staticLibInfo = `# 创建静态库
add_library(${LIB_NAME} ${SOURCES})`

var dynamicLibInfo = `# 创建动态库
add_library(${LIB_NAME} SHARED ${SOURCES})`

var launch = `{
    // 使用 IntelliSense 了解相关属性。 
    // 悬停以查看现有属性的描述。
    // 欲了解更多信息，请访问: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        //{{ .configuration }}
    ]
}
`

var testContent = `#include <gtest/gtest.h>

TEST({{ .testName }}, class) {
 

}`

var testLaunch = `{
            "type": "lldb",
            "request": "launch",
            "name": "{{ .testName }}-test",
            "program": "${workspaceFolder}/build/test/{{ .testName }}-test.exe",
            "args": [],
            "cwd": "${workspaceFolder}"
        },
        //{{ .configuration }}
`

var imagesRC = `<!DOCTYPE RCC><RCC version="1.0">
 <qresource>
     <!-- <file>../images/logo.ico</file> -->
 </qresource>
 </RCC>
`

var logoRc = `// IDI_ICON1 ICON "../images/logo.ico"`

var mainWindowUI = `<?xml version="1.0" encoding="UTF-8"?>
<ui version="4.0">
 <class>MainWindow</class>
 <widget class="QMainWindow" name="MainWindow">
  <property name="geometry">
   <rect>
    <x>0</x>
    <y>0</y>
    <width>800</width>
    <height>600</height>
   </rect>
  </property>
  <property name="windowTitle">
   <string>MainWindow</string>
  </property>
  <widget class="QWidget" name="centralwidget">
   <widget class="QPushButton" name="push_btn">
    <property name="geometry">
     <rect>
      <x>240</x>
      <y>160</y>
      <width>211</width>
      <height>171</height>
     </rect>
    </property>
    <property name="text">
     <string>PushButton</string>
    </property>
   </widget>
  </widget>
  <widget class="QMenuBar" name="menubar">
   <property name="geometry">
    <rect>
     <x>0</x>
     <y>0</y>
     <width>800</width>
     <height>23</height>
    </rect>
   </property>
  </widget>
  <widget class="QStatusBar" name="statusbar"/>
 </widget>
 <resources/>
 <connections/>
</ui>`

var templateUI = `<?xml version="1.0" encoding="UTF-8"?>
<ui version="4.0">
 <class>Template</class>
 <widget class="QWidget" name="Template">
  <property name="geometry">
   <rect>
    <x>0</x>
    <y>0</y>
    <width>800</width>
    <height>600</height>
   </rect>
  </property>
  <property name="windowTitle">
   <string>Template</string>
  </property>
 </widget>
 <resources/>
 <connections/>
</ui>

`

var qtMainCPP = `#include "app/main_window.h"

#include <QApplication>
#pragma comment(lib, "user32.lib")

int main(int argc, char *argv[]) {

    // 设置高DPI
    QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);

    QApplication a(argc, argv);
    MainWindow   w;
    w.show();
    return a.exec();
}
`

var appHeader = `#pragma once
#include "ui_main_window.h"
#include <QMainWindow>

class MainWindow : public QMainWindow {
    Q_OBJECT

  public:
    MainWindow(QMainWindow *parent = nullptr);
    ~MainWindow();

  private:
    // 初始化窗口
    void initWindow();

    // 初始化UI
    void initUI();

    /// @brief 初始化信号槽
    void initSignalSlot();

  private:
    Ui_MainWindow *ui_;
};
`
var appCPP = `#include "main_window.h"

MainWindow::MainWindow(QMainWindow *parent)
    : QMainWindow(parent)
    , ui_(new Ui_MainWindow) {
    ui_->setupUi(this);

    initWindow();

    initUI();

    initSignalSlot();
}

MainWindow::~MainWindow() { delete ui_; }

void MainWindow::initWindow() {

    // 设置窗口标题
    setWindowTitle("Template");
}

void MainWindow::initUI() {
    // 插入图片
    QPixmap pixmap(":/image/data.png");
    ui_->push_btn->setIcon(pixmap);
    ui_->push_btn->setIconSize(pixmap.size());
    ui_->push_btn->setFixedSize(pixmap.size());
}

void MainWindow::initSignalSlot() {}
`

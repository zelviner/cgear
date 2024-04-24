package new

var clangFormat = `# Run manually to reformat a file:
# clang-format -i --style=file <file>

Language: Cpp

BasedOnStyle: LLVM  #基于那个配置文件
PointerAlignment: Right  #指针的*的挨着哪边
IndentWidth: 4  #缩进宽度
ColumnLimit: 160  #行最大长度
SortIncludes: true  #允许排序#include
MaxEmptyLinesToKeep: 1  #连续的空行保留几行
ObjCSpaceAfterProperty: true  #在 @property 后面添加空格, \@property (readonly) 而不是 \@property(readonly).
ObjCBlockIndentWidth: 4  #OC block后面的缩进
AllowShortFunctionsOnASingleLine: true  #是否允许短方法单行
AllowShortIfStatementsOnASingleLine: true  #是否允许短if单行 If true, if (a) return; 可以放到同一行
AlignTrailingComments: true  #注释对齐
AlignOperands: false  #换行的时候对齐操作符
SpacesInSquareBrackets: false  #中括号两边空格 [] 
SpacesInParentheses : false  #小括号两边添加空格
AlignConsecutiveDeclarations: true  #多行声明语句按照=对齐
AlignConsecutiveAssignments: true  #连续的赋值语句以 = 为中心对齐
SpaceBeforeAssignmentOperators: true  #等号两边的空格
SpacesInContainerLiterals: true  #容器类的空格 例如 OC的字典
IndentWrappedFunctionNames: true  #缩进
KeepEmptyLinesAtTheStartOfBlocks: true  #在block从空行开始
BreakConstructorInitializersBeforeComma: true  #在构造函数初始化时按逗号断行，并以冒号对齐
AllowAllParametersOfDeclarationOnNextLine: true  #函数参数换行
SpaceAfterCStyleCast: true  #括号后添加空格
TabWidth: 4  #tab键盘的宽度
`

var readme = `# {{.ProjectName}}
`

var buildBat = `@echo off
setlocal

rem 避免中文乱码
chcp 65001

rem 检查是否提供了可执行文件名称参数
if "%1"=="" (
    echo.
    echo 请提供可执行文件名称作为参数。例如： .\run.bat hello
    goto :eof
)

rem 设置 CMake 和编译目录
set BUILD_DIR=build
set EXECUTABLE_PATH=bin
set EXECUTABLE_NAME=%1.exe

rem 设置 Visual Studio 环境变量
set VS_PATH="D:\Development\Microsoft Visual Studio\VC\Auxiliary\Build\vcvars32.bat"
if not exist %VS_PATH% (
    set VS_PATH="C:\Development\Microsoft Visual Studio\VC\Auxiliary\Build\vcvars32.bat"
)
call %VS_PATH%

rem 检查编译目录是否存在，如果存在则删除，然后创建
if not exist %BUILD_DIR% (
    mkdir %BUILD_DIR%
) else (
    rem 最后一个参数为 rebuild 表示删除原有的编译目录
    if "%2"=="rebuild" (
        echo Rebuilding...
        rmdir /s /q %BUILD_DIR%
        mkdir %BUILD_DIR%
    )
)

rem 进入编译目录
cd %BUILD_DIR%

rem 使用 CMake 生成项目文件
"cmake" -G "Ninja" -DCMAKE_EXPORT_COMPILE_COMMANDS=ON ..

rem 使用 Ninja 进行编译
ninja

endlocal
`

var projectCmakeLists = `# 最低版本
cmake_minimum_required(VERSION 3.20.2) 

# 设置项目名称
project({{.ProjectName}})

# 采用C++11标准
set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

if(MSVC)
    add_definitions(-D_CRT_SECURE_NO_WARNINGS -D_CRT_NONSTDC_NO_DEPRECATE)
    # Specify MSVC UTF-8 encoding   
    add_compile_options("$<$<C_COMPILER_ID:MSVC>:/utf-8>")
    add_compile_options("$<$<CXX_COMPILER_ID:MSVC>:/utf-8>")
    set(CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS_RELEASE} /MD")    
endif()

set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)

# 添加子工程
add_subdirectory(vendor)
add_subdirectory(src)
add_subdirectory(tests)
`

var srcCmakeLists = `# 查找源文件
file(GLOB_RECURSE SOURCES ${CMAKE_CURRENT_LIST_DIR}/*.cpp ${CMAKE_CURRENT_LIST_DIR}/*.hpp)

# 查找头文件
file(GLOB_RECURSE HEADERS ${CMAKE_CURRENT_LIST_DIR}/*.h)

#  编译静态库
add_library(${PROJECT_NAME} "")

target_sources(${PROJECT_NAME}
PRIVATE
    ${SOURCES}
PUBLIC
    ${HEADERS}
)

# 添加头文件
target_include_directories(${PROJECT_NAME}
PUBLIC
    ${CMAKE_CURRENT_LIST_DIR}
   
)

# 为target添加库文件目录
target_link_directories(${PROJECT_NAME}
PUBLIC
  
)



# 为target添加需要链接的共享库
TARGET_LINK_LIBRARIES(${PROJECT_NAME}
PUBLIC
   
)

# 添加自定义命令，用于复制静态库文件到指定目录
# add_custom_command(TARGET ${PROJECT_NAME} POST_BUILD
#     COMMAND ${CMAKE_COMMAND} -E copy
#         $<TARGET_FILE:${PROJECT_NAME}>
#         "填写要发送的路径 如: D:/Workspaces/C++/Vendor/xhlanguage/lib"
#     COMMENT "Copying ${PROJECT_NAME} static library to destination directory"
# )


`

var vendorCmakeLists = `
`

var testsCmakeLists = `
`

var utilsHeader = `#pragma once
`

var utilsCPP = `#include "utils.h"
`

var testContent = `#include <iostream>

int main() {

    std::cout << "This is a test program" << std::endl;

    return 0;
}`

var testCmakeLists = `# ---------- {{ .TestName }} ----------
add_executable({{ .TestName }}-test {{ .TestName }}/{{ .TestName }}_test.cpp)
target_include_directories({{ .TestName }}-test 
PUBLIC
)
target_link_libraries({{ .TestName }}-test
PUBLIC
)
`

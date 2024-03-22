package new

var projectCmakeList = `# 最低版本
cmake_minimum_required(VERSION 3.20.2) 

# 设置项目名称
project({{.ProjectName}})

# 采用C++17标准
set(CMAKE_CXX_STANDARD 17)
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
add_subdirectory(test)
`

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

var runBat = `@echo off
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

rem 如果编译成功，运行生成的可执行文件
if %ERRORLEVEL%==0 (
    rem 清空屏幕
    cls
    rem 进入可执行文件所在目录
    cd ..\%EXECUTABLE_PATH%
    echo Running %EXECUTABLE_NAME% ...
    echo.
    %EXECUTABLE_NAME%
    echo.
    echo Process exited with code %ERRORLEVEL%
) else (
    echo.
    echo Build failed!
    echo.
)
endlocal
`

var utilsHeader = `

`

var utilsCPP = `
`

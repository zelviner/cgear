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

var projectCmakeLists = `# 最低版本
cmake_minimum_required(VERSION 3.14) 

# 设置项目名称
project(zeltest)

# 采用C++14标准
set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(ZELPATH $ENV{ZELPATH}/${CMAKE_BUILD_TYPE})

# 设置安装路径
set(CMAKE_INSTALL_PREFIX ${ZELPATH})

if(WIN32)
    set(WINDOWS_EXPORT_ALL_SYMBOLS ON)
endif()

if(MSVC)
    add_definitions(-D_CRT_SECURE_NO_WARNINGS -D_CRT_NONSTDC_NO_DEPRECATE)
    # Specify MSVC UTF-8 encoding   
    add_compile_options("$<$<C_COMPILER_ID:MSVC>:/utf-8>")
    add_compile_options("$<$<CXX_COMPILER_ID:MSVC>:/utf-8>")
    set(CMAKE_CXX_FLAGS_RELEASE "${CMAKE_CXX_FLAGS_RELEASE} /MD")    
endif()

# 设置三方库的安装路径, 搜索路径, 链接路径
list(APPEND CMAKE_PREFIX_PATH ${ZELPATH})
include_directories(${ZELPATH}/include)
link_directories(${ZELPATH}/lib)

enable_testing()

# 添加子工程
add_subdirectory(src)
add_subdirectory(test)
`

var srcCmakeLists = `set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)

# 查找源文件
file(GLOB_RECURSE SOURCES ${CMAKE_CURRENT_LIST_DIR}/*.cpp)

# 查找头文件
file(GLOB_RECURSE HEADERS ${CMAKE_CURRENT_LIST_DIR}/*.h ${CMAKE_CURRENT_LIST_DIR}/*.hpp)

# 编译静态库
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

# 为 target 添加库文件目录
# 如果有需要，可以填入库文件目录路径
# target_link_directories(${PROJECT_NAME}
#     PUBLIC
#         path/to/libraries
# )

# 为 target 添加需要链接的共享库
# 如果有需要，可以填入共享库名字
# TARGET_LINK_LIBRARIES(${PROJECT_NAME}
#     PUBLIC
#         library_name
# )

# 安装目标文件
install(TARGETS ${PROJECT_NAME}
    ARCHIVE DESTINATION lib
    LIBRARY DESTINATION lib
    RUNTIME DESTINATION bin
)

# 安装目录
install(DIRECTORY ${CMAKE_CURRENT_LIST_DIR}/ DESTINATION include/${PROJECT_NAME}
    FILES_MATCHING PATTERN "*.h"
    PATTERN "*.hpp"
)
`

var testCmakeLists = `function(add_test_executable name)
    file(GLOB_RECURSE files ${name}/*.cpp)
    add_executable(${name}-test ${files})
    target_include_directories(${name}-test 
        PUBLIC
    )
    target_link_libraries(${name}-test
        PUBLIC
            ${PROJECT_NAME}
            GTest::gtest_main
            ${ARGN}
    )
endfunction(add_test_executable name)

find_package(GTest REQUIRED)

enable_testing()

# 添加测试
`

var utilsHeader = `#pragma once
`

var utilsCPP = `#include "utils.h"
`

var testContent = `#include <gtest/gtest.h>

// Demonstrate some basic assertions.
TEST({{ .testFile }}, demo) {
  // Expect two strings not to be equal.
  EXPECT_STRNE("hello", "ZEL");
  // Expect equality.
  EXPECT_EQ(7 * 6, 42);
}`

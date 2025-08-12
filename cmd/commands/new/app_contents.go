package new

var appSrcCMakeLists = `# [1] 基础配置 -----------------------------------------------------
set(APP_NAME {{ .ProjectName }})
string(TOUPPER ${APP_NAME} UPPER_LIB_NAME)

# [2] 输出目录配置 --------------------------------------------------
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)

# [3] 收集源码 -----------------------------------------------------
file(GLOB_RECURSE SOURCES CONFIGURE_DEPENDS "*.cpp")

# [4] 查找依赖 ------------------------------------------------------
# find_package(fmt CONFIG REQUIRED)

# [5] 添加可执行文件 -------------------------------------------------
add_executable(${APP_NAME} ${SOURCES})

# [6] 链接依赖库 ----------------------------------------------------
target_link_libraries(${APP_NAME} PUBLIC 
    # fmt::fmt  # fmt库，用于格式化输出
)

# [7] 添加宏定义 ----------------------------------------------------
target_compile_definitions(${APP_NAME}
PUBLIC
    NOLFS  # 可能用于禁用某些与LFS（Large File Storage）相关的功能
    _CRT_SECURE_NO_WARNINGS  # 禁用对不安全函数的警告
    _WINSOCK_DEPRECATED_NO_WARNINGS  # 禁用对已弃用Winsock功能的警告
)`

var appMainCPP = `#include <iostream>

int main(int argc, char *argv[]) {

    std::cout << "Welcome to zel!" << std::endl;

    return 0;
}
`

var appUtilsHeader = `#pragma once

namespace {{ .ProjectName }} {

    void print_hello();

}
`

var appUtilsCPP = `#include "utils.h"
#include <iostream>

namespace {{ .ProjectName }} {

void print_hello() { std::cout << "Hello, world!" << std::endl; }

} // namespace {{ .ProjectName }}
`

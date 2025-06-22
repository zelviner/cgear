package new

var staticLibInfo = `# [5] 创建库目标 ----------------------------------------------------
add_library(${LIB_NAME} ${SOURCES})
add_library(${LIB_NAME}::${LIB_NAME} ALIAS ${LIB_NAME})`

var dynamicLibInfo = `# [5] 创建库目标 ----------------------------------------------------
add_library(${LIB_NAME} SHARED ${SOURCES})
add_library(${LIB_NAME}::${LIB_NAME} ALIAS ${LIB_NAME})`

var libSrcCMakeLists = `# [1] 库基础配置 ----------------------------------------------------
set(LIB_NAME {{ .ProjectName }})

# [2] 输出目录配置 --------------------------------------------------
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)

# [3] 收集源码 ------------------------------------------------------
file(GLOB_RECURSE SOURCES CONFIGURE_DEPENDS "*.cpp")

# [4] 查找依赖 ------------------------------------------------------
find_package(fmt CONFIG REQUIRED)

{{ .LibInfo }}

# [5.1] 设置输出 DLL 名（根据配置不同使用 zeld.dll 或 zel.dll）----------
set_target_properties(${LIB_NAME} PROPERTIES
    OUTPUT_NAME_DEBUG "${LIB_NAME}d"
    OUTPUT_NAME_RELEASE "${LIB_NAME}"
    OUTPUT_NAME_RELWITHDEBINFO "${LIB_NAME}"
    OUTPUT_NAME_MINSIZEREL "${LIB_NAME}"
)

# [6] 生成导出头文件（确保安装后路径正确）-------------------------------
include(GenerateExportHeader)
generate_export_header(${LIB_NAME}
  BASE_NAME ${UPPER_LIB_NAME}
  EXPORT_FILE_NAME "${CMAKE_BINARY_DIR}/include/${LIB_NAME}/export.h"
)

target_include_directories(${LIB_NAME}
  PUBLIC 
    "$<BUILD_INTERFACE:${CMAKE_BINARY_DIR}/include>"
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>"
)

# [6] 链接依赖库 ----------------------------------------------------
target_link_libraries(${LIB_NAME} PUBLIC 
  openssl/VC/libcrypto32MD
  openssl/VC/libssl32MD
  mysql/libmysql
  Crypt32
  fmt::fmt
)

# [7] 添加宏定义 --------------------------------------------------
target_compile_definitions(${LIB_NAME} PRIVATE 
  ${UPPER_LIB_NAME}_EXPORTS        # 定义一个宏，用于区分导出和导入
  NOLFS                            # 可能用于禁用某些与LFS（Large File Storage）相关的功能
  _CRT_SECURE_NO_WARNINGS          # 禁用对不安全函数的警告
  _WINSOCK_DEPRECATED_NO_WARNINGS  # 禁用对已弃用Winsock功能的警告
  BUILDING_DLL                     # 定义一个宏，用于区分动态库和静态库
)

# [8] 安装规则 ------------------------------------------------------
install(TARGETS ${LIB_NAME} EXPORT ${LIB_NAME}Targets
  RUNTIME DESTINATION "$<$<CONFIG:Debug>:debug/>bin"
  LIBRARY DESTINATION "$<$<CONFIG:Debug>:debug/>lib"
  ARCHIVE DESTINATION "$<$<CONFIG:Debug>:debug/>lib"
)

install(DIRECTORY "${CMAKE_CURRENT_SOURCE_DIR}/"
  DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}/${LIB_NAME}
  FILES_MATCHING PATTERN "*.h" PATTERN "*.hpp"
)

install(FILES "${CMAKE_BINARY_DIR}/include/${LIB_NAME}/export.h"
  DESTINATION ${CMAKE_INSTALL_INCLUDEDIR}/${LIB_NAME}
)

# [9] 导出配置 ------------------------------------------------------
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
  NAMESPACE ${LIB_NAME}::
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

var libUtilsHeader = `#pragma once

#include "{{ .ProjectName }}/export.h"

namespace {{ .ProjectName }} {

    {{ .ProjectNameUpper }}_EXPORT void print_hello();

}
`

var libUtilsCPP = `#include "utils.h"
#include <iostream>

namespace {{ .ProjectName }} {

void print_hello() { std::cout << "Hello, world!" << std::endl; }

} // namespace {{ .ProjectName }}
`

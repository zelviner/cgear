package test

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

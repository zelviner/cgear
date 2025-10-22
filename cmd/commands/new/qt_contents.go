package new

var qtSrcCMakeLists = `# [1] 设置应用程序名称 ------------------------------------------------
set(APP_NAME {{ .ProjectName }})
string(TOUPPER ${APP_NAME} UPPER_APP_NAME)

# [2] 输出目录配置 -----------------------------------------------------
set(CMAKE_ARCHIVE_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/lib)
set(CMAKE_LIBRARY_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)
set(CMAKE_RUNTIME_OUTPUT_DIRECTORY ${CMAKE_SOURCE_DIR}/bin)

# [3] Qt 编译设置 -----------------------------------------------------
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)
set(CMAKE_INCLUDE_CURRENT_DIR ON)

# [4] 查找 Qt5 组件 ---------------------------------------------------
find_package(Qt5 COMPONENTS Core Gui Widgets REQUIRED)
# find_package(fmt CONFIG REQUIRED)


# [5] UI 搜索路径 -----------------------------------------------------
list(APPEND CMAKE_AUTOUIC_SEARCH_PATHS ${CMAKE_SOURCE_DIR}/res/ui)

# [6] 查找源文件/资源 --------------------------------------------------
file(GLOB_RECURSE SOURCES CONFIGURE_DEPENDS 
    ${CMAKE_CURRENT_LIST_DIR}/*.cpp 
    ${CMAKE_CURRENT_LIST_DIR}/*.hpp
)

file(GLOB RESOURCES CONFIGURE_DEPENDS ${CMAKE_SOURCE_DIR}/res/rc/*)

# [7] 添加可执行文件 ---------------------------------------------------
add_executable(${APP_NAME}  ${SOURCES} ${RESOURCES}) #debug
# add_executable(${APP_NAME} WIN32 ${SOURCES} ${RESOURCES}  ${MY_VERSIONINFO_RC}) # release

# [8] 包含 Qt 的私有头（如果有需要，比如自定义窗口样式）--------------
include_directories(${Qt5Gui_PRIVATE_INCLUDE_DIRS})

# [9] 添加宏定义 -------------------------------------------------------
target_compile_definitions(${APP_NAME} PUBLIC
    NOLFS
    _CRT_SECURE_NO_WARNINGS
    _WINSOCK_DEPRECATED_NO_WARNINGS
)

# [10] 链接 Qt 库 ------------------------------------------------------
target_link_libraries(${APP_NAME} PRIVATE
    Qt5::Core
    Qt5::Gui
    Qt5::Widgets
    # fmt::fmt # fmt库，用于格式化输出
)`

var qtMainWindowUI = `<?xml version="1.0" encoding="UTF-8"?>
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

var qtTemplateUI = `<?xml version="1.0" encoding="UTF-8"?>
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

#include <qapplication>
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

var qtImageRC = `<!DOCTYPE RCC><RCC version="1.0">
 <qresource>
     <!-- <file>../image/logo.ico</file> -->
 </qresource>
 </RCC>
`

var qtLogoRc = `// IDI_ICON1 ICON "../image/logo.ico"`

var qtMainWindowHeader = `#pragma once
#include "ui_main_window.h"
#include <qmainwindow>

class MainWindow : public QMainWindow {
    Q_OBJECT

  public:
    MainWindow(QMainWindow *parent = nullptr);
    ~MainWindow();

  private:
    /// @brief 初始化窗口
    void init_window();

    /// @brief 初始化 UI
    void init_ui();

    /// @brief 初始化信号槽
    void init_signals_slots();

  private:
    Ui_MainWindow *ui_;
};
`
var qtMainWindowCPP = `#include "main_window.h"

MainWindow::MainWindow(QMainWindow *parent)
    : QMainWindow(parent)
    , ui_(new Ui_MainWindow) {
    ui_->setupUi(this);

    init_window();

    init_ui();

    init_signals_slots();
}

MainWindow::~MainWindow() { delete ui_; }

void MainWindow::init_window() {
    // 设置窗口标题
    setWindowTitle("Cgear Window");
}

void MainWindow::init_ui() {
    // 插入图片
    // QPixmap pixmap(":/image/data.png");
    // ui_->push_btn->setIcon(pixmap);
    // ui_->push_btn->setIconSize(pixmap.size());
    // ui_->push_btn->setFixedSize(pixmap.size());
    ui_->push_btn->setText("欢迎使用 Cgear C++ 脚手架");
}

void MainWindow::init_signals_slots() {}
`

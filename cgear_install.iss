[Setup]
AppName=Cgear
AppVersion=2.0
DefaultDirName={pf}\Cgear
DefaultGroupName=Cgear
OutputDir=.
OutputBaseFilename=cgear-x86_64-pc-windows-setup
Compression=lzma
SolidCompression=yes
AlwaysRestart=yes
; PrivilegesRequired=admin

[Files]
Source: "cgear.exe"; DestDir: {app}\bin

[Run]
Filename: "{app}\bin\cgear.exe"; Description: "Run Cgear"; Flags: nowait postinstall skipifsilent

[Registry]
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}\bin";
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "CGEAR_HOME"; ValueData: "{app}"; Flags: createvalueifdoesntexist uninsdeletevalue

[Code]
function CheckGitInstalled: Boolean;
var
  ResultCode: Integer;
begin
  Result := False;
  if Exec('git', '--version', '', SW_HIDE, ewWaitUntilTerminated, ResultCode) then
    Result := (ResultCode = 0);
end;

function CopyFolderOrFile(SourcePath, DestPath: string): Boolean;
var
  ErrorCode: Integer;
begin
  Result := ShellExec('', 'xcopy', '"' + SourcePath + '" "' + DestPath + '" /E /I /Y /Q', '', SW_HIDE, ewWaitUntilTerminated, ErrorCode);
  if not Result or (ErrorCode <> 0) then
  begin
    MsgBox('复制失败：' + SourcePath + ' → ' + DestPath + '，错误代码：' + IntToStr(ErrorCode), mbError, MB_OK);
    Result := False;
  end;
end;

function DeleteFolderOrFile(TargetPath: string): Boolean;
var
  ErrorCode: Integer;
begin
  if DirExists(TargetPath) then
    Result := ShellExec('', 'cmd.exe', '/C rd /S /Q "' + TargetPath + '"', '', SW_HIDE, ewWaitUntilTerminated, ErrorCode)
  else if FileExists(TargetPath) then
    Result := ShellExec('', 'cmd.exe', '/C del /F /Q "' + TargetPath + '"', '', SW_HIDE, ewWaitUntilTerminated, ErrorCode)
  else
  begin
    Result := True;
    exit;
  end;

  if not Result or (ErrorCode <> 0) then
  begin
    MsgBox('删除失败：' + TargetPath + '，错误代码：' + IntToStr(ErrorCode), mbError, MB_OK);
    Result := False;
  end;
end;

procedure DownloadFile(const URL, DestFile: string);
var
  ResultCode: Integer;
begin
  if not Exec('powershell',
    '-Command "Invoke-WebRequest -Uri ''' + URL + ''' -OutFile ''' + DestFile + ''' -UseBasicParsing"',
    '', SW_HIDE, ewWaitUntilTerminated, ResultCode) then
  begin
    MsgBox('无法执行下载: ' + URL, mbError, MB_OK);
  end
  else if ResultCode <> 0 then
  begin
    MsgBox('下载失败: ' + URL + ' 错误代码 ' + IntToStr(ResultCode), mbError, MB_OK);
  end;
end;

procedure ExtractZip(const ZipFile, DestDir: string);
var
  ResultCode: Integer;
begin
  if not Exec('powershell',
    '-Command "Expand-Archive -Path ''' + ZipFile + ''' -DestinationPath ''' + DestDir + ''' -Force"',
    '', SW_HIDE, ewWaitUntilTerminated, ResultCode) then
  begin
    MsgBox('解压失败: ' + ZipFile, mbError, MB_OK);
  end
  else if ResultCode <> 0 then
  begin
    MsgBox('解压失败: ' + ZipFile + ' 错误代码 ' + IntToStr(ResultCode), mbError, MB_OK);
  end;
end;

procedure CurStepChanged(CurStep: TSetupStep);
var
  cmakeZip, ninjaZip, vcpkgDir: string;
  ResultCode: Integer;
begin
  if CurStep = ssInstall then
  begin
    // 下载 CMake
    // WizardForm.StatusLabel.Caption := '正在下载 CMake...';
    // WizardForm.ProgressGauge.Position := 0;
    // cmakeZip := ExpandConstant('{tmp}\cmake.zip');
    // DownloadFile('https://github.com/Kitware/CMake/releases/download/v3.29.0/cmake-3.29.0-windows-x86_64.zip', cmakeZip);
    // WizardForm.ProgressGauge.Position := 20;
    // ExtractZip(cmakeZip, ExpandConstant('{app}'));

    // 下载 Ninja
    WizardForm.StatusLabel.Caption := '正在下载 Ninja...';
    WizardForm.ProgressGauge.Position := 0;
    ninjaZip := ExpandConstant('{tmp}\ninja.zip');
    DownloadFile('https://github.com/ninja-build/ninja/releases/download/v1.12.1/ninja-win.zip', ninjaZip);
    WizardForm.ProgressGauge.Position := 40;
    ExtractZip(ninjaZip, ExpandConstant('{app}\bin'));

    // 下载 vcpkg
    WizardForm.StatusLabel.Caption := '正在下载 vcpkg...';
    if not CheckGitInstalled then
    begin
      MsgBox('未检测到 Git 安装，请先安装 Git。', mbError, MB_OK);
      exit;
    end;

    vcpkgDir := ExpandConstant('{app}\vcpkg');
    if not Exec('git', 'clone https://github.com/microsoft/vcpkg.git "' + vcpkgDir + '"', '', SW_HIDE, ewWaitUntilTerminated, ResultCode) then
    begin
      MsgBox('无法执行 Git 克隆命令', mbError, MB_OK);
      exit;
    end;
    WizardForm.ProgressGauge.Position := 70;

    // 编译 vcpkg
    if not Exec(ExpandConstant(vcpkgDir + '\bootstrap-vcpkg.bat'), '', vcpkgDir, SW_HIDE, ewWaitUntilTerminated, ResultCode) then
      MsgBox('执行 bootstrap-vcpkg.bat 失败', mbError, MB_OK)
    else if ResultCode <> 0 then
      MsgBox('执行 bootstrap-vcpkg.bat 失败，错误代码: ' + IntToStr(ResultCode), mbError, MB_OK);

    // 移动 vcpkg.exe 到 bin
    if not FileCopy(ExpandConstant(vcpkgDir + '\vcpkg.exe'), ExpandConstant('{app}\bin\vcpkg.exe'), False) then
    begin
      MsgBox('无法复制 vcpkg.exe', mbError, MB_OK);
    end;
    
    if not CopyFolderOrFile(ExpandConstant(vcpkgDir), ExpandConstant('{app}')) then
    begin
      MsgBox('无法复制 vcpkg', mbError, MB_OK);
    end;

    // 清理
    DeleteFolderOrFile(vcpkgDir);
    DeleteFolderOrFile(ExpandConstant('{app}\vcpkg.exe'));
    WizardForm.ProgressGauge.Position := 100;
    WizardForm.StatusLabel.Caption := '安装完成';
  end;
end;


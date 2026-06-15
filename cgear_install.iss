#define MyAppName "cgear"
#define MyAppVersion "2.1.0"
#define MyAppExeName "cgear.exe"
#define NinjaVersion "1.12.1"
#define VcpkgRef "2026.06.01"

[Setup]
AppId={{8CB91DF4-334C-4A9A-89AF-6430EF4D9E30}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
DefaultDirName={localappdata}\Programs\{#MyAppName}
DefaultGroupName={#MyAppName}
DisableProgramGroupPage=yes
PrivilegesRequired=lowest
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
OutputDir=bin
OutputBaseFilename=cgear-v{#MyAppVersion}-x86_64-pc-windows-setup
Compression=lzma2
SolidCompression=yes
WizardStyle=modern
ChangesEnvironment=yes
CloseApplications=yes
RestartApplications=no
SetupLogging=yes
UninstallDisplayIcon={app}\bin\{#MyAppExeName}

[Dirs]
Name: "{app}\bin"
Name: "{app}\installed"
Name: "{app}\pkg"

[Files]
Source: "bin\{#MyAppExeName}"; DestDir: "{app}\bin"; Flags: ignoreversion
Source: "README.md"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\Cgear Command Prompt"; Filename: "{cmd}"; Parameters: "/K cd /D ""{app}"" && set ""CGEAR_HOME={app}"" && set ""PATH={app}\bin;%PATH%"""; WorkingDir: "{app}"

[UninstallDelete]
Type: files; Name: "{app}\bin\ninja.exe"
Type: files; Name: "{app}\bin\vcpkg.exe"
Type: files; Name: "{app}\vcpkg.exe"
Type: files; Name: "{app}\.vcpkg-root"
Type: filesandordirs; Name: "{app}\ports"
Type: filesandordirs; Name: "{app}\scripts"
Type: filesandordirs; Name: "{app}\triplets"
Type: filesandordirs; Name: "{app}\versions"

[Code]
const
  EnvironmentKey = 'Environment';
  BackupKey = 'Software\cgear\EnvironmentBackup';

function SamePath(const Left, Right: string): Boolean;
var
  A, B: string;
begin
  A := RemoveBackslashUnlessRoot(Trim(Left));
  B := RemoveBackslashUnlessRoot(Trim(Right));
  Result := CompareText(A, B) = 0;
end;

function PathContains(const PathValue, Entry: string): Boolean;
var
  Remaining, Item: string;
  Separator: Integer;
begin
  Result := False;
  Remaining := PathValue;
  while Remaining <> '' do
  begin
    Separator := Pos(';', Remaining);
    if Separator = 0 then
    begin
      Item := Remaining;
      Remaining := '';
    end
    else
    begin
      Item := Copy(Remaining, 1, Separator - 1);
      Delete(Remaining, 1, Separator);
    end;

    if SamePath(Item, Entry) then
    begin
      Result := True;
      Exit;
    end;
  end;
end;

function RemovePathEntry(const PathValue, Entry: string): string;
var
  Remaining, Item: string;
  Separator: Integer;
begin
  Result := '';
  Remaining := PathValue;
  while Remaining <> '' do
  begin
    Separator := Pos(';', Remaining);
    if Separator = 0 then
    begin
      Item := Remaining;
      Remaining := '';
    end
    else
    begin
      Item := Copy(Remaining, 1, Separator - 1);
      Delete(Remaining, 1, Separator);
    end;

    Item := Trim(Item);
    if (Item <> '') and not SamePath(Item, Entry) then
    begin
      if Result <> '' then
        Result := Result + ';';
      Result := Result + Item;
    end;
  end;
end;

procedure BackupEnvironment;
var
  PreviousValue: string;
begin
  if RegValueExists(HKCU, BackupKey, 'BackupCreated') then
    Exit;

  RegWriteDWordValue(HKCU, BackupKey, 'BackupCreated', 1);
  if RegQueryStringValue(HKCU, EnvironmentKey, 'CGEAR_HOME', PreviousValue) then
  begin
    RegWriteDWordValue(HKCU, BackupKey, 'PreviousCgearHomeExists', 1);
    RegWriteStringValue(HKCU, BackupKey, 'PreviousCgearHome', PreviousValue);
  end
  else
    RegWriteDWordValue(HKCU, BackupKey, 'PreviousCgearHomeExists', 0);
end;

procedure ConfigureEnvironment;
var
  PathValue, BinPath: string;
begin
  BackupEnvironment;
  BinPath := ExpandConstant('{app}\bin');

  if not RegQueryStringValue(HKCU, EnvironmentKey, 'Path', PathValue) then
    PathValue := '';
  if not PathContains(PathValue, BinPath) then
  begin
    if (PathValue <> '') and (PathValue[Length(PathValue)] <> ';') then
      PathValue := PathValue + ';';
    RegWriteExpandStringValue(HKCU, EnvironmentKey, 'Path', PathValue + BinPath);
  end;

  RegWriteStringValue(HKCU, EnvironmentKey, 'CGEAR_HOME', ExpandConstant('{app}'));
end;

procedure RestoreEnvironment;
var
  PathValue, CurrentHome, PreviousHome, BinPath: string;
  PreviousExists: Cardinal;
begin
  BinPath := ExpandConstant('{app}\bin');
  if RegQueryStringValue(HKCU, EnvironmentKey, 'Path', PathValue) then
    RegWriteExpandStringValue(HKCU, EnvironmentKey, 'Path', RemovePathEntry(PathValue, BinPath));

  if RegQueryStringValue(HKCU, EnvironmentKey, 'CGEAR_HOME', CurrentHome) and
     SamePath(CurrentHome, ExpandConstant('{app}')) then
  begin
    if RegQueryDWordValue(HKCU, BackupKey, 'PreviousCgearHomeExists', PreviousExists) and
       (PreviousExists = 1) and
       RegQueryStringValue(HKCU, BackupKey, 'PreviousCgearHome', PreviousHome) then
      RegWriteStringValue(HKCU, EnvironmentKey, 'CGEAR_HOME', PreviousHome)
    else
      RegDeleteValue(HKCU, EnvironmentKey, 'CGEAR_HOME');
  end;

  RegDeleteKeyIncludingSubkeys(HKCU, BackupKey);
end;

function CheckGitInstalled: Boolean;
var
  ResultCode: Integer;
begin
  Result := Exec('git.exe', '--version', '', SW_HIDE, ewWaitUntilTerminated, ResultCode) and
    (ResultCode = 0);
end;

function DependenciesInstalled: Boolean;
begin
  Result := FileExists(ExpandConstant('{app}\bin\ninja.exe')) and
    FileExists(ExpandConstant('{app}\scripts\buildsystems\vcpkg.cmake')) and
    FileExists(ExpandConstant('{app}\vcpkg.exe'));
end;

function RunChecked(const FileName, Parameters, WorkingDir, FailureMessage: string): Boolean;
var
  ResultCode: Integer;
begin
  ResultCode := -1;
  Result := Exec(FileName, Parameters, WorkingDir, SW_HIDE, ewWaitUntilTerminated, ResultCode);
  if not Result or (ResultCode <> 0) then
  begin
    MsgBox(FailureMessage + ' (exit code ' + IntToStr(ResultCode) + ')', mbError, MB_OK);
    Result := False;
  end;
end;

function CopyDirectory(const SourceDir, DestDir: string): Boolean;
begin
  ForceDirectories(DestDir);
  Result := RunChecked(ExpandConstant('{cmd}'),
    '/D /C xcopy "' + SourceDir + '\*" "' + DestDir + '" /E /I /Y /Q /H',
    '', 'Unable to copy ' + SourceDir);
end;

function DownloadFile(const URL, DestFile: string): Boolean;
var
  PowerShell, Parameters: string;
begin
  PowerShell := ExpandConstant('{sys}\WindowsPowerShell\v1.0\powershell.exe');
  Parameters := '-NoProfile -ExecutionPolicy Bypass -Command "' +
    '$ProgressPreference = ''SilentlyContinue''; ' +
    'Invoke-WebRequest -UseBasicParsing -Uri ''' + URL + ''' -OutFile ''' + DestFile + '''"';
  Result := RunChecked(PowerShell, Parameters, '', 'Download failed: ' + URL);
end;

function ExtractZip(const ZipFile, DestDir: string): Boolean;
var
  PowerShell, Parameters: string;
begin
  PowerShell := ExpandConstant('{sys}\WindowsPowerShell\v1.0\powershell.exe');
  Parameters := '-NoProfile -ExecutionPolicy Bypass -Command "' +
    'Expand-Archive -Path ''' + ZipFile + ''' -DestinationPath ''' + DestDir + ''' -Force"';
  Result := RunChecked(PowerShell, Parameters, '', 'Unable to extract ' + ZipFile);
end;

function InstallDependencies: Boolean;
var
  NinjaZip, VcpkgDir: string;
begin
  Result := False;
  ForceDirectories(ExpandConstant('{app}\bin'));

  WizardForm.StatusLabel.Caption := 'Downloading Ninja...';
  NinjaZip := ExpandConstant('{tmp}\ninja.zip');
  if not DownloadFile(
    'https://github.com/ninja-build/ninja/releases/download/v{#NinjaVersion}/ninja-win.zip',
    NinjaZip) then
    Exit;
  if not ExtractZip(NinjaZip, ExpandConstant('{app}\bin')) then
    Exit;

  WizardForm.StatusLabel.Caption := 'Downloading vcpkg...';
  VcpkgDir := ExpandConstant('{tmp}\cgear-vcpkg');
  DelTree(VcpkgDir, True, True, True);
  if not RunChecked('git.exe',
    'clone --depth 1 --branch {#VcpkgRef} https://github.com/microsoft/vcpkg.git "' + VcpkgDir + '"',
    '', 'Unable to clone vcpkg') then
    Exit;

  WizardForm.StatusLabel.Caption := 'Bootstrapping vcpkg...';
  if not RunChecked(ExpandConstant('{cmd}'),
    '/D /C ""' + VcpkgDir + '\bootstrap-vcpkg.bat" -disableMetrics"',
    VcpkgDir, 'Unable to bootstrap vcpkg') then
    Exit;

  if not CopyDirectory(VcpkgDir + '\ports', ExpandConstant('{app}\ports')) then
    Exit;
  if not CopyDirectory(VcpkgDir + '\scripts', ExpandConstant('{app}\scripts')) then
    Exit;
  if not CopyDirectory(VcpkgDir + '\triplets', ExpandConstant('{app}\triplets')) then
    Exit;
  if not CopyDirectory(VcpkgDir + '\versions', ExpandConstant('{app}\versions')) then
    Exit;

  if not FileCopy(VcpkgDir + '\.vcpkg-root', ExpandConstant('{app}\.vcpkg-root'), False) or
     not FileCopy(VcpkgDir + '\vcpkg.exe', ExpandConstant('{app}\vcpkg.exe'), False) or
     not FileCopy(VcpkgDir + '\vcpkg.exe', ExpandConstant('{app}\bin\vcpkg.exe'), False) then
  begin
    MsgBox('Unable to copy the vcpkg executable or root marker.', mbError, MB_OK);
    Exit;
  end;

  Result := True;
end;

function PrepareToInstall(var NeedsRestart: Boolean): String;
begin
  Result := '';
  if not DependenciesInstalled and not CheckGitInstalled then
    Result := 'Git is required to download vcpkg. Install Git and run setup again.';
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if (CurStep = ssInstall) and not DependenciesInstalled then
  begin
    if not InstallDependencies then
      RaiseException('Cgear dependencies could not be installed.');
  end;

  if CurStep = ssPostInstall then
    ConfigureEnvironment;
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usUninstall then
    RestoreEnvironment;
end;

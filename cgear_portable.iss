[Setup]
AppName=cgear
AppVersion=2.0.2
DefaultDirName={pf}\cgear
DefaultGroupName=cgear
OutputDir=bin
OutputBaseFilename=cgear-v2.0.2-x86_64-pc-windows-setup
Compression=lzma
SolidCompression=yes
AlwaysRestart=yes
; PrivilegesRequired=admin

[Files]
; 把整个便携目录（例如 portable 文件夹）下的所有内容打包
Source: "bin\portable\*"; DestDir: "{app}"; Flags: recursesubdirs createallsubdirs

[Run]
Filename: "{app}\bin\cgear.exe"; Description: "Run cgear"; Flags: nowait postinstall skipifsilent

[Registry]
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}\bin";
Root: HKCU; Subkey: "Environment"; ValueType: expandsz; ValueName: "CGEAR_HOME"; ValueData: "{app}"; Flags: createvalueifdoesntexist uninsdeletevalue
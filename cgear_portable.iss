#define MyAppName "cgear"
#define MyAppVersion "2.1.0"

[Setup]
AppName={#MyAppName} Portable
AppVersion={#MyAppVersion}
DefaultDirName={userdesktop}\cgear-portable
PrivilegesRequired=lowest
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
OutputDir=bin
OutputBaseFilename=cgear-v{#MyAppVersion}-x86_64-pc-windows-portable
Compression=lzma2
SolidCompression=yes
WizardStyle=modern
DisableProgramGroupPage=yes
Uninstallable=no
CreateUninstallRegKey=no

[Files]
Source: "bin\portable\*"; DestDir: "{app}"; Flags: ignoreversion recursesubdirs createallsubdirs
Source: "README.md"; DestDir: "{app}"; Flags: ignoreversion

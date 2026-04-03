; NSIS 3+ — adjust Source paths before building.
!define APPNAME "moonBASIC"
!define EXE "moonbasic.exe"
OutFile "moonbasic-setup.exe"
InstallDir "$PROGRAMFILES64\${APPNAME}"
Name "${APPNAME}"
RequestExecutionLevel admin

Section ""
  SetOutPath $INSTDIR
  File "${EXE}"
  ; File /nonfatal *.dll
  SetOutPath "$INSTDIR\assets\fonts"
  File /r "..\..\assets\fonts\*.*"
  SetOutPath "$INSTDIR\examples"
  File /r "..\..\examples\*.*"
  WriteUninstaller "$INSTDIR\uninstall.exe"
  CreateShortcut "$SMPROGRAMS\${APPNAME}.lnk" "$INSTDIR\${EXE}"
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\uninstall.exe"
  RMDir /r "$INSTDIR"
SectionEnd

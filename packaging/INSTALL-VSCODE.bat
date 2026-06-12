@echo off
setlocal
cd /d "%~dp0"
if exist moonbasic.exe (
  moonbasic.exe install-vscode
) else (
  echo Could not find moonbasic.exe next to this script.
  echo Extract the full moonBASIC release zip first, then double-click INSTALL-VSCODE.bat again.
  pause
  exit /b 1
)
echo.
echo Open VS Code or Cursor and open any .mb file.
pause

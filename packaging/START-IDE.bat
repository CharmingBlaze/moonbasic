@echo off
cd /d "%~dp0"
if not exist "moonbasic-ide.exe" (
  echo moonbasic-ide.exe not found in this folder.
  pause
  exit /b 1
)
start "" "moonbasic-ide.exe"

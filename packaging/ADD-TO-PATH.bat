@echo off
REM Add this folder to your user PATH so "moonbasic" and "moonrun" work in any terminal.
set "HERE=%~dp0"
set "HERE=%HERE:~0,-1%"

echo.
echo moonBASIC — add this folder to your user PATH:
echo   %HERE%
echo.
echo This lets you run:  moonbasic --check game.mb
echo                     moonrun game.mb
echo from any folder. The IDE does NOT need PATH (it finds binaries beside itself).
echo.

set "KEY=HKCU\Environment"
for /f "tokens=2*" %%A in ('reg query "%KEY%" /v Path 2^>nul') do set "CUR=%%B"
echo %CUR% | find /I "%HERE%" >nul
if not errorlevel 1 (
  echo Already on PATH.
  pause
  exit /b 0
)

if defined CUR (
  setx Path "%CUR%;%HERE%" >nul
) else (
  setx Path "%HERE%" >nul
)
if errorlevel 1 (
  echo Failed to update PATH. Add the folder manually in System Properties → Environment Variables.
  pause
  exit /b 1
)
echo Done. Open a NEW terminal for the change to take effect.
pause

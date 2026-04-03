# Distribution artifacts

- **windows/** — NSIS script `moonbasic.nsi` builds an installer. Run on Windows with NSIS 3.x after placing `moonbasic.exe` and required MinGW DLLs next to the script (see CI release job).
- **linux/** — `build-appimage.sh` and `build-deb.sh` expect a staged tree under `dist/stage/` with `bin/moonbasic`, `share/moonbasic/{examples,assets}`.

Release CI (`.github/workflows/release.yml`) produces portable zips per OS. Full NSIS/AppImage integration may require local paths adjusted for your Raylib/MinGW layout.

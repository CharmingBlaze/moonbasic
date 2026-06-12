# moonBASIC (VS Code extension)

Syntax highlighting, snippets, **LSP**, **check / compile / run**, and a **debugger** for moonBASIC (`.mb` / `.mbc`).

## Install (easiest)

### If you have moonBASIC from a release zip

```bash
moonbasic install-vscode
```

Or:

- **Windows:** double-click **`INSTALL-VSCODE.bat`** in the extracted folder
- **Linux / macOS:** **`./INSTALL-VSCODE.sh`**

The release zip includes the `.vsix` and the installer configures **`moonbasic`** / **`moonrun`** paths automatically.

### If you only have VS Code / Cursor

1. Download **`moonbasic-*-vscode.vsix`** from [GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest), **or** run **`moonbasic install-vscode`** (downloads the VSIX if needed).
2. The installer uses **`code --install-extension`** or **`cursor --install-extension`**.

### From this repository (developers)

```bash
./scripts/install-vscode-extension.sh
# or: powershell -File scripts/install-vscode-extension.ps1
```

## After install

Open any **`.mb`** file:

| Action | Shortcut |
|--------|----------|
| Run | **Ctrl+F5** |
| Check | **Ctrl+Shift+C** |
| Help at cursor | **Alt+H** |
| Debug | **Run and Debug** → **Debug moonBASIC** |

## Settings (usually automatic)

| Setting | Description |
|--------|-------------|
| **`moonbasic.languageServerPath`** | Path to `moonbasic` (set by `install-vscode`) |
| **`moonbasic.moonrunPath`** | Path to `moonrun` (set by `install-vscode`) |
| **`moonbasic.checkOnSave`** | Run `moonbasic --check` on save (default: false) |

## Build VSIX locally

```bash
cd editors/vscode-moonbasic
npm install
npm run compile
npm run package
```

User guide: **[docs/GETTING_STARTED.md](../../docs/GETTING_STARTED.md#vs-code-syntax-and-lsp)**.

import * as cp from "child_process";
import * as path from "path";
import * as vscode from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient | undefined;
const cliDiagnostics = vscode.languages.createDiagnosticCollection("moonbasic-cli");

const MOON_ERROR_RE =
  /^\[moonBASIC\].+ in (.+) line (\d+) col (\d+):$/;

function moonbasicPath(): string {
  const configured = vscode.workspace
    .getConfiguration("moonbasic")
    .get<string>("languageServerPath", "")
    .trim();
  return configured.length > 0 ? configured : "moonbasic";
}

function moonrunPath(): string {
  const configured = vscode.workspace
    .getConfiguration("moonbasic")
    .get<string>("moonrunPath", "")
    .trim();
  return configured.length > 0 ? configured : "moonrun";
}

function workspaceFolderFor(uri: vscode.Uri): string {
  const folder = vscode.workspace.getWorkspaceFolder(uri);
  return folder?.uri.fsPath ?? path.dirname(uri.fsPath);
}

function activeMoonbasicDocument(): vscode.TextDocument | undefined {
  const editor = vscode.window.activeTextEditor;
  if (!editor || editor.document.languageId !== "moonbasic") {
    return undefined;
  }
  return editor.document;
}

function runExecutable(
  command: string,
  args: string[],
  cwd: string
): Promise<{ stdout: string; stderr: string; code: number }> {
  return new Promise((resolve) => {
    cp.execFile(command, args, { cwd, maxBuffer: 8 * 1024 * 1024 }, (err, stdout, stderr) => {
      const code =
        err && typeof (err as cp.ExecException).code === "number"
          ? (err as cp.ExecException).code!
          : 0;
      resolve({ stdout, stderr, code });
    });
  });
}

interface MoonDiagnostic {
  uri: vscode.Uri;
  diagnostic: vscode.Diagnostic;
}

function parseMoonErrors(text: string, workspaceFolder: string): MoonDiagnostic[] {
  const lines = text.split(/\r?\n/);
  const results: MoonDiagnostic[] = [];
  for (let i = 0; i < lines.length; i++) {
    const match = MOON_ERROR_RE.exec(lines[i]);
    if (!match) {
      continue;
    }
    const filePath = match[1];
    const line = parseInt(match[2], 10);
    const col = parseInt(match[3], 10);
    let message = "moonBASIC error";
    if (i + 1 < lines.length && lines[i + 1].trim().length > 0) {
      message = lines[i + 1].trim();
    }
    for (let j = i + 1; j < lines.length; j++) {
      if (MOON_ERROR_RE.test(lines[j])) {
        break;
      }
      if (lines[j].startsWith("  Hint: ")) {
        message += " — " + lines[j].replace(/^  Hint: /, "");
        break;
      }
    }
    const resolved = path.isAbsolute(filePath)
      ? filePath
      : path.join(workspaceFolder, filePath);
    const uri = vscode.Uri.file(resolved);
    const range = new vscode.Range(
      Math.max(0, line - 1),
      Math.max(0, col - 1),
      Math.max(0, line - 1),
      Math.max(0, col - 1) + 1
    );
    results.push({
      uri,
      diagnostic: new vscode.Diagnostic(range, message, vscode.DiagnosticSeverity.Error),
    });
  }
  return results;
}

function applyCliDiagnostics(
  doc: vscode.TextDocument,
  items: MoonDiagnostic[]
): void {
  const byUri = new Map<string, vscode.Diagnostic[]>();
  for (const item of items) {
    const key = item.uri.toString();
    const list = byUri.get(key) ?? [];
    list.push(item.diagnostic);
    byUri.set(key, list);
  }
  cliDiagnostics.clear();
  for (const [key, diags] of byUri) {
    cliDiagnostics.set(vscode.Uri.parse(key), diags);
  }
  if (items.length === 0) {
    cliDiagnostics.delete(doc.uri);
  }
}

function showCompilerOutput(text: string): void {
  const channel = vscode.window.createOutputChannel("moonBASIC");
  channel.appendLine(text);
  channel.show(true);
}

async function checkDocument(
  doc: vscode.TextDocument,
  showOkMessage = true
): Promise<boolean> {
  const cwd = workspaceFolderFor(doc.uri);
  const moonbasic = moonbasicPath();
  const { stderr, code } = await runExecutable(moonbasic, ["--check", doc.fileName], cwd);
  const output = stderr.trim();
  if (code === 0) {
    cliDiagnostics.delete(doc.uri);
    if (showOkMessage) {
      vscode.window.showInformationMessage("moonBASIC: Check OK");
    }
    return true;
  }
  const parsed = parseMoonErrors(output, cwd);
  if (parsed.length > 0) {
    applyCliDiagnostics(doc, parsed);
    vscode.window.showErrorMessage("moonBASIC: check failed");
  } else if (output) {
    showCompilerOutput(output);
    vscode.window.showErrorMessage("moonBASIC: check failed (see Output)");
  } else {
    vscode.window.showErrorMessage("moonBASIC: check failed");
  }
  return false;
}

async function compileDocument(doc: vscode.TextDocument): Promise<boolean> {
  const cwd = workspaceFolderFor(doc.uri);
  const moonbasic = moonbasicPath();
  const { stderr, code } = await runExecutable(moonbasic, [doc.fileName], cwd);
  if (code === 0) {
    const mbc = doc.fileName.replace(/\.mb$/i, ".mbc");
    vscode.window.showInformationMessage(`moonBASIC: wrote ${path.basename(mbc)}`);
    return true;
  }
  const output = stderr.trim();
  const parsed = parseMoonErrors(output, cwd);
  if (parsed.length > 0) {
    applyCliDiagnostics(doc, parsed);
  } else if (output) {
    showCompilerOutput(output);
  }
  vscode.window.showErrorMessage("moonBASIC: compile failed");
  return false;
}

function runDocument(doc: vscode.TextDocument): void {
  const cwd = workspaceFolderFor(doc.uri);
  const moonrun = moonrunPath();
  const term = vscode.window.terminals.find((t) => t.name === "moonBASIC Run");
  const terminal =
    term ??
    vscode.window.createTerminal({
      name: "moonBASIC Run",
      cwd,
    });
  terminal.show();
  const quoted = doc.fileName.includes(" ") ? `"${doc.fileName}"` : doc.fileName;
  terminal.sendText(`${moonrun} ${quoted}`, true);
}

async function showHelpAtCursor(): Promise<void> {
  await vscode.commands.executeCommand("editor.action.showHover");
}

async function discoverSiblingBinaries(): Promise<{
  moonbasic?: string;
  moonrun?: string;
}> {
  const moonbasicName =
    process.platform === "win32" ? "moonbasic.exe" : "moonbasic";
  const moonrunName =
    process.platform === "win32" ? "moonrun.exe" : "moonrun";
  const searchRoots: string[] = [];
  const folders = vscode.workspace.workspaceFolders;
  if (folders) {
    for (const f of folders) {
      searchRoots.push(f.uri.fsPath);
    }
  }
  for (const root of searchRoots) {
    let cur = root;
    for (let depth = 0; depth < 5; depth++) {
      const mb = path.join(cur, moonbasicName);
      const mr = path.join(cur, moonrunName);
      try {
        await vscode.workspace.fs.stat(vscode.Uri.file(mb));
        const out: { moonbasic?: string; moonrun?: string } = {
          moonbasic: mb,
        };
        try {
          await vscode.workspace.fs.stat(vscode.Uri.file(mr));
          out.moonrun = mr;
        } catch {
          // moonrun optional for LSP-only workflows
        }
        return out;
      } catch {
        // not in this folder
      }
      const parent = path.dirname(cur);
      if (parent === cur) {
        break;
      }
      cur = parent;
    }
  }
  return {};
}

async function ensureToolPaths(): Promise<void> {
  const cfg = vscode.workspace.getConfiguration("moonbasic");
  const mb = cfg.get<string>("languageServerPath", "").trim();
  const mr = cfg.get<string>("moonrunPath", "").trim();
  if (mb && mr) {
    return;
  }
  const found = await discoverSiblingBinaries();
  if (!mb && found.moonbasic) {
    await cfg.update(
      "languageServerPath",
      found.moonbasic,
      vscode.ConfigurationTarget.Global
    );
  }
  if (!mr && found.moonrun) {
    await cfg.update(
      "moonrunPath",
      found.moonrun,
      vscode.ConfigurationTarget.Global
    );
  }
}

async function openDocs(): Promise<void> {
  const folders = vscode.workspace.workspaceFolders;
  if (folders && folders.length > 0) {
    const beginHere = vscode.Uri.joinPath(folders[0].uri, "docs", "BEGIN_HERE.md");
    try {
      await vscode.workspace.fs.readFile(beginHere);
      const doc = await vscode.workspace.openTextDocument(beginHere);
      await vscode.window.showTextDocument(doc);
      return;
    } catch {
      // fall through
    }
    const guides = vscode.Uri.joinPath(folders[0].uri, "docs", "systems", "GUIDES.md");
    try {
      await vscode.workspace.fs.readFile(guides);
      const doc = await vscode.workspace.openTextDocument(guides);
      await vscode.window.showTextDocument(doc);
      return;
    } catch {
      // fall through
    }
  }
  await vscode.env.openExternal(
    vscode.Uri.parse(
      "https://github.com/CharmingBlaze/moonbasic-compiler/blob/main/docs/BEGIN_HERE.md"
    )
  );
}

export async function activate(context: vscode.ExtensionContext): Promise<void> {
  await ensureToolPaths();

  const config = vscode.workspace.getConfiguration("moonbasic");
  const configured = config.get<string>("languageServerPath", "").trim();
  const command = configured.length > 0 ? configured : "moonbasic";

  const serverOptions: ServerOptions = {
    run: { command, args: ["--lsp"], transport: TransportKind.stdio },
    debug: { command, args: ["--lsp"], transport: TransportKind.stdio },
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "moonbasic" },
    ],
  };

  client = new LanguageClient(
    "moonbasic",
    "moonBASIC Language Server",
    serverOptions,
    clientOptions
  );

  context.subscriptions.push(cliDiagnostics);
  context.subscriptions.push({
    dispose: () => {
      void client?.stop();
    },
  });
  void client.start();

  context.subscriptions.push(
    vscode.debug.registerDebugAdapterDescriptorFactory("moonbasic", {
      createDebugAdapterDescriptor(session: vscode.DebugSession) {
        const cfg = session.configuration as { moonrunPath?: string };
        const wsCfg = vscode.workspace.getConfiguration("moonbasic");
        const runConfigured =
          (cfg.moonrunPath ?? wsCfg.get<string>("moonrunPath", "")).trim();
        const runCommand = runConfigured.length > 0 ? runConfigured : "moonrun";
        return new vscode.DebugAdapterExecutable(runCommand, ["--dap"]);
      },
    })
  );

  context.subscriptions.push(
    vscode.commands.registerCommand("moonbasic.check", async () => {
      const doc = activeMoonbasicDocument();
      if (!doc) {
        vscode.window.showWarningMessage("Open a .mb file to check.");
        return;
      }
      await checkDocument(doc);
    }),
    vscode.commands.registerCommand("moonbasic.compile", async () => {
      const doc = activeMoonbasicDocument();
      if (!doc) {
        vscode.window.showWarningMessage("Open a .mb file to compile.");
        return;
      }
      await compileDocument(doc);
    }),
    vscode.commands.registerCommand("moonbasic.run", () => {
      const doc = activeMoonbasicDocument();
      if (!doc) {
        vscode.window.showWarningMessage("Open a .mb file to run.");
        return;
      }
      runDocument(doc);
    }),
    vscode.commands.registerCommand("moonbasic.help", () => showHelpAtCursor()),
    vscode.commands.registerCommand("moonbasic.openDocs", () => openDocs())
  );

  if (config.get<boolean>("checkOnSave", false)) {
    context.subscriptions.push(
      vscode.workspace.onDidSaveTextDocument((doc) => {
        if (doc.languageId === "moonbasic") {
          void checkDocument(doc, false);
        }
      })
    );
  }

  if (!context.globalState.get<boolean>("moonbasic.welcomed")) {
    void context.globalState.update("moonbasic.welcomed", true);
    void vscode.window
      .showInformationMessage(
        "moonBASIC is ready — Ctrl+F5 run, Ctrl+Shift+C check, Alt+H help at cursor.",
        "Open docs"
      )
      .then((choice) => {
        if (choice === "Open docs") {
          void openDocs();
        }
      });
  }
}

export function deactivate(): Thenable<void> | undefined {
  return client?.stop();
}

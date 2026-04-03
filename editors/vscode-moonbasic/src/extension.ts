import * as vscode from "vscode";
import {
  LanguageClient,
  LanguageClientOptions,
  ServerOptions,
  TransportKind,
} from "vscode-languageclient/node";

let client: LanguageClient | undefined;

export function activate(context: vscode.ExtensionContext): void {
  const config = vscode.workspace.getConfiguration("moonbasic");
  const configured = config.get<string>("languageServerPath", "").trim();
  const command = configured.length > 0 ? configured : "moonbasic";

  const serverOptions: ServerOptions = {
    run: { command, args: ["--lsp"], transport: TransportKind.stdio },
    debug: { command, args: ["--lsp"], transport: TransportKind.stdio },
  };

  const clientOptions: LanguageClientOptions = {
    documentSelector: [{ scheme: "file", language: "moonbasic" }],
  };

  client = new LanguageClient(
    "moonbasic",
    "moonBASIC Language Server",
    serverOptions,
    clientOptions
  );

  context.subscriptions.push({
    dispose: () => {
      void client?.stop();
    },
  });
  void client.start();
}

export function deactivate(): Thenable<void> | undefined {
  return client?.stop();
}

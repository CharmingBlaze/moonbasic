# Documentation Style Guide

All moonBASIC command reference documents should follow this standardized structure to ensure clarity and consistency.

## Header
Use a Level 1 heading with the module name.
`# [Namespace] Commands`

## Description
A brief one-sentence description of the module's purpose.

## Core Workflow (Optional)
If the module requires a specific sequence of calls (e.g., Init -> Loop -> Close), include a numbered list and a concise code example.

## Command Blocks
Group commands by sub-topic using Level 2 headings.

### Signatures
- Use Level 3 headings for each command signature.
- **PascalCase** naming: `Module.Command()`.
- **Parentheses**: Always include `()` even if there are no arguments.
- **Parameters**: List parameter names clearly within the parentheses.

### Explanations
- Provide a concise description of what the command does.
- Use bullet points for parameter descriptions if there are multiple.
- Include short code snippets for usage where helpful.

## Visual Separators
Use horizontal rules `---` between major logical groups of commands.

## Example Document Structure

```markdown
# Module Commands

Commands for managing [feature].

## Core Workflow

1. **Initialize**: Call `Module.Init()`.
2. **Main Loop**: Use `Module.Update()`.
3. **Cleanup**: Call `Module.Close()`.

```basic
Module.Init()
WHILE NOT Window.ShouldClose()
    Module.Update()
WEND
Module.Close()
\```

---

## Management

### `Module.Init(param1, param2)`
Initializes the system.
- `param1`: Description.
- `param2`: Description.

### `Module.Close()`
Releases all resources.

---

## Operations

### `Module.DoAction(id, amount)`
Performs a specific action.
```

---

## Final Consistency Check
Always verify signatures against `compiler/builtinmanifest/commands.json` before finalizing documentation.

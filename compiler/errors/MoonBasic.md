ROLE: You are a senior compiler engineer with 20 years of experience
designing production language runtimes. You have shipped compilers for
game scripting languages at commercial studios. You think in terms of
correctness first, performance second, and developer experience third.
You do not cut corners. You do not write placeholder code. You write
the real thing.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PROJECT: moonBASIC Compiler + Runtime
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
You are designing and implementing moonBASIC — a modern 2D and 3D game
programming language. The spiritual successor to BlitzBasic 3D and
DarkBASIC Professional. Built for game developers who want power
without ceremony.

Implementation language: Go 1.22+
Rendering:     raylib-go  github.com/gen2brain/raylib-go
2D Physics:    box2d      github.com/ByteArena/box2d
3D Physics:    jolt-go    github.com/bbitechnologies/jolt-go
Networking:    go-enet    github.com/codecat/go-enet
Target OS:     Windows x64 + Linux x64

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
GROUND TRUTH (READ FIRST)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
This file is the **long-form language + engine specification** for contributors.
For **enforced repo layout, stable APIs, and anti-revert rules**,
the mandatory companion is:

  **`ARCHITECTURE.md`** (repository root, next to `go.mod`)

If anything here disagrees with **`ARCHITECTURE.md`**, treat **`ARCHITECTURE.md`**
as authoritative for **what the codebase must do today**. This document may
describe **planned** modules (e.g. per-namespace `runtime/modules/*`) that are
not implemented yet — those sections are aspirational.

**Inside this file:** Language identity → Syntax → Reference program → **Pipeline
& repo tree** → File sizes → Go standards → Lexer → Parser → **VM bytecode**
→ Runtime dispatch → Errors → **CLI** → Implementation order → Cross-platform
→ Problem report format.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
LANGUAGE IDENTITY
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
moonBASIC feels like BlitzBasic 3D but cleaner, more powerful, and
built for modern hardware. The target user is a game developer who
wants to write game logic immediately without boilerplate, class
hierarchies, or build system configuration.

Core principles inherited from BlitzBasic 3D:
  - Everything is a handle — models, textures, bodies, sounds,
    peers are all integers under the hood
  - No boilerplate — first line of code is game code
  - Commands feel like English — readable without documentation
  - The game loop is obvious — setup, loop, render, done
  - Assets are trivial — one command to load anything

Core principles added for the modern era:
  - Namespace dot-call syntax: Physics3D.Start() Body3D.Make()
  - Mandatory parentheses on every call — no ambiguity ever
  - Named key constants: KEY_ESCAPE not integer literals
  - Built-in physics, networking, audio — no external plugins
  - Fully case agnostic — compiler normalises everything internally

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CASE AGNOSTICISM — THE FIRST LAW
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
moonBASIC is completely and unconditionally case agnostic. This means:

  - Every identifier, keyword, command name, namespace name, method
    name, variable name, function name, type name, field name, and
    constant name is treated as identical regardless of capitalisation
  - The lexer calls strings.ToUpper() on every identifier token at
    scan time before any other processing
  - The symbol table, dispatch table, keyword map, and alias table
    all store and look up keys in uppercase only
  - String literal CONTENTS are the only thing never uppercased
  - There is no "convention" the compiler enforces — the programmer
    can write in any style they prefer and the compiler accepts it

  All of these are identical and interchangeable:
    LOADMODEL("x")   loadmodel("x")   LoadModel("x")   lOaDmOdEl("x")
    PHYSICS3D.START()  physics3d.start()  Physics3D.Start()
    WHILE  while  While  wHiLe
    KEY_ESCAPE  key_escape  Key_Escape

  The normaliseCommand function handles namespace commands:
    "render.frame"    → "RENDER.FRAME"
    "Physics3D.Start" → "PHYSICS3D.START"
    "DRAW.TEXT"       → "DRAW.TEXT"
    "math.sin"        → "MATH.SIN"

  Implementation in the lexer — one line, no exceptions:
    ident = strings.ToUpper(rawIdent)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
SYNTAX SPECIFICATION — THE LAW
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

PARENTHESES — NON-NEGOTIABLE:
  Every command and function call requires parentheses. Always.
  Zero-argument commands still require parentheses. This rule has
  no exceptions anywhere in the language.

  CORRECT:  Window.Open(1280, 720, "My Game")
  CORRECT:  Render.Frame()
  CORRECT:  model = Model.Load("robot.mdl")
  CORRECT:  IF Input.KeyDown(KEY_W) THEN Camera.Move(0, 0, 0.1)
  CORRECT:  y# = Terrain.Height(Camera.X(), Camera.Z())
  CORRECT:  window.open(1280, 720, "My Game")
  CORRECT:  render.frame()
  WRONG:    Window.Open 1280, 720, "My Game"
  WRONG:    Render.Frame
  WRONG:    model = Model.Load "robot.mdl"

  When parens are missing the compiler emits this exact error:

    [moonBASIC] Error in {file} line {line} col {col}:
      Expected '(' after '{NAME}'

      {line_number} | {exact source line text}
                    | {spaces}^
      Hint: All commands require parentheses: {NAME}({args})

NAMESPACE DOT-CALL SYNTAX:
  Commands belong to namespaces separated by a dot.
  The entire command including namespace is uppercased by the lexer.
  The programmer can write it in any capitalisation they choose.

  Defined namespaces:
    WINDOW      RENDER      CAMERA      CAMERA2D
    MODEL       MESH        TEXTURE     SHADER
    LIGHT       TERRAIN     SPRITE      FONT
    DRAW        DRAW3D      INPUT       AUDIO
    PHYSICS3D   PHYSICS2D   BODY3D      BODY2D
    NET         PEER        MATH        VEC2
    VEC3        COLOR       FILE        MEM
    TIME        SYSTEM      DEBUG       SCENE
    TWEEN       NAV         JSON

  Global commands with no namespace:
    PRINT(v)   INPUT(prompt$)   DIM arr(n)   INCLUDE "file.mbc"
    STR$(v)    INT(v)           FLOAT(v)     LEN(s$)

VARIABLE SUFFIXES:
  score    = 0        ; INT    — plain, no suffix
  speed#   = 3.14     ; FLOAT  — # suffix
  name$    = "player" ; STRING — $ suffix
  alive?   = TRUE     ; BOOL   — ? suffix, TRUE/FALSE/NULL

  Suffixes participate fully in case agnosticism:
    y# and Y# are identical — stored internally as "Y#"
    name$ and NAME$ are identical — stored as "NAME$"
    alive? and ALIVE? are identical — stored as "ALIVE?"

  The suffix is consumed as part of the identifier token
  by the lexer before uppercasing.

  Implementation alignment (symbols):
  Treat **`NAME` and `NAME#` / `NAME$` / `NAME?` as distinct identifiers**
  (suffix is part of the unique string), matching classic Blitz-style typing.
  The lexer still uppercases the **whole** token (e.g. `y#` → `Y#`).

ARRAY ACCESS:
  Both () and [] are valid array index syntax.
  The lexer maps [ → TOK_LPAREN and ] → TOK_RPAREN.
  Parser and codegen see identical token streams for both.

  DIM enemies(100)
  enemies(0) = 42     ; correct
  enemies[0] = 42     ; also correct — lexer normalises

OPERATORS:
  Arithmetic:   +  -  *  /  MOD  ^
  Comparison:   =  <>  <  >  <=  >=
  Logical:      AND  OR  NOT  XOR
  Compound:     +=  -=  *=  /=
  String:       + (concatenation)
  Power:        ^ (right associative)

  NOT vs OR / AND (permanent — matches parser and ARCHITECTURE.md §7):
    NOT binds TIGHTER than OR, XOR, and AND. So:
      NOT a OR b     parses as   (NOT a) OR b
      NOT a AND b    parses as   (NOT a) AND b
    To negate a whole disjunction, parenthesize:
      NOT (a OR b)
    Game loop (ESC or OS close) MUST use:
      WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
    The unparenthesised OR form after NOT breaks window-close handling.

MULTIPLE STATEMENTS ON ONE LINE:
  Use colon as separator:
    x = 1 : y = 2 : z = 3
    IF dead? THEN score = 0 : GOTO respawn

CONTROL FLOW:

  ; Single-line IF — THEN required
  IF Input.KeyDown(KEY_W) THEN Camera.Move(0, 0, 0.1)

  ; Multi-line IF — THEN required on opening line
  IF hp < 0 THEN
      Audio.Play(deathSnd)
  ELSEIF hp < 20 THEN
      Audio.Play(lowHpSnd)
  ELSE
      Audio.Play(normalSnd)
  ENDIF

  ; Two-word block endings — all are identical to one-word form
  END IF       = ENDIF
  END FUNCTION = ENDFUNCTION
  END WHILE    = WEND        ; also ENDWHILE is valid
  END SELECT   = ENDSELECT
  END TYPE     = ENDTYPE

  ; WHILE / WEND — key-only
  WHILE NOT Input.KeyDown(KEY_ESCAPE)
      Render.Frame()
  WEND

  ; WHILE with OS close button — parentheses required (see NOT vs OR above)
  WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
      Render.Clear(20, 20, 30)
      Render.Frame()
  WEND

  ; FOR / NEXT with optional STEP
  FOR i = 1 TO 10
      PRINT(i)
  NEXT

  FOR i = 10 TO 1 STEP -1
      PRINT(i)
  NEXT

  ; REPEAT / UNTIL
  REPEAT
      Render.Frame()
  UNTIL Input.KeyDown(KEY_ESCAPE)

  ; SELECT / CASE
  SELECT state
      CASE 0
          RunMenu()
      CASE 1
          RunGame()
      DEFAULT
          QUIT()
  ENDSELECT

  ; GOTO / GOSUB — dot prefix on label definition
  GOTO mainloop
  GOSUB loadAssets
  RETURN()

  .mainloop
  .loadAssets

FUNCTIONS:
  FUNCTION SpawnEnemy(x#, y#, z#)
      e = Body3D.Make("dynamic")
      Body3D.AddCapsule(e, 0.4, 1.8)
      Body3D.SetPos(e, x#, y#, z#)
      RETURN(e)
  ENDFUNCTION

  ; Void function — no RETURN needed
  FUNCTION ResetGame()
      score = 0
      hp = 100
  ENDFUNCTION

  ; Call always uses parentheses
  enemy = SpawnEnemy(10.0, 0.0, 5.0)
  ResetGame()

  ; Case agnostic — all identical:
  enemy = SpawnEnemy(10.0, 0.0, 5.0)
  enemy = spawnenemy(10.0, 0.0, 5.0)
  enemy = SPAWNENEMY(10.0, 0.0, 5.0)

USER DEFINED TYPES:
  TYPE Enemy
      FIELD mesh
      FIELD body
      FIELD hp
      FIELD x#
      FIELD z#
  ENDTYPE

  e = NEW(Enemy)
  e.hp   = 100
  e.mesh = Model.Load("enemy.mdl")
  e.body = Body3D.Make("dynamic")

  FOR e = EACH(Enemy)
      IF e.hp <= 0 THEN
          Model.Free(e.mesh)
          Body3D.Free(e.body)
          DELETE(e)
      ENDIF
  NEXT

INCLUDES:
  INCLUDE "weapons.mbc"
  INCLUDE "enemy_ai.mbc"

COMMENTS:
  ; Semicolon starts a line comment unconditionally
  ; Everything after ; on the same line is ignored
  ; This is the only comment syntax — no block comments
  ; String scanner never calls the comment handler so
  ; semicolons inside strings are literal characters

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CANONICAL REFERENCE PROGRAM
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
This program is the acceptance test. Every feature of the compiler
must support this running correctly. It does not change.
The programmer chose mixed case — the compiler must accept it.

  Window.Open(1280, 720, "moonBASIC Demo")
  Window.SetFPS(60)

  Render.SetSkybox("assets/sky/sunset.ktx2")
  Render.SetIBLIntensity(1.0)
  Render.SetIBLSplit(0.6, 1.2)
  Render.SetShadowMapSize(2048)
  Render.SetAmbient(0.1, 0.1, 0.2)

  cam = Camera.Make()
  cam.SetPos(0, 5, -12)
  cam.SetRot(20, 0, 0)
  cam.SetFOV(75)

  sun = Light.Make("directional")
  sun.SetDir(-1, -2, -1)
  sun.SetIntensity(1.0)
  sun.SetShadow(TRUE)

  terrain = Terrain.Make(256, 256)
  terrain.FillPerlin(0.03, 20)
  terrain.SetLayer(0, Texture.Load("assets/grass.png"))
  terrain.SetLayer(1, Texture.Load("assets/rock.png"))
  terrain.SetSplat(Texture.Load("assets/splat.png"))
  terrain.SetLayerTile(4.0)

  robot = Model.Load("assets/robot.mdl")
  robot.SetTexture(Texture.Load("assets/robot.png"))

  Physics3D.Start()
  Physics3D.SetGravity(0, -9.8, 0)

  ground = Body3D.Make("static")
  ground.AddBox(100, 1, 100)
  ground.SetPos(0, -0.5, 0)

  player = Body3D.Make("dynamic")
  player.AddCapsule(0.4, 1.8)
  player.SetPos(0, 5, 0)
  player.SetMass(80)

  Physics3D.OnCollision(player, ground, "OnLand")

  mus = Audio.LoadMusic("assets/theme.ogg")
  mus.Play()

  ; For OS close as well as ESC: WHILE NOT (Input.KeyDown(KEY_ESCAPE) OR Window.ShouldClose())
  WHILE NOT Input.KeyDown(KEY_ESCAPE)

      dt# = Time.Delta()

      IF Input.KeyDown(KEY_W) THEN cam.Move(0, 0,  5.0 * dt#)
      IF Input.KeyDown(KEY_S) THEN cam.Move(0, 0, -5.0 * dt#)
      IF Input.KeyDown(KEY_A) THEN cam.Move(-5.0 * dt#, 0, 0)
      IF Input.KeyDown(KEY_D) THEN cam.Move( 5.0 * dt#, 0, 0)

      IF Input.KeyPressed(KEY_SPACE) THEN
          player.ApplyImpulse(0, 600, 0)
      ENDIF

      Physics3D.Step()
      mus.Update()

      robot.SetPos(player.X(), player.Y(), player.Z())
      robot.Rotate(0, 45 * dt#, 0)

      y# = terrain.Height(cam.X(), cam.Z())
      cam.SetY(y# + 1.8)

      Render.Clear(10, 10, 20)
      terrain.Draw()
      robot.Draw()
      Render.Frame()

  WEND

  robot.Free()
  terrain.Free()
  sun.Free()
  mus.Free()
  Physics3D.Stop()
  Window.Close()

  FUNCTION OnLand(a, b)
      Audio.Play(landSnd)
  ENDFUNCTION

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
COMPILER PIPELINE ARCHITECTURE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
**Source text** (e.g. `.mb`, `.mbc` extension is convention only for sources)
→ **Lexer** → **Parser** → **AST** → **Semantic** (fold + manifest checks)
→ **CodeGen** → **`opcode.Program`** (in-memory bytecode)
→ optional **`vm/moon` encode** → **`.mbc` file** (MOON container)

**Execute:** **`DecodeMOON`** (if from disk) → **VM** → **`runtime.Registry`**
(builtin dispatch + heap).

Dependency rules (target architecture):
  - **`compiler/*`** must not depend on raylib/jolt/enet CGo.
  - **`compiler/pipeline`** may import **`runtime`** and **`vm`** so embedders get
    **`RunProgram`** in one call (see **`ARCHITECTURE.md`** — this is intentional).
  - **`vm`** uses **`runtime`** only for **`Registry`** / **`BuiltinFn`** dispatch,
    not for engine modules that do not exist yet.

**Mandatory codegen rule:** `CallStmtNode` must lower to **`emitCallStmt`** (see
`codegen_stmts.go` + `codegen_calls.go`). An empty `case` for call statements is
a critical bug (`PRINT` would emit no bytecode).

Current repository layout (as implemented):

  moonbasic/
  ├── main.go
  ├── ARCHITECTURE.md          ← stable API + AI ground truth
  ├── Masterplan.md
  ├── go.mod
  ├── compiler/
  │   ├── token/
  │   ├── lexer/
  │   ├── ast/
  │   ├── parser/
  │   ├── symtable/
  │   ├── semantic/            ← analysis + strict dotted builtin checks
  │   ├── builtinmanifest/     ← commands.json (semantic + runtime stubs)
  │   ├── codegen/             ← codegen.go, codegen_expr/stmts/calls.go
  │   ├── pipeline/            ← Compile*, Check*, RunProgram, EncodeMOON, DecodeMOON
  │   └── errors/              ← errors.go, MoonBasic.md (this document)
  ├── vm/
  │   ├── opcode/
  │   ├── value/
  │   ├── heap/
  │   ├── callstack/
  │   ├── moon/                ← MOON .mbc binary format (ship bytecode)
  │   ├── vm.go
  │   ├── vm_control.go
  │   ├── vm_dispatch.go
  │   └── *_test.go
  └── runtime/
      ├── runtime.go           ← Registry, RegisterFromManifest (manifest stubs)
      └── core.go              ← PRINT, STR$, INT, FLOAT, LEN, …

**Planned (Phase B — not necessarily present yet):** split native engine into
`runtime/modules/{window,render,physics3d,...}` as in the original mega-tree.
Until then, new builtins are registered from **`core.go`** or small additions
that **`Bind`** into the same flat **`Commands`** map.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
FILE SIZE GUIDELINES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
These are targets not hard limits. Split a file when it becomes
difficult to navigate or when it mixes distinct concerns.
A file can be larger if it is genuinely cohesive.

  main.go                       target  100 lines
  compiler/token/token.go       target  200 lines
  compiler/lexer/*.go           target  300 lines each
  compiler/ast/ast.go           target  500 lines
  compiler/parser/*.go          target  400 lines each
  compiler/symtable/*.go        target  250 lines each
  compiler/codegen/*.go         target  400 lines each
  compiler/errors/*.go          target  200 lines
  vm/opcode/opcode.go           target  150 lines
  vm/value/value.go             target  150 lines
  vm/heap/*.go                  target  300 lines each
  vm/*.go (vm package)          target  400 lines each
  vm/moon/*.go                  target  400 lines each
  runtime/runtime.go            target  200 lines
  runtime/core.go               target  200 lines
  Any future runtime module     target  400 lines

  GOD FILE RULE: If a single file exceeds 800 lines and is not
  cohesive — meaning it mixes multiple distinct concerns — split it.
  A 600 line file that does one thing well is fine.
  A 300 line file that mixes parsing, codegen, and runtime is not.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
GO CODING STANDARDS
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

NAMING:
  Packages         lowercase single word        lexer, parser, mathmod
  Exported types   PascalCase                   Lexer, Parser, Value
  Exported fns     PascalCase                   Scan(), Parse(), Execute()
  Unexported       camelCase                    scanIdent(), emitOp()
  Constants        PascalCase or ALLCAPS        OpAdd, MaxStackDepth
  Test files       _test.go suffix              lexer_test.go

FUNCTIONS:
  Target 40 lines per function — split if logic warrants it
  Max 5 parameters — use a config struct if more needed
  Every exported function has a godoc comment
  Return early to avoid deep nesting
  No magic numbers — use named constants or explanatory comments

ERROR HANDLING:
  Never use -1 or 0 as error sentinel values — use the error interface
  Every compiler error includes file name + line number + column number
  Every runtime error includes the moonBASIC source location
  Use the MoonError type in errors.go for all structured errors
  Wrap errors with context: fmt.Errorf("parsing IF: %w", err)
  Never panic in runtime code — return error and let the VM handle it

CORE INTERFACES:

  // HeapObject is implemented by every object stored in the heap.
  // Free() releases the underlying external resource (raylib, jolt etc).
  type HeapObject interface {
      Free()
      TypeName() string
  }

  // BuiltinFn is the signature for every runtime command.
  // args are already type-checked by the VM before dispatch.
  type BuiltinFn func(args []Value) (Value, error)

  // Module is implemented by every runtime module (when split into packages).
  // Register fills the Registry's command map; Shutdown releases resources.
  type Module interface {
      Register(reg map[string]BuiltinFn)
      Shutdown()
  }

MEMORY AND RESOURCE MANAGEMENT:
  Every object created by a moonBASIC LOAD*/CREATE*/MAKE* command
  is registered in the Heap with an integer handle.
  The VM calls Heap.FreeAll() on shutdown — no leaks ever.
  raylib objects (Model, Texture, Sound) call their Unload function.
  jolt-go bodies and shapes call their Destroy function in Free().
  go-enet hosts and peers call their Destroy function in Free().
  Box2D worlds are destroyed on Physics2D.Stop().

CONCURRENCY:
  Use sync.RWMutex on all shared state in physics and network modules.
  Use sync.Once for one-time library initialisation (jolt, enet).
  The VM execution loop is single-threaded — no goroutines inside it.
  Physics3D.Step() and Net.Update() are called from the main loop.
  All raylib calls must happen on the main OS thread.
  Lock the main goroutine to the OS thread with runtime.LockOSThread()
  in main() before calling any raylib functions.

GO SPECIFICS:
  Use strings.ToUpper() in the lexer — not manual ASCII arithmetic.
  Use fmt.Errorf("context: %w", err) for error wrapping.
  No init() functions that do real work — use explicit Register().
  Prefer value semantics for small structs (Token, Value, Instruction).
  Use sync.Pool for frequently allocated temporary objects in hot paths.
  Profile before optimising — correctness first.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
LEXER SPECIFICATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

Token types required:
  TOK_IDENT         identifier or keyword — always uppercase after scan
  TOK_INTEGER       42  -42
  TOK_FLOAT         3.14  -3.14
  TOK_STRING        "hello world"
  TOK_NEWLINE       significant — ends a statement
  TOK_LPAREN        ( and [  — both map to this
  TOK_RPAREN        ) and ]  — both map to this
  TOK_COMMA         ,
  TOK_COLON         :
  TOK_DOT           .
  TOK_HASH          #  float suffix
  TOK_DOLLAR        $  string suffix
  TOK_QUESTION      ?  bool suffix
  TOK_ASSIGN        =  in assignment context
  TOK_PLUS          +
  TOK_MINUS         -
  TOK_STAR          *
  TOK_SLASH         /
  TOK_CARET         ^
  TOK_EQ            =  in expression context — parser disambiguates
  TOK_NEQ           <>
  TOK_LT            
  TOK_GT            >
  TOK_LTE           <=
  TOK_GTE           >=
  TOK_PLUS_EQ       +=
  TOK_MINUS_EQ      -=
  TOK_STAR_EQ       *=
  TOK_SLASH_EQ      /=
  TOK_EOF

Keywords — map[string]TokenType with all keys uppercase:
  IF THEN ELSE ELSEIF ENDIF
  WHILE WEND ENDWHILE
  FOR TO STEP NEXT
  REPEAT UNTIL
  SELECT CASE DEFAULT ENDSELECT
  FUNCTION ENDFUNCTION RETURN
  TYPE FIELD ENDTYPE
  NEW DELETE EACH
  GOTO GOSUB
  AND OR NOT XOR MOD
  DIM REDIM
  LOCAL GLOBAL CONST
  INCLUDE
  TRUE FALSE NULL
  END

Lexer rules — implement exactly in this order:

  1.  Track line (starts at 1) and column (starts at 1) for every byte.
  2.  Skip spaces and tabs. Do NOT skip newlines.
  3.  \n and \r\n both emit TOK_NEWLINE and increment the line counter.
      Reset column to 1 after each newline.
  4.  Semicolons begin a line comment. Skip all bytes until \n.
      The string scanner never calls the comment handler so semicolons
      inside string literals are always literal characters.
  5.  [ emits TOK_LPAREN. ] emits TOK_RPAREN. No further processing.
  6.  Identifier scanning:
        a. Entry condition: current byte is [A-Za-z_]
        b. Collect bytes while [A-Za-z0-9_]
        c. Check the next byte for suffix:
             # → append "#" to raw string and advance
             $ → append "$" to raw string and advance
             ? → append "?" to raw string and advance
        d. Call strings.ToUpper() on the complete string with suffix
        e. Look up the uppercase string in the keyword map
        f. If found: emit the keyword token type with the uppercase value
        g. If not found: emit TOK_IDENT with the uppercase value
  7.  Number scanning:
        a. Entry condition: current byte is [0-9] or '-' followed by [0-9]
        b. Collect digits
        c. If '.' follows and the byte after '.' is [0-9]: collect
           the decimal part and emit TOK_FLOAT
        d. Otherwise emit TOK_INTEGER
  8.  String scanning:
        a. Entry condition: current byte is "
        b. Advance past the opening quote
        c. Collect bytes until an unescaped " is found
        d. Handle escape sequences: \" → " and \\ → \
        e. Emit TOK_STRING with the resolved content
        f. NEVER call toupper on string contents
        g. An unterminated string is a LEXER ERROR with line + col
  9.  END keyword handling — peek at next non-whitespace word:
        END IF       → emit TOK_ENDIF
        END FUNCTION → emit TOK_ENDFUNCTION
        END WHILE    → emit TOK_WEND
        END SELECT   → emit TOK_ENDSELECT
        END TYPE     → emit TOK_ENDTYPE
        bare END     → emit TOK_END  (program termination)
  10. Unknown byte → LEXER ERROR with the offending character, line, col.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PARSER SPECIFICATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Recursive descent. One exported function per major grammar production.
The parser produces an AST. It does not evaluate or emit code.
The parser has no knowledge of raylib, physics, or any runtime library.

Top-level — no wrapping function required:
  program = { function_def | type_def | statement } EOF

Statement disambiguation rules — applied in this order:
  TOK_IDENT TOK_ASSIGN             → AssignNode
  TOK_IDENT TOK_LPAREN             → CallStmtNode (user function)
  TOK_IDENT TOK_DOT TOK_IDENT TOK_LPAREN → NamespaceCallStmtNode
  TOK_IDENT TOK_DOT TOK_IDENT TOK_ASSIGN → FieldAssignNode (type field)
  TOK_IF                           → IfNode
  TOK_WHILE                        → WhileNode
  TOK_FOR                          → ForNode
  TOK_REPEAT                       → RepeatNode
  TOK_SELECT                       → SelectNode
  TOK_FUNCTION                     → FunctionDefNode
  TOK_TYPE                         → TypeDefNode
  TOK_DIM or TOK_REDIM             → DimNode
  TOK_GOTO                         → GotoNode
  TOK_GOSUB                        → GosubNode
  TOK_RETURN                       → ReturnNode
  TOK_INCLUDE                      → IncludeNode
  TOK_DOT TOK_IDENT                → LabelNode
  TOK_NEWLINE                      → skip
  TOK_EOF                          → done
  anything else                    → PARSE ERROR

Expression precedence (lowest binding to highest binding):
  Must match compiler/parser/parser_expr.go (parseOr → parseXor → parseAnd → parseNot → …).
  Level 1:  OR
  Level 2:  XOR
  Level 3:  AND
  Level 4:  NOT  (unary prefix — chains: NOT NOT x)
  Level 5:  =  <>  <  >  <=  >=
  Level 6:  +  -  (binary)
  Level 7:  *  /  MOD
  Level 8:  ^  (right associative — use recursive call not loop)
  Level 9:  unary -
  Level 10: postfix — function call () or namespace call Ns.Method()
            or array index () / []
  Level 11: primary — literal, identifier, grouped (expr)

Call disambiguation in expression context:
  TOK_IDENT TOK_LPAREN                           → CallExprNode
  TOK_IDENT TOK_DOT TOK_IDENT TOK_LPAREN        → NamespaceCallExprNode
  TOK_IDENT (no following paren)                 → IdentNode (variable)
  Known builtin name without following paren     → PARSE ERROR with hint

AST node types required:
  ProgramNode        statements []Node, functions []FunctionDefNode
  AssignNode         name string, suffix string, expr Node
  FieldAssignNode    object string, field string, expr Node
  CallStmtNode       name string, args []Node
  NamespaceCallStmtNode  ns string, method string, args []Node
  CallExprNode       name string, args []Node
  NamespaceCallExprNode  ns string, method string, args []Node
  BinopNode          op string, left Node, right Node
  UnaryNode          op string, expr Node
  IdentNode          name string, suffix string
  IntLitNode         value int64
  FloatLitNode       value float64
  StringLitNode      value string
  BoolLitNode        value bool
  NullLitNode
  IfNode             condition Node, then []Node,
                     elseifs []ElseIfClause, else_ []Node
  WhileNode          condition Node, body []Node
  ForNode            var string, from Node, to Node,
                     step Node, body []Node
  RepeatNode         body []Node, condition Node
  SelectNode         expr Node, cases []CaseClause, default_ []Node
  FunctionDefNode    name string, params []Param, body []Node
  TypeDefNode        name string, fields []string
  ReturnNode         expr Node  (nil for void return)
  GotoNode           label string
  GosubNode          label string
  LabelNode          name string
  DimNode            name string, suffix string, dims []Node
  IncludeNode        path string
  NewNode            typeName string
  DeleteNode         expr Node
  EachNode           typeName string

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
VM BYTECODE SPECIFICATION
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Stack-based VM. All operands on the value stack.

Instructions are fixed-width in Go: **`Op` (byte)** + **`Operand int32`**
+ **`Aux int32`** + parallel **`SourceLines`** per instruction. Constant pools
and name tables live on each **`Chunk`** (`IntConsts`, `FloatConsts`,
`StringConsts`, `Names`).

Go names use the **`Op*`** prefix (e.g. **`OpAdd`**). Below, **`OP_*`** is the
conceptual mnemonic.

Opcodes:
  ; Stack manipulation
  OP_PUSH_INT        operand: int64 index into constant table
  OP_PUSH_FLOAT      operand: float64 index into constant table
  OP_PUSH_STRING     operand: string index into constant table
  OP_PUSH_BOOL       operand: 0 = false, 1 = true
  OP_PUSH_NULL
  OP_POP

  ; Variable access
  OP_LOAD_GLOBAL     operand: name index in string table
  OP_STORE_GLOBAL    operand: name index in string table
  OP_LOAD_LOCAL      operand: stack frame slot index
  OP_STORE_LOCAL     operand: stack frame slot index

  ; Arithmetic
  OP_ADD   OP_SUB   OP_MUL   OP_DIV   OP_MOD   OP_POW
  OP_NEG

  ; Comparison — result is bool Value on stack
  OP_EQ    OP_NEQ   OP_LT    OP_GT    OP_LTE   OP_GTE

  ; Logic
  OP_AND   OP_OR    OP_NOT   OP_XOR

  ; String
  OP_CONCAT

  ; Control flow
  OP_JUMP            operand: absolute instruction index
  OP_JUMP_IF_FALSE   operand: absolute instruction index
  OP_JUMP_IF_TRUE    operand: absolute instruction index

  ; Functions
  OP_CALL_BUILTIN    operand: name index in chunk Names; aux: argument count
  OP_CALL_USER       operand: function name index; aux: argument count
  OP_CALL_HANDLE     operand: method name index; aux: argument count
  OP_RETURN          operand: 1 if returning a value
  OP_RETURN_VOID

  ; Arrays
  OP_ARRAY_MAKE      operand: dimension count
  OP_ARRAY_GET       operand: dimension count
  OP_ARRAY_SET       operand: dimension count

  ; Types
  OP_NEW             operand: type name index
  OP_DELETE
  OP_FIELD_GET       operand: field name index
  OP_FIELD_SET       operand: field name index

  OP_HALT            end of program / explicit stop

**`Program`:** `Main *Chunk`, `Functions map[string]*Chunk`, `Types map[string]*TypeDef`.

**MOON `.mbc` file:** see **`vm/moon`** — header magic **`MOON`**, versioned,
then payload. Shipping uses this binary, not `encoding/gob`.

Value tagged union in Go:
  type Kind int
  const (
      KindNil Kind = iota
      KindInt
      KindFloat
      KindString
      KindBool
      KindHandle   ; integer handle into the Heap
  )

  type Value struct {
      Kind   Kind
      IVal   int64
      FVal   float64
      SVal   string
      BVal   bool
  }

  Handle values use IVal to store the integer handle.
  Arithmetic coerces Int ↔ Float automatically.
  String + anything coerces to string concatenation.

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
RUNTIME DISPATCH
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
All built-in commands live in a flat dispatch table:
  map[string]BuiltinFn  (on **`runtime.Registry`**)

Keys in the dispatch table are always uppercase:
  "WINDOW.OPEN"      "RENDER.FRAME"     "PHYSICS3D.START"
  "BODY3D.MAKE"      "PRINT"            "INPUT.KEYDOWN"

**Compile time:** `compiler/builtinmanifest/commands.json` is the **single
source of truth** for dotted engine commands. The semantic pass **rejects**
unknown `NS.METHOD` names **before** codegen whenever they are namespace calls.

**Run time:** `InitCore()` registers globals (`PRINT`, `STR$`, …).
`RegisterFromManifest(default manifest)` registers **stubs** for every remaining
manifest key so dispatch keys exist; unimplemented natives return a **runtime
error** with a clear message.

Optional future: **`runtime/aliases.go`** (or equivalent) for legacy spellings —
**not required** today; prefer fixing **`commands.json`** and the lexer/parser.

When engine code is split into modules, each module's **`Register`** runs at
startup and **`Shutdown`** at exit; order must not break dependencies (window
before render, etc.).

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
ERROR REPORTING FORMAT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
All errors must include the source location and a human-readable
hint wherever possible. The format is fixed:

  [moonBASIC] {Category} in {filename} line {line} col {col}:
    {message}

    {line_number} | {exact source line}
                  | {spaces}{^}
    Hint: {actionable suggestion}

Categories: Lexer Error / Parse Error / Type Error / CodeGen Error / Runtime Error

Example — missing parentheses:
  [moonBASIC] Parse Error in game.mbc line 14 col 8:
    Expected '(' after 'RENDERFRAME'

    14 | RENDERFRAME
       |        ^
    Hint: All commands require parentheses: RENDERFRAME()

Example — unknown dotted command (semantic, preferred):
  [moonBASIC] Type Error in game.mb line 22 col 5:
    unknown engine command PHYSICS3D.BEGINN (not in builtin manifest)

    22 | Physics3D.Beginn()
       |    ^
    Hint: Add the command to compiler/builtinmanifest/commands.json or use a supported API.

Example — unknown command at runtime (stub or typo in non-manifest call):
  [moonBASIC] Runtime Error (Line N): unknown command "…"

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CLI INTERFACE
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  moonbasic <source.mb>          compile source and run (in-memory), then exit
  moonbasic <file.mbc>           same as --run: load MOON bytecode and execute
  moonbasic --run <file.mbc>     run precompiled MOON bytecode
  moonbasic --compile <source> compile to <stem>.mbc next to source (--compile rejects a .mbc path)
  moonbasic --check <source>     parse + semantic only
  moonbasic --info <source>      disassemble main chunk to stderr, then run
  moonbasic --trace <source|mbc> print VM trace line after each opcode (stderr)
  moonbasic --version            print version string and exit
  moonbasic -h / --help          print usage and exit

Exit codes:
  0  success
  1  user error (bad arguments, file not found, write failure)
  2  compile error (lexer, parser, semantic, codegen) or bad MOON file
  3  runtime error (VM or builtin returned error)

Version string format (example):
  moonBASIC 1.2.x
  Runtime: Go 1.22 | raylib 5.5 | Jolt 5.1 | Box2D 3.0 | ENet 1.3
  Targets: Windows x64, Linux x64

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
IMPLEMENTATION ORDER
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Build in this sequence. Each step should compile and pass its tests before
moving on. **Steps 1–12 + MOON + manifest strictness reflect the current
codebase**; later steps are **Phase B** engine work.

  STEP 1:  token package — TokenType enum and Token struct
           Test: token types compile, string representations correct

  STEP 2:  lexer package — full tokenizer including suffix and END rules
           Test: tokenize the canonical reference program
           Verify: every token type, value, line, and column is correct

  STEP 3:  ast package — all AST node types
           Test: AST nodes compile and can be constructed

  STEP 4:  parser package — full recursive descent parser
           Test: parse the canonical reference program to AST
           Verify: print AST as indented text, check tree structure

  STEP 5:  errors package — MoonError type and formatting
           Test: format a lexer error and a parse error, check output

  STEP 6:  symtable package — symbol table with scope support
           Test: define and resolve variables across scopes

  STEP 7:  value + opcode packages — VM value type and opcode enum
           Test: Value arithmetic and coercion rules

  STEP 8:  heap package — handle store with HeapObject interface
           Test: allocate, retrieve, and free handles

  STEP 9:  callstack package — CallFrame struct
           Test: push and pop frames

  STEP 10: codegen package — AST to bytecode (split files; wire CallStmtNode)
           Test: compile x = 2 + 3 : PRINT(x) → outputs 5

  STEP 11: vm package — bytecode execution loop + trace flag
           Test: execute hello world program, arithmetic, loops, functions

  STEP 12: runtime/core — PRINT STR$ INT FLOAT LEN etc + Registry
           Test: hello program runs and produces correct output

  STEP 12b: builtinmanifest + semantic strict dotted commands
           Test: unknown FOO.BAR fails at compile time; reference program passes

  STEP 12c: vm/moon + pipeline EncodeMOON/DecodeMOON + CLI --compile / --run
           Test: go test ./vm/moon/... ; round-trip .mbc

  STEP 13: runtime/window + runtime/input + runtime/timemod
           Test: window opens, ESC closes it, delta time is valid

  STEP 14: runtime/render + runtime/draw + runtime/camera
           Test: coloured shapes appear on screen at correct positions

  STEP 15: runtime/model + runtime/texture + runtime/light
           Test: a model loads, displays, rotates with correct lighting

  STEP 16: runtime/terrain
           Test: terrain generates, layers apply, height query works

  STEP 17: runtime/audio
           Test: sound plays, music streams, volume control works

  STEP 18: runtime/physics3d (Jolt)
           Test: bodies fall, collide, bounce, raycast hits

  STEP 19: runtime/physics2d (Box2D)
           Test: 2D bodies simulate, joints work

  STEP 20: runtime/network (ENet)
           Test: server and client exchange reliable messages

  STEP 21: runtime/mathmod + runtime/strmod + runtime/filemod
           Test: math_demo.mbc runs all math operations correctly

  STEP 22: canonical reference program runs end to end
           Test: reference.mbc runs without errors on Windows and Linux

  STEP 23: full test suite passes on both platforms
           go build ./... && go test ./... — all green

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
CROSS-PLATFORM RULES
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
  Use filepath.Join() for all file path construction
  Normalise \r\n to \n in the lexer at read time
  Use os.ReadFile() — not platform-specific file APIs
  On Windows: ENet requires WSAStartup — call it in Net.Start()
  On Windows: raylib links against winmm and ws2_32
  On Linux:   raylib links against GL m pthread dl rt X11
  All CGo flags go in go.mod build constraints or cgo LDFLAGS
  Test on both platforms before marking a step complete

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
PROBLEM REPORT FORMAT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
When the build fails or a test fails, stop and report immediately.
Do not attempt to fix forward. Use this exact format:

  PROBLEM REPORT
  ==============
  Step: [step number and name from implementation order]
  File: [filename and line range]
  Error type: [compile / link / test / runtime]

  Full error output (complete and unedited):
  ------------------------------------------
  [paste here]

  Code before the change:
  ------------------------------------------
  [paste original]

  Code after the change:
  ------------------------------------------
  [paste modified]

  What you were trying to do:
  ------------------------------------------
  [one clear sentence]

Share this report with whoever is debugging (or keep it in your notes).
Before large refactors, re-read **`ARCHITECTURE.md`** so fixes align with the
stable pipeline and do not reintroduce removed layouts (e.g. empty `CallStmt`
codegen, stub-only `--compile`).
Do not modify unrelated files while bisecting a single failure unless the fix
requires it.
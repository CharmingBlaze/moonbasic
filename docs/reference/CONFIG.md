# Config file store (`CONFIG.*`)

**`CONFIG.LOAD`**, **`CONFIG.SAVE`**, **`CONFIG.GET`**, **`CONFIG.SET`**, etc., are implemented in **`runtime/mbgame/config.go`** against a **module-local** file-backed store (not a separate heap **`Config`** object). Paths and semantics match the **`CONFIG.*`** registrations in the manifest.

Use this for **settings.ini**-style persistence in small games; for larger data prefer **`JSON.*`** / **`FILE.*`** or your own format.

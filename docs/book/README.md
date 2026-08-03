# moonBASIC — The Book (downloadable editions)

Funny, mildly sweary, actually useful guide to the language and making games.

| Format | File |
|--------|------|
| Markdown (source of truth in repo) | [`../THE_MOONBASIC_BOOK.md`](../THE_MOONBASIC_BOOK.md) |
| **Word** | [`moonBASIC-The-Book.docx`](moonBASIC-The-Book.docx) |
| **PDF** | [`moonBASIC-The-Book.pdf`](moonBASIC-The-Book.pdf) |
| HTML (build intermediate) | [`moonBASIC-The-Book.html`](moonBASIC-The-Book.html) |

## Rebuild

From the engine repo root (needs Python packages `python-docx`, `markdown`, `xhtml2pdf`):

```bash
python tools/book/build_book_formats.py
```

Design: cream pages, navy tables, amber moonlight accents, Georgia body + Consolas code.

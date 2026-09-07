# How to View the Diagrams

## 📄 File Locations

All documentation files are in `C:\PROJECT\medix-be\docs\`:

| File | Size | Description |
|------|------|-------------|
| `architecture.mmd` | 894 bytes | System architecture diagram |
| `er-diagram.mmd` | 1,289 bytes | Database ER diagram |
| `openapi.yaml` | 29 KB | OpenAPI 3.0 specification |
| `README.md` | 1.1 KB | Usage instructions |

## 🖥️ How to Open/View Each File

### 1. **VS Code** (Recommended)

Install the **Mermaid** extension and **YAML** extension, then:

```bash
code docs/architecture.mmd   # Opens architecture diagram
code docs/er-diagram.mmd     # Opens ER diagram
code docs/openapi.yaml       # Opens OpenAPI spec (YAML highlighting)
```

VS Code will render Mermaid diagrams inline and provide YAML syntax highlighting.

### 2. **GitHub** (No install needed)

Simply visit the `docs/` directory in your GitHub repo. GitHub renders Mermaid diagrams automatically on markdown files (`.md`), but for `.mmd` files you can:

- Rename to `.md` and commit, OR
- Use GitHub's online editor: <https://mermaid.live>

### 3. **Mermaid Live Editor**

1. Go to <https://mermaid.live>
2. Copy the content of any `.mmd` file
3. Paste into the editor and view the rendered diagram
4. Export as PNG/SVG if needed

### 4. **Online YAML Viewers**

For `openapi.yaml`:

- <https://editor.openapi-generator.dev/>
- <https://yamlviewer.com/>
- VS Code (as mentioned above)

### 5. **Browser Extensions**

- **Mermaid Graphviz** (Chrome/Firefox) - renders `.mmd` files directly
- **YAML Viewer** extensions for browser

## 🛠️ Quick Commands

```bash
# Open in VS Code
code docs/architecture.mmd
code docs/er-diagram.mmd

# Render to PNG (requires Mermaid CLI)
npx @mermaid-js/mermaid-cli -i docs/architecture.mmd -o docs/architecture.png
npx @mermaid-js/mermaid-cli -i docs/er-diagram.mmd -o docs/er-diagram.png
```

## 📋 File Contents Summary

### `architecture.mmd`
Shows the high-level system architecture with layers: Client → API → Domain Modules → Data (PostgreSQL).

### `er-diagram.mmd`
Shows the database entity-relationship diagram with Users, Obat, Jenis Obat, Transaksi, and Detail Pembelian tables.

### `openapi.yaml`
Full OpenAPI 3.0 specification covering all 50+ endpoints across Auth, User, Medicine, Drug Type, Transaction, and Report modules.

### `README.md`
Instructions on how to view/render all files.

## 💡 Tips

- **VS Code tip**: Press `Ctrl+Shift+V` (or `Cmd+Shift+V` on Mac) to **preview** the Markdown/Mermaid rendering
- **GitHub tip**: Push to GitHub and view the `.md` versions in the repo for automatic rendering
- **Offline**: All files are plain text, no build step or external services required beyond a text editor
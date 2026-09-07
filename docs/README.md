# Medix BE Documentation

## 📁 Generated Files

| File | Description | Render With |
|------|-------------|-------------|
| `architecture.mmd` | System architecture diagram | Mermaid Live, VS Code, GitHub |
| `er-diagram.mmd` | Database ER diagram | Mermaid Live, VS Code, GitHub |
| `openapi.yaml` | OpenAPI 3.0 specification | Swagger UI, Redoc, Postman |

## 🚀 Quick View

### Architecture Diagram
```bash
# View in VS Code with Mermaid extension
code docs/architecture.mmd
```

### ER Diagram
```bash
code docs/er-diagram.mmd
```

### OpenAPI Spec (Swagger UI)
```bash
# Serve with swagger-ui
npx swagger-ui-serve docs/openapi.yaml

# Or use Docker
docker run -p 8080:8080 -v ${PWD}/docs/openapi.yaml:/app/openapi.yaml swaggerapi/swagger-ui
```

## 📊 Diagrams Preview

### Architecture
![Architecture](architecture.mmd)

### Entity Relationship
![ER Diagram](er-diagram.mmd)

## 🔧 Usage

- **Mermaid diagrams**: Paste into [Mermaid Live Editor](https://mermaid.live) or view in GitHub/VS Code
- **OpenAPI spec**: Import into Postman, or render with Swagger UI/Redoc
- **All files**: Plain text, version-controlled, no build step required
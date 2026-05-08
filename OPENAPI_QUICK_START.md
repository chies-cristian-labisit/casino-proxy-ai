# OpenAPI Viewers - Quick Start

## 🚀 Start in 30 Seconds

```bash
docker-compose -f docker-compose.openapi.yml up -d
```

## 🌐 Access Points

Open these URLs in your browser:

### Main API (Casino Proxy - 8 Providers)
- **Edit specs:** http://localhost:8080 (Swagger Editor)
- **Test API:** http://localhost:8081 (Swagger UI - Try It Out)
- **Read docs:** http://localhost:8082 (Redoc - Beautiful)

### Admin API (Operator Management)
- **Edit specs:** http://localhost:8083 (Swagger Editor)
- **Test API:** http://localhost:8084 (Swagger UI)
- **Read docs:** http://localhost:8085 (Redoc)

## 📋 What Each Tool Does

| Tool | URL | Use For |
|------|-----|---------|
| **Swagger Editor** | 8080/8083 | Editing & validating OpenAPI YAML specs |
| **Swagger UI** | 8081/8084 | Testing endpoints with "Try It Out" button |
| **Redoc** | 8082/8085 | Beautiful documentation (share with others) |

## 💡 Quick Examples

### Test an Endpoint
1. Open http://localhost:8081 (Swagger UI)
2. Find endpoint (e.g., "Create Operator")
3. Click "Try It Out"
4. Fill in parameters
5. Click "Execute"
6. See response in real-time

### Edit Documentation
1. Open http://localhost:8080 (Swagger Editor)
2. Edit YAML on left side
3. See validation errors and preview on right
4. Download updated spec when done

### Share with Team
1. Open http://localhost:8082 (Redoc)
2. Share the URL with your team
3. They see beautiful, professional documentation
4. No "Try It Out" - read-only

## 🛑 Stop Viewers

```bash
docker-compose -f docker-compose.openapi.yml down
```

## 🔍 Check Status

```bash
docker-compose -f docker-compose.openapi.yml ps
```

## 📚 Full Documentation

See [docs/OPENAPI_VIEWERS.md](docs/OPENAPI_VIEWERS.md) for:
- Detailed tool guides
- Common workflows
- Troubleshooting
- Advanced usage
- CI/CD integration

---

**That's it!** Your OpenAPI documentation is now running locally.

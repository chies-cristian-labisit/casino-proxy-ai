# OpenAPI Viewers Module

Local Docker-based OpenAPI documentation viewers for the Casino Proxy API.

## 📚 What's in This Module

This module provides a complete, self-contained setup for viewing and editing OpenAPI specifications locally using three industry-standard tools:

### Tools Included
1. **Swagger Editor** - Interactive OpenAPI spec editor with live validation
2. **Swagger UI** - Beautiful API documentation with "Try It Out" testing
3. **Redoc** - Professional API documentation for sharing

### APIs Documented
- **Main API** - Casino Proxy gateway (all 8 gaming providers)
- **Admin API** - Operator management and system endpoints

## 🚀 Quick Start

```bash
# Navigate to this directory
cd application/docs/openapi-viewers

# Start all services
docker-compose -f docker-compose.yml up -d

# Access in browser
open http://localhost:8080  # Swagger Editor (Main)
open http://localhost:8081  # Swagger UI (Main)
open http://localhost:8082  # Redoc (Main)
```

## 📂 Module Structure

```
application/docs/openapi-viewers/
├── docker-compose.yml          # Docker Compose configuration
├── docker/
│   ├── swagger-ui-main.env    # Nginx config for main API UI
│   └── swagger-ui-admin.env   # Nginx config for admin API UI
├── README.md                   # This file
├── QUICK_START.md              # 30-second setup guide
└── OPENAPI_VIEWERS.md          # Comprehensive documentation
```

## 🔌 Services & Ports

### Main API Viewers (casino-proxy-api.yaml auto-loaded)
| Port | Container Name | Tool | URL | Spec |
|------|---|------|-----|------|
| 8080 | casino-proxy-swagger-editor-main | Swagger Editor | http://localhost:8080 | casino-proxy-api.yaml |
| 8081 | casino-proxy-swagger-ui-main | Swagger UI | http://localhost:8081 | casino-proxy-api.yaml |
| 8082 | casino-proxy-redoc-main | Redoc | http://localhost:8082 | casino-proxy-api.yaml |

### Admin API Viewers (admin-api-spec.yaml auto-loaded)
| Port | Container Name | Tool | URL | Spec |
|------|---|------|-----|------|
| 8083 | casino-proxy-swagger-editor-admin | Swagger Editor | http://localhost:8083 | admin-api-spec.yaml |
| 8084 | casino-proxy-swagger-ui-admin | Swagger UI | http://localhost:8084 | admin-api-spec.yaml |
| 8085 | casino-proxy-redoc-admin | Redoc | http://localhost:8085 | admin-api-spec.yaml |

## ✨ Auto-Loaded Specs

All containers automatically load their OpenAPI specification files on startup:
- **Swagger Editor** reads spec from `SWAGGER_FILE` environment variable
- **Swagger UI** loads `SWAGGER_JSON` environment variable
- **Redoc** loads spec from command-line argument

## 📖 Documentation

- **[QUICK_START.md](QUICK_START.md)** - 30-second setup and basic usage
- **[OPENAPI_VIEWERS.md](OPENAPI_VIEWERS.md)** - Comprehensive guide
  - Tool guides and features
  - Common workflows
  - Docker commands reference
  - Troubleshooting
  - Advanced usage

## 🎯 Use Cases

### For Developers
- Edit OpenAPI specifications
- Test API endpoints locally
- Generate server stubs and client SDKs
- Debug API responses

### For Architects
- View complete API documentation
- Validate specification structure
- Design new endpoints
- Review existing integrations

### For Stakeholders
- Share professional API documentation
- Explore endpoint specifications
- Understand data models
- Review authentication requirements

## ✨ Key Features

✅ **6 Services** - 3 viewers for main API, 3 for admin API  
✅ **30-Second Setup** - Docker Compose with one command  
✅ **Live Validation** - Swagger Editor with instant error feedback  
✅ **Interactive Testing** - Swagger UI "Try It Out" functionality  
✅ **Professional Docs** - Redoc for stakeholder sharing  
✅ **Health Checks** - Automatic container health monitoring  
✅ **Network Isolation** - Services on dedicated Docker network  

## 🛠 Commands

### Start All Services
```bash
docker-compose -f docker-compose.yml up -d
```

### Stop All Services
```bash
docker-compose -f docker-compose.yml down
```

### View Logs
```bash
docker-compose -f docker-compose.yml logs -f
```

### Check Status
```bash
docker-compose -f docker-compose.yml ps
```

### Restart Services
```bash
docker-compose -f docker-compose.yml restart
```

## 📝 Notes

- All volume mounts use current directory (.)
- Services are on the `casino-docs` network
- Health checks configured for all containers
- Swagger UI runs on ports 8081 and 8084 (mapped from 8081 in container)
- All services have labels for easy identification

## 🔗 Related Documentation

- **Main API Spec:** `casino-proxy-api.yaml`
- **Admin API Spec:** `admin-api-spec.yaml`
- **API Guides:** Parent folder (`../api-guide/`)
  - `authentication.md` - Authentication methods
  - `webhook-integration.md` - Webhook patterns
  - `error-handling.md` - Error scenarios
  - `admin-operations.md` - Admin endpoints

## 🚨 Troubleshooting

### Port Already in Use
Edit `docker-compose.yml` and change port numbers, or kill existing process:
```bash
lsof -i :8080
kill -9 <PID>
```

### Container Won't Start
Check logs:
```bash
docker-compose -f docker-compose.yml logs swagger-editor-main
```

### Spec File Not Found
Verify files exist in current directory:
```bash
ls -la casino-proxy-api.yaml
ls -la admin-api-spec.yaml
```

For more troubleshooting, see [OPENAPI_VIEWERS.md](OPENAPI_VIEWERS.md).

## 📚 External Resources

- [Swagger Documentation](https://swagger.io/docs/)
- [Redoc GitHub](https://github.com/Redocly/redoc)
- [OpenAPI 3.0 Specification](https://spec.openapis.org/oas/v3.0.3)
- [Docker Compose Reference](https://docs.docker.com/compose/)

---

**Last Updated:** 2026-05-08  
**Status:** ✅ Production Ready

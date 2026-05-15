# Casino Proxy OpenAPI Viewers & Editors

This document explains how to set up and use the local OpenAPI documentation viewers for the Casino Proxy API.

## Quick Start

### Start All Viewers

```bash
docker-compose -f docker-compose.openapi.yml up -d
```

### Stop All Viewers

```bash
docker-compose -f docker-compose.openapi.yml down
```

### View Logs

```bash
docker-compose -f docker-compose.openapi.yml logs -f
```

### Check Health Status

```bash
docker-compose -f docker-compose.openapi.yml ps
```

---

## Access Points

### Main API Documentation (Casino Proxy - All 8 Providers)

| Tool | URL | Purpose | Best For |
|------|-----|---------|----------|
| **Swagger Editor** | http://localhost:8080 | Edit & validate specs | Modifying OpenAPI YAML |
| **Swagger UI** | http://localhost:8081 | View & test endpoints | Testing API calls, "Try It Out" |
| **Redoc** | http://localhost:8082 | Beautiful documentation | Sharing with stakeholders |

### Admin API Documentation

| Tool | URL | Purpose | Best For |
|------|-----|---------|----------|
| **Swagger Editor** | http://localhost:8083 | Edit & validate admin spec | Modifying admin API |
| **Swagger UI** | http://localhost:8084 | View & test admin endpoints | Testing operator management |
| **Redoc** | http://localhost:8085 | Beautiful admin documentation | Admin documentation sharing |

---

## Tool Guide

### 1. Swagger Editor (Port 8080 & 8083)

**URL:** http://localhost:8080 (Main) | http://localhost:8083 (Admin)

**Purpose:** Edit OpenAPI specifications with live validation

**Features:**
- ✅ Edit YAML/JSON directly
- ✅ Live validation and error highlighting
- ✅ Generate server stubs (Node.js, Go, Python, etc.)
- ✅ Generate client SDKs
- ✅ Download spec in multiple formats
- ✅ Schema validation

**Use Cases:**
1. **Modify API documentation**
   - Add new endpoints
   - Update field descriptions
   - Fix validation rules
   - Add examples

2. **Validate OpenAPI spec syntax**
   - Instant feedback on errors
   - Check schema references
   - Validate security schemes

3. **Generate code**
   - Generate server stubs
   - Generate client libraries
   - Create mock servers

**Example Workflow:**
```
1. Open http://localhost:8080
2. Edit the YAML spec on the left
3. See validation errors on the right
4. Preview rendered documentation in the center
5. Export modified spec when done
```

### 2. Swagger UI (Port 8081 & 8084)

**URL:** http://localhost:8081 (Main) | http://localhost:8084 (Admin)

**Purpose:** Interactive API documentation with "Try It Out" functionality

**Features:**
- ✅ Beautiful, responsive documentation
- ✅ "Try It Out" to test endpoints
- ✅ Generate curl commands
- ✅ View request/response examples
- ✅ Explore all endpoints and schemas
- ✅ Filter by provider (if applicable)
- ✅ Real-time API testing

**Use Cases:**
1. **Test API endpoints**
   ```
   1. Navigate to endpoint
   2. Click "Try It Out"
   3. Fill in parameters
   4. Click "Execute"
   5. See response in real-time
   ```

2. **Share with developers**
   - Send URL to team members
   - They can explore and test endpoints
   - View examples and schemas

3. **Integration testing**
   - Verify endpoint behavior
   - Test authentication
   - Test error scenarios
   - Check response formats

**Example Testing Scenario:**
```
Testing: POST /v1/internal/operator/store

1. Navigate to "Operators" section
2. Find "Create Operator" endpoint
3. Click "Try It Out"
4. Fill in:
   - slug: "test-operator"
   - name: "Test Casino"
   - url: "https://test.example.com/webhooks"
5. Click "Execute"
6. See 201 response with created operator
7. Copy curl command for automation
```

### 3. Redoc (Port 8082 & 8085)

**URL:** http://localhost:8082 (Main) | http://localhost:8085 (Admin)

**Purpose:** Modern, beautiful API documentation

**Features:**
- ✅ Responsive design (mobile-friendly)
- ✅ Clean, professional documentation
- ✅ Excellent for stakeholder sharing
- ✅ No "Try It Out" (read-only)
- ✅ Excellent schema exploration
- ✅ Fast rendering
- ✅ Search functionality

**Use Cases:**
1. **Share with stakeholders**
   - Non-technical users can read docs
   - Professional appearance
   - Mobile-friendly

2. **Create documentation website**
   - Run Redoc in CI/CD pipeline
   - Generate static HTML docs
   - Deploy to documentation site

3. **API portal**
   - Beautiful endpoint reference
   - Schema documentation
   - Example payloads

**Example: Sharing with External Partners**
```
1. Run Redoc on http://localhost:8082
2. Share link with partners
3. They see professional API documentation
4. They can search and explore schemas
5. No confusion from "Try It Out" button
```

---

## Common Workflows

### Workflow 1: Add a New Admin Endpoint

```
1. Open Swagger Editor for Admin API
   → http://localhost:8083

2. Edit admin-api-spec.yaml
   - Add new endpoint under paths
   - Define request body schema
   - Define response schemas
   - Add error responses

3. Review in Swagger UI
   → http://localhost:8084
   - Check documentation rendering
   - Test with "Try It Out"

4. Preview in Redoc
   → http://localhost:8085
   - Verify professional appearance
   - Check mobile rendering

5. Export modified spec
   - Download YAML
   - Update application/docs/openapi/admin-api-spec.yaml
   - Commit changes
```

### Workflow 2: Document a New Provider

```
1. Create provider-spec.yaml in Swagger Editor
   - Define all webhook endpoints
   - Add authentication scheme
   - Add request/response examples
   - Add error codes

2. Test examples in Swagger UI
   - Use "Try It Out" to verify payloads
   - Generate curl commands
   - Document API behavior

3. Add to master spec
   - Merge provider spec into casino-proxy-api.yaml
   - Update provider list
   - Verify references

4. Share documentation
   - Generate via Redoc
   - Publish to documentation portal
```

### Workflow 3: Debug API Response

```
1. Go to Swagger UI
   → http://localhost:8081

2. Find the problematic endpoint
   - Navigate to provider section
   - Find the endpoint

3. Click "Try It Out"
   - Fill in request parameters
   - Copy the generated curl command

4. Execute request
   - See actual response
   - Check status code
   - Review error messages

5. Compare with documentation
   - Verify response structure matches spec
   - Check if documentation needs update
   - Update Swagger Editor if needed
```

### Workflow 4: Generate Client SDK

```
1. Open Swagger Editor
   → http://localhost:8080

2. Click "Generate Code" (top menu)
   - Select language (Go, Python, Node.js, etc.)

3. Select target platform
   - Choose client or server
   - Choose language/framework

4. Generate and download
   - Client code is generated
   - Download as ZIP file
   - Integrate into your project
```

---

## Configuration Files

### docker-compose.openapi.yml
Main configuration file defining all services:
- 3 viewers for main API (Editor, UI, Redoc)
- 3 viewers for admin API (Editor, UI, Redoc)
- Networks for container communication
- Health checks for all services
- Volume mounts for spec files

### docker/openapi/swagger-ui-main.env
Nginx configuration for main API Swagger UI:
- Port binding
- Spec file serving
- Static file handling

### docker/openapi/swagger-ui-admin.env
Nginx configuration for admin API Swagger UI:
- Port binding
- Spec file serving
- Static file handling

---

## Docker Commands Reference

### Start Services
```bash
# Start all services
docker-compose -f docker-compose.openapi.yml up -d

# Start specific service
docker-compose -f docker-compose.openapi.yml up -d swagger-ui-main

# Start in foreground (see logs)
docker-compose -f docker-compose.openapi.yml up
```

### Stop Services
```bash
# Stop all services
docker-compose -f docker-compose.openapi.yml down

# Stop and remove volumes
docker-compose -f docker-compose.openapi.yml down -v

# Stop specific service
docker-compose -f docker-compose.openapi.yml stop swagger-editor-main
```

### View Logs
```bash
# View logs for all services
docker-compose -f docker-compose.openapi.yml logs -f

# View logs for specific service
docker-compose -f docker-compose.openapi.yml logs -f swagger-ui-main

# View last 100 lines
docker-compose -f docker-compose.openapi.yml logs --tail=100
```

### Restart Services
```bash
# Restart all services
docker-compose -f docker-compose.openapi.yml restart

# Restart specific service
docker-compose -f docker-compose.openapi.yml restart swagger-editor-main
```

### Health Check
```bash
# Check service status
docker-compose -f docker-compose.openapi.yml ps

# Get detailed health info
docker inspect casino-swagger-editor-main | grep -A 10 '"Health"'
```

---

## Troubleshooting

### Port Already in Use

**Problem:** Port 8080 already in use

**Solution 1 - Change port in docker-compose.openapi.yml:**
```yaml
swagger-editor-main:
  ports:
    - "9000:8080"  # Changed from 8080 to 9000
```

**Solution 2 - Kill existing process:**
```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Spec File Not Found

**Problem:** Volume mount error

**Solution:** Verify file paths
```bash
# Check if files exist
ls -la application/docs/openapi/casino-proxy-api.yaml
ls -la application/docs/openapi/admin-api-spec.yaml

# Verify docker-compose.openapi.yml paths are correct
cat docker-compose.openapi.yml | grep volumes
```

### Container Won't Start

**Problem:** Container exits immediately

**Solution:** Check logs
```bash
docker-compose -f docker-compose.openapi.yml logs swagger-editor-main
```

Common causes:
- Port already in use
- Invalid volume mount
- Image not found (pull images first)

**Fix:**
```bash
# Pull latest images
docker-compose -f docker-compose.openapi.yml pull

# Restart with no cache
docker-compose -f docker-compose.openapi.yml up --force-recreate
```

### Swagger UI Not Loading Spec

**Problem:** Spec file appears blank in Swagger UI

**Solution:** Verify YAML syntax
```bash
# Use yamllint to validate spec
docker run -v $(pwd):/app -w /app sdesbure/yamllint application/docs/openapi/casino-proxy-api.yaml

# Or manually check for syntax errors
cat application/docs/openapi/casino-proxy-api.yaml | head -20
```

---

## Advanced Usage

### Mount Writable Volume for Editing

To edit specs directly in containers and save back to host:

```yaml
swagger-editor-main:
  volumes:
    - ./application/docs/openapi:/tmp/specs  # Remove :ro to allow writes
```

**Warning:** Ensure container user has write permissions.

### Custom Swagger UI Configuration

Create custom Swagger UI HTML:

```html
<!-- docker/openapi/swagger-ui.html -->
<!DOCTYPE html>
<html>
  <head>
    <title>Casino Proxy API</title>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@3/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@3/swagger-ui-bundle.js"></script>
    <script>
    SwaggerUIBundle({
        url: "/specs/casino-proxy-api.yaml",
        dom_id: '#swagger-ui',
        presets: [
            SwaggerUIBundle.presets.apis,
            SwaggerUIBundle.SwaggerUIStandalonePreset
        ],
        tryItOutEnabled: true,
        persistAuthorization: true
    })
    </script>
  </body>
</html>
```

### Generate Static Documentation

```bash
# Generate HTML documentation from Redoc
docker run -v $(pwd):/app -w /app redocly/redoc \
  bundle application/docs/openapi/casino-proxy-api.yaml \
  -o docs/api-reference.html

# Serve static documentation
python -m http.server 8000 --directory docs/
```

---

## Integration with CI/CD

### GitHub Actions Example

```yaml
name: Validate OpenAPI Specs

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Validate OpenAPI specs
        run: |
          docker run -v $(pwd):/app swaggerapi/swagger-ui \
            curl -f http://localhost:8080 || exit 1
      
      - name: Generate documentation
        run: |
          docker run -v $(pwd):/app redocly/redoc \
            bundle /app/application/docs/openapi/casino-proxy-api.yaml \
            -o /app/docs/api-reference.html
      
      - name: Upload documentation
        uses: actions/upload-artifact@v2
        with:
          name: api-documentation
          path: docs/api-reference.html
```

---

## Performance Tips

### Reduce Memory Usage

```yaml
swagger-editor-main:
  deploy:
    resources:
      limits:
        memory: 256M
      reservations:
        memory: 128M
```

### Disable Health Checks (if needed)

```yaml
swagger-editor-main:
  healthcheck:
    disable: true
```

### Use Lighter Images

```yaml
# Replace swaggerapi/swagger-editor with prism (lighter)
prism-mock:
  image: stoplight/prism-http:latest
  ports:
    - "8000:4010"
  volumes:
    - ./application/docs/openapi:/specs:ro
  command: ["mock", "-h", "0.0.0.0", "/specs/casino-proxy-api.yaml"]
```

---

## Best Practices

1. **Always validate YAML before pushing**
   - Use yamllint locally
   - Catch syntax errors early

2. **Keep specs in sync with code**
   - Update OpenAPI spec when API changes
   - Document examples with actual responses
   - Add new endpoints to spec before coding

3. **Use semantic versioning for APIs**
   - Version in info.version
   - Document breaking changes
   - Deprecate old endpoints

4. **Document with examples**
   - Add realistic request examples
   - Include error response examples
   - Show authentication headers

5. **Test with Try It Out**
   - Verify endpoints work as documented
   - Catch response format mismatches
   - Validate error codes

6. **Share via Redoc**
   - Use for stakeholder communication
   - Professional appearance
   - Mobile-friendly documentation

---

## Next Steps

1. Start the viewers:
   ```bash
   docker-compose -f docker-compose.openapi.yml up -d
   ```

2. Open in browser:
   - **Editing:** http://localhost:8080 (Swagger Editor)
   - **Testing:** http://localhost:8081 (Swagger UI)
   - **Sharing:** http://localhost:8082 (Redoc)

3. Explore the API documentation

4. Test endpoints with "Try It Out"

5. Share links with your team

---

## Support

For issues or questions:
- **Swagger Documentation:** https://swagger.io/docs/
- **Redoc GitHub:** https://github.com/Redocly/redoc
- **OpenAPI Spec:** https://spec.openapis.org/oas/v3.0.3

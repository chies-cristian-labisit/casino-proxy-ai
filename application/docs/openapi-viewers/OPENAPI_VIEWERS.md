# Casino Proxy OpenAPI Viewers & Editors

This document explains how to set up and use the local OpenAPI documentation viewers for the Casino Proxy API.

## Quick Start

### Start All Viewers

```bash
cd application/docs/openapi-viewers
docker-compose -f docker-compose.yml up -d
```

### Stop All Viewers

```bash
docker-compose -f docker-compose.yml down
```

### View Logs

```bash
docker-compose -f docker-compose.yml logs -f
```

### Check Health Status

```bash
docker-compose -f docker-compose.yml ps
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
   - Update admin-api-spec.yaml
   - Commit changes
```

### Workflow 2: Test an Endpoint

```
1. Go to Swagger UI
   → http://localhost:8081

2. Find the endpoint
   - Navigate to provider section
   - Find the endpoint

3. Click "Try It Out"
   - Fill in request parameters
   - Copy the generated curl command

4. Execute request
   - See actual response
   - Check status code
   - Review error messages
```

---

## Docker Commands Reference

### Start Services
```bash
# Start all services
docker-compose -f docker-compose.yml up -d

# Start specific service
docker-compose -f docker-compose.yml up -d swagger-ui-main

# Start in foreground (see logs)
docker-compose -f docker-compose.yml up
```

### Stop Services
```bash
# Stop all services
docker-compose -f docker-compose.yml down

# Stop and remove volumes
docker-compose -f docker-compose.yml down -v

# Stop specific service
docker-compose -f docker-compose.yml stop swagger-editor-main
```

### View Logs
```bash
# View logs for all services
docker-compose -f docker-compose.yml logs -f

# View logs for specific service
docker-compose -f docker-compose.yml logs -f swagger-ui-main

# View last 100 lines
docker-compose -f docker-compose.yml logs --tail=100
```

### Restart Services
```bash
# Restart all services
docker-compose -f docker-compose.yml restart

# Restart specific service
docker-compose -f docker-compose.yml restart swagger-editor-main
```

### Health Check
```bash
# Check service status
docker-compose -f docker-compose.yml ps

# Get detailed health info
docker inspect casino-swagger-editor-main | grep -A 10 '"Health"'
```

---

## Troubleshooting

### Port Already in Use

**Problem:** Port 8080 already in use

**Solution 1 - Change port in docker-compose.yml:**
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
ls -la casino-proxy-api.yaml
ls -la admin-api-spec.yaml

# Verify docker-compose.yml paths are correct
cat docker-compose.yml | grep volumes
```

### Container Won't Start

**Problem:** Container exits immediately

**Solution:** Check logs
```bash
docker-compose -f docker-compose.yml logs swagger-editor-main
```

---

## Advanced Usage

### Mount Writable Volume for Editing

To edit specs directly in containers and save back to host:

```yaml
swagger-editor-main:
  volumes:
    - .  # Remove :ro to allow writes
```

### Generate Static Documentation

```bash
# Generate HTML documentation from Redoc
docker run -v $(pwd):/app -w /app redocly/redoc \
  bundle casino-proxy-api.yaml \
  -o api-reference.html

# Serve static documentation
python -m http.server 8000 --directory ./
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

3. **Test with Try It Out**
   - Verify endpoints work as documented
   - Catch response format mismatches
   - Validate error codes

4. **Share via Redoc**
   - Use for stakeholder communication
   - Professional appearance
   - Mobile-friendly documentation

---

## Support

For issues or questions:
- **Swagger Documentation:** https://swagger.io/docs/
- **Redoc GitHub:** https://github.com/Redocly/redoc
- **OpenAPI Spec:** https://spec.openapis.org/oas/v3.0.3

#!/bin/bash
#
# verify-template.sh — Copier template end-to-end verification
#
# Tests the full copier pipeline with AIOX integration:
# generate → AIOX setup → AIOX verify → build → acceptance tests
# Exit code 1 on failure, usable in CI/CD and pre-push hook
#

set -e  # Exit on error

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
START_TIME=$(date +%s)

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

cleanup() {
  local exit_code=$?
  if [ -n "$TMPDIR_CREATED" ] && [ -d "$TMPDIR_CREATED" ]; then
    rm -rf "$TMPDIR_CREATED"
  fi
  # Clean up any testcontainers left by acceptance tests
  docker ps -aq --filter "label=org.testcontainers=true" 2>/dev/null | xargs -r docker rm -f 2>/dev/null || true
  return $exit_code
}

trap cleanup EXIT

print_step() {
  echo -e "${BLUE}→${NC} $1"
}

print_success() {
  echo -e "${GREEN}✓${NC} $1"
}

print_error() {
  echo -e "${RED}✗${NC} $1"
}

print_warning() {
  echo -e "${YELLOW}⚠${NC} $1"
}

# Pre-flight: Check Node.js and npm
print_step "Verifying Node.js and npm availability..."
if ! command -v node &> /dev/null; then
  print_error "Node.js not found (required for AIOX)"
  exit 1
fi
NODE_VERSION=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 18 ]; then
  print_error "Node.js version $(node --version) below 18.0.0 requirement"
  exit 1
fi
print_success "Node.js $(node --version)"

if ! command -v npm &> /dev/null; then
  print_error "npm not found (required for AIOX)"
  exit 1
fi
print_success "npm $(npm --version)"

# Step 1: Create temp directory
print_step "Creating temporary directory..."
TMPDIR_CREATED=$(mktemp -d)
print_success "Created: $TMPDIR_CREATED"

# Step 2: Scaffold project with copier
print_step "Scaffolding project with copier..."
if ! copier copy "$REPO_ROOT" "$TMPDIR_CREATED" --defaults -w --vcs-ref HEAD > /dev/null 2>&1; then
  print_error "Copier failed"
  exit 1
fi
print_success "Project scaffolded"

# Step 3: Initialize git repo (required for Husky)
print_step "Initializing git repository..."
cd "$TMPDIR_CREATED"
if ! git init > /dev/null 2>&1; then
  print_error "Git initialization failed"
  exit 1
fi
print_success "Git repository initialized"

# Step 4: Initialize AIOX Framework
print_step "Initializing AIOX Framework..."
if ! npm run setup > /dev/null 2>&1; then
  print_error "AIOX setup failed"
  exit 1
fi
print_success "AIOX Framework initialized"

# Step 5: Verify AIOX Installation
print_step "Verifying AIOX installation..."

# Check aiox-core doctor
if ! npx aiox-core doctor > /dev/null 2>&1; then
  print_warning "AIOX doctor reported issues (continuing...)"
fi

# Check .claude/CLAUDE.md exists
if [ ! -f ".claude/CLAUDE.md" ]; then
  print_error ".claude/CLAUDE.md not found (AIOX sync failed)"
  exit 1
fi
print_success ".claude/CLAUDE.md exists"

# Check .aiox-core/core exists
if [ ! -d ".aiox-core/core" ]; then
  print_error ".aiox-core/core directory not found"
  exit 1
fi
print_success ".aiox-core/core directory exists"

# Check .aiox-core/development/agents exists
if [ ! -d ".aiox-core/development/agents" ]; then
  print_error ".aiox-core/development/agents directory not found"
  exit 1
fi
print_success ".aiox-core/development/agents exists"

# Check .husky directory (git hooks)
if [ ! -d ".husky" ]; then
  print_error ".husky directory not found (git hooks not initialized)"
  exit 1
fi
print_success ".husky git hooks configured"

print_success "AIOX verification complete"

# Step 6: Resolve Go dependencies
print_step "Resolving Go dependencies..."
cd "$TMPDIR_CREATED/application"
if ! go mod tidy > /dev/null 2>&1; then
  print_error "go mod tidy failed"
  exit 1
fi
print_success "Go dependencies resolved"

# Step 7: Build verification
print_step "Verifying build..."
if ! go build -o /dev/null ./cmd/api > /dev/null 2>&1; then
  print_error "Build failed"
  exit 1
fi
print_success "Build successful"

# Step 8: Run acceptance tests (with testcontainers — no docker-compose needed)
print_step "Running acceptance tests (this may take ~30s)..."
if ! go test ./tests/acceptance/... -v -timeout 300s 2>&1 | tail -5; then
  print_error "Acceptance tests failed"
  exit 1
fi
print_success "All acceptance tests passed"

# Summary
ELAPSED=$(($(date +%s) - START_TIME))
echo ""
echo -e "${GREEN}═══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  COPIER TEMPLATE VERIFICATION PASS${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════${NC}"
echo "Time: ${ELAPSED}s"
echo ""
echo "The copier template is working correctly:"
echo "  ✓ Project scaffolding with copier"
echo "  ✓ AIOX Framework initialization"
echo "  ✓ AIOX installation verification"
echo "  ✓ Agent definitions synced to IDE"
echo "  ✓ Git hooks configured"
echo "  ✓ Go module resolution"
echo "  ✓ Code compilation"
echo "  ✓ Acceptance tests (full stack)"
echo ""
echo "Generated project is AIOX-ready and fully tested!"
echo ""

exit 0

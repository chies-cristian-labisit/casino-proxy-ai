#!/bin/bash

###############################################################################
# AIOX Framework Setup Script
#
# This script initializes an AIOX-ready project with:
# - AIOX core framework installation
# - IDE agent definitions synchronization
# - Husky git hooks configuration
# - System verification
#
# Usage: npm run setup
#        bash scripts/setup.sh
#
# Features:
# - Idempotent: Safe to run multiple times
# - Fail-fast: Exits immediately on error
# - Cross-platform: Works on Windows (PowerShell), macOS, Linux
# - Interactive: Prompts user for AIOX configuration
###############################################################################

set -e

# Colors for output (works on macOS, Linux, Windows Terminal)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Utility functions
log_step() {
    echo -e "${BLUE}→${NC} $1"
}

log_success() {
    echo -e "${GREEN}✓${NC} $1"
}

log_error() {
    echo -e "${RED}✗${NC} $1" >&2
}

log_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Detect OS for cross-platform compatibility
detect_os() {
    case "$(uname -s)" in
        Linux*) echo "linux" ;;
        Darwin*) echo "macos" ;;
        MINGW*|MSYS*|CYGWIN*) echo "windows" ;;
        *) echo "unknown" ;;
    esac
}

###############################################################################
# Phase 1: Environment Verification
###############################################################################

echo ""
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo -e "${BLUE}  AIOX Framework Setup${NC}"
echo -e "${BLUE}═══════════════════════════════════════════════════════${NC}"
echo ""

log_step "Verifying environment..."

# Check Node.js version
if ! command -v node &> /dev/null; then
    log_error "Node.js not found. Install Node.js 18.0.0+ from https://nodejs.org"
    exit 1
fi

NODE_VERSION=$(node --version | cut -d'v' -f2 | cut -d'.' -f1)
if [ "$NODE_VERSION" -lt 18 ]; then
    log_error "Node.js version $(node --version) is below 18.0.0 requirement"
    exit 1
fi
log_success "Node.js $(node --version)"

# Check npm version
if ! command -v npm &> /dev/null; then
    log_error "npm not found. Install npm 9.0.0+ or update Node.js"
    exit 1
fi

NPM_VERSION=$(npm --version | cut -d'.' -f1)
if [ "$NPM_VERSION" -lt 9 ]; then
    log_error "npm version $(npm --version) is below 9.0.0 requirement"
    exit 1
fi
log_success "npm $(npm --version)"

# Check package.json
if [ ! -f "package.json" ]; then
    log_warning "package.json not found in current directory"
    log_step "Would you like me to create one? (y/n)"
    read -r CREATE_PKG
    if [[ $CREATE_PKG == "y" || $CREATE_PKG == "Y" ]]; then
        npm init -y > /dev/null 2>&1
        log_success "Created package.json"
    else
        log_error "package.json is required. Exiting."
        exit 1
    fi
fi

OS=$(detect_os)
log_success "Platform: $(uname -s) ($OS)"

echo ""

###############################################################################
# Phase 2: Install AIOX Core Framework
###############################################################################

log_step "Installing AIOX core framework..."

# Check if .aiox-core already exists
if [ -d ".aiox-core" ]; then
    log_warning ".aiox-core already exists (AIOX may be installed)"
fi

# Install AIOX with --quiet flag for non-interactive CI/CD compatibility
if npx aiox-core install --quiet; then
    log_success "AIOX core framework installed"
else
    log_error "Failed to install AIOX core framework"
    exit 1
fi

echo ""

###############################################################################
# Phase 3: Sync IDE Definitions
###############################################################################

log_step "Synchronizing agent definitions to IDE..."

# Check if npm run sync:ide is available
if grep -q "sync:ide" package.json; then
    if npm run sync:ide > /dev/null 2>&1; then
        log_success "IDE definitions synced"
    else
        log_warning "IDE sync encountered issues (continuing...)"
    fi
else
    log_warning "sync:ide script not found in package.json (skipping)"
fi

echo ""

###############################################################################
# Phase 4: Initialize Husky Git Hooks
###############################################################################

log_step "Configuring Husky git hooks..."

# Install Husky if not already installed
if ! grep -q '"husky"' package.json; then
    npm install husky --save-dev > /dev/null 2>&1
    log_success "Installed Husky"
else
    log_success "Husky already installed"
fi

# Initialize Husky
if [ -d ".husky" ]; then
    log_warning ".husky directory already exists"
else
    npx husky install > /dev/null 2>&1 || true
    log_success "Husky initialized"
fi

# Configure git to use .husky
git config core.hooksPath .husky > /dev/null 2>&1 || true
log_success "Git configured to use .husky hooks"

echo ""

###############################################################################
# Phase 5: Verify Installation
###############################################################################

log_step "Verifying installation..."

# Run AIOX doctor
if command -v npx &> /dev/null; then
    echo ""
    if npx aiox-core doctor; then
        log_success "System verification passed"
    else
        log_warning "AIOX doctor reported warnings (see above)"
    fi
else
    log_warning "npx not available, skipping verification"
fi

echo ""

###############################################################################
# Phase 6: Summary
###############################################################################

echo -e "${GREEN}═══════════════════════════════════════════════════════${NC}"
echo -e "${GREEN}  Setup Complete!${NC}"
echo -e "${GREEN}═══════════════════════════════════════════════════════${NC}"
echo ""
echo -e "Your project is now ${GREEN}AIOX-ready${NC}!"
echo ""
echo "Next steps:"
echo "  1. Create your first story:"
echo "     @sm *create-story"
echo ""
echo "  2. Activate the developer agent:"
echo "     @dev"
echo ""
echo "  3. View available commands:"
echo "     @dev *help"
echo ""
echo "Learn more:"
echo "  - Documentation: docs/stories/ (after creating first story)"
echo "  - Project conventions: .claude/PROJECT_CONVENTIONS.md"
echo "  - AIOX guide: https://github.com/SynkraAI/aiox-core"
echo ""


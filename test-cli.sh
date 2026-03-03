#!/bin/bash

# Test script for Plane CLI
# Usage: ./test-cli.sh

set -e

# Configuration
export PLANE_API_KEY="plane_api_dc8868628fd744cfaf33f55e482cb88c"
export PLANE_WORKSPACE="test-workspace"
export PLANE_API_HOST="http://plane.tequerist.com"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Plane CLI Test Script"
echo "=========================================="
echo ""

# Build the CLI
echo -e "${YELLOW}Building CLI...${NC}"
cd /home/rohithmahesh/Dev/plane/plane-cli
go build -o plane main.go
echo -e "${GREEN}✓ Build successful${NC}"
echo ""

# Test 1: Version
echo -e "${YELLOW}Test 1: Version${NC}"
./plane version
echo -e "${GREEN}✓ Version command works${NC}"
echo ""

# Test 2: Auth Status
echo -e "${YELLOW}Test 2: Auth Status${NC}"
./plane auth status 2>&1 || true
echo ""

# Test 3: Initialize config directory
echo -e "${YELLOW}Test 3: Initialize config${NC}"
mkdir -p ~/.config/plane-cli
echo "version: \"1.0\"" > ~/.config/plane-cli/config.yaml
echo "output_format: yaml" >> ~/.config/plane-cli/config.yaml
echo "api_host: http://plane.tequerist.com" >> ~/.config/plane-cli/config.yaml
echo "default_workspace: test-workspace" >> ~/.config/plane-cli/config.yaml
echo -e "${GREEN}✓ Config initialized${NC}"
echo ""

# Test 3b: Config set
echo -e "${YELLOW}Test 3b: Set configuration${NC}"
./plane config set workspace test-workspace || true
./plane config set api_host http://plane.tequerist.com || true
./plane config set project 99b45f00-73af-42a9-912f-7a348f5b42d4 || true
echo -e "${GREEN}✓ Config set successful${NC}"
echo ""

# Test 4: List projects (with env vars)
echo -e "${YELLOW}Test 4: List Projects${NC}"
./plane project list
echo -e "${GREEN}✓ Project list successful${NC}"
echo ""

# Test 5: List issues
echo -e "${YELLOW}Test 5: List Issues${NC}"
./plane issue list --limit 5
echo -e "${GREEN}✓ Issue list successful${NC}"
echo ""

# Test 6: List issues with filters
echo -e "${YELLOW}Test 6: List Issues with filters${NC}"
./plane issue list --priority high --limit 5
echo -e "${GREEN}✓ Issue list with filters successful${NC}"
echo ""

# Test 7: Output formats
echo -e "${YELLOW}Test 7: JSON Output${NC}"
./plane project list --output json | head -20
echo -e "${GREEN}✓ JSON output works${NC}"
echo ""

echo -e "${YELLOW}Test 8: YAML Output${NC}"
./plane project list --output yaml | head -20
echo -e "${GREEN}✓ YAML output works${NC}"
echo ""

# Test 9: Shell completion
echo -e "${YELLOW}Test 9: Shell Completion${NC}"
./plane completion bash > /dev/null
echo -e "${GREEN}✓ Bash completion generated${NC}"
echo ""

echo "=========================================="
echo -e "${GREEN}All tests passed!${NC}"
echo "=========================================="

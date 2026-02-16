#!/bin/bash
# Customization script for the Omada Terraform Provider
# This script replaces generic placeholders with your specific values

set -e

echo "╔════════════════════════════════════════════════════════════════════╗"
echo "║   TERRAFORM PROVIDER CUSTOMIZATION                                 ║"
echo "╚════════════════════════════════════════════════════════════════════╝"
echo ""

# Check if running from project root
if [ ! -f "go.mod" ]; then
    echo "❌ Error: Must run from project root directory"
    exit 1
fi

# Prompt for GitHub username/organization
echo "Enter your GitHub username or organization name:"
read -p "(e.g., mycompany, johnsmith): " USERNAME

if [ -z "$USERNAME" ]; then
    echo "❌ Error: Username cannot be empty"
    exit 1
fi

echo ""
echo "You entered: $USERNAME"
read -p "Is this correct? (y/n): " CONFIRM

if [ "$CONFIRM" != "y" ] && [ "$CONFIRM" != "Y" ]; then
    echo "❌ Cancelled"
    exit 1
fi

echo ""
echo "🔄 Updating files..."

# Update go.mod
echo "  • Updating go.mod..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s|your-org|$USERNAME|g" go.mod
else
    sed -i "s|your-org|$USERNAME|g" go.mod
fi

# Update all Go files
echo "  • Updating Go files..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    find . -name "*.go" ! -path "./.git/*" -exec sed -i '' "s|your-org|$USERNAME|g" {} +
else
    find . -name "*.go" ! -path "./.git/*" -exec sed -i "s|your-org|$USERNAME|g" {} +
fi

# Update all Markdown files
echo "  • Updating documentation..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    find . -name "*.md" ! -path "./.git/*" ! -path "./.claude/*" -exec sed -i '' "s|your-org|$USERNAME|g" {} +
else
    find . -name "*.md" ! -path "./.git/*" ! -path "./.claude/*" -exec sed -i "s|your-org|$USERNAME|g" {} +
fi

# Update all Terraform files
echo "  • Updating examples..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    find . -name "*.tf" ! -path "./.git/*" -exec sed -i '' "s|your-org|$USERNAME|g" {} +
else
    find . -name "*.tf" ! -path "./.git/*" -exec sed -i "s|your-org|$USERNAME|g" {} +
fi

# Update workflows
echo "  • Updating workflows..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    find .github -name "*.yml" -exec sed -i '' "s|your-org|$USERNAME|g" {} +
else
    find .github -name "*.yml" -exec sed -i "s|your-org|$USERNAME|g" {} +
fi

# Update Makefile
echo "  • Updating Makefile..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s|your-org|$USERNAME|g" Makefile
else
    sed -i "s|your-org|$USERNAME|g" Makefile
fi

echo ""
echo "🔧 Updating dependencies..."
go mod tidy

echo ""
echo "🏗️  Testing build..."
if go build -o terraform-provider-omada; then
    echo "✅ Build successful!"
    rm -f terraform-provider-omada
else
    echo "❌ Build failed - please check errors above"
    exit 1
fi

echo ""
echo "🔍 Verifying customization..."
REMAINING=$(grep -r "your-org" --include="*.go" --include="*.md" --include="*.tf" --include="*.yml" . 2>/dev/null | grep -v ".git" | wc -l)

if [ "$REMAINING" -gt 0 ]; then
    echo "⚠️  Warning: Found $REMAINING remaining 'your-org' references:"
    grep -r "your-org" --include="*.go" --include="*.md" --include="*.tf" --include="*.yml" . 2>/dev/null | grep -v ".git" | head -5
    echo ""
    echo "You may need to manually update these files."
else
    echo "✅ No remaining placeholders found!"
fi

echo ""
echo "╔════════════════════════════════════════════════════════════════════╗"
echo "║   CUSTOMIZATION COMPLETE                                           ║"
echo "╚════════════════════════════════════════════════════════════════════╝"
echo ""
echo "📋 Next Steps:"
echo "  1. Review changes: git diff"
echo "  2. Run verification: ./scripts/verify-release.sh"
echo "  3. Commit changes: git add . && git commit -m 'Customize for $USERNAME'"
echo "  4. Follow PUBLISHING.md to publish your provider"
echo ""
echo "✨ Your provider is now customized for: $USERNAME"

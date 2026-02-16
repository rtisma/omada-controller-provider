#!/bin/bash
# Verification script to check if the provider is ready for Terraform Registry publication

set -e

echo "╔════════════════════════════════════════════════════════════════════╗"
echo "║   TERRAFORM REGISTRY PUBLICATION VERIFICATION                      ║"
echo "╚════════════════════════════════════════════════════════════════════╝"
echo ""

ERRORS=0
WARNINGS=0

# Function to print check results
check_pass() {
    echo "✅ $1"
}

check_fail() {
    echo "❌ $1"
    ((ERRORS++))
}

check_warn() {
    echo "⚠️  $1"
    ((WARNINGS++))
}

echo "📋 Checking Repository Structure..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check required files
[ -f "LICENSE" ] && check_pass "LICENSE file exists" || check_fail "LICENSE file missing"
[ -f "README.md" ] && check_pass "README.md exists" || check_fail "README.md missing"
[ -f "go.mod" ] && check_pass "go.mod exists" || check_fail "go.mod missing"
[ -f "main.go" ] && check_pass "main.go exists" || check_fail "main.go missing"
[ -f ".goreleaser.yml" ] && check_pass ".goreleaser.yml exists" || check_fail ".goreleaser.yml missing"
[ -f "terraform-registry-manifest.json" ] && check_pass "Registry manifest exists" || check_fail "Registry manifest missing"

echo ""
echo "📖 Checking Documentation..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check docs structure
[ -d "docs" ] && check_pass "docs/ directory exists" || check_fail "docs/ directory missing"
[ -f "docs/index.md" ] && check_pass "docs/index.md exists" || check_fail "docs/index.md missing"
[ -d "docs/resources" ] && check_pass "docs/resources/ exists" || check_fail "docs/resources/ missing"
[ -d "docs/data-sources" ] && check_pass "docs/data-sources/ exists" || check_fail "docs/data-sources/ missing"

# Count documentation files
RESOURCE_DOCS=$(find docs/resources -name "*.md" 2>/dev/null | wc -l)
DATASOURCE_DOCS=$(find docs/data-sources -name "*.md" 2>/dev/null | wc -l)

echo "   Found $RESOURCE_DOCS resource docs"
echo "   Found $DATASOURCE_DOCS data source docs"

if [ $RESOURCE_DOCS -gt 0 ]; then
    check_pass "Resource documentation exists ($RESOURCE_DOCS files)"
else
    check_warn "No resource documentation found"
fi

if [ $DATASOURCE_DOCS -gt 0 ]; then
    check_pass "Data source documentation exists ($DATASOURCE_DOCS files)"
else
    check_warn "No data source documentation found"
fi

echo ""
echo "📚 Checking Examples..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

[ -d "examples" ] && check_pass "examples/ directory exists" || check_fail "examples/ directory missing"
[ -d "examples/provider" ] && check_pass "examples/provider/ exists" || check_fail "examples/provider/ missing"
[ -d "examples/resources" ] && check_pass "examples/resources/ exists" || check_warn "examples/resources/ missing"
[ -d "examples/data-sources" ] && check_pass "examples/data-sources/ exists" || check_warn "examples/data-sources/ missing"

# Count example files
EXAMPLE_COUNT=$(find examples -name "*.tf" 2>/dev/null | wc -l)
echo "   Found $EXAMPLE_COUNT .tf example files"

if [ $EXAMPLE_COUNT -gt 0 ]; then
    check_pass "Example files exist ($EXAMPLE_COUNT files)"
else
    check_warn "No example .tf files found"
fi

echo ""
echo "🔧 Checking GitHub Actions..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

[ -d ".github/workflows" ] && check_pass ".github/workflows/ exists" || check_fail ".github/workflows/ missing"
[ -f ".github/workflows/release.yml" ] && check_pass "Release workflow exists" || check_fail "Release workflow missing"
[ -f ".github/workflows/test.yml" ] && check_pass "Test workflow exists" || check_fail "Test workflow missing"

echo ""
echo "🏗️  Checking Build..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if go build -o terraform-provider-omada 2>/dev/null; then
    check_pass "Provider builds successfully"
    rm -f terraform-provider-omada
else
    check_fail "Provider build failed"
fi

echo ""
echo "📝 Checking Documentation Format..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Check if docs have proper frontmatter
if [ -f "docs/index.md" ]; then
    if head -1 docs/index.md | grep -q "^---$"; then
        check_pass "docs/index.md has frontmatter"
    else
        check_warn "docs/index.md missing frontmatter"
    fi
fi

# Check resource docs
if [ -d "docs/resources" ]; then
    for doc in docs/resources/*.md; do
        if [ -f "$doc" ]; then
            if head -1 "$doc" | grep -q "^---$"; then
                check_pass "$(basename "$doc") has frontmatter"
            else
                check_warn "$(basename "$doc") missing frontmatter"
            fi
        fi
    done
fi

echo ""
echo "🔍 Checking Naming Convention..."
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if grep -q "module github.com/.*/terraform-provider-" go.mod; then
    check_pass "Repository follows naming convention"
else
    check_fail "Repository name should be terraform-provider-{name}"
fi

echo ""
echo "📊 Summary"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ $ERRORS -eq 0 ] && [ $WARNINGS -eq 0 ]; then
    echo "🎉 All checks passed! Provider is ready for publication."
    echo ""
    echo "Next steps:"
    echo "  1. Ensure GitHub repository is public"
    echo "  2. Configure GPG key in GitHub Secrets"
    echo "  3. Register provider at https://registry.terraform.io/publish/provider"
    echo "  4. Create and push a version tag (e.g., git tag v0.1.0 && git push origin v0.1.0)"
    exit 0
elif [ $ERRORS -eq 0 ]; then
    echo "⚠️  $WARNINGS warning(s) found. Provider can be published but consider fixing warnings."
    echo ""
    echo "To publish:"
    echo "  1. Review warnings above"
    echo "  2. Ensure GitHub repository is public"
    echo "  3. Configure GPG key in GitHub Secrets"
    echo "  4. Register provider at https://registry.terraform.io/publish/provider"
    echo "  5. Create and push a version tag"
    exit 0
else
    echo "❌ $ERRORS error(s) and $WARNINGS warning(s) found. Fix errors before publishing."
    echo ""
    echo "Common fixes:"
    echo "  - Ensure all required files exist"
    echo "  - Check documentation structure"
    echo "  - Verify provider builds successfully"
    echo "  - Add frontmatter to documentation files"
    exit 1
fi

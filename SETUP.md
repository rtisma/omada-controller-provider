# Provider Setup Guide

Before publishing your provider, you need to customize it for your organization.

## Required Customizations

### 1. Update Organization/Namespace

Replace `your-org` with your actual GitHub organization or username in these files:

**Go Module:**
```bash
# In go.mod
module github.com/YOUR-USERNAME/terraform-provider-omada
```

**Main Provider:**
```go
// In main.go
opts := providerserver.ServeOpts{
    Address: "registry.terraform.io/YOUR-USERNAME/omada",
    Debug:   debug,
}
```

**Makefile:**
```makefile
# In Makefile
NAMESPACE=YOUR-USERNAME
```

**All Import Statements:**
```bash
# Update all Go files
find . -name "*.go" -exec sed -i 's|your-org|YOUR-USERNAME|g' {} +
```

**Documentation & Examples:**
```bash
# Update all docs and examples
find . -name "*.md" -name "*.tf" -exec sed -i 's|your-org|YOUR-USERNAME|g' {} +
```

### 2. Update URLs in Documentation

Replace repository URLs in:
- `README.md`
- `PUBLISHING.md`
- `REGISTRY_READY.md`
- `CONTRIBUTING.md`
- `docs/index.md`

Change:
```
https://github.com/your-org/terraform-provider-omada
```

To:
```
https://github.com/YOUR-USERNAME/terraform-provider-omada
```

### 3. Verify All References

Run this command to find remaining placeholders:

```bash
grep -r "your-org" --include="*.go" --include="*.md" --include="*.tf" --include="*.yml" .
```

Should return no results after customization.

## Quick Setup Script

```bash
#!/bin/bash
# Replace YOUR_USERNAME with your GitHub username/organization

USERNAME="YOUR_USERNAME"

# Update go.mod
sed -i "s|your-org|$USERNAME|g" go.mod

# Update all Go files
find . -name "*.go" -exec sed -i "s|your-org|$USERNAME|g" {} +

# Update all documentation
find . -name "*.md" -exec sed -i "s|your-org|$USERNAME|g" {} +

# Update all examples
find . -name "*.tf" -exec sed -i "s|your-org|$USERNAME|g" {} +

# Update workflows
find .github -name "*.yml" -exec sed -i "s|your-org|$USERNAME|g" {} +

# Update Makefile
sed -i "s|your-org|$USERNAME|g" Makefile

# Tidy dependencies
go mod tidy

# Verify build
go build -o terraform-provider-omada

echo "✅ Setup complete! Verify with: grep -r 'your-org' ."
```

## Verification Checklist

After customization, verify:

- [ ] `go.mod` has correct module path
- [ ] `main.go` has correct registry address
- [ ] `Makefile` has correct namespace
- [ ] All Go imports updated
- [ ] All documentation updated
- [ ] All examples updated
- [ ] Provider builds: `go build`
- [ ] Tests pass: `go test ./...`
- [ ] No "your-org" remains: `grep -r "your-org" .`

## Example Values

**If your GitHub username is `johnsmith`:**

```
your-org → johnsmith
github.com/your-org → github.com/johnsmith
registry.terraform.io/your-org → registry.terraform.io/johnsmith
```

**If your organization is `acme-corp`:**

```
your-org → acme-corp
github.com/your-org → github.com/acme-corp
registry.terraform.io/your-org → registry.terraform.io/acme-corp
```

## After Customization

1. Run `go mod tidy`
2. Run `go build`
3. Run `./scripts/verify-release.sh`
4. Commit changes
5. Follow `PUBLISHING.md` for publication steps

## Need Help?

See `PUBLISHING.md` for detailed publication instructions.

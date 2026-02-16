# Terraform Registry Publication Checklist

Use this checklist when publishing your provider to the Terraform Registry.

## ☐ Pre-Publication Checklist

### Repository Setup
- [ ] Repository is public on GitHub
- [ ] Repository name follows convention: `terraform-provider-{NAME}`
- [ ] Repository has description and topics configured
- [ ] Default branch is `main` or `master`

### Required Files
- [ ] `LICENSE` file present (MPL 2.0)
- [ ] `README.md` with provider overview
- [ ] `go.mod` with correct module path
- [ ] `main.go` with provider entry point
- [ ] `.goreleaser.yml` configured
- [ ] `terraform-registry-manifest.json` present

### Documentation
- [ ] `docs/index.md` exists with provider documentation
- [ ] `docs/index.md` has proper YAML frontmatter
- [ ] All resource docs in `docs/resources/` have frontmatter
- [ ] All data source docs in `docs/data-sources/` have frontmatter
- [ ] Documentation follows Terraform Registry format

### Examples
- [ ] `examples/` directory exists
- [ ] `examples/provider/` has provider configuration
- [ ] Example .tf files for all resources
- [ ] Example .tf files for all data sources
- [ ] Examples are tested and work

### GitHub Actions
- [ ] `.github/workflows/release.yml` configured
- [ ] `.github/workflows/test.yml` configured
- [ ] Workflows have proper permissions

### Testing
- [ ] Provider builds successfully: `make build`
- [ ] Unit tests pass: `make test`
- [ ] No syntax errors in Go code
- [ ] Terraform examples validate

### Verification
- [ ] Run `./scripts/verify-release.sh` - all checks pass
- [ ] No errors in verification output
- [ ] Documentation format verified

## ☐ GPG Key Setup

### Generate Key (if needed)
- [ ] GPG key generated with 4096 bits
- [ ] Key has no expiration (or appropriate expiration)
- [ ] Key email matches your GitHub email

### Export Keys
- [ ] Public key exported: `gpg --armor --export KEY_ID`
- [ ] Private key exported: `gpg --armor --export-secret-keys KEY_ID`
- [ ] Passphrase documented (if set)

### GitHub Configuration
- [ ] Added `GPG_PRIVATE_KEY` to GitHub Secrets
- [ ] Added `PASSPHRASE` to GitHub Secrets (if key has passphrase)
- [ ] Secrets are in the correct repository

## ☐ Terraform Registry Setup

### Registry Account
- [ ] Signed up at https://registry.terraform.io/
- [ ] Logged in with GitHub account
- [ ] Accepted Terms of Service
- [ ] Profile information complete

### Provider Registration
- [ ] Visited https://registry.terraform.io/publish/provider
- [ ] Selected correct GitHub repository
- [ ] Provider successfully registered
- [ ] Webhook configured (automatic)

### Public Key Upload
- [ ] Uploaded GPG public key to Registry
- [ ] Key fingerprint matches local key
- [ ] Key verified in Registry

## ☐ First Release

### Pre-Release
- [ ] All changes committed and pushed
- [ ] CHANGELOG.md updated with v0.1.0
- [ ] README.md reviewed and complete
- [ ] Documentation reviewed
- [ ] Examples tested
- [ ] Version number decided (e.g., v0.1.0)

### Create Release
- [ ] Created tag: `git tag -a v0.1.0 -m "First release"`
- [ ] Pushed tag: `git push origin v0.1.0`
- [ ] GitHub Actions workflow started
- [ ] Workflow completed successfully

### Verify Release
- [ ] GitHub release created
- [ ] Binaries present in release
- [ ] Checksums file present
- [ ] GPG signature present
- [ ] Release notes added

## ☐ Post-Publication

### Registry Verification
- [ ] Provider appears on Terraform Registry
- [ ] Documentation visible on Registry
- [ ] Examples visible on Registry
- [ ] Version badge shows correct version

### Test Installation
- [ ] Created test Terraform configuration
- [ ] Ran `terraform init` - provider downloads
- [ ] Provider installs from Registry
- [ ] No errors during installation

### Documentation Check
- [ ] All resources documented on Registry
- [ ] All data sources documented on Registry
- [ ] Examples render correctly
- [ ] Links work correctly

## ☐ Ongoing Maintenance

### For Each Release
- [ ] Update CHANGELOG.md
- [ ] Increment version following semver
- [ ] Create and push tag
- [ ] Verify GitHub Actions success
- [ ] Check Registry updates
- [ ] Test installation

### Community
- [ ] Monitor GitHub issues
- [ ] Respond to questions
- [ ] Review pull requests
- [ ] Update documentation as needed

---

## Quick Commands

```bash
# Verify readiness
./scripts/verify-release.sh

# Build
make build

# Run tests
make test

# Create release
git tag -a v0.1.0 -m "Release v0.1.0"
git push origin v0.1.0

# Check release
make release-check VERSION=v0.1.0
```

## Useful Links

- Registry: https://registry.terraform.io/
- Publish: https://registry.terraform.io/publish/provider
- Docs: https://developer.hashicorp.com/terraform/registry/providers/publishing
- GitHub Actions: https://github.com/[YOUR_ORG]/terraform-provider-omada/actions

---

**Status:** Ready for publication ✅
**Last Updated:** $(date +%Y-%m-%d)

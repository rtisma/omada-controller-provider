# Publishing to Terraform Registry

This guide walks through the steps to publish the Omada provider to the Terraform Registry.

## Prerequisites

Before publishing, ensure you have:

1. ✅ A public GitHub repository under your organization/user account
2. ✅ GPG key for signing releases
3. ✅ GitHub account with appropriate permissions
4. ✅ All tests passing
5. ✅ Complete documentation

## Step 1: Prepare Your GPG Key

### Generate a GPG Key (if you don't have one)

```bash
gpg --full-generate-key
```

Choose:
- RSA and RSA (default)
- 4096 bits
- No expiration (or your preference)
- Your real name and email

### Export Your GPG Key

```bash
# List your keys to find the key ID
gpg --list-secret-keys --keyid-format=long

# Export the public key (for Terraform Registry)
gpg --armor --export YOUR_KEY_ID > gpg-public-key.asc

# Export the private key (for GitHub Actions)
gpg --armor --export-secret-keys YOUR_KEY_ID
```

## Step 2: Configure GitHub Repository

### Add GitHub Secrets

Go to your repository Settings → Secrets and variables → Actions and add:

1. **GPG_PRIVATE_KEY**: Your private GPG key (the output from the export command)
2. **PASSPHRASE**: Your GPG key passphrase (if you set one)

### Verify GitHub Actions

Ensure these workflow files exist:
- `.github/workflows/release.yml` - For creating releases
- `.github/workflows/test.yml` - For running tests

## Step 3: Sign Up for Terraform Registry

1. Go to https://registry.terraform.io/
2. Sign in with your GitHub account
3. Accept the Terms of Service

## Step 4: Publish Your Provider

### Register the Provider

1. Go to https://registry.terraform.io/publish/provider
2. Select your GitHub repository: `your-org/terraform-provider-omada`
3. Click "Publish Provider"

The Terraform Registry will:
- Verify your repository structure
- Check for required files
- Set up webhooks for automatic updates

### Required Repository Structure

Your repository must have:

```
terraform-provider-omada/
├── .github/
│   └── workflows/
│       ├── release.yml         # ✅ Release automation
│       └── test.yml            # ✅ Test automation
├── .goreleaser.yml             # ✅ GoReleaser config
├── docs/                       # ✅ Documentation
│   ├── index.md
│   ├── resources/
│   └── data-sources/
├── examples/                   # ✅ Usage examples
│   ├── provider/
│   ├── resources/
│   └── data-sources/
├── internal/                   # ✅ Provider code
│   ├── client/
│   └── provider/
├── go.mod                      # ✅ Go module
├── LICENSE                     # ✅ License file
├── main.go                     # ✅ Entry point
├── terraform-registry-manifest.json  # ✅ Registry manifest
└── README.md                   # ✅ Repository README
```

All of these are already in place! ✅

## Step 5: Create Your First Release

### Tag a Release

```bash
# Ensure you're on main branch with latest changes
git checkout main
git pull

# Create and push a tag (must start with 'v')
git tag v0.1.0
git push origin v0.1.0
```

### What Happens Next

1. GitHub Actions triggers the release workflow
2. GoReleaser builds binaries for all platforms
3. Binaries are signed with your GPG key
4. Release is created on GitHub
5. Terraform Registry is notified via webhook
6. Documentation is automatically published

### Monitor the Release

1. Check GitHub Actions: `https://github.com/your-org/terraform-provider-omada/actions`
2. Verify the release: `https://github.com/your-org/terraform-provider-omada/releases`
3. Check the Registry: `https://registry.terraform.io/providers/your-org/omada`

## Step 6: Verify the Published Provider

After a few minutes, your provider should appear on the Registry.

Test it works:

```hcl
terraform {
  required_providers {
    omada = {
      source  = "your-org/omada"
      version = "~> 0.1"
    }
  }
}

provider "omada" {
  host     = "https://your-controller:8043"
  username = "admin"
  password = var.password
  site_id  = "Default"
  insecure = true
}
```

```bash
terraform init
# Should download the provider from the registry
```

## Ongoing Releases

For subsequent releases:

1. Make your changes and commit them
2. Update `CHANGELOG.md` with the new version
3. Create and push a new tag:
   ```bash
   git tag v0.2.0
   git push origin v0.2.0
   ```
4. GitHub Actions automatically handles the rest

## Semantic Versioning

Follow semantic versioning (https://semver.org/):

- **v0.x.y**: Initial development (breaking changes allowed)
- **v1.0.0**: First stable release
- **MAJOR**: Incompatible API changes
- **MINOR**: New functionality (backwards compatible)
- **PATCH**: Bug fixes (backwards compatible)

Examples:
- v0.1.0 → v0.2.0: New resources added
- v0.2.0 → v0.2.1: Bug fix
- v0.9.9 → v1.0.0: First stable release
- v1.0.0 → v2.0.0: Breaking changes

## Troubleshooting

### Release Fails

Check:
1. GPG key is properly configured in GitHub Secrets
2. GITHUB_TOKEN has sufficient permissions
3. All builds pass in the test workflow
4. Tag name starts with 'v' (e.g., v0.1.0, not 0.1.0)

### Documentation Doesn't Appear

Verify:
1. `docs/` directory structure is correct
2. All `.md` files have proper frontmatter
3. Examples are in the `examples/` directory
4. Wait 10-15 minutes for Registry to update

### Provider Not Found

Ensure:
1. Repository is public
2. Provider was successfully published to Registry
3. Using correct source: `your-org/omada`
4. Tag was successfully created and released

## Registry Requirements Checklist

Before your first release:

- ✅ Repository is public on GitHub
- ✅ Repository follows naming convention: `terraform-provider-{NAME}`
- ✅ LICENSE file present (MPL 2.0)
- ✅ Documentation in `docs/` directory
- ✅ Examples in `examples/` directory
- ✅ `.goreleaser.yml` configured
- ✅ GitHub Actions workflows configured
- ✅ GPG key configured in GitHub Secrets
- ✅ `terraform-registry-manifest.json` present
- ✅ Version tag starts with 'v'

All requirements are met! ✅

## Additional Resources

- [Terraform Registry - Publishing Providers](https://developer.hashicorp.com/terraform/registry/providers/publishing)
- [Provider Requirements](https://developer.hashicorp.com/terraform/registry/providers/requirements)
- [GoReleaser Documentation](https://goreleaser.com/)
- [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)

## Support

If you encounter issues:

1. Check [Terraform Registry Documentation](https://developer.hashicorp.com/terraform/registry)
2. Review [GitHub Actions logs](https://github.com/your-org/terraform-provider-omada/actions)
3. Consult [GoReleaser Troubleshooting](https://goreleaser.com/errors/)
4. Open an issue in the repository

## First Release Checklist

Before running `git tag v0.1.0`:

- [ ] All code is committed and pushed
- [ ] Tests are passing (`make test`)
- [ ] Provider builds successfully (`make build`)
- [ ] Documentation is complete
- [ ] CHANGELOG.md is updated
- [ ] Examples are working and tested
- [ ] GPG key is configured in GitHub
- [ ] Repository is public
- [ ] Provider is registered on Terraform Registry

Ready to release? Run:

```bash
git tag -a v0.1.0 -m "First release"
git push origin v0.1.0
```

Then watch the magic happen! 🚀

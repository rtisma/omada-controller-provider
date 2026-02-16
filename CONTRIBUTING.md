# Contributing to Terraform Provider for Omada Controller

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to the project.

## Development Environment Setup

### Prerequisites

- [Go](https://golang.org/doc/install) >= 1.24
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- Access to an Omada Controller 5.x instance for testing

### Clone and Build

```bash
git clone https://github.com/your-org/terraform-provider-omada.git
cd terraform-provider-omada
make deps
make build
```

### Install Locally for Testing

```bash
make install
```

This installs the provider to your local Terraform plugins directory.

## Making Changes

### Code Style

- Follow standard Go conventions and idioms
- Run `make fmt` before committing to format code
- Run `make lint` to check for common issues
- Keep functions focused and well-documented

### Testing

1. **Unit Tests**: Test individual functions and components
   ```bash
   make test
   ```

2. **Acceptance Tests**: Test against a real Omada Controller
   ```bash
   # Set up test environment variables
   export TF_ACC=1
   export OMADA_HOST="https://your-controller:8043"
   export OMADA_USERNAME="admin"
   export OMADA_PASSWORD="password"
   export OMADA_SITE="Default"

   make testacc
   ```

3. **Manual Testing**: Use the examples in the `examples/` directory

### Adding New Resources

1. Create the API client methods in `internal/client/`
2. Create the resource implementation in `internal/provider/`
3. Register the resource in `internal/provider/provider.go`
4. Add examples in `examples/resources/<resource_name>/`
5. Update the README.md with the new resource
6. Add tests for the resource

### Adding New Data Sources

Follow the same pattern as resources, but use data source interfaces.

## Pull Request Process

1. **Fork the repository** and create a new branch for your feature or bugfix
   ```bash
   git checkout -b feature/my-new-feature
   ```

2. **Make your changes** following the guidelines above

3. **Test your changes** thoroughly
   ```bash
   make test
   make build
   ```

4. **Update documentation**:
   - Update README.md if adding new features
   - Add or update examples
   - Update CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/) format

5. **Commit your changes** with clear, descriptive commit messages
   ```bash
   git commit -m "Add support for port forwarding rules"
   ```

6. **Push to your fork** and submit a pull request
   ```bash
   git push origin feature/my-new-feature
   ```

7. **Fill out the PR template** with:
   - Description of the changes
   - Related issue numbers (if applicable)
   - Testing performed
   - Screenshots (if UI-related)

## API Documentation

Since TP-Link's official API documentation for Omada Controller is limited:

- Reference the [unofficial API examples](https://gist.github.com/mbentley/03c198077c81d52cb029b825e9a6dc18)
- Use browser developer tools to inspect API calls
- Document any new endpoints or behaviors you discover
- Be cautious about breaking changes between Omada Controller versions

## Debugging

### Enable Debug Logging

```bash
export TF_LOG=DEBUG
export TF_LOG_PATH=./terraform.log
terraform plan
```

### Run Provider in Debug Mode

```bash
go run main.go -debug
```

This will output instructions for attaching Terraform to the debug session.

## Coding Guidelines

### Error Handling

- Always return descriptive errors
- Use `fmt.Errorf` for error wrapping
- Provide context in error messages to help with debugging

### API Client

- Keep the API client independent of Terraform logic
- Handle authentication and session management transparently
- Implement retry logic for transient failures
- Parse and return structured error responses

### Provider Resources

- Use computed fields for read-only attributes
- Implement proper state refresh in Read()
- Handle resources deleted outside Terraform gracefully
- Support resource import where appropriate

## Release Process

Releases are managed by project maintainers:

1. Update CHANGELOG.md with release notes
2. Create a new tag following semantic versioning
3. GitHub Actions will build and publish the release

## Getting Help

- Open an issue for bugs or feature requests
- Join discussions in existing issues
- Ask questions in pull requests

## Code of Conduct

Be respectful and professional in all interactions. We're all here to build something useful together.

## License

By contributing, you agree that your contributions will be licensed under the Mozilla Public License 2.0.

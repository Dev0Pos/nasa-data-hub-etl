# Contributing to NASA Data Hub ETL

Thank you for your interest in contributing to NASA Data Hub ETL! 🚀

## 🚀 Quick Start

### Prerequisites
- Go 1.25+
- Docker
- Git

### Setup Development Environment

```bash
# Clone the repository
git clone https://github.com/your-username/nasa-data-hub-etl.git
cd nasa-data-hub-etl

# Install dependencies
make deps

# Set up environment variables
cp env.example .env
# Edit .env with your configuration

# Run tests
make test

# Build the application
make build
```

## 🔄 Development Workflow

### 1. Create a Branch

```bash
# Create and switch to a new branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/your-bug-description
```

### 2. Make Changes

- Write your code following Go best practices
- Add tests for new functionality
- Update documentation if needed
- Ensure all tests pass: `make test`

### 3. Commit Changes

```bash
# Add your changes
git add .

# Commit with descriptive message
git commit -m "feat: add new feature description"

# Push to your fork
git push origin feature/your-feature-name
```

### 4. Create Pull Request

1. Go to GitHub and create a Pull Request
2. Fill out the PR template
3. Request review from maintainers
4. Address any feedback

## 📝 Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### Examples:
```
feat(api): add retry logic for NASA EONET API calls
fix(database): resolve connection pool leak
docs: update README with new configuration options
test(etl): add unit tests for data transformation
```

## 🧪 Testing

### Run Tests
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run specific test
go test ./internal/etl/...
```

### Test Requirements
- All new code must have tests
- Maintain or improve test coverage
- Tests should be fast and reliable

## 🐳 Docker Development

### Build and Test Docker Image
```bash
# Build Docker image
make docker-build

# Run container locally
make docker-run

# Test with environment variables
export DATABASE_HOST="localhost"
export DATABASE_PASSWORD="test-password"
make docker-run
```

## 📋 Code Review Process

### Before Submitting PR:
- [ ] Code follows Go best practices
- [ ] All tests pass
- [ ] Code is properly formatted (`make fmt`)
- [ ] No linting errors (`make lint`)
- [ ] Documentation is updated
- [ ] Commit messages follow convention

### Review Checklist:
- [ ] Code quality and style
- [ ] Test coverage
- [ ] Security considerations
- [ ] Performance impact
- [ ] Documentation accuracy

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment:**
   - Go version
   - OS and version
   - Docker version (if applicable)

2. **Steps to Reproduce:**
   - Clear, numbered steps
   - Expected vs actual behavior

3. **Logs:**
   - Relevant error messages
   - Application logs

4. **Additional Context:**
   - Configuration files (remove sensitive data)
   - Environment variables (remove sensitive data)

## 💡 Feature Requests

When requesting features:

1. **Describe the feature:**
   - What problem does it solve?
   - How should it work?

2. **Provide context:**
   - Use cases
   - Expected behavior
   - Any relevant examples

3. **Consider implementation:**
   - Breaking changes
   - Backward compatibility
   - Performance impact

## 🔒 Security

### Reporting Security Issues

**DO NOT** create public issues for security vulnerabilities.

Instead:
1. Email security concerns to: [your-email@domain.com]
2. Include detailed description
3. Provide steps to reproduce
4. We'll respond within 48 hours

### Security Best Practices

- Never commit secrets or sensitive data
- Use environment variables for configuration
- Follow principle of least privilege
- Keep dependencies updated

## 📚 Documentation

### Code Documentation
- Use Go doc comments for public APIs
- Include examples for complex functions
- Document configuration options

### README Updates
- Update README.md for user-facing changes
- Include new configuration options
- Update examples and usage

## 🏷️ Release Process

### Versioning
We follow [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Schedule
- **Major releases**: As needed
- **Minor releases**: Monthly
- **Patch releases**: As needed

## 🤝 Community Guidelines

### Be Respectful
- Use welcoming and inclusive language
- Be respectful of differing viewpoints
- Accept constructive criticism gracefully

### Be Collaborative
- Help others when possible
- Share knowledge and best practices
- Contribute to discussions constructively

### Be Professional
- Stay on topic
- Provide clear, actionable feedback
- Follow the project's code of conduct

## 📞 Getting Help

- **Documentation**: Check README.md and code comments
- **Issues**: Search existing issues before creating new ones
- **Discussions**: Use GitHub Discussions for questions
- **Email**: [your-email@domain.com] for private matters

## 🎉 Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes
- Project documentation

Thank you for contributing to NASA Data Hub ETL! 🌟

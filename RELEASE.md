# Release Process

## 🚀 How to Release NASA Data Hub ETL

### 1. Prepare for Release

```bash
# Ensure you're on main branch
git checkout main
git pull origin main

# Run tests locally
make test

# Build and test Docker image
make docker-build
```

### 2. Update Version

Update version in:
- `go.mod` (if needed)
- `README.md` (version references)
- `Makefile` (VERSION variable)

### 3. Create Release

#### Option A: GitHub Web Interface (Recommended)
1. Go to your GitHub repository
2. Click "Releases" → "Create a new release"
3. Choose a tag version (e.g., `v1.0.0`)
4. Set release title: `NASA Data Hub ETL v1.0.0`
5. Add release notes (see template below)
6. Click "Publish release"

#### Option B: GitHub CLI
```bash
# Install GitHub CLI if not installed
# https://cli.github.com/

# Create release
gh release create v1.0.0 \
  --title "NASA Data Hub ETL v1.0.0" \
  --notes-file RELEASE_NOTES.md \
  --latest
```

### 4. Automated Process

Once you create a release, GitHub Actions will automatically:
- ✅ Run all tests
- ✅ Build Docker image
- ✅ Push to GitHub Container Registry
- ✅ Run security scans
- ✅ Create release artifacts

### 5. Verify Release

Check that:
- [ ] Docker image is available: `ghcr.io/your-username/nasa-data-hub-etl:v1.0.0`
- [ ] Release notes are published
- [ ] All CI/CD checks passed

## 📝 Release Notes Template

```markdown
## 🚀 NASA Data Hub ETL v1.0.0

### ✨ Features
- High-performance ETL pipeline for NASA EONET data
- VerticaDB integration with connection pooling
- Docker containerization with security best practices
- Health checks and monitoring endpoints
- Environment variable configuration
- Clean architecture with proper separation of concerns

### 🔧 Technical Details
- **Go Version:** 1.25
- **Base Image:** Distroless (non-root user)
- **Architecture:** Multi-stage Docker build
- **Security:** No hardcoded secrets, environment variables only

### 🐳 Docker Usage

```bash
# Pull the image
docker pull ghcr.io/your-username/nasa-data-hub-etl:v1.0.0

# Run with environment variables
docker run --rm \
  -e DATABASE_HOST="your-database-host" \
  -e DATABASE_PASSWORD="your-password" \
  -p 8080:8080 \
  ghcr.io/your-username/nasa-data-hub-etl:v1.0.0
```

### 📋 Requirements
- Go 1.25+
- Docker
- VerticaDB access
- NASA EONET API access

### 🔗 Links
- [Documentation](README.md)
- [Docker Image](https://github.com/your-username/nasa-data-hub-etl/pkgs/container/nasa-data-hub-etl)
- [Source Code](https://github.com/your-username/nasa-data-hub-etl)
```

## 🏷️ Versioning Strategy

We follow [Semantic Versioning](https://semver.org/):
- **MAJOR** (1.0.0): Breaking changes
- **MINOR** (0.1.0): New features, backward compatible
- **PATCH** (0.0.1): Bug fixes, backward compatible

## 🔄 Release Schedule

- **Major releases:** As needed for breaking changes
- **Minor releases:** Monthly or when significant features are added
- **Patch releases:** Weekly or as needed for bug fixes

## 🚨 Rollback Process

If a release has issues:

1. **Mark release as pre-release:**
   - Go to GitHub releases
   - Edit the release
   - Check "Set as a pre-release"

2. **Create hotfix:**
   ```bash
   git checkout main
   git pull origin main
   # Make fixes
   git commit -m "fix: critical issue"
   git push origin main
   ```

3. **Create patch release:**
   - Follow normal release process
   - Use next patch version (e.g., v1.0.1)

# Release Process for agent-def

This document outlines the process for creating and publishing new releases of agent-def.

## Release Workflow Overview

1. Update version and changelog  
2. Create and push a new git tag  
3. GitHub Actions automatically builds binaries and creates a release  
4. Verify the release artifacts  

## Detailed Steps

### 1. Update Version and Changelog

1. Ensure your local main branch is up to date:
   ```bash
   git checkout main
   git pull origin main
   ```

2. Update `CHANGELOG.md` with the new version and changes:
   - Add a new section at the top for the new version  
   - Follow the [Keep a Changelog](https://keepachangelog.com/) format  
   - Group changes into Added, Changed, Deprecated, Removed, Fixed, and Security  
   - Reference PRs or issues where appropriate  
   - Set the release date  

3. Commit the changes:
   ```bash
   git add CHANGELOG.md
   git commit -m "Prepare for release vX.Y.Z"
   git push origin main
   ```

### 2. Create and Push a Git Tag

1. Create a new tag using semantic versioning:
   ```bash
   git tag -a vX.Y.Z -m "Release vX.Y.Z"
   ```

2. Push the tag to GitHub:
   ```bash
   git push origin vX.Y.Z
   ```

### 3. Automated Release Process

Once the tag is pushed, the GitHub Actions release workflow is triggered automatically:

1. The CI workflow runs tests to ensure everything is working  
2. The release workflow:
   - Builds binaries for all supported platforms (Linux, macOS, Windows)  
   - Creates distributable archives (tar.gz for Linux/macOS, zip for Windows)  
   - Generates release notes from the commits since the previous tag  
   - Creates a GitHub release with the binaries as assets  

### 4. Verify the Release

After the automated process completes:

1. Go to the [Releases page](https://github.com/user/agent-def/releases) on GitHub  
2. Verify the release notes are correct  
3. Verify all binary artifacts are present and download correctly  
4. Verify the version information is correctly embedded in the binaries  

## Versioning Guidelines

We follow [Semantic Versioning](https://semver.org/) (vMAJOR.MINOR.PATCH):

- **MAJOR**: Increment for incompatible API changes  
- **MINOR**: Increment for backward-compatible new features  
- **PATCH**: Increment for backward-compatible bug fixes  

### Examples:

- **v1.0.0**: Initial stable release  
- **v1.1.0**: Added new template helpers  
- **v1.1.1**: Fixed bug in template engine  
- **v2.0.0**: Incompatible change to command syntax  

## Distribution Channels

After a successful GitHub release, consider updating:

1. The Go package: `go install github.com/user/agent-def/cmd/agent-def@latest`  
2. Update the Homebrew formula if applicable  

## Troubleshooting

If the release workflow fails:

1. Check the workflow logs in GitHub Actions  
2. Fix any issues in a new commit  
3. Delete the problematic tag: `git push --delete origin vX.Y.Z`  
4. Re-tag and push once issues are resolved  
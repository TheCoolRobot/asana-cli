set -e

# Get version from user
if [ -z "$1" ]; then
  echo "Usage: $0 <version>"
  echo "Example: $0 v0.2.0"
  exit 1
fi

VERSION=$1

# Validate version format
if ! [[ $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Invalid version format. Use: v0.0.0"
  exit 1
fi

echo "📦 Preparing release for $VERSION"

# Build and test
echo "🔨 Building..."
make clean
make build

echo "✅ Running tests..."
make test

# Update CHANGELOG
echo "📝 Remember to update CHANGELOG.md with changes for $VERSION"
echo "Press enter when done..."
read

# Create tag
echo "🏷️  Creating git tag..."
git add CHANGELOG.md
git commit -m "Release $VERSION" || true
git tag -a "$VERSION" -m "Release $VERSION"

# Push
echo "📤 Pushing to GitHub..."
git push origin main
git push origin "$VERSION"

echo "✨ Release $VERSION initiated!"
echo "GitHub Actions will now build and publish to Homebrew."
echo ""
echo "Check progress at: https://github.com/TheCoolRobot/asana-cli/actions"
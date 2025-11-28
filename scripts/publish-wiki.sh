#!/bin/bash
# Script to publish wiki pages to GitHub Wiki

set -e

WIKI_REPO="https://github.com/VexoaXYZ/InkWash.wiki.git"
WIKI_DIR="wiki-temp"

echo "Publishing Wiki Pages to GitHub..."

# Clone the wiki repository
echo "Cloning wiki repository..."
if [ -d "$WIKI_DIR" ]; then
    rm -rf "$WIKI_DIR"
fi
git clone "$WIKI_REPO" "$WIKI_DIR"

# Copy wiki files
echo "Copying wiki pages..."
cp -f wiki/*.md "$WIKI_DIR/"

# Commit and push
cd "$WIKI_DIR"
git add .

# Check if there are changes to commit
if git diff --staged --quiet; then
    echo "No changes to publish - wiki is up to date!"
    cd ..
    rm -rf "$WIKI_DIR"
    exit 0
fi

# Commit with timestamp
TIMESTAMP=$(date "+%Y-%m-%d %H:%M:%S")
git commit -m "Update wiki documentation ($TIMESTAMP)

- Updated FiveM key URL to portal.cfx.re/servers/registration-keys
- Improved documentation clarity
- Fixed broken links and outdated information"

echo "Pushing to GitHub Wiki..."
git push origin master

cd ..
rm -rf "$WIKI_DIR"

echo "Wiki pages published successfully!"
echo "View at: https://github.com/VexoaXYZ/InkWash/wiki"

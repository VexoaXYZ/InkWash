#!/bin/bash
# Script to publish wiki pages to GitHub Wiki

set -e

WIKI_REPO="https://github.com/VexoaXYZ/InkWash.wiki.git"
WIKI_DIR="wiki-temp"

echo "📚 Publishing Wiki Pages to GitHub..."

# Clone the wiki repository
echo "Cloning wiki repository..."
if [ -d "$WIKI_DIR" ]; then
    rm -rf "$WIKI_DIR"
fi
git clone "$WIKI_REPO" "$WIKI_DIR"

# Copy wiki files
echo "Copying wiki pages..."
cp wiki/*.md "$WIKI_DIR/"

# Commit and push
cd "$WIKI_DIR"
git add .
git commit -m "Add comprehensive wiki documentation

- Home page with quick links
- Installation Guide for beginners
- Creating Your First Server guide
- Converting GTA5 Mods guide
- Troubleshooting guide

All guides written for accessibility to young users."

echo "Pushing to GitHub Wiki..."
git push origin master

cd ..
rm -rf "$WIKI_DIR"

echo "✅ Wiki pages published successfully!"
echo "View at: https://github.com/VexoaXYZ/InkWash/wiki"

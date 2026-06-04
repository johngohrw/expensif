#!/usr/bin/env bash
set -euo pipefail

# bag-of-skills installer
# Drop this in any project root and run: bash install.sh

REPO_URL="${SKILLS_REPO:-https://github.com/johngohrw/bag-of-skills.git}"
TARGET_DIR=".skills"

echo "📦 Installing skills from $REPO_URL ..."

if [ -f "AGENTS.md" ]; then
  echo "❌ AGENTS.md already exists in this directory. Aborting to avoid overwriting."
  exit 1
fi

if [ -d "$TARGET_DIR/.git" ]; then
  echo "⚠️  $TARGET_DIR/ already exists and is a git repo."
  echo "   Run 'cd $TARGET_DIR && git pull' to update."
  exit 0
fi

if [ -d "$TARGET_DIR" ]; then
  echo "❌ $TARGET_DIR/ exists but is not a git repo. Aborting to avoid overwriting."
  exit 1
fi

git clone --depth 1 "$REPO_URL" "$TARGET_DIR"

echo ""
echo "✅ Skills installed to $TARGET_DIR/"

# Copy AGENTS.md into project root
cp "$TARGET_DIR/AGENTS.md" "AGENTS.md"
echo "📋 Copied AGENTS.md to project root"

echo ""
echo "Done. Your project now has its own .skills/ instance."
echo "Commit .skills/ and AGENTS.md to your repo to share the setup with your team."
echo "To update later, run: cd $TARGET_DIR && git pull"

# Self-destruct
rm -f "$0"
echo "🧹 Removed install.sh"

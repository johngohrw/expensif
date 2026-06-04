#!/usr/bin/env bash
# Run this from your project root as: .skills/update.sh
# Or cd into .skills/ and run: git pull

set -euo pipefail

cd "$(dirname "$0")"
echo "🔄 Updating skills..."
git pull
echo "✅ Skills updated."

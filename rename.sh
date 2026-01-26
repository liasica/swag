#!/usr/bin/env bash

# Find and replace github.com/swaggo/swag with github.com/liasica/swag in all files

set -e

OLD_MODULE="github.com/swaggo/swag"
NEW_MODULE="github.com/liasica/swag"

echo "Starting to replace ${OLD_MODULE} with ${NEW_MODULE}..."
echo ""

# Find all files containing the target string (all text files)
echo "Searching for files containing '${OLD_MODULE}'..."
files=$(grep -rl "${OLD_MODULE}" . \
    --exclude-dir=.git \
    --exclude-dir=vendor \
    --exclude-dir=node_modules \
    --exclude=rename.sh \
    --binary-files=without-match \
    2>/dev/null || true)

if [ -z "$files" ]; then
    echo "No files containing '${OLD_MODULE}' found"
    exit 0
fi

echo "Found the following files:"
echo "$files"
echo ""
echo "Starting replacement..."

# Replace all files at once
echo "$files" | xargs sed -i '' "s|${OLD_MODULE}|${NEW_MODULE}|g"

echo ""
echo "✅ Replacement completed!"
echo ""
echo "📝 Number of affected files: $(echo "$files" | wc -l | tr -d ' ')"
echo ""
echo "🔧 Please run the following command to update dependencies:"
echo "  go mod tidy"
echo ""
echo "📦 If you have subprojects, go to the subproject directory and run:"
echo "  cd example/celler && go mod tidy"

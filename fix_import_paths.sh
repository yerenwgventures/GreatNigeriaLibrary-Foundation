#!/bin/bash

echo "Ì¥ß Fixing all import paths in Go files..."

# Fix all files that import from the wrong path
find backend/ -name "*.go" -type f -exec sed -i 's|github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/|github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/internal/|g' {} \;

find backend/ -name "*.go" -type f -exec sed -i 's|github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/|github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/backend/pkg/|g' {} \;

echo "‚úÖ Import path fixes completed"

# Count remaining issues
echo "Ì¥ç Checking for remaining import issues..."
grep -r "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/" backend/ | wc -l
grep -r "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/" backend/ | wc -l

echo "Ì≥ä Import path fix summary:"
echo "Files processed: $(find backend/ -name "*.go" | wc -l)"
echo "Remaining old internal imports: $(grep -r "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/internal/" backend/ | wc -l || echo 0)"
echo "Remaining old pkg imports: $(grep -r "github.com/yerenwgventures/GreatNigeriaLibrary-Foundation/pkg/" backend/ | wc -l || echo 0)"

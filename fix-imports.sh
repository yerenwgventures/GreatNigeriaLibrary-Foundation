#!/bin/bash

# Script to fix all import paths in foundation to use local imports
# This prepares the foundation for GitHub deployment

echo "🔧 Fixing import paths for GitHub deployment..."

# Create backup directory
BACKUP_DIR="import-fixes-backup-$(date +%Y%m%d-%H%M%S)"
mkdir -p "$BACKUP_DIR"

# Counter for files processed
count=0

# Find all Go files in foundation backend
files=$(find backend -name "*.go" -type f)

for file in $files; do
    echo "Processing: $file"
    
    # Create backup
    cp "$file" "$BACKUP_DIR/$(basename $file).backup"
    
    # Fix the import paths - replace GitHub module with local paths
    sed -i 's|"github\.com/yerenwgventures/GreatNigeriaLibrary/|"../../../|g' "$file"
    
    # Fix specific patterns for different directory levels
    # For files in internal/auth/* referencing pkg/*
    if [[ "$file" == *"internal/auth/"* ]]; then
        sed -i 's|"../../../pkg/|"../../pkg/|g' "$file"
        sed -i 's|"../../../internal/|"../|g' "$file"
    fi
    
    # For files in internal/content/* referencing pkg/*
    if [[ "$file" == *"internal/content/"* ]]; then
        sed -i 's|"../../../pkg/|"../../pkg/|g' "$file"
        sed -i 's|"../../../internal/|"../|g' "$file"
    fi
    
    # For files in internal/discussion/* referencing pkg/*
    if [[ "$file" == *"internal/discussion/"* ]]; then
        sed -i 's|"../../../pkg/|"../../pkg/|g' "$file"
        sed -i 's|"../../../internal/|"../|g' "$file"
    fi
    
    # For files in internal/groups/* referencing pkg/*
    if [[ "$file" == *"internal/groups/"* ]]; then
        sed -i 's|"../../../pkg/|"../../pkg/|g' "$file"
        sed -i 's|"../../../internal/|"../|g' "$file"
    fi
    
    # For files in cmd/* referencing internal/* and pkg/*
    if [[ "$file" == *"cmd/"* ]]; then
        sed -i 's|"../../../internal/|"../../internal/|g' "$file"
        sed -i 's|"../../../pkg/|"../../pkg/|g' "$file"
    fi
    
    count=$((count + 1))
done

echo "✅ Fixed import paths in $count files"
echo "📁 Backups created in: $BACKUP_DIR"

# Test compilation
echo "🧪 Testing compilation..."
go mod tidy
if go build -o foundation-test ./main.go; then
    echo "✅ Foundation compiles successfully!"
    rm -f foundation-test
else
    echo "❌ Compilation failed - check import paths"
fi

echo "📋 Import fix completed!"

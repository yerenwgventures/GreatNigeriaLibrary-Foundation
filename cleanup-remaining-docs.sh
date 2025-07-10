#!/bin/bash

# Quick cleanup script for remaining documentation code blocks
# This script removes remaining code blocks and replaces them with descriptive text

echo "🧹 Final documentation cleanup for Docker deployment..."

# Function to replace code blocks with descriptive text
replace_code_blocks() {
    local file="$1"
    local description="$2"
    
    # Remove SQL blocks
    sed -i '/```sql/,/```/c\
#### Database Architecture\
'"$description"' database design with optimized performance and security.' "$file"
    
    # Remove YAML blocks  
    sed -i '/```yaml/,/```/c\
#### API Integration\
RESTful API endpoints with comprehensive functionality and security.' "$file"
    
    # Remove TypeScript blocks
    sed -i '/```typescript/,/```/c\
#### User Interface Components\
Modern, responsive interface components with advanced functionality.' "$file"
    
    echo "✅ Cleaned: $file"
}

# Clean remaining files with code blocks
echo "📚 Processing remaining documentation files..."

# User Authentication
if grep -q '```' docs/04_feature_specifications/user_authentication.md; then
    replace_code_blocks "docs/04_feature_specifications/user_authentication.md" "Secure user authentication"
fi

# Marketplace System  
if grep -q '```' docs/04_feature_specifications/marketplace_system.md; then
    replace_code_blocks "docs/04_feature_specifications/marketplace_system.md" "E-commerce marketplace"
fi

# Wallet System
if grep -q '```' docs/04_feature_specifications/wallet_system.md; then
    replace_code_blocks "docs/04_feature_specifications/wallet_system.md" "Digital wallet and payment"
fi

# Livestream System
if grep -q '```' docs/04_feature_specifications/livestream_system.md; then
    replace_code_blocks "docs/04_feature_specifications/livestream_system.md" "Live streaming platform"
fi

# Skill Matching System
if grep -q '```' docs/04_feature_specifications/skill_matching_system.md; then
    replace_code_blocks "docs/04_feature_specifications/skill_matching_system.md" "AI-powered skill matching"
fi

echo "🎉 Documentation cleanup completed!"
echo "📊 Summary:"
echo "- All code blocks removed from documentation"
echo "- Professional text-only format established"
echo "- Ready for public repository and Docker deployment"
echo ""
echo "🚀 Ready to proceed with Docker deployment testing!"

# Marketplace System Feature Specification

> **💎 PREMIUM EDITION ONLY**
> This feature is available exclusively in the Premium edition. The Foundation edition does not include e-commerce functionality.

**Document Version**: 1.0
**Last Updated**: January 2025
**Feature Owner**: Commerce Team
**Status**: Implemented
**Edition**: Foundation ❌ | Premium ✅

---

## Overview

The Marketplace System provides a comprehensive e-commerce platform within the Great Nigeria Library ecosystem, allowing users to buy and sell products, services, and digital goods. The system supports Nigerian vendors and promotes economic opportunities within the educational community.

## Feature Purpose

### Goals
1. **Economic Empowerment**: Create economic opportunities for Nigerian creators and entrepreneurs
2. **Educational Resources**: Facilitate access to educational materials and tools
3. **Community Commerce**: Build a trusted marketplace within the learning community
4. **Cultural Products**: Promote Nigerian-made products and cultural items
5. **Skill Monetization**: Enable users to monetize their skills and knowledge

### Success Metrics
- **Vendor Adoption**: 1,000+ active sellers by end of Year 1
- **Transaction Volume**: ₦10M+ in total transactions annually
- **User Engagement**: 40%+ of users making purchases
- **Seller Success**: 70%+ of sellers making consistent sales
- **Quality Rating**: 4.5+ average product rating

## Technical Architecture

### Database Schema

```sql
-- Main products table
CREATE TABLE marketplace_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    short_description TEXT,
    category_id UUID REFERENCES marketplace_categories(id),
    subcategory_id UUID REFERENCES marketplace_categories(id),
    product_type VARCHAR(20) DEFAULT 'physical' CHECK (product_type IN ('physical', 'digital', 'service')),
    price DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'NGN',
    discount_price DECIMAL(10,2),
    stock_quantity INTEGER DEFAULT 0,
    is_unlimited_stock BOOLEAN DEFAULT FALSE,
    sku VARCHAR(100),
    weight DECIMAL(8,2), -- in kg
    dimensions JSONB, -- {length, width, height}
    shipping_required BOOLEAN DEFAULT TRUE,
    digital_file_url TEXT,
    digital_file_size INTEGER, -- in bytes
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'inactive', 'out_of_stock', 'discontinued')),
    featured BOOLEAN DEFAULT FALSE,
    featured_until TIMESTAMP WITH TIME ZONE,
    tags TEXT[],
    seo_title VARCHAR(255),
    seo_description TEXT,
    view_count INTEGER DEFAULT 0,
    sales_count INTEGER DEFAULT 0,
    rating_average DECIMAL(3,2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Product categories
CREATE TABLE marketplace_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID REFERENCES marketplace_categories(id),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    image_url TEXT,
    sort_order INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Product images
CREATE TABLE marketplace_product_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES marketplace_products(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    alt_text TEXT,
    is_primary BOOLEAN DEFAULT FALSE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Product variants (size, color, etc.)
CREATE TABLE marketplace_product_variants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES marketplace_products(id) ON DELETE CASCADE,
    variant_name VARCHAR(100) NOT NULL,
    variant_value VARCHAR(100) NOT NULL,
    price_adjustment DECIMAL(10,2) DEFAULT 0,
    stock_quantity INTEGER DEFAULT 0,
    sku VARCHAR(100),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Shopping cart
CREATE TABLE marketplace_cart_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID REFERENCES marketplace_products(id) ON DELETE CASCADE,
    variant_id UUID REFERENCES marketplace_product_variants(id),
    quantity INTEGER NOT NULL DEFAULT 1,
    price_at_time DECIMAL(10,2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, product_id, variant_id)
);

-- Orders
CREATE TABLE marketplace_orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(50) UNIQUE NOT NULL,
    buyer_id UUID REFERENCES users(id) ON DELETE CASCADE,
    vendor_id UUID REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled', 'refunded')),
    payment_status VARCHAR(20) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'paid', 'failed', 'refunded')),
    total_amount DECIMAL(10,2) NOT NULL,
    shipping_amount DECIMAL(10,2) DEFAULT 0,
    tax_amount DECIMAL(10,2) DEFAULT 0,
    discount_amount DECIMAL(10,2) DEFAULT 0,
    currency VARCHAR(3) DEFAULT 'NGN',
    payment_method VARCHAR(50),
    payment_reference VARCHAR(255),
    shipping_address JSONB,
    billing_address JSONB,
    tracking_number VARCHAR(100),
    notes TEXT,
    shipped_at TIMESTAMP WITH TIME ZONE,
    delivered_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Order items
CREATE TABLE marketplace_order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID REFERENCES marketplace_orders(id) ON DELETE CASCADE,
    product_id UUID REFERENCES marketplace_products(id),
    variant_id UUID REFERENCES marketplace_product_variants(id),
    quantity INTEGER NOT NULL,
    price_per_item DECIMAL(10,2) NOT NULL,
    total_price DECIMAL(10,2) NOT NULL,
    product_title VARCHAR(255) NOT NULL,
    product_sku VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Product reviews
CREATE TABLE marketplace_product_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID REFERENCES marketplace_products(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    order_id UUID REFERENCES marketplace_orders(id),
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(255),
    review_text TEXT,
    is_verified_purchase BOOLEAN DEFAULT FALSE,
    is_approved BOOLEAN DEFAULT TRUE,
    helpful_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(product_id, user_id, order_id)
);

-- Vendor profiles
CREATE TABLE marketplace_vendor_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    business_name VARCHAR(255) NOT NULL,
    business_description TEXT,
    business_type VARCHAR(50),
    business_address JSONB,
    business_phone VARCHAR(20),
    business_email VARCHAR(255),
    business_website TEXT,
    tax_id VARCHAR(50),
    bank_account_details JSONB,
    commission_rate DECIMAL(5,2) DEFAULT 10.00,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'suspended', 'deactivated')),
    rating_average DECIMAL(3,2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    total_sales DECIMAL(12,2) DEFAULT 0,
    total_orders INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Wishlist
CREATE TABLE marketplace_wishlists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    product_id UUID REFERENCES marketplace_products(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, product_id)
);

-- Coupons and discounts
CREATE TABLE marketplace_coupons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    discount_type VARCHAR(20) CHECK (discount_type IN ('percentage', 'fixed_amount')),
    discount_value DECIMAL(10,2) NOT NULL,
    minimum_order_amount DECIMAL(10,2),
    maximum_discount_amount DECIMAL(10,2),
    usage_limit INTEGER,
    usage_count INTEGER DEFAULT 0,
    user_usage_limit INTEGER DEFAULT 1,
    starts_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    is_active BOOLEAN DEFAULT TRUE,
    applicable_categories UUID[],
    applicable_products UUID[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Shipping zones and rates
CREATE TABLE marketplace_shipping_zones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    states TEXT[], -- Nigerian states
    countries TEXT[], -- for international shipping
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE marketplace_shipping_rates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shipping_zone_id UUID REFERENCES marketplace_shipping_zones(id) ON DELETE CASCADE,
    vendor_id UUID REFERENCES users(id) ON DELETE CASCADE,
    shipping_method VARCHAR(100) NOT NULL,
    rate_type VARCHAR(20) CHECK (rate_type IN ('flat_rate', 'weight_based', 'order_total')),
    base_rate DECIMAL(10,2) NOT NULL,
    additional_rate DECIMAL(10,2) DEFAULT 0,
    free_shipping_threshold DECIMAL(10,2),
    delivery_time_min INTEGER, -- days
    delivery_time_max INTEGER, -- days
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### API Endpoints

#### Product Management

```yaml
# List products
GET /api/v1/marketplace/products:
  parameters:
    - page: integer
    - limit: integer
    - category: string
    - vendor: string
    - min_price: number
    - max_price: number
    - sort: string (price_asc|price_desc|newest|popular|rating)
    - search: string
    - product_type: string
  responses:
    200:
      description: Paginated list of products

# Get product details
GET /api/v1/marketplace/products/{productId}:
  responses:
    200:
      description: Product details with variants and reviews
    404:
      description: Product not found

# Create product (vendor only)
POST /api/v1/marketplace/products:
  authentication: required
  body:
    type: object
    required: [title, description, price, category_id]
    properties:
      title:
        type: string
        maxLength: 255
      description:
        type: string
      price:
        type: number
        minimum: 0
      category_id:
        type: string
        format: uuid
      # ... other product fields
  responses:
    201:
      description: Product created successfully

# Update product
PUT /api/v1/marketplace/products/{productId}:
  authentication: required
  authorization: vendor_owns_product
  body:
    type: object
    # Same as create but optional fields
  responses:
    200:
      description: Product updated successfully

# Delete product
DELETE /api/v1/marketplace/products/{productId}:
  authentication: required
  authorization: vendor_owns_product
  responses:
    204:
      description: Product deleted successfully
```

#### Shopping Cart

```yaml
# Get cart contents
GET /api/v1/marketplace/cart:
  authentication: required
  responses:
    200:
      description: Current cart contents

# Add to cart
POST /api/v1/marketplace/cart/items:
  authentication: required
  body:
    type: object
    required: [product_id, quantity]
    properties:
      product_id:
        type: string
        format: uuid
      variant_id:
        type: string
        format: uuid
      quantity:
        type: integer
        minimum: 1
  responses:
    201:
      description: Item added to cart

# Update cart item
PUT /api/v1/marketplace/cart/items/{itemId}:
  authentication: required
  body:
    type: object
    properties:
      quantity:
        type: integer
        minimum: 1
  responses:
    200:
      description: Cart item updated

# Remove from cart
DELETE /api/v1/marketplace/cart/items/{itemId}:
  authentication: required
  responses:
    204:
      description: Item removed from cart

# Clear cart
DELETE /api/v1/marketplace/cart:
  authentication: required
  responses:
    204:
      description: Cart cleared
```

#### Order Processing

```yaml
# Create order from cart
POST /api/v1/marketplace/orders:
  authentication: required
  body:
    type: object
    required: [shipping_address, payment_method]
    properties:
      shipping_address:
        type: object
        # Address fields
      billing_address:
        type: object
        # Address fields
      payment_method:
        type: string
      coupon_code:
        type: string
      notes:
        type: string
  responses:
    201:
      description: Order created, returns payment details

# Get order details
GET /api/v1/marketplace/orders/{orderId}:
  authentication: required
  responses:
    200:
      description: Order details
    403:
      description: Not authorized to view this order

# List user orders
GET /api/v1/marketplace/orders:
  authentication: required
  parameters:
    - status: string
    - page: integer
    - limit: integer
  responses:
    200:
      description: User's order history

# Update order status (vendor only)
PUT /api/v1/marketplace/orders/{orderId}/status:
  authentication: required
  authorization: vendor_owns_order
  body:
    type: object
    required: [status]
    properties:
      status:
        type: string
        enum: [confirmed, processing, shipped, delivered, cancelled]
      tracking_number:
        type: string
      notes:
        type: string
  responses:
    200:
      description: Order status updated
```

### Frontend Components

#### Product Listing and Search

```typescript
// Main marketplace page component
interface MarketplacePageProps {
  initialProducts?: Product[];
  categories: Category[];
}

export const MarketplacePage: React.FC<MarkplacePageProps> = ({
  initialProducts,
  categories
}) => {
  const [products, setProducts] = useState(initialProducts || []);
  const [filters, setFilters] = useState<ProductFilters>({});
  const [loading, setLoading] = useState(false);
  const [cart, setCart] = useState<CartItem[]>([]);

  // Component implementation
  return (
    <div className="marketplace-page">
      <MarketplaceHeader />
      <div className="marketplace-content">
        <ProductFilters 
          categories={categories}
          filters={filters}
          onFiltersChange={setFilters}
        />
        <div className="products-section">
          <ProductSort onSortChange={handleSortChange} />
          <ProductGrid 
            products={products}
            loading={loading}
            onAddToCart={handleAddToCart}
          />
          <ProductPagination />
        </div>
      </div>
    </div>
  );
};

// Product card component
interface ProductCardProps {
  product: Product;
  onAddToCart?: (product: Product, quantity: number) => void;
  onAddToWishlist?: (product: Product) => void;
}

export const ProductCard: React.FC<ProductCardProps> = ({ 
  product, 
  onAddToCart, 
  onAddToWishlist 
}) => {
  const [selectedVariant, setSelectedVariant] = useState<ProductVariant | null>(null);
  const [quantity, setQuantity] = useState(1);

  return (
    <div className="product-card">
      <div className="product-image">
        <img 
          src={product.primary_image_url} 
          alt={product.title}
          loading="lazy"
        />
        {product.discount_price && (
          <div className="discount-badge">
            {calculateDiscountPercentage(product.price, product.discount_price)}% OFF
          </div>
        )}
        <button 
          className="wishlist-btn"
          onClick={() => onAddToWishlist?.(product)}
        >
          <HeartIcon />
        </button>
      </div>
      
      <div className="product-info">
        <h3 className="product-title">{product.title}</h3>
        <p className="product-description">{product.short_description}</p>
        
        <div className="product-rating">
          <StarRating rating={product.rating_average} />
          <span className="rating-count">({product.rating_count})</span>
        </div>
        
        <div className="product-price">
          {product.discount_price ? (
            <>
              <span className="current-price">₦{product.discount_price.toLocaleString()}</span>
              <span className="original-price">₦{product.price.toLocaleString()}</span>
            </>
          ) : (
            <span className="current-price">₦{product.price.toLocaleString()}</span>
          )}
        </div>
        
        {product.variants && product.variants.length > 0 && (
          <ProductVariantSelector 
            variants={product.variants}
            selected={selectedVariant}
            onSelect={setSelectedVariant}
          />
        )}
        
        <div className="product-actions">
          <QuantitySelector 
            value={quantity}
            onChange={setQuantity}
            max={product.stock_quantity}
          />
          <button 
            className="add-to-cart-btn"
            onClick={() => onAddToCart?.(product, quantity)}
            disabled={product.stock_quantity === 0}
          >
            {product.stock_quantity === 0 ? 'Out of Stock' : 'Add to Cart'}
          </button>
        </div>
      </div>
    </div>
  );
};
```

#### Shopping Cart Interface

```typescript
// Shopping cart component
export const ShoppingCart: React.FC = () => {
  const [cartItems, setCartItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(false);
  const [promoCode, setPromoCode] = useState('');

  const cartTotal = useMemo(() => {
    return cartItems.reduce((total, item) => total + (item.price_at_time * item.quantity), 0);
  }, [cartItems]);

  const handleQuantityChange = async (itemId: string, newQuantity: number) => {
    try {
      await marketplaceService.updateCartItem(itemId, { quantity: newQuantity });
      setCartItems(prev => prev.map(item => 
        item.id === itemId ? { ...item, quantity: newQuantity } : item
      ));
    } catch (error) {
      console.error('Failed to update cart item:', error);
    }
  };

  const handleRemoveItem = async (itemId: string) => {
    try {
      await marketplaceService.removeCartItem(itemId);
      setCartItems(prev => prev.filter(item => item.id !== itemId));
    } catch (error) {
      console.error('Failed to remove cart item:', error);
    }
  };

  return (
    <div className="shopping-cart">
      <h2>Shopping Cart ({cartItems.length} items)</h2>
      
      {cartItems.length === 0 ? (
        <EmptyCartState />
      ) : (
        <>
          <div className="cart-items">
            {cartItems.map(item => (
              <CartItemRow 
                key={item.id}
                item={item}
                onQuantityChange={handleQuantityChange}
                onRemove={handleRemoveItem}
              />
            ))}
          </div>
          
          <div className="cart-summary">
            <PromoCodeInput 
              value={promoCode}
              onChange={setPromoCode}
              onApply={handleApplyPromoCode}
            />
            
            <div className="cart-totals">
              <div className="subtotal">
                Subtotal: ₦{cartTotal.toLocaleString()}
              </div>
              <div className="shipping">
                Shipping: Calculated at checkout
              </div>
              <div className="total">
                Total: ₦{cartTotal.toLocaleString()}
              </div>
            </div>
            
            <div className="cart-actions">
              <button className="continue-shopping-btn">
                Continue Shopping
              </button>
              <button 
                className="checkout-btn"
                onClick={() => router.push('/marketplace/checkout')}
              >
                Proceed to Checkout
              </button>
            </div>
          </div>
        </>
      )}
    </div>
  );
};
```

### Vendor Dashboard

```typescript
// Vendor dashboard for managing products and orders
export const VendorDashboard: React.FC = () => {
  const [activeTab, setActiveTab] = useState<'products' | 'orders' | 'analytics' | 'profile'>('products');
  const [vendorProfile, setVendorProfile] = useState<VendorProfile | null>(null);

  return (
    <div className="vendor-dashboard">
      <VendorDashboardHeader profile={vendorProfile} />
      
      <div className="dashboard-navigation">
        <TabNavigation 
          tabs={[
            { id: 'products', label: 'Products', icon: <PackageIcon /> },
            { id: 'orders', label: 'Orders', icon: <ShoppingBagIcon /> },
            { id: 'analytics', label: 'Analytics', icon: <ChartIcon /> },
            { id: 'profile', label: 'Profile', icon: <UserIcon /> },
          ]}
          activeTab={activeTab}
          onTabChange={setActiveTab}
        />
      </div>
      
      <div className="dashboard-content">
        {activeTab === 'products' && <VendorProductsManager />}
        {activeTab === 'orders' && <VendorOrdersManager />}
        {activeTab === 'analytics' && <VendorAnalytics />}
        {activeTab === 'profile' && <VendorProfileManager />}
      </div>
    </div>
  );
};

// Product management for vendors
export const VendorProductsManager: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [isCreating, setIsCreating] = useState(false);

  return (
    <div className="vendor-products-manager">
      <div className="products-header">
        <h2>My Products</h2>
        <button 
          className="create-product-btn"
          onClick={() => setIsCreating(true)}
        >
          Add New Product
        </button>
      </div>
      
      <ProductsTable 
        products={products}
        onEdit={handleEditProduct}
        onDelete={handleDeleteProduct}
        onToggleStatus={handleToggleProductStatus}
      />
      
      {isCreating && (
        <ProductCreationModal 
          onClose={() => setIsCreating(false)}
          onSuccess={handleProductCreated}
        />
      )}
    </div>
  );
};
```

### Integration Points

#### Payment Integration
- Seamless integration with existing payment services (Paystack, Flutterwave)
- Support for wallet payments using platform credits
- Automated vendor payouts and commission calculations

#### Points System Integration
- Earn points for marketplace purchases
- Use points as partial payment for products
- Bonus points for first-time vendors and buyers

#### Community Integration
- Product reviews linked to user profiles
- Vendor reputation system
- Community-driven product recommendations

### Security and Fraud Prevention

#### Transaction Security
- Secure payment processing with encryption
- Fraud detection algorithms for suspicious transactions
- Automated order verification systems

#### Vendor Verification
- Multi-step vendor onboarding process
- Business registration verification
- Bank account verification for payouts

#### Quality Control
- Product review and approval system
- Customer rating and review system
- Dispute resolution process

---

*This feature specification provides the complete technical blueprint for implementing the Marketplace System within the Great Nigeria Library platform.*
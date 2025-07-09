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
#### E-commerce Transaction System
Comprehensive order and payment processing:

- **Shopping Cart Management**: Persistent cart functionality with item quantity and variant tracking
- **Order Processing**: Complete order lifecycle from creation to fulfillment
- **Payment Integration**: Multiple payment gateway support with secure transaction processing
- **Order Tracking**: Real-time order status updates and delivery tracking
- **Digital Product Delivery**: Automated delivery system for digital products and services
- **Refund Management**: Streamlined refund and return processing with automated workflows

#### Search and Discovery Engine
Advanced product discovery capabilities:

- **Full-Text Search**: Comprehensive search across product titles, descriptions, and tags
- **Category Filtering**: Hierarchical category navigation with advanced filtering options
- **Price Range Filtering**: Dynamic price filtering with currency conversion support
- **Vendor Filtering**: Search and filter products by specific vendors or seller ratings
- **Recommendation Engine**: AI-powered product recommendations based on user behavior
- **Trending Products**: Dynamic trending and featured product highlighting

#### Review and Rating System
Community-driven quality assurance:

- **Product Reviews**: Detailed customer review system with text and rating components
- **Vendor Ratings**: Comprehensive vendor performance tracking and public ratings
- **Review Moderation**: Automated and manual review moderation for quality control
- **Verified Purchase Reviews**: Enhanced credibility through verified purchase validation
- **Review Analytics**: Detailed analytics for vendors to improve product quality
- **Community Feedback**: User-driven feedback system for continuous improvement

### Business Intelligence and Analytics

#### Sales Performance Tracking
Comprehensive business analytics for vendors and platform administrators:

- **Revenue Analytics**: Detailed sales tracking with revenue breakdowns by product, category, and time period
- **Performance Metrics**: Key performance indicators including conversion rates, average order value, and customer lifetime value
- **Vendor Dashboards**: Personalized analytics dashboards for vendors to track their business performance
- **Market Insights**: Platform-wide analytics showing trending products, popular categories, and market opportunities
- **Financial Reporting**: Automated financial reports for commission tracking, tax reporting, and business planning
- **Customer Behavior Analysis**: Detailed insights into customer purchasing patterns and preferences

#### Inventory and Supply Chain Management
Advanced inventory control and logistics:

- **Stock Management**: Real-time inventory tracking with automated low-stock alerts and reorder notifications
- **Supplier Integration**: Direct integration with suppliers for automated inventory replenishment
- **Warehouse Management**: Multi-warehouse support with location-based inventory tracking
- **Shipping Integration**: Automated shipping calculations and carrier integration for seamless fulfillment
- **Return Management**: Streamlined return processing with automated refund and restocking workflows
- **Quality Control**: Inventory quality tracking with batch management and expiration date monitoring

### Marketing and Promotion Tools

#### Discount and Coupon Management
Comprehensive promotional campaign system:

- **Coupon Creation**: Flexible coupon system with percentage and fixed amount discounts
- **Promotional Campaigns**: Time-limited promotional campaigns with automatic activation and expiration
- **Bulk Discount Rules**: Volume-based pricing with automatic discount application
- **Category Promotions**: Category-wide promotional campaigns for seasonal sales and special events
- **Vendor Promotions**: Individual vendor promotional tools for independent marketing campaigns
- **Customer Segmentation**: Targeted promotions based on customer behavior and purchase history

#### Shipping and Logistics Integration
Advanced shipping and delivery management:

- **Shipping Zone Configuration**: Flexible shipping zones covering Nigerian states and international destinations
- **Dynamic Rate Calculation**: Intelligent shipping rate calculation based on weight, distance, and order value
- **Carrier Integration**: Direct integration with major Nigerian and international shipping carriers
- **Delivery Tracking**: Real-time package tracking with customer notifications and updates
- **Free Shipping Thresholds**: Configurable free shipping options to encourage larger orders
- **Express Delivery Options**: Premium delivery services with guaranteed delivery timeframes

### API Integration and Services

#### RESTful API Architecture
Comprehensive API system for marketplace functionality:

- **Product Management APIs**: Complete CRUD operations for product creation, updating, and management
- **Search and Filtering APIs**: Advanced search capabilities with multiple filtering and sorting options
- **Category Management APIs**: Hierarchical category management with nested category support
- **Vendor Management APIs**: Vendor registration, profile management, and performance tracking
- **Authentication and Authorization**: Secure API access with role-based permissions and vendor ownership validation
- **Pagination and Performance**: Optimized API responses with efficient pagination and caching strategies

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
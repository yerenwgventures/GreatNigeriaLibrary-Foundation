# Affiliate System Feature Specification

**Document Version**: 1.0  
**Last Updated**: January 2025  
**Feature Owner**: Marketing Team  
**Status**: Implemented

---

## Overview

The Affiliate System enables users to earn commissions by promoting the Great Nigeria Library platform and its products. The system includes referral tracking, commission management, and comprehensive analytics for affiliates to optimize their earning potential.

## Feature Purpose

### Goals
1. **User Acquisition**: Leverage user networks to acquire new platform members
2. **Revenue Sharing**: Provide earning opportunities for active community members
3. **Organic Growth**: Foster word-of-mouth marketing through incentives
4. **Community Building**: Encourage users to become platform advocates
5. **Performance Tracking**: Provide transparent analytics for affiliate success

### Success Metrics
- **Affiliate Adoption**: 5,000+ active affiliates by end of Year 1
- **Referral Conversion**: 15%+ conversion rate from referrals to premium members
- **Commission Payouts**: ₦2M+ in affiliate commissions annually
- **Quality Referrals**: 80%+ of referred users remain active after 30 days
- **Top Performer Growth**: 100+ affiliates earning ₦10,000+ monthly

## Technical Architecture

### Database Schema

```sql
-- Affiliate profiles
CREATE TABLE affiliate_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    affiliate_code VARCHAR(20) UNIQUE NOT NULL,
    custom_link_slug VARCHAR(50) UNIQUE,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('pending', 'active', 'suspended', 'terminated')),
    commission_tier INTEGER DEFAULT 1, -- 1=basic, 2=silver, 3=gold, 4=platinum
    total_referrals INTEGER DEFAULT 0,
    active_referrals INTEGER DEFAULT 0,
    total_earnings DECIMAL(15,2) DEFAULT 0,
    pending_earnings DECIMAL(15,2) DEFAULT 0,
    paid_earnings DECIMAL(15,2) DEFAULT 0,
    current_month_earnings DECIMAL(10,2) DEFAULT 0,
    last_month_earnings DECIMAL(10,2) DEFAULT 0,
    conversion_rate DECIMAL(5,2) DEFAULT 0,
    payment_method VARCHAR(50) DEFAULT 'bank_transfer',
    bank_details JSONB,
    paypal_email VARCHAR(255),
    tax_id VARCHAR(50),
    marketing_materials_access BOOLEAN DEFAULT TRUE,
    analytics_access_level INTEGER DEFAULT 1,
    special_promotions_eligible BOOLEAN DEFAULT TRUE,
    notes TEXT,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    last_activity_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Referral tracking
CREATE TABLE affiliate_referrals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    affiliate_id UUID REFERENCES affiliate_profiles(id) ON DELETE CASCADE,
    referred_user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    referral_code VARCHAR(20) NOT NULL,
    source_url TEXT,
    utm_source VARCHAR(100),
    utm_medium VARCHAR(100),
    utm_campaign VARCHAR(100),
    utm_content VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    device_type VARCHAR(50),
    browser VARCHAR(100),
    country VARCHAR(2),
    city VARCHAR(100),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'converted', 'rejected')),
    conversion_type VARCHAR(50), -- registration, premium_purchase, course_enrollment
    conversion_value DECIMAL(10,2) DEFAULT 0,
    commission_rate DECIMAL(5,2) NOT NULL,
    commission_amount DECIMAL(10,2) DEFAULT 0,
    commission_status VARCHAR(20) DEFAULT 'pending' CHECK (commission_status IN ('pending', 'approved', 'paid', 'cancelled')),
    first_click_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    conversion_at TIMESTAMP WITH TIME ZONE,
    commission_approved_at TIMESTAMP WITH TIME ZONE,
    commission_paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Commission structure and tiers
CREATE TABLE commission_tiers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tier_level INTEGER UNIQUE NOT NULL,
    tier_name VARCHAR(50) NOT NULL,
    min_referrals INTEGER NOT NULL,
    min_monthly_sales DECIMAL(10,2),
    base_commission_rate DECIMAL(5,2) NOT NULL,
    bonus_commission_rate DECIMAL(5,2) DEFAULT 0,
    requirements JSONB,
    benefits JSONB,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Commission transactions
CREATE TABLE affiliate_commissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    affiliate_id UUID REFERENCES affiliate_profiles(id) ON DELETE CASCADE,
    referral_id UUID REFERENCES affiliate_referrals(id) ON DELETE CASCADE,
    transaction_type VARCHAR(50) NOT NULL, -- earned, bonus, adjustment, payout
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'NGN',
    description TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'paid', 'cancelled')),
    payout_id UUID REFERENCES affiliate_payouts(id),
    metadata JSONB,
    processed_by UUID REFERENCES users(id),
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Payout management
CREATE TABLE affiliate_payouts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    affiliate_id UUID REFERENCES affiliate_profiles(id) ON DELETE CASCADE,
    payout_reference VARCHAR(100) UNIQUE NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'NGN',
    payment_method VARCHAR(50) NOT NULL,
    payment_details JSONB NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'cancelled')),
    failure_reason TEXT,
    transaction_id VARCHAR(255),
    fees DECIMAL(10,2) DEFAULT 0,
    net_amount DECIMAL(10,2) NOT NULL,
    requested_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Marketing materials and links
CREATE TABLE affiliate_marketing_materials (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    material_type VARCHAR(50) NOT NULL, -- banner, text_ad, email_template, social_post
    content JSONB NOT NULL,
    file_url TEXT,
    thumbnail_url TEXT,
    dimensions VARCHAR(20), -- for banners (e.g., "728x90")
    target_audience VARCHAR(100),
    campaign_id VARCHAR(100),
    click_count INTEGER DEFAULT 0,
    conversion_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Link tracking
CREATE TABLE affiliate_link_clicks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    affiliate_id UUID REFERENCES affiliate_profiles(id) ON DELETE CASCADE,
    link_id UUID REFERENCES affiliate_marketing_materials(id),
    referral_code VARCHAR(20) NOT NULL,
    clicked_url TEXT NOT NULL,
    destination_url TEXT NOT NULL,
    ip_address INET,
    user_agent TEXT,
    referer_url TEXT,
    utm_parameters JSONB,
    device_info JSONB,
    location_info JSONB,
    session_id VARCHAR(255),
    converted BOOLEAN DEFAULT FALSE,
    conversion_at TIMESTAMP WITH TIME ZONE,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Affiliate targets and goals
CREATE TABLE affiliate_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    affiliate_id UUID REFERENCES affiliate_profiles(id) ON DELETE CASCADE,
    target_type VARCHAR(50) NOT NULL, -- monthly_referrals, monthly_sales, quarterly_goals
    target_value DECIMAL(10,2) NOT NULL,
    current_value DECIMAL(10,2) DEFAULT 0,
    target_period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    target_period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    reward_type VARCHAR(50), -- bonus_commission, tier_upgrade, special_access
    reward_value DECIMAL(10,2),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'achieved', 'failed', 'cancelled')),
    achieved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Promotional campaigns
CREATE TABLE affiliate_campaigns (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    campaign_name VARCHAR(255) NOT NULL,
    description TEXT,
    campaign_type VARCHAR(50) NOT NULL, -- seasonal, product_launch, special_offer
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date TIMESTAMP WITH TIME ZONE NOT NULL,
    target_audience VARCHAR(100),
    commission_bonus DECIMAL(5,2) DEFAULT 0,
    special_rules JSONB,
    marketing_materials UUID[] REFERENCES affiliate_marketing_materials(id),
    eligible_tiers INTEGER[],
    max_participants INTEGER,
    current_participants INTEGER DEFAULT 0,
    total_clicks INTEGER DEFAULT 0,
    total_conversions INTEGER DEFAULT 0,
    total_commission_paid DECIMAL(15,2) DEFAULT 0,
    status VARCHAR(20) DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'paused', 'completed', 'cancelled')),
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_affiliate_profiles_user_id ON affiliate_profiles(user_id);
CREATE INDEX idx_affiliate_profiles_code ON affiliate_profiles(affiliate_code);
CREATE INDEX idx_affiliate_referrals_affiliate_id ON affiliate_referrals(affiliate_id);
CREATE INDEX idx_affiliate_referrals_referred_user ON affiliate_referrals(referred_user_id);
CREATE INDEX idx_affiliate_referrals_code ON affiliate_referrals(referral_code);
CREATE INDEX idx_affiliate_referrals_status ON affiliate_referrals(status);
CREATE INDEX idx_affiliate_commissions_affiliate_id ON affiliate_commissions(affiliate_id);
CREATE INDEX idx_affiliate_commissions_status ON affiliate_commissions(status);
CREATE INDEX idx_affiliate_payouts_affiliate_id ON affiliate_payouts(affiliate_id);
CREATE INDEX idx_affiliate_payouts_status ON affiliate_payouts(status);
CREATE INDEX idx_affiliate_link_clicks_affiliate_id ON affiliate_link_clicks(affiliate_id);
CREATE INDEX idx_affiliate_link_clicks_clicked_at ON affiliate_link_clicks(clicked_at DESC);
```

### API Endpoints

#### Affiliate Management

```yaml
# Join affiliate program
POST /api/v1/affiliate/join:
  authentication: required
  body:
    type: object
    required: [preferred_code]
    properties:
      preferred_code:
        type: string
        pattern: '^[A-Z0-9]{6,20}$'
      marketing_experience:
        type: string
      promotion_channels:
        type: array
        items:
          type: string
      bank_details:
        type: object
  responses:
    201:
      description: Affiliate application submitted
    400:
      description: Code already taken or invalid

# Get affiliate dashboard
GET /api/v1/affiliate/dashboard:
  authentication: required
  authorization: affiliate_member
  responses:
    200:
      description: Comprehensive affiliate dashboard data
      schema:
        type: object
        properties:
          profile:
            type: object
          earnings_summary:
            type: object
          recent_referrals:
            type: array
          performance_metrics:
            type: object
          targets:
            type: array

# Update affiliate profile
PUT /api/v1/affiliate/profile:
  authentication: required
  authorization: affiliate_member
  body:
    type: object
    properties:
      custom_link_slug:
        type: string
      payment_method:
        type: string
      bank_details:
        type: object
      marketing_bio:
        type: string
  responses:
    200:
      description: Profile updated successfully

# Generate affiliate link
POST /api/v1/affiliate/links:
  authentication: required
  authorization: affiliate_member
  body:
    type: object
    required: [destination_url]
    properties:
      destination_url:
        type: string
        format: url
      utm_campaign:
        type: string
      utm_content:
        type: string
      custom_parameters:
        type: object
  responses:
    201:
      description: Affiliate link generated
```

#### Referral Tracking

```yaml
# Track referral click
POST /api/v1/affiliate/track-click:
  body:
    type: object
    required: [affiliate_code, destination_url]
    properties:
      affiliate_code:
        type: string
      destination_url:
        type: string
      source_url:
        type: string
      utm_parameters:
        type: object
      device_info:
        type: object
  responses:
    200:
      description: Click tracked successfully

# Get referral history
GET /api/v1/affiliate/referrals:
  authentication: required
  authorization: affiliate_member
  parameters:
    - status: string
    - conversion_type: string
    - from_date: string
    - to_date: string
    - page: integer
    - limit: integer
  responses:
    200:
      description: Paginated referral history

# Get referral details
GET /api/v1/affiliate/referrals/{referralId}:
  authentication: required
  authorization: affiliate_member
  responses:
    200:
      description: Detailed referral information
    404:
      description: Referral not found

# Convert referral
POST /api/v1/affiliate/referrals/{referralId}/convert:
  authentication: internal_service
  body:
    type: object
    required: [conversion_type, conversion_value]
    properties:
      conversion_type:
        type: string
        enum: [registration, premium_purchase, course_enrollment]
      conversion_value:
        type: number
      additional_data:
        type: object
  responses:
    200:
      description: Referral conversion recorded
```

#### Commission Management

```yaml
# Get commission history
GET /api/v1/affiliate/commissions:
  authentication: required
  authorization: affiliate_member
  parameters:
    - transaction_type: string
    - status: string
    - from_date: string
    - to_date: string
    - page: integer
    - limit: integer
  responses:
    200:
      description: Commission transaction history

# Request payout
POST /api/v1/affiliate/payouts:
  authentication: required
  authorization: affiliate_member
  body:
    type: object
    required: [amount, payment_method]
    properties:
      amount:
        type: number
        minimum: 1000 # Minimum payout ₦1,000
      payment_method:
        type: string
        enum: [bank_transfer, paypal]
      payment_details:
        type: object
  responses:
    201:
      description: Payout request submitted
    400:
      description: Insufficient balance or invalid details

# Get payout history
GET /api/v1/affiliate/payouts:
  authentication: required
  authorization: affiliate_member
  parameters:
    - status: string
    - page: integer
    - limit: integer
  responses:
    200:
      description: Payout request history

# Cancel payout request
DELETE /api/v1/affiliate/payouts/{payoutId}:
  authentication: required
  authorization: affiliate_member
  responses:
    204:
      description: Payout request cancelled
    400:
      description: Cannot cancel processed payout
```

#### Marketing Materials

```yaml
# Get marketing materials
GET /api/v1/affiliate/materials:
  authentication: required
  authorization: affiliate_member
  parameters:
    - material_type: string
    - campaign_id: string
    - target_audience: string
  responses:
    200:
      description: Available marketing materials

# Download marketing material
GET /api/v1/affiliate/materials/{materialId}/download:
  authentication: required
  authorization: affiliate_member
  responses:
    200:
      description: Material file download
    404:
      description: Material not found

# Track material usage
POST /api/v1/affiliate/materials/{materialId}/track:
  authentication: required
  authorization: affiliate_member
  body:
    type: object
    required: [usage_type]
    properties:
      usage_type:
        type: string
        enum: [view, click, download, share]
      platform:
        type: string
      additional_data:
        type: object
  responses:
    200:
      description: Usage tracked successfully
```

### Frontend Components

#### Affiliate Dashboard

```typescript
// Main affiliate dashboard component
export const AffiliateDashboard: React.FC = () => {
  const [dashboardData, setDashboardData] = useState<AffiliateDashboard | null>(null);
  const [activeTab, setActiveTab] = useState<'overview' | 'referrals' | 'commissions' | 'materials' | 'payouts'>('overview');
  const [dateRange, setDateRange] = useState<'7d' | '30d' | '90d' | '1y'>('30d');

  const affiliateTabs = [
    { id: 'overview', label: 'Overview', icon: <DashboardIcon /> },
    { id: 'referrals', label: 'Referrals', icon: <UsersIcon /> },
    { id: 'commissions', label: 'Commissions', icon: <MoneyIcon /> },
    { id: 'materials', label: 'Materials', icon: <ImageIcon /> },
    { id: 'payouts', label: 'Payouts', icon: <BankIcon /> },
  ];

  useEffect(() => {
    fetchDashboardData();
  }, [dateRange]);

  return (
    <div className="affiliate-dashboard">
      <AffiliateHeader 
        profile={dashboardData?.profile}
        onDateRangeChange={setDateRange}
      />
      
      <div className="dashboard-navigation">
        <TabNavigation 
          tabs={affiliateTabs}
          activeTab={activeTab}
          onTabChange={setActiveTab}
        />
      </div>
      
      <div className="dashboard-content">
        {activeTab === 'overview' && <AffiliateOverview data={dashboardData} />}
        {activeTab === 'referrals' && <ReferralManagement />}
        {activeTab === 'commissions' && <CommissionHistory />}
        {activeTab === 'materials' && <MarketingMaterials />}
        {activeTab === 'payouts' && <PayoutManagement />}
      </div>
    </div>
  );
};

// Affiliate overview component
export const AffiliateOverview: React.FC<{ data: AffiliateDashboard | null }> = ({ data }) => {
  if (!data) return <DashboardSkeleton />;

  return (
    <div className="affiliate-overview">
      <div className="earnings-summary">
        <div className="earnings-cards">
          <EarningsCard 
            title="Total Earnings"
            amount={data.profile.total_earnings}
            currency="NGN"
            icon={<WalletIcon />}
            trend={data.earnings_trend}
          />
          <EarningsCard 
            title="Pending Commissions"
            amount={data.profile.pending_earnings}
            currency="NGN"
            icon={<ClockIcon />}
            description="Awaiting approval"
          />
          <EarningsCard 
            title="This Month"
            amount={data.profile.current_month_earnings}
            currency="NGN"
            icon={<CalendarIcon />}
            comparison={data.profile.last_month_earnings}
          />
          <EarningsCard 
            title="Conversion Rate"
            amount={data.profile.conversion_rate}
            suffix="%"
            icon={<TrendIcon />}
            color="blue"
          />
        </div>
      </div>

      <div className="performance-metrics">
        <div className="metrics-grid">
          <MetricCard 
            title="Total Referrals"
            value={data.profile.total_referrals}
            icon={<UsersIcon />}
            change={data.referrals_change}
          />
          <MetricCard 
            title="Active Referrals"
            value={data.profile.active_referrals}
            icon={<UserCheckIcon />}
            subtitle="Currently active"
          />
          <MetricCard 
            title="Commission Tier"
            value={getTierName(data.profile.commission_tier)}
            icon={<StarIcon />}
            color="gold"
          />
          <MetricCard 
            title="Next Payout"
            value={formatDate(data.next_payout_date)}
            icon={<CalendarIcon />}
            subtitle="Estimated date"
          />
        </div>
      </div>

      <div className="quick-actions">
        <h3>Quick Actions</h3>
        <div className="action-buttons">
          <QuickActionButton 
            icon={<LinkIcon />}
            label="Generate Link"
            onClick={() => openLinkGenerator()}
          />
          <QuickActionButton 
            icon={<ShareIcon />}
            label="Share Materials"
            onClick={() => openMaterialsLibrary()}
          />
          <QuickActionButton 
            icon={<MoneyIcon />}
            label="Request Payout"
            onClick={() => openPayoutRequest()}
            disabled={data.profile.pending_earnings < 1000}
          />
          <QuickActionButton 
            icon={<ChartIcon />}
            label="View Analytics"
            onClick={() => openDetailedAnalytics()}
          />
        </div>
      </div>

      <div className="recent-activity">
        <div className="activity-section">
          <h3>Recent Referrals</h3>
          <RecentReferralsTable referrals={data.recent_referrals} />
        </div>
        
        <div className="activity-section">
          <h3>Performance Chart</h3>
          <PerformanceChart data={data.performance_chart} />
        </div>
      </div>

      <div className="affiliate-tools">
        <div className="link-generator">
          <h3>Quick Link Generator</h3>
          <LinkGeneratorForm onGenerate={handleLinkGeneration} />
        </div>
        
        <div className="tier-progress">
          <h3>Tier Progress</h3>
          <TierProgressIndicator 
            currentTier={data.profile.commission_tier}
            progress={data.tier_progress}
          />
        </div>
      </div>
    </div>
  );
};
```

#### Referral Management

```typescript
// Referral tracking and management
export const ReferralManagement: React.FC = () => {
  const [referrals, setReferrals] = useState<Referral[]>([]);
  const [filters, setFilters] = useState<ReferralFilters>({});
  const [selectedReferral, setSelectedReferral] = useState<Referral | null>(null);
  const [linkGenerator, setLinkGenerator] = useState(false);

  const referralColumns = [
    {
      key: 'referred_user',
      title: 'Referred User',
      render: (referral: Referral) => (
        <div className="user-info">
          <img src={referral.referred_user.avatar_url} alt="" />
          <div>
            <div className="user-name">{referral.referred_user.full_name}</div>
            <div className="user-email">{referral.referred_user.email}</div>
          </div>
        </div>
      ),
    },
    {
      key: 'conversion_type',
      title: 'Conversion',
      render: (referral: Referral) => (
        <div className="conversion-info">
          <div className={`conversion-badge ${referral.status}`}>
            {referral.conversion_type || 'Pending'}
          </div>
          {referral.conversion_value > 0 && (
            <div className="conversion-value">
              ₦{referral.conversion_value.toLocaleString()}
            </div>
          )}
        </div>
      ),
    },
    {
      key: 'commission',
      title: 'Commission',
      render: (referral: Referral) => (
        <div className="commission-info">
          <div className="commission-amount">
            ₦{referral.commission_amount.toLocaleString()}
          </div>
          <div className="commission-rate">
            {referral.commission_rate}% rate
          </div>
        </div>
      ),
    },
    {
      key: 'dates',
      title: 'Timeline',
      render: (referral: Referral) => (
        <div className="date-info">
          <div className="first-click">
            Clicked: {formatDate(referral.first_click_at)}
          </div>
          {referral.conversion_at && (
            <div className="conversion-date">
              Converted: {formatDate(referral.conversion_at)}
            </div>
          )}
        </div>
      ),
    },
    {
      key: 'actions',
      title: 'Actions',
      render: (referral: Referral) => (
        <div className="referral-actions">
          <button 
            onClick={() => setSelectedReferral(referral)}
            className="btn-sm btn-outline"
          >
            View Details
          </button>
          {referral.status === 'pending' && (
            <button 
              onClick={() => sendFollowUp(referral)}
              className="btn-sm btn-primary"
            >
              Follow Up
            </button>
          )}
        </div>
      ),
    },
  ];

  return (
    <div className="referral-management">
      <div className="referral-header">
        <h2>Referral Management</h2>
        <div className="header-actions">
          <button 
            onClick={() => setLinkGenerator(true)}
            className="btn-primary"
          >
            Generate New Link
          </button>
          <ReferralFilters 
            filters={filters}
            onFiltersChange={setFilters}
          />
        </div>
      </div>

      <div className="referral-stats">
        <ReferralStatsCards />
      </div>

      <div className="referral-table">
        <DataTable 
          columns={referralColumns}
          data={referrals}
          pagination={true}
          sorting={true}
          loading={loading}
        />
      </div>

      {selectedReferral && (
        <ReferralDetailModal 
          referral={selectedReferral}
          onClose={() => setSelectedReferral(null)}
        />
      )}

      {linkGenerator && (
        <LinkGeneratorModal 
          onClose={() => setLinkGenerator(false)}
          onGenerate={handleLinkGeneration}
        />
      )}
    </div>
  );
};
```

### Integration Points

#### Payment Integration
- Automated commission calculations
- Secure payout processing through multiple payment gateways
- Integration with wallet system for internal transfers

#### Analytics Integration
- Real-time performance tracking
- Conversion attribution
- ROI calculation for marketing efforts

#### User Management Integration
- Seamless affiliate profile creation
- Permission-based access control
- User activity correlation

---

*This feature specification provides the complete technical blueprint for implementing the Affiliate System within the Great Nigeria Library platform.*
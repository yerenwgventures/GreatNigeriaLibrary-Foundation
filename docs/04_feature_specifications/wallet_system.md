# Wallet System Feature Specification

**Document Version**: 1.0  
**Last Updated**: January 2025  
**Feature Owner**: Payment Team  
**Status**: Implemented

---

## Overview

The Wallet System provides a comprehensive digital wallet solution within the Great Nigeria Library platform, enabling users to store, transfer, and manage digital currency (points and Naira equivalents). The system supports both earned points from platform activities and purchased credits for marketplace transactions.

## Feature Purpose

### Goals
1. **Digital Economy**: Create a seamless digital payment experience within the platform
2. **Points Management**: Provide a central hub for earning, storing, and redeeming points
3. **Secure Transactions**: Ensure safe and transparent financial transactions
4. **Multi-Currency Support**: Support both points and real currency (NGN)
5. **User Empowerment**: Give users control over their digital assets

### Success Metrics
- **User Adoption**: 80%+ of active users with funded wallets
- **Transaction Volume**: ₦50M+ in wallet transactions annually
- **Security**: Zero successful fraud attempts
- **User Satisfaction**: 4.8+ wallet experience rating
- **Transaction Speed**: <2 seconds for internal transfers

## System Architecture

### Core Wallet Components

#### Multi-Currency Wallet Management
Comprehensive digital wallet system supporting multiple currencies:

- **Points Balance Management**: Earned points from platform activities with detailed tracking and history
- **Naira Balance Support**: Real Nigerian currency storage and management with secure transactions
- **Multi-Currency Capability**: Support for USD and other international currencies for global users
- **Frozen Funds Management**: Secure fund freezing for pending transactions and dispute resolution
- **Balance Tracking**: Comprehensive tracking of earned, spent, deposited, and withdrawn amounts
- **Security Features**: PIN-based wallet protection with attempt limiting and security monitoring

#### Transaction Processing Engine
Advanced transaction management system:

- **Real-Time Processing**: Instant transaction processing for internal platform transfers
- **Transaction History**: Detailed transaction logs with comprehensive metadata and tracking
- **Batch Processing**: Efficient batch processing for bulk transactions and automated payments
- **Transaction Validation**: Multi-layer validation system ensuring transaction integrity and security
- **Rollback Capability**: Secure transaction rollback for failed or disputed transactions
- **Audit Trail**: Complete audit trail for all wallet activities and administrative actions

#### Payment Gateway Integration
Comprehensive payment processing with multiple Nigerian providers:

- **Multi-Provider Support**: Integration with Paystack, Flutterwave, and Squad for diverse payment options
- **Deposit Management**: Seamless fund deposits with real-time processing and confirmation
- **Withdrawal Processing**: Secure withdrawal system with bank account verification and processing
- **Fee Management**: Transparent fee calculation and automatic deduction for all transactions
- **Provider Redundancy**: Automatic failover between payment providers for maximum uptime
- **Callback Handling**: Robust webhook processing for real-time transaction status updates

#### User-to-User Transfer System
Peer-to-peer transfer capabilities within the platform:

- **Instant Transfers**: Real-time transfers between platform users with immediate balance updates
- **Transfer Limits**: Configurable daily and monthly transfer limits based on KYC verification levels
- **PIN Security**: Mandatory PIN verification for all outgoing transfers and sensitive operations
- **Transfer History**: Comprehensive transfer history with detailed transaction records
- **Recipient Verification**: User verification system to prevent accidental transfers
- **Transfer Fees**: Configurable transfer fee structure with transparent cost display

#### Security and Compliance Framework
Advanced security measures and regulatory compliance:

- **KYC Integration**: Multi-level Know Your Customer verification with document upload and verification
- **Spending Limits**: Configurable daily and monthly spending limits based on verification status
- **Security Monitoring**: Real-time fraud detection and suspicious activity monitoring
- **PIN Management**: Secure PIN creation, change, and recovery with attempt limiting
- **Account Freezing**: Administrative tools for account suspension and fund freezing when necessary
- **Compliance Reporting**: Automated compliance reporting for regulatory requirements

#### Scheduled Transaction Management
Automated and recurring transaction capabilities:
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID REFERENCES wallet_accounts(id) ON DELETE CASCADE,
    transaction_type VARCHAR(50) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    description TEXT NOT NULL,
    schedule_type VARCHAR(20) CHECK (schedule_type IN ('daily', 'weekly', 'monthly', 'quarterly')),
    next_execution TIMESTAMP WITH TIME ZONE NOT NULL,
    last_execution TIMESTAMP WITH TIME ZONE,
    execution_count INTEGER DEFAULT 0,
    max_executions INTEGER,
    is_active BOOLEAN DEFAULT TRUE,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Points earning rules and activities
CREATE TABLE wallet_points_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID REFERENCES wallet_accounts(id) ON DELETE CASCADE,
    activity_type VARCHAR(50) NOT NULL, -- reading, commenting, sharing, quiz, daily_login
    activity_reference VARCHAR(255),
    points_earned DECIMAL(10,2) NOT NULL,
    multiplier DECIMAL(5,2) DEFAULT 1.0,
    bonus_points DECIMAL(10,2) DEFAULT 0,
    bonus_reason VARCHAR(255),
    earned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Wallet security logs
CREATE TABLE wallet_security_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_id UUID REFERENCES wallet_accounts(id) ON DELETE CASCADE,
    action VARCHAR(100) NOT NULL,
    ip_address INET,
    user_agent TEXT,
    device_fingerprint VARCHAR(255),
    location_data JSONB,
    status VARCHAR(20) CHECK (status IN ('success', 'failed', 'suspicious')),
    failure_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes for performance
CREATE INDEX idx_wallet_transactions_wallet_id ON wallet_transactions(wallet_id);
CREATE INDEX idx_wallet_transactions_created_at ON wallet_transactions(created_at DESC);
CREATE INDEX idx_wallet_transactions_type ON wallet_transactions(transaction_type);
CREATE INDEX idx_wallet_transactions_status ON wallet_transactions(status);
CREATE INDEX idx_wallet_payment_transactions_provider_ref ON wallet_payment_transactions(provider_reference);
CREATE INDEX idx_wallet_transfers_sender ON wallet_transfers(sender_wallet_id);
CREATE INDEX idx_wallet_transfers_recipient ON wallet_transfers(recipient_wallet_id);
CREATE INDEX idx_wallet_points_activities_wallet_id ON wallet_points_activities(wallet_id);
CREATE INDEX idx_wallet_security_logs_wallet_id ON wallet_security_logs(wallet_id);
```

### API Endpoints

#### Wallet Management

```yaml
# Get wallet balance and details
GET /api/v1/wallet:
  authentication: required
  responses:
    200:
      description: Wallet details with balances
      schema:
        type: object
        properties:
          id:
            type: string
            format: uuid
          user_id:
            type: string
            format: uuid
          points_balance:
            type: number
          naira_balance:
            type: number
          usd_balance:
            type: number
          frozen_points:
            type: number
          frozen_naira:
            type: number
          total_earned_points:
            type: number
          daily_spend_limit:
            type: number
          monthly_spend_limit:
            type: number
          kyc_level:
            type: integer

# Get transaction history
GET /api/v1/wallet/transactions:
  authentication: required
  parameters:
    - page: integer
    - limit: integer
    - type: string
    - currency: string
    - from_date: string (ISO date)
    - to_date: string (ISO date)
  responses:
    200:
      description: Paginated transaction history

# Create wallet PIN
POST /api/v1/wallet/pin:
  authentication: required
  body:
    type: object
    required: [pin, confirm_pin]
    properties:
      pin:
        type: string
        pattern: '^[0-9]{4,6}$'
      confirm_pin:
        type: string
        pattern: '^[0-9]{4,6}$'
  responses:
    200:
      description: PIN created successfully

# Verify wallet PIN
POST /api/v1/wallet/pin/verify:
  authentication: required
  body:
    type: object
    required: [pin]
    properties:
      pin:
        type: string
        pattern: '^[0-9]{4,6}$'
  responses:
    200:
      description: PIN verified successfully
    400:
      description: Invalid PIN
    423:
      description: Wallet locked due to too many attempts

# Change wallet PIN
PUT /api/v1/wallet/pin:
  authentication: required
  body:
    type: object
    required: [current_pin, new_pin, confirm_new_pin]
    properties:
      current_pin:
        type: string
      new_pin:
        type: string
        pattern: '^[0-9]{4,6}$'
      confirm_new_pin:
        type: string
        pattern: '^[0-9]{4,6}$'
  responses:
    200:
      description: PIN changed successfully
```

#### Money Transfers

```yaml
# Transfer funds to another user
POST /api/v1/wallet/transfer:
  authentication: required
  body:
    type: object
    required: [recipient_identifier, amount, currency, pin]
    properties:
      recipient_identifier:
        type: string
        description: "Email, phone, or username"
      amount:
        type: number
        minimum: 1
      currency:
        type: string
        enum: [points, NGN, USD]
      pin:
        type: string
      description:
        type: string
        maxLength: 500
  responses:
    200:
      description: Transfer initiated successfully
    400:
      description: Insufficient funds or invalid recipient
    423:
      description: Wallet locked

# Get transfer history
GET /api/v1/wallet/transfers:
  authentication: required
  parameters:
    - type: string (sent|received|all)
    - status: string
    - page: integer
    - limit: integer
  responses:
    200:
      description: Transfer history

# Get transfer details
GET /api/v1/wallet/transfers/{transferId}:
  authentication: required
  responses:
    200:
      description: Transfer details
    404:
      description: Transfer not found
```

#### Deposits and Withdrawals

```yaml
# Initiate deposit
POST /api/v1/wallet/deposit:
  authentication: required
  body:
    type: object
    required: [amount, currency, payment_method]
    properties:
      amount:
        type: number
        minimum: 100
        maximum: 1000000
      currency:
        type: string
        enum: [NGN, USD]
      payment_method:
        type: string
        enum: [card, bank_transfer, ussd]
      callback_url:
        type: string
        format: url
  responses:
    200:
      description: Deposit initiated, returns payment details

# Confirm deposit (webhook handler)
POST /api/v1/wallet/deposit/confirm:
  authentication: webhook_signature
  body:
    type: object
    # Provider-specific webhook payload
  responses:
    200:
      description: Deposit processed

# Initiate withdrawal
POST /api/v1/wallet/withdraw:
  authentication: required
  body:
    type: object
    required: [amount, currency, bank_details, pin]
    properties:
      amount:
        type: number
        minimum: 500
      currency:
        type: string
        enum: [NGN, USD]
      bank_details:
        type: object
        required: [account_number, bank_code, account_name]
        properties:
          account_number:
            type: string
            pattern: '^[0-9]{10}$'
          bank_code:
            type: string
          account_name:
            type: string
      pin:
        type: string
      reason:
        type: string
  responses:
    200:
      description: Withdrawal request submitted
    400:
      description: Insufficient funds or invalid bank details

# Get withdrawal status
GET /api/v1/wallet/withdrawals/{withdrawalId}:
  authentication: required
  responses:
    200:
      description: Withdrawal status and details
```

#### Points System Integration

```yaml
# Earn points from activity
POST /api/v1/wallet/points/earn:
  authentication: required
  body:
    type: object
    required: [activity_type, activity_reference, points_amount]
    properties:
      activity_type:
        type: string
        enum: [reading, commenting, sharing, quiz, daily_login, referral]
      activity_reference:
        type: string
      points_amount:
        type: number
        minimum: 1
      multiplier:
        type: number
        default: 1.0
      bonus_points:
        type: number
        default: 0
      bonus_reason:
        type: string
  responses:
    200:
      description: Points awarded successfully

# Redeem points
POST /api/v1/wallet/points/redeem:
  authentication: required
  body:
    type: object
    required: [redemption_type, points_amount, pin]
    properties:
      redemption_type:
        type: string
        enum: [cash_out, marketplace_credit, premium_access]
      points_amount:
        type: number
        minimum: 100
      pin:
        type: string
      redemption_details:
        type: object
  responses:
    200:
      description: Points redeemed successfully

# Get points earning history
GET /api/v1/wallet/points/history:
  authentication: required
  parameters:
    - activity_type: string
    - from_date: string
    - to_date: string
    - page: integer
    - limit: integer
  responses:
    200:
      description: Points earning history

# Check points earning eligibility
GET /api/v1/wallet/points/eligibility:
  authentication: required
  parameters:
    - activity_type: string
    - activity_reference: string
  responses:
    200:
      description: Eligibility status and potential points
```

### Frontend Components

#### Wallet Dashboard

```typescript
// Main wallet page component
export const WalletPage: React.FC = () => {
  const [wallet, setWallet] = useState<WalletAccount | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [activeTab, setActiveTab] = useState<'overview' | 'transactions' | 'transfer' | 'deposit' | 'withdraw'>('overview');
  const [loading, setLoading] = useState(true);

  const walletTabs = [
    { id: 'overview', label: 'Overview', icon: <DashboardIcon /> },
    { id: 'transactions', label: 'Transactions', icon: <HistoryIcon /> },
    { id: 'transfer', label: 'Transfer', icon: <SendIcon /> },
    { id: 'deposit', label: 'Add Money', icon: <AddIcon /> },
    { id: 'withdraw', label: 'Withdraw', icon: <DownloadIcon /> },
  ];

  return (
    <div className="wallet-page">
      <WalletHeader wallet={wallet} />
      
      <div className="wallet-navigation">
        <TabNavigation 
          tabs={walletTabs}
          activeTab={activeTab}
          onTabChange={setActiveTab}
        />
      </div>
      
      <div className="wallet-content">
        {activeTab === 'overview' && <WalletOverview wallet={wallet} />}
        {activeTab === 'transactions' && <TransactionHistory transactions={transactions} />}
        {activeTab === 'transfer' && <MoneyTransfer wallet={wallet} />}
        {activeTab === 'deposit' && <DepositMoney wallet={wallet} />}
        {activeTab === 'withdraw' && <WithdrawMoney wallet={wallet} />}
      </div>
    </div>
  );
};

// Wallet overview component
export const WalletOverview: React.FC<{ wallet: WalletAccount | null }> = ({ wallet }) => {
  if (!wallet) return <WalletSkeleton />;

  return (
    <div className="wallet-overview">
      <div className="balance-cards">
        <BalanceCard 
          title="Points Balance"
          amount={wallet.points_balance}
          currency="points"
          icon={<StarIcon />}
          color="purple"
        />
        <BalanceCard 
          title="Naira Balance"
          amount={wallet.naira_balance}
          currency="NGN"
          icon={<NairaIcon />}
          color="green"
        />
        {wallet.usd_balance > 0 && (
          <BalanceCard 
            title="USD Balance"
            amount={wallet.usd_balance}
            currency="USD"
            icon={<DollarIcon />}
            color="blue"
          />
        )}
      </div>
      
      <div className="quick-actions">
        <h3>Quick Actions</h3>
        <div className="action-buttons">
          <QuickActionButton 
            icon={<SendIcon />}
            label="Send Money"
            onClick={() => setActiveTab('transfer')}
          />
          <QuickActionButton 
            icon={<AddIcon />}
            label="Add Money"
            onClick={() => setActiveTab('deposit')}
          />
          <QuickActionButton 
            icon={<DownloadIcon />}
            label="Withdraw"
            onClick={() => setActiveTab('withdraw')}
          />
          <QuickActionButton 
            icon={<GiftIcon />}
            label="Redeem Points"
            onClick={() => openPointsRedemption()}
          />
        </div>
      </div>
      
      <div className="recent-activity">
        <h3>Recent Activity</h3>
        <RecentTransactions limit={5} />
      </div>
      
      <div className="wallet-stats">
        <StatCard 
          title="Total Earned Points"
          value={wallet.total_earned_points}
          subtitle="All time points earned"
        />
        <StatCard 
          title="This Month Spending"
          value={calculateMonthlySpending(wallet)}
          subtitle={`Limit: ₦${wallet.monthly_spend_limit.toLocaleString()}`}
        />
        <StatCard 
          title="KYC Level"
          value={wallet.kyc_level}
          subtitle={getKYCLevelDescription(wallet.kyc_level)}
        />
      </div>
    </div>
  );
};
```

#### Money Transfer Interface

```typescript
// Money transfer component
export const MoneyTransfer: React.FC<{ wallet: WalletAccount }> = ({ wallet }) => {
  const [recipientType, setRecipientType] = useState<'email' | 'phone' | 'username'>('email');
  const [recipient, setRecipient] = useState('');
  const [amount, setAmount] = useState('');
  const [currency, setCurrency] = useState<'points' | 'NGN'>('NGN');
  const [description, setDescription] = useState('');
  const [pin, setPin] = useState('');
  const [loading, setLoading] = useState(false);
  const [recipientInfo, setRecipientInfo] = useState<User | null>(null);

  const handleRecipientLookup = async () => {
    try {
      const result = await walletService.lookupRecipient(recipientType, recipient);
      setRecipientInfo(result);
    } catch (error) {
      setRecipientInfo(null);
      console.error('Recipient lookup failed:', error);
    }
  };

  const handleTransfer = async () => {
    if (!recipientInfo || !amount || !pin) return;

    setLoading(true);
    try {
      const transfer = await walletService.initiateTransfer({
        recipient_identifier: recipient,
        amount: parseFloat(amount),
        currency,
        description,
        pin,
      });

      // Show success message and reset form
      toast.success('Transfer completed successfully!');
      resetForm();
    } catch (error) {
      toast.error(error.message || 'Transfer failed');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="money-transfer">
      <div className="transfer-form">
        <h2>Send Money</h2>
        
        <div className="recipient-section">
          <label>Send to:</label>
          <div className="recipient-type-selector">
            <button 
              className={recipientType === 'email' ? 'active' : ''}
              onClick={() => setRecipientType('email')}
            >
              Email
            </button>
            <button 
              className={recipientType === 'phone' ? 'active' : ''}
              onClick={() => setRecipientType('phone')}
            >
              Phone
            </button>
            <button 
              className={recipientType === 'username' ? 'active' : ''}
              onClick={() => setRecipientType('username')}
            >
              Username
            </button>
          </div>
          
          <div className="recipient-input">
            <input 
              type={recipientType === 'email' ? 'email' : 'text'}
              placeholder={`Enter ${recipientType}`}
              value={recipient}
              onChange={(e) => setRecipient(e.target.value)}
              onBlur={handleRecipientLookup}
            />
            {recipientInfo && (
              <div className="recipient-info">
                <img src={recipientInfo.avatar_url} alt={recipientInfo.full_name} />
                <span>{recipientInfo.full_name}</span>
                <VerifiedIcon />
              </div>
            )}
          </div>
        </div>
        
        <div className="amount-section">
          <label>Amount:</label>
          <div className="amount-input">
            <CurrencySelector 
              value={currency}
              onChange={setCurrency}
              balances={{
                points: wallet.points_balance,
                NGN: wallet.naira_balance,
              }}
            />
            <input 
              type="number"
              placeholder="0.00"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              min="1"
              max={currency === 'points' ? wallet.points_balance : wallet.naira_balance}
            />
          </div>
          <div className="available-balance">
            Available: {currency === 'points' ? 
              `${wallet.points_balance.toLocaleString()} points` : 
              `₦${wallet.naira_balance.toLocaleString()}`
            }
          </div>
        </div>
        
        <div className="description-section">
          <label>Description (optional):</label>
          <textarea 
            placeholder="What's this transfer for?"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            maxLength={500}
          />
        </div>
        
        <div className="security-section">
          <label>Enter your wallet PIN:</label>
          <PinInput 
            value={pin}
            onChange={setPin}
            masked={true}
            length={4}
          />
        </div>
        
        <div className="transfer-summary">
          <div className="summary-row">
            <span>Amount:</span>
            <span>
              {currency === 'points' ? 
                `${parseFloat(amount || '0').toLocaleString()} points` : 
                `₦${parseFloat(amount || '0').toLocaleString()}`
              }
            </span>
          </div>
          <div className="summary-row">
            <span>Transfer Fee:</span>
            <span>₦0.00</span>
          </div>
          <div className="summary-row total">
            <span>Total:</span>
            <span>
              {currency === 'points' ? 
                `${parseFloat(amount || '0').toLocaleString()} points` : 
                `₦${parseFloat(amount || '0').toLocaleString()}`
              }
            </span>
          </div>
        </div>
        
        <button 
          className="transfer-button"
          onClick={handleTransfer}
          disabled={loading || !recipientInfo || !amount || !pin}
        >
          {loading ? 'Processing...' : 'Send Money'}
        </button>
      </div>
    </div>
  );
};
```

### Security Features

#### PIN Management
- 4-6 digit PIN for transaction authorization
- PIN hashing with bcrypt
- Failed attempt tracking and temporary lockout
- PIN change functionality with current PIN verification

#### Transaction Monitoring
- Real-time fraud detection algorithms
- Unusual activity pattern detection
- IP address and device tracking
- Automated alerts for suspicious activities

#### KYC Integration
- Multi-level KYC verification
- Higher limits for verified users
- BVN/NIN verification integration
- Document upload and verification

---

*This feature specification provides the complete technical blueprint for implementing the Wallet System within the Great Nigeria Library platform.*
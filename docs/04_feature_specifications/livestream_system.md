# Livestream System Feature Specification

**Document Version**: 1.0  
**Last Updated**: January 2025  
**Feature Owner**: Media Team  
**Status**: Implemented

---

## Overview

The Livestream System provides comprehensive live streaming capabilities within the Great Nigeria Library platform, enabling users to broadcast educational content, host discussions, and engage with their audience in real-time. The system includes TikTok-style virtual gifting and monetization features.

## Feature Purpose

### Goals
1. **Educational Broadcasting**: Enable live educational content delivery
2. **Community Engagement**: Foster real-time interaction between creators and audience
3. **Creator Monetization**: Provide revenue streams through virtual gifts and subscriptions
4. **Cultural Celebration**: Support Nigerian cultural events and celebrations
5. **Knowledge Sharing**: Facilitate live Q&A sessions and interactive learning

### Success Metrics
- **Creator Adoption**: 1,000+ active streamers by end of Year 1
- **Viewer Engagement**: 70%+ concurrent viewership during peak hours
- **Revenue Generation**: ₦5M+ in virtual gift transactions monthly
- **Technical Performance**: 99.9% stream uptime with <2 second latency
- **Content Quality**: 4.5+ average stream rating

## Technical Architecture

### Live Streaming Platform Infrastructure
Comprehensive live streaming system with advanced creator economy features:
-- Live stream sessions
CREATE TABLE live_streams (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    thumbnail_url TEXT,
    stream_key VARCHAR(255) UNIQUE NOT NULL,
    rtmp_url TEXT NOT NULL,
    hls_url TEXT,
    stream_status VARCHAR(20) DEFAULT 'offline' CHECK (stream_status IN ('offline', 'live', 'ended', 'scheduled')),
    category VARCHAR(100),
    tags TEXT[],
    is_featured BOOLEAN DEFAULT FALSE,
    is_premium BOOLEAN DEFAULT FALSE,
    scheduled_start_time TIMESTAMP WITH TIME ZONE,
    actual_start_time TIMESTAMP WITH TIME ZONE,
    end_time TIMESTAMP WITH TIME ZONE,
    max_viewers INTEGER DEFAULT 0,
    current_viewers INTEGER DEFAULT 0,
    total_views INTEGER DEFAULT 0,
    total_watch_time_seconds BIGINT DEFAULT 0,
    total_gifts_received DECIMAL(15,2) DEFAULT 0,
    total_gift_count INTEGER DEFAULT 0,
    language VARCHAR(10) DEFAULT 'en',
    visibility VARCHAR(20) DEFAULT 'public' CHECK (visibility IN ('public', 'unlisted', 'private')),
    allow_chat BOOLEAN DEFAULT TRUE,
    allow_gifts BOOLEAN DEFAULT TRUE,
    moderate_chat BOOLEAN DEFAULT FALSE,
    record_stream BOOLEAN DEFAULT FALSE,
    recording_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Stream viewers tracking
CREATE TABLE live_stream_viewers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID REFERENCES live_streams(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    left_at TIMESTAMP WITH TIME ZONE,
    watch_duration_seconds INTEGER DEFAULT 0,
    ip_address INET,
    user_agent TEXT,
    is_anonymous BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Live chat messages
CREATE TABLE live_stream_chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID REFERENCES live_streams(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    message_type VARCHAR(20) DEFAULT 'text' CHECK (message_type IN ('text', 'emote', 'system', 'gift')),
    content TEXT NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    is_highlighted BOOLEAN DEFAULT FALSE,
    mentions UUID[], -- array of user IDs mentioned
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Virtual gifts system
CREATE TABLE virtual_gifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    image_url TEXT NOT NULL,
    animation_url TEXT,
    sound_url TEXT,
    cost_in_coins INTEGER NOT NULL,
    naira_value DECIMAL(10,2) NOT NULL,
    category VARCHAR(50), -- traditional, modern, premium, exclusive
    rarity VARCHAR(20) DEFAULT 'common' CHECK (rarity IN ('common', 'rare', 'epic', 'legendary')),
    effects JSONB, -- animation effects config
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Virtual currency (coins) packages
CREATE TABLE virtual_currency_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    coin_amount INTEGER NOT NULL,
    naira_price DECIMAL(10,2) NOT NULL,
    usd_price DECIMAL(10,2),
    bonus_coins INTEGER DEFAULT 0,
    is_featured BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User virtual currency balances
CREATE TABLE virtual_currency_balances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE UNIQUE,
    coin_balance INTEGER DEFAULT 0,
    total_coins_purchased INTEGER DEFAULT 0,
    total_coins_spent INTEGER DEFAULT 0,
    total_coins_earned INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Gift transactions
CREATE TABLE live_stream_gifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID REFERENCES live_streams(id) ON DELETE CASCADE,
    gift_id UUID REFERENCES virtual_gifts(id),
    sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    recipient_id UUID REFERENCES users(id) ON DELETE CASCADE,
    quantity INTEGER DEFAULT 1,
    total_cost_coins INTEGER NOT NULL,
    total_naira_value DECIMAL(10,2) NOT NULL,
    message TEXT,
    is_anonymous BOOLEAN DEFAULT FALSE,
    combo_count INTEGER DEFAULT 1,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Creator revenue tracking
CREATE TABLE creator_revenues (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id UUID REFERENCES users(id) ON DELETE CASCADE,
    revenue_type VARCHAR(50) NOT NULL, -- gifts, subscriptions, ads
    source_id UUID, -- reference to gift, subscription, etc.
    amount_naira DECIMAL(10,2) NOT NULL,
    platform_fee DECIMAL(10,2) NOT NULL,
    creator_share DECIMAL(10,2) NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'paid')),
    payout_date TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Stream categories
CREATE TABLE stream_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    icon_url TEXT,
    color_code VARCHAR(7), -- hex color
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Stream reports and moderation
CREATE TABLE stream_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID REFERENCES live_streams(id) ON DELETE CASCADE,
    reporter_id UUID REFERENCES users(id) ON DELETE SET NULL,
    report_type VARCHAR(50) NOT NULL, -- inappropriate_content, spam, harassment, copyright
    description TEXT,
    timestamp_in_stream INTEGER, -- seconds from stream start
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'reviewing', 'resolved', 'dismissed')),
    moderator_id UUID REFERENCES users(id) ON DELETE SET NULL,
    moderator_notes TEXT,
    action_taken VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Stream analytics
CREATE TABLE stream_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    stream_id UUID REFERENCES live_streams(id) ON DELETE CASCADE,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL(15,2) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    metadata JSONB
);

-- Follower/subscription system for creators
CREATE TABLE creator_followers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    creator_id UUID REFERENCES users(id) ON DELETE CASCADE,
    follower_id UUID REFERENCES users(id) ON DELETE CASCADE,
    notification_enabled BOOLEAN DEFAULT TRUE,
    followed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(creator_id, follower_id)
);

-- Indexes for performance
CREATE INDEX idx_live_streams_creator_id ON live_streams(creator_id);
CREATE INDEX idx_live_streams_status ON live_streams(stream_status);
CREATE INDEX idx_live_streams_created_at ON live_streams(created_at DESC);
CREATE INDEX idx_live_stream_viewers_stream_id ON live_stream_viewers(stream_id);
CREATE INDEX idx_live_stream_chat_stream_id ON live_stream_chat_messages(stream_id);
CREATE INDEX idx_live_stream_chat_timestamp ON live_stream_chat_messages(timestamp DESC);
CREATE INDEX idx_live_stream_gifts_stream_id ON live_stream_gifts(stream_id);
CREATE INDEX idx_live_stream_gifts_sender_id ON live_stream_gifts(sender_id);
CREATE INDEX idx_live_stream_gifts_timestamp ON live_stream_gifts(timestamp DESC);
CREATE INDEX idx_creator_revenues_creator_id ON creator_revenues(creator_id);
CREATE INDEX idx_creator_followers_creator_id ON creator_followers(creator_id);
CREATE INDEX idx_creator_followers_follower_id ON creator_followers(follower_id);
```

### API Endpoints

#### Stream Management

```yaml
# Create new stream
POST /api/v1/streams:
  authentication: required
  body:
    type: object
    required: [title]
    properties:
      title:
        type: string
        maxLength: 255
      description:
        type: string
        maxLength: 2000
      category:
        type: string
      tags:
        type: array
        items:
          type: string
      scheduled_start_time:
        type: string
        format: date-time
      is_premium:
        type: boolean
      allow_chat:
        type: boolean
      allow_gifts:
        type: boolean
      visibility:
        type: string
        enum: [public, unlisted, private]
  responses:
    201:
      description: Stream created with RTMP details

# Get stream details
GET /api/v1/streams/{streamId}:
  responses:
    200:
      description: Stream details including stats
    404:
      description: Stream not found

# Update stream
PUT /api/v1/streams/{streamId}:
  authentication: required
  authorization: creator_owns_stream
  body:
    type: object
    properties:
      title:
        type: string
      description:
        type: string
      category:
        type: string
      tags:
        type: array
      # ... other updatable fields
  responses:
    200:
      description: Stream updated successfully

# Start/stop stream
POST /api/v1/streams/{streamId}/status:
  authentication: required
  authorization: creator_owns_stream
  body:
    type: object
    required: [action]
    properties:
      action:
        type: string
        enum: [start, stop]
  responses:
    200:
      description: Stream status updated

# List active streams
GET /api/v1/streams:
  parameters:
    - status: string (live|scheduled|all)
    - category: string
    - creator: string
    - page: integer
    - limit: integer
    - sort: string (viewers|newest|trending)
  responses:
    200:
      description: List of streams with pagination
```

#### Live Chat

```yaml
# WebSocket connection for real-time chat
WS /api/v1/streams/{streamId}/chat:
  authentication: optional
  description: Real-time bidirectional chat communication

# Send chat message
POST /api/v1/streams/{streamId}/chat:
  authentication: required
  body:
    type: object
    required: [content]
    properties:
      content:
        type: string
        maxLength: 500
      message_type:
        type: string
        enum: [text, emote]
      mentions:
        type: array
        items:
          type: string
          format: uuid
  responses:
    201:
      description: Message sent successfully

# Get chat history
GET /api/v1/streams/{streamId}/chat:
  parameters:
    - before: string (timestamp)
    - limit: integer
  responses:
    200:
      description: Chat message history

# Delete chat message (moderator)
DELETE /api/v1/streams/{streamId}/chat/{messageId}:
  authentication: required
  authorization: creator_or_moderator
  responses:
    204:
      description: Message deleted
```

#### Virtual Gifts

```yaml
# Get available gifts
GET /api/v1/gifts:
  parameters:
    - category: string
    - rarity: string
  responses:
    200:
      description: List of available virtual gifts

# Send gift to stream
POST /api/v1/streams/{streamId}/gifts:
  authentication: required
  body:
    type: object
    required: [gift_id, quantity]
    properties:
      gift_id:
        type: string
        format: uuid
      quantity:
        type: integer
        minimum: 1
        maximum: 100
      message:
        type: string
        maxLength: 200
      is_anonymous:
        type: boolean
  responses:
    201:
      description: Gift sent successfully
    400:
      description: Insufficient coins or invalid gift

# Get gift history for stream
GET /api/v1/streams/{streamId}/gifts:
  parameters:
    - limit: integer
  responses:
    200:
      description: Recent gifts sent to stream

# Get user's gift history
GET /api/v1/users/{userId}/gifts:
  authentication: required
  parameters:
    - type: string (sent|received)
    - page: integer
    - limit: integer
  responses:
    200:
      description: User's gift transaction history
```

#### Virtual Currency

```yaml
# Get user's coin balance
GET /api/v1/currency/balance:
  authentication: required
  responses:
    200:
      description: Current coin balance and statistics

# Get coin packages
GET /api/v1/currency/packages:
  responses:
    200:
      description: Available coin packages for purchase

# Purchase coins
POST /api/v1/currency/purchase:
  authentication: required
  body:
    type: object
    required: [package_id, payment_method]
    properties:
      package_id:
        type: string
        format: uuid
      payment_method:
        type: string
        enum: [card, bank_transfer, wallet]
      callback_url:
        type: string
        format: url
  responses:
    200:
      description: Purchase initiated, returns payment details

# Get coin transaction history
GET /api/v1/currency/transactions:
  authentication: required
  parameters:
    - type: string (purchase|gift|earn)
    - page: integer
    - limit: integer
  responses:
    200:
      description: Coin transaction history
```

### Frontend Components

#### Live Stream Player

```typescript
// Main stream viewing component
export const LiveStreamPlayer: React.FC<{ streamId: string }> = ({ streamId }) => {
  const [stream, setStream] = useState<LiveStream | null>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [volume, setVolume] = useState(1);
  const [fullscreen, setFullscreen] = useState(false);
  const [chatVisible, setChatVisible] = useState(true);
  const [giftPanelOpen, setGiftPanelOpen] = useState(false);
  
  const videoRef = useRef<HTMLVideoElement>(null);
  const chatRef = useRef<WebSocket>(null);

  useEffect(() => {
    // Initialize HLS player
    if (stream?.hls_url && videoRef.current) {
      const video = videoRef.current;
      if (Hls.isSupported()) {
        const hls = new Hls();
        hls.loadSource(stream.hls_url);
        hls.attachMedia(video);
      } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
        video.src = stream.hls_url;
      }
    }

    // Initialize WebSocket for chat
    chatRef.current = new WebSocket(`ws://localhost:5000/api/v1/streams/${streamId}/chat`);
    chatRef.current.onmessage = handleChatMessage;

    return () => {
      chatRef.current?.close();
    };
  }, [stream, streamId]);

  const handleGiftSend = async (giftId: string, quantity: number, message?: string) => {
    try {
      await livestreamService.sendGift(streamId, {
        gift_id: giftId,
        quantity,
        message,
      });
      toast.success('Gift sent successfully!');
    } catch (error) {
      toast.error('Failed to send gift');
    }
  };

  return (
    <div className={`stream-player ${fullscreen ? 'fullscreen' : ''}`}>
      <div className="video-container">
        <video
          ref={videoRef}
          controls
          autoPlay
          muted={false}
          volume={volume}
          onPlay={() => setIsPlaying(true)}
          onPause={() => setIsPlaying(false)}
        />
        
        <div className="video-overlay">
          <StreamInfo stream={stream} />
          <ViewerCount count={stream?.current_viewers || 0} />
          <StreamControls 
            onToggleFullscreen={() => setFullscreen(!fullscreen)}
            onToggleChat={() => setChatVisible(!chatVisible)}
            onOpenGifts={() => setGiftPanelOpen(true)}
          />
        </div>
        
        {stream?.stream_status === 'offline' && (
          <div className="offline-overlay">
            <h3>Stream is currently offline</h3>
            <p>Check back later or follow the creator for notifications</p>
          </div>
        )}
      </div>

      <div className={`chat-section ${chatVisible ? 'visible' : 'hidden'}`}>
        <LiveChat 
          streamId={streamId} 
          websocket={chatRef.current}
          onSendMessage={handleChatMessage}
        />
      </div>

      {giftPanelOpen && (
        <VirtualGiftPanel 
          onSendGift={handleGiftSend}
          onClose={() => setGiftPanelOpen(false)}
        />
      )}
    </div>
  );
};

// Live chat component
export const LiveChat: React.FC<{
  streamId: string;
  websocket: WebSocket | null;
  onSendMessage: (message: string) => void;
}> = ({ streamId, websocket, onSendMessage }) => {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [chatSettings, setChatSettings] = useState({
    emoteOnly: false,
    slowMode: 0,
    followersOnly: false,
  });

  const messagesEndRef = useRef<HTMLDivElement>(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const handleSendMessage = () => {
    if (!newMessage.trim() || !websocket) return;

    const message = {
      type: 'chat_message',
      content: newMessage,
      stream_id: streamId,
    };

    websocket.send(JSON.stringify(message));
    setNewMessage('');
  };

  const handleGiftReceived = (giftData: GiftMessage) => {
    // Add gift animation to chat
    setMessages(prev => [...prev, {
      id: generateId(),
      type: 'gift',
      content: `sent ${giftData.quantity}x ${giftData.gift_name}`,
      user: giftData.sender,
      timestamp: new Date(),
      gift_data: giftData,
    }]);
  };

  return (
    <div className="live-chat">
      <div className="chat-header">
        <h4>Live Chat</h4>
        <div className="chat-viewer-count">
          {messages.length} viewers
        </div>
      </div>

      <div className="chat-messages">
        {messages.map(message => (
          <ChatMessage 
            key={message.id}
            message={message}
            onUserClick={handleUserClick}
          />
        ))}
        <div ref={messagesEndRef} />
      </div>

      <div className="chat-input">
        <input
          type="text"
          placeholder="Say something..."
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
          maxLength={500}
        />
        <button 
          onClick={handleSendMessage}
          disabled={!newMessage.trim()}
        >
          Send
        </button>
      </div>
    </div>
  );
};
```

#### Virtual Gift Panel

```typescript
// Virtual gift sending interface
export const VirtualGiftPanel: React.FC<{
  onSendGift: (giftId: string, quantity: number, message?: string) => void;
  onClose: () => void;
}> = ({ onSendGift, onClose }) => {
  const [gifts, setGifts] = useState<VirtualGift[]>([]);
  const [selectedGift, setSelectedGift] = useState<VirtualGift | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [message, setMessage] = useState('');
  const [userBalance, setUserBalance] = useState(0);
  const [activeCategory, setActiveCategory] = useState('all');

  const giftCategories = [
    { id: 'all', name: 'All', icon: '🎁' },
    { id: 'traditional', name: 'Traditional', icon: '🏺' },
    { id: 'modern', name: 'Modern', icon: '✨' },
    { id: 'premium', name: 'Premium', icon: '💎' },
  ];

  const filteredGifts = gifts.filter(gift => 
    activeCategory === 'all' || gift.category === activeCategory
  );

  const totalCost = selectedGift ? selectedGift.cost_in_coins * quantity : 0;
  const canAfford = totalCost <= userBalance;

  const handleSendGift = () => {
    if (!selectedGift || !canAfford) return;

    onSendGift(selectedGift.id, quantity, message);
    onClose();
  };

  return (
    <div className="virtual-gift-panel">
      <div className="gift-panel-header">
        <h3>Send Gift</h3>
        <div className="user-balance">
          <CoinIcon /> {userBalance.toLocaleString()} coins
        </div>
        <button className="close-btn" onClick={onClose}>
          <CloseIcon />
        </button>
      </div>

      <div className="gift-categories">
        {giftCategories.map(category => (
          <button
            key={category.id}
            className={`category-btn ${activeCategory === category.id ? 'active' : ''}`}
            onClick={() => setActiveCategory(category.id)}
          >
            <span className="category-icon">{category.icon}</span>
            <span className="category-name">{category.name}</span>
          </button>
        ))}
      </div>

      <div className="gifts-grid">
        {filteredGifts.map(gift => (
          <div
            key={gift.id}
            className={`gift-item ${selectedGift?.id === gift.id ? 'selected' : ''}`}
            onClick={() => setSelectedGift(gift)}
          >
            <div className="gift-image">
              <img src={gift.image_url} alt={gift.name} />
              {gift.rarity !== 'common' && (
                <div className={`rarity-badge ${gift.rarity}`}>
                  {gift.rarity}
                </div>
              )}
            </div>
            <div className="gift-info">
              <h4>{gift.name}</h4>
              <div className="gift-cost">
                <CoinIcon /> {gift.cost_in_coins}
              </div>
            </div>
          </div>
        ))}
      </div>

      {selectedGift && (
        <div className="gift-send-section">
          <div className="selected-gift-info">
            <img src={selectedGift.image_url} alt={selectedGift.name} />
            <div className="gift-details">
              <h4>{selectedGift.name}</h4>
              <p>{selectedGift.description}</p>
              <div className="gift-value">
                ≈ ₦{selectedGift.naira_value.toLocaleString()}
              </div>
            </div>
          </div>

          <div className="quantity-selector">
            <label>Quantity:</label>
            <div className="quantity-controls">
              <button 
                onClick={() => setQuantity(Math.max(1, quantity - 1))}
                disabled={quantity <= 1}
              >
                -
              </button>
              <input
                type="number"
                value={quantity}
                onChange={(e) => setQuantity(Math.max(1, parseInt(e.target.value) || 1))}
                min="1"
                max="100"
              />
              <button 
                onClick={() => setQuantity(quantity + 1)}
                disabled={quantity >= 100}
              >
                +
              </button>
            </div>
          </div>

          <div className="gift-message">
            <label>Message (optional):</label>
            <input
              type="text"
              placeholder="Say something nice..."
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              maxLength={200}
            />
          </div>

          <div className="send-summary">
            <div className="total-cost">
              Total: <CoinIcon /> {totalCost.toLocaleString()} coins
            </div>
            {!canAfford && (
              <div className="insufficient-funds">
                Insufficient coins. <button>Buy more</button>
              </div>
            )}
          </div>

          <button
            className="send-gift-btn"
            onClick={handleSendGift}
            disabled={!canAfford}
          >
            Send Gift
          </button>
        </div>
      )}
    </div>
  );
};
```

### Integration Points

#### Payment Integration
- Seamless coin purchase through existing payment gateways
- Wallet integration for coin management
- Automated creator revenue distribution

#### Community Features
- Stream notifications for followers
- Integration with user profiles and achievements
- Community-driven content discovery

#### Content Management
- Stream recording and VOD system
- Automated content moderation
- Analytics and reporting tools

---

*This feature specification provides the complete technical blueprint for implementing the Livestream System within the Great Nigeria Library platform.*
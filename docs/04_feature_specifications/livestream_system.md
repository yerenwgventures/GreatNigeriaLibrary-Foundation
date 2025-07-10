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
#### Stream Management System
Advanced live streaming infrastructure with professional broadcasting capabilities:

- **Stream Sessions**: Complete stream lifecycle management from scheduling to archival with metadata tracking
- **Broadcasting Infrastructure**: Professional RTMP and HLS streaming with adaptive bitrate and quality optimization
- **Viewer Management**: Real-time viewer tracking with concurrent user limits and engagement analytics
- **Stream Categories**: Organized content categorization with custom icons, colors, and discovery optimization
- **Scheduling System**: Advanced stream scheduling with timezone support and automated notifications
- **Recording System**: Automatic stream recording with cloud storage and on-demand playback capabilities
- **Quality Control**: Multi-bitrate streaming with automatic quality adjustment based on viewer connection
- **Stream Analytics**: Comprehensive streaming analytics with viewer engagement and performance metrics

#### Creator Economy Platform
Monetization and creator support ecosystem:

- **Creator Profiles**: Professional creator profiles with branding, social links, and streaming schedules
- **Follower System**: Creator-follower relationships with notification preferences and engagement tracking
- **Monetization Tools**: Multiple revenue streams including premium streams, virtual gifts, and subscription models
- **Virtual Gifts**: Comprehensive virtual gift system with rarity levels, animations, and creator revenue sharing
- **Coin Economy**: Digital currency system for gift purchases with secure transaction processing
- **Creator Analytics**: Detailed creator performance analytics with revenue tracking and audience insights
- **Verification System**: Creator verification program with badge display and enhanced features
- **Revenue Management**: Automated revenue calculation and distribution with transparent reporting

#### Real-Time Interaction System
Advanced real-time communication and engagement features:

- **Live Chat**: Real-time chat system with emoji support, moderation tools, and message filtering
- **Gift Animations**: Interactive gift animations with sound effects and visual celebrations
- **Viewer Engagement**: Real-time viewer interaction tracking with engagement scoring and rewards
- **Moderation Tools**: Comprehensive moderation system with automated filtering and manual review capabilities
- **Community Features**: Follower notifications, stream alerts, and community building tools
- **Social Integration**: Social media sharing and cross-platform promotion capabilities
- **Mobile Optimization**: Mobile-first design with touch-friendly interactions and optimized streaming
- **Accessibility**: Full accessibility compliance with screen reader support and keyboard navigation

#### Content Management and Discovery
Advanced content organization and discovery system:

- **Content Categorization**: Hierarchical content organization with tags, categories, and search optimization
- **Discovery Engine**: AI-powered content discovery with personalized recommendations and trending algorithms
- **Search Functionality**: Advanced search with filters, sorting, and real-time suggestions
- **Content Moderation**: Automated and manual content moderation with community reporting and review systems
- **Archive Management**: Stream archive organization with searchable metadata and highlight extraction
- **Thumbnail Management**: Custom thumbnail support with automatic generation and optimization
- **SEO Optimization**: Search engine optimization for stream discovery and content visibility
- **Analytics Integration**: Content performance analytics with engagement metrics and optimization insights

### API Integration

#### Stream Management API
Comprehensive RESTful API for stream lifecycle management:
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
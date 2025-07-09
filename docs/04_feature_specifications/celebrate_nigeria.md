# Celebrate Nigeria Feature Specification

**Document Version**: 1.0  
**Last Updated**: January 2025  
**Feature Owner**: Cultural Heritage Team  
**Status**: Implemented

---

## Overview

The Celebrate Nigeria feature is a comprehensive digital repository showcasing Nigerian excellence across people, places, events, and cultural elements. It serves as an educational resource and a platform for celebrating Nigeria's rich heritage and achievements.

## Feature Purpose

### Goals
1. **Showcase Excellence**: Highlight accomplished Nigerians, iconic locations, and significant events
2. **Preserve Heritage**: Document and preserve Nigeria's diverse cultural heritage and landmarks
3. **Promote Unity**: Foster national unity by celebrating achievements across all regions
4. **Educate**: Provide accessible, engaging information about Nigerian history and culture
5. **Build Pride**: Cultivate national pride through celebration of collective achievements

### Success Metrics
- **Content Growth**: 1,000+ entries by end of Year 1
- **User Engagement**: 50,000+ monthly active users
- **Educational Impact**: Integration with 100+ educational curricula
- **Community Contribution**: 30% of content from user submissions

## Technical Architecture

### Database Schema

#### Core Tables

```sql
-- Main celebration entries table
CREATE TABLE celebration_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entry_type VARCHAR(20) NOT NULL CHECK (entry_type IN ('person', 'place', 'event')),
    slug VARCHAR(100) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    short_description TEXT NOT NULL,
    full_description TEXT NOT NULL,
    primary_image_url TEXT,
    location VARCHAR(255),
    featured_rank INTEGER DEFAULT 0,
    status VARCHAR(20) DEFAULT 'published' CHECK (status IN ('draft', 'published', 'archived')),
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Categories for organizing entries
CREATE TABLE celebration_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_id UUID REFERENCES celebration_categories(id),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    image_url TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Many-to-many relationship between entries and categories
CREATE TABLE celebration_entry_categories (
    entry_id UUID REFERENCES celebration_entries(id) ON DELETE CASCADE,
    category_id UUID REFERENCES celebration_categories(id) ON DELETE CASCADE,
    PRIMARY KEY (entry_id, category_id)
);

-- Type-specific data for people
CREATE TABLE person_entries (
    celebration_entry_id UUID PRIMARY KEY REFERENCES celebration_entries(id) ON DELETE CASCADE,
    birth_date DATE,
    death_date DATE,
    profession VARCHAR(255),
    achievements TEXT,
    contributions TEXT,
    education TEXT,
    related_links JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Type-specific data for places
CREATE TABLE place_entries (
    celebration_entry_id UUID PRIMARY KEY REFERENCES celebration_entries(id) ON DELETE CASCADE,
    place_type VARCHAR(100),
    latitude DECIMAL(10, 8),
    longitude DECIMAL(11, 8),
    address TEXT,
    visiting_hours TEXT,
    visiting_fees TEXT,
    accessibility TEXT,
    history TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Type-specific data for events
CREATE TABLE event_entries (
    celebration_entry_id UUID PRIMARY KEY REFERENCES celebration_entries(id) ON DELETE CASCADE,
    event_type VARCHAR(100),
    start_date DATE,
    end_date DATE,
    is_recurring BOOLEAN DEFAULT FALSE,
    recurrence_pattern VARCHAR(100),
    organizer VARCHAR(255),
    contact_info TEXT,
    event_history TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Key facts about entries
CREATE TABLE entry_facts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    celebration_entry_id UUID REFERENCES celebration_entries(id) ON DELETE CASCADE,
    label VARCHAR(100) NOT NULL,
    value TEXT NOT NULL,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Media items (images, videos, documents)
CREATE TABLE entry_media (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    celebration_entry_id UUID REFERENCES celebration_entries(id) ON DELETE CASCADE,
    media_type VARCHAR(20) NOT NULL CHECK (media_type IN ('image', 'video', 'audio', 'document')),
    url TEXT NOT NULL,
    caption TEXT,
    alt_text TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User comments on entries
CREATE TABLE entry_comments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    celebration_entry_id UUID REFERENCES celebration_entries(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    parent_comment_id UUID REFERENCES entry_comments(id),
    content TEXT NOT NULL,
    status VARCHAR(20) DEFAULT 'approved' CHECK (status IN ('pending', 'approved', 'rejected')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User votes (likes/dislikes) on entries
CREATE TABLE entry_votes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    celebration_entry_id UUID REFERENCES celebration_entries(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    vote_type VARCHAR(10) NOT NULL CHECK (vote_type IN ('upvote', 'downvote')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(celebration_entry_id, user_id)
);

-- User submissions for new entries
CREATE TABLE entry_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    entry_type VARCHAR(20) NOT NULL CHECK (entry_type IN ('person', 'place', 'event')),
    target_entry_id UUID REFERENCES celebration_entries(id),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    data JSONB NOT NULL,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    reviewer_id UUID REFERENCES users(id),
    reviewer_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### API Endpoints

#### Public Endpoints

```yaml
# List celebration entries
GET /api/v1/celebrate/entries:
  parameters:
    - page: integer
    - limit: integer
    - type: string (person|place|event)
    - category: string
    - featured: boolean
    - search: string
  responses:
    200:
      description: Paginated list of entries
      schema:
        type: object
        properties:
          data:
            type: array
            items:
              $ref: '#/components/schemas/CelebrationEntry'
          pagination:
            $ref: '#/components/schemas/Pagination'

# Get specific entry
GET /api/v1/celebrate/entries/{entryId}:
  parameters:
    - entryId: string (UUID)
  responses:
    200:
      description: Entry details with related content
    404:
      description: Entry not found

# List categories
GET /api/v1/celebrate/categories:
  responses:
    200:
      description: Hierarchical list of categories

# Search entries
GET /api/v1/celebrate/search:
  parameters:
    - q: string (search query)
    - type: string
    - category: string
    - filters: object
  responses:
    200:
      description: Search results with facets
```

#### Authenticated Endpoints

```yaml
# Vote on entry
POST /api/v1/celebrate/entries/{entryId}/vote:
  authentication: required
  body:
    type: object
    properties:
      vote_type:
        type: string
        enum: [upvote, downvote]
  responses:
    201:
      description: Vote recorded
    400:
      description: Invalid vote type

# Comment on entry
POST /api/v1/celebrate/entries/{entryId}/comments:
  authentication: required
  body:
    type: object
    properties:
      content:
        type: string
        minLength: 1
        maxLength: 1000
      parent_comment_id:
        type: string
        format: uuid
  responses:
    201:
      description: Comment created

# Submit new entry
POST /api/v1/celebrate/submissions:
  authentication: required
  body:
    type: object
    properties:
      entry_type:
        type: string
        enum: [person, place, event]
      title:
        type: string
        minLength: 3
        maxLength: 255
      description:
        type: string
        minLength: 50
      data:
        type: object
        description: Type-specific data
  responses:
    201:
      description: Submission created for review
```

### Frontend Components

#### Core Components

```typescript
// Main celebration page component
interface CelebrationPageProps {
  initialEntries?: CelebrationEntry[];
  categories: Category[];
}

export const CelebrationPage: React.FC<CelebrationPageProps> = ({
  initialEntries,
  categories
}) => {
  const [entries, setEntries] = useState(initialEntries || []);
  const [filters, setFilters] = useState<CelebrationFilters>({});
  const [loading, setLoading] = useState(false);

  // Component implementation
  return (
    <div className="celebration-page">
      <CelebrationHeader />
      <CelebrationFilters 
        categories={categories}
        filters={filters}
        onFiltersChange={setFilters}
      />
      <CelebrationGrid 
        entries={entries}
        loading={loading}
      />
      <CelebrationPagination />
    </div>
  );
};

// Entry card component
interface EntryCardProps {
  entry: CelebrationEntry;
  onClick?: (entry: CelebrationEntry) => void;
}

export const EntryCard: React.FC<EntryCardProps> = ({ entry, onClick }) => {
  return (
    <div 
      className="entry-card"
      onClick={() => onClick?.(entry)}
    >
      <div className="entry-image">
        <img 
          src={entry.primary_image_url} 
          alt={entry.title}
          loading="lazy"
        />
        <div className="entry-type-badge">
          {entry.entry_type}
        </div>
      </div>
      <div className="entry-content">
        <h3 className="entry-title">{entry.title}</h3>
        <p className="entry-description">{entry.short_description}</p>
        <div className="entry-stats">
          <span className="views">{entry.view_count} views</span>
          <span className="likes">{entry.like_count} likes</span>
        </div>
      </div>
    </div>
  );
};

// Entry detail page component
interface EntryDetailPageProps {
  entryId: string;
}

export const EntryDetailPage: React.FC<EntryDetailPageProps> = ({ entryId }) => {
  const [entry, setEntry] = useState<CelebrationEntry | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadEntryDetails();
  }, [entryId]);

  const loadEntryDetails = async () => {
    try {
      const response = await celebrationService.getEntry(entryId);
      setEntry(response.data);
      setComments(response.data.comments || []);
    } catch (error) {
      console.error('Failed to load entry:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <EntryDetailSkeleton />;
  if (!entry) return <EntryNotFound />;

  return (
    <div className="entry-detail-page">
      <EntryHeader entry={entry} />
      <EntryContent entry={entry} />
      <EntryFacts facts={entry.facts} />
      <EntryMedia media={entry.media} />
      <EntryActions entry={entry} />
      <EntryComments 
        comments={comments}
        onCommentAdd={handleCommentAdd}
      />
      <RelatedEntries entries={entry.related_entries} />
    </div>
  );
};
```

#### Redux State Management

```typescript
// Celebration slice
interface CelebrationState {
  entries: CelebrationEntry[];
  currentEntry: CelebrationEntry | null;
  categories: Category[];
  filters: CelebrationFilters;
  pagination: Pagination;
  loading: boolean;
  error: string | null;
}

const initialState: CelebrationState = {
  entries: [],
  currentEntry: null,
  categories: [],
  filters: {},
  pagination: { page: 1, limit: 20, total: 0, pages: 0 },
  loading: false,
  error: null,
};

export const celebrationSlice = createSlice({
  name: 'celebration',
  initialState,
  reducers: {
    setLoading: (state, action) => {
      state.loading = action.payload;
    },
    setEntries: (state, action) => {
      state.entries = action.payload;
      state.loading = false;
    },
    setCurrentEntry: (state, action) => {
      state.currentEntry = action.payload;
      state.loading = false;
    },
    setFilters: (state, action) => {
      state.filters = { ...state.filters, ...action.payload };
    },
    updateEntryVotes: (state, action) => {
      const { entryId, voteCount } = action.payload;
      if (state.currentEntry?.id === entryId) {
        state.currentEntry.like_count = voteCount;
      }
      const entryIndex = state.entries.findIndex(e => e.id === entryId);
      if (entryIndex !== -1) {
        state.entries[entryIndex].like_count = voteCount;
      }
    },
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchEntries.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchEntries.fulfilled, (state, action) => {
        state.entries = action.payload.data;
        state.pagination = action.payload.pagination;
        state.loading = false;
      })
      .addCase(fetchEntries.rejected, (state, action) => {
        state.error = action.error.message || 'Failed to fetch entries';
        state.loading = false;
      });
  },
});

// Async thunks
export const fetchEntries = createAsyncThunk(
  'celebration/fetchEntries',
  async (params: GetEntriesParams) => {
    const response = await celebrationService.getEntries(params);
    return response;
  }
);

export const voteForEntry = createAsyncThunk(
  'celebration/voteForEntry',
  async ({ entryId, voteType }: { entryId: string; voteType: string }) => {
    await celebrationService.voteForEntry(entryId, voteType);
    const updatedEntry = await celebrationService.getEntry(entryId);
    return { entryId, voteCount: updatedEntry.data.like_count };
  }
);
```

### Content Guidelines

#### Entry Creation Standards

**People Entries**:
- Must be verifiable public figures with documented achievements
- Minimum 200 words for full description
- At least 3 key facts with sources
- Professional headshot or historical photo required
- Birth/death dates must be accurate and sourced

**Places Entries**:
- Geographical locations within Nigeria or Nigerian diaspora
- Historical, cultural, or educational significance required
- GPS coordinates for physical locations
- High-quality landscape or architectural photos
- Visiting information and accessibility details

**Events Entries**:
- Cultural, historical, or educational significance
- Accurate dates and recurring pattern information
- Official organizer information when applicable
- Multiple photos showing event highlights
- Historical context and cultural importance

#### Content Moderation

**Approval Process**:
1. User submits entry with required information
2. Automated content filtering for inappropriate material
3. Editorial review for accuracy and quality
4. Community feedback period (72 hours)
5. Final approval and publication
6. Ongoing monitoring for user reports

**Quality Standards**:
- All factual claims must be verifiable
- Multiple reliable sources required
- Cultural sensitivity review
- Language appropriateness check
- Image quality and copyright compliance

### User Interaction Features

#### Voting System
- Users can upvote or downvote entries
- Vote counts influence featured entry selection
- Trending algorithm considers recency and votes
- Vote history tracked for abuse prevention

#### Comment System
- Threaded conversations on each entry
- Moderation queue for new comments
- User reputation system affects comment visibility
- Report system for inappropriate content

#### Submission System
- Guided forms for each entry type
- Rich text editor for descriptions
- Image upload with automatic optimization
- Draft saving and collaborative editing
- Submission tracking and status updates

### Integration Points

#### Educational Platform Integration
- Entry suggestions based on current reading material
- Curriculum mapping to relevant cultural content
- Quiz questions generated from entry facts
- Citation system for academic use

#### Search Integration
- Full-text search across all entry content
- Faceted search by type, category, location
- Auto-complete for entry titles and people names
- Related content suggestions

#### Social Features
- Share entries on social media platforms
- Create personal collections of favorite entries
- Follow other users' submission activity
- Achievement badges for content contributions

### Performance Considerations

#### Caching Strategy
- Entry details cached for 1 hour
- Entry lists cached for 15 minutes
- Images served via CDN with aggressive caching
- Search results cached for 5 minutes

#### Database Optimization
- Indexes on frequently queried fields
- Full-text search indexes for content
- Materialized views for complex aggregations
- Connection pooling for high concurrency

#### Mobile Optimization
- Progressive Web App capabilities
- Offline reading for cached entries
- Image lazy loading and compression
- Touch-optimized interface elements

---

*This feature specification serves as the definitive guide for the Celebrate Nigeria platform and should be referenced for all development and content creation activities.* 
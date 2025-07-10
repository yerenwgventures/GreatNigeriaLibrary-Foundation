// User Types - Updated to match backend models
export interface User {
  id: number;                    // Changed from string to number to match backend uint
  username: string;              // Added username field
  email: string;
  full_name: string;             // Changed from name to full_name
  display_name?: string;         // Added display_name
  bio?: string;
  profile_image?: string;        // Changed from avatar to profile_image
  profile_image_url?: string;    // Added for compatibility
  is_active: boolean;            // Added is_active
  is_verified: boolean;          // Added is_verified
  membership_level: string;      // Added membership_level
  points_balance: number;        // Changed from points to points_balance
  reputation: number;            // Added reputation
  role: number;                  // Added role
  last_login: string;            // Added last_login
  last_seen_at: string;          // Added last_seen_at
  created_at: string;            // Changed from createdAt to created_at
  updated_at: string;            // Added updated_at
}

export interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface LoginCredentials {
  email: string;
  password: string;
}

export interface RegisterData {
  username: string;    // Added username field (required by backend)
  full_name: string;   // Changed from name to full_name
  email: string;
  password: string;
}

export interface AuthResponse {
  user: User;
  token: string;
}

// Book Types - Updated to match backend models
export interface Book {
  id: number;          // Changed from string to number to match backend uint
  title: string;
  description: string;
  cover_image: string;
  author: string;
  published_date: string;
  chapters_count?: number;
  total_pages?: number;
}

export interface Chapter {
  id: string;
  title: string;
  order: number;
  book_id?: string;
  sections: Section[];
}

export interface Section {
  id: string;
  title: string;
  order: number;
  chapter_id?: string;
  content?: string;
  next_section_id?: string | null;
  prev_section_id?: string | null;
  has_subsections?: boolean;
}

export interface Subsection {
  id: string;
  title: string;
  number: number;
  section_id: string;
  content?: string;
  format?: string;
  time_to_read?: number;
  published?: boolean;
  next_subsection_id?: string | null;
  prev_subsection_id?: string | null;
}

export interface ReadingProgress {
  book_id: string;
  section_id: string;
  percentage: number;
  last_read_at: string;
}

export interface Bookmark {
  id: string;
  book_id: string;
  book_title: string;
  section_id: string;
  section_title: string;
  created_at: string;
  note?: string;
}

// Interactive Book Elements Types
export interface AudioBookResponse {
  audioUrl: string;
  duration: number;
  title: string;
  generatedAt: string;
}

export interface PhotoBookResponse {
  photoUrls: string[];
  count: number;
  title: string;
  generatedAt: string;
}

export interface VideoBookResponse {
  videoUrl: string;
  duration: number;
  title: string;
  thumbnailUrl: string;
  generatedAt: string;
}

export interface PDFBookResponse {
  pdfUrl: string;
  pageCount: number;
  title: string;
  generatedAt: string;
}

export interface ShareableLinkResponse {
  shareableLink: string;
  mediaUrl: string;
  expiresAt?: string;
}

export interface ForumTopicItem {
  id: string;
  title: string;
  description: string;
  responseCount: number;
  lastResponseAt?: string;
}

export interface ActionStepItem {
  id: string;
  title: string;
  description: string;
  completed: boolean;
  points: number;
}

export interface QuizQuestion {
  id: string;
  question: string;
  options: string[];
  correctOptionIndex?: number;
  explanation?: string;
}

// Profile Types
export interface ReadingStats {
  books_started: number;
  books_completed: number;
  total_pages_read: number;
  reading_streak: number;
  average_reading_time: number;
  reading_history: ReadingHistoryItem[];
}

export interface ReadingHistoryItem {
  date: string;
  pages_read: number;
}

export interface UserActivity {
  id: string;
  type: 'reading' | 'forum' | 'bookmark' | 'note';
  description: string;
  created_at: string;
  book_id?: string;
  book_title?: string;
  topic_id?: string;
  topic_title?: string;
}

// Forum Types
export interface ForumCategory {
  id: string;
  name: string;
  description: string;
  topics_count: number;
}

export interface ForumTopic {
  id: string;
  title: string;
  content?: string;
  author: {
    id: string;
    name: string;
    avatar?: string;
  };
  created_at: string;
  category?: {
    id: string;
    name: string;
  };
  replies_count?: number;
  last_reply_at?: string;
  replies?: ForumReply[];
}

export interface ForumReply {
  id: string;
  content: string;
  author: {
    id: string;
    name: string;
    avatar?: string;
  };
  created_at: string;
  votes: number;
}

export interface NewTopic {
  title: string;
  content: string;
  category_id: string;
}

export interface NewReply {
  content: string;
}

// Resource Types
export interface ResourceCategory {
  id: string;
  name: string;
  description: string;
  resources_count: number;
}

export interface Resource {
  id: string;
  title: string;
  description: string;
  content?: string;
  file_type: string;
  file_size: number;
  download_url: string;
  created_at: string;
  downloads_count: number;
  category?: {
    id: string;
    name: string;
  };
  related_resources?: {
    id: string;
    title: string;
  }[];
}

// Celebrate Nigeria Types
export interface CelebrationEntry {
  id: string;
  type: 'person' | 'place' | 'event';
  name: string;
  slug: string;
  image_url: string;
  summary: string;
  featured?: boolean;
  description?: string;
  birth_date?: string;
  death_date?: string;
  achievements?: string[];
  facts?: {
    title: string;
    content: string;
  }[];
  media?: {
    type: 'image' | 'video';
    url: string;
    caption: string;
  }[];
  related_entries?: {
    id: string;
    type: 'person' | 'place' | 'event';
    name: string;
    slug: string;
  }[];
}

export interface CelebrationCategory {
  id: string;
  name: string;
  entries_count: number;
}

export interface NewCelebrationEntry {
  type: 'person' | 'place' | 'event';
  name: string;
  summary: string;
  description: string;
  birth_date?: string;
  death_date?: string;
  achievements?: string[];
  facts?: {
    title: string;
    content: string;
  }[];
  category_ids: string[];
}

// API Error Types
export interface ApiError {
  code: string;
  message: string;
  details?: string;
}

// Redux Store Types
export interface RootState {
  auth: AuthState;
  books: BooksState;
  forum: ForumState;
  resources: ResourcesState;
  celebrate: CelebrateState;
  profile: ProfileState;
}

export interface BooksState {
  books: Book[];
  currentBook: Book | null;
  chapters: Chapter[];
  currentChapter: Chapter | null;
  currentSection: Section | null;
  currentSubsections: Subsection[];
  currentSubsection: Subsection | null;
  readingProgress: ReadingProgress | null;
  bookmarks: Bookmark[];
  forumTopics: ForumTopicItem[];
  actionSteps: ActionStepItem[];
  quizQuestions: QuizQuestion[];
  audioBook: AudioBookResponse | null;
  photoBook: PhotoBookResponse | null;
  videoBook: VideoBookResponse | null;
  pdfBook: PDFBookResponse | null;
  isLoadingAudio: boolean;
  isLoadingPhoto: boolean;
  isLoadingVideo: boolean;
  isLoadingPdf: boolean;
  isLoading: boolean;
  error: string | null;
}

export interface ForumState {
  categories: ForumCategory[];
  topics: ForumTopic[];
  currentTopic: ForumTopic | null;
  isLoading: boolean;
  error: string | null;
}

export interface ResourcesState {
  categories: ResourceCategory[];
  resources: Resource[];
  currentResource: Resource | null;
  isLoading: boolean;
  error: string | null;
}

export interface CelebrateState {
  featuredEntries: CelebrationEntry[];
  currentEntry: CelebrationEntry | null;
  searchResults: CelebrationEntry[];
  categories: CelebrationCategory[];
  isLoading: boolean;
  error: string | null;
}

export interface ProfileState {
  profile: User | null;
  readingStats: ReadingStats | null;
  activities: UserActivity[];
  bookmarks: Bookmark[];
  isLoading: boolean;
  error: string | null;
}

# Learning Platform Feature Specification

**Document Version**: 1.0  
**Last Updated**: January 2025  
**Feature Owner**: Educational Technology Team  
**Status**: Implemented

---

## Overview

The Learning Platform is the core educational experience of the Great Nigeria Library, providing interactive reading, progress tracking, assessments, and personalized learning paths. It transforms traditional digital reading into an engaging, social, and measurable learning experience.

## Feature Purpose

### Goals
1. **Enhanced Learning**: Transform passive reading into active, interactive learning
2. **Progress Tracking**: Provide detailed analytics on learning progress and outcomes
3. **Personalization**: Adapt content and pacing to individual learning styles
4. **Social Learning**: Enable collaborative learning through discussions and study groups
5. **Assessment**: Validate learning through quizzes, assignments, and peer review

### Success Metrics
- **Completion Rate**: 70%+ of started books completed
- **Learning Retention**: 80%+ quiz success rate
- **Engagement**: 30+ minutes average session duration
- **Social Interaction**: 50%+ of learners participate in discussions
- **Achievement**: 60%+ earn completion certificates

## Technical Architecture

### Database Schema

```sql
-- Books and educational content
CREATE TABLE books (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(500) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    subtitle VARCHAR(500),
    description TEXT,
    cover_image_url TEXT,
    author VARCHAR(255),
    publisher VARCHAR(255),
    isbn VARCHAR(20),
    publication_date DATE,
    language VARCHAR(10) DEFAULT 'en',
    category_id UUID REFERENCES categories(id),
    difficulty_level VARCHAR(20) CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
    estimated_reading_time INTEGER, -- in minutes
    page_count INTEGER,
    word_count INTEGER,
    content_type VARCHAR(20) DEFAULT 'text' CHECK (content_type IN ('text', 'pdf', 'epub', 'audio', 'video')),
    content_url TEXT,
    content_metadata JSONB,
    tags TEXT[],
    featured BOOLEAN DEFAULT FALSE,
    published BOOLEAN DEFAULT FALSE,
    view_count INTEGER DEFAULT 0,
    download_count INTEGER DEFAULT 0,
    rating_average DECIMAL(3,2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Book chapters/sections
CREATE TABLE book_chapters (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    chapter_number INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(200) NOT NULL,
    content TEXT,
    content_html TEXT,
    word_count INTEGER DEFAULT 0,
    estimated_reading_time INTEGER, -- in minutes
    learning_objectives TEXT[],
    key_concepts TEXT[],
    sort_order INTEGER DEFAULT 0,
    published BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(book_id, slug)
);

-- User reading progress
CREATE TABLE reading_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    chapter_id UUID REFERENCES book_chapters(id) ON DELETE CASCADE,
    progress_percentage DECIMAL(5,2) DEFAULT 0,
    reading_position INTEGER DEFAULT 0, -- character position or page number
    time_spent INTEGER DEFAULT 0, -- in seconds
    last_read_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed BOOLEAN DEFAULT FALSE,
    completed_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    bookmarks JSONB DEFAULT '[]',
    highlights JSONB DEFAULT '[]',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, book_id, chapter_id)
);

-- Learning paths and curricula
CREATE TABLE learning_paths (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(200) UNIQUE NOT NULL,
    description TEXT,
    image_url TEXT,
    difficulty_level VARCHAR(20) CHECK (difficulty_level IN ('beginner', 'intermediate', 'advanced')),
    estimated_duration INTEGER, -- in hours
    category_id UUID REFERENCES categories(id),
    instructor_id UUID REFERENCES users(id),
    published BOOLEAN DEFAULT FALSE,
    featured BOOLEAN DEFAULT FALSE,
    enrollment_count INTEGER DEFAULT 0,
    completion_count INTEGER DEFAULT 0,
    rating_average DECIMAL(3,2) DEFAULT 0,
    rating_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Learning path items (books, quizzes, assignments)
CREATE TABLE learning_path_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    learning_path_id UUID REFERENCES learning_paths(id) ON DELETE CASCADE,
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('book', 'quiz', 'assignment', 'discussion', 'video')),
    item_id UUID NOT NULL, -- references books, quizzes, etc.
    title VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order INTEGER NOT NULL,
    required BOOLEAN DEFAULT TRUE,
    unlock_criteria JSONB, -- conditions to unlock this item
    estimated_time INTEGER, -- in minutes
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User enrollment in learning paths
CREATE TABLE learning_path_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    learning_path_id UUID REFERENCES learning_paths(id) ON DELETE CASCADE,
    enrolled_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    progress_percentage DECIMAL(5,2) DEFAULT 0,
    current_item_id UUID REFERENCES learning_path_items(id),
    status VARCHAR(20) DEFAULT 'enrolled' CHECK (status IN ('enrolled', 'active', 'completed', 'paused', 'dropped')),
    UNIQUE(user_id, learning_path_id)
);

-- Quizzes and assessments
CREATE TABLE quizzes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    book_id UUID REFERENCES books(id),
    chapter_id UUID REFERENCES book_chapters(id),
    quiz_type VARCHAR(20) DEFAULT 'knowledge_check' CHECK (quiz_type IN ('knowledge_check', 'comprehension', 'final_exam')),
    time_limit INTEGER, -- in minutes, NULL for unlimited
    attempts_allowed INTEGER DEFAULT 3,
    passing_score DECIMAL(5,2) DEFAULT 70.00,
    randomize_questions BOOLEAN DEFAULT TRUE,
    show_correct_answers BOOLEAN DEFAULT TRUE,
    published BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Quiz questions
CREATE TABLE quiz_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    quiz_id UUID REFERENCES quizzes(id) ON DELETE CASCADE,
    question_type VARCHAR(20) NOT NULL CHECK (question_type IN ('multiple_choice', 'true_false', 'short_answer', 'essay')),
    question_text TEXT NOT NULL,
    explanation TEXT,
    points INTEGER DEFAULT 1,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Quiz question options (for multiple choice)
CREATE TABLE quiz_question_options (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id UUID REFERENCES quiz_questions(id) ON DELETE CASCADE,
    option_text TEXT NOT NULL,
    is_correct BOOLEAN DEFAULT FALSE,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User quiz attempts
CREATE TABLE quiz_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    quiz_id UUID REFERENCES quizzes(id) ON DELETE CASCADE,
    attempt_number INTEGER NOT NULL,
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    score DECIMAL(5,2),
    max_score DECIMAL(5,2),
    percentage DECIMAL(5,2),
    passed BOOLEAN,
    time_taken INTEGER, -- in seconds
    answers JSONB, -- user's answers
    status VARCHAR(20) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'timed_out')),
    UNIQUE(user_id, quiz_id, attempt_number)
);

-- Study sessions
CREATE TABLE study_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_type VARCHAR(20) NOT NULL CHECK (session_type IN ('reading', 'quiz', 'discussion', 'video')),
    item_id UUID NOT NULL, -- book_id, quiz_id, etc.
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ended_at TIMESTAMP WITH TIME ZONE,
    duration INTEGER, -- in seconds
    progress_before DECIMAL(5,2),
    progress_after DECIMAL(5,2),
    actions_taken JSONB, -- highlights, notes, bookmarks, etc.
    device_info JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User achievements and badges
CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    badge_icon_url TEXT,
    category VARCHAR(50),
    criteria JSONB NOT NULL, -- conditions to earn this achievement
    points_reward INTEGER DEFAULT 0,
    rarity VARCHAR(20) DEFAULT 'common' CHECK (rarity IN ('common', 'uncommon', 'rare', 'epic', 'legendary')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- User earned achievements
CREATE TABLE user_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID REFERENCES achievements(id) ON DELETE CASCADE,
    earned_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    progress_data JSONB, -- data that led to earning the achievement
    UNIQUE(user_id, achievement_id)
);

-- Book ratings and reviews
CREATE TABLE book_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    book_id UUID REFERENCES books(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    review_text TEXT,
    helpful_votes INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, book_id)
);
```

### API Endpoints

#### Book and Content Endpoints

```yaml
# List books
GET /api/v1/books:
  parameters:
    - page: integer
    - limit: integer
    - category: string
    - difficulty: string
    - featured: boolean
    - search: string
    - sort: string (title|author|date|rating|popularity)
  responses:
    200:
      description: Paginated list of books
      schema:
        type: object
        properties:
          data:
            type: array
            items:
              $ref: '#/components/schemas/Book'
          pagination:
            $ref: '#/components/schemas/Pagination'

# Get book details
GET /api/v1/books/{bookId}:
  parameters:
    - bookId: string (UUID)
  responses:
    200:
      description: Book details with chapters
    404:
      description: Book not found

# Get book chapter
GET /api/v1/books/{bookId}/chapters/{chapterId}:
  authentication: required
  parameters:
    - bookId: string (UUID)
    - chapterId: string (UUID)
  responses:
    200:
      description: Chapter content with user progress
    403:
      description: Access denied (subscription required)

# Get book content (for reading)
GET /api/v1/books/{bookId}/content:
  authentication: required
  parameters:
    - bookId: string (UUID)
    - chapter: integer (optional)
  responses:
    200:
      description: Book content with navigation
    403:
      description: Access denied

# Rate book
POST /api/v1/books/{bookId}/rating:
  authentication: required
  body:
    type: object
    required: [rating]
    properties:
      rating:
        type: integer
        minimum: 1
        maximum: 5
      review:
        type: string
        maxLength: 2000
  responses:
    201:
      description: Rating submitted
    400:
      description: Invalid rating
```

#### Reading Progress Endpoints

```yaml
# Get reading progress
GET /api/v1/reading/progress:
  authentication: required
  parameters:
    - book_id: string (UUID, optional)
    - status: string (reading|completed|paused)
  responses:
    200:
      description: User's reading progress

# Update reading progress
PUT /api/v1/reading/progress:
  authentication: required
  body:
    type: object
    required: [book_id, chapter_id, progress_percentage]
    properties:
      book_id:
        type: string
        format: uuid
      chapter_id:
        type: string
        format: uuid
      progress_percentage:
        type: number
        minimum: 0
        maximum: 100
      reading_position:
        type: integer
      time_spent:
        type: integer
        description: Time spent in seconds
      notes:
        type: string
      bookmarks:
        type: array
        items:
          type: object
      highlights:
        type: array
        items:
          type: object
  responses:
    200:
      description: Progress updated

# Mark chapter as completed
POST /api/v1/reading/progress/{progressId}/complete:
  authentication: required
  responses:
    200:
      description: Chapter marked as completed
```

#### Learning Path Endpoints

```yaml
# List learning paths
GET /api/v1/learning-paths:
  parameters:
    - page: integer
    - limit: integer
    - category: string
    - difficulty: string
    - featured: boolean
  responses:
    200:
      description: Paginated list of learning paths

# Get learning path details
GET /api/v1/learning-paths/{pathId}:
  parameters:
    - pathId: string (UUID)
  responses:
    200:
      description: Learning path with items and progress

# Enroll in learning path
POST /api/v1/learning-paths/{pathId}/enroll:
  authentication: required
  responses:
    201:
      description: Successfully enrolled
    409:
      description: Already enrolled

# Get learning path progress
GET /api/v1/learning-paths/{pathId}/progress:
  authentication: required
  responses:
    200:
      description: Detailed progress information

# Start/resume learning path
POST /api/v1/learning-paths/{pathId}/start:
  authentication: required
  responses:
    200:
      description: Path started, returns next item
```

#### Quiz and Assessment Endpoints

```yaml
# Get quiz details
GET /api/v1/quizzes/{quizId}:
  authentication: required
  responses:
    200:
      description: Quiz information (without answers)

# Start quiz attempt
POST /api/v1/quizzes/{quizId}/attempts:
  authentication: required
  responses:
    201:
      description: Quiz attempt started
      schema:
        type: object
        properties:
          attempt_id:
            type: string
            format: uuid
          questions:
            type: array
            items:
              $ref: '#/components/schemas/QuizQuestion'

# Submit quiz answer
PUT /api/v1/quizzes/{quizId}/attempts/{attemptId}/answers:
  authentication: required
  body:
    type: object
    required: [question_id, answer]
    properties:
      question_id:
        type: string
        format: uuid
      answer:
        oneOf:
          - type: string
          - type: array
            items:
              type: string
  responses:
    200:
      description: Answer saved

# Submit quiz attempt
POST /api/v1/quizzes/{quizId}/attempts/{attemptId}/submit:
  authentication: required
  responses:
    200:
      description: Quiz submitted and graded
      schema:
        type: object
        properties:
          score:
            type: number
          percentage:
            type: number
          passed:
            type: boolean
          results:
            type: array
            items:
              $ref: '#/components/schemas/QuizResult'
```

### Frontend Components

#### Reading Interface

```typescript
// Main reading component
interface BookReaderProps {
  bookId: string;
  chapterId?: string;
}

export const BookReader: React.FC<BookReaderProps> = ({ bookId, chapterId }) => {
  const [book, setBook] = useState<Book | null>(null);
  const [currentChapter, setCurrentChapter] = useState<Chapter | null>(null);
  const [progress, setProgress] = useState<ReadingProgress | null>(null);
  const [readingSettings, setReadingSettings] = useState<ReadingSettings>({
    fontSize: 16,
    lineHeight: 1.6,
    theme: 'light',
    fontFamily: 'serif',
  });

  const [isHighlighting, setIsHighlighting] = useState(false);
  const [selectedText, setSelectedText] = useState('');
  const [showNotes, setShowNotes] = useState(false);

  useEffect(() => {
    loadBookContent();
    startReadingSession();
    return () => endReadingSession();
  }, [bookId, chapterId]);

  const loadBookContent = async () => {
    try {
      const [bookResponse, progressResponse] = await Promise.all([
        contentService.getBookDetails(bookId),
        progressService.getReadingProgress(bookId)
      ]);

      setBook(bookResponse.data);
      setProgress(progressResponse.data);

      // Load specific chapter or current chapter
      const targetChapter = chapterId 
        ? bookResponse.data.chapters.find(c => c.id === chapterId)
        : progressResponse.data.current_chapter || bookResponse.data.chapters[0];
      
      setCurrentChapter(targetChapter);
    } catch (error) {
      console.error('Failed to load book content:', error);
    }
  };

  const handleTextSelection = () => {
    const selection = window.getSelection();
    if (selection && selection.toString().trim()) {
      setSelectedText(selection.toString());
      setIsHighlighting(true);
    }
  };

  const addHighlight = async (color: string) => {
    if (!selectedText || !currentChapter) return;

    const range = window.getSelection()?.getRangeAt(0);
    if (!range) return;

    const highlight = {
      id: generateId(),
      text: selectedText,
      color,
      startOffset: range.startOffset,
      endOffset: range.endOffset,
      startContainer: range.startContainer,
      endContainer: range.endContainer,
      created_at: new Date().toISOString(),
    };

    try {
      await progressService.addHighlight(bookId, currentChapter.id, highlight);
      setProgress(prev => prev ? {
        ...prev,
        highlights: [...prev.highlights, highlight]
      } : null);
      setIsHighlighting(false);
      setSelectedText('');
    } catch (error) {
      console.error('Failed to add highlight:', error);
    }
  };

  const updateReadingProgress = async (newProgress: number) => {
    if (!currentChapter) return;

    try {
      await progressService.updateProgress({
        book_id: bookId,
        chapter_id: currentChapter.id,
        progress_percentage: newProgress,
        time_spent: getSessionTime(),
      });

      setProgress(prev => prev ? {
        ...prev,
        progress_percentage: newProgress
      } : null);
    } catch (error) {
      console.error('Failed to update progress:', error);
    }
  };

  return (
    <div className={`book-reader theme-${readingSettings.theme}`}>
      <ReaderHeader 
        book={book}
        currentChapter={currentChapter}
        progress={progress}
        onSettingsChange={setReadingSettings}
      />
      
      <div className="reader-content">
        <ReaderSidebar 
          book={book}
          currentChapter={currentChapter}
          progress={progress}
          onChapterSelect={setCurrentChapter}
          showNotes={showNotes}
          onToggleNotes={setShowNotes}
        />

        <div 
          className="reader-main"
          style={{
            fontSize: `${readingSettings.fontSize}px`,
            lineHeight: readingSettings.lineHeight,
            fontFamily: readingSettings.fontFamily,
          }}
          onMouseUp={handleTextSelection}
        >
          {currentChapter && (
            <ChapterContent 
              chapter={currentChapter}
              highlights={progress?.highlights || []}
              bookmarks={progress?.bookmarks || []}
              onProgressUpdate={updateReadingProgress}
            />
          )}
        </div>

        {showNotes && (
          <NotesPanel 
            bookId={bookId}
            chapterId={currentChapter?.id}
            notes={progress?.notes}
            onNotesUpdate={(notes) => setProgress(prev => prev ? {...prev, notes} : null)}
          />
        )}
      </div>

      {isHighlighting && (
        <HighlightToolbar 
          selectedText={selectedText}
          onHighlight={addHighlight}
          onCancel={() => setIsHighlighting(false)}
        />
      )}

      <ReaderFooter 
        book={book}
        currentChapter={currentChapter}
        progress={progress}
        onNavigate={setCurrentChapter}
      />
    </div>
  );
};

// Progress tracking component
interface ProgressTrackerProps {
  userId: string;
  books?: Book[];
}

export const ProgressTracker: React.FC<ProgressTrackerProps> = ({ userId, books }) => {
  const [progressData, setProgressData] = useState<ProgressSummary | null>(null);
  const [timeframe, setTimeframe] = useState<'week' | 'month' | 'year'>('month');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadProgressData();
  }, [userId, timeframe]);

  const loadProgressData = async () => {
    try {
      const response = await progressService.getProgressSummary(userId, timeframe);
      setProgressData(response.data);
    } catch (error) {
      console.error('Failed to load progress data:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <ProgressSkeleton />;

  return (
    <div className="progress-tracker">
      <div className="progress-header">
        <h2>Your Learning Progress</h2>
        <TimeframeSelector 
          value={timeframe}
          onChange={setTimeframe}
        />
      </div>

      <div className="progress-stats">
        <StatCard 
          title="Books Read"
          value={progressData?.books_completed || 0}
          icon="📚"
          trend={progressData?.books_trend}
        />
        <StatCard 
          title="Hours Read"
          value={Math.round((progressData?.time_spent || 0) / 3600)}
          icon="⏱️"
          trend={progressData?.time_trend}
        />
        <StatCard 
          title="Quiz Score Avg"
          value={`${progressData?.average_quiz_score || 0}%`}
          icon="🎯"
          trend={progressData?.quiz_trend}
        />
        <StatCard 
          title="Streak"
          value={`${progressData?.reading_streak || 0} days`}
          icon="🔥"
          trend={progressData?.streak_trend}
        />
      </div>

      <div className="progress-charts">
        <ReadingChart data={progressData?.daily_reading || []} />
        <CategoryBreakdown data={progressData?.category_progress || []} />
      </div>

      <div className="current-books">
        <h3>Currently Reading</h3>
        <div className="book-progress-list">
          {progressData?.current_books?.map(book => (
            <BookProgressCard 
              key={book.id}
              book={book}
              progress={book.progress}
              onContinue={() => router.push(`/books/${book.id}/read`)}
            />
          ))}
        </div>
      </div>

      <div className="achievements">
        <h3>Recent Achievements</h3>
        <AchievementsList achievements={progressData?.recent_achievements || []} />
      </div>
    </div>
  );
};
```

#### Quiz Interface

```typescript
// Quiz taking component
interface QuizInterfaceProps {
  quizId: string;
  onComplete?: (result: QuizResult) => void;
}

export const QuizInterface: React.FC<QuizInterfaceProps> = ({ quizId, onComplete }) => {
  const [quiz, setQuiz] = useState<Quiz | null>(null);
  const [currentAttempt, setCurrentAttempt] = useState<QuizAttempt | null>(null);
  const [currentQuestionIndex, setCurrentQuestionIndex] = useState(0);
  const [answers, setAnswers] = useState<Record<string, any>>({});
  const [timeRemaining, setTimeRemaining] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadQuiz();
  }, [quizId]);

  useEffect(() => {
    if (quiz?.time_limit && currentAttempt) {
      const timer = setInterval(() => {
        setTimeRemaining(prev => {
          if (prev === null) return null;
          if (prev <= 0) {
            submitQuiz(true); // Auto-submit on timeout
            return 0;
          }
          return prev - 1;
        });
      }, 1000);

      return () => clearInterval(timer);
    }
  }, [quiz, currentAttempt]);

  const loadQuiz = async () => {
    try {
      const response = await quizService.getQuiz(quizId);
      setQuiz(response.data);
      
      // Start new attempt
      const attemptResponse = await quizService.startAttempt(quizId);
      setCurrentAttempt(attemptResponse.data);
      
      if (response.data.time_limit) {
        setTimeRemaining(response.data.time_limit * 60); // Convert minutes to seconds
      }
    } catch (error) {
      console.error('Failed to load quiz:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAnswerChange = async (questionId: string, answer: any) => {
    setAnswers(prev => ({ ...prev, [questionId]: answer }));

    // Auto-save answer
    if (currentAttempt) {
      try {
        await quizService.saveAnswer(quizId, currentAttempt.id, questionId, answer);
      } catch (error) {
        console.error('Failed to save answer:', error);
      }
    }
  };

  const submitQuiz = async (isTimeout = false) => {
    if (!currentAttempt) return;

    try {
      const result = await quizService.submitAttempt(quizId, currentAttempt.id);
      onComplete?.(result.data);
    } catch (error) {
      console.error('Failed to submit quiz:', error);
    }
  };

  const navigateQuestion = (direction: 'prev' | 'next') => {
    if (!quiz) return;

    const newIndex = direction === 'next' 
      ? Math.min(currentQuestionIndex + 1, quiz.questions.length - 1)
      : Math.max(currentQuestionIndex - 1, 0);
    
    setCurrentQuestionIndex(newIndex);
  };

  if (loading) return <QuizSkeleton />;
  if (!quiz || !currentAttempt) return <div>Quiz not found</div>;

  const currentQuestion = quiz.questions[currentQuestionIndex];
  const isLastQuestion = currentQuestionIndex === quiz.questions.length - 1;

  return (
    <div className="quiz-interface">
      <QuizHeader 
        quiz={quiz}
        questionNumber={currentQuestionIndex + 1}
        totalQuestions={quiz.questions.length}
        timeRemaining={timeRemaining}
      />

      <div className="quiz-content">
        <QuizProgress 
          current={currentQuestionIndex + 1}
          total={quiz.questions.length}
          answered={Object.keys(answers).length}
        />

        <QuestionCard 
          question={currentQuestion}
          answer={answers[currentQuestion.id]}
          onAnswerChange={(answer) => handleAnswerChange(currentQuestion.id, answer)}
        />

        <div className="quiz-navigation">
          <button 
            onClick={() => navigateQuestion('prev')}
            disabled={currentQuestionIndex === 0}
            className="nav-button prev"
          >
            Previous
          </button>

          {isLastQuestion ? (
            <button 
              onClick={() => submitQuiz()}
              className="submit-button"
              disabled={Object.keys(answers).length < quiz.questions.length}
            >
              Submit Quiz
            </button>
          ) : (
            <button 
              onClick={() => navigateQuestion('next')}
              className="nav-button next"
            >
              Next
            </button>
          )}
        </div>
      </div>

      <QuestionOverview 
        questions={quiz.questions}
        answers={answers}
        currentQuestion={currentQuestionIndex}
        onQuestionSelect={setCurrentQuestionIndex}
      />
    </div>
  );
};
```

### Learning Analytics

#### Progress Analytics Service

```typescript
// Learning analytics service
export class LearningAnalyticsService {
  static async calculateLearningMetrics(userId: string, timeframe: string) {
    const endDate = new Date();
    const startDate = this.getStartDate(timeframe);

    const metrics = await Promise.all([
      this.getReadingTime(userId, startDate, endDate),
      this.getCompletionRates(userId, startDate, endDate),
      this.getQuizPerformance(userId, startDate, endDate),
      this.getLearningStreak(userId),
      this.getCategoryProgress(userId, startDate, endDate),
    ]);

    return {
      reading_time: metrics[0],
      completion_rates: metrics[1],
      quiz_performance: metrics[2],
      learning_streak: metrics[3],
      category_progress: metrics[4],
      recommendations: await this.generateRecommendations(userId, metrics),
    };
  }

  static async generateRecommendations(userId: string, metrics: any[]) {
    const recommendations = [];

    // Analyze reading patterns
    if (metrics[0].daily_average < 30) { // Less than 30 minutes per day
      recommendations.push({
        type: 'reading_time',
        priority: 'medium',
        message: 'Try reading for at least 30 minutes daily to improve retention',
        action: 'set_reading_goal',
      });
    }

    // Analyze quiz performance
    if (metrics[2].average_score < 70) {
      recommendations.push({
        type: 'quiz_performance',
        priority: 'high',
        message: 'Review quiz explanations to improve understanding',
        action: 'review_quizzes',
      });
    }

    // Suggest learning paths based on interests
    const interests = await this.getUserInterests(userId);
    const suggestedPaths = await this.getRecommendedLearningPaths(interests);
    
    recommendations.push({
      type: 'learning_path',
      priority: 'low',
      message: 'Explore these learning paths based on your interests',
      action: 'browse_paths',
      data: suggestedPaths,
    });

    return recommendations;
  }

  private static getStartDate(timeframe: string): Date {
    const date = new Date();
    switch (timeframe) {
      case 'week':
        date.setDate(date.getDate() - 7);
        break;
      case 'month':
        date.setMonth(date.getMonth() - 1);
        break;
      case 'year':
        date.setFullYear(date.getFullYear() - 1);
        break;
    }
    return date;
  }
}
```

### Gamification Features

#### Achievement System

```typescript
// Achievement engine
export class AchievementEngine {
  private static readonly ACHIEVEMENT_RULES = [
    {
      id: 'first_book',
      name: 'Getting Started',
      description: 'Complete your first book',
      criteria: { books_completed: 1 },
      points: 100,
      badge: 'first_book.svg',
    },
    {
      id: 'speed_reader',
      name: 'Speed Reader',
      description: 'Read 5 books in one month',
      criteria: { books_per_month: 5 },
      points: 500,
      badge: 'speed_reader.svg',
    },
    {
      id: 'quiz_master',
      name: 'Quiz Master',
      description: 'Score 100% on 10 quizzes',
      criteria: { perfect_quizzes: 10 },
      points: 300,
      badge: 'quiz_master.svg',
    },
    // ... more achievements
  ];

  static async checkAchievements(userId: string, activityType: string, activityData: any) {
    const userStats = await this.getUserStats(userId);
    const newAchievements = [];

    for (const rule of this.ACHIEVEMENT_RULES) {
      const hasAchievement = await this.userHasAchievement(userId, rule.id);
      if (!hasAchievement && this.meetsAchievementCriteria(rule, userStats, activityData)) {
        await this.awardAchievement(userId, rule.id);
        newAchievements.push(rule);
      }
    }

    return newAchievements;
  }

  private static meetsAchievementCriteria(rule: any, userStats: any, activityData: any): boolean {
    // Implement achievement criteria checking logic
    return false;
  }

  private static async awardAchievement(userId: string, achievementId: string) {
    // Award achievement and points to user
    await db.user_achievements.create({
      data: {
        user_id: userId,
        achievement_id: achievementId,
        earned_at: new Date(),
      },
    });

    // Add points to user account
    await pointsService.addPoints(userId, 'achievement', achievementId);
  }
}
```

---

*This feature specification provides the complete technical blueprint for implementing the Learning Platform core functionality.*

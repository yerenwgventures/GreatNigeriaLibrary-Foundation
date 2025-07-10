import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { BookService } from '../../api';
import { BooksState } from '../../types';

// Initial state
const initialState: BooksState = {
  books: [],
  currentBook: null,
  chapters: [],
  currentChapter: null,
  currentSection: null,
  currentSubsections: [],
  currentSubsection: null,
  readingProgress: null,
  bookmarks: [],
  forumTopics: [],
  actionSteps: [],
  quizQuestions: [],
  audioBook: null,
  photoBook: null,
  videoBook: null,
  pdfBook: null,
  isLoadingAudio: false,
  isLoadingPhoto: false,
  isLoadingVideo: false,
  isLoadingPdf: false,
  isLoading: false,
  error: null,
};

// Async thunks
export const fetchBooks = createAsyncThunk('books/fetchBooks', async (_, { rejectWithValue }) => {
  try {
    return await BookService.getBooks();
  } catch (error: any) {
    return rejectWithValue(error.message || 'Failed to fetch books');
  }
});

export const fetchBookById = createAsyncThunk(
  'books/fetchBookById',
  async (bookId: string, { rejectWithValue }) => {
    try {
      return await BookService.getBookById(bookId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch book');
    }
  }
);

export const fetchBookChapters = createAsyncThunk(
  'books/fetchBookChapters',
  async (bookId: string, { rejectWithValue }) => {
    try {
      return await BookService.getBookChapters(bookId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch book chapters');
    }
  }
);

export const fetchChapterById = createAsyncThunk(
  'books/fetchChapterById',
  async (chapterId: string, { rejectWithValue }) => {
    try {
      return await BookService.getChapterById(chapterId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch chapter');
    }
  }
);

export const fetchSectionById = createAsyncThunk(
  'books/fetchSectionById',
  async (sectionId: string, { rejectWithValue }) => {
    try {
      return await BookService.getSectionById(sectionId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch section');
    }
  }
);

export const saveReadingProgress = createAsyncThunk(
  'books/saveReadingProgress',
  async ({ bookId, sectionId }: { bookId: string; sectionId: string }, { rejectWithValue }) => {
    try {
      await BookService.saveReadingProgress(bookId, sectionId);
      return { bookId, sectionId };
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to save reading progress');
    }
  }
);

export const fetchReadingProgress = createAsyncThunk(
  'books/fetchReadingProgress',
  async (bookId: string, { rejectWithValue }) => {
    try {
      return await BookService.getReadingProgress(bookId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch reading progress');
    }
  }
);

export const fetchBookmarks = createAsyncThunk(
  'books/fetchBookmarks',
  async (bookId: string, { rejectWithValue }) => {
    try {
      return await BookService.getBookmarks(bookId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch bookmarks');
    }
  }
);

export const addBookmark = createAsyncThunk(
  'books/addBookmark',
  async (
    { bookId, sectionId, note }: { bookId: string; sectionId: string; note?: string },
    { rejectWithValue }
  ) => {
    try {
      return await BookService.addBookmark(bookId, sectionId, note);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to add bookmark');
    }
  }
);

export const deleteBookmark = createAsyncThunk(
  'books/deleteBookmark',
  async (bookmarkId: string, { rejectWithValue }) => {
    try {
      await BookService.deleteBookmark(bookmarkId);
      return bookmarkId;
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to delete bookmark');
    }
  }
);

export const fetchSubsectionsBySection = createAsyncThunk(
  'books/fetchSubsectionsBySection',
  async (sectionId: string, { rejectWithValue }) => {
    try {
      return await BookService.getSubsectionsBySection(sectionId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch subsections');
    }
  }
);

export const fetchSubsectionById = createAsyncThunk(
  'books/fetchSubsectionById',
  async (subsectionId: string, { rejectWithValue }) => {
    try {
      return await BookService.getSubsectionById(subsectionId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch subsection');
    }
  }
);

// Interactive elements thunks
export const fetchForumTopics = createAsyncThunk(
  'books/fetchForumTopics',
  async (sectionId: string, { rejectWithValue }) => {
    try {
      return await BookService.getForumTopics(sectionId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch forum topics');
    }
  }
);

export const fetchActionSteps = createAsyncThunk(
  'books/fetchActionSteps',
  async (sectionId: string, { rejectWithValue }) => {
    try {
      return await BookService.getActionSteps(sectionId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch action steps');
    }
  }
);

export const fetchQuizQuestions = createAsyncThunk(
  'books/fetchQuizQuestions',
  async (sectionId: string, { rejectWithValue }) => {
    try {
      return await BookService.getQuizQuestions(sectionId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch quiz questions');
    }
  }
);

export const generateAudioBook = createAsyncThunk(
  'books/generateAudioBook',
  async (sectionId: string, { rejectWithValue }) => {
    try {
      return await BookService.generateAudio(sectionId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to generate audio book');
    }
  }
);

export const generatePhotoBook = createAsyncThunk(
  'books/generatePhotoBook',
  async (sectionId: string, { rejectWithValue }) => {
    try {
      return await BookService.generatePhotoBook(sectionId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to generate photo book');
    }
  }
);

export const generateVideoBook = createAsyncThunk(
  'books/generateVideoBook',
  async (sectionId: string, { rejectWithValue }) => {
    try {
      return await BookService.generateVideo(sectionId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to generate video book');
    }
  }
);

export const generatePdfBook = createAsyncThunk(
  'books/generatePdfBook',
  async (sectionId: string, { rejectWithValue }) => {
    try {
      return await BookService.generatePdf(sectionId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to generate PDF book');
    }
  }
);

export const getShareableLink = createAsyncThunk(
  'books/getShareableLink',
  async ({ sectionId, mediaType }: { sectionId: string; mediaType: string }, { rejectWithValue }) => {
    try {
      return await BookService.getShareableLink(sectionId, mediaType);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to get shareable link');
    }
  }
);

// Slice
const booksSlice = createSlice({
  name: 'books',
  initialState,
  reducers: {
    clearCurrentBook: (state) => {
      state.currentBook = null;
      state.chapters = [];
      state.currentChapter = null;
      state.currentSection = null;
      state.currentSubsections = [];
      state.currentSubsection = null;
      state.readingProgress = null;
      state.bookmarks = [];
      state.forumTopics = [];
      state.actionSteps = [];
      state.quizQuestions = [];
      state.audioBook = null;
      state.photoBook = null;
      state.videoBook = null;
      state.pdfBook = null;
    },
    clearError: (state) => {
      state.error = null;
    },
    clearInteractiveElements: (state) => {
      state.forumTopics = [];
      state.actionSteps = [];
      state.quizQuestions = [];
      state.audioBook = null;
      state.photoBook = null;
      state.videoBook = null;
      state.pdfBook = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch books
      .addCase(fetchBooks.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchBooks.fulfilled, (state, action) => {
        state.isLoading = false;
        state.books = action.payload;
      })
      .addCase(fetchBooks.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch book by ID
      .addCase(fetchBookById.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchBookById.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentBook = action.payload;
      })
      .addCase(fetchBookById.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch book chapters
      .addCase(fetchBookChapters.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchBookChapters.fulfilled, (state, action) => {
        state.isLoading = false;
        state.chapters = action.payload;
      })
      .addCase(fetchBookChapters.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch chapter by ID
      .addCase(fetchChapterById.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchChapterById.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentChapter = action.payload;
      })
      .addCase(fetchChapterById.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch section by ID
      .addCase(fetchSectionById.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchSectionById.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentSection = action.payload;
      })
      .addCase(fetchSectionById.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch reading progress
      .addCase(fetchReadingProgress.fulfilled, (state, action) => {
        state.readingProgress = action.payload;
      })
      // Fetch bookmarks
      .addCase(fetchBookmarks.fulfilled, (state, action) => {
        state.bookmarks = action.payload;
      })
      // Add bookmark
      .addCase(addBookmark.fulfilled, (state, action) => {
        state.bookmarks.push(action.payload);
      })
      // Delete bookmark
      .addCase(deleteBookmark.fulfilled, (state, action) => {
        state.bookmarks = state.bookmarks.filter((bookmark) => bookmark.id !== action.payload);
      })
      // Fetch subsections by section
      .addCase(fetchSubsectionsBySection.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchSubsectionsBySection.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentSubsections = action.payload;
      })
      .addCase(fetchSubsectionsBySection.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch subsection by ID
      .addCase(fetchSubsectionById.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchSubsectionById.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentSubsection = action.payload;
      })
      .addCase(fetchSubsectionById.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })

      // Forum topics
      .addCase(fetchForumTopics.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchForumTopics.fulfilled, (state, action) => {
        state.isLoading = false;
        state.forumTopics = action.payload;
      })
      .addCase(fetchForumTopics.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })

      // Action steps
      .addCase(fetchActionSteps.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchActionSteps.fulfilled, (state, action) => {
        state.isLoading = false;
        state.actionSteps = action.payload;
      })
      .addCase(fetchActionSteps.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })

      // Quiz questions
      .addCase(fetchQuizQuestions.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchQuizQuestions.fulfilled, (state, action) => {
        state.isLoading = false;
        state.quizQuestions = action.payload;
      })
      .addCase(fetchQuizQuestions.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })

      // Audio book
      .addCase(generateAudioBook.pending, (state) => {
        state.isLoadingAudio = true;
        state.error = null;
      })
      .addCase(generateAudioBook.fulfilled, (state, action) => {
        state.isLoadingAudio = false;
        state.audioBook = action.payload;
      })
      .addCase(generateAudioBook.rejected, (state, action) => {
        state.isLoadingAudio = false;
        state.error = action.payload as string;
      })

      // Photo book
      .addCase(generatePhotoBook.pending, (state) => {
        state.isLoadingPhoto = true;
        state.error = null;
      })
      .addCase(generatePhotoBook.fulfilled, (state, action) => {
        state.isLoadingPhoto = false;
        state.photoBook = action.payload;
      })
      .addCase(generatePhotoBook.rejected, (state, action) => {
        state.isLoadingPhoto = false;
        state.error = action.payload as string;
      })

      // Video book
      .addCase(generateVideoBook.pending, (state) => {
        state.isLoadingVideo = true;
        state.error = null;
      })
      .addCase(generateVideoBook.fulfilled, (state, action) => {
        state.isLoadingVideo = false;
        state.videoBook = action.payload;
      })
      .addCase(generateVideoBook.rejected, (state, action) => {
        state.isLoadingVideo = false;
        state.error = action.payload as string;
      })

      // PDF book
      .addCase(generatePdfBook.pending, (state) => {
        state.isLoadingPdf = true;
        state.error = null;
      })
      .addCase(generatePdfBook.fulfilled, (state, action) => {
        state.isLoadingPdf = false;
        state.pdfBook = action.payload;
      })
      .addCase(generatePdfBook.rejected, (state, action) => {
        state.isLoadingPdf = false;
        state.error = action.payload as string;
      });
  },
});

export const { clearCurrentBook, clearError, clearInteractiveElements } = booksSlice.actions;
export default booksSlice.reducer;

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
    },
    clearError: (state) => {
      state.error = null;
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
      });
  },
});

export const { clearCurrentBook, clearError } = booksSlice.actions;
export default booksSlice.reducer;

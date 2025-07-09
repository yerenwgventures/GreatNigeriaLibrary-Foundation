import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { RootState } from '../../store';
import searchService, { SearchRequest, SearchResponse, SearchResult } from '../../api/searchService';

// Define types
interface SearchState {
  results: SearchResult[];
  totalResults: number;
  hasMore: boolean;
  page: number;
  pageSize: number;
  query: string;
  types: string[];
  sortBy: string;
  sortOrder: 'asc' | 'desc';
  tags: string[];
  dateFrom: string | null;
  dateTo: string | null;
  recentSearches: string[];
  popularSearches: { query: string; count: number }[];
  isLoading: boolean;
  error: string | null;
}

// Initial state
const initialState: SearchState = {
  results: [],
  totalResults: 0,
  hasMore: false,
  page: 1,
  pageSize: 10,
  query: '',
  types: [],
  sortBy: 'relevance',
  sortOrder: 'desc',
  tags: [],
  dateFrom: null,
  dateTo: null,
  recentSearches: [],
  popularSearches: [],
  isLoading: false,
  error: null,
};

// Async thunks
export const searchContent = createAsyncThunk(
  'search/searchContent',
  async (request: SearchRequest, { rejectWithValue }) => {
    try {
      const response = await searchService.search(request);
      
      // Save the search query to recent searches
      if (request.query.trim()) {
        await searchService.saveRecentSearch(request.query);
      }
      
      return response;
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to search content');
    }
  }
);

export const fetchRecentSearches = createAsyncThunk(
  'search/fetchRecentSearches',
  async (_, { rejectWithValue }) => {
    try {
      return await searchService.getRecentSearches();
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch recent searches');
    }
  }
);

export const clearRecentSearches = createAsyncThunk(
  'search/clearRecentSearches',
  async (_, { rejectWithValue }) => {
    try {
      await searchService.clearRecentSearches();
      return true;
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to clear recent searches');
    }
  }
);

export const fetchPopularSearches = createAsyncThunk(
  'search/fetchPopularSearches',
  async (_, { rejectWithValue }) => {
    try {
      return await searchService.getPopularSearches();
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch popular searches');
    }
  }
);

// Create slice
const searchSlice = createSlice({
  name: 'search',
  initialState,
  reducers: {
    setQuery: (state, action) => {
      state.query = action.payload;
    },
    setTypes: (state, action) => {
      state.types = action.payload;
    },
    setPage: (state, action) => {
      state.page = action.payload;
    },
    setPageSize: (state, action) => {
      state.pageSize = action.payload;
    },
    setSortBy: (state, action) => {
      state.sortBy = action.payload;
    },
    setSortOrder: (state, action) => {
      state.sortOrder = action.payload;
    },
    setTags: (state, action) => {
      state.tags = action.payload;
    },
    setDateRange: (state, action) => {
      state.dateFrom = action.payload.from;
      state.dateTo = action.payload.to;
    },
    resetFilters: (state) => {
      state.types = [];
      state.sortBy = 'relevance';
      state.sortOrder = 'desc';
      state.tags = [];
      state.dateFrom = null;
      state.dateTo = null;
    },
    clearSearchResults: (state) => {
      state.results = [];
      state.totalResults = 0;
      state.hasMore = false;
      state.page = 1;
    },
    clearSearchError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // searchContent
      .addCase(searchContent.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(searchContent.fulfilled, (state, action) => {
        const response = action.payload as SearchResponse;
        state.isLoading = false;
        
        // If it's the first page, replace results; otherwise, append
        if (state.page === 1) {
          state.results = response.results;
        } else {
          state.results = [...state.results, ...response.results];
        }
        
        state.totalResults = response.totalResults;
        state.hasMore = response.hasMore;
        state.page = response.page;
        state.pageSize = response.pageSize;
      })
      .addCase(searchContent.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      
      // fetchRecentSearches
      .addCase(fetchRecentSearches.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchRecentSearches.fulfilled, (state, action) => {
        state.isLoading = false;
        state.recentSearches = action.payload;
      })
      .addCase(fetchRecentSearches.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      
      // clearRecentSearches
      .addCase(clearRecentSearches.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(clearRecentSearches.fulfilled, (state) => {
        state.isLoading = false;
        state.recentSearches = [];
      })
      .addCase(clearRecentSearches.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      
      // fetchPopularSearches
      .addCase(fetchPopularSearches.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchPopularSearches.fulfilled, (state, action) => {
        state.isLoading = false;
        state.popularSearches = action.payload;
      })
      .addCase(fetchPopularSearches.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });
  },
});

// Export actions
export const {
  setQuery,
  setTypes,
  setPage,
  setPageSize,
  setSortBy,
  setSortOrder,
  setTags,
  setDateRange,
  resetFilters,
  clearSearchResults,
  clearSearchError,
} = searchSlice.actions;

// Export selectors
export const selectSearchResults = (state: RootState) => state.search.results;
export const selectTotalResults = (state: RootState) => state.search.totalResults;
export const selectHasMore = (state: RootState) => state.search.hasMore;
export const selectPage = (state: RootState) => state.search.page;
export const selectPageSize = (state: RootState) => state.search.pageSize;
export const selectQuery = (state: RootState) => state.search.query;
export const selectTypes = (state: RootState) => state.search.types;
export const selectSortBy = (state: RootState) => state.search.sortBy;
export const selectSortOrder = (state: RootState) => state.search.sortOrder;
export const selectTags = (state: RootState) => state.search.tags;
export const selectDateRange = (state: RootState) => ({
  from: state.search.dateFrom,
  to: state.search.dateTo,
});
export const selectRecentSearches = (state: RootState) => state.search.recentSearches;
export const selectPopularSearches = (state: RootState) => state.search.popularSearches;
export const selectSearchLoading = (state: RootState) => state.search.isLoading;
export const selectSearchError = (state: RootState) => state.search.error;

// Export reducer
export default searchSlice.reducer;

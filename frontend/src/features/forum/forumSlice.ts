import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { ForumService } from '../../api';
import { ForumState, NewReply, NewTopic } from '../../types';

// Initial state
const initialState: ForumState = {
  categories: [],
  topics: [],
  currentTopic: null,
  isLoading: false,
  error: null,
};

// Async thunks
export const fetchCategories = createAsyncThunk(
  'forum/fetchCategories',
  async (_, { rejectWithValue }) => {
    try {
      return await ForumService.getCategories();
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch categories');
    }
  }
);

export const fetchTopicsByCategory = createAsyncThunk(
  'forum/fetchTopicsByCategory',
  async (categoryId: string, { rejectWithValue }) => {
    try {
      return await ForumService.getTopicsByCategory(categoryId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch topics');
    }
  }
);

export const fetchTopicById = createAsyncThunk(
  'forum/fetchTopicById',
  async (topicId: string, { rejectWithValue }) => {
    try {
      return await ForumService.getTopicById(topicId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch topic');
    }
  }
);

export const createTopic = createAsyncThunk(
  'forum/createTopic',
  async (newTopic: NewTopic, { rejectWithValue }) => {
    try {
      return await ForumService.createTopic(newTopic);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to create topic');
    }
  }
);

export const createReply = createAsyncThunk(
  'forum/createReply',
  async (
    { topicId, newReply }: { topicId: string; newReply: NewReply },
    { rejectWithValue }
  ) => {
    try {
      return await ForumService.createReply(topicId, newReply);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to create reply');
    }
  }
);

export const voteReply = createAsyncThunk(
  'forum/voteReply',
  async (
    { replyId, vote }: { replyId: string; vote: 'up' | 'down' },
    { rejectWithValue }
  ) => {
    try {
      const response = await ForumService.voteReply(replyId, vote);
      return { replyId, votes: response.votes };
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to vote on reply');
    }
  }
);

export const searchTopics = createAsyncThunk(
  'forum/searchTopics',
  async (query: string, { rejectWithValue }) => {
    try {
      return await ForumService.searchTopics(query);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to search topics');
    }
  }
);

export const deleteTopic = createAsyncThunk(
  'forum/deleteTopic',
  async (topicId: string, { rejectWithValue }) => {
    try {
      await ForumService.deleteTopic(topicId);
      return topicId;
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to delete topic');
    }
  }
);

export const deleteReply = createAsyncThunk(
  'forum/deleteReply',
  async (
    { topicId, replyId }: { topicId: string; replyId: string },
    { rejectWithValue }
  ) => {
    try {
      await ForumService.deleteReply(replyId);
      return { topicId, replyId };
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to delete reply');
    }
  }
);

// Slice
const forumSlice = createSlice({
  name: 'forum',
  initialState,
  reducers: {
    clearCurrentTopic: (state) => {
      state.currentTopic = null;
    },
    clearError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch categories
      .addCase(fetchCategories.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchCategories.fulfilled, (state, action) => {
        state.isLoading = false;
        state.categories = action.payload;
      })
      .addCase(fetchCategories.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch topics by category
      .addCase(fetchTopicsByCategory.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchTopicsByCategory.fulfilled, (state, action) => {
        state.isLoading = false;
        state.topics = action.payload;
      })
      .addCase(fetchTopicsByCategory.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch topic by ID
      .addCase(fetchTopicById.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchTopicById.fulfilled, (state, action) => {
        state.isLoading = false;
        state.currentTopic = action.payload;
      })
      .addCase(fetchTopicById.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Create topic
      .addCase(createTopic.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(createTopic.fulfilled, (state, action) => {
        state.isLoading = false;
        state.topics.unshift(action.payload);
      })
      .addCase(createTopic.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Create reply
      .addCase(createReply.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(createReply.fulfilled, (state, action) => {
        state.isLoading = false;
        if (state.currentTopic) {
          if (!state.currentTopic.replies) {
            state.currentTopic.replies = [];
          }
          state.currentTopic.replies.push(action.payload);
        }
      })
      .addCase(createReply.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Vote reply
      .addCase(voteReply.fulfilled, (state, action) => {
        if (state.currentTopic && state.currentTopic.replies) {
          const reply = state.currentTopic.replies.find(
            (r) => r.id === action.payload.replyId
          );
          if (reply) {
            reply.votes = action.payload.votes;
          }
        }
      })
      // Search topics
      .addCase(searchTopics.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(searchTopics.fulfilled, (state, action) => {
        state.isLoading = false;
        state.topics = action.payload;
      })
      .addCase(searchTopics.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Delete topic
      .addCase(deleteTopic.fulfilled, (state, action) => {
        state.topics = state.topics.filter((topic) => topic.id !== action.payload);
        if (state.currentTopic && state.currentTopic.id === action.payload) {
          state.currentTopic = null;
        }
      })
      // Delete reply
      .addCase(deleteReply.fulfilled, (state, action) => {
        if (
          state.currentTopic &&
          state.currentTopic.id === action.payload.topicId &&
          state.currentTopic.replies
        ) {
          state.currentTopic.replies = state.currentTopic.replies.filter(
            (reply) => reply.id !== action.payload.replyId
          );
        }
      });
  },
});

export const { clearCurrentTopic, clearError } = forumSlice.actions;
export default forumSlice.reducer;

import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import { UserService } from '../../api';
import { ProfileState, User } from '../../types';

// Initial state
const initialState: ProfileState = {
  profile: null,
  readingStats: null,
  activities: [],
  bookmarks: [],
  isLoading: false,
  error: null,
};

// Async thunks
export const fetchUserProfile = createAsyncThunk(
  'profile/fetchUserProfile',
  async (userId: string, { rejectWithValue }) => {
    try {
      return await UserService.getUserProfile(userId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch user profile');
    }
  }
);

export const updateUserProfile = createAsyncThunk(
  'profile/updateUserProfile',
  async ({ userId, profileData }: { userId: string; profileData: Partial<User> }, { rejectWithValue }) => {
    try {
      return await UserService.updateUserProfile(userId, profileData);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to update user profile');
    }
  }
);

export const fetchReadingStats = createAsyncThunk(
  'profile/fetchReadingStats',
  async (userId: string, { rejectWithValue }) => {
    try {
      return await UserService.getReadingStats(userId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch reading stats');
    }
  }
);

export const fetchUserBookmarks = createAsyncThunk(
  'profile/fetchUserBookmarks',
  async (userId: string, { rejectWithValue }) => {
    try {
      return await UserService.getUserBookmarks(userId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch user bookmarks');
    }
  }
);

export const fetchUserActivities = createAsyncThunk(
  'profile/fetchUserActivities',
  async (userId: string, { rejectWithValue }) => {
    try {
      return await UserService.getUserActivities(userId);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to fetch user activities');
    }
  }
);

export const changePassword = createAsyncThunk(
  'profile/changePassword',
  async (
    {
      userId,
      currentPassword,
      newPassword,
    }: { userId: string; currentPassword: string; newPassword: string },
    { rejectWithValue }
  ) => {
    try {
      await UserService.changePassword(userId, currentPassword, newPassword);
      return true;
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to change password');
    }
  }
);

export const uploadAvatar = createAsyncThunk(
  'profile/uploadAvatar',
  async ({ userId, file }: { userId: string; file: File }, { rejectWithValue }) => {
    try {
      return await UserService.uploadAvatar(userId, file);
    } catch (error: any) {
      return rejectWithValue(error.message || 'Failed to upload avatar');
    }
  }
);

// Slice
const profileSlice = createSlice({
  name: 'profile',
  initialState,
  reducers: {
    clearProfileError: (state) => {
      state.error = null;
    },
  },
  extraReducers: (builder) => {
    builder
      // Fetch user profile
      .addCase(fetchUserProfile.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchUserProfile.fulfilled, (state, action) => {
        state.isLoading = false;
        state.profile = action.payload;
      })
      .addCase(fetchUserProfile.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Update user profile
      .addCase(updateUserProfile.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(updateUserProfile.fulfilled, (state, action) => {
        state.isLoading = false;
        state.profile = action.payload;
      })
      .addCase(updateUserProfile.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch reading stats
      .addCase(fetchReadingStats.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchReadingStats.fulfilled, (state, action) => {
        state.isLoading = false;
        state.readingStats = action.payload;
      })
      .addCase(fetchReadingStats.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch user bookmarks
      .addCase(fetchUserBookmarks.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchUserBookmarks.fulfilled, (state, action) => {
        state.isLoading = false;
        state.bookmarks = action.payload;
      })
      .addCase(fetchUserBookmarks.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Fetch user activities
      .addCase(fetchUserActivities.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(fetchUserActivities.fulfilled, (state, action) => {
        state.isLoading = false;
        state.activities = action.payload;
      })
      .addCase(fetchUserActivities.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Change password
      .addCase(changePassword.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(changePassword.fulfilled, (state) => {
        state.isLoading = false;
      })
      .addCase(changePassword.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      })
      // Upload avatar
      .addCase(uploadAvatar.pending, (state) => {
        state.isLoading = true;
        state.error = null;
      })
      .addCase(uploadAvatar.fulfilled, (state, action) => {
        state.isLoading = false;
        if (state.profile) {
          state.profile.avatar = action.payload.avatar_url;
        }
      })
      .addCase(uploadAvatar.rejected, (state, action) => {
        state.isLoading = false;
        state.error = action.payload as string;
      });
  },
});

export const { clearProfileError } = profileSlice.actions;
export default profileSlice.reducer;

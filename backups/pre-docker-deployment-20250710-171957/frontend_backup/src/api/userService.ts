import apiClient from './client';
import { Bookmark, ReadingStats, User, UserActivity } from '../types';

const UserService = {
  /**
   * Get user profile
   */
  getUserProfile: async (userId: string): Promise<User> => {
    const response = await apiClient.get<User>(`/users/${userId}/profile`);
    return response.data;
  },

  /**
   * Update user profile
   */
  updateUserProfile: async (userId: string, profileData: Partial<User>): Promise<User> => {
    const response = await apiClient.put<User>(`/users/${userId}/profile`, profileData);
    return response.data;
  },

  /**
   * Get user reading statistics
   */
  getReadingStats: async (userId: string): Promise<ReadingStats> => {
    const response = await apiClient.get<ReadingStats>(`/users/${userId}/reading-stats`);
    return response.data;
  },

  /**
   * Get user bookmarks
   */
  getUserBookmarks: async (userId: string): Promise<Bookmark[]> => {
    const response = await apiClient.get<Bookmark[]>(`/users/${userId}/bookmarks`);
    return response.data;
  },

  /**
   * Get user activities
   */
  getUserActivities: async (userId: string): Promise<UserActivity[]> => {
    const response = await apiClient.get<UserActivity[]>(`/users/${userId}/activities`);
    return response.data;
  },

  /**
   * Change user password
   */
  changePassword: async (
    userId: string,
    currentPassword: string,
    newPassword: string
  ): Promise<void> => {
    await apiClient.post(`/users/${userId}/change-password`, {
      currentPassword,
      newPassword,
    });
  },

  /**
   * Upload user avatar
   */
  uploadAvatar: async (userId: string, file: File): Promise<{ avatar_url: string }> => {
    const formData = new FormData();
    formData.append('avatar', file);

    const response = await apiClient.post<{ avatar_url: string }>(
      `/users/${userId}/avatar`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    );

    return response.data;
  },
};

export default UserService;

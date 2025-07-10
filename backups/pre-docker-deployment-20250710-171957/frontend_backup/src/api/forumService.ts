import apiClient from './client';
import { ForumCategory, ForumReply, ForumTopic, NewReply, NewTopic } from '../types';

const ForumService = {
  /**
   * Get all forum categories
   */
  getCategories: async (): Promise<ForumCategory[]> => {
    const response = await apiClient.get<ForumCategory[]>('/forum/categories');
    return response.data;
  },

  /**
   * Get topics by category
   */
  getTopicsByCategory: async (categoryId: string): Promise<ForumTopic[]> => {
    const response = await apiClient.get<ForumTopic[]>(`/forum/categories/${categoryId}/topics`);
    return response.data;
  },

  /**
   * Get topic by ID
   */
  getTopicById: async (topicId: string): Promise<ForumTopic> => {
    const response = await apiClient.get<ForumTopic>(`/forum/topics/${topicId}`);
    return response.data;
  },

  /**
   * Create new topic
   */
  createTopic: async (newTopic: NewTopic): Promise<ForumTopic> => {
    const response = await apiClient.post<ForumTopic>('/forum/topics', newTopic);
    return response.data;
  },

  /**
   * Create reply to topic
   */
  createReply: async (topicId: string, newReply: NewReply): Promise<ForumReply> => {
    const response = await apiClient.post<ForumReply>(`/forum/topics/${topicId}/replies`, newReply);
    return response.data;
  },

  /**
   * Vote on a reply
   */
  voteReply: async (replyId: string, vote: 'up' | 'down'): Promise<{ votes: number }> => {
    const response = await apiClient.post<{ votes: number }>(`/forum/replies/${replyId}/vote`, {
      vote,
    });
    return response.data;
  },

  /**
   * Search topics
   */
  searchTopics: async (query: string): Promise<ForumTopic[]> => {
    const response = await apiClient.get<ForumTopic[]>('/forum/search', {
      params: { q: query },
    });
    return response.data;
  },

  /**
   * Delete topic (admin or author only)
   */
  deleteTopic: async (topicId: string): Promise<void> => {
    await apiClient.delete(`/forum/topics/${topicId}`);
  },

  /**
   * Delete reply (admin or author only)
   */
  deleteReply: async (replyId: string): Promise<void> => {
    await apiClient.delete(`/forum/replies/${replyId}`);
  },
};

export default ForumService;

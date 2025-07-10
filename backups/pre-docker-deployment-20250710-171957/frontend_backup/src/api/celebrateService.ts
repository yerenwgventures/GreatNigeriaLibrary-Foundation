import apiClient from './client';
import { CelebrationCategory, CelebrationEntry, NewCelebrationEntry } from '../types';

const CelebrateService = {
  /**
   * Get featured entries
   */
  getFeaturedEntries: async (): Promise<CelebrationEntry[]> => {
    const response = await apiClient.get<CelebrationEntry[]>('/celebrate/featured');
    return response.data;
  },

  /**
   * Get entry by type and slug
   */
  getEntryByTypeAndSlug: async (
    type: 'person' | 'place' | 'event',
    slug: string
  ): Promise<CelebrationEntry> => {
    const response = await apiClient.get<CelebrationEntry>(`/celebrate/${type}/${slug}`);
    return response.data;
  },

  /**
   * Search entries
   */
  searchEntries: async (
    query: string,
    type?: 'person' | 'place' | 'event',
    categoryId?: string
  ): Promise<CelebrationEntry[]> => {
    const params: Record<string, string> = { q: query };
    if (type) params.type = type;
    if (categoryId) params.category = categoryId;

    const response = await apiClient.get<CelebrationEntry[]>('/celebrate/search', { params });
    return response.data;
  },

  /**
   * Get all categories
   */
  getCategories: async (): Promise<CelebrationCategory[]> => {
    const response = await apiClient.get<CelebrationCategory[]>('/celebrate/categories');
    return response.data;
  },

  /**
   * Submit new entry
   */
  submitEntry: async (entry: NewCelebrationEntry): Promise<{ entry_id: string }> => {
    const response = await apiClient.post<{ success: boolean; message: string; entry_id: string }>(
      '/celebrate/submit',
      entry
    );
    return { entry_id: response.data.entry_id };
  },

  /**
   * Vote for an entry
   */
  voteForEntry: async (entryId: string): Promise<{ votes: number }> => {
    const response = await apiClient.post<{ votes: number }>(`/celebrate/entries/${entryId}/vote`);
    return response.data;
  },

  /**
   * Comment on an entry
   */
  commentOnEntry: async (
    entryId: string,
    comment: string
  ): Promise<{ id: string; created_at: string }> => {
    const response = await apiClient.post<{ id: string; created_at: string }>(
      `/celebrate/entries/${entryId}/comments`,
      { comment }
    );
    return response.data;
  },

  /**
   * Get random entry
   */
  getRandomEntry: async (): Promise<CelebrationEntry> => {
    const response = await apiClient.get<CelebrationEntry>('/celebrate/random');
    return response.data;
  },
};

export default CelebrateService;

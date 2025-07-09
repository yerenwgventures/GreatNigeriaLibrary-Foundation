import apiClient from './client';
import { Resource, ResourceCategory } from '../types';

const ResourceService = {
  /**
   * Get all resource categories
   */
  getCategories: async (): Promise<ResourceCategory[]> => {
    const response = await apiClient.get<ResourceCategory[]>('/resources/categories');
    return response.data;
  },

  /**
   * Get resources by category
   */
  getResourcesByCategory: async (categoryId: string): Promise<Resource[]> => {
    const response = await apiClient.get<Resource[]>(`/resources/categories/${categoryId}/resources`);
    return response.data;
  },

  /**
   * Get resource by ID
   */
  getResourceById: async (resourceId: string): Promise<Resource> => {
    const response = await apiClient.get<Resource>(`/resources/${resourceId}`);
    return response.data;
  },

  /**
   * Download resource
   * This returns the URL to download the resource
   */
  getDownloadUrl: (resourceId: string): string => {
    return `${apiClient.defaults.baseURL}/resources/${resourceId}/download`;
  },

  /**
   * Track resource download
   * Call this when a user downloads a resource to track download count
   */
  trackDownload: async (resourceId: string): Promise<void> => {
    await apiClient.post(`/resources/${resourceId}/track-download`);
  },

  /**
   * Search resources
   */
  searchResources: async (query: string): Promise<Resource[]> => {
    const response = await apiClient.get<Resource[]>('/resources/search', {
      params: { q: query },
    });
    return response.data;
  },
};

export default ResourceService;

import axios from 'axios';
import { API_BASE_URL } from '../config';

// Define search result types
export interface SearchResult {
  id: number;
  title: string;
  description: string;
  type: string;
  url: string;
  imageUrl?: string;
  author?: string;
  createdAt?: string;
  tags?: string[];
  relevanceScore?: number;
}

export interface SearchResponse {
  results: SearchResult[];
  totalResults: number;
  hasMore: boolean;
  page: number;
  pageSize: number;
}

export interface SearchRequest {
  query: string;
  types?: string[];
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortOrder?: 'asc' | 'desc';
  tags?: string[];
  dateFrom?: string;
  dateTo?: string;
}

// Create the search service
const searchService = {
  search: async (request: SearchRequest): Promise<SearchResponse> => {
    const response = await axios.post(`${API_BASE_URL}/api/search`, request);
    return response.data;
  },

  getRecentSearches: async (): Promise<string[]> => {
    const response = await axios.get(`${API_BASE_URL}/api/search/recent`);
    return response.data;
  },

  saveRecentSearch: async (query: string): Promise<void> => {
    await axios.post(`${API_BASE_URL}/api/search/recent`, { query });
  },

  clearRecentSearches: async (): Promise<void> => {
    await axios.delete(`${API_BASE_URL}/api/search/recent`);
  },

  getPopularSearches: async (): Promise<{ query: string; count: number }[]> => {
    const response = await axios.get(`${API_BASE_URL}/api/search/popular`);
    return response.data;
  },
};

export default searchService;

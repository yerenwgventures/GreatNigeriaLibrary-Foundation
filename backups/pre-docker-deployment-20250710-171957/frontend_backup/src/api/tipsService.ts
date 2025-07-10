import axios from 'axios';
import { API_BASE_URL } from '../config';

// Define types
export interface Tip {
  id: number;
  title: string;
  content: string;
  category: string;
  trigger: string;
  triggerData: string;
  priority: number;
  imageUrl?: string;
  actionUrl?: string;
  actionText?: string;
  active: boolean;
  startDate?: string;
  endDate?: string;
  createdAt: string;
  updatedAt: string;
}

export interface TipRule {
  id: number;
  tipId: number;
  contextType: string;
  contextValue: string;
  condition?: string;
  priority: number;
}

export interface TipRequest {
  userId?: number;
  contextType: string;
  contextId: string;
  pageUrl: string;
  action?: string;
}

export interface TipResponse {
  tips: Tip[];
}

export interface TipFeedback {
  tipId: number;
  helpful: boolean;
  feedback?: string;
}

export interface TipStatistics {
  tipId: number;
  viewCount: number;
  dismissCount: number;
  clickCount: number;
  helpfulCount: number;
  effectivenessRate: number;
}

// Create the tips service
const tipsService = {
  // Get all tips
  getAllTips: async (): Promise<Tip[]> => {
    const response = await axios.get(`${API_BASE_URL}/api/tips`);
    return response.data;
  },
  
  // Get tip by ID
  getTipById: async (id: number): Promise<Tip> => {
    const response = await axios.get(`${API_BASE_URL}/api/tips/${id}`);
    return response.data;
  },
  
  // Get tips by category
  getTipsByCategory: async (category: string): Promise<Tip[]> => {
    const response = await axios.get(`${API_BASE_URL}/api/tips/category/${category}`);
    return response.data;
  },
  
  // Get contextual tips
  getContextualTips: async (request: TipRequest): Promise<TipResponse> => {
    const response = await axios.post(`${API_BASE_URL}/api/tips/contextual`, request);
    return response.data;
  },
  
  // Create tip (admin only)
  createTip: async (tip: Omit<Tip, 'id' | 'createdAt' | 'updatedAt'>): Promise<Tip> => {
    const response = await axios.post(`${API_BASE_URL}/api/tips`, tip);
    return response.data;
  },
  
  // Update tip (admin only)
  updateTip: async (id: number, tip: Partial<Tip>): Promise<Tip> => {
    const response = await axios.put(`${API_BASE_URL}/api/tips/${id}`, tip);
    return response.data;
  },
  
  // Delete tip (admin only)
  deleteTip: async (id: number): Promise<void> => {
    await axios.delete(`${API_BASE_URL}/api/tips/${id}`);
  },
  
  // Get tip rules by tip ID (admin only)
  getTipRulesByTipId: async (tipId: number): Promise<TipRule[]> => {
    const response = await axios.get(`${API_BASE_URL}/api/tips/rules/${tipId}`);
    return response.data;
  },
  
  // Create tip rule (admin only)
  createTipRule: async (rule: Omit<TipRule, 'id'>): Promise<TipRule> => {
    const response = await axios.post(`${API_BASE_URL}/api/tips/rules`, rule);
    return response.data;
  },
  
  // Update tip rule (admin only)
  updateTipRule: async (id: number, rule: Partial<TipRule>): Promise<TipRule> => {
    const response = await axios.put(`${API_BASE_URL}/api/tips/rules/${id}`, rule);
    return response.data;
  },
  
  // Delete tip rule (admin only)
  deleteTipRule: async (id: number): Promise<void> => {
    await axios.delete(`${API_BASE_URL}/api/tips/rules/${id}`);
  },
  
  // Record tip view
  recordTipView: async (id: number): Promise<void> => {
    await axios.post(`${API_BASE_URL}/api/tips/view/${id}`);
  },
  
  // Record tip dismiss
  recordTipDismiss: async (id: number): Promise<void> => {
    await axios.post(`${API_BASE_URL}/api/tips/dismiss/${id}`);
  },
  
  // Record tip click
  recordTipClick: async (id: number): Promise<void> => {
    await axios.post(`${API_BASE_URL}/api/tips/click/${id}`);
  },
  
  // Submit tip feedback
  submitTipFeedback: async (feedback: TipFeedback): Promise<void> => {
    await axios.post(`${API_BASE_URL}/api/tips/feedback`, feedback);
  },
  
  // Get tip statistics (admin only)
  getTipStatistics: async (id: number): Promise<TipStatistics> => {
    const response = await axios.get(`${API_BASE_URL}/api/tips/statistics/${id}`);
    return response.data;
  },
  
  // Get all tip statistics (admin only)
  getAllTipStatistics: async (): Promise<TipStatistics[]> => {
    const response = await axios.get(`${API_BASE_URL}/api/tips/statistics`);
    return response.data;
  },
};

export default tipsService;

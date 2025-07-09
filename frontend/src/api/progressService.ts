import axios from 'axios';
import { client } from './client';

// Define types for API responses
interface Activity {
  id: number;
  type: string;
  name: string;
  progress: number;
  date: string;
}

interface UserProgress {
  overallCompletion: number;
  pointsEarned: number;
  streak: number;
  level: number;
  recentActivities: Activity[];
}

interface Milestone {
  id: number;
  name: string;
  description: string;
  completed: boolean;
  date?: string;
  progress?: number;
  icon: string;
}

interface Achievement {
  id: number;
  name: string;
  description: string;
  earned: boolean;
  date?: string;
  progress?: number;
  icon: string;
}

interface HistoricalData {
  month: string;
  progress: number;
  activities: number;
}

interface SkillData {
  name: string;
  value: number;
}

// Define request types
interface LogActivityRequest {
  type: string;
  name: string;
  progress: number;
}

interface UpdateProgressRequest {
  progress: number;
}

// Progress API service
const progressService = {
  // Get user progress
  getUserProgress: async (): Promise<UserProgress> => {
    const response = await client.get('/api/progress/user');
    return response.data;
  },

  // Get milestones
  getMilestones: async (): Promise<Milestone[]> => {
    const response = await client.get('/api/progress/milestones');
    return response.data;
  },

  // Get achievements
  getAchievements: async (): Promise<Achievement[]> => {
    const response = await client.get('/api/progress/achievements');
    return response.data;
  },

  // Get historical data
  getHistoricalData: async (): Promise<HistoricalData[]> => {
    const response = await client.get('/api/progress/history');
    return response.data;
  },

  // Get skills data
  getSkillsData: async (): Promise<SkillData[]> => {
    const response = await client.get('/api/progress/skills');
    return response.data;
  },

  // Update milestone progress
  updateMilestoneProgress: async (milestoneId: number, progress: number): Promise<void> => {
    const request: UpdateProgressRequest = { progress };
    await client.put(`/api/progress/milestones/${milestoneId}`, request);
  },

  // Update achievement progress
  updateAchievementProgress: async (achievementId: number, progress: number): Promise<void> => {
    const request: UpdateProgressRequest = { progress };
    await client.put(`/api/progress/achievements/${achievementId}`, request);
  },

  // Log activity
  logActivity: async (type: string, name: string, progress: number): Promise<void> => {
    const request: LogActivityRequest = { type, name, progress };
    await client.post('/api/progress/activities', request);
  }
};

export default progressService;

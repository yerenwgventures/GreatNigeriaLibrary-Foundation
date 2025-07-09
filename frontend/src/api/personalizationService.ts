import axios from 'axios';
import { API_BASE_URL } from '../config';

// Define types
export interface LearningStyle {
  id?: number;
  userId?: number;
  visual: number;
  auditory: number;
  readWrite: number;
  kinesthetic: number;
  social: number;
  solitary: number;
  logical: number;
  primaryStyle?: string;
  secondaryStyle?: string;
  assessmentTaken: boolean;
  assessmentDate?: string;
  lastUpdated?: string;
}

export interface LearningPreference {
  id?: number;
  userId?: number;
  preferredTopics: string[];
  avoidTopics: string[];
  preferredFormats: string[];
  difficultyLevel: string;
  learningGoals: string[];
  lastUpdated?: string;
}

export interface PersonalizedPath {
  id?: number;
  userId?: number;
  name: string;
  description: string;
  createdAt?: string;
  updatedAt?: string;
  isActive: boolean;
  completionRate: number;
}

export interface PathItem {
  id?: number;
  pathId?: number;
  itemType: string;
  itemId: number;
  title: string;
  description: string;
  order: number;
  isCompleted: boolean;
  completedAt?: string;
  estimatedDuration: number;
}

export interface AssessmentQuestion {
  id: number;
  question: string;
  options: string[];
  styleDimension: string;
  weight: number;
  isActive: boolean;
}

export interface AssessmentResponse {
  questionId: number;
  selectedOption: number;
}

export interface ContentRecommendation {
  id?: number;
  userId?: number;
  contentType: string;
  contentId: number;
  title: string;
  description: string;
  recommendationScore: number;
  reasonCodes: string[];
  isViewed: boolean;
  isSaved: boolean;
  isRejected: boolean;
  createdAt?: string;
}

export interface UserPerformance {
  id?: number;
  userId?: number;
  overallScore: number;
  topicScores: Record<string, number>;
  quizzesTaken: number;
  quizAvgScore: number;
  contentCompleted: number;
  lastUpdated?: string;
}

export interface AssessmentResult {
  userId: number;
  learningStyle: LearningStyle;
  recommendations: ContentRecommendation[];
  suggestedPaths: PersonalizedPath[];
}

export interface PersonalizationRequest {
  contentType?: string;
  topic?: string;
  count?: number;
}

export interface PersonalizationResponse {
  recommendations: ContentRecommendation[];
  hasMoreResults: boolean;
  userHasStyle: boolean;
}

// Create the personalization service
const personalizationService = {
  // Learning Style methods
  getLearningStyle: async (): Promise<LearningStyle> => {
    const response = await axios.get(`${API_BASE_URL}/api/personalization/learning-style`);
    return response.data;
  },
  
  saveLearningStyle: async (style: LearningStyle): Promise<LearningStyle> => {
    const response = await axios.post(`${API_BASE_URL}/api/personalization/learning-style`, style);
    return response.data;
  },
  
  // Learning Preference methods
  getLearningPreference: async (): Promise<LearningPreference> => {
    const response = await axios.get(`${API_BASE_URL}/api/personalization/preferences`);
    return response.data;
  },
  
  saveLearningPreference: async (pref: LearningPreference): Promise<LearningPreference> => {
    const response = await axios.post(`${API_BASE_URL}/api/personalization/preferences`, pref);
    return response.data;
  },
  
  // Assessment methods
  getAssessmentQuestions: async (): Promise<AssessmentQuestion[]> => {
    const response = await axios.get(`${API_BASE_URL}/api/personalization/assessment/questions`);
    return response.data;
  },
  
  submitAssessment: async (responses: AssessmentResponse[]): Promise<AssessmentResult> => {
    const response = await axios.post(`${API_BASE_URL}/api/personalization/assessment/submit`, responses);
    return response.data;
  },
  
  // Personalized Path methods
  getPersonalizedPaths: async (): Promise<PersonalizedPath[]> => {
    const response = await axios.get(`${API_BASE_URL}/api/personalization/paths`);
    return response.data;
  },
  
  getPersonalizedPathWithItems: async (pathId: number): Promise<{ path: PersonalizedPath; items: PathItem[] }> => {
    const response = await axios.get(`${API_BASE_URL}/api/personalization/paths/${pathId}`);
    return response.data;
  },
  
  createPersonalizedPath: async (path: PersonalizedPath, items: PathItem[]): Promise<PersonalizedPath> => {
    const response = await axios.post(`${API_BASE_URL}/api/personalization/paths`, { path, items });
    return response.data;
  },
  
  updatePersonalizedPath: async (pathId: number, path: Partial<PersonalizedPath>): Promise<PersonalizedPath> => {
    const response = await axios.put(`${API_BASE_URL}/api/personalization/paths/${pathId}`, path);
    return response.data;
  },
  
  deletePersonalizedPath: async (pathId: number): Promise<void> => {
    await axios.delete(`${API_BASE_URL}/api/personalization/paths/${pathId}`);
  },
  
  // Path Item methods
  addPathItem: async (pathId: number, item: Omit<PathItem, 'pathId'>): Promise<PathItem> => {
    const response = await axios.post(`${API_BASE_URL}/api/personalization/paths/${pathId}/items`, item);
    return response.data;
  },
  
  updatePathItem: async (itemId: number, item: Partial<PathItem>): Promise<PathItem> => {
    const response = await axios.put(`${API_BASE_URL}/api/personalization/paths/items/${itemId}`, item);
    return response.data;
  },
  
  deletePathItem: async (itemId: number): Promise<void> => {
    await axios.delete(`${API_BASE_URL}/api/personalization/paths/items/${itemId}`);
  },
  
  markPathItemComplete: async (itemId: number, completed: boolean): Promise<void> => {
    await axios.put(`${API_BASE_URL}/api/personalization/paths/items/${itemId}/complete`, { completed });
  },
  
  // Recommendation methods
  getRecommendations: async (request: PersonalizationRequest): Promise<PersonalizationResponse> => {
    const response = await axios.post(`${API_BASE_URL}/api/personalization/recommendations`, request);
    return response.data;
  },
  
  updateRecommendationStatus: async (recId: number, viewed: boolean, saved: boolean, rejected: boolean): Promise<void> => {
    await axios.put(`${API_BASE_URL}/api/personalization/recommendations/${recId}/status`, { viewed, saved, rejected });
  },
  
  // User Performance methods
  getUserPerformance: async (): Promise<UserPerformance> => {
    const response = await axios.get(`${API_BASE_URL}/api/personalization/performance`);
    return response.data;
  },
  
  updateUserPerformance: async (performance: UserPerformance): Promise<UserPerformance> => {
    const response = await axios.put(`${API_BASE_URL}/api/personalization/performance`, performance);
    return response.data;
  },
};

export default personalizationService;

import axios from 'axios';
import { Feature } from '../components/features/FeatureToggle';
import { API_URL } from '../config';

// Create axios instance with base URL
const api = axios.create({
  baseURL: API_URL
});

// Add request interceptor to include auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Get all features for the current user
const getFeatures = async (): Promise<Feature[]> => {
  try {
    const response = await api.get('/features');
    return response.data;
  } catch (error) {
    // For development/demo purposes, return mock data if API fails
    if (process.env.NODE_ENV === 'development') {
      console.warn('Using mock feature data');
      return getMockFeatures();
    }
    throw error;
  }
};

// Update a feature's enabled status
const updateFeature = async (featureId: string, enabled: boolean): Promise<Feature> => {
  try {
    const response = await api.patch(`/features/${featureId}`, { enabled });
    return response.data;
  } catch (error) {
    // For development/demo purposes, return mock data if API fails
    if (process.env.NODE_ENV === 'development') {
      console.warn('Using mock feature update');
      return getMockFeatureUpdate(featureId, enabled);
    }
    throw error;
  }
};

// Reset features to default
const resetFeatures = async (): Promise<Feature[]> => {
  try {
    const response = await api.post('/features/reset');
    return response.data;
  } catch (error) {
    // For development/demo purposes, return mock data if API fails
    if (process.env.NODE_ENV === 'development') {
      console.warn('Using mock feature reset');
      return getMockFeatures();
    }
    throw error;
  }
};

// Mock data for development/demo purposes
const getMockFeatures = (): Feature[] => {
  return [
    {
      id: 'social_networking',
      name: 'Social Networking',
      description: 'Connect with other users, follow their activities, and build your network.',
      icon: '👥',
      enabled: true,
      category: 'Social',
      beta: false,
      premium: false,
      dependencies: [],
      incompatibleWith: []
    },
    {
      id: 'friends',
      name: 'Friends System',
      description: 'Send and receive friend requests, manage your friends list.',
      icon: '🤝',
      enabled: true,
      category: 'Social',
      beta: false,
      premium: false,
      dependencies: ['social_networking'],
      incompatibleWith: []
    },
    {
      id: 'groups',
      name: 'Groups & Communities',
      description: 'Create and join groups around specific interests or topics.',
      icon: '👪',
      enabled: true,
      category: 'Social',
      beta: false,
      premium: false,
      dependencies: ['social_networking'],
      incompatibleWith: []
    },
    {
      id: 'messaging',
      name: 'Private Messaging',
      description: 'Send private messages to other users.',
      icon: '✉️',
      enabled: true,
      category: 'Communication',
      beta: false,
      premium: false,
      dependencies: ['social_networking'],
      incompatibleWith: []
    },
    {
      id: 'video_calls',
      name: 'Video Calls',
      description: 'Make video calls with other users.',
      icon: '📹',
      enabled: false,
      category: 'Communication',
      beta: true,
      premium: true,
      dependencies: ['messaging'],
      incompatibleWith: []
    },
    {
      id: 'livestreaming',
      name: 'Livestreaming',
      description: 'Stream live video to your followers.',
      icon: '🎬',
      enabled: true,
      category: 'Content',
      beta: false,
      premium: false,
      dependencies: [],
      incompatibleWith: []
    },
    {
      id: 'gifting',
      name: 'Virtual Gifts',
      description: 'Send and receive virtual gifts during livestreams.',
      icon: '🎁',
      enabled: true,
      category: 'Content',
      beta: false,
      premium: false,
      dependencies: ['livestreaming'],
      incompatibleWith: []
    },
    {
      id: 'content_creation',
      name: 'Content Creation',
      description: 'Create and publish articles, posts, and other content.',
      icon: '✏️',
      enabled: true,
      category: 'Content',
      beta: false,
      premium: false,
      dependencies: [],
      incompatibleWith: []
    },
    {
      id: 'marketplace',
      name: 'Marketplace',
      description: 'Buy and sell items within the platform.',
      icon: '🛒',
      enabled: false,
      category: 'Economic',
      beta: true,
      premium: false,
      dependencies: [],
      incompatibleWith: []
    },
    {
      id: 'jobs',
      name: 'Jobs & Opportunities',
      description: 'Find and post job opportunities.',
      icon: '💼',
      enabled: false,
      category: 'Economic',
      beta: true,
      premium: false,
      dependencies: [],
      incompatibleWith: []
    },
    {
      id: 'events',
      name: 'Events',
      description: 'Create and join events.',
      icon: '📅',
      enabled: true,
      category: 'Social',
      beta: false,
      premium: false,
      dependencies: [],
      incompatibleWith: []
    },
    {
      id: 'advanced_analytics',
      name: 'Advanced Analytics',
      description: 'Get detailed analytics about your content and engagement.',
      icon: '📊',
      enabled: false,
      category: 'Premium',
      beta: false,
      premium: true,
      dependencies: [],
      incompatibleWith: []
    },
    {
      id: 'dark_mode',
      name: 'Dark Mode',
      description: 'Enable dark mode for the platform.',
      icon: '🌙',
      enabled: true,
      category: 'Appearance',
      beta: false,
      premium: false,
      dependencies: [],
      incompatibleWith: []
    },
    {
      id: 'notifications',
      name: 'Notifications',
      description: 'Receive notifications about activities relevant to you.',
      icon: '🔔',
      enabled: true,
      category: 'Communication',
      beta: false,
      premium: false,
      dependencies: [],
      incompatibleWith: []
    },
    {
      id: 'ai_recommendations',
      name: 'AI Recommendations',
      description: 'Get personalized content recommendations powered by AI.',
      icon: '🤖',
      enabled: false,
      category: 'Premium',
      beta: true,
      premium: true,
      dependencies: [],
      incompatibleWith: []
    }
  ];
};

// Mock feature update for development/demo purposes
const getMockFeatureUpdate = (featureId: string, enabled: boolean): Feature => {
  const features = getMockFeatures();
  const feature = features.find(f => f.id === featureId);
  
  if (!feature) {
    throw new Error(`Feature with ID ${featureId} not found`);
  }
  
  return {
    ...feature,
    enabled
  };
};

const featuresService = {
  getFeatures,
  updateFeature,
  resetFeatures
};

export default featuresService;

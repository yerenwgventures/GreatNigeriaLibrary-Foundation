import { client } from './client';

// Define the badge types
export interface Badge {
  id: string;
  name: string;
  description: string;
  category: string;
  level: string;
  imageUrl?: string;
  isPublic: boolean;
  isRare: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface UserBadge {
  id: string;
  userId: string;
  badgeId: string;
  badge: Badge;
  awardedAt: string;
  awardedBy?: string;
  reason?: string;
  isHidden: boolean;
  isPublic: boolean;
  progress?: number;
  requiredProgress?: number;
  isNew?: boolean;
}

// Get all badges
export const getAllBadges = async (): Promise<Badge[]> => {
  const response = await client.get('/badges');
  return response.data;
};

// Get user badges
export const getUserBadges = async (userId: string): Promise<UserBadge[]> => {
  const response = await client.get(`/users/${userId}/badges`);
  return response.data;
};

// Get badge by ID
export const getBadgeById = async (badgeId: string): Promise<Badge> => {
  const response = await client.get(`/badges/${badgeId}`);
  return response.data;
};

// Mark badge as viewed (no longer new)
export const markBadgeAsViewed = async (userId: string, badgeId: string): Promise<void> => {
  await client.post(`/users/${userId}/badges/${badgeId}/viewed`);
};

// Hide/unhide badge
export const toggleBadgeVisibility = async (userId: string, badgeId: string, isHidden: boolean): Promise<void> => {
  await client.patch(`/users/${userId}/badges/${badgeId}`, { isHidden });
};

// Set badge as featured
export const setFeaturedBadge = async (userId: string, badgeId: string): Promise<void> => {
  await client.post(`/users/${userId}/badges/${badgeId}/featured`);
};

// Get user's badge progress
export const getBadgeProgress = async (userId: string, badgeId: string): Promise<{ progress: number; requiredProgress: number }> => {
  const response = await client.get(`/users/${userId}/badges/${badgeId}/progress`);
  return response.data;
};

export default {
  getAllBadges,
  getUserBadges,
  getBadgeById,
  markBadgeAsViewed,
  toggleBadgeVisibility,
  setFeaturedBadge,
  getBadgeProgress
};

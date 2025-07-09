import { client } from './client';

// Affiliate interfaces
export interface ReferralCode {
  id: string;
  userId: string;
  code: string;
  description?: string;
  isActive: boolean;
  createdAt: string;
  updatedAt: string;
  usageCount: number;
  conversionCount: number;
}

export interface Referral {
  id: string;
  referrerUserId: string;
  referredUserId: string;
  referralCodeId: string;
  referralCode: string;
  status: 'pending' | 'converted' | 'expired';
  referralType: 'membership' | 'product';
  productId?: string;
  membershipPlanId?: string;
  conversionValue?: number;
  createdAt: string;
  updatedAt: string;
  convertedAt?: string;
  expiresAt?: string;
}

export interface Commission {
  id: string;
  userId: string;
  referralId: string;
  amount: number;
  currency: string;
  status: 'pending' | 'approved' | 'paid' | 'rejected';
  commissionType: 'membership' | 'product';
  productId?: string;
  membershipPlanId?: string;
  tier: number;
  createdAt: string;
  updatedAt: string;
  approvedAt?: string;
  paidAt?: string;
  rejectedAt?: string;
  rejectionReason?: string;
}

export interface AffiliateStats {
  totalReferrals: number;
  activeReferrals: number;
  pendingReferrals: number;
  totalCommissions: number;
  pendingCommissions: number;
  paidCommissions: number;
  conversionRate: number;
  totalEarnings: number;
  currentMonthEarnings: number;
  referralsByMonth: {
    month: string;
    count: number;
  }[];
  commissionsByMonth: {
    month: string;
    amount: number;
  }[];
}

export interface AffiliateSettings {
  commissionRates: {
    registration: number;
    purchase: number;
    subscription: number;
    contentSale: number;
  };
  multiTierLevels: number;
  tierRates: {
    [key: number]: number;
  };
  minimumPayout: number;
  referralValidDays: number;
  bonusThresholds: {
    threshold: number;
    bonusAmount: number;
  }[];
}

export interface MembershipPlan {
  id: string;
  name: string;
  price: number;
  durationDays: number;
  affiliateCommissionPercentage: number;
  description?: string;
  isActive: boolean;
}

export interface ProductAffiliateSettings {
  id: string;
  productId: string;
  sellerId: string;
  isAffiliateEnabled: boolean;
  commissionPercentage: number;
  cookieDurationDays: number;
  termsAndConditions?: string;
  createdAt: string;
  updatedAt: string;
}

export interface AffiliateProduct {
  id: string;
  title: string;
  description: string;
  price: number;
  currency: string;
  category: string;
  mediaUrls: string[];
  seller: {
    id: string;
    username: string;
    name: string;
    profileImage?: string;
  };
  affiliateSettings: ProductAffiliateSettings;
}

export interface CreateReferralCodeRequest {
  description?: string;
}

export interface UpdateReferralCodeRequest {
  description?: string;
  isActive?: boolean;
}

export interface WithdrawCommissionRequest {
  amount: number;
  paymentMethod: 'wallet' | 'bank_transfer' | 'mobile_money';
  accountDetails?: {
    accountName?: string;
    accountNumber?: string;
    bankCode?: string;
    bankName?: string;
    phoneNumber?: string;
    provider?: string;
  };
}

export interface UpdateProductAffiliateSettingsRequest {
  isAffiliateEnabled: boolean;
  commissionPercentage: number;
  cookieDurationDays: number;
  termsAndConditions?: string;
}

// API functions for referral codes
export const getReferralCodes = async (): Promise<ReferralCode[]> => {
  const response = await client.get('/affiliate/referral-codes');
  return response.data;
};

export const getReferralCodeById = async (id: string): Promise<ReferralCode> => {
  const response = await client.get(`/affiliate/referral-codes/${id}`);
  return response.data;
};

export const createReferralCode = async (request: CreateReferralCodeRequest): Promise<ReferralCode> => {
  const response = await client.post('/affiliate/referral-codes', request);
  return response.data;
};

export const updateReferralCode = async (id: string, request: UpdateReferralCodeRequest): Promise<ReferralCode> => {
  const response = await client.put(`/affiliate/referral-codes/${id}`, request);
  return response.data;
};

export const deleteReferralCode = async (id: string): Promise<void> => {
  await client.delete(`/affiliate/referral-codes/${id}`);
};

// API functions for referrals
export const getReferrals = async (
  page = 1,
  limit = 10,
  status?: string,
  type?: 'membership' | 'product'
): Promise<{ referrals: Referral[]; total: number; page: number; totalPages: number }> => {
  const response = await client.get('/affiliate/referrals', {
    params: { page, limit, status, type }
  });
  return response.data;
};

export const getReferralById = async (id: string): Promise<Referral> => {
  const response = await client.get(`/affiliate/referrals/${id}`);
  return response.data;
};

// API functions for commissions
export const getCommissions = async (
  page = 1,
  limit = 10,
  status?: string,
  type?: 'membership' | 'product'
): Promise<{ commissions: Commission[]; total: number; page: number; totalPages: number }> => {
  const response = await client.get('/affiliate/commissions', {
    params: { page, limit, status, type }
  });
  return response.data;
};

export const getCommissionById = async (id: string): Promise<Commission> => {
  const response = await client.get(`/affiliate/commissions/${id}`);
  return response.data;
};

export const withdrawCommission = async (request: WithdrawCommissionRequest): Promise<{ 
  transactionId: string; 
  status: string;
  amount: number;
  fee: number;
  netAmount: number;
  estimatedArrivalDate: string;
}> => {
  const response = await client.post('/affiliate/commissions/withdraw', request);
  return response.data;
};

// API functions for affiliate stats
export const getAffiliateStats = async (type?: 'membership' | 'product'): Promise<AffiliateStats> => {
  const response = await client.get('/affiliate/stats', {
    params: { type }
  });
  return response.data;
};

// API functions for affiliate settings
export const getAffiliateSettings = async (): Promise<AffiliateSettings> => {
  const response = await client.get('/affiliate/settings');
  return response.data;
};

// API function to apply a referral code
export const applyReferralCode = async (code: string): Promise<{ success: boolean; message: string }> => {
  const response = await client.post('/affiliate/apply-code', { code });
  return response.data;
};

// API function to get referral link
export const getReferralLink = async (codeId?: string, productId?: string): Promise<{ link: string; code: string }> => {
  const response = await client.get('/affiliate/referral-link', {
    params: { codeId, productId }
  });
  return response.data;
};

// Membership affiliate functions
export const getMembershipPlans = async (): Promise<MembershipPlan[]> => {
  const response = await client.get('/affiliate/membership/plans');
  return response.data;
};

export const getMembershipAffiliateStats = async (): Promise<AffiliateStats> => {
  const response = await client.get('/affiliate/membership/stats');
  return response.data;
};

// Marketplace affiliate functions
export const getSellerProducts = async (): Promise<AffiliateProduct[]> => {
  const response = await client.get('/affiliate/marketplace/my-products');
  return response.data;
};

export const updateProductAffiliateSettings = async (
  productId: string, 
  settings: UpdateProductAffiliateSettingsRequest
): Promise<ProductAffiliateSettings> => {
  const response = await client.post(`/affiliate/marketplace/products/${productId}/settings`, settings);
  return response.data;
};

export const getProductsWithAffiliatePrograms = async (
  page = 1,
  limit = 10,
  category?: string
): Promise<{ products: AffiliateProduct[]; total: number; page: number; totalPages: number }> => {
  const response = await client.get('/affiliate/marketplace/products', {
    params: { page, limit, category }
  });
  return response.data;
};

export const getMarketplaceAffiliateStats = async (): Promise<AffiliateStats> => {
  const response = await client.get('/affiliate/marketplace/stats');
  return response.data;
};

// Combined dashboard
export const getAffiliateDashboard = async (): Promise<{
  membership: AffiliateStats;
  marketplace: AffiliateStats;
  combined: AffiliateStats;
}> => {
  const response = await client.get('/affiliate/dashboard');
  return response.data;
};

export default {
  getReferralCodes,
  getReferralCodeById,
  createReferralCode,
  updateReferralCode,
  deleteReferralCode,
  getReferrals,
  getReferralById,
  getCommissions,
  getCommissionById,
  withdrawCommission,
  getAffiliateStats,
  getAffiliateSettings,
  applyReferralCode,
  getReferralLink,
  getMembershipPlans,
  getMembershipAffiliateStats,
  getSellerProducts,
  updateProductAffiliateSettings,
  getProductsWithAffiliatePrograms,
  getMarketplaceAffiliateStats,
  getAffiliateDashboard
};

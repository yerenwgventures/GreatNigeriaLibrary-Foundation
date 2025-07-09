import apiClient from './client';

// Types
export interface Stream {
  id: number;
  creatorId: number;
  title: string;
  description: string;
  thumbnailUrl: string;
  status: 'scheduled' | 'live' | 'ended';
  scheduledStart: string;
  actualStart?: string;
  endTime?: string;
  viewerCount: number;
  peakViewerCount: number;
  totalGiftsValue: number;
  playbackUrl: string;
  isPrivate: boolean;
  allowGifting: boolean;
  categories: string;
  tags: string;
  createdAt: string;
  updatedAt: string;
}

export interface StreamViewer {
  id: number;
  streamId: number;
  userId: number;
  joinTime: string;
  leaveTime?: string;
  duration: number;
  isActive: boolean;
}

export interface CoinPackage {
  id: number;
  name: string;
  coinsAmount: number;
  priceNaira: number;
  bonusCoins: number;
  isActive: boolean;
  isPromotional: boolean;
  promotionEnds?: string;
  description: string;
  imageUrl: string;
}

export interface VirtualCurrency {
  id: number;
  userId: number;
  balance: number;
  totalBought: number;
  totalSpent: number;
}

export interface Gift {
  id: number;
  streamId: number;
  senderId: number;
  recipientId: number;
  giftId: number;
  giftName: string;
  coinsAmount: number;
  nairaValue: number;
  message?: string;
  isAnonymous: boolean;
  isHighlighted: boolean;
  comboCount: number;
  creatorRevenuePercent: number;
  creatorRevenueAmount: number;
  platformRevenueAmount: number;
  createdAt: string;
}

export interface GifterRanking {
  id: number;
  userId: number;
  streamId?: number;
  rankingPeriod: 'daily' | 'weekly' | 'monthly' | 'all_time';
  periodStart: string;
  periodEnd: string;
  totalGifts: number;
  totalCoins: number;
  totalNairaValue: number;
  rank: number;
  previousRank?: number;
  badgeLevel: 'bronze' | 'silver' | 'gold' | 'platinum' | 'diamond';
}

export interface CreatorRevenue {
  id: number;
  creatorId: number;
  streamId?: number;
  period: 'daily' | 'weekly' | 'monthly';
  periodStart: string;
  periodEnd: string;
  totalGifts: number;
  totalCoins: number;
  totalNairaValue: number;
  platformFee: number;
  netRevenue: number;
  isPaid: boolean;
  paymentDate?: string;
  paymentReference?: string;
}

export interface WithdrawalRequest {
  id: number;
  creatorId: number;
  amount: number;
  status: 'pending' | 'approved' | 'rejected' | 'completed';
  bankName: string;
  accountNumber: string;
  accountName: string;
  processedDate?: string;
  processedBy?: number;
  notes?: string;
  transactionReference?: string;
  createdAt: string;
}

export interface PurchaseCoinsRequest {
  packageId: number;
  paymentId: number;
}

export interface SendGiftRequest {
  streamId: number;
  recipientId: number;
  giftId: number;
  coinsAmount: number;
  message?: string;
  isAnonymous?: boolean;
}

export interface CreateStreamRequest {
  title: string;
  description: string;
  thumbnailUrl: string;
  scheduledStart: string;
  isPrivate: boolean;
  categories?: string;
  tags?: string;
}

export interface WithdrawalRequestData {
  amount: number;
  bankName: string;
  accountNumber: string;
  accountName: string;
}

// Livestream Service
const livestreamService = {
  // Stream operations
  getActiveStreams: async (page = 1, limit = 10): Promise<{ streams: Stream[], pagination: any }> => {
    const response = await apiClient.get(`/livestream/active?page=${page}&limit=${limit}`);
    return response.data;
  },

  getStreamById: async (streamId: number): Promise<Stream> => {
    const response = await apiClient.get(`/livestream/${streamId}`);
    return response.data;
  },

  createStream: async (streamData: CreateStreamRequest): Promise<Stream> => {
    const response = await apiClient.post('/livestream', streamData);
    return response.data;
  },

  updateStream: async (streamId: number, streamData: Partial<CreateStreamRequest>): Promise<Stream> => {
    const response = await apiClient.put(`/livestream/${streamId}`, streamData);
    return response.data;
  },

  startStream: async (streamId: number): Promise<Stream> => {
    const response = await apiClient.post(`/livestream/${streamId}/start`);
    return response.data;
  },

  endStream: async (streamId: number): Promise<Stream> => {
    const response = await apiClient.post(`/livestream/${streamId}/end`);
    return response.data;
  },

  getStreamViewers: async (streamId: number, page = 1, limit = 10): Promise<{ viewers: StreamViewer[], pagination: any }> => {
    const response = await apiClient.get(`/livestream/${streamId}/viewers?page=${page}&limit=${limit}`);
    return response.data;
  },

  // Virtual Currency operations
  getCoinPackages: async (): Promise<CoinPackage[]> => {
    const response = await apiClient.get('/currency/packages');
    return response.data;
  },

  getUserBalance: async (userId: number): Promise<VirtualCurrency> => {
    const response = await apiClient.get(`/currency/balance/${userId}`);
    return response.data;
  },

  purchaseCoins: async (data: PurchaseCoinsRequest): Promise<{ message: string, balance: VirtualCurrency }> => {
    const response = await apiClient.post('/currency/purchase', data);
    return response.data;
  },

  getUserTransactions: async (userId: number, page = 1, limit = 10): Promise<{ transactions: any[], pagination: any }> => {
    const response = await apiClient.get(`/currency/transactions/${userId}?page=${page}&limit=${limit}`);
    return response.data;
  },

  // Gift operations
  sendGift: async (data: SendGiftRequest): Promise<Gift> => {
    const response = await apiClient.post('/gifts/send', data);
    return response.data;
  },

  getStreamGifts: async (streamId: number, page = 1, limit = 10): Promise<{ gifts: Gift[], pagination: any }> => {
    const response = await apiClient.get(`/gifts/stream/${streamId}?page=${page}&limit=${limit}`);
    return response.data;
  },

  getUserSentGifts: async (userId: number, page = 1, limit = 10): Promise<{ gifts: Gift[], pagination: any }> => {
    const response = await apiClient.get(`/gifts/user/${userId}/sent?page=${page}&limit=${limit}`);
    return response.data;
  },

  getUserReceivedGifts: async (userId: number, page = 1, limit = 10): Promise<{ gifts: Gift[], pagination: any }> => {
    const response = await apiClient.get(`/gifts/user/${userId}/received?page=${page}&limit=${limit}`);
    return response.data;
  },

  // Ranking operations
  getStreamRankings: async (streamId: number, period = 'daily', limit = 10): Promise<{ rankings: GifterRanking[], period: string, streamId: number }> => {
    const response = await apiClient.get(`/rankings/stream/${streamId}?period=${period}&limit=${limit}`);
    return response.data;
  },

  getGlobalRankings: async (period = 'daily', limit = 10): Promise<{ rankings: GifterRanking[], period: string }> => {
    const response = await apiClient.get(`/rankings/global?period=${period}&limit=${limit}`);
    return response.data;
  },

  getUserRanking: async (userId: number, streamId?: number, period = 'daily'): Promise<GifterRanking> => {
    const streamParam = streamId ? `&streamId=${streamId}` : '';
    const response = await apiClient.get(`/rankings/user/${userId}?period=${period}${streamParam}`);
    return response.data;
  },

  // Revenue operations
  getCreatorRevenue: async (userId: number, period?: string, page = 1, limit = 10): Promise<{ revenues: CreatorRevenue[], pagination: any }> => {
    const periodParam = period ? `&period=${period}` : '';
    const response = await apiClient.get(`/revenue/creator/${userId}?page=${page}&limit=${limit}${periodParam}`);
    return response.data;
  },

  getRevenueSummary: async (userId: number): Promise<any> => {
    const response = await apiClient.get(`/revenue/summary/${userId}`);
    return response.data;
  },

  requestWithdrawal: async (data: WithdrawalRequestData): Promise<WithdrawalRequest> => {
    const response = await apiClient.post('/revenue/withdraw', data);
    return response.data;
  },

  // WebSocket connection
  getWebSocketUrl: (userId: number, streamId?: number): string => {
    const baseUrl = process.env.REACT_APP_WS_URL || window.location.origin.replace('http', 'ws');
    const streamParam = streamId ? `&streamId=${streamId}` : '';
    return `${baseUrl}/api/ws?userId=${userId}${streamParam}`;
  }
};

export default livestreamService;

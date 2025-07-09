import { client } from './client';

// Escrow interfaces
export interface EscrowTransaction {
  id: string;
  orderId: string;
  buyerId: string;
  sellerId: string;
  amount: number;
  currency: string;
  status: 'pending' | 'held' | 'released' | 'refunded' | 'disputed' | 'cancelled';
  releaseCondition: 'manual' | 'auto' | 'timed';
  releaseAfter?: string; // ISO date string for timed release
  createdAt: string;
  updatedAt: string;
  releasedAt?: string;
  refundedAt?: string;
  disputedAt?: string;
  cancelledAt?: string;
}

export interface EscrowCreateRequest {
  orderId: string;
  sellerId: string;
  amount: number;
  currency?: string;
  releaseCondition: 'manual' | 'auto' | 'timed';
  releaseAfter?: string; // ISO date string for timed release
}

export interface EscrowUpdateRequest {
  status?: 'held' | 'released' | 'refunded' | 'disputed' | 'cancelled';
  releaseCondition?: 'manual' | 'auto' | 'timed';
  releaseAfter?: string;
}

// Dispute interfaces
export interface Dispute {
  id: string;
  escrowTransactionId: string;
  orderId: string;
  initiatorId: string;
  respondentId: string;
  reason: string;
  status: 'open' | 'under_review' | 'resolved' | 'closed';
  resolution?: 'release' | 'refund' | 'partial_release' | 'split';
  resolutionDetails?: string;
  resolutionAmount?: number;
  createdAt: string;
  updatedAt: string;
  resolvedAt?: string;
  closedAt?: string;
  evidence: DisputeEvidence[];
  messages: DisputeMessage[];
}

export interface DisputeEvidence {
  id: string;
  disputeId: string;
  userId: string;
  type: 'image' | 'document' | 'text' | 'video' | 'audio';
  content: string;
  fileUrl?: string;
  createdAt: string;
}

export interface DisputeMessage {
  id: string;
  disputeId: string;
  userId: string;
  message: string;
  isAdminMessage: boolean;
  createdAt: string;
  readAt?: string;
}

export interface DisputeCreateRequest {
  escrowTransactionId: string;
  reason: string;
  initialEvidence?: {
    type: 'image' | 'document' | 'text' | 'video' | 'audio';
    content: string;
    fileUrl?: string;
  }[];
}

export interface DisputeEvidenceCreateRequest {
  type: 'image' | 'document' | 'text' | 'video' | 'audio';
  content: string;
  fileUrl?: string;
}

export interface DisputeMessageCreateRequest {
  message: string;
}

export interface DisputeResolutionRequest {
  resolution: 'release' | 'refund' | 'partial_release' | 'split';
  resolutionDetails?: string;
  resolutionAmount?: number;
}

// API functions for escrow
export const getEscrowTransactions = async (
  page = 1,
  limit = 10,
  status?: string
): Promise<{ transactions: EscrowTransaction[]; total: number; page: number; totalPages: number }> => {
  const response = await client.get('/escrow/transactions', {
    params: { page, limit, status }
  });
  return response.data;
};

export const getEscrowTransactionById = async (id: string): Promise<EscrowTransaction> => {
  const response = await client.get(`/escrow/transactions/${id}`);
  return response.data;
};

export const getEscrowTransactionByOrderId = async (orderId: string): Promise<EscrowTransaction> => {
  const response = await client.get(`/escrow/transactions/order/${orderId}`);
  return response.data;
};

export const createEscrowTransaction = async (request: EscrowCreateRequest): Promise<EscrowTransaction> => {
  const response = await client.post('/escrow/transactions', request);
  return response.data;
};

export const updateEscrowTransaction = async (
  id: string,
  request: EscrowUpdateRequest
): Promise<EscrowTransaction> => {
  const response = await client.put(`/escrow/transactions/${id}`, request);
  return response.data;
};

export const releaseEscrowFunds = async (id: string): Promise<EscrowTransaction> => {
  const response = await client.post(`/escrow/transactions/${id}/release`);
  return response.data;
};

export const refundEscrowFunds = async (id: string): Promise<EscrowTransaction> => {
  const response = await client.post(`/escrow/transactions/${id}/refund`);
  return response.data;
};

export const cancelEscrowTransaction = async (id: string): Promise<EscrowTransaction> => {
  const response = await client.post(`/escrow/transactions/${id}/cancel`);
  return response.data;
};

// API functions for disputes
export const getDisputes = async (
  page = 1,
  limit = 10,
  status?: string
): Promise<{ disputes: Dispute[]; total: number; page: number; totalPages: number }> => {
  const response = await client.get('/escrow/disputes', {
    params: { page, limit, status }
  });
  return response.data;
};

export const getDisputeById = async (id: string): Promise<Dispute> => {
  const response = await client.get(`/escrow/disputes/${id}`);
  return response.data;
};

export const getDisputeByEscrowTransactionId = async (escrowTransactionId: string): Promise<Dispute> => {
  const response = await client.get(`/escrow/disputes/transaction/${escrowTransactionId}`);
  return response.data;
};

export const createDispute = async (request: DisputeCreateRequest): Promise<Dispute> => {
  const response = await client.post('/escrow/disputes', request);
  return response.data;
};

export const addDisputeEvidence = async (
  disputeId: string,
  evidence: DisputeEvidenceCreateRequest
): Promise<DisputeEvidence> => {
  const response = await client.post(`/escrow/disputes/${disputeId}/evidence`, evidence);
  return response.data;
};

export const addDisputeMessage = async (
  disputeId: string,
  message: DisputeMessageCreateRequest
): Promise<DisputeMessage> => {
  const response = await client.post(`/escrow/disputes/${disputeId}/messages`, message);
  return response.data;
};

export const resolveDispute = async (
  disputeId: string,
  resolution: DisputeResolutionRequest
): Promise<Dispute> => {
  const response = await client.post(`/escrow/disputes/${disputeId}/resolve`, resolution);
  return response.data;
};

export const closeDispute = async (disputeId: string): Promise<Dispute> => {
  const response = await client.post(`/escrow/disputes/${disputeId}/close`);
  return response.data;
};

export const getDisputeEvidence = async (disputeId: string): Promise<DisputeEvidence[]> => {
  const response = await client.get(`/escrow/disputes/${disputeId}/evidence`);
  return response.data;
};

export const getDisputeMessages = async (disputeId: string): Promise<DisputeMessage[]> => {
  const response = await client.get(`/escrow/disputes/${disputeId}/messages`);
  return response.data;
};

export const markDisputeMessageAsRead = async (disputeId: string, messageId: string): Promise<void> => {
  await client.post(`/escrow/disputes/${disputeId}/messages/${messageId}/read`);
};

export default {
  getEscrowTransactions,
  getEscrowTransactionById,
  getEscrowTransactionByOrderId,
  createEscrowTransaction,
  updateEscrowTransaction,
  releaseEscrowFunds,
  refundEscrowFunds,
  cancelEscrowTransaction,
  getDisputes,
  getDisputeById,
  getDisputeByEscrowTransactionId,
  createDispute,
  addDisputeEvidence,
  addDisputeMessage,
  resolveDispute,
  closeDispute,
  getDisputeEvidence,
  getDisputeMessages,
  markDisputeMessageAsRead
};

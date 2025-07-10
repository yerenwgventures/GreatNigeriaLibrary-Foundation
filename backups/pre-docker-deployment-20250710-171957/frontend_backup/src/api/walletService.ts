import { client } from './client';

// Wallet interfaces
export interface Wallet {
  id: string;
  userId: string;
  balance: number;
  currency: string;
  status: 'active' | 'suspended' | 'closed';
  createdAt: string;
  updatedAt: string;
}

export interface Transaction {
  id: string;
  userId: string;
  walletId: string;
  type: 'deposit' | 'withdrawal' | 'payment' | 'refund' | 'transfer' | 'fee' | 'other';
  amount: number;
  currency: string;
  status: 'pending' | 'completed' | 'failed' | 'cancelled';
  description: string;
  referenceId?: string;
  metadata?: Record<string, any>;
  createdAt: string;
  updatedAt: string;
}

export interface DepositRequest {
  amount: number;
  currency?: string;
  paymentMethod: 'card' | 'bank_transfer' | 'ussd' | 'mobile_money';
  returnUrl?: string;
}

export interface WithdrawalRequest {
  amount: number;
  currency?: string;
  withdrawalMethod: 'bank_transfer' | 'mobile_money';
  accountDetails: {
    accountName?: string;
    accountNumber?: string;
    bankCode?: string;
    bankName?: string;
    phoneNumber?: string;
    provider?: string;
  };
}

export interface TransferRequest {
  amount: number;
  recipientUserId: string;
  description?: string;
}

export interface PaymentRequest {
  amount: number;
  itemId: string;
  itemType: 'product' | 'service' | 'subscription' | 'other';
  description?: string;
}

// Payment method interfaces
export interface PaymentMethod {
  id: string;
  userId: string;
  type: 'card' | 'bank_account' | 'mobile_money';
  isDefault: boolean;
  details: {
    last4?: string;
    brand?: string;
    expiryMonth?: string;
    expiryYear?: string;
    accountName?: string;
    bankName?: string;
    phoneNumber?: string;
    provider?: string;
  };
  createdAt: string;
  updatedAt: string;
}

export interface AddPaymentMethodRequest {
  type: 'card' | 'bank_account' | 'mobile_money';
  token?: string; // For card tokenization
  details?: {
    accountName?: string;
    accountNumber?: string;
    bankCode?: string;
    phoneNumber?: string;
    provider?: string;
  };
  setAsDefault?: boolean;
}

// API functions for wallet
export const getWallet = async (): Promise<Wallet> => {
  const response = await client.get('/wallet');
  return response.data;
};

export const getTransactions = async (
  page = 1,
  limit = 20,
  type?: string,
  status?: string,
  startDate?: string,
  endDate?: string
): Promise<{ transactions: Transaction[]; total: number; page: number; totalPages: number }> => {
  const response = await client.get('/wallet/transactions', {
    params: { page, limit, type, status, startDate, endDate }
  });
  return response.data;
};

export const initiateDeposit = async (request: DepositRequest): Promise<{ redirectUrl: string; reference: string }> => {
  const response = await client.post('/wallet/deposit', request);
  return response.data;
};

export const requestWithdrawal = async (request: WithdrawalRequest): Promise<{ reference: string; status: string }> => {
  const response = await client.post('/wallet/withdraw', request);
  return response.data;
};

export const transferFunds = async (request: TransferRequest): Promise<Transaction> => {
  const response = await client.post('/wallet/transfer', request);
  return response.data;
};

export const makePayment = async (request: PaymentRequest): Promise<Transaction> => {
  const response = await client.post('/wallet/pay', request);
  return response.data;
};

// API functions for payment methods
export const getPaymentMethods = async (): Promise<PaymentMethod[]> => {
  const response = await client.get('/wallet/payment-methods');
  return response.data;
};

export const addPaymentMethod = async (request: AddPaymentMethodRequest): Promise<PaymentMethod> => {
  const response = await client.post('/wallet/payment-methods', request);
  return response.data;
};

export const deletePaymentMethod = async (id: string): Promise<void> => {
  await client.delete(`/wallet/payment-methods/${id}`);
};

export const setDefaultPaymentMethod = async (id: string): Promise<PaymentMethod> => {
  const response = await client.put(`/wallet/payment-methods/${id}/default`);
  return response.data;
};

export default {
  getWallet,
  getTransactions,
  initiateDeposit,
  requestWithdrawal,
  transferFunds,
  makePayment,
  getPaymentMethods,
  addPaymentMethod,
  deletePaymentMethod,
  setDefaultPaymentMethod
};

import { client } from './client';

// Product interfaces
export interface Product {
  id: string;
  userId: string;
  title: string;
  description: string;
  price: number;
  currency: string;
  category: string;
  subcategory?: string;
  status: 'active' | 'pending' | 'sold' | 'inactive';
  location?: string;
  mediaUrls: string[];
  tags?: string[];
  createdAt: string;
  updatedAt: string;
  seller?: User;
  rating?: number;
  reviewCount?: number;
}

export interface ProductCreateRequest {
  title: string;
  description: string;
  price: number;
  currency?: string;
  category: string;
  subcategory?: string;
  location?: string;
  mediaUrls?: string[];
  tags?: string[];
}

export interface ProductUpdateRequest {
  title?: string;
  description?: string;
  price?: number;
  currency?: string;
  category?: string;
  subcategory?: string;
  status?: 'active' | 'pending' | 'sold' | 'inactive';
  location?: string;
  mediaUrls?: string[];
  tags?: string[];
}

// Service interfaces
export interface Service {
  id: string;
  userId: string;
  title: string;
  description: string;
  priceType: 'fixed' | 'hourly' | 'daily' | 'custom';
  price: number;
  currency: string;
  category: string;
  subcategory?: string;
  status: 'active' | 'pending' | 'inactive';
  location?: string;
  isRemote: boolean;
  mediaUrls: string[];
  tags?: string[];
  createdAt: string;
  updatedAt: string;
  provider?: User;
  rating?: number;
  reviewCount?: number;
}

export interface ServiceCreateRequest {
  title: string;
  description: string;
  priceType: 'fixed' | 'hourly' | 'daily' | 'custom';
  price: number;
  currency?: string;
  category: string;
  subcategory?: string;
  location?: string;
  isRemote: boolean;
  mediaUrls?: string[];
  tags?: string[];
}

export interface ServiceUpdateRequest {
  title?: string;
  description?: string;
  priceType?: 'fixed' | 'hourly' | 'daily' | 'custom';
  price?: number;
  currency?: string;
  category?: string;
  subcategory?: string;
  status?: 'active' | 'pending' | 'inactive';
  location?: string;
  isRemote?: boolean;
  mediaUrls?: string[];
  tags?: string[];
}

// Job interfaces
export interface Job {
  id: string;
  userId: string;
  title: string;
  description: string;
  company: string;
  locationType: 'remote' | 'onsite' | 'hybrid';
  location?: string;
  salaryMin?: number;
  salaryMax?: number;
  salaryCurrency?: string;
  salaryPeriod?: 'hourly' | 'daily' | 'weekly' | 'monthly' | 'yearly';
  category: string;
  status: 'active' | 'filled' | 'expired' | 'inactive';
  applicationUrl?: string;
  applicationEmail?: string;
  applicationDeadline?: string;
  tags?: string[];
  createdAt: string;
  updatedAt: string;
  poster?: User;
}

export interface JobCreateRequest {
  title: string;
  description: string;
  company: string;
  locationType: 'remote' | 'onsite' | 'hybrid';
  location?: string;
  salaryMin?: number;
  salaryMax?: number;
  salaryCurrency?: string;
  salaryPeriod?: 'hourly' | 'daily' | 'weekly' | 'monthly' | 'yearly';
  category: string;
  applicationUrl?: string;
  applicationEmail?: string;
  applicationDeadline?: string;
  tags?: string[];
}

export interface JobUpdateRequest {
  title?: string;
  description?: string;
  company?: string;
  locationType?: 'remote' | 'onsite' | 'hybrid';
  location?: string;
  salaryMin?: number;
  salaryMax?: number;
  salaryCurrency?: string;
  salaryPeriod?: 'hourly' | 'daily' | 'weekly' | 'monthly' | 'yearly';
  category?: string;
  status?: 'active' | 'filled' | 'expired' | 'inactive';
  applicationUrl?: string;
  applicationEmail?: string;
  applicationDeadline?: string;
  tags?: string[];
}

// Review interfaces
export interface Review {
  id: string;
  userId: string;
  itemId: string;
  itemType: 'product' | 'service' | 'job';
  rating: number;
  comment: string;
  createdAt: string;
  updatedAt: string;
  reviewer?: User;
}

export interface ReviewCreateRequest {
  itemId: string;
  itemType: 'product' | 'service' | 'job';
  rating: number;
  comment: string;
}

// Category interfaces
export interface Category {
  id: string;
  name: string;
  description?: string;
  parentId?: string;
  icon?: string;
  displayOrder?: number;
  subcategories?: Category[];
}

// User interface (simplified)
interface User {
  id: string;
  username: string;
  name: string;
  profileImage?: string;
  rating?: number;
  reviewCount?: number;
}

// Search interfaces
export interface SearchParams {
  query?: string;
  category?: string;
  subcategory?: string;
  minPrice?: number;
  maxPrice?: number;
  location?: string;
  tags?: string[];
  status?: string;
  sortBy?: 'newest' | 'oldest' | 'price_low' | 'price_high' | 'rating';
  page?: number;
  limit?: number;
}

export interface SearchResult<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// API functions for products
export const getProducts = async (params?: SearchParams): Promise<SearchResult<Product>> => {
  const response = await client.get('/marketplace/products', { params });
  return response.data;
};

export const getProductById = async (id: string): Promise<Product> => {
  const response = await client.get(`/marketplace/products/${id}`);
  return response.data;
};

export const createProduct = async (product: ProductCreateRequest): Promise<Product> => {
  const response = await client.post('/marketplace/products', product);
  return response.data;
};

export const updateProduct = async (id: string, product: ProductUpdateRequest): Promise<Product> => {
  const response = await client.put(`/marketplace/products/${id}`, product);
  return response.data;
};

export const deleteProduct = async (id: string): Promise<void> => {
  await client.delete(`/marketplace/products/${id}`);
};

// API functions for services
export const getServices = async (params?: SearchParams): Promise<SearchResult<Service>> => {
  const response = await client.get('/marketplace/services', { params });
  return response.data;
};

export const getServiceById = async (id: string): Promise<Service> => {
  const response = await client.get(`/marketplace/services/${id}`);
  return response.data;
};

export const createService = async (service: ServiceCreateRequest): Promise<Service> => {
  const response = await client.post('/marketplace/services', service);
  return response.data;
};

export const updateService = async (id: string, service: ServiceUpdateRequest): Promise<Service> => {
  const response = await client.put(`/marketplace/services/${id}`, service);
  return response.data;
};

export const deleteService = async (id: string): Promise<void> => {
  await client.delete(`/marketplace/services/${id}`);
};

// API functions for jobs
export const getJobs = async (params?: SearchParams): Promise<SearchResult<Job>> => {
  const response = await client.get('/marketplace/jobs', { params });
  return response.data;
};

export const getJobById = async (id: string): Promise<Job> => {
  const response = await client.get(`/marketplace/jobs/${id}`);
  return response.data;
};

export const createJob = async (job: JobCreateRequest): Promise<Job> => {
  const response = await client.post('/marketplace/jobs', job);
  return response.data;
};

export const updateJob = async (id: string, job: JobUpdateRequest): Promise<Job> => {
  const response = await client.put(`/marketplace/jobs/${id}`, job);
  return response.data;
};

export const deleteJob = async (id: string): Promise<void> => {
  await client.delete(`/marketplace/jobs/${id}`);
};

// API functions for reviews
export const getReviews = async (itemId: string, itemType: 'product' | 'service' | 'job'): Promise<Review[]> => {
  const response = await client.get(`/marketplace/reviews`, { params: { itemId, itemType } });
  return response.data;
};

export const createReview = async (review: ReviewCreateRequest): Promise<Review> => {
  const response = await client.post('/marketplace/reviews', review);
  return response.data;
};

// API functions for categories
export const getCategories = async (): Promise<Category[]> => {
  const response = await client.get('/marketplace/categories');
  return response.data;
};

export default {
  getProducts,
  getProductById,
  createProduct,
  updateProduct,
  deleteProduct,
  getServices,
  getServiceById,
  createService,
  updateService,
  deleteService,
  getJobs,
  getJobById,
  createJob,
  updateJob,
  deleteJob,
  getReviews,
  createReview,
  getCategories
};

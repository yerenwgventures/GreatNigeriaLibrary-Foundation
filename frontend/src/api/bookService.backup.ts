import apiClient from './client';
import { Book, Bookmark, Chapter, ReadingProgress, Section, Subsection } from '../types';

const BookService = {
  /**
   * Get all books
   */
  getBooks: async (): Promise<Book[]> => {
    const response = await apiClient.get<Book[]>('/books');
    return response.data;
  },

  /**
   * Get book by ID
   */
  getBookById: async (id: string): Promise<Book> => {
    const response = await apiClient.get<Book>(`/books/${id}`);
    return response.data;
  },

  /**
   * Get chapters for a book
   */
  getBookChapters: async (bookId: string): Promise<Chapter[]> => {
    const response = await apiClient.get<Chapter[]>(`/books/${bookId}/chapters`);
    return response.data;
  },

  /**
   * Get chapter by ID
   */
  getChapterById: async (chapterId: string): Promise<Chapter> => {
    const response = await apiClient.get<Chapter>(`/books/chapters/${chapterId}`);
    return response.data;
  },

  /**
   * Get section by ID
   */
  getSectionById: async (sectionId: string): Promise<Section> => {
    const response = await apiClient.get<Section>(`/books/sections/${sectionId}`);
    return response.data;
  },

  /**
   * Save reading progress
   */
  saveReadingProgress: async (bookId: string, sectionId: string): Promise<void> => {
    await apiClient.post(`/books/${bookId}/progress`, { sectionId });
  },

  /**
   * Get reading progress for a book
   */
  getReadingProgress: async (bookId: string): Promise<ReadingProgress> => {
    const response = await apiClient.get<ReadingProgress>(`/books/${bookId}/progress`);
    return response.data;
  },

  /**
   * Add bookmark
   */
  addBookmark: async (bookId: string, sectionId: string, note?: string): Promise<Bookmark> => {
    const response = await apiClient.post<Bookmark>(`/books/${bookId}/bookmarks`, {
      sectionId,
      note,
    });
    return response.data;
  },

  /**
   * Get bookmarks for a book
   */
  getBookmarks: async (bookId: string): Promise<Bookmark[]> => {
    const response = await apiClient.get<Bookmark[]>(`/books/${bookId}/bookmarks`);
    return response.data;
  },

  /**
   * Delete bookmark
   */
  deleteBookmark: async (bookmarkId: string): Promise<void> => {
    await apiClient.delete(`/books/bookmarks/${bookmarkId}`);
  },

  /**
   * Get subsections for a section
   */
  getSubsectionsBySection: async (sectionId: string): Promise<Subsection[]> => {
    const response = await apiClient.get<Subsection[]>(`/books/sections/${sectionId}/subsections`);
    return response.data;
  },

  /**
   * Get subsection by ID
   */
  getSubsectionById: async (subsectionId: string): Promise<Subsection> => {
    const response = await apiClient.get<Subsection>(`/books/subsections/${subsectionId}`);
    return response.data;
  },
};

export default BookService;

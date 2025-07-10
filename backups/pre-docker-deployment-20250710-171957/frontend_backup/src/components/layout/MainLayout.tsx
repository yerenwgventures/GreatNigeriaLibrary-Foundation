import React from 'react';
import { Box } from '@mui/material';
import Header from './Header';
import Footer from './Footer';
import ContextualTipsComponent from '../tips/ContextualTipsComponent';
import { useLocation } from 'react-router-dom';

interface MainLayoutProps {
  children: React.ReactNode;
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const location = useLocation();
  
  // Determine context type and ID based on the current route
  const getContextInfo = () => {
    const path = location.pathname;
    
    // Book viewer page
    if (path.startsWith('/books/') && path.split('/').length > 2) {
      const bookId = path.split('/')[2];
      return { type: 'book', id: bookId };
    }
    
    // Forum topic page
    if (path.startsWith('/forum/topic/') && path.split('/').length > 3) {
      const topicId = path.split('/')[3];
      return { type: 'forum_topic', id: topicId };
    }
    
    // Marketplace page
    if (path.startsWith('/marketplace')) {
      return { type: 'marketplace', id: 'main' };
    }
    
    // Profile page
    if (path.startsWith('/profile')) {
      return { type: 'profile', id: 'main' };
    }
    
    // Default to page context
    return { type: 'page', id: path };
  };
  
  const { type, id } = getContextInfo();
  
  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <Header />
      <Box component="main" sx={{ flexGrow: 1, py: 3 }}>
        {children}
      </Box>
      <Footer />
      
      {/* Contextual Tips Component */}
      <ContextualTipsComponent 
        contextType={type} 
        contextId={id} 
        position="bottom" 
      />
    </Box>
  );
};

export default MainLayout;

import React from 'react';
import { Outlet } from 'react-router-dom';
import { Box, Container, useTheme, useMediaQuery } from '@mui/material';
import ResponsiveHeader from './ResponsiveHeader';
import ResponsiveFooter from './ResponsiveFooter';
import ContextualTipsComponent from '../tips/ContextualTipsComponent';
import { useLocation } from 'react-router-dom';

const ResponsiveLayout: React.FC = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
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
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        minHeight: '100vh',
        backgroundColor: theme.palette.background.default,
      }}
    >
      <ResponsiveHeader />
      
      <Box
        component="main"
        sx={{
          flexGrow: 1,
          py: { xs: 2, md: 4 },
          px: { xs: 1, sm: 2, md: 3 },
        }}
      >
        <Outlet />
      </Box>
      
      {/* Contextual Tips Component */}
      <ContextualTipsComponent 
        contextType={type} 
        contextId={id} 
        position={isMobile ? "bottom" : "right"} 
      />
      
      <ResponsiveFooter />
    </Box>
  );
};

export default ResponsiveLayout;

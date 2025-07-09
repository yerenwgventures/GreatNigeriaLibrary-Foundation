import React from 'react';
import {
  Box,
  Typography,
  Chip,
  Avatar,
  Grid,
  Divider
} from '@mui/material';
import {
  Person as PersonIcon,
  Category as CategoryIcon,
  LocalOffer as TagIcon,
  Visibility as ViewIcon,
  CardGiftcard as GiftIcon,
  Public as PublicIcon,
  Lock as PrivateIcon
} from '@mui/icons-material';
import { Stream } from '../../api/livestreamService';

interface StreamInfoPanelProps {
  stream: Stream;
}

const StreamInfoPanel: React.FC<StreamInfoPanelProps> = ({ stream }) => {
  const tags = stream.tags ? stream.tags.split(',').filter(Boolean) : [];
  const categories = stream.categories ? stream.categories.split(',').filter(Boolean) : [];
  
  return (
    <Box>
      <Grid container spacing={2}>
        <Grid item xs={12} md={8}>
          <Typography variant="body1" paragraph>
            {stream.description || 'No description provided.'}
          </Typography>
          
          <Box display="flex" alignItems="center" mb={1}>
            <Avatar sx={{ width: 24, height: 24, mr: 1 }}>
              <PersonIcon fontSize="small" />
            </Avatar>
            <Typography variant="body2">
              Creator #{stream.creatorId}
            </Typography>
          </Box>
          
          {categories.length > 0 && (
            <Box display="flex" alignItems="center" mb={1}>
              <CategoryIcon fontSize="small" sx={{ mr: 1 }} />
              <Box display="flex" flexWrap="wrap" gap={0.5}>
                {categories.map((category, index) => (
                  <Chip 
                    key={index} 
                    label={category.trim()} 
                    size="small" 
                    color="primary" 
                    variant="outlined" 
                  />
                ))}
              </Box>
            </Box>
          )}
          
          {tags.length > 0 && (
            <Box display="flex" alignItems="center">
              <TagIcon fontSize="small" sx={{ mr: 1 }} />
              <Box display="flex" flexWrap="wrap" gap={0.5}>
                {tags.map((tag, index) => (
                  <Chip 
                    key={index} 
                    label={tag.trim()} 
                    size="small" 
                    variant="outlined" 
                  />
                ))}
              </Box>
            </Box>
          )}
        </Grid>
        
        <Grid item xs={12} md={4}>
          <Box>
            <Typography variant="subtitle2" gutterBottom>
              Stream Stats
            </Typography>
            
            <Box display="flex" alignItems="center" mb={1}>
              <ViewIcon fontSize="small" sx={{ mr: 1 }} />
              <Typography variant="body2">
                {stream.viewerCount} current viewers (peak: {stream.peakViewerCount})
              </Typography>
            </Box>
            
            <Box display="flex" alignItems="center" mb={1}>
              <GiftIcon fontSize="small" sx={{ mr: 1 }} />
              <Typography variant="body2">
                ₦{stream.totalGiftsValue.toLocaleString()} in gifts
              </Typography>
            </Box>
            
            <Box display="flex" alignItems="center">
              {stream.isPrivate ? (
                <>
                  <PrivateIcon fontSize="small" sx={{ mr: 1 }} />
                  <Typography variant="body2">
                    Private stream
                  </Typography>
                </>
              ) : (
                <>
                  <PublicIcon fontSize="small" sx={{ mr: 1 }} />
                  <Typography variant="body2">
                    Public stream
                  </Typography>
                </>
              )}
            </Box>
          </Box>
        </Grid>
      </Grid>
      
      <Divider sx={{ my: 2 }} />
      
      <Box>
        <Typography variant="subtitle2" gutterBottom>
          Stream Timeline
        </Typography>
        
        <Box display="flex" flexDirection="column" gap={1}>
          <Box display="flex" alignItems="center">
            <Chip 
              label="Created" 
              size="small" 
              color="default" 
              sx={{ mr: 1, minWidth: 80 }} 
            />
            <Typography variant="body2">
              {new Date(stream.createdAt).toLocaleString()}
            </Typography>
          </Box>
          
          <Box display="flex" alignItems="center">
            <Chip 
              label="Scheduled" 
              size="small" 
              color="info" 
              sx={{ mr: 1, minWidth: 80 }} 
            />
            <Typography variant="body2">
              {new Date(stream.scheduledStart).toLocaleString()}
            </Typography>
          </Box>
          
          {stream.actualStart && (
            <Box display="flex" alignItems="center">
              <Chip 
                label="Started" 
                size="small" 
                color="success" 
                sx={{ mr: 1, minWidth: 80 }} 
              />
              <Typography variant="body2">
                {new Date(stream.actualStart).toLocaleString()}
              </Typography>
            </Box>
          )}
          
          {stream.endTime && (
            <Box display="flex" alignItems="center">
              <Chip 
                label="Ended" 
                size="small" 
                color="error" 
                sx={{ mr: 1, minWidth: 80 }} 
              />
              <Typography variant="body2">
                {new Date(stream.endTime).toLocaleString()}
              </Typography>
            </Box>
          )}
        </Box>
      </Box>
    </Box>
  );
};

export default StreamInfoPanel;

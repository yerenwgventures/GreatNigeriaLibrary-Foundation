import React from 'react';
import { 
  Card, 
  CardActionArea, 
  CardContent, 
  CardMedia, 
  Typography, 
  Box, 
  Chip, 
  Avatar 
} from '@mui/material';
import { 
  Person as PersonIcon, 
  Visibility as ViewIcon 
} from '@mui/icons-material';
import { Stream } from '../../api/livestreamService';

interface StreamCardProps {
  stream: Stream;
  onClick: () => void;
}

const StreamCard: React.FC<StreamCardProps> = ({ stream, onClick }) => {
  const isLive = stream.status === 'live';
  const placeholderImage = 'https://via.placeholder.com/320x180?text=No+Thumbnail';
  
  return (
    <Card 
      elevation={3} 
      sx={{ 
        height: '100%', 
        display: 'flex', 
        flexDirection: 'column',
        transition: 'transform 0.2s',
        '&:hover': {
          transform: 'translateY(-5px)'
        }
      }}
    >
      <CardActionArea onClick={onClick} sx={{ flexGrow: 1 }}>
        <Box sx={{ position: 'relative' }}>
          <CardMedia
            component="img"
            height="180"
            image={stream.thumbnailUrl || placeholderImage}
            alt={stream.title}
          />
          
          {isLive && (
            <Chip
              label="LIVE"
              color="error"
              size="small"
              sx={{
                position: 'absolute',
                top: 10,
                left: 10,
                fontWeight: 'bold'
              }}
            />
          )}
          
          {stream.viewerCount > 0 && (
            <Box
              sx={{
                position: 'absolute',
                bottom: 10,
                right: 10,
                bgcolor: 'rgba(0, 0, 0, 0.6)',
                color: 'white',
                px: 1,
                py: 0.5,
                borderRadius: 1,
                display: 'flex',
                alignItems: 'center'
              }}
            >
              <ViewIcon fontSize="small" sx={{ mr: 0.5 }} />
              <Typography variant="caption">
                {stream.viewerCount}
              </Typography>
            </Box>
          )}
        </Box>
        
        <CardContent>
          <Typography variant="h6" component="div" noWrap>
            {stream.title}
          </Typography>
          
          <Box display="flex" alignItems="center" mt={1}>
            <Avatar sx={{ width: 24, height: 24, mr: 1 }}>
              <PersonIcon fontSize="small" />
            </Avatar>
            <Typography variant="body2" color="text.secondary">
              Creator #{stream.creatorId}
            </Typography>
          </Box>
          
          <Box display="flex" justifyContent="space-between" alignItems="center" mt={1}>
            <Typography variant="caption" color="text.secondary">
              {isLive 
                ? `Started ${new Date(stream.actualStart!).toLocaleString()}`
                : stream.status === 'scheduled'
                  ? `Scheduled for ${new Date(stream.scheduledStart).toLocaleString()}`
                  : `Ended ${new Date(stream.endTime!).toLocaleString()}`
              }
            </Typography>
            
            {stream.totalGiftsValue > 0 && (
              <Chip
                label={`₦${stream.totalGiftsValue.toLocaleString()}`}
                size="small"
                color="primary"
                variant="outlined"
              />
            )}
          </Box>
        </CardContent>
      </CardActionArea>
    </Card>
  );
};

export default StreamCard;

import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Box,
  Typography,
  Avatar,
  Button,
  Chip,
  Divider,
  Grid,
  Paper,
  Tabs,
  Tab,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  CircularProgress
} from '@mui/material';
import {
  Person as PersonIcon,
  Videocam as VideoIcon,
  CardGiftcard as GiftIcon,
  EmojiEvents as TrophyIcon,
  MonetizationOn as RevenueIcon
} from '@mui/icons-material';
import { RootState } from '../../store';
import { fetchUserSentGifts, fetchUserReceivedGifts, fetchUserRanking } from '../../features/livestream/livestreamSlice';
import { useNavigate } from 'react-router-dom';

interface StreamerProfileProps {
  userId: number;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`profile-tabpanel-${index}`}
      aria-labelledby={`profile-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 2 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

const StreamerProfile: React.FC<StreamerProfileProps> = ({ userId }) => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { user } = useSelector((state: RootState) => state.auth);
  const { 
    sent: sentGifts, 
    received: receivedGifts, 
    loading: giftsLoading 
  } = useSelector((state: RootState) => state.livestream.gifts);
  const { 
    userRanking, 
    loading: rankingLoading 
  } = useSelector((state: RootState) => state.livestream.rankings);
  
  const [tabValue, setTabValue] = useState(0);
  const [isFollowing, setIsFollowing] = useState(false);
  
  useEffect(() => {
    dispatch(fetchUserSentGifts({ userId, page: 1, limit: 5 }) as any);
    dispatch(fetchUserReceivedGifts({ userId, page: 1, limit: 5 }) as any);
    dispatch(fetchUserRanking({ userId, period: 'all_time' }) as any);
  }, [dispatch, userId]);
  
  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };
  
  const handleFollowToggle = () => {
    // In a real app, this would call an API to follow/unfollow
    setIsFollowing(!isFollowing);
  };
  
  const handleViewStreams = () => {
    // Navigate to a filtered view of this user's streams
    navigate(`/livestream?creator=${userId}`);
  };
  
  const isCurrentUser = user && user.id === userId;
  
  // Badge colors based on level
  const getBadgeColor = (level: string) => {
    switch (level) {
      case 'diamond':
        return '#B9F2FF';
      case 'platinum':
        return '#E5E4E2';
      case 'gold':
        return '#FFD700';
      case 'silver':
        return '#C0C0C0';
      case 'bronze':
      default:
        return '#CD7F32';
    }
  };
  
  return (
    <Box>
      <Paper elevation={2} sx={{ p: 3, mb: 3 }}>
        <Grid container spacing={3}>
          <Grid item xs={12} md={3}>
            <Box display="flex" flexDirection="column" alignItems="center">
              <Avatar
                sx={{ width: 120, height: 120, mb: 2 }}
              >
                <PersonIcon fontSize="large" />
              </Avatar>
              
              {!isCurrentUser && (
                <Button
                  variant={isFollowing ? "outlined" : "contained"}
                  color={isFollowing ? "secondary" : "primary"}
                  onClick={handleFollowToggle}
                  fullWidth
                  sx={{ mb: 1 }}
                >
                  {isFollowing ? "Unfollow" : "Follow"}
                </Button>
              )}
              
              <Button
                variant="outlined"
                startIcon={<VideoIcon />}
                onClick={handleViewStreams}
                fullWidth
              >
                View Streams
              </Button>
            </Box>
          </Grid>
          
          <Grid item xs={12} md={9}>
            <Box>
              <Typography variant="h4" gutterBottom>
                User #{userId}
                {isCurrentUser && (
                  <Chip 
                    label="You" 
                    color="primary" 
                    size="small" 
                    sx={{ ml: 1 }} 
                  />
                )}
              </Typography>
              
              <Typography variant="body1" paragraph>
                {isCurrentUser ? 'This is your profile' : `This is the profile of User #${userId}`}
              </Typography>
              
              <Box display="flex" flexWrap="wrap" gap={2} mb={2}>
                {rankingLoading ? (
                  <CircularProgress size={24} />
                ) : userRanking ? (
                  <Chip
                    icon={<TrophyIcon />}
                    label={`${userRanking.badgeLevel.toUpperCase()} Gifter`}
                    sx={{
                      bgcolor: getBadgeColor(userRanking.badgeLevel),
                      color: 'text.primary',
                      fontWeight: 'bold'
                    }}
                  />
                ) : null}
                
                <Chip
                  icon={<VideoIcon />}
                  label="Streamer"
                  color="primary"
                  variant="outlined"
                />
                
                {isCurrentUser && (
                  <Chip
                    icon={<RevenueIcon />}
                    label="Creator"
                    color="secondary"
                    variant="outlined"
                  />
                )}
              </Box>
              
              <Grid container spacing={2}>
                <Grid item xs={4}>
                  <Box textAlign="center">
                    <Typography variant="h6">
                      {userRanking?.totalGifts || 0}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Gifts Sent
                    </Typography>
                  </Box>
                </Grid>
                
                <Grid item xs={4}>
                  <Box textAlign="center">
                    <Typography variant="h6">
                      {receivedGifts.length || 0}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Gifts Received
                    </Typography>
                  </Box>
                </Grid>
                
                <Grid item xs={4}>
                  <Box textAlign="center">
                    <Typography variant="h6">
                      {userRanking?.rank || '-'}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      Gifter Rank
                    </Typography>
                  </Box>
                </Grid>
              </Grid>
            </Box>
          </Grid>
        </Grid>
      </Paper>
      
      <Paper elevation={2}>
        <Tabs
          value={tabValue}
          onChange={handleTabChange}
          variant="fullWidth"
        >
          <Tab 
            icon={<GiftIcon />} 
            label="Gifts Received" 
            id="profile-tab-0" 
            aria-controls="profile-tabpanel-0" 
          />
          <Tab 
            icon={<GiftIcon />} 
            label="Gifts Sent" 
            id="profile-tab-1" 
            aria-controls="profile-tabpanel-1" 
          />
        </Tabs>
        
        <TabPanel value={tabValue} index={0}>
          <Typography variant="h6" gutterBottom>
            Recent Gifts Received
          </Typography>
          
          {giftsLoading ? (
            <Box display="flex" justifyContent="center" my={4}>
              <CircularProgress />
            </Box>
          ) : receivedGifts.length > 0 ? (
            <List>
              {receivedGifts.map((gift) => (
                <React.Fragment key={gift.id}>
                  <ListItem alignItems="flex-start">
                    <ListItemAvatar>
                      <Avatar>
                        <PersonIcon />
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText
                      primary={
                        <Typography variant="body1">
                          {gift.isAnonymous ? 'Anonymous' : `User #${gift.senderId}`} sent a {gift.giftName}
                        </Typography>
                      }
                      secondary={
                        <>
                          <Typography component="span" variant="body2" color="text.primary">
                            {gift.coinsAmount} coins
                          </Typography>
                          {gift.message && (
                            <Typography component="p" variant="body2">
                              "{gift.message}"
                            </Typography>
                          )}
                          <Typography component="span" variant="caption" color="text.secondary">
                            {new Date(gift.createdAt).toLocaleString()}
                          </Typography>
                        </>
                      }
                    />
                  </ListItem>
                  <Divider variant="inset" component="li" />
                </React.Fragment>
              ))}
            </List>
          ) : (
            <Typography variant="body1" color="text.secondary" align="center" sx={{ my: 4 }}>
              No gifts received yet
            </Typography>
          )}
          
          {receivedGifts.length > 0 && (
            <Box display="flex" justifyContent="center" mt={2}>
              <Button variant="text" color="primary">
                View All
              </Button>
            </Box>
          )}
        </TabPanel>
        
        <TabPanel value={tabValue} index={1}>
          <Typography variant="h6" gutterBottom>
            Recent Gifts Sent
          </Typography>
          
          {giftsLoading ? (
            <Box display="flex" justifyContent="center" my={4}>
              <CircularProgress />
            </Box>
          ) : sentGifts.length > 0 ? (
            <List>
              {sentGifts.map((gift) => (
                <React.Fragment key={gift.id}>
                  <ListItem alignItems="flex-start">
                    <ListItemAvatar>
                      <Avatar>
                        <PersonIcon />
                      </Avatar>
                    </ListItemAvatar>
                    <ListItemText
                      primary={
                        <Typography variant="body1">
                          Sent a {gift.giftName} to User #{gift.recipientId}
                        </Typography>
                      }
                      secondary={
                        <>
                          <Typography component="span" variant="body2" color="text.primary">
                            {gift.coinsAmount} coins
                          </Typography>
                          {gift.message && (
                            <Typography component="p" variant="body2">
                              "{gift.message}"
                            </Typography>
                          )}
                          <Typography component="span" variant="caption" color="text.secondary">
                            {new Date(gift.createdAt).toLocaleString()}
                          </Typography>
                        </>
                      }
                    />
                  </ListItem>
                  <Divider variant="inset" component="li" />
                </React.Fragment>
              ))}
            </List>
          ) : (
            <Typography variant="body1" color="text.secondary" align="center" sx={{ my: 4 }}>
              No gifts sent yet
            </Typography>
          )}
          
          {sentGifts.length > 0 && (
            <Box display="flex" justifyContent="center" mt={2}>
              <Button variant="text" color="primary">
                View All
              </Button>
            </Box>
          )}
        </TabPanel>
      </Paper>
    </Box>
  );
};

export default StreamerProfile;

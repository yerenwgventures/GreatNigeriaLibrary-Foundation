import React, { useState } from 'react';
import { useDispatch } from 'react-redux';
import {
  Box,
  Typography,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Avatar,
  Divider,
  Tabs,
  Tab,
  Chip,
  Badge,
  Paper
} from '@mui/material';
import {
  EmojiEvents as TrophyIcon,
  ArrowUpward as UpIcon,
  ArrowDownward as DownIcon,
  Remove as SameIcon,
  Person as PersonIcon
} from '@mui/icons-material';
import { GifterRanking } from '../../api/livestreamService';
import { fetchStreamRankings } from '../../features/livestream/livestreamSlice';

interface LeaderboardPanelProps {
  rankings: GifterRanking[];
  streamId: number;
}

const LeaderboardPanel: React.FC<LeaderboardPanelProps> = ({ rankings, streamId }) => {
  const dispatch = useDispatch();
  const [period, setPeriod] = useState('daily');
  
  const handlePeriodChange = (_event: React.SyntheticEvent, newPeriod: string) => {
    setPeriod(newPeriod);
    dispatch(fetchStreamRankings({ streamId, period: newPeriod, limit: 10 }) as any);
  };
  
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
  
  // Rank change icon
  const getRankChangeIcon = (current: number, previous: number | undefined) => {
    if (previous === undefined) return null;
    
    if (current < previous) {
      return <UpIcon fontSize="small" color="success" />;
    } else if (current > previous) {
      return <DownIcon fontSize="small" color="error" />;
    } else {
      return <SameIcon fontSize="small" color="disabled" />;
    }
  };
  
  return (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Tabs
        value={period}
        onChange={handlePeriodChange}
        variant="fullWidth"
        sx={{ mb: 2 }}
      >
        <Tab label="Daily" value="daily" />
        <Tab label="Weekly" value="weekly" />
        <Tab label="Monthly" value="monthly" />
        <Tab label="All Time" value="all_time" />
      </Tabs>
      
      {rankings.length > 0 ? (
        <List disablePadding sx={{ flexGrow: 1, overflow: 'auto' }}>
          {rankings.map((ranking, index) => (
            <React.Fragment key={ranking.id}>
              <ListItem alignItems="flex-start">
                <ListItemAvatar>
                  <Badge
                    overlap="circular"
                    anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                    badgeContent={
                      <Avatar
                        sx={{
                          width: 22,
                          height: 22,
                          bgcolor: index < 3 ? ['gold', 'silver', '#cd7f32'][index] : 'grey.500',
                          border: '2px solid white'
                        }}
                      >
                        {ranking.rank}
                      </Avatar>
                    }
                  >
                    <Avatar sx={{ bgcolor: getBadgeColor(ranking.badgeLevel) }}>
                      <PersonIcon />
                    </Avatar>
                  </Badge>
                </ListItemAvatar>
                <ListItemText
                  primary={
                    <Box display="flex" alignItems="center">
                      <Typography variant="body1" fontWeight={index < 3 ? 'bold' : 'normal'}>
                        User #{ranking.userId}
                      </Typography>
                      <Box ml={1} display="flex" alignItems="center">
                        {getRankChangeIcon(ranking.rank, ranking.previousRank)}
                      </Box>
                    </Box>
                  }
                  secondary={
                    <>
                      <Box display="flex" alignItems="center" mt={0.5}>
                        <Chip
                          label={ranking.badgeLevel.toUpperCase()}
                          size="small"
                          sx={{
                            bgcolor: getBadgeColor(ranking.badgeLevel),
                            color: 'text.primary',
                            mr: 1
                          }}
                        />
                        <Typography variant="body2" color="text.secondary">
                          {ranking.totalGifts} gifts
                        </Typography>
                      </Box>
                      <Typography variant="body2" color="primary" fontWeight="bold">
                        {ranking.totalCoins} coins (₦{ranking.totalNairaValue.toLocaleString()})
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
        <Paper 
          variant="outlined" 
          sx={{ 
            flexGrow: 1, 
            display: 'flex', 
            flexDirection: 'column', 
            alignItems: 'center', 
            justifyContent: 'center',
            p: 3
          }}
        >
          <TrophyIcon fontSize="large" color="disabled" />
          <Typography variant="body1" color="text.secondary" align="center" mt={2}>
            No rankings available yet
          </Typography>
          <Typography variant="body2" color="text.secondary" align="center">
            Send gifts to appear on the leaderboard!
          </Typography>
        </Paper>
      )}
    </Box>
  );
};

export default LeaderboardPanel;

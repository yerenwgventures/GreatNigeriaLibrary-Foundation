import React, { useState, useEffect } from 'react';
import { 
  Box, 
  Typography, 
  Grid, 
  Paper, 
  Tooltip, 
  Badge, 
  Avatar, 
  Dialog, 
  DialogTitle, 
  DialogContent, 
  DialogActions, 
  Button,
  useTheme,
  alpha,
  Chip
} from '@mui/material';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  EmojiEvents as TrophyIcon,
  School as EducationIcon,
  Whatshot as StreakIcon,
  Forum as CommunityIcon,
  Create as CreativeIcon,
  Psychology as KnowledgeIcon,
  Lightbulb as InsightIcon,
  Favorite as ContributionIcon,
  Star as StarIcon,
  Lock as LockIcon
} from '@mui/icons-material';

export interface UserBadge {
  id: string;
  name: string;
  description: string;
  category: 'achievement' | 'education' | 'streak' | 'community' | 'creative' | 'knowledge' | 'insight' | 'contribution';
  level: 'bronze' | 'silver' | 'gold' | 'platinum';
  icon?: React.ReactNode;
  earnedAt?: string;
  progress?: number;
  requiredProgress?: number;
  isNew?: boolean;
}

interface BadgeSystemProps {
  badges: UserBadge[];
  onBadgeClick?: (badge: UserBadge) => void;
  showNewBadgeNotification?: boolean;
  onNewBadgeDismiss?: () => void;
}

const BadgeSystem: React.FC<BadgeSystemProps> = ({
  badges,
  onBadgeClick,
  showNewBadgeNotification = false,
  onNewBadgeDismiss
}) => {
  const theme = useTheme();
  const [selectedBadge, setSelectedBadge] = useState<UserBadge | null>(null);
  const [newBadge, setNewBadge] = useState<UserBadge | null>(null);
  const [showDialog, setShowDialog] = useState(false);
  const [showNewBadgeDialog, setShowNewBadgeDialog] = useState(false);

  // Find the newest badge that is marked as new
  useEffect(() => {
    if (showNewBadgeNotification) {
      const newestBadge = badges.find(badge => badge.isNew);
      if (newestBadge) {
        setNewBadge(newestBadge);
        setShowNewBadgeDialog(true);
      }
    }
  }, [badges, showNewBadgeNotification]);

  const handleBadgeClick = (badge: UserBadge) => {
    setSelectedBadge(badge);
    setShowDialog(true);
    if (onBadgeClick) {
      onBadgeClick(badge);
    }
  };

  const handleCloseDialog = () => {
    setShowDialog(false);
    setSelectedBadge(null);
  };

  const handleCloseNewBadgeDialog = () => {
    setShowNewBadgeDialog(false);
    setNewBadge(null);
    if (onNewBadgeDismiss) {
      onNewBadgeDismiss();
    }
  };

  const getBadgeIcon = (badge: UserBadge) => {
    if (badge.icon) return badge.icon;

    switch (badge.category) {
      case 'achievement':
        return <TrophyIcon />;
      case 'education':
        return <EducationIcon />;
      case 'streak':
        return <StreakIcon />;
      case 'community':
        return <CommunityIcon />;
      case 'creative':
        return <CreativeIcon />;
      case 'knowledge':
        return <KnowledgeIcon />;
      case 'insight':
        return <InsightIcon />;
      case 'contribution':
        return <ContributionIcon />;
      default:
        return <StarIcon />;
    }
  };

  const getBadgeColor = (level: string) => {
    switch (level) {
      case 'bronze':
        return '#CD7F32';
      case 'silver':
        return '#C0C0C0';
      case 'gold':
        return '#FFD700';
      case 'platinum':
        return '#E5E4E2';
      default:
        return theme.palette.grey[500];
    }
  };

  const getProgressColor = (progress: number, required: number) => {
    const percentage = (progress / required) * 100;
    if (percentage >= 100) return theme.palette.success.main;
    if (percentage >= 75) return theme.palette.warning.main;
    if (percentage >= 50) return theme.palette.info.main;
    if (percentage >= 25) return theme.palette.primary.main;
    return theme.palette.grey[500];
  };

  // Group badges by category
  const groupedBadges = badges.reduce((acc, badge) => {
    if (!acc[badge.category]) {
      acc[badge.category] = [];
    }
    acc[badge.category].push(badge);
    return acc;
  }, {} as Record<string, UserBadge[]>);

  // Get category label
  const getCategoryLabel = (category: string) => {
    switch (category) {
      case 'achievement':
        return 'Achievements';
      case 'education':
        return 'Education';
      case 'streak':
        return 'Streaks';
      case 'community':
        return 'Community';
      case 'creative':
        return 'Creative';
      case 'knowledge':
        return 'Knowledge';
      case 'insight':
        return 'Insights';
      case 'contribution':
        return 'Contributions';
      default:
        return category.charAt(0).toUpperCase() + category.slice(1);
    }
  };

  return (
    <Box>
      {/* Badge Categories */}
      {Object.keys(groupedBadges).map((category) => (
        <Box key={category} mb={4}>
          <Typography variant="h6" gutterBottom>
            {getCategoryLabel(category)}
          </Typography>
          <Grid container spacing={2}>
            {groupedBadges[category].map((badge) => (
              <Grid item xs={6} sm={4} md={3} lg={2} key={badge.id}>
                <Tooltip
                  title={badge.earnedAt ? `${badge.name} - ${badge.description}` : 'Locked Badge'}
                  arrow
                >
                  <Paper
                    elevation={badge.earnedAt ? 2 : 0}
                    sx={{
                      p: 2,
                      textAlign: 'center',
                      cursor: 'pointer',
                      transition: 'all 0.3s ease',
                      backgroundColor: badge.earnedAt 
                        ? alpha(getBadgeColor(badge.level), 0.1)
                        : alpha(theme.palette.grey[500], 0.05),
                      border: `1px solid ${badge.earnedAt ? getBadgeColor(badge.level) : theme.palette.grey[300]}`,
                      '&:hover': {
                        transform: badge.earnedAt ? 'translateY(-5px)' : 'none',
                        boxShadow: badge.earnedAt ? theme.shadows[4] : theme.shadows[0],
                      },
                      position: 'relative',
                      overflow: 'hidden'
                    }}
                    onClick={() => badge.earnedAt && handleBadgeClick(badge)}
                  >
                    {badge.isNew && badge.earnedAt && (
                      <Box
                        sx={{
                          position: 'absolute',
                          top: 0,
                          right: 0,
                          backgroundColor: theme.palette.error.main,
                          color: theme.palette.error.contrastText,
                          px: 1,
                          py: 0.5,
                          borderBottomLeftRadius: 8,
                          fontSize: '0.7rem',
                          fontWeight: 'bold',
                          zIndex: 1
                        }}
                      >
                        NEW
                      </Box>
                    )}
                    
                    <Badge
                      overlap="circular"
                      anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                      badgeContent={
                        badge.level === 'platinum' ? (
                          <StarIcon sx={{ color: getBadgeColor(badge.level), fontSize: '1rem' }} />
                        ) : null
                      }
                    >
                      <Avatar
                        sx={{
                          width: 56,
                          height: 56,
                          margin: '0 auto',
                          backgroundColor: badge.earnedAt 
                            ? getBadgeColor(badge.level)
                            : theme.palette.grey[300],
                          color: badge.earnedAt 
                            ? theme.palette.getContrastText(getBadgeColor(badge.level))
                            : theme.palette.grey[500],
                          mb: 1,
                          filter: badge.earnedAt ? 'none' : 'grayscale(100%)'
                        }}
                      >
                        {badge.earnedAt ? getBadgeIcon(badge) : <LockIcon />}
                      </Avatar>
                    </Badge>
                    
                    <Typography 
                      variant="subtitle2" 
                      sx={{ 
                        fontWeight: 'bold',
                        color: badge.earnedAt ? 'text.primary' : 'text.disabled'
                      }}
                    >
                      {badge.name}
                    </Typography>
                    
                    {badge.progress !== undefined && badge.requiredProgress !== undefined && (
                      <Box sx={{ mt: 1, width: '100%' }}>
                        <Box 
                          sx={{ 
                            height: 4, 
                            width: '100%', 
                            backgroundColor: theme.palette.grey[300],
                            borderRadius: 2,
                            overflow: 'hidden'
                          }}
                        >
                          <Box 
                            sx={{ 
                              height: '100%', 
                              width: `${Math.min((badge.progress / badge.requiredProgress) * 100, 100)}%`,
                              backgroundColor: getProgressColor(badge.progress, badge.requiredProgress),
                              borderRadius: 2,
                              transition: 'width 1s ease-in-out'
                            }} 
                          />
                        </Box>
                        <Typography variant="caption" color="text.secondary">
                          {badge.progress}/{badge.requiredProgress}
                        </Typography>
                      </Box>
                    )}
                    
                    {badge.earnedAt && (
                      <Typography variant="caption" color="text.secondary" display="block" sx={{ mt: 0.5 }}>
                        Earned {new Date(badge.earnedAt).toLocaleDateString()}
                      </Typography>
                    )}
                  </Paper>
                </Tooltip>
              </Grid>
            ))}
          </Grid>
        </Box>
      ))}

      {/* Badge Detail Dialog */}
      <Dialog
        open={showDialog}
        onClose={handleCloseDialog}
        maxWidth="xs"
        fullWidth
      >
        {selectedBadge && (
          <>
            <DialogTitle sx={{ textAlign: 'center', pb: 1 }}>
              {selectedBadge.name}
            </DialogTitle>
            <DialogContent sx={{ textAlign: 'center', pt: 1 }}>
              <Avatar
                sx={{
                  width: 80,
                  height: 80,
                  margin: '0 auto 16px auto',
                  backgroundColor: getBadgeColor(selectedBadge.level),
                  color: theme.palette.getContrastText(getBadgeColor(selectedBadge.level))
                }}
              >
                {getBadgeIcon(selectedBadge)}
              </Avatar>
              
              <Typography variant="body1" paragraph>
                {selectedBadge.description}
              </Typography>
              
              <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', mb: 2 }}>
                <Chip 
                  label={selectedBadge.level.toUpperCase()} 
                  sx={{ 
                    backgroundColor: getBadgeColor(selectedBadge.level),
                    color: theme.palette.getContrastText(getBadgeColor(selectedBadge.level)),
                    fontWeight: 'bold'
                  }} 
                />
              </Box>
              
              {selectedBadge.earnedAt && (
                <Typography variant="body2" color="text.secondary">
                  Earned on {new Date(selectedBadge.earnedAt).toLocaleDateString()}
                </Typography>
              )}
              
              {selectedBadge.progress !== undefined && selectedBadge.requiredProgress !== undefined && (
                <Box sx={{ mt: 2, width: '100%' }}>
                  <Typography variant="body2" color="text.secondary" gutterBottom>
                    Progress: {selectedBadge.progress}/{selectedBadge.requiredProgress}
                  </Typography>
                  <Box 
                    sx={{ 
                      height: 8, 
                      width: '100%', 
                      backgroundColor: theme.palette.grey[300],
                      borderRadius: 4,
                      overflow: 'hidden'
                    }}
                  >
                    <Box 
                      sx={{ 
                        height: '100%', 
                        width: `${Math.min((selectedBadge.progress / selectedBadge.requiredProgress) * 100, 100)}%`,
                        backgroundColor: getProgressColor(selectedBadge.progress, selectedBadge.requiredProgress),
                        borderRadius: 4,
                        transition: 'width 1s ease-in-out'
                      }} 
                    />
                  </Box>
                </Box>
              )}
            </DialogContent>
            <DialogActions>
              <Button onClick={handleCloseDialog} color="primary">
                Close
              </Button>
              <Button 
                color="primary" 
                variant="contained"
                onClick={handleCloseDialog}
              >
                Share Badge
              </Button>
            </DialogActions>
          </>
        )}
      </Dialog>

      {/* New Badge Notification Dialog */}
      <Dialog
        open={showNewBadgeDialog}
        onClose={handleCloseNewBadgeDialog}
        maxWidth="xs"
        fullWidth
        PaperProps={{
          sx: {
            borderRadius: 4,
            overflow: 'hidden'
          }
        }}
      >
        <AnimatePresence>
          {newBadge && (
            <motion.div
              initial={{ opacity: 0, scale: 0.8 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.8 }}
              transition={{ duration: 0.5 }}
            >
              <Box 
                sx={{ 
                  p: 3, 
                  textAlign: 'center',
                  background: `linear-gradient(135deg, ${alpha(getBadgeColor(newBadge.level), 0.2)} 0%, ${alpha(getBadgeColor(newBadge.level), 0.1)} 100%)`,
                }}
              >
                <Typography variant="h5" gutterBottom fontWeight="bold">
                  New Badge Earned!
                </Typography>
                
                <motion.div
                  initial={{ y: 20, opacity: 0 }}
                  animate={{ y: 0, opacity: 1 }}
                  transition={{ delay: 0.3, duration: 0.5 }}
                >
                  <Avatar
                    sx={{
                      width: 100,
                      height: 100,
                      margin: '16px auto',
                      backgroundColor: getBadgeColor(newBadge.level),
                      color: theme.palette.getContrastText(getBadgeColor(newBadge.level)),
                      border: `4px solid ${theme.palette.background.paper}`,
                      boxShadow: theme.shadows[10]
                    }}
                  >
                    {getBadgeIcon(newBadge)}
                  </Avatar>
                </motion.div>
                
                <motion.div
                  initial={{ y: 20, opacity: 0 }}
                  animate={{ y: 0, opacity: 1 }}
                  transition={{ delay: 0.5, duration: 0.5 }}
                >
                  <Typography variant="h6" gutterBottom>
                    {newBadge.name}
                  </Typography>
                  
                  <Typography variant="body1" paragraph>
                    {newBadge.description}
                  </Typography>
                  
                  <Chip 
                    label={newBadge.level.toUpperCase()} 
                    sx={{ 
                      backgroundColor: getBadgeColor(newBadge.level),
                      color: theme.palette.getContrastText(getBadgeColor(newBadge.level)),
                      fontWeight: 'bold',
                      mb: 2
                    }} 
                  />
                </motion.div>
                
                {/* Confetti effect */}
                {[...Array(30)].map((_, i) => (
                  <motion.div
                    key={i}
                    style={{
                      position: 'absolute',
                      width: Math.random() * 10 + 5,
                      height: Math.random() * 10 + 5,
                      borderRadius: '50%',
                      backgroundColor: [
                        '#FFC700', '#FF0000', '#2E7D32', '#0288D1', '#9C27B0', getBadgeColor(newBadge.level)
                      ][i % 6],
                      zIndex: 0
                    }}
                    initial={{ 
                      x: 0, 
                      y: 0,
                      opacity: 1
                    }}
                    animate={{ 
                      x: (Math.random() - 0.5) * 300, 
                      y: Math.random() * 300,
                      opacity: 0
                    }}
                    transition={{ 
                      duration: 2 + Math.random() * 2,
                      ease: "easeOut",
                      delay: Math.random() * 0.5
                    }}
                  />
                ))}
                
                <Box sx={{ mt: 2, display: 'flex', justifyContent: 'center', gap: 2 }}>
                  <Button 
                    variant="outlined" 
                    onClick={handleCloseNewBadgeDialog}
                  >
                    Close
                  </Button>
                  <Button 
                    variant="contained" 
                    color="primary"
                    onClick={handleCloseNewBadgeDialog}
                  >
                    View All Badges
                  </Button>
                </Box>
              </Box>
            </motion.div>
          )}
        </AnimatePresence>
      </Dialog>
    </Box>
  );
};

export default BadgeSystem;

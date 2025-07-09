import React, { useState, useEffect } from 'react';
import { 
  Dialog, 
  DialogContent, 
  Typography, 
  Box, 
  Button, 
  Avatar, 
  Grow,
  useTheme
} from '@mui/material';
import { motion } from 'framer-motion';
import Confetti from 'react-confetti';
import { 
  EmojiEvents as TrophyIcon,
  Star as StarIcon,
  School as EducationIcon,
  LocalLibrary as BookIcon,
  Lightbulb as IdeaIcon
} from '@mui/icons-material';

interface MilestoneAchievementProps {
  open: boolean;
  onClose: () => void;
  title: string;
  description: string;
  type: 'chapter' | 'book' | 'points' | 'streak' | 'custom';
  level?: 'bronze' | 'silver' | 'gold' | 'platinum';
  pointsAwarded?: number;
  icon?: React.ReactNode;
}

const MilestoneAchievement: React.FC<MilestoneAchievementProps> = ({
  open,
  onClose,
  title,
  description,
  type,
  level = 'bronze',
  pointsAwarded = 0,
  icon
}) => {
  const theme = useTheme();
  const [showConfetti, setShowConfetti] = useState(false);
  const [windowDimensions, setWindowDimensions] = useState({
    width: window.innerWidth,
    height: window.innerHeight
  });

  // Update window dimensions when window is resized
  useEffect(() => {
    const handleResize = () => {
      setWindowDimensions({
        width: window.innerWidth,
        height: window.innerHeight
      });
    };

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  // Show confetti when dialog opens
  useEffect(() => {
    if (open) {
      setShowConfetti(true);
      const timer = setTimeout(() => {
        setShowConfetti(false);
      }, 5000);
      return () => clearTimeout(timer);
    }
  }, [open]);

  // Get icon based on achievement type
  const getIcon = () => {
    if (icon) return icon;
    
    switch (type) {
      case 'chapter':
        return <BookIcon fontSize="large" />;
      case 'book':
        return <LocalLibrary fontSize="large" />;
      case 'points':
        return <StarIcon fontSize="large" />;
      case 'streak':
        return <IdeaIcon fontSize="large" />;
      default:
        return <TrophyIcon fontSize="large" />;
    }
  };

  // Get color based on achievement level
  const getLevelColor = () => {
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
        return '#CD7F32';
    }
  };

  // Get background gradient based on achievement level
  const getBackgroundGradient = () => {
    switch (level) {
      case 'bronze':
        return 'linear-gradient(135deg, #CD7F32 0%, #8B4513 100%)';
      case 'silver':
        return 'linear-gradient(135deg, #C0C0C0 0%, #A9A9A9 100%)';
      case 'gold':
        return 'linear-gradient(135deg, #FFD700 0%, #FFA500 100%)';
      case 'platinum':
        return 'linear-gradient(135deg, #E5E4E2 0%, #AFAFAF 100%)';
      default:
        return 'linear-gradient(135deg, #CD7F32 0%, #8B4513 100%)';
    }
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      maxWidth="sm"
      fullWidth
      TransitionComponent={Grow}
      PaperProps={{
        sx: {
          borderRadius: 4,
          overflow: 'hidden',
          backgroundImage: 'radial-gradient(circle, rgba(0,0,0,0.8) 0%, rgba(0,0,0,0.9) 100%)',
          color: 'white',
          boxShadow: '0 8px 32px rgba(0, 0, 0, 0.3)'
        }
      }}
    >
      {showConfetti && (
        <Confetti
          width={windowDimensions.width}
          height={windowDimensions.height}
          recycle={false}
          numberOfPieces={200}
          gravity={0.15}
        />
      )}
      
      <DialogContent sx={{ p: 4, textAlign: 'center' }}>
        <motion.div
          initial={{ scale: 0 }}
          animate={{ scale: 1 }}
          transition={{ type: 'spring', stiffness: 260, damping: 20 }}
        >
          <Avatar
            sx={{
              width: 120,
              height: 120,
              margin: '0 auto 24px',
              background: getBackgroundGradient(),
              color: 'white',
              boxShadow: '0 8px 16px rgba(0, 0, 0, 0.3)'
            }}
          >
            {getIcon()}
          </Avatar>
        </motion.div>
        
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.3 }}
        >
          <Typography variant="h4" component="h2" gutterBottom fontWeight="bold">
            {title}
          </Typography>
        </motion.div>
        
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
        >
          <Typography variant="body1" paragraph>
            {description}
          </Typography>
        </motion.div>
        
        {pointsAwarded > 0 && (
          <motion.div
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ delay: 0.9 }}
          >
            <Box
              sx={{
                display: 'inline-flex',
                alignItems: 'center',
                bgcolor: 'rgba(255, 255, 255, 0.1)',
                borderRadius: 2,
                px: 2,
                py: 1,
                mb: 3
              }}
            >
              <StarIcon sx={{ color: theme.palette.warning.main, mr: 1 }} />
              <Typography variant="h6" color={theme.palette.warning.main}>
                +{pointsAwarded} Points Awarded
              </Typography>
            </Box>
          </motion.div>
        )}
        
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 1.2 }}
        >
          <Button
            variant="contained"
            color="primary"
            size="large"
            onClick={onClose}
            sx={{
              mt: 2,
              px: 4,
              py: 1.5,
              borderRadius: 8,
              fontWeight: 'bold',
              textTransform: 'none',
              fontSize: '1.1rem'
            }}
          >
            Continue
          </Button>
        </motion.div>
      </DialogContent>
    </Dialog>
  );
};

export default MilestoneAchievement;

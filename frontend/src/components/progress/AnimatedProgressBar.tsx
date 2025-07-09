import React, { useEffect, useState } from 'react';
import { Box, Typography, LinearProgress, Tooltip, useTheme } from '@mui/material';
import { motion, AnimatePresence } from 'framer-motion';
import { CheckCircle, Star, EmojiEvents } from '@mui/icons-material';

interface AnimatedProgressBarProps {
  value: number;
  total: number;
  label?: string;
  showPercentage?: boolean;
  height?: number;
  milestones?: number[];
  onMilestoneReached?: (milestone: number) => void;
  color?: 'primary' | 'secondary' | 'success' | 'info' | 'warning' | 'error';
}

const AnimatedProgressBar: React.FC<AnimatedProgressBarProps> = ({
  value,
  total,
  label,
  showPercentage = true,
  height = 10,
  milestones = [],
  onMilestoneReached,
  color = 'primary'
}) => {
  const theme = useTheme();
  const [prevValue, setPrevValue] = useState(value);
  const [animatedValue, setAnimatedValue] = useState(value);
  const [reachedMilestones, setReachedMilestones] = useState<number[]>([]);
  const [showCelebration, setShowCelebration] = useState(false);
  const [celebrationMilestone, setCelebrationMilestone] = useState<number | null>(null);

  // Calculate percentage
  const percentage = Math.min(Math.round((value / total) * 100), 100);
  const animatedPercentage = Math.min(Math.round((animatedValue / total) * 100), 100);

  // Check for milestone achievements
  useEffect(() => {
    if (value > prevValue) {
      // Animate the progress bar
      setPrevValue(value);
      
      // Animate from previous value to new value
      const animationDuration = 1000; // 1 second
      const startTime = Date.now();
      const startValue = prevValue;
      const endValue = value;
      
      const animateProgress = () => {
        const elapsed = Date.now() - startTime;
        const progress = Math.min(elapsed / animationDuration, 1);
        const currentValue = startValue + progress * (endValue - startValue);
        
        setAnimatedValue(currentValue);
        
        if (progress < 1) {
          requestAnimationFrame(animateProgress);
        }
      };
      
      requestAnimationFrame(animateProgress);
      
      // Check for newly reached milestones
      const newMilestones = milestones.filter(
        milestone => 
          !reachedMilestones.includes(milestone) && 
          (value / total) * 100 >= milestone
      );
      
      if (newMilestones.length > 0) {
        setReachedMilestones([...reachedMilestones, ...newMilestones]);
        
        // Trigger celebration for the highest milestone reached
        const highestMilestone = Math.max(...newMilestones);
        setCelebrationMilestone(highestMilestone);
        setShowCelebration(true);
        
        // Notify parent component
        if (onMilestoneReached) {
          newMilestones.forEach(milestone => onMilestoneReached(milestone));
        }
        
        // Hide celebration after 3 seconds
        setTimeout(() => {
          setShowCelebration(false);
          setCelebrationMilestone(null);
        }, 3000);
      }
    }
  }, [value, total, prevValue, milestones, reachedMilestones, onMilestoneReached]);

  // Get milestone icon based on percentage
  const getMilestoneIcon = (milestone: number) => {
    if (milestone >= 100) return <EmojiEvents fontSize="small" />;
    if (milestone >= 75) return <Star fontSize="small" />;
    return <CheckCircle fontSize="small" />;
  };

  // Get milestone color based on percentage
  const getMilestoneColor = (milestone: number) => {
    if (milestone >= 100) return theme.palette.success.main;
    if (milestone >= 75) return theme.palette.warning.main;
    if (milestone >= 50) return theme.palette.info.main;
    if (milestone >= 25) return theme.palette.primary.main;
    return theme.palette.grey[500];
  };

  return (
    <Box sx={{ width: '100%', position: 'relative' }}>
      {label && (
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={0.5}>
          <Typography variant="body2" color="text.secondary">
            {label}
          </Typography>
          {showPercentage && (
            <Typography variant="body2" color="text.secondary">
              {animatedPercentage}%
            </Typography>
          )}
        </Box>
      )}
      
      <Box sx={{ position: 'relative', height: height }}>
        <LinearProgress
          variant="determinate"
          value={animatedPercentage}
          color={color}
          sx={{ 
            height: height,
            borderRadius: height / 2,
            '& .MuiLinearProgress-bar': {
              transition: 'none' // Disable default transition for custom animation
            }
          }}
        />
        
        {/* Milestone markers */}
        {milestones.map((milestone) => (
          <Tooltip 
            key={milestone} 
            title={`${milestone}% Complete`} 
            placement="top"
            arrow
          >
            <Box
              sx={{
                position: 'absolute',
                left: `${milestone}%`,
                top: '50%',
                transform: 'translate(-50%, -50%)',
                width: height * 1.5,
                height: height * 1.5,
                borderRadius: '50%',
                bgcolor: reachedMilestones.includes(milestone) 
                  ? getMilestoneColor(milestone)
                  : theme.palette.grey[300],
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
                color: reachedMilestones.includes(milestone) 
                  ? theme.palette.common.white
                  : theme.palette.grey[500],
                zIndex: 1,
                transition: 'all 0.3s ease',
                cursor: 'pointer'
              }}
            >
              {getMilestoneIcon(milestone)}
            </Box>
          </Tooltip>
        ))}
      </Box>
      
      {/* Celebration animation */}
      <AnimatePresence>
        {showCelebration && celebrationMilestone && (
          <motion.div
            initial={{ opacity: 0, scale: 0.5 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.5 }}
            style={{
              position: 'absolute',
              top: -50,
              left: '50%',
              transform: 'translateX(-50%)',
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              zIndex: 10
            }}
          >
            <Box
              sx={{
                bgcolor: getMilestoneColor(celebrationMilestone),
                color: 'white',
                py: 1,
                px: 2,
                borderRadius: 2,
                display: 'flex',
                alignItems: 'center',
                gap: 1
              }}
            >
              {getMilestoneIcon(celebrationMilestone)}
              <Typography variant="body1" fontWeight="bold">
                {celebrationMilestone}% Milestone Reached!
              </Typography>
            </Box>
            
            {/* Confetti effect */}
            {[...Array(20)].map((_, i) => (
              <motion.div
                key={i}
                style={{
                  position: 'absolute',
                  width: 10,
                  height: 10,
                  borderRadius: '50%',
                  backgroundColor: [
                    '#FFC700', '#FF0000', '#2E7D32', '#0288D1', '#9C27B0'
                  ][i % 5],
                  zIndex: -1
                }}
                initial={{ 
                  x: 0, 
                  y: 0 
                }}
                animate={{ 
                  x: (Math.random() - 0.5) * 200, 
                  y: Math.random() * 100 + 50,
                  opacity: [1, 0]
                }}
                transition={{ 
                  duration: 2,
                  ease: "easeOut"
                }}
              />
            ))}
          </motion.div>
        )}
      </AnimatePresence>
    </Box>
  );
};

export default AnimatedProgressBar;

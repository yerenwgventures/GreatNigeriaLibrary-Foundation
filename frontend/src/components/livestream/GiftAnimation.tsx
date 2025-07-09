import React, { useEffect, useState } from 'react';
import { Box, Typography, Avatar } from '@mui/material';
import { motion, AnimatePresence } from 'framer-motion';
import { Gift } from '../../api/livestreamService';

interface GiftAnimationProps {
  gift: Gift | null;
  onAnimationComplete: () => void;
}

// Sample gift animations data
const GIFT_ANIMATIONS = {
  1: { // Heart
    image: 'https://via.placeholder.com/100?text=❤️',
    duration: 3000,
    animation: 'float'
  },
  2: { // Star
    image: 'https://via.placeholder.com/100?text=⭐',
    duration: 3000,
    animation: 'spin'
  },
  3: { // Trophy
    image: 'https://via.placeholder.com/100?text=🏆',
    duration: 3500,
    animation: 'bounce'
  },
  4: { // Diamond
    image: 'https://via.placeholder.com/100?text=💎',
    duration: 4000,
    animation: 'pulse'
  },
  5: { // Crown
    image: 'https://via.placeholder.com/100?text=👑',
    duration: 4500,
    animation: 'zoom'
  },
  6: { // Rocket
    image: 'https://via.placeholder.com/100?text=🚀',
    duration: 5000,
    animation: 'rocket'
  }
};

const GiftAnimation: React.FC<GiftAnimationProps> = ({ gift, onAnimationComplete }) => {
  const [isVisible, setIsVisible] = useState(false);
  
  useEffect(() => {
    if (gift) {
      setIsVisible(true);
      
      // Get animation duration based on gift type
      const duration = GIFT_ANIMATIONS[gift.giftId as keyof typeof GIFT_ANIMATIONS]?.duration || 3000;
      
      // Hide animation after duration
      const timer = setTimeout(() => {
        setIsVisible(false);
        
        // Notify parent component that animation is complete
        setTimeout(() => {
          onAnimationComplete();
        }, 300); // Small delay to ensure exit animation completes
      }, duration);
      
      return () => clearTimeout(timer);
    }
  }, [gift, onAnimationComplete]);
  
  if (!gift) return null;
  
  const giftAnimation = GIFT_ANIMATIONS[gift.giftId as keyof typeof GIFT_ANIMATIONS];
  if (!giftAnimation) return null;
  
  // Different animation variants based on gift type
  const getAnimationVariants = () => {
    switch (giftAnimation.animation) {
      case 'float':
        return {
          initial: { opacity: 0, y: 100 },
          animate: { 
            opacity: 1, 
            y: -100,
            transition: { 
              duration: giftAnimation.duration / 1000,
              ease: 'easeOut'
            }
          },
          exit: { opacity: 0 }
        };
      case 'spin':
        return {
          initial: { opacity: 0, scale: 0.5 },
          animate: { 
            opacity: 1, 
            scale: 1.5, 
            rotate: 360,
            transition: { 
              duration: giftAnimation.duration / 1000,
              ease: 'easeOut'
            }
          },
          exit: { opacity: 0, scale: 0 }
        };
      case 'bounce':
        return {
          initial: { opacity: 0, y: 100 },
          animate: { 
            opacity: 1, 
            y: 0,
            transition: { 
              duration: giftAnimation.duration / 1000,
              type: 'spring',
              stiffness: 100
            }
          },
          exit: { opacity: 0, y: -100 }
        };
      case 'pulse':
        return {
          initial: { opacity: 0, scale: 0.5 },
          animate: { 
            opacity: 1, 
            scale: [1, 1.2, 1, 1.2, 1],
            transition: { 
              duration: giftAnimation.duration / 1000,
              times: [0, 0.25, 0.5, 0.75, 1]
            }
          },
          exit: { opacity: 0, scale: 0 }
        };
      case 'zoom':
        return {
          initial: { opacity: 0, scale: 0 },
          animate: { 
            opacity: 1, 
            scale: 1.5,
            transition: { 
              duration: giftAnimation.duration / 1000 * 0.5,
              ease: 'easeOut'
            }
          },
          exit: { 
            opacity: 0, 
            scale: 2,
            transition: { 
              duration: giftAnimation.duration / 1000 * 0.5
            }
          }
        };
      case 'rocket':
        return {
          initial: { opacity: 0, x: -300, y: 100 },
          animate: { 
            opacity: 1, 
            x: 300,
            y: -100,
            transition: { 
              duration: giftAnimation.duration / 1000,
              ease: 'easeOut'
            }
          },
          exit: { opacity: 0 }
        };
      default:
        return {
          initial: { opacity: 0 },
          animate: { opacity: 1 },
          exit: { opacity: 0 }
        };
    }
  };
  
  return (
    <AnimatePresence>
      {isVisible && (
        <Box
          sx={{
            position: 'absolute',
            top: 0,
            left: 0,
            width: '100%',
            height: '100%',
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            justifyContent: 'center',
            pointerEvents: 'none',
            zIndex: 10
          }}
        >
          <motion.div
            initial="initial"
            animate="animate"
            exit="exit"
            variants={getAnimationVariants()}
          >
            <Box
              sx={{
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'center',
                position: 'relative'
              }}
            >
              <Box
                component="img"
                src={giftAnimation.image}
                alt={gift.giftName}
                sx={{ width: 100, height: 100 }}
              />
              
              <Box
                sx={{
                  display: 'flex',
                  alignItems: 'center',
                  bgcolor: 'rgba(0, 0, 0, 0.7)',
                  borderRadius: 2,
                  px: 2,
                  py: 0.5,
                  mt: 1
                }}
              >
                <Avatar sx={{ width: 24, height: 24, mr: 1 }} />
                <Typography variant="body2" color="white">
                  {gift.isAnonymous ? 'Anonymous' : `User #${gift.senderId}`}
                </Typography>
                <Typography variant="body2" color="white" sx={{ mx: 1 }}>
                  sent
                </Typography>
                <Typography variant="body2" color="white" fontWeight="bold">
                  {gift.giftName}
                  {gift.comboCount > 1 && ` x${gift.comboCount}`}
                </Typography>
              </Box>
            </Box>
          </motion.div>
        </Box>
      )}
    </AnimatePresence>
  );
};

export default GiftAnimation;

import React, { useState, useEffect } from 'react';
import GiftAnimation from './GiftAnimation';
import { Gift } from '../../api/livestreamService';

interface GiftAnimationManagerProps {
  gifts: Gift[];
}

const GiftAnimationManager: React.FC<GiftAnimationManagerProps> = ({ gifts }) => {
  const [giftQueue, setGiftQueue] = useState<Gift[]>([]);
  const [currentGift, setCurrentGift] = useState<Gift | null>(null);
  const [isAnimating, setIsAnimating] = useState(false);
  
  // Process new gifts
  useEffect(() => {
    if (gifts.length > 0) {
      // Get the latest gift
      const latestGift = gifts[0];
      
      // Check if it's already in the queue or currently animating
      const isInQueue = giftQueue.some(gift => gift.id === latestGift.id);
      const isCurrentlyAnimating = currentGift && currentGift.id === latestGift.id;
      
      if (!isInQueue && !isCurrentlyAnimating) {
        // Add to queue
        setGiftQueue(prevQueue => [...prevQueue, latestGift]);
      }
    }
  }, [gifts, giftQueue, currentGift]);
  
  // Process gift queue
  useEffect(() => {
    if (giftQueue.length > 0 && !isAnimating) {
      // Get the next gift from the queue
      const nextGift = giftQueue[0];
      
      // Remove it from the queue
      setGiftQueue(prevQueue => prevQueue.slice(1));
      
      // Set as current gift
      setCurrentGift(nextGift);
      setIsAnimating(true);
    }
  }, [giftQueue, isAnimating]);
  
  // Handle animation completion
  const handleAnimationComplete = () => {
    setCurrentGift(null);
    setIsAnimating(false);
  };
  
  return (
    <GiftAnimation 
      gift={currentGift} 
      onAnimationComplete={handleAnimationComplete} 
    />
  );
};

export default GiftAnimationManager;

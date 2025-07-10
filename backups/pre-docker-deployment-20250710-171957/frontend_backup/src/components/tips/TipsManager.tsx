import React, { useState, useEffect } from 'react';
import { Box } from '@mui/material';
import ContextualTip, { TipType } from './ContextualTip';

export interface Tip {
  id: string;
  type: TipType;
  title: string;
  content: string;
  trigger: string;
  seen?: boolean;
  dismissed?: boolean;
  dontShowAgain?: boolean;
  autoHideDuration?: number;
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right' | 'center';
  showDontShowAgain?: boolean;
  priority?: 'low' | 'medium' | 'high';
  condition?: () => boolean;
}

interface TipsManagerProps {
  tips: Tip[];
  currentTrigger?: string;
  onTipDismiss?: (tipId: string, dontShowAgain: boolean) => void;
  maxVisibleTips?: number;
}

const TipsManager: React.FC<TipsManagerProps> = ({
  tips,
  currentTrigger,
  onTipDismiss,
  maxVisibleTips = 3
}) => {
  const [visibleTips, setVisibleTips] = useState<Tip[]>([]);
  const [dismissedTips, setDismissedTips] = useState<Record<string, boolean>>({});

  // Update visible tips when triggers change
  useEffect(() => {
    if (currentTrigger) {
      // Find tips that match the current trigger and haven't been dismissed
      const matchingTips = tips.filter(tip => 
        tip.trigger === currentTrigger && 
        !dismissedTips[tip.id] && 
        (!tip.condition || tip.condition())
      );
      
      // Add new matching tips to visible tips, respecting the max limit
      if (matchingTips.length > 0) {
        setVisibleTips(prevTips => {
          // Filter out any tips that are already visible
          const newTips = matchingTips.filter(
            newTip => !prevTips.some(prevTip => prevTip.id === newTip.id)
          );
          
          // Combine existing and new tips, then sort by priority
          const allTips = [...prevTips, ...newTips].sort((a, b) => {
            const priorityOrder = { high: 0, medium: 1, low: 2 };
            return (priorityOrder[a.priority || 'medium'] || 1) - (priorityOrder[b.priority || 'medium'] || 1);
          });
          
          // Limit to max visible tips
          return allTips.slice(0, maxVisibleTips);
        });
      }
    }
  }, [currentTrigger, tips, dismissedTips, maxVisibleTips]);

  // Handle tip dismissal
  const handleTipDismiss = (tipId: string, dontShowAgain: boolean) => {
    // Remove the tip from visible tips
    setVisibleTips(prevTips => prevTips.filter(tip => tip.id !== tipId));
    
    // Mark the tip as dismissed
    setDismissedTips(prevDismissed => ({
      ...prevDismissed,
      [tipId]: dontShowAgain
    }));
    
    // Call the onTipDismiss callback if provided
    if (onTipDismiss) {
      onTipDismiss(tipId, dontShowAgain);
    }
  };

  return (
    <Box>
      {visibleTips.map(tip => (
        <ContextualTip
          key={tip.id}
          id={tip.id}
          type={tip.type}
          title={tip.title}
          content={tip.content}
          onDismiss={handleTipDismiss}
          autoHideDuration={tip.autoHideDuration}
          position={tip.position}
          showDontShowAgain={tip.showDontShowAgain}
          priority={tip.priority}
        />
      ))}
    </Box>
  );
};

export default TipsManager;

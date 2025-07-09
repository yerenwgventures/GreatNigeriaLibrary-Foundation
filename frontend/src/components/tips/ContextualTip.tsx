import React, { useState, useEffect } from 'react';
import { 
  Paper, 
  Typography, 
  Box, 
  IconButton, 
  Collapse, 
  Checkbox, 
  FormControlLabel,
  useTheme,
  Fade
} from '@mui/material';
import { 
  Close as CloseIcon,
  Lightbulb as TipIcon,
  Info as InfoIcon,
  Star as StarIcon,
  Bolt as BoltIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

export type TipType = 'info' | 'tip' | 'suggestion' | 'achievement';

interface ContextualTipProps {
  id: string;
  type: TipType;
  title: string;
  content: string;
  onDismiss: (id: string, dontShowAgain: boolean) => void;
  autoHideDuration?: number;
  position?: 'top-left' | 'top-right' | 'bottom-left' | 'bottom-right' | 'center';
  showDontShowAgain?: boolean;
  priority?: 'low' | 'medium' | 'high';
}

const ContextualTip: React.FC<ContextualTipProps> = ({
  id,
  type,
  title,
  content,
  onDismiss,
  autoHideDuration = 0,
  position = 'bottom-right',
  showDontShowAgain = true,
  priority = 'medium'
}) => {
  const theme = useTheme();
  const [open, setOpen] = useState(true);
  const [expanded, setExpanded] = useState(false);
  const [dontShowAgain, setDontShowAgain] = useState(false);

  // Auto-hide tip after specified duration
  useEffect(() => {
    if (autoHideDuration > 0) {
      const timer = setTimeout(() => {
        handleClose();
      }, autoHideDuration);
      return () => clearTimeout(timer);
    }
  }, [autoHideDuration]);

  // Handle close button click
  const handleClose = () => {
    setOpen(false);
    // Allow animation to complete before calling onDismiss
    setTimeout(() => {
      onDismiss(id, dontShowAgain);
    }, 300);
  };

  // Get icon based on tip type
  const getIcon = () => {
    switch (type) {
      case 'info':
        return <InfoIcon />;
      case 'tip':
        return <TipIcon />;
      case 'suggestion':
        return <StarIcon />;
      case 'achievement':
        return <BoltIcon />;
      default:
        return <TipIcon />;
    }
  };

  // Get color based on tip type
  const getColor = () => {
    switch (type) {
      case 'info':
        return theme.palette.info.main;
      case 'tip':
        return theme.palette.primary.main;
      case 'suggestion':
        return theme.palette.warning.main;
      case 'achievement':
        return theme.palette.success.main;
      default:
        return theme.palette.primary.main;
    }
  };

  // Get position styles
  const getPositionStyles = () => {
    switch (position) {
      case 'top-left':
        return { top: 16, left: 16 };
      case 'top-right':
        return { top: 16, right: 16 };
      case 'bottom-left':
        return { bottom: 16, left: 16 };
      case 'bottom-right':
        return { bottom: 16, right: 16 };
      case 'center':
        return { 
          top: '50%', 
          left: '50%', 
          transform: 'translate(-50%, -50%)' 
        };
      default:
        return { bottom: 16, right: 16 };
    }
  };

  // Get z-index based on priority
  const getZIndex = () => {
    switch (priority) {
      case 'low':
        return 1000;
      case 'medium':
        return 1100;
      case 'high':
        return 1200;
      default:
        return 1000;
    }
  };

  return (
    <Fade in={open}>
      <Box
        sx={{
          position: 'fixed',
          ...getPositionStyles(),
          zIndex: getZIndex(),
          maxWidth: 350,
          width: '100%'
        }}
      >
        <motion.div
          initial={{ scale: 0.8, opacity: 0 }}
          animate={{ scale: 1, opacity: 1 }}
          exit={{ scale: 0.8, opacity: 0 }}
          transition={{ type: 'spring', stiffness: 500, damping: 30 }}
        >
          <Paper
            elevation={6}
            sx={{
              borderRadius: 2,
              overflow: 'hidden',
              border: `1px solid ${getColor()}`,
              boxShadow: `0 4px 20px rgba(0, 0, 0, 0.1), 0 0 0 1px ${getColor()}22`
            }}
          >
            <Box
              sx={{
                display: 'flex',
                alignItems: 'center',
                bgcolor: getColor(),
                color: 'white',
                px: 2,
                py: 1,
                cursor: 'pointer'
              }}
              onClick={() => setExpanded(!expanded)}
            >
              <Box sx={{ mr: 1 }}>
                {getIcon()}
              </Box>
              <Typography variant="subtitle1" fontWeight="medium" sx={{ flexGrow: 1 }}>
                {title}
              </Typography>
              <IconButton 
                size="small" 
                onClick={(e) => {
                  e.stopPropagation();
                  handleClose();
                }}
                sx={{ color: 'white' }}
              >
                <CloseIcon fontSize="small" />
              </IconButton>
            </Box>
            
            <Collapse in={expanded} collapsedSize={40}>
              <Box sx={{ p: 2 }}>
                <Typography variant="body2" color="text.secondary">
                  {content}
                </Typography>
                
                {showDontShowAgain && (
                  <Box sx={{ mt: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                    <FormControlLabel
                      control={
                        <Checkbox
                          size="small"
                          checked={dontShowAgain}
                          onChange={(e) => setDontShowAgain(e.target.checked)}
                        />
                      }
                      label={
                        <Typography variant="caption" color="text.secondary">
                          Don't show again
                        </Typography>
                      }
                    />
                  </Box>
                )}
              </Box>
            </Collapse>
            
            {!expanded && (
              <Box 
                sx={{ 
                  px: 2, 
                  py: 0.5, 
                  cursor: 'pointer',
                  color: 'text.secondary',
                  fontSize: '0.75rem',
                  textAlign: 'center'
                }}
                onClick={() => setExpanded(true)}
              >
                Click to expand
              </Box>
            )}
          </Paper>
        </motion.div>
      </Box>
    </Fade>
  );
};

export default ContextualTip;

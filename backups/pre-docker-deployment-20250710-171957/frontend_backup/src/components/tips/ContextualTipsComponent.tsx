import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { useLocation } from 'react-router-dom';
import { 
  Box, 
  Card, 
  CardContent, 
  CardActions, 
  Typography, 
  Button, 
  IconButton, 
  Snackbar, 
  Alert, 
  Slide, 
  Fade, 
  Tooltip,
  useTheme,
  useMediaQuery,
  styled
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import ThumbUpIcon from '@mui/icons-material/ThumbUp';
import ThumbDownIcon from '@mui/icons-material/ThumbDown';
import LightbulbIcon from '@mui/icons-material/Lightbulb';
import { 
  fetchContextualTips, 
  recordTipView, 
  recordTipDismiss, 
  recordTipClick, 
  submitTipFeedback,
  selectContextualTips,
  selectTipsLoading,
  clearContextualTips
} from '../../features/tips/tipsSlice';
import { useAuth } from '../../hooks/useAuth';
import { AppDispatch } from '../../store';
import { Tip } from '../../api/tipsService';

// Styled components
const TipCard = styled(Card)(({ theme }) => ({
  position: 'relative',
  maxWidth: 400,
  margin: theme.spacing(1),
  boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
  borderRadius: theme.spacing(1),
  overflow: 'visible',
  transition: 'transform 0.3s ease, box-shadow 0.3s ease',
  '&:hover': {
    transform: 'translateY(-4px)',
    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.2)',
  },
}));

const TipIcon = styled(Box)(({ theme }) => ({
  position: 'absolute',
  top: -16,
  left: 16,
  width: 32,
  height: 32,
  borderRadius: '50%',
  backgroundColor: theme.palette.primary.main,
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: theme.palette.primary.contrastText,
  boxShadow: '0 2px 8px rgba(0, 0, 0, 0.2)',
}));

const FeedbackContainer = styled(Box)(({ theme }) => ({
  display: 'flex',
  alignItems: 'center',
  marginLeft: 'auto',
}));

interface ContextualTipsComponentProps {
  contextType?: string;
  contextId?: string;
  action?: string;
  position?: 'top' | 'bottom' | 'left' | 'right';
  maxTips?: number;
}

const ContextualTipsComponent: React.FC<ContextualTipsComponentProps> = ({
  contextType,
  contextId,
  action,
  position = 'bottom',
  maxTips = 1,
}) => {
  const dispatch = useDispatch<AppDispatch>();
  const location = useLocation();
  const { user } = useAuth();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  
  const tips = useSelector(selectContextualTips);
  const isLoading = useSelector(selectTipsLoading);
  
  const [currentTipIndex, setCurrentTipIndex] = useState(0);
  const [feedbackSubmitted, setFeedbackSubmitted] = useState<Record<number, boolean>>({});
  const [showSnackbar, setShowSnackbar] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState('');
  
  // Fetch contextual tips when component mounts or context changes
  useEffect(() => {
    const fetchTips = async () => {
      const pageUrl = location.pathname;
      
      await dispatch(fetchContextualTips({
        userId: user?.id,
        contextType: contextType || 'page',
        contextId: contextId || pageUrl,
        pageUrl,
        action,
      }));
    };
    
    fetchTips();
    
    // Cleanup
    return () => {
      dispatch(clearContextualTips());
    };
  }, [dispatch, contextType, contextId, action, location.pathname, user?.id]);
  
  // Record tip view when tips change
  useEffect(() => {
    if (tips.length > 0 && user?.id) {
      dispatch(recordTipView(tips[currentTipIndex].id));
    }
  }, [dispatch, tips, currentTipIndex, user?.id]);
  
  const handleDismiss = (tip: Tip) => {
    if (user?.id) {
      dispatch(recordTipDismiss(tip.id));
    }
    
    // Move to next tip or hide if last
    if (currentTipIndex < tips.length - 1) {
      setCurrentTipIndex(currentTipIndex + 1);
    } else {
      dispatch(clearContextualTips());
    }
  };
  
  const handleAction = (tip: Tip) => {
    if (user?.id) {
      dispatch(recordTipClick(tip.id));
    }
    
    // If there's an action URL, navigate to it
    if (tip.actionUrl) {
      window.open(tip.actionUrl, '_blank');
    }
  };
  
  const handleFeedback = (tip: Tip, helpful: boolean) => {
    if (user?.id) {
      dispatch(submitTipFeedback({
        tipId: tip.id,
        helpful,
      }));
      
      setFeedbackSubmitted({
        ...feedbackSubmitted,
        [tip.id]: true,
      });
      
      setSnackbarMessage(helpful ? 'Thank you for your feedback!' : 'We\'ll try to improve our suggestions.');
      setShowSnackbar(true);
    }
  };
  
  const handleCloseSnackbar = () => {
    setShowSnackbar(false);
  };
  
  // Don't render if loading or no tips
  if (isLoading || tips.length === 0 || currentTipIndex >= tips.length) {
    return null;
  }
  
  // Only show up to maxTips
  const visibleTips = tips.slice(0, maxTips);
  const currentTip = tips[currentTipIndex];
  
  // Determine position styles
  let positionStyles = {};
  switch (position) {
    case 'top':
      positionStyles = {
        position: 'fixed',
        top: 16,
        left: '50%',
        transform: 'translateX(-50%)',
        zIndex: 1000,
      };
      break;
    case 'bottom':
      positionStyles = {
        position: 'fixed',
        bottom: 16,
        left: '50%',
        transform: 'translateX(-50%)',
        zIndex: 1000,
      };
      break;
    case 'left':
      positionStyles = {
        position: 'fixed',
        left: 16,
        top: '50%',
        transform: 'translateY(-50%)',
        zIndex: 1000,
      };
      break;
    case 'right':
      positionStyles = {
        position: 'fixed',
        right: 16,
        top: '50%',
        transform: 'translateY(-50%)',
        zIndex: 1000,
      };
      break;
    default:
      positionStyles = {
        position: 'fixed',
        bottom: 16,
        right: 16,
        zIndex: 1000,
      };
  }
  
  return (
    <>
      <Fade in={true} timeout={500}>
        <Box sx={{ ...positionStyles, maxWidth: isMobile ? '90%' : '400px' }}>
          <TipCard>
            <TipIcon>
              <LightbulbIcon fontSize="small" />
            </TipIcon>
            <CardContent sx={{ pt: 3 }}>
              <Typography variant="h6" component="div" gutterBottom>
                {currentTip.title}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {currentTip.content}
              </Typography>
            </CardContent>
            <CardActions sx={{ justifyContent: 'space-between', flexWrap: 'wrap' }}>
              {currentTip.actionText && currentTip.actionUrl && (
                <Button 
                  size="small" 
                  color="primary" 
                  onClick={() => handleAction(currentTip)}
                >
                  {currentTip.actionText}
                </Button>
              )}
              
              <Box sx={{ display: 'flex', alignItems: 'center', ml: 'auto' }}>
                {!feedbackSubmitted[currentTip.id] && (
                  <FeedbackContainer>
                    <Tooltip title="Helpful">
                      <IconButton 
                        size="small" 
                        onClick={() => handleFeedback(currentTip, true)}
                        aria-label="Helpful"
                      >
                        <ThumbUpIcon fontSize="small" />
                      </IconButton>
                    </Tooltip>
                    <Tooltip title="Not helpful">
                      <IconButton 
                        size="small" 
                        onClick={() => handleFeedback(currentTip, false)}
                        aria-label="Not helpful"
                      >
                        <ThumbDownIcon fontSize="small" />
                      </IconButton>
                    </Tooltip>
                  </FeedbackContainer>
                )}
                
                <Tooltip title="Dismiss">
                  <IconButton 
                    size="small" 
                    onClick={() => handleDismiss(currentTip)}
                    aria-label="Dismiss"
                  >
                    <CloseIcon fontSize="small" />
                  </IconButton>
                </Tooltip>
              </Box>
            </CardActions>
          </TipCard>
        </Box>
      </Fade>
      
      <Snackbar
        open={showSnackbar}
        autoHideDuration={3000}
        onClose={handleCloseSnackbar}
        TransitionComponent={Slide}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
      >
        <Alert onClose={handleCloseSnackbar} severity="success" sx={{ width: '100%' }}>
          {snackbarMessage}
        </Alert>
      </Snackbar>
    </>
  );
};

export default ContextualTipsComponent;

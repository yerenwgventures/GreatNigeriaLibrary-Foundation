import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Box,
  Button,
  Card,
  CardActions,
  CardContent,
  CardMedia,
  Chip,
  CircularProgress,
  Grid,
  IconButton,
  Tooltip,
  Typography,
  useTheme,
} from '@mui/material';
import {
  BookmarkAdd as BookmarkAddIcon,
  BookmarkAdded as BookmarkAddedIcon,
  Close as CloseIcon,
  Info as InfoIcon,
} from '@mui/icons-material';
import {
  fetchRecommendations,
  updateRecommendationStatus,
  selectRecommendations,
  selectPersonalizationLoading,
  selectPersonalizationError,
  selectHasMoreRecommendations,
  selectUserHasStyle,
} from '../../features/personalization/personalizationSlice';
import { ContentRecommendation, PersonalizationRequest } from '../../api/personalizationService';
import { AppDispatch } from '../../store';
import { useNavigate } from 'react-router-dom';
import { styled } from '@mui/system';

// Styled components
const RecommendationCard = styled(Card)(({ theme }) => ({
  height: '100%',
  display: 'flex',
  flexDirection: 'column',
  transition: 'transform 0.3s ease, box-shadow 0.3s ease',
  '&:hover': {
    transform: 'translateY(-4px)',
    boxShadow: '0 8px 24px rgba(0, 0, 0, 0.15)',
  },
}));

const RecommendationCardContent = styled(CardContent)({
  flexGrow: 1,
});

const ReasonChip = styled(Chip)(({ theme }) => ({
  margin: theme.spacing(0.5),
  fontSize: '0.7rem',
}));

interface PersonalizedRecommendationsProps {
  contentType?: string;
  topic?: string;
  limit?: number;
  showHeader?: boolean;
}

const PersonalizedRecommendations: React.FC<PersonalizedRecommendationsProps> = ({
  contentType,
  topic,
  limit = 3,
  showHeader = true,
}) => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const theme = useTheme();
  
  const recommendations = useSelector(selectRecommendations);
  const loading = useSelector(selectPersonalizationLoading);
  const error = useSelector(selectPersonalizationError);
  const hasMore = useSelector(selectHasMoreRecommendations);
  const userHasStyle = useSelector(selectUserHasStyle);
  
  const [expanded, setExpanded] = useState(false);
  
  useEffect(() => {
    const request: PersonalizationRequest = {
      contentType,
      topic,
      count: expanded ? 10 : limit,
    };
    
    dispatch(fetchRecommendations(request));
  }, [dispatch, contentType, topic, limit, expanded]);
  
  const handleSaveRecommendation = (rec: ContentRecommendation) => {
    dispatch(
      updateRecommendationStatus({
        recId: rec.id!,
        viewed: true,
        saved: !rec.isSaved,
        rejected: false,
      })
    );
  };
  
  const handleRejectRecommendation = (rec: ContentRecommendation) => {
    dispatch(
      updateRecommendationStatus({
        recId: rec.id!,
        viewed: true,
        saved: false,
        rejected: true,
      })
    );
  };
  
  const handleNavigateToContent = (rec: ContentRecommendation) => {
    // Mark as viewed
    dispatch(
      updateRecommendationStatus({
        recId: rec.id!,
        viewed: true,
        saved: rec.isSaved,
        rejected: rec.isRejected,
      })
    );
    
    // Navigate to the appropriate page based on content type
    switch (rec.contentType.toLowerCase()) {
      case 'book':
        navigate(`/books/${rec.contentId}`);
        break;
      case 'video':
        navigate(`/videos/${rec.contentId}`);
        break;
      case 'course':
        navigate(`/courses/${rec.contentId}`);
        break;
      case 'tutorial':
        navigate(`/tutorials/${rec.contentId}`);
        break;
      default:
        navigate(`/${rec.contentType.toLowerCase()}s/${rec.contentId}`);
    }
  };
  
  const handleTakeAssessment = () => {
    navigate('/learning-style-assessment');
  };
  
  const handleToggleExpand = () => {
    setExpanded(!expanded);
  };
  
  const visibleRecommendations = expanded
    ? recommendations
    : recommendations.slice(0, limit);
  
  if (loading && recommendations.length === 0) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" p={3}>
        <CircularProgress size={24} />
      </Box>
    );
  }
  
  if (error) {
    return (
      <Box p={3}>
        <Typography color="error" variant="body2">
          Error loading recommendations: {error}
        </Typography>
      </Box>
    );
  }
  
  if (!userHasStyle) {
    return (
      <Box p={3} textAlign="center">
        <Typography variant="body1" paragraph>
          Take the learning style assessment to get personalized recommendations.
        </Typography>
        <Button variant="contained" color="primary" onClick={handleTakeAssessment}>
          Take Assessment
        </Button>
      </Box>
    );
  }
  
  if (recommendations.length === 0) {
    return (
      <Box p={3}>
        <Typography variant="body2" color="textSecondary">
          No recommendations available at this time.
        </Typography>
      </Box>
    );
  }
  
  return (
    <Box>
      {showHeader && (
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="h6">Recommended For You</Typography>
          <Tooltip title="Based on your learning style and preferences">
            <IconButton size="small">
              <InfoIcon />
            </IconButton>
          </Tooltip>
        </Box>
      )}
      
      <Grid container spacing={2}>
        {visibleRecommendations.map((rec) => (
          <Grid item xs={12} sm={6} md={4} key={rec.id}>
            <RecommendationCard>
              <CardMedia
                component="img"
                height="140"
                image={`https://source.unsplash.com/random/300x200?${rec.contentType}`}
                alt={rec.title}
              />
              <RecommendationCardContent>
                <Typography variant="h6" gutterBottom>
                  {rec.title}
                </Typography>
                <Typography variant="body2" color="textSecondary" paragraph>
                  {rec.description}
                </Typography>
                <Box display="flex" flexWrap="wrap" mt={1}>
                  <Chip
                    label={rec.contentType}
                    size="small"
                    color="primary"
                    sx={{ mr: 0.5, mb: 0.5 }}
                  />
                  {rec.recommendationScore >= 90 && (
                    <Chip
                      label="Highly Recommended"
                      size="small"
                      color="secondary"
                      sx={{ mr: 0.5, mb: 0.5 }}
                    />
                  )}
                </Box>
                <Box display="flex" flexWrap="wrap" mt={1}>
                  {rec.reasonCodes.map((reason, index) => (
                    <ReasonChip key={index} label={reason} size="small" variant="outlined" />
                  ))}
                </Box>
              </RecommendationCardContent>
              <CardActions>
                <Button size="small" onClick={() => handleNavigateToContent(rec)}>
                  View
                </Button>
                <Tooltip title={rec.isSaved ? 'Remove from saved' : 'Save for later'}>
                  <IconButton
                    size="small"
                    onClick={() => handleSaveRecommendation(rec)}
                    color={rec.isSaved ? 'primary' : 'default'}
                  >
                    {rec.isSaved ? <BookmarkAddedIcon /> : <BookmarkAddIcon />}
                  </IconButton>
                </Tooltip>
                <Tooltip title="Not interested">
                  <IconButton
                    size="small"
                    onClick={() => handleRejectRecommendation(rec)}
                    sx={{ ml: 'auto' }}
                  >
                    <CloseIcon />
                  </IconButton>
                </Tooltip>
              </CardActions>
            </RecommendationCard>
          </Grid>
        ))}
      </Grid>
      
      {hasMore && (
        <Box display="flex" justifyContent="center" mt={3}>
          <Button onClick={handleToggleExpand}>
            {expanded ? 'Show Less' : 'Show More'}
          </Button>
        </Box>
      )}
    </Box>
  );
};

export default PersonalizedRecommendations;

import React from 'react';
import { 
  Card, 
  CardContent, 
  CardActions, 
  Typography, 
  Button, 
  Box, 
  Chip, 
  Skeleton,
  useTheme,
  alpha
} from '@mui/material';
import { 
  LocationOn as LocationIcon,
  Business as BusinessIcon,
  Public as RemoteIcon,
  Work as WorkIcon,
  AccessTime as TimeIcon,
  Bookmark as BookmarkIcon,
  BookmarkBorder as BookmarkBorderIcon
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import { Job } from '../../api/marketplaceService';

interface JobCardProps {
  job: Job;
  loading?: boolean;
  onBookmark?: (job: Job, isBookmarked: boolean) => void;
  isBookmarked?: boolean;
}

const JobCard: React.FC<JobCardProps> = ({
  job,
  loading = false,
  onBookmark,
  isBookmarked = false
}) => {
  const theme = useTheme();

  // Format salary range
  const formatSalaryRange = (min?: number, max?: number, currency?: string, period?: string) => {
    if (!min && !max) return 'Salary not specified';
    
    const formatter = new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency: currency || 'NGN',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    });
    
    let salaryText = '';
    
    if (min && max) {
      salaryText = `${formatter.format(min)} - ${formatter.format(max)}`;
    } else if (min) {
      salaryText = `From ${formatter.format(min)}`;
    } else if (max) {
      salaryText = `Up to ${formatter.format(max)}`;
    }
    
    if (period) {
      switch (period) {
        case 'hourly':
          return `${salaryText}/hr`;
        case 'daily':
          return `${salaryText}/day`;
        case 'weekly':
          return `${salaryText}/week`;
        case 'monthly':
          return `${salaryText}/month`;
        case 'yearly':
          return `${salaryText}/year`;
        default:
          return salaryText;
      }
    }
    
    return salaryText;
  };

  // Format location type
  const formatLocationType = (type: string) => {
    switch (type) {
      case 'remote':
        return 'Remote';
      case 'onsite':
        return 'On-site';
      case 'hybrid':
        return 'Hybrid';
      default:
        return type;
    }
  };

  // Calculate days remaining until deadline
  const getDaysRemaining = (deadline?: string) => {
    if (!deadline) return null;
    
    const deadlineDate = new Date(deadline);
    const today = new Date();
    
    // Reset time to compare just the dates
    today.setHours(0, 0, 0, 0);
    deadlineDate.setHours(0, 0, 0, 0);
    
    const diffTime = deadlineDate.getTime() - today.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays < 0) return 'Expired';
    if (diffDays === 0) return 'Today';
    if (diffDays === 1) return '1 day left';
    return `${diffDays} days left`;
  };

  if (loading) {
    return (
      <Card 
        elevation={2} 
        sx={{ 
          height: '100%', 
          display: 'flex', 
          flexDirection: 'column',
          transition: 'transform 0.3s, box-shadow 0.3s',
          '&:hover': {
            transform: 'translateY(-5px)',
            boxShadow: theme.shadows[8]
          },
          borderRadius: 2,
          overflow: 'hidden'
        }}
      >
        <CardContent sx={{ flexGrow: 1 }}>
          <Skeleton variant="text" height={28} width="80%" animation="wave" />
          <Skeleton variant="text" height={20} width="60%" animation="wave" />
          <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
            <Skeleton variant="circular" width={20} height={20} animation="wave" />
            <Skeleton variant="text" height={20} width="40%" sx={{ ml: 1 }} animation="wave" />
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
            <Skeleton variant="circular" width={20} height={20} animation="wave" />
            <Skeleton variant="text" height={20} width="40%" sx={{ ml: 1 }} animation="wave" />
          </Box>
          <Box sx={{ mt: 1 }}>
            <Skeleton variant="text" height={24} width="30%" animation="wave" />
          </Box>
        </CardContent>
        <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
          <Skeleton variant="rectangular" height={36} width={100} animation="wave" />
          <Skeleton variant="rectangular" height={36} width={36} animation="wave" />
        </CardActions>
      </Card>
    );
  }

  const daysRemaining = job.applicationDeadline ? getDaysRemaining(job.applicationDeadline) : null;
  const isExpired = daysRemaining === 'Expired';

  return (
    <Card 
      elevation={2} 
      sx={{ 
        height: '100%', 
        display: 'flex', 
        flexDirection: 'column',
        transition: 'transform 0.3s, box-shadow 0.3s',
        '&:hover': {
          transform: 'translateY(-5px)',
          boxShadow: theme.shadows[8]
        },
        borderRadius: 2,
        overflow: 'hidden',
        position: 'relative',
        bgcolor: isExpired ? alpha(theme.palette.action.disabledBackground, 0.2) : 'background.paper'
      }}
    >
      {job.status !== 'active' && (
        <Box 
          sx={{ 
            position: 'absolute', 
            top: 0, 
            left: 0, 
            width: '100%', 
            bgcolor: alpha(theme.palette.background.paper, 0.8),
            color: theme.palette.text.primary,
            p: 1,
            zIndex: 1,
            textAlign: 'center',
            fontWeight: 'bold',
            textTransform: 'uppercase'
          }}
        >
          {job.status}
        </Box>
      )}
      
      <CardContent sx={{ flexGrow: 1, pt: 3 }}>
        <Typography 
          variant="h6" 
          component={Link} 
          to={`/marketplace/jobs/${job.id}`}
          sx={{ 
            textDecoration: 'none', 
            color: 'inherit',
            display: '-webkit-box',
            WebkitLineClamp: 2,
            WebkitBoxOrient: 'vertical',
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            height: 48,
            '&:hover': {
              color: theme.palette.primary.main
            }
          }}
        >
          {job.title}
        </Typography>
        
        <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
          <BusinessIcon fontSize="small" color="action" />
          <Typography variant="body2" color="text.secondary" sx={{ ml: 0.5 }}>
            {job.company}
          </Typography>
        </Box>
        
        <Box sx={{ display: 'flex', alignItems: 'center', mt: 0.5 }}>
          {job.locationType === 'remote' ? (
            <>
              <RemoteIcon fontSize="small" color="success" />
              <Typography variant="body2" color="success.main" sx={{ ml: 0.5 }}>
                Remote
              </Typography>
            </>
          ) : (
            <>
              <LocationIcon fontSize="small" color="action" />
              <Typography variant="body2" color="text.secondary" sx={{ ml: 0.5 }}>
                {job.location ? `${job.location} (${formatLocationType(job.locationType)})` : formatLocationType(job.locationType)}
              </Typography>
            </>
          )}
        </Box>
        
        <Box sx={{ display: 'flex', alignItems: 'center', mt: 0.5 }}>
          <WorkIcon fontSize="small" color="action" />
          <Typography variant="body2" color="text.secondary" sx={{ ml: 0.5 }}>
            {formatSalaryRange(job.salaryMin, job.salaryMax, job.salaryCurrency, job.salaryPeriod)}
          </Typography>
        </Box>
        
        {daysRemaining && (
          <Box sx={{ display: 'flex', alignItems: 'center', mt: 0.5 }}>
            <TimeIcon fontSize="small" color={isExpired ? "error" : "action"} />
            <Typography 
              variant="body2" 
              color={isExpired ? "error.main" : "text.secondary"} 
              sx={{ ml: 0.5 }}
            >
              {daysRemaining}
            </Typography>
          </Box>
        )}
        
        {job.category && (
          <Chip 
            label={job.category} 
            size="small" 
            sx={{ mt: 1, mr: 0.5 }} 
            color="primary" 
            variant="outlined" 
          />
        )}
        
        {job.tags && job.tags.length > 0 && (
          <Box sx={{ mt: 1, display: 'flex', flexWrap: 'wrap', gap: 0.5 }}>
            {job.tags.slice(0, 3).map((tag, index) => (
              <Chip 
                key={index}
                label={tag} 
                size="small" 
                color="default" 
                variant="outlined" 
              />
            ))}
            {job.tags.length > 3 && (
              <Chip 
                label={`+${job.tags.length - 3}`} 
                size="small" 
                color="default" 
                variant="outlined" 
              />
            )}
          </Box>
        )}
      </CardContent>
      
      <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
        <Button 
          size="small" 
          variant="contained" 
          color="primary" 
          component={Link}
          to={`/marketplace/jobs/${job.id}`}
          disabled={isExpired}
        >
          View Details
        </Button>
        
        {onBookmark && (
          <Button 
            size="small" 
            color="primary" 
            onClick={() => onBookmark(job, !isBookmarked)}
            sx={{ minWidth: 'auto', p: 1 }}
          >
            {isBookmarked ? <BookmarkIcon color="primary" /> : <BookmarkBorderIcon />}
          </Button>
        )}
      </CardActions>
    </Card>
  );
};

export default JobCard;

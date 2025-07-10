import React from 'react';
import { 
  Card, 
  CardMedia, 
  CardContent, 
  CardActions, 
  Typography, 
  Button, 
  Box, 
  Chip, 
  Rating, 
  Skeleton,
  useTheme,
  alpha
} from '@mui/material';
import { 
  LocationOn as LocationIcon,
  Person as PersonIcon,
  Public as RemoteIcon,
  Bookmark as BookmarkIcon,
  BookmarkBorder as BookmarkBorderIcon
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import { Service } from '../../api/marketplaceService';

interface ServiceCardProps {
  service: Service;
  loading?: boolean;
  onBookmark?: (service: Service, isBookmarked: boolean) => void;
  isBookmarked?: boolean;
}

const ServiceCard: React.FC<ServiceCardProps> = ({
  service,
  loading = false,
  onBookmark,
  isBookmarked = false
}) => {
  const theme = useTheme();

  // Format price with currency
  const formatPrice = (price: number, currency: string, priceType: string) => {
    const formatter = new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency: currency || 'NGN',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    });
    
    const formattedPrice = formatter.format(price);
    
    switch (priceType) {
      case 'hourly':
        return `${formattedPrice}/hr`;
      case 'daily':
        return `${formattedPrice}/day`;
      case 'fixed':
        return formattedPrice;
      case 'custom':
        return 'Custom pricing';
      default:
        return formattedPrice;
    }
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
        <Skeleton variant="rectangular" height={200} animation="wave" />
        <CardContent sx={{ flexGrow: 1 }}>
          <Skeleton variant="text" height={28} width="80%" animation="wave" />
          <Skeleton variant="text" height={20} width="60%" animation="wave" />
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
        position: 'relative'
      }}
    >
      {service.status !== 'active' && (
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
          {service.status}
        </Box>
      )}
      
      <CardMedia
        component="img"
        height="200"
        image={service.mediaUrls?.[0] || '/images/placeholder-service.jpg'}
        alt={service.title}
        sx={{ objectFit: 'cover' }}
      />
      
      <CardContent sx={{ flexGrow: 1 }}>
        <Typography 
          variant="h6" 
          component={Link} 
          to={`/marketplace/services/${service.id}`}
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
          {service.title}
        </Typography>
        
        <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
          <PersonIcon fontSize="small" color="action" />
          <Typography variant="body2" color="text.secondary" sx={{ ml: 0.5 }}>
            {service.provider?.name || 'Unknown Provider'}
          </Typography>
        </Box>
        
        <Box sx={{ display: 'flex', alignItems: 'center', mt: 0.5 }}>
          {service.isRemote ? (
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
                {service.location || 'On-site'}
              </Typography>
            </>
          )}
        </Box>
        
        {(service.rating !== undefined && service.reviewCount !== undefined) && (
          <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
            <Rating 
              value={service.rating} 
              readOnly 
              size="small" 
              precision={0.5} 
            />
            <Typography variant="body2" color="text.secondary" sx={{ ml: 0.5 }}>
              ({service.reviewCount})
            </Typography>
          </Box>
        )}
        
        {service.category && (
          <Chip 
            label={service.category} 
            size="small" 
            sx={{ mt: 1, mr: 0.5 }} 
            color="primary" 
            variant="outlined" 
          />
        )}
        
        {service.subcategory && (
          <Chip 
            label={service.subcategory} 
            size="small" 
            sx={{ mt: 1 }} 
            color="secondary" 
            variant="outlined" 
          />
        )}
      </CardContent>
      
      <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
        <Typography variant="h6" color="primary" fontWeight="bold">
          {formatPrice(service.price, service.currency, service.priceType)}
        </Typography>
        
        <Box>
          {onBookmark && (
            <Button 
              size="small" 
              color="primary" 
              onClick={() => onBookmark(service, !isBookmarked)}
              sx={{ minWidth: 'auto', p: 1 }}
            >
              {isBookmarked ? <BookmarkIcon color="primary" /> : <BookmarkBorderIcon />}
            </Button>
          )}
          
          <Button 
            size="small" 
            variant="contained" 
            color="primary" 
            component={Link}
            to={`/marketplace/services/${service.id}`}
            sx={{ ml: 1 }}
          >
            Details
          </Button>
        </Box>
      </CardActions>
    </Card>
  );
};

export default ServiceCard;

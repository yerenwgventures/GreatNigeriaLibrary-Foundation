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
  ShoppingCart as CartIcon,
  Favorite as FavoriteIcon,
  FavoriteBorder as FavoriteBorderIcon
} from '@mui/icons-material';
import { Link } from 'react-router-dom';
import { Product } from '../../api/marketplaceService';

interface ProductCardProps {
  product: Product;
  loading?: boolean;
  onAddToCart?: (product: Product) => void;
  onToggleFavorite?: (product: Product, isFavorite: boolean) => void;
  isFavorite?: boolean;
}

const ProductCard: React.FC<ProductCardProps> = ({
  product,
  loading = false,
  onAddToCart,
  onToggleFavorite,
  isFavorite = false
}) => {
  const theme = useTheme();

  // Format price with currency
  const formatPrice = (price: number, currency: string) => {
    const formatter = new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency: currency || 'NGN',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    });
    return formatter.format(price);
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
      {product.status !== 'active' && (
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
          {product.status}
        </Box>
      )}
      
      <CardMedia
        component="img"
        height="200"
        image={product.mediaUrls?.[0] || '/images/placeholder-product.jpg'}
        alt={product.title}
        sx={{ objectFit: 'cover' }}
      />
      
      <CardContent sx={{ flexGrow: 1 }}>
        <Typography 
          variant="h6" 
          component={Link} 
          to={`/marketplace/products/${product.id}`}
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
          {product.title}
        </Typography>
        
        <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
          <PersonIcon fontSize="small" color="action" />
          <Typography variant="body2" color="text.secondary" sx={{ ml: 0.5 }}>
            {product.seller?.name || 'Unknown Seller'}
          </Typography>
        </Box>
        
        {product.location && (
          <Box sx={{ display: 'flex', alignItems: 'center', mt: 0.5 }}>
            <LocationIcon fontSize="small" color="action" />
            <Typography variant="body2" color="text.secondary" sx={{ ml: 0.5 }}>
              {product.location}
            </Typography>
          </Box>
        )}
        
        {(product.rating !== undefined && product.reviewCount !== undefined) && (
          <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
            <Rating 
              value={product.rating} 
              readOnly 
              size="small" 
              precision={0.5} 
            />
            <Typography variant="body2" color="text.secondary" sx={{ ml: 0.5 }}>
              ({product.reviewCount})
            </Typography>
          </Box>
        )}
        
        {product.category && (
          <Chip 
            label={product.category} 
            size="small" 
            sx={{ mt: 1, mr: 0.5 }} 
            color="primary" 
            variant="outlined" 
          />
        )}
        
        {product.subcategory && (
          <Chip 
            label={product.subcategory} 
            size="small" 
            sx={{ mt: 1 }} 
            color="secondary" 
            variant="outlined" 
          />
        )}
      </CardContent>
      
      <CardActions sx={{ justifyContent: 'space-between', px: 2, pb: 2 }}>
        <Typography variant="h6" color="primary" fontWeight="bold">
          {formatPrice(product.price, product.currency)}
        </Typography>
        
        <Box>
          {onToggleFavorite && (
            <Button 
              size="small" 
              color="primary" 
              onClick={() => onToggleFavorite(product, !isFavorite)}
              sx={{ minWidth: 'auto', p: 1 }}
            >
              {isFavorite ? <FavoriteIcon color="error" /> : <FavoriteBorderIcon />}
            </Button>
          )}
          
          {onAddToCart && (
            <Button 
              size="small" 
              variant="contained" 
              color="primary" 
              startIcon={<CartIcon />}
              onClick={() => onAddToCart(product)}
              sx={{ ml: 1 }}
            >
              Add
            </Button>
          )}
        </Box>
      </CardActions>
    </Card>
  );
};

export default ProductCard;

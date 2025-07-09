import React, { useState } from 'react';
import {
  Card,
  CardMedia,
  CardContent,
  CardActions,
  Typography,
  Box,
  Button,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  IconButton,
  Tooltip,
  useTheme,
  alpha
} from '@mui/material';
import {
  ContentCopy as CopyIcon,
  Share as ShareIcon,
  Check as CheckIcon,
  MonetizationOn as CommissionIcon,
  AccessTime as TimeIcon,
  Person as SellerIcon
} from '@mui/icons-material';
import { AffiliateProduct } from '../../api/affiliateService';

interface AffiliateProductCardProps {
  product: AffiliateProduct;
  referralCode: string;
  onGetLink: (productId: string, referralCode: string) => string;
}

const AffiliateProductCard: React.FC<AffiliateProductCardProps> = ({
  product,
  referralCode,
  onGetLink
}) => {
  const theme = useTheme();
  const [linkDialogOpen, setLinkDialogOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  
  const affiliateLink = onGetLink(product.id, referralCode);
  
  const handleCopyLink = () => {
    navigator.clipboard.writeText(affiliateLink);
    setCopied(true);
    
    setTimeout(() => {
      setCopied(false);
    }, 2000);
  };
  
  const handleShare = () => {
    if (navigator.share) {
      navigator.share({
        title: product.title,
        text: `Check out this product: ${product.title}`,
        url: affiliateLink
      });
    } else {
      setLinkDialogOpen(true);
    }
  };
  
  // Format currency
  const formatCurrency = (amount: number, currency: string = 'NGN') => {
    return new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  };
  
  // Calculate commission amount
  const calculateCommission = () => {
    return (product.price * product.affiliateSettings.commissionPercentage) / 100;
  };
  
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
      <CardMedia
        component="img"
        height="200"
        image={product.mediaUrls?.[0] || '/images/placeholder-product.jpg'}
        alt={product.title}
      />
      
      <CardContent sx={{ flexGrow: 1 }}>
        <Typography variant="h6" gutterBottom>
          {product.title}
        </Typography>
        
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          {product.description.length > 100 
            ? `${product.description.substring(0, 100)}...` 
            : product.description}
        </Typography>
        
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Typography variant="h6" color="primary" fontWeight="bold">
            {formatCurrency(product.price, product.currency)}
          </Typography>
          
          <Chip 
            icon={<CommissionIcon />} 
            label={`${product.affiliateSettings.commissionPercentage}% Commission`} 
            color="success" 
            size="small" 
          />
        </Box>
        
        <Box sx={{ 
          p: 2, 
          bgcolor: alpha(theme.palette.success.main, 0.1), 
          borderRadius: 2,
          display: 'flex',
          alignItems: 'center',
          mb: 2
        }}>
          <CommissionIcon color="success" sx={{ mr: 1 }} />
          <Typography variant="body2">
            You earn: <strong>{formatCurrency(calculateCommission(), product.currency)}</strong> per sale
          </Typography>
        </Box>
        
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
          <TimeIcon fontSize="small" color="action" sx={{ mr: 1 }} />
          <Typography variant="body2" color="text.secondary">
            {product.affiliateSettings.cookieDurationDays} day cookie duration
          </Typography>
        </Box>
        
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <SellerIcon fontSize="small" color="action" sx={{ mr: 1 }} />
          <Typography variant="body2" color="text.secondary">
            Seller: {product.seller.name || product.seller.username}
          </Typography>
        </Box>
      </CardContent>
      
      <CardActions sx={{ p: 2, borderTop: `1px solid ${theme.palette.divider}` }}>
        <Button
          variant="outlined"
          color="primary"
          startIcon={<CopyIcon />}
          onClick={handleCopyLink}
          fullWidth
          sx={{ mr: 1 }}
        >
          Get Link
        </Button>
        
        <Button
          variant="contained"
          color="primary"
          startIcon={<ShareIcon />}
          onClick={handleShare}
          fullWidth
        >
          Share
        </Button>
      </CardActions>
      
      {/* Affiliate Link Dialog */}
      <Dialog open={linkDialogOpen} onClose={() => setLinkDialogOpen(false)}>
        <DialogTitle>Your Affiliate Link</DialogTitle>
        <DialogContent>
          <Typography variant="body2" paragraph>
            Share this link to earn {product.affiliateSettings.commissionPercentage}% commission on sales.
          </Typography>
          
          <TextField
            fullWidth
            value={affiliateLink}
            InputProps={{
              readOnly: true,
              endAdornment: (
                <Tooltip title={copied ? "Copied!" : "Copy"}>
                  <IconButton onClick={handleCopyLink} edge="end">
                    {copied ? <CheckIcon color="success" /> : <CopyIcon />}
                  </IconButton>
                </Tooltip>
              )
            }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setLinkDialogOpen(false)}>Close</Button>
        </DialogActions>
      </Dialog>
    </Card>
  );
};

export default AffiliateProductCard;

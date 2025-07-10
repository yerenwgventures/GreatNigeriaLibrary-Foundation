import React, { useState } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Switch,
  TextField,
  Button,
  Slider,
  FormControlLabel,
  Divider,
  Alert,
  CircularProgress,
  useTheme,
  alpha
} from '@mui/material';
import {
  MonetizationOn as CommissionIcon,
  AccessTime as TimeIcon,
  Description as TermsIcon
} from '@mui/icons-material';
import { AffiliateProduct, UpdateProductAffiliateSettingsRequest } from '../../api/affiliateService';

interface SellerAffiliateSettingsProps {
  product: AffiliateProduct;
  onSave: (productId: string, settings: UpdateProductAffiliateSettingsRequest) => void;
  loading?: boolean;
}

const SellerAffiliateSettings: React.FC<SellerAffiliateSettingsProps> = ({
  product,
  onSave,
  loading = false
}) => {
  const theme = useTheme();
  
  const [isEnabled, setIsEnabled] = useState(
    product.affiliateSettings?.isAffiliateEnabled || false
  );
  const [commission, setCommission] = useState(
    product.affiliateSettings?.commissionPercentage || 5
  );
  const [cookieDays, setCookieDays] = useState(
    product.affiliateSettings?.cookieDurationDays || 30
  );
  const [terms, setTerms] = useState(
    product.affiliateSettings?.termsAndConditions || ''
  );

  const handleSave = () => {
    onSave(product.id, {
      isAffiliateEnabled: isEnabled,
      commissionPercentage: commission,
      cookieDurationDays: cookieDays,
      termsAndConditions: terms
    });
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
  const calculateCommissionAmount = () => {
    return (product.price * commission) / 100;
  };

  return (
    <Card 
      elevation={2} 
      sx={{ 
        mb: 3, 
        borderRadius: 2,
        position: 'relative',
        overflow: 'hidden'
      }}
    >
      {isEnabled && (
        <Box 
          sx={{ 
            position: 'absolute', 
            top: 0, 
            right: 0, 
            width: 8, 
            height: '100%', 
            bgcolor: theme.palette.success.main
          }} 
        />
      )}
      
      <CardContent sx={{ p: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Box>
            <Typography variant="h6">{product.title}</Typography>
            <Typography variant="body2" color="text.secondary">
              Price: {formatCurrency(product.price, product.currency)}
            </Typography>
          </Box>
          
          <FormControlLabel
            control={
              <Switch
                checked={isEnabled}
                onChange={(e) => setIsEnabled(e.target.checked)}
                color="primary"
              />
            }
            label={isEnabled ? "Affiliate Enabled" : "Affiliate Disabled"}
          />
        </Box>
        
        <Divider sx={{ mb: 3 }} />
        
        <Box sx={{ mb: 3 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
            <CommissionIcon color={isEnabled ? "success" : "disabled"} sx={{ mr: 1 }} />
            <Typography variant="subtitle1" color={isEnabled ? "text.primary" : "text.disabled"}>
              Commission Rate: {commission}%
            </Typography>
          </Box>
          
          <Slider
            value={commission}
            onChange={(_, value) => setCommission(value as number)}
            min={1}
            max={50}
            step={0.5}
            valueLabelDisplay="auto"
            disabled={!isEnabled}
            sx={{ mb: 1 }}
          />
          
          <Box sx={{ 
            p: 2, 
            bgcolor: isEnabled ? alpha(theme.palette.success.main, 0.1) : alpha(theme.palette.grey[500], 0.1), 
            borderRadius: 2,
            display: 'flex',
            alignItems: 'center'
          }}>
            <Typography variant="body2" color={isEnabled ? "text.primary" : "text.disabled"}>
              Affiliates will earn <strong>{formatCurrency(calculateCommissionAmount(), product.currency)}</strong> per sale
            </Typography>
          </Box>
          
          <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
            This is the percentage of the product price that affiliates will earn for each sale.
          </Typography>
        </Box>
        
        <Box sx={{ mb: 3 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
            <TimeIcon color={isEnabled ? "primary" : "disabled"} sx={{ mr: 1 }} />
            <Typography variant="subtitle1" color={isEnabled ? "text.primary" : "text.disabled"}>
              Cookie Duration: {cookieDays} days
            </Typography>
          </Box>
          
          <Slider
            value={cookieDays}
            onChange={(_, value) => setCookieDays(value as number)}
            min={1}
            max={90}
            step={1}
            valueLabelDisplay="auto"
            disabled={!isEnabled}
            sx={{ mb: 1 }}
          />
          
          <Typography variant="body2" color="text.secondary">
            The number of days a referral will be credited to an affiliate after a user clicks their link.
          </Typography>
        </Box>
        
        <Box sx={{ mb: 3 }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
            <TermsIcon color={isEnabled ? "primary" : "disabled"} sx={{ mr: 1 }} />
            <Typography variant="subtitle1" color={isEnabled ? "text.primary" : "text.disabled"}>
              Terms and Conditions (Optional)
            </Typography>
          </Box>
          
          <TextField
            multiline
            rows={4}
            fullWidth
            value={terms}
            onChange={(e) => setTerms(e.target.value)}
            disabled={!isEnabled}
            placeholder="Specify any special terms for your affiliate program..."
            sx={{ mb: 1 }}
          />
          
          <Typography variant="body2" color="text.secondary">
            Add any specific terms or conditions for affiliates promoting this product.
          </Typography>
        </Box>
        
        <Alert severity={isEnabled ? "info" : "warning"} sx={{ mb: 3 }}>
          {isEnabled 
            ? `Affiliates will earn ${commission}% commission (${formatCurrency(calculateCommissionAmount(), product.currency)}) on sales of this product for ${cookieDays} days after a user clicks their link.`
            : "Enable the affiliate program to allow other users to promote your product and earn commissions."
          }
        </Alert>
        
        <Button
          variant="contained"
          color="primary"
          onClick={handleSave}
          disabled={loading}
          fullWidth
        >
          {loading ? <CircularProgress size={24} /> : "Save Settings"}
        </Button>
      </CardContent>
    </Card>
  );
};

export default SellerAffiliateSettings;

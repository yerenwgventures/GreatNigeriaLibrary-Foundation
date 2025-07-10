import React, { useState } from 'react';
import {
  Card,
  CardContent,
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
  CardMembership as MembershipIcon
} from '@mui/icons-material';
import { MembershipPlan } from '../../api/affiliateService';

interface MembershipAffiliateCardProps {
  plan: MembershipPlan;
  referralCode: string;
  onGetLink: (planId: string, referralCode: string) => string;
}

const MembershipAffiliateCard: React.FC<MembershipAffiliateCardProps> = ({
  plan,
  referralCode,
  onGetLink
}) => {
  const theme = useTheme();
  const [linkDialogOpen, setLinkDialogOpen] = useState(false);
  const [copied, setCopied] = useState(false);
  
  const affiliateLink = onGetLink(plan.id, referralCode);
  
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
        title: `Join Great Nigeria - ${plan.name} Plan`,
        text: `Join Great Nigeria with the ${plan.name} membership plan and get access to premium features!`,
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
    return (plan.price * plan.affiliateCommissionPercentage) / 100;
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
        overflow: 'hidden',
        position: 'relative'
      }}
    >
      <Box 
        sx={{ 
          position: 'absolute', 
          top: 0, 
          left: 0, 
          width: '100%', 
          height: 8, 
          bgcolor: theme.palette.primary.main
        }} 
      />
      
      <CardContent sx={{ flexGrow: 1, pt: 4 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <MembershipIcon color="primary" sx={{ mr: 1, fontSize: 28 }} />
            <Typography variant="h5" fontWeight="bold">
              {plan.name}
            </Typography>
          </Box>
          
          <Chip 
            label={plan.isActive ? 'Active' : 'Inactive'} 
            color={plan.isActive ? 'success' : 'default'} 
            size="small" 
          />
        </Box>
        
        <Typography variant="body2" color="text.secondary" sx={{ mb: 3 }}>
          {plan.description || `${plan.name} membership plan for ${plan.durationDays} days.`}
        </Typography>
        
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
          <Typography variant="h5" color="primary" fontWeight="bold">
            {formatCurrency(plan.price)}
          </Typography>
          
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <TimeIcon fontSize="small" color="action" sx={{ mr: 1 }} />
            <Typography variant="body2" color="text.secondary">
              {plan.durationDays} days
            </Typography>
          </Box>
        </Box>
        
        <Box sx={{ 
          p: 2, 
          bgcolor: alpha(theme.palette.success.main, 0.1), 
          borderRadius: 2,
          mb: 3
        }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
            <CommissionIcon color="success" sx={{ mr: 1 }} />
            <Typography variant="subtitle2">
              {plan.affiliateCommissionPercentage}% Commission
            </Typography>
          </Box>
          
          <Typography variant="body2">
            You earn <strong>{formatCurrency(calculateCommission())}</strong> for each new member who joins with your link
          </Typography>
        </Box>
      </CardContent>
      
      <Box sx={{ p: 2, borderTop: `1px solid ${theme.palette.divider}` }}>
        <Box sx={{ display: 'flex', gap: 2 }}>
          <Button
            variant="outlined"
            color="primary"
            startIcon={<CopyIcon />}
            onClick={handleCopyLink}
            fullWidth
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
        </Box>
      </Box>
      
      {/* Affiliate Link Dialog */}
      <Dialog open={linkDialogOpen} onClose={() => setLinkDialogOpen(false)}>
        <DialogTitle>Your Membership Affiliate Link</DialogTitle>
        <DialogContent>
          <Typography variant="body2" paragraph>
            Share this link to earn {plan.affiliateCommissionPercentage}% commission ({formatCurrency(calculateCommission())}) when someone joins with the {plan.name} plan.
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

export default MembershipAffiliateCard;

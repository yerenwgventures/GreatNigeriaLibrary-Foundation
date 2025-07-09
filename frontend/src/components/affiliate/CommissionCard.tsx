import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Chip,
  Button,
  Divider,
  useTheme,
  alpha
} from '@mui/material';
import {
  AccountBalance as BankIcon,
  Person as PersonIcon,
  ShoppingCart as PurchaseIcon,
  CardMembership as SubscriptionIcon,
  PersonAdd as RegistrationIcon,
  Article as ContentIcon,
  ArrowUpward as DirectIcon,
  SubdirectoryArrowRight as IndirectIcon,
  EmojiEvents as BonusIcon
} from '@mui/icons-material';
import { Commission } from '../../api/affiliateService';

interface CommissionCardProps {
  commission: Commission;
  onWithdraw?: (commission: Commission) => void;
}

const CommissionCard: React.FC<CommissionCardProps> = ({
  commission,
  onWithdraw
}) => {
  const theme = useTheme();

  // Format currency
  const formatCurrency = (amount: number, currency: string = 'NGN') => {
    return new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency,
      minimumFractionDigits: 2
    }).format(amount);
  };

  // Format date
  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-NG', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  // Get status chip
  const getStatusChip = (status: string) => {
    switch (status) {
      case 'pending':
        return <Chip size="small" label="Pending" color="warning" />;
      case 'approved':
        return <Chip size="small" label="Approved" color="info" />;
      case 'paid':
        return <Chip size="small" label="Paid" color="success" />;
      case 'rejected':
        return <Chip size="small" label="Rejected" color="error" />;
      default:
        return <Chip size="small" label={status} />;
    }
  };

  // Get source icon
  const getSourceIcon = (source: string) => {
    switch (source) {
      case 'registration':
        return <RegistrationIcon />;
      case 'purchase':
        return <PurchaseIcon />;
      case 'subscription':
        return <SubscriptionIcon />;
      case 'content_sale':
        return <ContentIcon />;
      default:
        return <PurchaseIcon />;
    }
  };

  // Get type icon
  const getTypeIcon = (type: string) => {
    switch (type) {
      case 'direct':
        return <DirectIcon />;
      case 'indirect':
        return <IndirectIcon />;
      case 'bonus':
        return <BonusIcon />;
      default:
        return <DirectIcon />;
    }
  };

  // Get source label
  const getSourceLabel = (source: string) => {
    switch (source) {
      case 'registration':
        return 'User Registration';
      case 'purchase':
        return 'Product Purchase';
      case 'subscription':
        return 'Subscription';
      case 'content_sale':
        return 'Content Sale';
      default:
        return source;
    }
  };

  // Get type label
  const getTypeLabel = (type: string) => {
    switch (type) {
      case 'direct':
        return 'Direct Commission';
      case 'indirect':
        return `Indirect (Tier ${commission.tier})`;
      case 'bonus':
        return 'Bonus Commission';
      default:
        return type;
    }
  };

  // Check if commission can be withdrawn
  const canWithdraw = commission.status === 'approved' && onWithdraw;

  return (
    <Card 
      elevation={2} 
      sx={{ 
        height: '100%', 
        display: 'flex', 
        flexDirection: 'column',
        borderRadius: 2,
        overflow: 'hidden',
        position: 'relative'
      }}
    >
      <Box 
        sx={{ 
          position: 'absolute', 
          top: 0, 
          right: 0, 
          width: 8, 
          height: '100%', 
          bgcolor: 
            commission.status === 'paid' ? theme.palette.success.main :
            commission.status === 'approved' ? theme.palette.info.main :
            commission.status === 'pending' ? theme.palette.warning.main :
            theme.palette.error.main
        }} 
      />
      
      <CardContent sx={{ flexGrow: 1, p: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
          <Typography variant="h5" fontWeight="bold" color="primary">
            {formatCurrency(commission.amount, commission.currency)}
          </Typography>
          
          {getStatusChip(commission.status)}
        </Box>
        
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
          <Box sx={{ 
            mr: 2, 
            p: 1, 
            borderRadius: '50%', 
            bgcolor: alpha(theme.palette.primary.main, 0.1),
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center'
          }}>
            {getSourceIcon(commission.source)}
          </Box>
          
          <Box>
            <Typography variant="subtitle1">
              {getSourceLabel(commission.source)}
            </Typography>
            <Typography variant="body2" color="text.secondary">
              {commission.sourceId || 'N/A'}
            </Typography>
          </Box>
        </Box>
        
        <Divider sx={{ my: 2 }} />
        
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            {getTypeIcon(commission.type)}
            <Typography variant="body2" sx={{ ml: 1 }}>
              {getTypeLabel(commission.type)}
            </Typography>
          </Box>
          
          <Typography variant="body2" color="text.secondary">
            {formatDate(commission.createdAt)}
          </Typography>
        </Box>
        
        {commission.status === 'paid' && commission.paidAt && (
          <Box sx={{ 
            p: 2, 
            bgcolor: alpha(theme.palette.success.main, 0.1), 
            borderRadius: 1,
            display: 'flex',
            alignItems: 'center'
          }}>
            <BankIcon color="success" fontSize="small" sx={{ mr: 1 }} />
            <Typography variant="body2">
              Paid on {formatDate(commission.paidAt)}
            </Typography>
          </Box>
        )}
        
        {commission.status === 'rejected' && commission.rejectionReason && (
          <Box sx={{ 
            p: 2, 
            bgcolor: alpha(theme.palette.error.main, 0.1), 
            borderRadius: 1
          }}>
            <Typography variant="body2" color="error">
              Rejected: {commission.rejectionReason}
            </Typography>
          </Box>
        )}
      </CardContent>
      
      {canWithdraw && (
        <Box sx={{ p: 2, borderTop: `1px solid ${theme.palette.divider}` }}>
          <Button 
            variant="contained" 
            color="primary" 
            fullWidth
            onClick={() => onWithdraw && onWithdraw(commission)}
          >
            Withdraw
          </Button>
        </Box>
      )}
    </Card>
  );
};

export default CommissionCard;

import React from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Grid,
  Divider,
  LinearProgress,
  useTheme,
  alpha
} from '@mui/material';
import {
  TrendingUp as TrendingUpIcon,
  People as PeopleIcon,
  MonetizationOn as MoneyIcon,
  Payments as PaymentsIcon
} from '@mui/icons-material';
import { AffiliateStats } from '../../api/affiliateService';

interface ReferralStatsCardProps {
  stats: AffiliateStats;
  loading?: boolean;
}

const ReferralStatsCard: React.FC<ReferralStatsCardProps> = ({
  stats,
  loading = false
}) => {
  const theme = useTheme();

  // Format currency
  const formatCurrency = (amount: number, currency: string = 'NGN') => {
    return new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount);
  };

  // Format percentage
  const formatPercentage = (value: number) => {
    return `${Math.round(value * 100)}%`;
  };

  return (
    <Card 
      elevation={2} 
      sx={{ 
        borderRadius: 2,
        overflow: 'hidden',
        height: '100%'
      }}
    >
      <CardContent sx={{ p: 3 }}>
        <Typography variant="h6" gutterBottom>
          Affiliate Performance
        </Typography>
        
        {loading ? (
          <Box sx={{ width: '100%', mt: 2 }}>
            <LinearProgress />
          </Box>
        ) : (
          <>
            <Grid container spacing={3}>
              <Grid item xs={12} sm={6}>
                <Box sx={{ mb: 3 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <PeopleIcon color="primary" sx={{ mr: 1 }} />
                    <Typography variant="body2" color="text.secondary">
                      Total Referrals
                    </Typography>
                  </Box>
                  <Typography variant="h4" fontWeight="bold">
                    {stats.totalReferrals}
                  </Typography>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
                    <Typography variant="caption" color="text.secondary">
                      Active: {stats.activeReferrals}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      Pending: {stats.pendingReferrals}
                    </Typography>
                  </Box>
                </Box>
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <Box sx={{ mb: 3 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <TrendingUpIcon color="success" sx={{ mr: 1 }} />
                    <Typography variant="body2" color="text.secondary">
                      Conversion Rate
                    </Typography>
                  </Box>
                  <Typography variant="h4" fontWeight="bold" color="success.main">
                    {formatPercentage(stats.conversionRate)}
                  </Typography>
                  <Box sx={{ 
                    height: 8, 
                    width: '100%', 
                    bgcolor: theme.palette.grey[200],
                    borderRadius: 4,
                    overflow: 'hidden',
                    mt: 1
                  }}>
                    <Box 
                      sx={{ 
                        height: '100%', 
                        width: `${stats.conversionRate * 100}%`,
                        bgcolor: theme.palette.success.main,
                        borderRadius: 4
                      }} 
                    />
                  </Box>
                </Box>
              </Grid>
            </Grid>
            
            <Divider sx={{ my: 2 }} />
            
            <Grid container spacing={3}>
              <Grid item xs={12} sm={6}>
                <Box sx={{ mb: 3 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <MoneyIcon color="primary" sx={{ mr: 1 }} />
                    <Typography variant="body2" color="text.secondary">
                      Total Earnings
                    </Typography>
                  </Box>
                  <Typography variant="h4" fontWeight="bold" color="primary">
                    {formatCurrency(stats.totalEarnings)}
                  </Typography>
                  <Typography variant="caption" color="text.secondary">
                    This Month: {formatCurrency(stats.currentMonthEarnings)}
                  </Typography>
                </Box>
              </Grid>
              
              <Grid item xs={12} sm={6}>
                <Box sx={{ mb: 3 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                    <PaymentsIcon color="info" sx={{ mr: 1 }} />
                    <Typography variant="body2" color="text.secondary">
                      Commission Status
                    </Typography>
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                    <Box>
                      <Typography variant="h5" fontWeight="bold" color="info.main">
                        {stats.pendingCommissions}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        Pending
                      </Typography>
                    </Box>
                    <Box>
                      <Typography variant="h5" fontWeight="bold" color="success.main">
                        {stats.paidCommissions}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        Paid
                      </Typography>
                    </Box>
                    <Box>
                      <Typography variant="h5" fontWeight="bold">
                        {stats.totalCommissions}
                      </Typography>
                      <Typography variant="caption" color="text.secondary">
                        Total
                      </Typography>
                    </Box>
                  </Box>
                </Box>
              </Grid>
            </Grid>
            
            <Box sx={{ 
              p: 2, 
              bgcolor: alpha(theme.palette.primary.main, 0.1), 
              borderRadius: 2,
              mt: 1
            }}>
              <Typography variant="subtitle2" gutterBottom>
                Monthly Performance
              </Typography>
              
              <Grid container spacing={2}>
                <Grid item xs={6}>
                  <Typography variant="caption" color="text.secondary">
                    Referrals by Month
                  </Typography>
                  {stats.referralsByMonth.slice(-3).map((item, index) => (
                    <Box key={index} sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
                      <Typography variant="body2">
                        {item.month}
                      </Typography>
                      <Typography variant="body2" fontWeight="bold">
                        {item.count}
                      </Typography>
                    </Box>
                  ))}
                </Grid>
                
                <Grid item xs={6}>
                  <Typography variant="caption" color="text.secondary">
                    Commissions by Month
                  </Typography>
                  {stats.commissionsByMonth.slice(-3).map((item, index) => (
                    <Box key={index} sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
                      <Typography variant="body2">
                        {item.month}
                      </Typography>
                      <Typography variant="body2" fontWeight="bold">
                        {formatCurrency(item.amount)}
                      </Typography>
                    </Box>
                  ))}
                </Grid>
              </Grid>
            </Box>
          </>
        )}
      </CardContent>
    </Card>
  );
};

export default ReferralStatsCard;

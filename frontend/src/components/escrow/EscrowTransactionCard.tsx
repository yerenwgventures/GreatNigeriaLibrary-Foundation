import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Chip,
  Button,
  Divider,
  Grid,
  CircularProgress,
  Alert,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  useTheme,
  alpha
} from '@mui/material';
import {
  AccessTime as TimeIcon,
  AccountBalance as BankIcon,
  CheckCircle as CheckIcon,
  Cancel as CancelIcon,
  Gavel as DisputeIcon,
  ArrowForward as ReleaseIcon,
  ArrowBack as RefundIcon,
  Info as InfoIcon
} from '@mui/icons-material';
import { EscrowTransaction } from '../../api/escrowService';

interface EscrowTransactionCardProps {
  transaction: EscrowTransaction;
  isBuyer: boolean;
  isSeller: boolean;
  onRelease?: (transaction: EscrowTransaction) => void;
  onRefund?: (transaction: EscrowTransaction) => void;
  onCancel?: (transaction: EscrowTransaction) => void;
  onDispute?: (transaction: EscrowTransaction) => void;
  onViewDetails?: (transaction: EscrowTransaction) => void;
  loading?: boolean;
}

const EscrowTransactionCard: React.FC<EscrowTransactionCardProps> = ({
  transaction,
  isBuyer,
  isSeller,
  onRelease,
  onRefund,
  onCancel,
  onDispute,
  onViewDetails,
  loading = false
}) => {
  const theme = useTheme();
  const [confirmDialogOpen, setConfirmDialogOpen] = useState(false);
  const [confirmAction, setConfirmAction] = useState<'release' | 'refund' | 'cancel' | null>(null);
  const [confirmReason, setConfirmReason] = useState('');

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
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  // Get status chip
  const getStatusChip = (status: string) => {
    switch (status) {
      case 'pending':
        return <Chip size="small" label="Pending" color="warning" />;
      case 'held':
        return <Chip size="small" label="Funds Held" color="info" />;
      case 'released':
        return <Chip size="small" label="Released" color="success" />;
      case 'refunded':
        return <Chip size="small" label="Refunded" color="secondary" />;
      case 'disputed':
        return <Chip size="small" label="Disputed" color="error" />;
      case 'cancelled':
        return <Chip size="small" label="Cancelled" color="default" />;
      default:
        return <Chip size="small" label={status} />;
    }
  };

  // Handle confirm dialog
  const handleOpenConfirmDialog = (action: 'release' | 'refund' | 'cancel') => {
    setConfirmAction(action);
    setConfirmDialogOpen(true);
  };

  const handleCloseConfirmDialog = () => {
    setConfirmDialogOpen(false);
    setConfirmAction(null);
    setConfirmReason('');
  };

  const handleConfirmAction = () => {
    if (!confirmAction) return;

    switch (confirmAction) {
      case 'release':
        onRelease && onRelease(transaction);
        break;
      case 'refund':
        onRefund && onRefund(transaction);
        break;
      case 'cancel':
        onCancel && onCancel(transaction);
        break;
    }

    handleCloseConfirmDialog();
  };

  // Check if actions are available
  const canRelease = isSeller && transaction.status === 'held';
  const canRefund = isSeller && transaction.status === 'held';
  const canCancel = (isBuyer || isSeller) && transaction.status === 'pending';
  const canDispute = isBuyer && (transaction.status === 'held' || transaction.status === 'pending');

  return (
    <Paper 
      elevation={2} 
      sx={{ 
        p: 3, 
        borderRadius: 2,
        position: 'relative',
        overflow: 'hidden'
      }}
    >
      {/* Status indicator */}
      <Box 
        sx={{ 
          position: 'absolute', 
          top: 0, 
          right: 0, 
          width: 8, 
          height: '100%', 
          bgcolor: 
            transaction.status === 'released' ? theme.palette.success.main :
            transaction.status === 'refunded' ? theme.palette.secondary.main :
            transaction.status === 'disputed' ? theme.palette.error.main :
            transaction.status === 'held' ? theme.palette.info.main :
            transaction.status === 'pending' ? theme.palette.warning.main :
            theme.palette.grey[500]
        }} 
      />

      <Grid container spacing={2}>
        <Grid item xs={12} sm={8}>
          <Typography variant="h6" gutterBottom>
            Order #{transaction.orderId}
          </Typography>
          
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
            <TimeIcon fontSize="small" color="action" sx={{ mr: 1 }} />
            <Typography variant="body2" color="text.secondary">
              Created: {formatDate(transaction.createdAt)}
            </Typography>
          </Box>
          
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
            <BankIcon fontSize="small" color="action" sx={{ mr: 1 }} />
            <Typography variant="body2" color="text.secondary">
              Release Condition: {
                transaction.releaseCondition === 'manual' ? 'Manual Release' :
                transaction.releaseCondition === 'auto' ? 'Automatic on Delivery' :
                transaction.releaseCondition === 'timed' ? `Timed (${transaction.releaseAfter ? formatDate(transaction.releaseAfter) : 'N/A'})` :
                transaction.releaseCondition
              }
            </Typography>
          </Box>
          
          {transaction.status === 'released' && transaction.releasedAt && (
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
              <CheckIcon fontSize="small" color="success" sx={{ mr: 1 }} />
              <Typography variant="body2" color="success.main">
                Released: {formatDate(transaction.releasedAt)}
              </Typography>
            </Box>
          )}
          
          {transaction.status === 'refunded' && transaction.refundedAt && (
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
              <RefundIcon fontSize="small" color="secondary" sx={{ mr: 1 }} />
              <Typography variant="body2" color="secondary.main">
                Refunded: {formatDate(transaction.refundedAt)}
              </Typography>
            </Box>
          )}
          
          {transaction.status === 'disputed' && transaction.disputedAt && (
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
              <DisputeIcon fontSize="small" color="error" sx={{ mr: 1 }} />
              <Typography variant="body2" color="error.main">
                Disputed: {formatDate(transaction.disputedAt)}
              </Typography>
            </Box>
          )}
          
          {transaction.status === 'cancelled' && transaction.cancelledAt && (
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
              <CancelIcon fontSize="small" color="action" sx={{ mr: 1 }} />
              <Typography variant="body2" color="text.secondary">
                Cancelled: {formatDate(transaction.cancelledAt)}
              </Typography>
            </Box>
          )}
        </Grid>
        
        <Grid item xs={12} sm={4} sx={{ display: 'flex', flexDirection: 'column', alignItems: { xs: 'flex-start', sm: 'flex-end' } }}>
          <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
            {getStatusChip(transaction.status)}
          </Box>
          
          <Typography variant="h5" color="primary" fontWeight="bold" sx={{ mb: 2 }}>
            {formatCurrency(transaction.amount, transaction.currency)}
          </Typography>
          
          <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
            {isBuyer ? 'Seller' : 'Buyer'} ID: {isBuyer ? transaction.sellerId : transaction.buyerId}
          </Typography>
        </Grid>
      </Grid>
      
      <Divider sx={{ my: 2 }} />
      
      <Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, justifyContent: 'flex-end' }}>
        {onViewDetails && (
          <Button
            variant="outlined"
            color="primary"
            size="small"
            onClick={() => onViewDetails(transaction)}
            startIcon={<InfoIcon />}
          >
            Details
          </Button>
        )}
        
        {canDispute && onDispute && (
          <Button
            variant="outlined"
            color="error"
            size="small"
            onClick={() => onDispute(transaction)}
            startIcon={<DisputeIcon />}
            disabled={loading}
          >
            Open Dispute
          </Button>
        )}
        
        {canCancel && onCancel && (
          <Button
            variant="outlined"
            color="secondary"
            size="small"
            onClick={() => handleOpenConfirmDialog('cancel')}
            startIcon={<CancelIcon />}
            disabled={loading}
          >
            Cancel
          </Button>
        )}
        
        {canRefund && onRefund && (
          <Button
            variant="outlined"
            color="secondary"
            size="small"
            onClick={() => handleOpenConfirmDialog('refund')}
            startIcon={<RefundIcon />}
            disabled={loading}
          >
            Refund
          </Button>
        )}
        
        {canRelease && onRelease && (
          <Button
            variant="contained"
            color="success"
            size="small"
            onClick={() => handleOpenConfirmDialog('release')}
            startIcon={<ReleaseIcon />}
            disabled={loading}
          >
            Release Funds
          </Button>
        )}
        
        {loading && <CircularProgress size={24} />}
      </Box>
      
      {/* Confirmation Dialog */}
      <Dialog open={confirmDialogOpen} onClose={handleCloseConfirmDialog}>
        <DialogTitle>
          {confirmAction === 'release' ? 'Release Funds' : 
           confirmAction === 'refund' ? 'Refund Funds' : 
           'Cancel Transaction'}
        </DialogTitle>
        <DialogContent>
          <Typography variant="body1" paragraph>
            {confirmAction === 'release' ? 
              'Are you sure you want to release the funds to the buyer? This action cannot be undone.' : 
             confirmAction === 'refund' ? 
              'Are you sure you want to refund the funds to the buyer? This action cannot be undone.' : 
              'Are you sure you want to cancel this transaction? This action cannot be undone.'}
          </Typography>
          
          <TextField
            label="Reason (Optional)"
            fullWidth
            multiline
            rows={3}
            value={confirmReason}
            onChange={(e) => setConfirmReason(e.target.value)}
            margin="normal"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseConfirmDialog}>Cancel</Button>
          <Button 
            onClick={handleConfirmAction} 
            color={
              confirmAction === 'release' ? 'success' : 
              confirmAction === 'refund' ? 'secondary' : 
              'error'
            }
            variant="contained"
          >
            Confirm
          </Button>
        </DialogActions>
      </Dialog>
    </Paper>
  );
};

export default EscrowTransactionCard;

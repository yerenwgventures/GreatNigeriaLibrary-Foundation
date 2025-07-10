import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  Box,
  TextField,
  CircularProgress,
  Alert,
  InputAdornment,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  FormHelperText
} from '@mui/material';
import { requestWithdrawal } from '../../features/livestream/livestreamSlice';
import { RootState } from '../../store';

interface WithdrawalModalProps {
  open: boolean;
  onClose: () => void;
  availableAmount: number;
}

// Sample Nigerian banks
const NIGERIAN_BANKS = [
  'Access Bank',
  'Fidelity Bank',
  'First Bank of Nigeria',
  'First City Monument Bank',
  'Guaranty Trust Bank',
  'Polaris Bank',
  'Stanbic IBTC Bank',
  'Sterling Bank',
  'Union Bank of Nigeria',
  'United Bank for Africa',
  'Wema Bank',
  'Zenith Bank'
];

const WithdrawalModal: React.FC<WithdrawalModalProps> = ({ open, onClose, availableAmount }) => {
  const dispatch = useDispatch();
  const { loading, error } = useSelector((state: RootState) => state.livestream.revenue);
  
  const [amount, setAmount] = useState('');
  const [bankName, setBankName] = useState('');
  const [accountNumber, setAccountNumber] = useState('');
  const [accountName, setAccountName] = useState('');
  
  // Form validation
  const [errors, setErrors] = useState({
    amount: '',
    bankName: '',
    accountNumber: '',
    accountName: ''
  });
  
  const validateForm = () => {
    const newErrors = {
      amount: '',
      bankName: '',
      accountNumber: '',
      accountName: ''
    };
    
    let isValid = true;
    
    // Validate amount
    const amountValue = parseFloat(amount);
    if (!amount) {
      newErrors.amount = 'Amount is required';
      isValid = false;
    } else if (isNaN(amountValue)) {
      newErrors.amount = 'Amount must be a number';
      isValid = false;
    } else if (amountValue <= 0) {
      newErrors.amount = 'Amount must be greater than zero';
      isValid = false;
    } else if (amountValue > availableAmount) {
      newErrors.amount = 'Amount exceeds available balance';
      isValid = false;
    }
    
    // Validate bank name
    if (!bankName) {
      newErrors.bankName = 'Bank name is required';
      isValid = false;
    }
    
    // Validate account number
    if (!accountNumber) {
      newErrors.accountNumber = 'Account number is required';
      isValid = false;
    } else if (!/^\d{10}$/.test(accountNumber)) {
      newErrors.accountNumber = 'Account number must be 10 digits';
      isValid = false;
    }
    
    // Validate account name
    if (!accountName) {
      newErrors.accountName = 'Account name is required';
      isValid = false;
    }
    
    setErrors(newErrors);
    return isValid;
  };
  
  const handleSubmit = () => {
    if (!validateForm()) return;
    
    dispatch(requestWithdrawal({
      amount: parseFloat(amount),
      bankName,
      accountNumber,
      accountName
    }) as any).then((result: any) => {
      if (requestWithdrawal.fulfilled.match(result)) {
        // Reset form and close modal on success
        resetForm();
        onClose();
      }
    });
  };
  
  const resetForm = () => {
    setAmount('');
    setBankName('');
    setAccountNumber('');
    setAccountName('');
    setErrors({
      amount: '',
      bankName: '',
      accountNumber: '',
      accountName: ''
    });
  };
  
  const handleClose = () => {
    resetForm();
    onClose();
  };
  
  return (
    <Dialog open={open} onClose={handleClose} maxWidth="sm" fullWidth>
      <DialogTitle>Withdraw Funds</DialogTitle>
      
      <DialogContent dividers>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        
        <Box mb={3}>
          <Typography variant="body1" gutterBottom>
            Available for withdrawal: <strong>₦{availableAmount.toLocaleString()}</strong>
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Withdrawals are processed within 1-3 business days.
          </Typography>
        </Box>
        
        <Box component="form" noValidate>
          <TextField
            label="Amount (₦)"
            fullWidth
            value={amount}
            onChange={(e) => setAmount(e.target.value)}
            margin="normal"
            required
            error={!!errors.amount}
            helperText={errors.amount}
            InputProps={{
              startAdornment: <InputAdornment position="start">₦</InputAdornment>,
            }}
          />
          
          <FormControl fullWidth margin="normal" required error={!!errors.bankName}>
            <InputLabel>Bank</InputLabel>
            <Select
              value={bankName}
              onChange={(e) => setBankName(e.target.value)}
              label="Bank"
            >
              {NIGERIAN_BANKS.map((bank) => (
                <MenuItem key={bank} value={bank}>
                  {bank}
                </MenuItem>
              ))}
            </Select>
            {errors.bankName && <FormHelperText>{errors.bankName}</FormHelperText>}
          </FormControl>
          
          <TextField
            label="Account Number"
            fullWidth
            value={accountNumber}
            onChange={(e) => setAccountNumber(e.target.value)}
            margin="normal"
            required
            error={!!errors.accountNumber}
            helperText={errors.accountNumber}
            inputProps={{ maxLength: 10 }}
          />
          
          <TextField
            label="Account Name"
            fullWidth
            value={accountName}
            onChange={(e) => setAccountName(e.target.value)}
            margin="normal"
            required
            error={!!errors.accountName}
            helperText={errors.accountName}
          />
        </Box>
        
        <Alert severity="info" sx={{ mt: 3 }}>
          Please ensure your bank details are correct. Incorrect details may result in failed transfers or delays.
        </Alert>
      </DialogContent>
      
      <DialogActions>
        <Button onClick={handleClose} disabled={loading}>
          Cancel
        </Button>
        <Button 
          variant="contained" 
          color="primary" 
          onClick={handleSubmit}
          disabled={loading || availableAmount <= 0}
        >
          {loading ? <CircularProgress size={24} /> : 'Request Withdrawal'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default WithdrawalModal;

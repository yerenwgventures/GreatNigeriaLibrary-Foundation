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
  Stepper,
  Step,
  StepLabel,
  Radio,
  RadioGroup,
  FormControlLabel,
  Divider,
  Paper
} from '@mui/material';
import { CoinPackage } from '../../api/livestreamService';
import { purchaseCoins } from '../../features/livestream/livestreamSlice';
import { RootState } from '../../store';

interface PurchaseCoinsModalProps {
  open: boolean;
  onClose: () => void;
  coinPackage: CoinPackage | null;
}

const PurchaseCoinsModal: React.FC<PurchaseCoinsModalProps> = ({ open, onClose, coinPackage }) => {
  const dispatch = useDispatch();
  const { loading, error } = useSelector((state: RootState) => state.livestream.virtualCurrency);
  
  const [activeStep, setActiveStep] = useState(0);
  const [paymentMethod, setPaymentMethod] = useState('card');
  const [cardNumber, setCardNumber] = useState('');
  const [expiryDate, setExpiryDate] = useState('');
  const [cvv, setCvv] = useState('');
  const [cardholderName, setCardholderName] = useState('');
  
  const steps = ['Select Payment Method', 'Enter Payment Details', 'Confirm Purchase'];
  
  const handleNext = () => {
    setActiveStep((prevStep) => prevStep + 1);
  };
  
  const handleBack = () => {
    setActiveStep((prevStep) => prevStep - 1);
  };
  
  const handlePaymentMethodChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setPaymentMethod(event.target.value);
  };
  
  const handlePurchase = () => {
    if (!coinPackage) return;
    
    // Simulate payment processing
    // In a real app, you would integrate with a payment gateway
    const paymentId = Math.floor(Math.random() * 1000000);
    
    dispatch(purchaseCoins({
      packageId: coinPackage.id,
      paymentId
    }) as any).then((result: any) => {
      if (purchaseCoins.fulfilled.match(result)) {
        // Reset form and close modal on success
        setActiveStep(0);
        setPaymentMethod('card');
        setCardNumber('');
        setExpiryDate('');
        setCvv('');
        setCardholderName('');
        onClose();
      }
    });
  };
  
  const isNextDisabled = () => {
    if (activeStep === 0) {
      return !paymentMethod;
    } else if (activeStep === 1) {
      if (paymentMethod === 'card') {
        return !cardNumber || !expiryDate || !cvv || !cardholderName;
      }
      return false;
    }
    return false;
  };
  
  return (
    <Dialog open={open} onClose={onClose} maxWidth="sm" fullWidth>
      <DialogTitle>Purchase Coins</DialogTitle>
      
      <DialogContent dividers>
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        
        <Stepper activeStep={activeStep} sx={{ mb: 4 }}>
          {steps.map((label) => (
            <Step key={label}>
              <StepLabel>{label}</StepLabel>
            </Step>
          ))}
        </Stepper>
        
        {coinPackage && (
          <Box mb={3}>
            <Typography variant="h6" gutterBottom>
              Selected Package
            </Typography>
            <Paper variant="outlined" sx={{ p: 2 }}>
              <Box display="flex" justifyContent="space-between" alignItems="center">
                <Box>
                  <Typography variant="h6" color="primary">
                    {coinPackage.name}
                  </Typography>
                  <Typography variant="body1">
                    {coinPackage.coinsAmount} Coins
                    {coinPackage.bonusCoins > 0 && (
                      <Typography component="span" color="secondary" fontWeight="bold">
                        {' '}+ {coinPackage.bonusCoins} Bonus
                      </Typography>
                    )}
                  </Typography>
                </Box>
                <Typography variant="h5" fontWeight="bold">
                  ₦{coinPackage.priceNaira.toLocaleString()}
                </Typography>
              </Box>
            </Paper>
          </Box>
        )}
        
        {activeStep === 0 && (
          <Box>
            <Typography variant="h6" gutterBottom>
              Select Payment Method
            </Typography>
            <RadioGroup
              value={paymentMethod}
              onChange={handlePaymentMethodChange}
            >
              <FormControlLabel 
                value="card" 
                control={<Radio />} 
                label="Credit/Debit Card" 
              />
              <FormControlLabel 
                value="bank" 
                control={<Radio />} 
                label="Bank Transfer" 
              />
              <FormControlLabel 
                value="ussd" 
                control={<Radio />} 
                label="USSD" 
              />
            </RadioGroup>
          </Box>
        )}
        
        {activeStep === 1 && (
          <Box>
            <Typography variant="h6" gutterBottom>
              Payment Details
            </Typography>
            
            {paymentMethod === 'card' && (
              <Box>
                <TextField
                  label="Card Number"
                  fullWidth
                  value={cardNumber}
                  onChange={(e) => setCardNumber(e.target.value)}
                  margin="normal"
                  required
                  placeholder="1234 5678 9012 3456"
                />
                <Box display="flex" gap={2}>
                  <TextField
                    label="Expiry Date"
                    value={expiryDate}
                    onChange={(e) => setExpiryDate(e.target.value)}
                    margin="normal"
                    required
                    placeholder="MM/YY"
                    sx={{ flex: 1 }}
                  />
                  <TextField
                    label="CVV"
                    value={cvv}
                    onChange={(e) => setCvv(e.target.value)}
                    margin="normal"
                    required
                    placeholder="123"
                    sx={{ flex: 1 }}
                  />
                </Box>
                <TextField
                  label="Cardholder Name"
                  fullWidth
                  value={cardholderName}
                  onChange={(e) => setCardholderName(e.target.value)}
                  margin="normal"
                  required
                />
              </Box>
            )}
            
            {paymentMethod === 'bank' && (
              <Box>
                <Alert severity="info" sx={{ mb: 2 }}>
                  Make a transfer to the account details below and your coins will be credited once payment is confirmed.
                </Alert>
                <Typography variant="body1" gutterBottom>
                  Bank: Great Nigeria Bank
                </Typography>
                <Typography variant="body1" gutterBottom>
                  Account Number: 0123456789
                </Typography>
                <Typography variant="body1" gutterBottom>
                  Account Name: Great Nigeria Library
                </Typography>
                <Typography variant="body1" gutterBottom>
                  Reference: COINS-{Math.floor(Math.random() * 1000000)}
                </Typography>
              </Box>
            )}
            
            {paymentMethod === 'ussd' && (
              <Box>
                <Alert severity="info" sx={{ mb: 2 }}>
                  Dial the USSD code below to complete your payment.
                </Alert>
                <Typography variant="h5" align="center" gutterBottom>
                  *737*000*{Math.floor(Math.random() * 1000000)}#
                </Typography>
                <Typography variant="body2" color="text.secondary" align="center">
                  Follow the prompts to complete your payment
                </Typography>
              </Box>
            )}
          </Box>
        )}
        
        {activeStep === 2 && (
          <Box>
            <Typography variant="h6" gutterBottom>
              Confirm Purchase
            </Typography>
            
            <Box mb={2}>
              <Typography variant="body1" gutterBottom>
                Please review your purchase details:
              </Typography>
              
              <Divider sx={{ my: 2 }} />
              
              <Box display="flex" justifyContent="space-between" mb={1}>
                <Typography variant="body1">Package:</Typography>
                <Typography variant="body1" fontWeight="bold">
                  {coinPackage?.name}
                </Typography>
              </Box>
              
              <Box display="flex" justifyContent="space-between" mb={1}>
                <Typography variant="body1">Coins:</Typography>
                <Typography variant="body1">
                  {coinPackage?.coinsAmount} + {coinPackage?.bonusCoins} Bonus
                </Typography>
              </Box>
              
              <Box display="flex" justifyContent="space-between" mb={1}>
                <Typography variant="body1">Payment Method:</Typography>
                <Typography variant="body1">
                  {paymentMethod === 'card' ? 'Credit/Debit Card' : 
                   paymentMethod === 'bank' ? 'Bank Transfer' : 'USSD'}
                </Typography>
              </Box>
              
              <Divider sx={{ my: 2 }} />
              
              <Box display="flex" justifyContent="space-between">
                <Typography variant="h6">Total:</Typography>
                <Typography variant="h6" color="primary" fontWeight="bold">
                  ₦{coinPackage?.priceNaira.toLocaleString()}
                </Typography>
              </Box>
            </Box>
            
            <Alert severity="info">
              By clicking "Complete Purchase", you agree to our Terms of Service and Privacy Policy.
            </Alert>
          </Box>
        )}
      </DialogContent>
      
      <DialogActions>
        <Button onClick={onClose} disabled={loading}>
          Cancel
        </Button>
        
        {activeStep > 0 && (
          <Button onClick={handleBack} disabled={loading}>
            Back
          </Button>
        )}
        
        {activeStep < steps.length - 1 ? (
          <Button 
            variant="contained" 
            onClick={handleNext}
            disabled={isNextDisabled()}
          >
            Next
          </Button>
        ) : (
          <Button 
            variant="contained" 
            color="primary" 
            onClick={handlePurchase}
            disabled={loading}
          >
            {loading ? <CircularProgress size={24} /> : 'Complete Purchase'}
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
};

export default PurchaseCoinsModal;

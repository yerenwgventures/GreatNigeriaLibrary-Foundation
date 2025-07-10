import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Box,
  Typography,
  Button,
  Grid,
  Card,
  CardContent,
  CardActionArea,
  Avatar,
  TextField,
  FormControlLabel,
  Checkbox,
  Divider,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  CircularProgress,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Alert
} from '@mui/material';
import {
  CardGiftcard as GiftIcon,
  Send as SendIcon,
  Person as PersonIcon,
  Visibility as ViewIcon,
  VisibilityOff as HideIcon
} from '@mui/icons-material';
import { RootState } from '../../store';
import { Gift } from '../../api/livestreamService';
import { sendGift, fetchUserBalance, fetchCoinPackages } from '../../features/livestream/livestreamSlice';
import { useNavigate } from 'react-router-dom';

interface GiftPanelProps {
  gifts: Gift[];
  streamId: number;
  recipientId: number;
  isLive: boolean;
}

// Sample gift options
const GIFT_OPTIONS = [
  { id: 1, name: 'Heart', coins: 10, imageUrl: 'https://via.placeholder.com/50?text=❤️' },
  { id: 2, name: 'Star', coins: 50, imageUrl: 'https://via.placeholder.com/50?text=⭐' },
  { id: 3, name: 'Trophy', coins: 100, imageUrl: 'https://via.placeholder.com/50?text=🏆' },
  { id: 4, name: 'Diamond', coins: 500, imageUrl: 'https://via.placeholder.com/50?text=💎' },
  { id: 5, name: 'Crown', coins: 1000, imageUrl: 'https://via.placeholder.com/50?text=👑' },
  { id: 6, name: 'Rocket', coins: 5000, imageUrl: 'https://via.placeholder.com/50?text=🚀' },
];

const GiftPanel: React.FC<GiftPanelProps> = ({ gifts, streamId, recipientId, isLive }) => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { user } = useSelector((state: RootState) => state.auth);
  const { balance, loading: balanceLoading } = useSelector((state: RootState) => state.livestream.virtualCurrency);
  const { loading: giftLoading, error: giftError } = useSelector((state: RootState) => state.livestream.gifts);
  
  const [selectedGift, setSelectedGift] = useState<typeof GIFT_OPTIONS[0] | null>(null);
  const [message, setMessage] = useState('');
  const [isAnonymous, setIsAnonymous] = useState(false);
  const [showPurchaseDialog, setShowPurchaseDialog] = useState(false);
  
  useEffect(() => {
    if (user) {
      dispatch(fetchUserBalance(user.id) as any);
      dispatch(fetchCoinPackages() as any);
    }
  }, [dispatch, user]);
  
  const handleGiftSelect = (gift: typeof GIFT_OPTIONS[0]) => {
    setSelectedGift(gift);
  };
  
  const handleSendGift = () => {
    if (!user) {
      // Redirect to login
      navigate('/login');
      return;
    }
    
    if (!selectedGift) return;
    
    if (!balance || balance.balance < selectedGift.coins) {
      setShowPurchaseDialog(true);
      return;
    }
    
    const giftData = {
      streamId,
      recipientId,
      giftId: selectedGift.id,
      coinsAmount: selectedGift.coins,
      message: message.trim() || undefined,
      isAnonymous
    };
    
    dispatch(sendGift(giftData) as any).then(() => {
      // Reset form after sending
      setSelectedGift(null);
      setMessage('');
    });
  };
  
  const handlePurchaseCoins = () => {
    setShowPurchaseDialog(false);
    navigate('/coins');
  };
  
  return (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      {/* Gift list */}
      <Box sx={{ flexGrow: 1, overflow: 'auto', mb: 2 }}>
        {gifts.length > 0 ? (
          <List>
            {gifts.map((gift) => (
              <React.Fragment key={gift.id}>
                <ListItem alignItems="flex-start">
                  <ListItemAvatar>
                    <Avatar>
                      {gift.isAnonymous ? <HideIcon /> : <PersonIcon />}
                    </Avatar>
                  </ListItemAvatar>
                  <ListItemText
                    primary={
                      <Typography component="span" variant="body1">
                        {gift.isAnonymous ? 'Anonymous' : `User #${gift.senderId}`} sent a {gift.giftName}
                        {gift.comboCount > 1 && (
                          <Chip 
                            label={`x${gift.comboCount}`} 
                            color="secondary" 
                            size="small" 
                            sx={{ ml: 1 }}
                          />
                        )}
                      </Typography>
                    }
                    secondary={
                      <>
                        <Typography component="span" variant="body2" color="text.primary">
                          {gift.coinsAmount} coins
                        </Typography>
                        {gift.message && (
                          <Typography component="p" variant="body2">
                            "{gift.message}"
                          </Typography>
                        )}
                        <Typography component="span" variant="caption" color="text.secondary">
                          {new Date(gift.createdAt).toLocaleString()}
                        </Typography>
                      </>
                    }
                  />
                </ListItem>
                <Divider variant="inset" component="li" />
              </React.Fragment>
            ))}
          </List>
        ) : (
          <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height="100%">
            <GiftIcon fontSize="large" color="disabled" />
            <Typography variant="body1" color="text.secondary" mt={1}>
              No gifts yet. Be the first to send a gift!
            </Typography>
          </Box>
        )}
      </Box>
      
      {/* Gift selection and sending */}
      {isLive ? (
        <Box>
          {giftError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {giftError}
            </Alert>
          )}
          
          <Typography variant="subtitle2" gutterBottom>
            Select a Gift
          </Typography>
          
          <Grid container spacing={1} sx={{ mb: 2 }}>
            {GIFT_OPTIONS.map((gift) => (
              <Grid item xs={4} key={gift.id}>
                <Card 
                  variant="outlined" 
                  sx={{ 
                    border: selectedGift?.id === gift.id ? 2 : 1,
                    borderColor: selectedGift?.id === gift.id ? 'primary.main' : 'divider'
                  }}
                >
                  <CardActionArea onClick={() => handleGiftSelect(gift)}>
                    <CardContent sx={{ p: 1, textAlign: 'center' }}>
                      <Box 
                        component="img" 
                        src={gift.imageUrl} 
                        alt={gift.name}
                        sx={{ width: 40, height: 40 }}
                      />
                      <Typography variant="caption" display="block">
                        {gift.name}
                      </Typography>
                      <Chip 
                        label={`${gift.coins} coins`} 
                        size="small" 
                        color="primary" 
                        variant="outlined"
                      />
                    </CardContent>
                  </CardActionArea>
                </Card>
              </Grid>
            ))}
          </Grid>
          
          <TextField
            label="Message (optional)"
            fullWidth
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            variant="outlined"
            size="small"
            margin="dense"
            inputProps={{ maxLength: 100 }}
          />
          
          <Box display="flex" justifyContent="space-between" alignItems="center" mt={1}>
            <FormControlLabel
              control={
                <Checkbox
                  checked={isAnonymous}
                  onChange={(e) => setIsAnonymous(e.target.checked)}
                  size="small"
                />
              }
              label={
                <Typography variant="body2">
                  Send anonymously
                </Typography>
              }
            />
            
            <Box display="flex" alignItems="center">
              {user && balance && (
                <Typography variant="body2" color="text.secondary" mr={1}>
                  Balance: {balance.balance} coins
                </Typography>
              )}
              
              <Button
                variant="contained"
                color="primary"
                startIcon={giftLoading ? <CircularProgress size={20} /> : <SendIcon />}
                disabled={!selectedGift || giftLoading || balanceLoading}
                onClick={handleSendGift}
              >
                Send Gift
              </Button>
            </Box>
          </Box>
        </Box>
      ) : (
        <Box textAlign="center" py={2}>
          <Typography variant="body1" color="text.secondary">
            Gifting is only available during live streams
          </Typography>
        </Box>
      )}
      
      {/* Purchase coins dialog */}
      <Dialog open={showPurchaseDialog} onClose={() => setShowPurchaseDialog(false)}>
        <DialogTitle>Insufficient Balance</DialogTitle>
        <DialogContent>
          <Typography variant="body1" gutterBottom>
            You don't have enough coins to send this gift.
          </Typography>
          <Typography variant="body2" color="text.secondary">
            Your current balance: {balance?.balance || 0} coins
          </Typography>
          {selectedGift && (
            <Typography variant="body2" color="text.secondary">
              Required for this gift: {selectedGift.coins} coins
            </Typography>
          )}
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowPurchaseDialog(false)}>
            Cancel
          </Button>
          <Button 
            variant="contained" 
            color="primary" 
            onClick={handlePurchaseCoins}
          >
            Purchase Coins
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default GiftPanel;

import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Box,
  Typography,
  Button,
  Grid,
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
  Alert,
  Tabs,
  Tab,
  IconButton,
  Badge,
  Tooltip,
  Menu,
  MenuItem
} from '@mui/material';
import {
  CardGiftcard as GiftIcon,
  Send as SendIcon,
  Person as PersonIcon,
  Visibility as ViewIcon,
  VisibilityOff as HideIcon,
  Bolt as EffectIcon,
  ShoppingCart as ShopIcon,
  EmojiEmotions as EmojiIcon,
  MoreVert as MoreIcon,
  Favorite as HeartIcon,
  Star as StarIcon,
  Whatshot as FireIcon
} from '@mui/icons-material';
import { RootState } from '../../store';
import { Gift } from '../../api/livestreamService';
import { sendGift, fetchUserBalance, fetchCoinPackages } from '../../features/livestream/livestreamSlice';
import { useNavigate } from 'react-router-dom';
import GiftShop from './GiftShop';
import GiftEffects from './GiftEffects';

interface EnhancedGiftPanelProps {
  gifts: Gift[];
  streamId: number;
  recipientId: number;
  isLive: boolean;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`gift-tabpanel-${index}`}
      aria-labelledby={`gift-tab-${index}`}
      {...other}
      style={{ height: '100%' }}
    >
      {value === index && (
        <Box sx={{ height: '100%' }}>
          {children}
        </Box>
      )}
    </div>
  );
}

const EnhancedGiftPanel: React.FC<EnhancedGiftPanelProps> = ({ gifts, streamId, recipientId, isLive }) => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { user } = useSelector((state: RootState) => state.auth);
  const { balance, loading: balanceLoading } = useSelector((state: RootState) => state.livestream.virtualCurrency);
  const { loading: giftLoading, error: giftError } = useSelector((state: RootState) => state.livestream.gifts);
  
  const [tabValue, setTabValue] = useState(0);
  const [selectedGift, setSelectedGift] = useState<any | null>(null);
  const [selectedEffect, setSelectedEffect] = useState<any | null>(null);
  const [message, setMessage] = useState('');
  const [isAnonymous, setIsAnonymous] = useState(false);
  const [showPurchaseDialog, setShowPurchaseDialog] = useState(false);
  const [isGiftShopOpen, setIsGiftShopOpen] = useState(false);
  const [isEffectsShopOpen, setIsEffectsShopOpen] = useState(false);
  const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
  
  // Recent gifts for quick access
  const [recentGifts, setRecentGifts] = useState<any[]>([
    { id: 1, name: 'Heart', coins: 10, imageUrl: 'https://via.placeholder.com/50?text=❤️', usageCount: 5 },
    { id: 3, name: 'Trophy', coins: 100, imageUrl: 'https://via.placeholder.com/50?text=🏆', usageCount: 2 }
  ]);
  
  // Popular gifts based on platform usage
  const [popularGifts, setPopularGifts] = useState<any[]>([
    { id: 2, name: 'Star', coins: 50, imageUrl: 'https://via.placeholder.com/50?text=⭐', usageCount: 1024 },
    { id: 4, name: 'Diamond', coins: 500, imageUrl: 'https://via.placeholder.com/50?text=💎', usageCount: 512 },
    { id: 5, name: 'Crown', coins: 1000, imageUrl: 'https://via.placeholder.com/50?text=👑', usageCount: 256 }
  ]);
  
  useEffect(() => {
    if (user) {
      dispatch(fetchUserBalance(user.id) as any);
      dispatch(fetchCoinPackages() as any);
    }
  }, [dispatch, user]);
  
  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };
  
  const handleGiftSelect = (gift: any) => {
    setSelectedGift(gift);
    
    // Update recent gifts
    if (!recentGifts.some(g => g.id === gift.id)) {
      setRecentGifts(prev => [gift, ...prev].slice(0, 5));
    } else {
      setRecentGifts(prev => 
        prev.map(g => g.id === gift.id ? { ...g, usageCount: g.usageCount + 1 } : g)
           .sort((a, b) => b.usageCount - a.usageCount)
      );
    }
  };
  
  const handleEffectSelect = (effect: any) => {
    setSelectedEffect(effect);
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
      isAnonymous,
      effectId: selectedEffect?.id
    };
    
    dispatch(sendGift(giftData) as any).then(() => {
      // Reset form after sending
      setMessage('');
    });
  };
  
  const handlePurchaseCoins = () => {
    setShowPurchaseDialog(false);
    navigate('/coins');
  };
  
  const handleOpenGiftShop = () => {
    setIsGiftShopOpen(true);
  };
  
  const handleOpenEffectsShop = () => {
    setIsEffectsShopOpen(true);
  };
  
  const handleMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setMenuAnchorEl(event.currentTarget);
  };
  
  const handleMenuClose = () => {
    setMenuAnchorEl(null);
  };
  
  return (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
        <Tabs
          value={tabValue}
          onChange={handleTabChange}
          variant="fullWidth"
        >
          <Tab 
            icon={<GiftIcon />} 
            label="Gifts" 
            id="gift-tab-0" 
            aria-controls="gift-tabpanel-0" 
          />
          <Tab 
            icon={<HeartIcon />} 
            label="Recent" 
            id="gift-tab-1" 
            aria-controls="gift-tabpanel-1" 
          />
          <Tab 
            icon={<StarIcon />} 
            label="Popular" 
            id="gift-tab-2" 
            aria-controls="gift-tabpanel-2" 
          />
        </Tabs>
      </Box>
      
      {/* Gift list */}
      <Box sx={{ flexGrow: 1, overflow: 'auto', mb: 2 }}>
        <TabPanel value={tabValue} index={0}>
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
        </TabPanel>
        
        <TabPanel value={tabValue} index={1}>
          {recentGifts.length > 0 ? (
            <List>
              {recentGifts.map((gift) => (
                <React.Fragment key={gift.id}>
                  <ListItem 
                    button 
                    onClick={() => handleGiftSelect(gift)}
                    selected={selectedGift?.id === gift.id}
                  >
                    <ListItemAvatar>
                      <Avatar src={gift.imageUrl} alt={gift.name} />
                    </ListItemAvatar>
                    <ListItemText
                      primary={gift.name}
                      secondary={`${gift.coins} coins • Used ${gift.usageCount} times`}
                    />
                    <Chip 
                      label={`${gift.coins} coins`} 
                      size="small" 
                      color="primary" 
                      variant="outlined"
                    />
                  </ListItem>
                  <Divider variant="inset" component="li" />
                </React.Fragment>
              ))}
            </List>
          ) : (
            <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height="100%">
              <HeartIcon fontSize="large" color="disabled" />
              <Typography variant="body1" color="text.secondary" mt={1}>
                No recent gifts. Start sending gifts to build your history!
              </Typography>
            </Box>
          )}
        </TabPanel>
        
        <TabPanel value={tabValue} index={2}>
          {popularGifts.length > 0 ? (
            <List>
              {popularGifts.map((gift) => (
                <React.Fragment key={gift.id}>
                  <ListItem 
                    button 
                    onClick={() => handleGiftSelect(gift)}
                    selected={selectedGift?.id === gift.id}
                  >
                    <ListItemAvatar>
                      <Avatar src={gift.imageUrl} alt={gift.name} />
                    </ListItemAvatar>
                    <ListItemText
                      primary={gift.name}
                      secondary={`${gift.coins} coins • ${gift.usageCount.toLocaleString()} users`}
                    />
                    <Chip 
                      label={`${gift.coins} coins`} 
                      size="small" 
                      color="primary" 
                      variant="outlined"
                    />
                  </ListItem>
                  <Divider variant="inset" component="li" />
                </React.Fragment>
              ))}
              <ListItem button onClick={handleOpenGiftShop}>
                <ListItemAvatar>
                  <Avatar>
                    <ShopIcon />
                  </Avatar>
                </ListItemAvatar>
                <ListItemText
                  primary="Browse Gift Shop"
                  secondary="Discover more gifts and special items"
                />
              </ListItem>
            </List>
          ) : (
            <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height="100%">
              <StarIcon fontSize="large" color="disabled" />
              <Typography variant="body1" color="text.secondary" mt={1}>
                No popular gifts available.
              </Typography>
            </Box>
          )}
        </TabPanel>
      </Box>
      
      {/* Gift selection and sending */}
      {isLive ? (
        <Box>
          {giftError && (
            <Alert severity="error" sx={{ mb: 2 }}>
              {giftError}
            </Alert>
          )}
          
          <Box display="flex" justifyContent="space-between" alignItems="center" mb={1}>
            <Typography variant="subtitle2">
              Selected Gift
            </Typography>
            
            <Box>
              <IconButton size="small" onClick={handleMenuOpen}>
                <MoreIcon />
              </IconButton>
              <Menu
                anchorEl={menuAnchorEl}
                open={Boolean(menuAnchorEl)}
                onClose={handleMenuClose}
              >
                <MenuItem onClick={() => {
                  handleMenuClose();
                  handleOpenGiftShop();
                }}>
                  <ShopIcon fontSize="small" sx={{ mr: 1 }} />
                  Gift Shop
                </MenuItem>
                <MenuItem onClick={() => {
                  handleMenuClose();
                  handleOpenEffectsShop();
                }}>
                  <EffectIcon fontSize="small" sx={{ mr: 1 }} />
                  Effects Shop
                </MenuItem>
              </Menu>
            </Box>
          </Box>
          
          <Box 
            sx={{ 
              p: 2, 
              border: '1px solid', 
              borderColor: 'divider', 
              borderRadius: 1,
              mb: 2,
              display: 'flex',
              alignItems: 'center',
              gap: 2
            }}
          >
            {selectedGift ? (
              <>
                <Badge
                  overlap="circular"
                  anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
                  badgeContent={
                    selectedEffect ? (
                      <Tooltip title={selectedEffect.name}>
                        <Avatar 
                          src={selectedEffect.imageUrl} 
                          sx={{ width: 22, height: 22, border: '2px solid white' }}
                        />
                      </Tooltip>
                    ) : null
                  }
                >
                  <Avatar 
                    src={selectedGift.imageUrl} 
                    alt={selectedGift.name}
                    sx={{ width: 56, height: 56 }}
                  />
                </Badge>
                <Box>
                  <Typography variant="body1" fontWeight="bold">
                    {selectedGift.name}
                    {selectedEffect && (
                      <Chip 
                        label={selectedEffect.name} 
                        size="small" 
                        color="secondary" 
                        sx={{ ml: 1 }}
                      />
                    )}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    {selectedGift.coins} coins
                  </Typography>
                </Box>
                <Box sx={{ ml: 'auto' }}>
                  <Button
                    variant="outlined"
                    size="small"
                    startIcon={<EffectIcon />}
                    onClick={handleOpenEffectsShop}
                    sx={{ mr: 1 }}
                  >
                    {selectedEffect ? 'Change Effect' : 'Add Effect'}
                  </Button>
                </Box>
              </>
            ) : (
              <Box display="flex" flexDirection="column" alignItems="center" width="100%">
                <Typography variant="body2" color="text.secondary">
                  No gift selected
                </Typography>
                <Button
                  variant="outlined"
                  size="small"
                  startIcon={<ShopIcon />}
                  onClick={handleOpenGiftShop}
                  sx={{ mt: 1 }}
                >
                  Browse Gift Shop
                </Button>
              </Box>
            )}
          </Box>
          
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
      
      {/* Gift Shop Dialog */}
      <Dialog 
        open={isGiftShopOpen} 
        onClose={() => setIsGiftShopOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>Gift Shop</DialogTitle>
        <DialogContent dividers sx={{ p: 0 }}>
          <GiftShop 
            onSelectGift={(gift) => {
              handleGiftSelect(gift);
              setIsGiftShopOpen(false);
            }}
            onClose={() => setIsGiftShopOpen(false)}
          />
        </DialogContent>
      </Dialog>
      
      {/* Effects Shop Dialog */}
      <Dialog 
        open={isEffectsShopOpen} 
        onClose={() => setIsEffectsShopOpen(false)}
        maxWidth="md"
        fullWidth
      >
        <DialogTitle>Effects Shop</DialogTitle>
        <DialogContent dividers sx={{ p: 0 }}>
          <GiftEffects 
            onSelectEffect={(effect) => {
              handleEffectSelect(effect);
              setIsEffectsShopOpen(false);
            }}
            onClose={() => setIsEffectsShopOpen(false)}
          />
        </DialogContent>
      </Dialog>
    </Box>
  );
};

export default EnhancedGiftPanel;

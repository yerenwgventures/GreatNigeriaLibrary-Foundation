import React, { useState, useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Box,
  Typography,
  Grid,
  Card,
  CardContent,
  CardMedia,
  CardActionArea,
  Button,
  Chip,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  CircularProgress,
  Tabs,
  Tab,
  Divider,
  Slider,
  IconButton,
  FormControlLabel,
  Checkbox,
  Alert,
  Tooltip
} from '@mui/material';
import {
  CardGiftcard as GiftIcon,
  Favorite as HeartIcon,
  Star as StarIcon,
  EmojiEmotions as EmojiIcon,
  ColorLens as ColorIcon,
  Save as SaveIcon,
  Delete as DeleteIcon,
  Edit as EditIcon,
  Add as AddIcon
} from '@mui/icons-material';
import { RootState } from '../../store';
import { fetchUserBalance } from '../../features/livestream/livestreamSlice';

// Sample gift data
const GIFT_CATEGORIES = [
  { id: 1, name: 'Popular', icon: <StarIcon /> },
  { id: 2, name: 'Animated', icon: <EmojiIcon /> },
  { id: 3, name: 'Special', icon: <HeartIcon /> },
  { id: 4, name: 'Custom', icon: <ColorIcon /> }
];

const SAMPLE_GIFTS = [
  { id: 1, name: 'Heart', category: 1, coins: 10, imageUrl: 'https://via.placeholder.com/100?text=❤️', isCustomizable: false },
  { id: 2, name: 'Star', category: 1, coins: 50, imageUrl: 'https://via.placeholder.com/100?text=⭐', isCustomizable: false },
  { id: 3, name: 'Trophy', category: 1, coins: 100, imageUrl: 'https://via.placeholder.com/100?text=🏆', isCustomizable: false },
  { id: 4, name: 'Diamond', category: 2, coins: 500, imageUrl: 'https://via.placeholder.com/100?text=💎', isCustomizable: true },
  { id: 5, name: 'Crown', category: 2, coins: 1000, imageUrl: 'https://via.placeholder.com/100?text=👑', isCustomizable: true },
  { id: 6, name: 'Rocket', category: 2, coins: 5000, imageUrl: 'https://via.placeholder.com/100?text=🚀', isCustomizable: true },
  { id: 7, name: 'Birthday Cake', category: 3, coins: 2000, imageUrl: 'https://via.placeholder.com/100?text=🎂', isCustomizable: true },
  { id: 8, name: 'Bouquet', category: 3, coins: 3000, imageUrl: 'https://via.placeholder.com/100?text=💐', isCustomizable: true },
  { id: 9, name: 'Custom Text', category: 4, coins: 10000, imageUrl: 'https://via.placeholder.com/100?text=ABC', isCustomizable: true }
];

// Sample user custom gifts
const SAMPLE_CUSTOM_GIFTS = [
  { id: 101, baseGiftId: 4, name: 'My Diamond', coins: 500, imageUrl: 'https://via.placeholder.com/100?text=💎', color: '#ff00ff', scale: 1.2, message: 'Special for you!' },
  { id: 102, baseGiftId: 6, name: 'Super Rocket', coins: 5000, imageUrl: 'https://via.placeholder.com/100?text=🚀', color: '#00ffff', scale: 1.5, message: 'To the moon!' }
];

interface GiftShopProps {
  onSelectGift: (gift: any) => void;
  onClose: () => void;
}

const GiftShop: React.FC<GiftShopProps> = ({ onSelectGift, onClose }) => {
  const dispatch = useDispatch();
  const { user } = useSelector((state: RootState) => state.auth);
  const { balance, loading } = useSelector((state: RootState) => state.livestream.virtualCurrency);
  
  const [selectedCategory, setSelectedCategory] = useState<number>(1);
  const [gifts, setGifts] = useState(SAMPLE_GIFTS);
  const [customGifts, setCustomGifts] = useState(SAMPLE_CUSTOM_GIFTS);
  const [selectedGift, setSelectedGift] = useState<any | null>(null);
  
  // Customization dialog
  const [isCustomizeDialogOpen, setIsCustomizeDialogOpen] = useState(false);
  const [customGiftName, setCustomGiftName] = useState('');
  const [customGiftColor, setCustomGiftColor] = useState('#ff0000');
  const [customGiftScale, setCustomGiftScale] = useState(1);
  const [customGiftMessage, setCustomGiftMessage] = useState('');
  const [isSaving, setIsSaving] = useState(false);
  
  useEffect(() => {
    if (user) {
      dispatch(fetchUserBalance(user.id) as any);
    }
  }, [dispatch, user]);
  
  const handleCategoryChange = (event: React.SyntheticEvent, newValue: number) => {
    setSelectedCategory(newValue);
  };
  
  const handleGiftSelect = (gift: any) => {
    setSelectedGift(gift);
    
    // If gift is not customizable, select it directly
    if (!gift.isCustomizable) {
      onSelectGift(gift);
      onClose();
    }
  };
  
  const handleCustomGiftSelect = (gift: any) => {
    onSelectGift(gift);
    onClose();
  };
  
  const handleCustomizeGift = () => {
    if (!selectedGift) return;
    
    setCustomGiftName(`Custom ${selectedGift.name}`);
    setCustomGiftColor('#ff0000');
    setCustomGiftScale(1);
    setCustomGiftMessage('');
    setIsCustomizeDialogOpen(true);
  };
  
  const handleSaveCustomGift = () => {
    if (!selectedGift) return;
    
    setIsSaving(true);
    
    // Simulate API call to save custom gift
    setTimeout(() => {
      const newCustomGift = {
        id: Date.now(),
        baseGiftId: selectedGift.id,
        name: customGiftName,
        coins: selectedGift.coins,
        imageUrl: selectedGift.imageUrl,
        color: customGiftColor,
        scale: customGiftScale,
        message: customGiftMessage
      };
      
      setCustomGifts([newCustomGift, ...customGifts]);
      setIsSaving(false);
      setIsCustomizeDialogOpen(false);
      
      // Select the newly created custom gift
      onSelectGift(newCustomGift);
      onClose();
    }, 1000);
  };
  
  const handleDeleteCustomGift = (giftId: number) => {
    setCustomGifts(customGifts.filter(gift => gift.id !== giftId));
  };
  
  // Filter gifts by selected category
  const filteredGifts = gifts.filter(gift => 
    selectedCategory === 4 ? gift.category === 4 : gift.category === selectedCategory
  );
  
  // Filter custom gifts by selected category
  const filteredCustomGifts = selectedCategory === 4 
    ? customGifts 
    : [];
  
  return (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        <Tabs
          value={selectedCategory}
          onChange={handleCategoryChange}
          variant="scrollable"
          scrollButtons="auto"
        >
          {GIFT_CATEGORIES.map(category => (
            <Tab 
              key={category.id} 
              label={category.name} 
              icon={category.icon} 
              value={category.id} 
            />
          ))}
        </Tabs>
      </Box>
      
      {loading ? (
        <Box display="flex" justifyContent="center" my={4}>
          <CircularProgress />
        </Box>
      ) : (
        <>
          {balance && (
            <Box display="flex" justifyContent="space-between" alignItems="center" mb={2} px={2}>
              <Typography variant="body1">
                Your Balance: <strong>{balance.balance} coins</strong>
              </Typography>
              
              <Button 
                variant="outlined" 
                size="small" 
                onClick={() => window.location.href = '/coins'}
              >
                Buy Coins
              </Button>
            </Box>
          )}
          
          {selectedCategory === 4 && customGifts.length > 0 && (
            <>
              <Typography variant="h6" gutterBottom sx={{ px: 2 }}>
                Your Custom Gifts
              </Typography>
              
              <Grid container spacing={2} sx={{ px: 2, mb: 3 }}>
                {filteredCustomGifts.map(gift => (
                  <Grid item xs={6} sm={4} key={gift.id}>
                    <Card 
                      variant="outlined" 
                      sx={{ 
                        position: 'relative',
                        height: '100%',
                        display: 'flex',
                        flexDirection: 'column'
                      }}
                    >
                      <CardActionArea 
                        onClick={() => handleCustomGiftSelect(gift)}
                        sx={{ flexGrow: 1 }}
                      >
                        <Box 
                          sx={{ 
                            position: 'relative',
                            pt: '100%',
                            overflow: 'hidden'
                          }}
                        >
                          <Box
                            component="img"
                            src={gift.imageUrl}
                            alt={gift.name}
                            sx={{
                              position: 'absolute',
                              top: '50%',
                              left: '50%',
                              transform: `translate(-50%, -50%) scale(${gift.scale})`,
                              filter: `drop-shadow(0 0 5px ${gift.color})`,
                              maxWidth: '80%',
                              maxHeight: '80%'
                            }}
                          />
                        </Box>
                        
                        <CardContent sx={{ p: 1 }}>
                          <Typography variant="body2" noWrap>
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
                      
                      <Box 
                        sx={{ 
                          position: 'absolute', 
                          top: 5, 
                          right: 5,
                          bgcolor: 'background.paper',
                          borderRadius: '50%'
                        }}
                      >
                        <IconButton 
                          size="small" 
                          onClick={() => handleDeleteCustomGift(gift.id)}
                        >
                          <DeleteIcon fontSize="small" />
                        </IconButton>
                      </Box>
                    </Card>
                  </Grid>
                ))}
              </Grid>
              
              <Divider sx={{ mb: 2 }} />
            </>
          )}
          
          <Typography variant="h6" gutterBottom sx={{ px: 2 }}>
            {selectedCategory === 4 ? 'Create Custom Gift' : 'Available Gifts'}
          </Typography>
          
          <Grid container spacing={2} sx={{ px: 2 }}>
            {filteredGifts.map(gift => (
              <Grid item xs={6} sm={4} key={gift.id}>
                <Card 
                  variant="outlined" 
                  sx={{ 
                    height: '100%',
                    display: 'flex',
                    flexDirection: 'column',
                    border: selectedGift?.id === gift.id ? 2 : 1,
                    borderColor: selectedGift?.id === gift.id ? 'primary.main' : 'divider'
                  }}
                >
                  <CardActionArea 
                    onClick={() => handleGiftSelect(gift)}
                    sx={{ flexGrow: 1 }}
                  >
                    <Box 
                      sx={{ 
                        position: 'relative',
                        pt: '100%',
                        overflow: 'hidden'
                      }}
                    >
                      <Box
                        component="img"
                        src={gift.imageUrl}
                        alt={gift.name}
                        sx={{
                          position: 'absolute',
                          top: '50%',
                          left: '50%',
                          transform: 'translate(-50%, -50%)',
                          maxWidth: '80%',
                          maxHeight: '80%'
                        }}
                      />
                    </Box>
                    
                    <CardContent sx={{ p: 1 }}>
                      <Typography variant="body2" noWrap>
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
                  
                  {gift.isCustomizable && (
                    <Box sx={{ p: 1, pt: 0 }}>
                      <Button
                        variant="outlined"
                        size="small"
                        fullWidth
                        startIcon={<EditIcon />}
                        onClick={handleCustomizeGift}
                        disabled={!selectedGift || selectedGift.id !== gift.id}
                      >
                        Customize
                      </Button>
                    </Box>
                  )}
                </Card>
              </Grid>
            ))}
          </Grid>
        </>
      )}
      
      {/* Customize Gift Dialog */}
      <Dialog 
        open={isCustomizeDialogOpen} 
        onClose={() => setIsCustomizeDialogOpen(false)}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>Customize Your Gift</DialogTitle>
        
        <DialogContent dividers>
          <Grid container spacing={3}>
            <Grid item xs={12} sm={6}>
              <Box 
                sx={{ 
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  justifyContent: 'center',
                  height: '100%'
                }}
              >
                <Box 
                  sx={{ 
                    position: 'relative',
                    width: 150,
                    height: 150,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    mb: 2
                  }}
                >
                  {selectedGift && (
                    <Box
                      component="img"
                      src={selectedGift.imageUrl}
                      alt={selectedGift.name}
                      sx={{
                        transform: `scale(${customGiftScale})`,
                        filter: `drop-shadow(0 0 5px ${customGiftColor})`,
                        maxWidth: '80%',
                        maxHeight: '80%'
                      }}
                    />
                  )}
                </Box>
                
                <Typography variant="body2" color="text.secondary" align="center">
                  Preview
                </Typography>
                
                {selectedGift && (
                  <Chip 
                    label={`${selectedGift.coins} coins`} 
                    size="small" 
                    color="primary" 
                    sx={{ mt: 1 }}
                  />
                )}
              </Box>
            </Grid>
            
            <Grid item xs={12} sm={6}>
              <TextField
                label="Gift Name"
                fullWidth
                value={customGiftName}
                onChange={(e) => setCustomGiftName(e.target.value)}
                margin="normal"
                required
              />
              
              <Box sx={{ mt: 2 }}>
                <Typography gutterBottom>Color</Typography>
                <Box 
                  sx={{ 
                    display: 'flex',
                    alignItems: 'center',
                    gap: 2
                  }}
                >
                  <Box
                    sx={{
                      width: 36,
                      height: 36,
                      bgcolor: customGiftColor,
                      borderRadius: 1,
                      border: '1px solid #ccc'
                    }}
                  />
                  <TextField
                    fullWidth
                    value={customGiftColor}
                    onChange={(e) => setCustomGiftColor(e.target.value)}
                    size="small"
                    placeholder="#ff0000"
                  />
                </Box>
              </Box>
              
              <Box sx={{ mt: 2 }}>
                <Typography gutterBottom>Size</Typography>
                <Slider
                  value={customGiftScale}
                  onChange={(_e, value) => setCustomGiftScale(value as number)}
                  min={0.5}
                  max={2}
                  step={0.1}
                  marks={[
                    { value: 0.5, label: 'Small' },
                    { value: 1, label: 'Normal' },
                    { value: 2, label: 'Large' }
                  ]}
                  valueLabelDisplay="auto"
                />
              </Box>
              
              <TextField
                label="Custom Message (optional)"
                fullWidth
                value={customGiftMessage}
                onChange={(e) => setCustomGiftMessage(e.target.value)}
                margin="normal"
                multiline
                rows={2}
                placeholder="Add a personal message to display with your gift"
              />
            </Grid>
          </Grid>
        </DialogContent>
        
        <DialogActions>
          <Button 
            onClick={() => setIsCustomizeDialogOpen(false)}
            disabled={isSaving}
          >
            Cancel
          </Button>
          <Button 
            variant="contained" 
            color="primary" 
            onClick={handleSaveCustomGift}
            disabled={isSaving || !customGiftName}
            startIcon={isSaving ? <CircularProgress size={20} /> : <SaveIcon />}
          >
            Save Custom Gift
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default GiftShop;

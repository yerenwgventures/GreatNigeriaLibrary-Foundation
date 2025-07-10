import React, { useState } from 'react';
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
  CircularProgress,
  Tabs,
  Tab,
  Divider,
  IconButton,
  Tooltip,
  Alert
} from '@mui/material';
import {
  Bolt as EffectIcon,
  Stars as StarsIcon,
  Whatshot as FireIcon,
  Water as WaterIcon,
  Flare as SparkleIcon,
  Visibility as PreviewIcon,
  ShoppingCart as CartIcon,
  Check as CheckIcon
} from '@mui/icons-material';
import { motion } from 'framer-motion';

// Sample effect categories
const EFFECT_CATEGORIES = [
  { id: 1, name: 'Basic', icon: <StarsIcon /> },
  { id: 2, name: 'Premium', icon: <FireIcon /> },
  { id: 3, name: 'Exclusive', icon: <SparkleIcon /> },
  { id: 4, name: 'Seasonal', icon: <WaterIcon /> }
];

// Sample effects
const SAMPLE_EFFECTS = [
  { 
    id: 1, 
    name: 'Sparkle', 
    category: 1, 
    coins: 500, 
    imageUrl: 'https://via.placeholder.com/100?text=✨', 
    description: 'Adds sparkling stars around your gift',
    owned: true,
    animation: 'sparkle'
  },
  { 
    id: 2, 
    name: 'Explosion', 
    category: 1, 
    coins: 1000, 
    imageUrl: 'https://via.placeholder.com/100?text=💥', 
    description: 'Creates a colorful explosion effect',
    owned: false,
    animation: 'explosion'
  },
  { 
    id: 3, 
    name: 'Rainbow', 
    category: 2, 
    coins: 2000, 
    imageUrl: 'https://via.placeholder.com/100?text=🌈', 
    description: 'Surrounds your gift with a rainbow trail',
    owned: true,
    animation: 'rainbow'
  },
  { 
    id: 4, 
    name: 'Fire', 
    category: 2, 
    coins: 3000, 
    imageUrl: 'https://via.placeholder.com/100?text=🔥', 
    description: 'Adds flames that burn brightly',
    owned: false,
    animation: 'fire'
  },
  { 
    id: 5, 
    name: 'Diamond Dust', 
    category: 3, 
    coins: 5000, 
    imageUrl: 'https://via.placeholder.com/100?text=💎', 
    description: 'Premium sparkling diamond particles',
    owned: false,
    animation: 'diamond'
  },
  { 
    id: 6, 
    name: 'Spotlight', 
    category: 3, 
    coins: 7500, 
    imageUrl: 'https://via.placeholder.com/100?text=🔦', 
    description: 'Highlights your gift with a spotlight effect',
    owned: false,
    animation: 'spotlight'
  },
  { 
    id: 7, 
    name: 'Christmas', 
    category: 4, 
    coins: 2500, 
    imageUrl: 'https://via.placeholder.com/100?text=🎄', 
    description: 'Festive Christmas theme with snow and ornaments',
    owned: false,
    animation: 'christmas'
  },
  { 
    id: 8, 
    name: 'Birthday', 
    category: 4, 
    coins: 2500, 
    imageUrl: 'https://via.placeholder.com/100?text=🎂', 
    description: 'Celebration theme with confetti and balloons',
    owned: true,
    animation: 'birthday'
  }
];

interface GiftEffectsProps {
  onSelectEffect: (effect: any) => void;
  onClose: () => void;
}

const GiftEffects: React.FC<GiftEffectsProps> = ({ onSelectEffect, onClose }) => {
  const [selectedCategory, setSelectedCategory] = useState<number>(1);
  const [effects, setEffects] = useState(SAMPLE_EFFECTS);
  const [selectedEffect, setSelectedEffect] = useState<any | null>(null);
  const [isPreviewDialogOpen, setIsPreviewDialogOpen] = useState(false);
  const [isPurchasing, setIsPurchasing] = useState(false);
  const [purchaseSuccess, setPurchaseSuccess] = useState(false);
  
  const handleCategoryChange = (event: React.SyntheticEvent, newValue: number) => {
    setSelectedCategory(newValue);
  };
  
  const handleEffectSelect = (effect: any) => {
    setSelectedEffect(effect);
    
    // If effect is owned, select it directly
    if (effect.owned) {
      onSelectEffect(effect);
      onClose();
    }
  };
  
  const handlePreviewEffect = () => {
    if (!selectedEffect) return;
    setIsPreviewDialogOpen(true);
  };
  
  const handlePurchaseEffect = () => {
    if (!selectedEffect) return;
    
    setIsPurchasing(true);
    
    // Simulate API call to purchase effect
    setTimeout(() => {
      // Update the effect to be owned
      setEffects(effects.map(effect => 
        effect.id === selectedEffect.id 
          ? { ...effect, owned: true } 
          : effect
      ));
      
      setIsPurchasing(false);
      setPurchaseSuccess(true);
      
      // Reset success message after 3 seconds
      setTimeout(() => {
        setPurchaseSuccess(false);
        
        // Select the newly purchased effect
        onSelectEffect({ ...selectedEffect, owned: true });
        onClose();
      }, 2000);
    }, 1500);
  };
  
  // Filter effects by selected category
  const filteredEffects = effects.filter(effect => effect.category === selectedCategory);
  
  // Get animation component based on effect type
  const getAnimationPreview = (animationType: string) => {
    switch (animationType) {
      case 'sparkle':
        return (
          <Box sx={{ position: 'relative', width: 200, height: 200 }}>
            <Box
              component="img"
              src="https://via.placeholder.com/100?text=✨"
              alt="Gift"
              sx={{
                position: 'absolute',
                top: '50%',
                left: '50%',
                transform: 'translate(-50%, -50%)',
                zIndex: 1
              }}
            />
            {[...Array(10)].map((_, i) => (
              <motion.div
                key={i}
                style={{
                  position: 'absolute',
                  width: 10,
                  height: 10,
                  borderRadius: '50%',
                  backgroundColor: '#ffeb3b',
                  top: '50%',
                  left: '50%',
                  zIndex: 0
                }}
                animate={{
                  x: Math.random() * 200 - 100,
                  y: Math.random() * 200 - 100,
                  opacity: [1, 0],
                  scale: [0, 1, 0]
                }}
                transition={{
                  duration: 2,
                  repeat: Infinity,
                  delay: i * 0.2
                }}
              />
            ))}
          </Box>
        );
      case 'explosion':
        return (
          <Box sx={{ position: 'relative', width: 200, height: 200 }}>
            <Box
              component="img"
              src="https://via.placeholder.com/100?text=💥"
              alt="Gift"
              sx={{
                position: 'absolute',
                top: '50%',
                left: '50%',
                transform: 'translate(-50%, -50%)',
                zIndex: 1
              }}
            />
            <motion.div
              style={{
                position: 'absolute',
                width: 100,
                height: 100,
                borderRadius: '50%',
                background: 'radial-gradient(circle, rgba(255,0,0,0.5) 0%, rgba(255,255,0,0) 70%)',
                top: '50%',
                left: '50%',
                transform: 'translate(-50%, -50%)',
                zIndex: 0
              }}
              animate={{
                scale: [0, 3],
                opacity: [1, 0]
              }}
              transition={{
                duration: 1.5,
                repeat: Infinity,
                repeatDelay: 0.5
              }}
            />
          </Box>
        );
      case 'rainbow':
        return (
          <Box sx={{ position: 'relative', width: 200, height: 200 }}>
            <Box
              component="img"
              src="https://via.placeholder.com/100?text=🌈"
              alt="Gift"
              sx={{
                position: 'absolute',
                top: '50%',
                left: '50%',
                transform: 'translate(-50%, -50%)',
                zIndex: 1
              }}
            />
            {['#ff0000', '#ff9900', '#ffff00', '#00ff00', '#0099ff', '#6633ff'].map((color, i) => (
              <motion.div
                key={i}
                style={{
                  position: 'absolute',
                  width: 120 + i * 20,
                  height: 120 + i * 20,
                  borderRadius: '50%',
                  border: `2px solid ${color}`,
                  top: '50%',
                  left: '50%',
                  transform: 'translate(-50%, -50%)',
                  zIndex: 0
                }}
                animate={{
                  scale: [0.8, 1.2, 0.8],
                  opacity: [0.7, 1, 0.7]
                }}
                transition={{
                  duration: 3,
                  repeat: Infinity,
                  delay: i * 0.2
                }}
              />
            ))}
          </Box>
        );
      case 'fire':
        return (
          <Box sx={{ position: 'relative', width: 200, height: 200 }}>
            <Box
              component="img"
              src="https://via.placeholder.com/100?text=🔥"
              alt="Gift"
              sx={{
                position: 'absolute',
                top: '50%',
                left: '50%',
                transform: 'translate(-50%, -50%)',
                zIndex: 1
              }}
            />
            {[...Array(15)].map((_, i) => (
              <motion.div
                key={i}
                style={{
                  position: 'absolute',
                  width: 10,
                  height: 20,
                  borderRadius: '50%',
                  background: 'linear-gradient(to top, #ff9900, #ff0000)',
                  bottom: '30%',
                  left: `${40 + i * 3}%`,
                  zIndex: 0
                }}
                animate={{
                  y: [-20, -60],
                  opacity: [1, 0],
                  scale: [1, 0.5]
                }}
                transition={{
                  duration: 1 + Math.random(),
                  repeat: Infinity,
                  repeatDelay: Math.random() * 0.5
                }}
              />
            ))}
          </Box>
        );
      default:
        return (
          <Box sx={{ position: 'relative', width: 200, height: 200 }}>
            <Box
              component="img"
              src={selectedEffect?.imageUrl || "https://via.placeholder.com/100"}
              alt="Gift"
              sx={{
                position: 'absolute',
                top: '50%',
                left: '50%',
                transform: 'translate(-50%, -50%)'
              }}
            />
          </Box>
        );
    }
  };
  
  return (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <Box sx={{ borderBottom: 1, borderColor: 'divider', mb: 2 }}>
        <Tabs
          value={selectedCategory}
          onChange={handleCategoryChange}
          variant="scrollable"
          scrollButtons="auto"
        >
          {EFFECT_CATEGORIES.map(category => (
            <Tab 
              key={category.id} 
              label={category.name} 
              icon={category.icon} 
              value={category.id} 
            />
          ))}
        </Tabs>
      </Box>
      
      {purchaseSuccess && (
        <Alert severity="success" sx={{ mx: 2, mb: 2 }}>
          Effect purchased successfully!
        </Alert>
      )}
      
      <Typography variant="h6" gutterBottom sx={{ px: 2 }}>
        {EFFECT_CATEGORIES.find(c => c.id === selectedCategory)?.name} Effects
      </Typography>
      
      <Grid container spacing={2} sx={{ px: 2, flexGrow: 1, overflow: 'auto' }}>
        {filteredEffects.map(effect => (
          <Grid item xs={6} sm={4} key={effect.id}>
            <Card 
              variant="outlined" 
              sx={{ 
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                border: selectedEffect?.id === effect.id ? 2 : 1,
                borderColor: selectedEffect?.id === effect.id ? 'primary.main' : 'divider'
              }}
            >
              <CardActionArea 
                onClick={() => handleEffectSelect(effect)}
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
                    src={effect.imageUrl}
                    alt={effect.name}
                    sx={{
                      position: 'absolute',
                      top: '50%',
                      left: '50%',
                      transform: 'translate(-50%, -50%)',
                      maxWidth: '80%',
                      maxHeight: '80%'
                    }}
                  />
                  
                  {effect.owned && (
                    <Box
                      sx={{
                        position: 'absolute',
                        top: 8,
                        right: 8,
                        bgcolor: 'success.main',
                        color: 'white',
                        borderRadius: '50%',
                        width: 24,
                        height: 24,
                        display: 'flex',
                        alignItems: 'center',
                        justifyContent: 'center'
                      }}
                    >
                      <CheckIcon fontSize="small" />
                    </Box>
                  )}
                </Box>
                
                <CardContent sx={{ p: 1 }}>
                  <Typography variant="body2" noWrap>
                    {effect.name}
                  </Typography>
                  <Typography variant="caption" color="text.secondary" display="block" noWrap>
                    {effect.description}
                  </Typography>
                  <Chip 
                    label={`${effect.coins} coins`} 
                    size="small" 
                    color="primary" 
                    variant="outlined"
                    sx={{ mt: 0.5 }}
                  />
                </CardContent>
              </CardActionArea>
              
              {selectedEffect?.id === effect.id && (
                <Box sx={{ p: 1, pt: 0, display: 'flex', gap: 1 }}>
                  <Tooltip title="Preview Effect">
                    <IconButton
                      size="small"
                      color="primary"
                      onClick={handlePreviewEffect}
                    >
                      <PreviewIcon fontSize="small" />
                    </IconButton>
                  </Tooltip>
                  
                  {!effect.owned && (
                    <Button
                      variant="contained"
                      size="small"
                      fullWidth
                      startIcon={isPurchasing ? <CircularProgress size={16} /> : <CartIcon />}
                      onClick={handlePurchaseEffect}
                      disabled={isPurchasing}
                    >
                      Purchase
                    </Button>
                  )}
                  
                  {effect.owned && (
                    <Button
                      variant="contained"
                      size="small"
                      fullWidth
                      color="success"
                      startIcon={<CheckIcon />}
                      onClick={() => {
                        onSelectEffect(effect);
                        onClose();
                      }}
                    >
                      Select
                    </Button>
                  )}
                </Box>
              )}
            </Card>
          </Grid>
        ))}
      </Grid>
      
      {/* Preview Dialog */}
      <Dialog
        open={isPreviewDialogOpen}
        onClose={() => setIsPreviewDialogOpen(false)}
        maxWidth="sm"
      >
        <DialogTitle>Effect Preview: {selectedEffect?.name}</DialogTitle>
        
        <DialogContent dividers>
          <Box 
            sx={{ 
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              py: 2
            }}
          >
            {selectedEffect && getAnimationPreview(selectedEffect.animation)}
            
            <Typography variant="body2" color="text.secondary" sx={{ mt: 2 }}>
              {selectedEffect?.description}
            </Typography>
          </Box>
        </DialogContent>
        
        <DialogActions>
          <Button onClick={() => setIsPreviewDialogOpen(false)}>
            Close
          </Button>
          
          {selectedEffect && !selectedEffect.owned && (
            <Button 
              variant="contained" 
              color="primary" 
              onClick={handlePurchaseEffect}
              disabled={isPurchasing}
              startIcon={isPurchasing ? <CircularProgress size={20} /> : <CartIcon />}
            >
              Purchase for {selectedEffect?.coins} coins
            </Button>
          )}
          
          {selectedEffect && selectedEffect.owned && (
            <Button 
              variant="contained" 
              color="success" 
              onClick={() => {
                setIsPreviewDialogOpen(false);
                onSelectEffect(selectedEffect);
                onClose();
              }}
            >
              Select This Effect
            </Button>
          )}
        </DialogActions>
      </Dialog>
    </Box>
  );
};

export default GiftEffects;

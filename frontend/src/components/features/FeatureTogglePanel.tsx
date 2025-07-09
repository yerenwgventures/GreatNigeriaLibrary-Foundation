import React, { useState, useEffect } from 'react';
import {
  Box,
  Typography,
  Paper,
  Tabs,
  Tab,
  TextField,
  InputAdornment,
  Button,
  Divider,
  Alert,
  CircularProgress,
  useTheme
} from '@mui/material';
import {
  Search as SearchIcon,
  Refresh as RefreshIcon,
  Settings as SettingsIcon
} from '@mui/icons-material';
import FeatureToggle, { Feature } from './FeatureToggle';

interface FeatureTogglePanelProps {
  features: Feature[];
  onFeatureToggle: (featureId: string, enabled: boolean) => void;
  onRefresh?: () => void;
  onResetToDefault?: () => void;
  loading?: boolean;
  error?: string | null;
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
      id={`feature-tabpanel-${index}`}
      aria-labelledby={`feature-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 2 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

const FeatureTogglePanel: React.FC<FeatureTogglePanelProps> = ({
  features,
  onFeatureToggle,
  onRefresh,
  onResetToDefault,
  loading = false,
  error = null
}) => {
  const theme = useTheme();
  const [tabValue, setTabValue] = useState(0);
  const [searchQuery, setSearchQuery] = useState('');
  const [filteredFeatures, setFilteredFeatures] = useState<Feature[]>(features);
  const [categories, setCategories] = useState<string[]>([]);

  // Extract unique categories from features
  useEffect(() => {
    const uniqueCategories = Array.from(new Set(features.map(feature => feature.category)));
    setCategories(uniqueCategories);
  }, [features]);

  // Filter features based on search query and selected tab
  useEffect(() => {
    const query = searchQuery.toLowerCase();
    let filtered = features;
    
    // Apply search filter
    if (query) {
      filtered = filtered.filter(feature => 
        feature.name.toLowerCase().includes(query) || 
        feature.description.toLowerCase().includes(query)
      );
    }
    
    // Apply category filter (if not on "All" tab)
    if (tabValue > 0 && categories.length > 0) {
      const selectedCategory = categories[tabValue - 1];
      filtered = filtered.filter(feature => feature.category === selectedCategory);
    }
    
    setFilteredFeatures(filtered);
  }, [searchQuery, tabValue, features, categories]);

  // Handle tab change
  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };

  // Handle feature toggle
  const handleFeatureToggle = (featureId: string, enabled: boolean) => {
    onFeatureToggle(featureId, enabled);
  };

  return (
    <Paper elevation={3} sx={{ borderRadius: 2, overflow: 'hidden' }}>
      <Box
        sx={{
          p: 2,
          bgcolor: 'primary.main',
          color: 'primary.contrastText',
          display: 'flex',
          alignItems: 'center'
        }}
      >
        <SettingsIcon sx={{ mr: 1 }} />
        <Typography variant="h6" component="h2">
          Feature Settings
        </Typography>
      </Box>
      
      <Box p={2}>
        <TextField
          fullWidth
          placeholder="Search features..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          InputProps={{
            startAdornment: (
              <InputAdornment position="start">
                <SearchIcon />
              </InputAdornment>
            ),
          }}
          variant="outlined"
          size="small"
          sx={{ mb: 2 }}
        />
        
        {error && (
          <Alert severity="error" sx={{ mb: 2 }}>
            {error}
          </Alert>
        )}
        
        <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
          <Typography variant="body2" color="text.secondary">
            {filteredFeatures.length} feature{filteredFeatures.length !== 1 ? 's' : ''} found
          </Typography>
          
          <Box>
            <Button
              size="small"
              startIcon={<RefreshIcon />}
              onClick={onRefresh}
              disabled={loading}
              sx={{ mr: 1 }}
            >
              {loading ? <CircularProgress size={16} /> : 'Refresh'}
            </Button>
            
            <Button
              size="small"
              color="secondary"
              onClick={onResetToDefault}
              disabled={loading}
            >
              Reset to Default
            </Button>
          </Box>
        </Box>
        
        <Divider sx={{ mb: 2 }} />
        
        <Tabs
          value={tabValue}
          onChange={handleTabChange}
          variant="scrollable"
          scrollButtons="auto"
          sx={{ mb: 2, borderBottom: 1, borderColor: 'divider' }}
        >
          <Tab label="All Features" />
          {categories.map((category, index) => (
            <Tab key={index} label={category} />
          ))}
        </Tabs>
        
        <TabPanel value={tabValue} index={0}>
          {filteredFeatures.length > 0 ? (
            filteredFeatures.map(feature => (
              <FeatureToggle
                key={feature.id}
                feature={feature}
                onChange={handleFeatureToggle}
                allFeatures={features}
              />
            ))
          ) : (
            <Box textAlign="center" py={4}>
              <Typography variant="body1" color="text.secondary">
                No features found matching your search
              </Typography>
            </Box>
          )}
        </TabPanel>
        
        {categories.map((category, index) => (
          <TabPanel key={index} value={tabValue} index={index + 1}>
            {filteredFeatures.length > 0 ? (
              filteredFeatures.map(feature => (
                <FeatureToggle
                  key={feature.id}
                  feature={feature}
                  onChange={handleFeatureToggle}
                  allFeatures={features}
                />
              ))
            ) : (
              <Box textAlign="center" py={4}>
                <Typography variant="body1" color="text.secondary">
                  No features found in this category
                </Typography>
              </Box>
            )}
          </TabPanel>
        ))}
      </Box>
    </Paper>
  );
};

export default FeatureTogglePanel;

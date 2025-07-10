import React, { useState } from 'react';
import {
  Box,
  Typography,
  Switch,
  Paper,
  Grid,
  Tooltip,
  IconButton,
  Collapse,
  Divider,
  Chip,
  useTheme,
  alpha
} from '@mui/material';
import {
  Info as InfoIcon,
  ExpandMore as ExpandMoreIcon,
  ExpandLess as ExpandLessIcon,
  Lock as LockIcon
} from '@mui/icons-material';

export interface Feature {
  id: string;
  name: string;
  description: string;
  icon: React.ReactNode;
  enabled: boolean;
  category: string;
  premium?: boolean;
  beta?: boolean;
  dependencies?: string[];
  incompatibleWith?: string[];
  locked?: boolean;
  lockedReason?: string;
}

interface FeatureToggleProps {
  feature: Feature;
  onChange: (featureId: string, enabled: boolean) => void;
  allFeatures: Feature[];
}

const FeatureToggle: React.FC<FeatureToggleProps> = ({
  feature,
  onChange,
  allFeatures
}) => {
  const theme = useTheme();
  const [expanded, setExpanded] = useState(false);

  // Check if feature can be enabled based on dependencies and incompatibilities
  const canToggle = () => {
    // If feature is locked, it cannot be toggled
    if (feature.locked) {
      return false;
    }
    
    // If feature is being enabled, check dependencies
    if (!feature.enabled && feature.dependencies && feature.dependencies.length > 0) {
      // Check if all dependencies are enabled
      const missingDependencies = feature.dependencies.filter(depId => {
        const dependency = allFeatures.find(f => f.id === depId);
        return !dependency || !dependency.enabled;
      });
      
      if (missingDependencies.length > 0) {
        return false;
      }
    }
    
    // If feature is being enabled, check incompatibilities
    if (!feature.enabled && feature.incompatibleWith && feature.incompatibleWith.length > 0) {
      // Check if any incompatible features are enabled
      const activeIncompatibles = feature.incompatibleWith.filter(incId => {
        const incompatible = allFeatures.find(f => f.id === incId);
        return incompatible && incompatible.enabled;
      });
      
      if (activeIncompatibles.length > 0) {
        return false;
      }
    }
    
    return true;
  };

  // Get tooltip message based on feature status
  const getTooltipMessage = () => {
    if (feature.locked) {
      return feature.lockedReason || 'This feature is currently locked';
    }
    
    if (!feature.enabled && feature.dependencies && feature.dependencies.length > 0) {
      const missingDependencies = feature.dependencies
        .filter(depId => {
          const dependency = allFeatures.find(f => f.id === depId);
          return !dependency || !dependency.enabled;
        })
        .map(depId => {
          const dependency = allFeatures.find(f => f.id === depId);
          return dependency ? dependency.name : depId;
        });
      
      if (missingDependencies.length > 0) {
        return `Requires: ${missingDependencies.join(', ')}`;
      }
    }
    
    if (!feature.enabled && feature.incompatibleWith && feature.incompatibleWith.length > 0) {
      const activeIncompatibles = feature.incompatibleWith
        .filter(incId => {
          const incompatible = allFeatures.find(f => f.id === incId);
          return incompatible && incompatible.enabled;
        })
        .map(incId => {
          const incompatible = allFeatures.find(f => f.id === incId);
          return incompatible ? incompatible.name : incId;
        });
      
      if (activeIncompatibles.length > 0) {
        return `Incompatible with: ${activeIncompatibles.join(', ')}`;
      }
    }
    
    return feature.enabled ? 'Click to disable this feature' : 'Click to enable this feature';
  };

  // Handle toggle change
  const handleToggle = () => {
    if (canToggle()) {
      onChange(feature.id, !feature.enabled);
    }
  };

  return (
    <Paper
      elevation={2}
      sx={{
        mb: 2,
        borderRadius: 2,
        overflow: 'hidden',
        border: feature.enabled 
          ? `1px solid ${alpha(theme.palette.primary.main, 0.5)}` 
          : '1px solid transparent',
        transition: 'all 0.3s ease'
      }}
    >
      <Box
        sx={{
          p: 2,
          display: 'flex',
          alignItems: 'center',
          bgcolor: feature.enabled 
            ? alpha(theme.palette.primary.main, 0.05)
            : 'background.paper'
        }}
      >
        <Box
          sx={{
            mr: 2,
            color: feature.enabled 
              ? 'primary.main' 
              : feature.locked 
                ? 'text.disabled' 
                : 'text.secondary',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            width: 40,
            height: 40,
            borderRadius: '50%',
            bgcolor: feature.enabled 
              ? alpha(theme.palette.primary.main, 0.1)
              : alpha(theme.palette.grey[500], 0.1)
          }}
        >
          {feature.icon}
        </Box>
        
        <Box sx={{ flexGrow: 1 }}>
          <Box display="flex" alignItems="center">
            <Typography variant="subtitle1" fontWeight="medium">
              {feature.name}
            </Typography>
            
            <Box ml={1} display="flex" gap={0.5}>
              {feature.premium && (
                <Chip 
                  label="Premium" 
                  size="small" 
                  color="secondary" 
                  variant="outlined" 
                />
              )}
              
              {feature.beta && (
                <Chip 
                  label="Beta" 
                  size="small" 
                  color="warning" 
                  variant="outlined" 
                />
              )}
              
              {feature.locked && (
                <Chip 
                  label="Locked" 
                  size="small" 
                  color="default" 
                  variant="outlined"
                  icon={<LockIcon fontSize="small" />}
                />
              )}
            </Box>
          </Box>
          
          <Typography variant="body2" color="text.secondary" noWrap>
            {feature.description}
          </Typography>
        </Box>
        
        <Box display="flex" alignItems="center">
          <Tooltip title={getTooltipMessage()} arrow>
            <Box>
              <Switch
                checked={feature.enabled}
                onChange={handleToggle}
                disabled={!canToggle()}
                color="primary"
              />
            </Box>
          </Tooltip>
          
          <Tooltip title="More information" arrow>
            <IconButton 
              size="small" 
              onClick={() => setExpanded(!expanded)}
              sx={{ ml: 1 }}
            >
              {expanded ? <ExpandLessIcon /> : <ExpandMoreIcon />}
            </IconButton>
          </Tooltip>
        </Box>
      </Box>
      
      <Collapse in={expanded}>
        <Divider />
        <Box p={2} bgcolor={alpha(theme.palette.background.default, 0.5)}>
          <Typography variant="body2" paragraph>
            {feature.description}
          </Typography>
          
          {feature.dependencies && feature.dependencies.length > 0 && (
            <Box mb={1}>
              <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
                Required Features:
              </Typography>
              <Box display="flex" flexWrap="wrap" gap={0.5}>
                {feature.dependencies.map(depId => {
                  const dependency = allFeatures.find(f => f.id === depId);
                  return (
                    <Chip
                      key={depId}
                      label={dependency ? dependency.name : depId}
                      size="small"
                      color={dependency?.enabled ? 'success' : 'default'}
                      variant="outlined"
                    />
                  );
                })}
              </Box>
            </Box>
          )}
          
          {feature.incompatibleWith && feature.incompatibleWith.length > 0 && (
            <Box>
              <Typography variant="caption" color="text.secondary" display="block" gutterBottom>
                Incompatible With:
              </Typography>
              <Box display="flex" flexWrap="wrap" gap={0.5}>
                {feature.incompatibleWith.map(incId => {
                  const incompatible = allFeatures.find(f => f.id === incId);
                  return (
                    <Chip
                      key={incId}
                      label={incompatible ? incompatible.name : incId}
                      size="small"
                      color={incompatible?.enabled ? 'error' : 'default'}
                      variant="outlined"
                    />
                  );
                })}
              </Box>
            </Box>
          )}
          
          {feature.locked && feature.lockedReason && (
            <Box mt={1} display="flex" alignItems="center" gap={1}>
              <InfoIcon color="info" fontSize="small" />
              <Typography variant="body2" color="info.main">
                {feature.lockedReason}
              </Typography>
            </Box>
          )}
        </Box>
      </Collapse>
    </Paper>
  );
};

export default FeatureToggle;

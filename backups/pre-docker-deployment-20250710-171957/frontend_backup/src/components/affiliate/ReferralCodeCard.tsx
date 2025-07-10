import React, { useState } from 'react';
import {
  Card,
  CardContent,
  Typography,
  Box,
  Chip,
  IconButton,
  Button,
  TextField,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Tooltip,
  useTheme,
  alpha
} from '@mui/material';
import {
  ContentCopy as CopyIcon,
  Edit as EditIcon,
  Delete as DeleteIcon,
  Share as ShareIcon,
  Check as CheckIcon,
  Link as LinkIcon
} from '@mui/icons-material';
import { ReferralCode } from '../../api/affiliateService';

interface ReferralCodeCardProps {
  referralCode: ReferralCode;
  onEdit?: (referralCode: ReferralCode, description: string) => void;
  onDelete?: (referralCode: ReferralCode) => void;
  onCopy?: (referralCode: ReferralCode) => void;
  onShare?: (referralCode: ReferralCode) => void;
  baseUrl?: string;
}

const ReferralCodeCard: React.FC<ReferralCodeCardProps> = ({
  referralCode,
  onEdit,
  onDelete,
  onCopy,
  onShare,
  baseUrl = 'https://greatnigeria.com/ref/'
}) => {
  const theme = useTheme();
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [description, setDescription] = useState(referralCode.description || '');
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    const fullUrl = `${baseUrl}${referralCode.code}`;
    navigator.clipboard.writeText(fullUrl);
    setCopied(true);
    
    if (onCopy) {
      onCopy(referralCode);
    }
    
    setTimeout(() => {
      setCopied(false);
    }, 2000);
  };

  const handleShare = () => {
    if (onShare) {
      onShare(referralCode);
    } else {
      const fullUrl = `${baseUrl}${referralCode.code}`;
      if (navigator.share) {
        navigator.share({
          title: 'Join Great Nigeria with my referral code',
          text: 'Use my referral code to join Great Nigeria and get special benefits!',
          url: fullUrl
        });
      } else {
        handleCopy();
      }
    }
  };

  const handleEditSubmit = () => {
    if (onEdit) {
      onEdit(referralCode, description);
    }
    setEditDialogOpen(false);
  };

  const handleDeleteSubmit = () => {
    if (onDelete) {
      onDelete(referralCode);
    }
    setDeleteDialogOpen(false);
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-NG', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  return (
    <Card 
      elevation={2} 
      sx={{ 
        height: '100%', 
        display: 'flex', 
        flexDirection: 'column',
        transition: 'transform 0.3s, box-shadow 0.3s',
        '&:hover': {
          transform: 'translateY(-5px)',
          boxShadow: theme.shadows[8]
        },
        borderRadius: 2,
        overflow: 'hidden',
        position: 'relative'
      }}
    >
      <Box 
        sx={{ 
          position: 'absolute', 
          top: 0, 
          right: 0, 
          width: 8, 
          height: '100%', 
          bgcolor: referralCode.isActive ? theme.palette.success.main : theme.palette.grey[500]
        }} 
      />
      
      <CardContent sx={{ flexGrow: 1, p: 3 }}>
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 2 }}>
          <Typography variant="h6" fontWeight="bold">
            {referralCode.code}
          </Typography>
          
          <Chip 
            label={referralCode.isActive ? 'Active' : 'Inactive'} 
            color={referralCode.isActive ? 'success' : 'default'} 
            size="small" 
          />
        </Box>
        
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2, minHeight: 40 }}>
          {referralCode.description || 'No description provided'}
        </Typography>
        
        <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
          <Box>
            <Typography variant="caption" color="text.secondary">
              Created
            </Typography>
            <Typography variant="body2">
              {formatDate(referralCode.createdAt)}
            </Typography>
          </Box>
          
          <Box sx={{ textAlign: 'right' }}>
            <Typography variant="caption" color="text.secondary">
              Usage
            </Typography>
            <Typography variant="body2">
              {referralCode.usageCount} uses
            </Typography>
          </Box>
        </Box>
        
        <Box sx={{ 
          p: 2, 
          bgcolor: alpha(theme.palette.primary.main, 0.1), 
          borderRadius: 1,
          display: 'flex',
          alignItems: 'center',
          mb: 2
        }}>
          <LinkIcon color="primary" fontSize="small" sx={{ mr: 1 }} />
          <Typography 
            variant="body2" 
            sx={{ 
              flexGrow: 1,
              overflow: 'hidden',
              textOverflow: 'ellipsis',
              whiteSpace: 'nowrap'
            }}
          >
            {`${baseUrl}${referralCode.code}`}
          </Typography>
          
          <Tooltip title={copied ? 'Copied!' : 'Copy link'}>
            <IconButton size="small" onClick={handleCopy} color={copied ? 'success' : 'primary'}>
              {copied ? <CheckIcon fontSize="small" /> : <CopyIcon fontSize="small" />}
            </IconButton>
          </Tooltip>
        </Box>
        
        <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <Box>
            <Typography variant="caption" color="text.secondary">
              Conversions
            </Typography>
            <Typography variant="body2" fontWeight="bold" color="primary">
              {referralCode.conversionCount}
            </Typography>
          </Box>
          
          <Box>
            <Typography variant="caption" color="text.secondary">
              Conversion Rate
            </Typography>
            <Typography variant="body2" fontWeight="bold" color="primary">
              {referralCode.usageCount > 0 
                ? `${Math.round((referralCode.conversionCount / referralCode.usageCount) * 100)}%` 
                : '0%'}
            </Typography>
          </Box>
        </Box>
      </CardContent>
      
      <Box sx={{ display: 'flex', justifyContent: 'space-between', p: 2, borderTop: `1px solid ${theme.palette.divider}` }}>
        {onEdit && (
          <Tooltip title="Edit">
            <IconButton size="small" onClick={() => setEditDialogOpen(true)}>
              <EditIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        )}
        
        {onDelete && (
          <Tooltip title="Delete">
            <IconButton size="small" color="error" onClick={() => setDeleteDialogOpen(true)}>
              <DeleteIcon fontSize="small" />
            </IconButton>
          </Tooltip>
        )}
        
        <Tooltip title="Share">
          <IconButton size="small" color="primary" onClick={handleShare}>
            <ShareIcon fontSize="small" />
          </IconButton>
        </Tooltip>
      </Box>
      
      {/* Edit Dialog */}
      <Dialog open={editDialogOpen} onClose={() => setEditDialogOpen(false)}>
        <DialogTitle>Edit Referral Code</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Description"
            fullWidth
            multiline
            rows={3}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleEditSubmit} color="primary">Save</Button>
        </DialogActions>
      </Dialog>
      
      {/* Delete Dialog */}
      <Dialog open={deleteDialogOpen} onClose={() => setDeleteDialogOpen(false)}>
        <DialogTitle>Delete Referral Code</DialogTitle>
        <DialogContent>
          <Typography>
            Are you sure you want to delete the referral code "{referralCode.code}"? This action cannot be undone.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDeleteDialogOpen(false)}>Cancel</Button>
          <Button onClick={handleDeleteSubmit} color="error">Delete</Button>
        </DialogActions>
      </Dialog>
    </Card>
  );
};

export default ReferralCodeCard;

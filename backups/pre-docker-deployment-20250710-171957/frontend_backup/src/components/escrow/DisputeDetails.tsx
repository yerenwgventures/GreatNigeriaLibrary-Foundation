import React, { useState, useRef } from 'react';
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
  TextField,
  List,
  ListItem,
  ListItemText,
  ListItemAvatar,
  Avatar,
  Card,
  CardContent,
  CardMedia,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  useTheme,
  alpha
} from '@mui/material';
import {
  Gavel as DisputeIcon,
  AccessTime as TimeIcon,
  Send as SendIcon,
  AttachFile as AttachIcon,
  Image as ImageIcon,
  Description as DocumentIcon,
  Videocam as VideoIcon,
  AudioFile as AudioIcon,
  TextFields as TextIcon,
  Close as CloseIcon,
  CheckCircle as ResolvedIcon,
  Cancel as CancelledIcon,
  Search as SearchIcon,
  ArrowForward as ReleaseIcon,
  ArrowBack as RefundIcon,
  Balance as SplitIcon,
  Person as PersonIcon,
  AdminPanelSettings as AdminIcon
} from '@mui/icons-material';
import { Dispute, DisputeEvidence, DisputeMessage, DisputeResolutionRequest } from '../../api/escrowService';

interface DisputeDetailsProps {
  dispute: Dispute;
  evidence: DisputeEvidence[];
  messages: DisputeMessage[];
  currentUserId: string;
  isAdmin?: boolean;
  onAddEvidence?: (evidence: { type: string; content: string; fileUrl?: string }) => void;
  onAddMessage?: (message: string) => void;
  onResolveDispute?: (resolution: DisputeResolutionRequest) => void;
  onCloseDispute?: () => void;
  loading?: boolean;
}

const DisputeDetails: React.FC<DisputeDetailsProps> = ({
  dispute,
  evidence,
  messages,
  currentUserId,
  isAdmin = false,
  onAddEvidence,
  onAddMessage,
  onResolveDispute,
  onCloseDispute,
  loading = false
}) => {
  const theme = useTheme();
  const fileInputRef = useRef<HTMLInputElement>(null);
  
  const [newMessage, setNewMessage] = useState('');
  const [newEvidence, setNewEvidence] = useState({
    type: 'text',
    content: '',
    fileUrl: ''
  });
  const [showEvidenceForm, setShowEvidenceForm] = useState(false);
  const [showResolutionForm, setShowResolutionForm] = useState(false);
  const [resolution, setResolution] = useState<DisputeResolutionRequest>({
    resolution: 'release',
    resolutionDetails: '',
    resolutionAmount: undefined
  });
  const [selectedEvidence, setSelectedEvidence] = useState<DisputeEvidence | null>(null);
  
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
  
  // Format currency
  const formatCurrency = (amount: number | undefined, currency: string = 'NGN') => {
    if (amount === undefined) return 'N/A';
    
    return new Intl.NumberFormat('en-NG', {
      style: 'currency',
      currency,
      minimumFractionDigits: 2
    }).format(amount);
  };
  
  // Get status chip
  const getStatusChip = (status: string) => {
    switch (status) {
      case 'open':
        return <Chip size="small" label="Open" color="warning" />;
      case 'under_review':
        return <Chip size="small" label="Under Review" color="info" />;
      case 'resolved':
        return <Chip size="small" label="Resolved" color="success" />;
      case 'closed':
        return <Chip size="small" label="Closed" color="default" />;
      default:
        return <Chip size="small" label={status} />;
    }
  };
  
  // Get resolution chip
  const getResolutionChip = (resolution?: string) => {
    if (!resolution) return null;
    
    switch (resolution) {
      case 'release':
        return <Chip size="small" label="Released to Seller" color="success" icon={<ReleaseIcon />} />;
      case 'refund':
        return <Chip size="small" label="Refunded to Buyer" color="secondary" icon={<RefundIcon />} />;
      case 'partial_release':
        return <Chip size="small" label="Partial Release" color="info" icon={<ReleaseIcon />} />;
      case 'split':
        return <Chip size="small" label="Split Between Parties" color="primary" icon={<SplitIcon />} />;
      default:
        return <Chip size="small" label={resolution} />;
    }
  };
  
  // Handle file upload
  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;
    
    // In a real app, you would upload the file to a server and get a URL back
    // For now, we'll just create a fake URL
    const fakeUrl = URL.createObjectURL(file);
    
    // Determine file type
    let type: 'image' | 'document' | 'video' | 'audio' = 'document';
    if (file.type.startsWith('image/')) {
      type = 'image';
    } else if (file.type.startsWith('video/')) {
      type = 'video';
    } else if (file.type.startsWith('audio/')) {
      type = 'audio';
    }
    
    setNewEvidence({
      ...newEvidence,
      type,
      fileUrl: fakeUrl,
      content: file.name
    });
  };
  
  // Handle submit evidence
  const handleSubmitEvidence = () => {
    if (!onAddEvidence) return;
    
    if (newEvidence.type === 'text' && !newEvidence.content) {
      return;
    }
    
    if (newEvidence.type !== 'text' && !newEvidence.fileUrl) {
      return;
    }
    
    onAddEvidence(newEvidence);
    
    setNewEvidence({
      type: 'text',
      content: '',
      fileUrl: ''
    });
    
    setShowEvidenceForm(false);
  };
  
  // Handle submit message
  const handleSubmitMessage = () => {
    if (!onAddMessage || !newMessage.trim()) return;
    
    onAddMessage(newMessage);
    setNewMessage('');
  };
  
  // Handle submit resolution
  const handleSubmitResolution = () => {
    if (!onResolveDispute) return;
    
    onResolveDispute(resolution);
    setShowResolutionForm(false);
  };
  
  // Get evidence icon
  const getEvidenceIcon = (type: string) => {
    switch (type) {
      case 'image':
        return <ImageIcon />;
      case 'document':
        return <DocumentIcon />;
      case 'video':
        return <VideoIcon />;
      case 'audio':
        return <AudioIcon />;
      case 'text':
        return <TextIcon />;
      default:
        return <AttachIcon />;
    }
  };
  
  // Check if user can add evidence or messages
  const canAddEvidence = dispute.status === 'open' || dispute.status === 'under_review';
  const canAddMessage = dispute.status === 'open' || dispute.status === 'under_review';
  const canResolve = isAdmin && (dispute.status === 'open' || dispute.status === 'under_review');
  const canClose = isAdmin && dispute.status === 'resolved';
  
  return (
    <Box>
      <Paper elevation={2} sx={{ p: 3, borderRadius: 2, mb: 3 }}>
        <Grid container spacing={3}>
          <Grid item xs={12} md={8}>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              <DisputeIcon color="error" sx={{ mr: 1 }} />
              <Typography variant="h5" component="h1">
                Dispute #{dispute.id}
              </Typography>
              <Box sx={{ ml: 2 }}>
                {getStatusChip(dispute.status)}
              </Box>
            </Box>
            
            <Typography variant="body1" paragraph>
              <strong>Reason:</strong> {dispute.reason}
            </Typography>
            
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
              <TimeIcon fontSize="small" color="action" sx={{ mr: 1 }} />
              <Typography variant="body2" color="text.secondary">
                Opened: {formatDate(dispute.createdAt)}
              </Typography>
            </Box>
            
            {dispute.resolvedAt && (
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                <ResolvedIcon fontSize="small" color="success" sx={{ mr: 1 }} />
                <Typography variant="body2" color="success.main">
                  Resolved: {formatDate(dispute.resolvedAt)}
                </Typography>
              </Box>
            )}
            
            {dispute.closedAt && (
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                <CancelledIcon fontSize="small" color="action" sx={{ mr: 1 }} />
                <Typography variant="body2" color="text.secondary">
                  Closed: {formatDate(dispute.closedAt)}
                </Typography>
              </Box>
            )}
            
            <Box sx={{ display: 'flex', alignItems: 'center', mt: 2 }}>
              <Typography variant="body2">
                <strong>Order ID:</strong> {dispute.orderId}
              </Typography>
            </Box>
            
            <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
              <Typography variant="body2">
                <strong>Escrow Transaction:</strong> {dispute.escrowTransactionId}
              </Typography>
            </Box>
          </Grid>
          
          <Grid item xs={12} md={4}>
            <Box sx={{ 
              p: 2, 
              bgcolor: alpha(theme.palette.primary.main, 0.1), 
              borderRadius: 2,
              border: `1px solid ${alpha(theme.palette.primary.main, 0.2)}`
            }}>
              <Typography variant="h6" gutterBottom>
                Resolution
              </Typography>
              
              {dispute.resolution ? (
                <>
                  <Box sx={{ mb: 1 }}>
                    {getResolutionChip(dispute.resolution)}
                  </Box>
                  
                  {dispute.resolutionDetails && (
                    <Typography variant="body2" paragraph>
                      <strong>Details:</strong> {dispute.resolutionDetails}
                    </Typography>
                  )}
                  
                  {dispute.resolutionAmount !== undefined && (
                    <Typography variant="body2">
                      <strong>Amount:</strong> {formatCurrency(dispute.resolutionAmount)}
                    </Typography>
                  )}
                </>
              ) : (
                <Typography variant="body2" color="text.secondary">
                  This dispute has not been resolved yet.
                </Typography>
              )}
              
              {canResolve && (
                <Button
                  variant="contained"
                  color="primary"
                  fullWidth
                  sx={{ mt: 2 }}
                  onClick={() => setShowResolutionForm(true)}
                  disabled={loading}
                >
                  Resolve Dispute
                </Button>
              )}
              
              {canClose && (
                <Button
                  variant="outlined"
                  color="primary"
                  fullWidth
                  sx={{ mt: 2 }}
                  onClick={onCloseDispute}
                  disabled={loading}
                >
                  Close Dispute
                </Button>
              )}
            </Box>
            
            <Box sx={{ 
              p: 2, 
              bgcolor: alpha(theme.palette.info.main, 0.1), 
              borderRadius: 2,
              border: `1px solid ${alpha(theme.palette.info.main, 0.2)}`,
              mt: 2
            }}>
              <Typography variant="h6" gutterBottom>
                Parties
              </Typography>
              
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                <PersonIcon fontSize="small" color="primary" sx={{ mr: 1 }} />
                <Typography variant="body2">
                  <strong>Initiator:</strong> {dispute.initiatorId}
                </Typography>
              </Box>
              
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <PersonIcon fontSize="small" color="secondary" sx={{ mr: 1 }} />
                <Typography variant="body2">
                  <strong>Respondent:</strong> {dispute.respondentId}
                </Typography>
              </Box>
            </Box>
          </Grid>
        </Grid>
      </Paper>
      
      <Grid container spacing={3}>
        <Grid item xs={12} md={6}>
          <Paper elevation={2} sx={{ p: 3, borderRadius: 2, height: '100%' }}>
            <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
              <Typography variant="h6">
                Evidence
              </Typography>
              
              {canAddEvidence && onAddEvidence && (
                <Button
                  variant="outlined"
                  color="primary"
                  size="small"
                  startIcon={<AttachIcon />}
                  onClick={() => setShowEvidenceForm(true)}
                  disabled={loading}
                >
                  Add Evidence
                </Button>
              )}
            </Box>
            
            {evidence.length === 0 ? (
              <Typography variant="body2" color="text.secondary" sx={{ py: 2 }}>
                No evidence has been submitted yet.
              </Typography>
            ) : (
              <List>
                {evidence.map((item) => (
                  <Card 
                    key={item.id} 
                    sx={{ 
                      mb: 2, 
                      cursor: item.type !== 'text' ? 'pointer' : 'default',
                      '&:hover': {
                        boxShadow: item.type !== 'text' ? theme.shadows[4] : undefined
                      }
                    }}
                    onClick={() => item.type !== 'text' && setSelectedEvidence(item)}
                  >
                    <CardContent sx={{ p: 2 }}>
                      <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                        <Avatar sx={{ bgcolor: theme.palette.primary.main, mr: 1 }}>
                          {getEvidenceIcon(item.type)}
                        </Avatar>
                        <Box>
                          <Typography variant="subtitle2">
                            {item.userId === currentUserId ? 'You' : 
                             item.userId === dispute.initiatorId ? 'Initiator' : 
                             item.userId === dispute.respondentId ? 'Respondent' : 
                             'Admin'}
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            {formatDate(item.createdAt)}
                          </Typography>
                        </Box>
                      </Box>
                      
                      {item.type === 'text' ? (
                        <Typography variant="body2" sx={{ mt: 1, whiteSpace: 'pre-wrap' }}>
                          {item.content}
                        </Typography>
                      ) : item.type === 'image' && item.fileUrl ? (
                        <Box sx={{ mt: 1, position: 'relative' }}>
                          <CardMedia
                            component="img"
                            height="140"
                            image={item.fileUrl}
                            alt={item.content}
                            sx={{ borderRadius: 1 }}
                          />
                          <Box 
                            sx={{ 
                              position: 'absolute', 
                              bottom: 8, 
                              left: 8, 
                              bgcolor: alpha(theme.palette.background.paper, 0.7),
                              px: 1,
                              py: 0.5,
                              borderRadius: 1
                            }}
                          >
                            <Typography variant="caption">
                              {item.content}
                            </Typography>
                          </Box>
                        </Box>
                      ) : (
                        <Box sx={{ 
                          display: 'flex', 
                          alignItems: 'center', 
                          mt: 1,
                          p: 1,
                          bgcolor: alpha(theme.palette.primary.main, 0.1),
                          borderRadius: 1
                        }}>
                          {getEvidenceIcon(item.type)}
                          <Typography variant="body2" sx={{ ml: 1 }}>
                            {item.content}
                          </Typography>
                        </Box>
                      )}
                    </CardContent>
                  </Card>
                ))}
              </List>
            )}
            
            {showEvidenceForm && (
              <Box sx={{ mt: 2, p: 2, bgcolor: alpha(theme.palette.background.default, 0.5), borderRadius: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                  <Typography variant="subtitle1" fontWeight="bold">
                    Add Evidence
                  </Typography>
                  <IconButton size="small" onClick={() => setShowEvidenceForm(false)}>
                    <CloseIcon />
                  </IconButton>
                </Box>
                
                <FormControl fullWidth size="small" sx={{ mb: 2 }}>
                  <InputLabel id="evidence-type-label">Type</InputLabel>
                  <Select
                    labelId="evidence-type-label"
                    value={newEvidence.type}
                    label="Type"
                    onChange={(e) => setNewEvidence({ ...newEvidence, type: e.target.value as any })}
                  >
                    <MenuItem value="text">Text</MenuItem>
                    <MenuItem value="image">Image</MenuItem>
                    <MenuItem value="document">Document</MenuItem>
                    <MenuItem value="video">Video</MenuItem>
                    <MenuItem value="audio">Audio</MenuItem>
                  </Select>
                </FormControl>
                
                {newEvidence.type === 'text' ? (
                  <TextField
                    fullWidth
                    multiline
                    rows={4}
                    label="Evidence Text"
                    value={newEvidence.content}
                    onChange={(e) => setNewEvidence({ ...newEvidence, content: e.target.value })}
                    sx={{ mb: 2 }}
                  />
                ) : (
                  <>
                    <input
                      type="file"
                      ref={fileInputRef}
                      style={{ display: 'none' }}
                      onChange={handleFileUpload}
                      accept={
                        newEvidence.type === 'image' ? 'image/*' :
                        newEvidence.type === 'document' ? '.pdf,.doc,.docx,.txt' :
                        newEvidence.type === 'video' ? 'video/*' :
                        newEvidence.type === 'audio' ? 'audio/*' :
                        undefined
                      }
                    />
                    
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                      <Button
                        variant="outlined"
                        onClick={() => fileInputRef.current?.click()}
                        startIcon={getEvidenceIcon(newEvidence.type)}
                      >
                        Select File
                      </Button>
                      
                      {newEvidence.fileUrl && (
                        <Typography variant="body2" sx={{ ml: 2 }}>
                          {newEvidence.content}
                        </Typography>
                      )}
                    </Box>
                  </>
                )}
                
                <Box sx={{ display: 'flex', justifyContent: 'flex-end' }}>
                  <Button
                    variant="contained"
                    color="primary"
                    onClick={handleSubmitEvidence}
                    disabled={
                      loading || 
                      (newEvidence.type === 'text' && !newEvidence.content) ||
                      (newEvidence.type !== 'text' && !newEvidence.fileUrl)
                    }
                  >
                    Submit Evidence
                  </Button>
                </Box>
              </Box>
            )}
          </Paper>
        </Grid>
        
        <Grid item xs={12} md={6}>
          <Paper elevation={2} sx={{ p: 3, borderRadius: 2, height: '100%', display: 'flex', flexDirection: 'column' }}>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Messages
            </Typography>
            
            <Box sx={{ flexGrow: 1, overflow: 'auto', maxHeight: 400 }}>
              {messages.length === 0 ? (
                <Typography variant="body2" color="text.secondary" sx={{ py: 2 }}>
                  No messages yet.
                </Typography>
              ) : (
                <List>
                  {messages.map((message) => (
                    <ListItem
                      key={message.id}
                      alignItems="flex-start"
                      sx={{
                        bgcolor: message.isAdminMessage ? alpha(theme.palette.warning.main, 0.1) : 'transparent',
                        borderRadius: 2,
                        mb: 1
                      }}
                    >
                      <ListItemAvatar>
                        <Avatar sx={{ 
                          bgcolor: message.isAdminMessage ? theme.palette.warning.main : 
                                  message.userId === currentUserId ? theme.palette.primary.main : 
                                  theme.palette.secondary.main 
                        }}>
                          {message.isAdminMessage ? <AdminIcon /> : 
                           message.userId === currentUserId ? 'You' : 
                           message.userId === dispute.initiatorId ? 'I' : 'R'}
                        </Avatar>
                      </ListItemAvatar>
                      <ListItemText
                        primary={
                          <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                            <Typography variant="subtitle2">
                              {message.isAdminMessage ? 'Admin' : 
                               message.userId === currentUserId ? 'You' : 
                               message.userId === dispute.initiatorId ? 'Initiator' : 
                               'Respondent'}
                            </Typography>
                            <Typography variant="caption" color="text.secondary">
                              {formatDate(message.createdAt)}
                            </Typography>
                          </Box>
                        }
                        secondary={
                          <Typography
                            variant="body2"
                            color="text.primary"
                            sx={{ mt: 1, whiteSpace: 'pre-wrap' }}
                          >
                            {message.message}
                          </Typography>
                        }
                      />
                    </ListItem>
                  ))}
                </List>
              )}
            </Box>
            
            {canAddMessage && onAddMessage && (
              <Box sx={{ mt: 2, display: 'flex', alignItems: 'flex-start' }}>
                <TextField
                  fullWidth
                  multiline
                  rows={3}
                  placeholder="Type your message..."
                  value={newMessage}
                  onChange={(e) => setNewMessage(e.target.value)}
                  disabled={loading}
                  sx={{ mr: 1 }}
                />
                <Button
                  variant="contained"
                  color="primary"
                  endIcon={<SendIcon />}
                  onClick={handleSubmitMessage}
                  disabled={loading || !newMessage.trim()}
                  sx={{ mt: 1, minWidth: 100 }}
                >
                  Send
                </Button>
              </Box>
            )}
          </Paper>
        </Grid>
      </Grid>
      
      {/* Resolution Form Dialog */}
      <Dialog open={showResolutionForm} onClose={() => setShowResolutionForm(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Resolve Dispute</DialogTitle>
        <DialogContent>
          <Typography variant="body2" paragraph>
            Please select a resolution for this dispute. This action cannot be undone.
          </Typography>
          
          <FormControl fullWidth margin="normal">
            <InputLabel id="resolution-type-label">Resolution Type</InputLabel>
            <Select
              labelId="resolution-type-label"
              value={resolution.resolution}
              label="Resolution Type"
              onChange={(e) => setResolution({ ...resolution, resolution: e.target.value as any })}
            >
              <MenuItem value="release">Release Funds to Seller</MenuItem>
              <MenuItem value="refund">Refund Funds to Buyer</MenuItem>
              <MenuItem value="partial_release">Partial Release to Seller</MenuItem>
              <MenuItem value="split">Split Between Parties</MenuItem>
            </Select>
          </FormControl>
          
          {(resolution.resolution === 'partial_release' || resolution.resolution === 'split') && (
            <TextField
              fullWidth
              type="number"
              label="Amount"
              margin="normal"
              value={resolution.resolutionAmount || ''}
              onChange={(e) => setResolution({ 
                ...resolution, 
                resolutionAmount: e.target.value ? Number(e.target.value) : undefined 
              })}
              helperText={
                resolution.resolution === 'partial_release' 
                  ? 'Amount to release to the seller' 
                  : 'Amount to give to the seller (remaining goes to buyer)'
              }
            />
          )}
          
          <TextField
            fullWidth
            multiline
            rows={4}
            label="Resolution Details"
            margin="normal"
            value={resolution.resolutionDetails || ''}
            onChange={(e) => setResolution({ ...resolution, resolutionDetails: e.target.value })}
            placeholder="Explain your decision..."
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setShowResolutionForm(false)}>Cancel</Button>
          <Button 
            variant="contained" 
            color="primary" 
            onClick={handleSubmitResolution}
            disabled={
              loading || 
              ((resolution.resolution === 'partial_release' || resolution.resolution === 'split') && 
               resolution.resolutionAmount === undefined)
            }
          >
            Resolve Dispute
          </Button>
        </DialogActions>
      </Dialog>
      
      {/* Evidence Viewer Dialog */}
      <Dialog 
        open={!!selectedEvidence} 
        onClose={() => setSelectedEvidence(null)} 
        maxWidth="md" 
        fullWidth
      >
        {selectedEvidence && (
          <>
            <DialogTitle>
              <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Typography variant="h6">
                  {selectedEvidence.content}
                </Typography>
                <IconButton onClick={() => setSelectedEvidence(null)}>
                  <CloseIcon />
                </IconButton>
              </Box>
            </DialogTitle>
            <DialogContent>
              {selectedEvidence.type === 'image' && selectedEvidence.fileUrl ? (
                <Box sx={{ display: 'flex', justifyContent: 'center' }}>
                  <img 
                    src={selectedEvidence.fileUrl} 
                    alt={selectedEvidence.content} 
                    style={{ maxWidth: '100%', maxHeight: '70vh' }} 
                  />
                </Box>
              ) : selectedEvidence.type === 'video' && selectedEvidence.fileUrl ? (
                <Box sx={{ display: 'flex', justifyContent: 'center' }}>
                  <video 
                    src={selectedEvidence.fileUrl} 
                    controls 
                    style={{ maxWidth: '100%', maxHeight: '70vh' }} 
                  />
                </Box>
              ) : selectedEvidence.type === 'audio' && selectedEvidence.fileUrl ? (
                <Box sx={{ display: 'flex', justifyContent: 'center', p: 3 }}>
                  <audio 
                    src={selectedEvidence.fileUrl} 
                    controls 
                    style={{ width: '100%' }} 
                  />
                </Box>
              ) : (
                <Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', flexDirection: 'column' }}>
                  {getEvidenceIcon(selectedEvidence.type)}
                  <Typography variant="body1" sx={{ mt: 2 }}>
                    {selectedEvidence.content}
                  </Typography>
                  {selectedEvidence.fileUrl && (
                    <Button
                      variant="contained"
                      color="primary"
                      href={selectedEvidence.fileUrl}
                      target="_blank"
                      rel="noopener noreferrer"
                      sx={{ mt: 2 }}
                    >
                      Download File
                    </Button>
                  )}
                </Box>
              )}
            </DialogContent>
          </>
        )}
      </Dialog>
    </Box>
  );
};

export default DisputeDetails;

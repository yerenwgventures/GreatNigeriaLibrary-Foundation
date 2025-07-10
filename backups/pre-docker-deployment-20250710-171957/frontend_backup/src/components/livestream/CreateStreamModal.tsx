import React, { useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { 
  Dialog, 
  DialogTitle, 
  DialogContent, 
  DialogActions, 
  Button, 
  TextField, 
  Grid, 
  FormControlLabel, 
  Switch, 
  Typography, 
  Box, 
  CircularProgress,
  Chip,
  InputAdornment,
  IconButton
} from '@mui/material';
import { 
  CalendarMonth as CalendarIcon,
  AddPhotoAlternate as PhotoIcon
} from '@mui/icons-material';
import { DateTimePicker } from '@mui/x-date-pickers/DateTimePicker';
import { LocalizationProvider } from '@mui/x-date-pickers/LocalizationProvider';
import { AdapterDateFns } from '@mui/x-date-pickers/AdapterDateFns';
import { createStream } from '../../features/livestream/livestreamSlice';
import { RootState } from '../../store';
import { useNavigate } from 'react-router-dom';

interface CreateStreamModalProps {
  open: boolean;
  onClose: () => void;
}

const CreateStreamModal: React.FC<CreateStreamModalProps> = ({ open, onClose }) => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { loading, error } = useSelector((state: RootState) => state.livestream.currentStream);
  
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [thumbnailUrl, setThumbnailUrl] = useState('');
  const [scheduledStart, setScheduledStart] = useState<Date>(new Date(Date.now() + 30 * 60000)); // 30 minutes from now
  const [isPrivate, setIsPrivate] = useState(false);
  const [categories, setCategories] = useState('');
  const [tags, setTags] = useState<string[]>([]);
  const [currentTag, setCurrentTag] = useState('');
  
  const handleAddTag = () => {
    if (currentTag.trim() && !tags.includes(currentTag.trim())) {
      setTags([...tags, currentTag.trim()]);
      setCurrentTag('');
    }
  };
  
  const handleDeleteTag = (tagToDelete: string) => {
    setTags(tags.filter((tag) => tag !== tagToDelete));
  };
  
  const handleSubmit = async () => {
    if (!title || !scheduledStart) return;
    
    const streamData = {
      title,
      description,
      thumbnailUrl,
      scheduledStart: scheduledStart.toISOString(),
      isPrivate,
      categories,
      tags: tags.join(',')
    };
    
    const resultAction = await dispatch(createStream(streamData) as any);
    
    if (createStream.fulfilled.match(resultAction)) {
      onClose();
      navigate(`/livestream/${resultAction.payload.id}`);
    }
  };
  
  return (
    <Dialog open={open} onClose={onClose} maxWidth="md" fullWidth>
      <DialogTitle>Create New Stream</DialogTitle>
      
      <DialogContent dividers>
        {error && (
          <Typography color="error" sx={{ mb: 2 }}>
            {error}
          </Typography>
        )}
        
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <TextField
              label="Stream Title"
              fullWidth
              required
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              variant="outlined"
              placeholder="Enter an engaging title for your stream"
            />
          </Grid>
          
          <Grid item xs={12}>
            <TextField
              label="Description"
              fullWidth
              multiline
              rows={4}
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              variant="outlined"
              placeholder="Describe what your stream will be about"
            />
          </Grid>
          
          <Grid item xs={12} sm={6}>
            <TextField
              label="Thumbnail URL"
              fullWidth
              value={thumbnailUrl}
              onChange={(e) => setThumbnailUrl(e.target.value)}
              variant="outlined"
              placeholder="Enter URL for stream thumbnail"
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton edge="end">
                      <PhotoIcon />
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          
          <Grid item xs={12} sm={6}>
            <LocalizationProvider dateAdapter={AdapterDateFns}>
              <DateTimePicker
                label="Scheduled Start Time"
                value={scheduledStart}
                onChange={(newValue) => newValue && setScheduledStart(newValue)}
                slotProps={{
                  textField: {
                    fullWidth: true,
                    required: true,
                    variant: "outlined",
                    InputProps: {
                      endAdornment: (
                        <InputAdornment position="end">
                          <CalendarIcon />
                        </InputAdornment>
                      ),
                    }
                  }
                }}
              />
            </LocalizationProvider>
          </Grid>
          
          <Grid item xs={12} sm={6}>
            <TextField
              label="Categories"
              fullWidth
              value={categories}
              onChange={(e) => setCategories(e.target.value)}
              variant="outlined"
              placeholder="e.g. Education, Literature, Discussion"
            />
          </Grid>
          
          <Grid item xs={12} sm={6}>
            <TextField
              label="Add Tags"
              fullWidth
              value={currentTag}
              onChange={(e) => setCurrentTag(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleAddTag()}
              variant="outlined"
              placeholder="Enter tags and press Enter"
              InputProps={{
                endAdornment: (
                  <InputAdornment position="end">
                    <Button onClick={handleAddTag} disabled={!currentTag.trim()}>
                      Add
                    </Button>
                  </InputAdornment>
                ),
              }}
            />
          </Grid>
          
          <Grid item xs={12}>
            <Box display="flex" flexWrap="wrap" gap={1}>
              {tags.map((tag) => (
                <Chip
                  key={tag}
                  label={tag}
                  onDelete={() => handleDeleteTag(tag)}
                  color="primary"
                  variant="outlined"
                />
              ))}
            </Box>
          </Grid>
          
          <Grid item xs={12}>
            <FormControlLabel
              control={
                <Switch
                  checked={isPrivate}
                  onChange={(e) => setIsPrivate(e.target.checked)}
                  color="primary"
                />
              }
              label="Private Stream (only accessible with direct link)"
            />
          </Grid>
        </Grid>
      </DialogContent>
      
      <DialogActions>
        <Button onClick={onClose} disabled={loading}>
          Cancel
        </Button>
        <Button 
          onClick={handleSubmit} 
          variant="contained" 
          color="primary" 
          disabled={loading || !title || !scheduledStart}
        >
          {loading ? <CircularProgress size={24} /> : 'Create Stream'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default CreateStreamModal;

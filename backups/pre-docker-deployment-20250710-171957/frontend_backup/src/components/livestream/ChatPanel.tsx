import React, { useState, useRef, useEffect } from 'react';
import { useSelector } from 'react-redux';
import {
  Box,
  Typography,
  TextField,
  Button,
  List,
  ListItem,
  ListItemAvatar,
  ListItemText,
  Avatar,
  Divider,
  Paper,
  CircularProgress,
  Chip
} from '@mui/material';
import {
  Send as SendIcon,
  Person as PersonIcon,
  SignalWifi4Bar as ConnectedIcon,
  SignalWifiOff as DisconnectedIcon
} from '@mui/icons-material';
import { RootState } from '../../store';
import { useNavigate } from 'react-router-dom';

interface ChatPanelProps {
  messages: any[];
  onSendMessage: (message: string) => void;
  isConnected: boolean;
  isLive: boolean;
}

const ChatPanel: React.FC<ChatPanelProps> = ({ messages, onSendMessage, isConnected, isLive }) => {
  const navigate = useNavigate();
  const { user } = useSelector((state: RootState) => state.auth);
  
  const [message, setMessage] = useState('');
  const [isSending, setIsSending] = useState(false);
  
  const messagesEndRef = useRef<HTMLDivElement>(null);
  
  // Scroll to bottom when new messages arrive
  useEffect(() => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, [messages]);
  
  const handleSendMessage = () => {
    if (!user) {
      navigate('/login');
      return;
    }
    
    if (!message.trim() || !isConnected || !isLive) return;
    
    setIsSending(true);
    onSendMessage(message.trim());
    setMessage('');
    
    // Simulate network delay
    setTimeout(() => {
      setIsSending(false);
    }, 300);
  };
  
  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };
  
  return (
    <Box sx={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      {/* Connection status */}
      <Box 
        display="flex" 
        alignItems="center" 
        justifyContent="center" 
        py={0.5}
        bgcolor={isConnected ? 'success.light' : 'error.light'}
        color="white"
        borderRadius={1}
        mb={1}
      >
        {isConnected ? (
          <>
            <ConnectedIcon fontSize="small" sx={{ mr: 0.5 }} />
            <Typography variant="caption">Connected</Typography>
          </>
        ) : (
          <>
            <DisconnectedIcon fontSize="small" sx={{ mr: 0.5 }} />
            <Typography variant="caption">Disconnected</Typography>
          </>
        )}
      </Box>
      
      {/* Messages list */}
      <Paper 
        variant="outlined" 
        sx={{ 
          flexGrow: 1, 
          overflow: 'auto', 
          mb: 2,
          p: 1
        }}
      >
        {messages.length > 0 ? (
          <List disablePadding>
            {messages.map((msg, index) => (
              <React.Fragment key={index}>
                <ListItem alignItems="flex-start" dense>
                  <ListItemAvatar>
                    <Avatar sx={{ width: 32, height: 32 }}>
                      <PersonIcon fontSize="small" />
                    </Avatar>
                  </ListItemAvatar>
                  <ListItemText
                    primary={
                      <Box display="flex" alignItems="center">
                        <Typography variant="body2" fontWeight="bold" component="span">
                          {msg.from === user?.id ? 'You' : `User #${msg.from}`}
                        </Typography>
                        {msg.from === user?.id && (
                          <Chip 
                            label="You" 
                            size="small" 
                            color="primary" 
                            variant="outlined"
                            sx={{ ml: 1, height: 20, fontSize: '0.6rem' }}
                          />
                        )}
                      </Box>
                    }
                    secondary={
                      <>
                        <Typography component="span" variant="body2">
                          {msg.content}
                        </Typography>
                        <Typography component="div" variant="caption" color="text.secondary">
                          {new Date(msg.time).toLocaleTimeString()}
                        </Typography>
                      </>
                    }
                  />
                </ListItem>
                <Divider variant="inset" component="li" />
              </React.Fragment>
            ))}
            <div ref={messagesEndRef} />
          </List>
        ) : (
          <Box display="flex" flexDirection="column" alignItems="center" justifyContent="center" height="100%">
            <Typography variant="body1" color="text.secondary">
              No messages yet
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Be the first to say hello!
            </Typography>
          </Box>
        )}
      </Paper>
      
      {/* Message input */}
      <Box>
        <TextField
          label="Type a message"
          fullWidth
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyPress={handleKeyPress}
          variant="outlined"
          size="small"
          disabled={!isConnected || !isLive}
          InputProps={{
            endAdornment: (
              <Button
                variant="contained"
                color="primary"
                size="small"
                endIcon={isSending ? <CircularProgress size={16} /> : <SendIcon />}
                disabled={!message.trim() || isSending || !isConnected || !isLive}
                onClick={handleSendMessage}
                sx={{ ml: 1 }}
              >
                Send
              </Button>
            ),
          }}
        />
        
        {!isLive && (
          <Typography variant="caption" color="text.secondary" sx={{ mt: 0.5, display: 'block' }}>
            Chat is only available during live streams
          </Typography>
        )}
      </Box>
    </Box>
  );
};

export default ChatPanel;

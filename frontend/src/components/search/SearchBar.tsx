import React, { useState, useEffect, useRef } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import {
  Box,
  TextField,
  InputAdornment,
  IconButton,
  Paper,
  List,
  ListItem,
  ListItemText,
  Typography,
  Chip,
  Divider,
  Popper,
  Grow,
  ClickAwayListener,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import {
  Search as SearchIcon,
  Close as CloseIcon,
  History as HistoryIcon,
  TrendingUp as TrendingUpIcon,
} from '@mui/icons-material';
import { useNavigate } from 'react-router-dom';
import { AppDispatch } from '../../store';
import {
  setQuery,
  clearSearchResults,
  fetchRecentSearches,
  selectQuery,
  selectRecentSearches,
  selectPopularSearches,
} from '../../features/search/searchSlice';

interface SearchBarProps {
  variant?: 'standard' | 'outlined' | 'filled';
  size?: 'small' | 'medium';
  fullWidth?: boolean;
  placeholder?: string;
  autoFocus?: boolean;
  onSearch?: (query: string) => void;
}

const SearchBar: React.FC<SearchBarProps> = ({
  variant = 'outlined',
  size = 'medium',
  fullWidth = false,
  placeholder = 'Search for books, courses, tutorials...',
  autoFocus = false,
  onSearch,
}) => {
  const dispatch = useDispatch<AppDispatch>();
  const navigate = useNavigate();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  
  const query = useSelector(selectQuery);
  const recentSearches = useSelector(selectRecentSearches);
  const popularSearches = useSelector(selectPopularSearches);
  
  const [open, setOpen] = useState(false);
  const [inputValue, setInputValue] = useState(query);
  const anchorRef = useRef<HTMLDivElement>(null);
  
  useEffect(() => {
    dispatch(fetchRecentSearches());
  }, [dispatch]);
  
  useEffect(() => {
    setInputValue(query);
  }, [query]);
  
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setInputValue(e.target.value);
  };
  
  const handleInputFocus = () => {
    setOpen(true);
  };
  
  const handleClickAway = () => {
    setOpen(false);
  };
  
  const handleClear = () => {
    setInputValue('');
    dispatch(setQuery(''));
    dispatch(clearSearchResults());
  };
  
  const handleSearch = (searchQuery: string = inputValue) => {
    if (searchQuery.trim()) {
      dispatch(setQuery(searchQuery));
      dispatch(clearSearchResults());
      
      if (onSearch) {
        onSearch(searchQuery);
      } else {
        navigate(`/search?q=${encodeURIComponent(searchQuery)}`);
      }
      
      setOpen(false);
    }
  };
  
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };
  
  const handleSuggestionClick = (suggestion: string) => {
    handleSearch(suggestion);
  };
  
  return (
    <Box ref={anchorRef} sx={{ position: 'relative', width: fullWidth ? '100%' : 'auto' }}>
      <TextField
        variant={variant}
        size={size}
        fullWidth={fullWidth}
        placeholder={placeholder}
        value={inputValue}
        onChange={handleInputChange}
        onFocus={handleInputFocus}
        onKeyDown={handleKeyDown}
        autoFocus={autoFocus}
        InputProps={{
          startAdornment: (
            <InputAdornment position="start">
              <SearchIcon color="action" />
            </InputAdornment>
          ),
          endAdornment: inputValue ? (
            <InputAdornment position="end">
              <IconButton
                aria-label="clear search"
                onClick={handleClear}
                edge="end"
                size="small"
              >
                <CloseIcon fontSize="small" />
              </IconButton>
            </InputAdornment>
          ) : null,
          sx: {
            borderRadius: 2,
            transition: theme.transitions.create(['width', 'box-shadow']),
            ...(open && {
              boxShadow: '0 4px 6px rgba(0, 0, 0, 0.1)',
            }),
          },
        }}
        sx={{
          width: fullWidth ? '100%' : (open && !isMobile ? 400 : 240),
          transition: theme.transitions.create('width'),
        }}
      />
      
      <Popper
        open={open}
        anchorEl={anchorRef.current}
        placement="bottom-start"
        transition
        disablePortal
        style={{
          width: fullWidth ? '100%' : (isMobile ? '100%' : 400),
          zIndex: theme.zIndex.modal,
        }}
      >
        {({ TransitionProps }) => (
          <Grow {...TransitionProps} style={{ transformOrigin: 'top left' }}>
            <Paper
              elevation={4}
              sx={{
                mt: 0.5,
                borderRadius: 2,
                maxHeight: 400,
                overflow: 'auto',
              }}
            >
              <ClickAwayListener onClickAway={handleClickAway}>
                <Box>
                  {recentSearches.length > 0 && (
                    <>
                      <Box px={2} py={1} display="flex" alignItems="center">
                        <HistoryIcon fontSize="small" sx={{ mr: 1, color: 'text.secondary' }} />
                        <Typography variant="subtitle2" color="text.secondary">
                          Recent Searches
                        </Typography>
                      </Box>
                      <List dense disablePadding>
                        {recentSearches.slice(0, 5).map((search, index) => (
                          <ListItem
                            key={`recent-${index}`}
                            button
                            onClick={() => handleSuggestionClick(search)}
                            sx={{
                              py: 0.5,
                              px: 2,
                              '&:hover': {
                                backgroundColor: 'action.hover',
                              },
                            }}
                          >
                            <ListItemText primary={search} />
                          </ListItem>
                        ))}
                      </List>
                      <Divider />
                    </>
                  )}
                  
                  {popularSearches.length > 0 && (
                    <>
                      <Box px={2} py={1} display="flex" alignItems="center">
                        <TrendingUpIcon fontSize="small" sx={{ mr: 1, color: 'text.secondary' }} />
                        <Typography variant="subtitle2" color="text.secondary">
                          Popular Searches
                        </Typography>
                      </Box>
                      <Box px={2} py={1} display="flex" flexWrap="wrap" gap={1}>
                        {popularSearches.slice(0, 8).map((item, index) => (
                          <Chip
                            key={`popular-${index}`}
                            label={item.query}
                            size="small"
                            onClick={() => handleSuggestionClick(item.query)}
                            sx={{ cursor: 'pointer' }}
                          />
                        ))}
                      </Box>
                      <Divider />
                    </>
                  )}
                  
                  {inputValue && (
                    <Box p={2}>
                      <Typography variant="body2">
                        Press Enter to search for "{inputValue}"
                      </Typography>
                    </Box>
                  )}
                </Box>
              </ClickAwayListener>
            </Paper>
          </Grow>
        )}
      </Popper>
    </Box>
  );
};

export default SearchBar;

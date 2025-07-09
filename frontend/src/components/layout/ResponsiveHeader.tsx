import React, { useState } from 'react';
import { Link as RouterLink } from 'react-router-dom';
import {
  AppBar,
  Box,
  Toolbar,
  IconButton,
  Typography,
  Menu,
  MenuItem,
  Container,
  Avatar,
  Button,
  Tooltip,
  Drawer,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Divider,
  Badge,
  useTheme,
  useMediaQuery,
} from '@mui/material';
import {
  Menu as MenuIcon,
  Notifications as NotificationsIcon,
  AccountCircle,
  Home as HomeIcon,
  Book as BookIcon,
  Forum as ForumIcon,
  School as SchoolIcon,
  Store as StoreIcon,
  LiveTv as LiveTvIcon,
  Wallet as WalletIcon,
  Settings as SettingsIcon,
  Logout as LogoutIcon,
} from '@mui/icons-material';
import { useSelector, useDispatch } from 'react-redux';
import { logout, selectUser, selectIsAuthenticated } from '../../features/auth/authSlice';
import { AppDispatch } from '../../store';
import SearchBar from '../search/SearchBar';
import ThemeToggle from '../theme/ThemeToggle';
import Logo from '../common/Logo';

const ResponsiveHeader: React.FC = () => {
  const dispatch = useDispatch<AppDispatch>();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const isSmall = useMediaQuery(theme.breakpoints.down('sm'));
  
  const isAuthenticated = useSelector(selectIsAuthenticated);
  const user = useSelector(selectUser);
  
  const [anchorElUser, setAnchorElUser] = useState<null | HTMLElement>(null);
  const [drawerOpen, setDrawerOpen] = useState(false);
  
  const handleOpenUserMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorElUser(event.currentTarget);
  };
  
  const handleCloseUserMenu = () => {
    setAnchorElUser(null);
  };
  
  const handleDrawerToggle = () => {
    setDrawerOpen(!drawerOpen);
  };
  
  const handleLogout = () => {
    dispatch(logout());
    handleCloseUserMenu();
  };
  
  const pages = [
    { name: 'Books', path: '/books', icon: <BookIcon /> },
    { name: 'Forum', path: '/forum', icon: <ForumIcon /> },
    { name: 'Courses', path: '/courses', icon: <SchoolIcon /> },
    { name: 'Marketplace', path: '/marketplace', icon: <StoreIcon /> },
    { name: 'Livestream', path: '/livestream', icon: <LiveTvIcon /> },
  ];
  
  const userMenuItems = [
    { name: 'Profile', path: '/profile', icon: <AccountCircle /> },
    { name: 'Wallet', path: '/wallet', icon: <WalletIcon /> },
    { name: 'Settings', path: '/settings', icon: <SettingsIcon /> },
  ];
  
  const drawer = (
    <Box sx={{ width: 250 }} role="presentation">
      <Box sx={{ p: 2, display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
        <Logo height={40} />
      </Box>
      <Divider />
      <List>
        <ListItem button component={RouterLink} to="/" onClick={handleDrawerToggle}>
          <ListItemIcon>
            <HomeIcon />
          </ListItemIcon>
          <ListItemText primary="Home" />
        </ListItem>
        {pages.map((page) => (
          <ListItem
            button
            key={page.name}
            component={RouterLink}
            to={page.path}
            onClick={handleDrawerToggle}
          >
            <ListItemIcon>{page.icon}</ListItemIcon>
            <ListItemText primary={page.name} />
          </ListItem>
        ))}
      </List>
      <Divider />
      {isAuthenticated ? (
        <List>
          {userMenuItems.map((item) => (
            <ListItem
              button
              key={item.name}
              component={RouterLink}
              to={item.path}
              onClick={handleDrawerToggle}
            >
              <ListItemIcon>{item.icon}</ListItemIcon>
              <ListItemText primary={item.name} />
            </ListItem>
          ))}
          <ListItem button onClick={handleLogout}>
            <ListItemIcon>
              <LogoutIcon />
            </ListItemIcon>
            <ListItemText primary="Logout" />
          </ListItem>
        </List>
      ) : (
        <List>
          <ListItem button component={RouterLink} to="/login" onClick={handleDrawerToggle}>
            <ListItemText primary="Login" />
          </ListItem>
          <ListItem button component={RouterLink} to="/register" onClick={handleDrawerToggle}>
            <ListItemText primary="Register" />
          </ListItem>
        </List>
      )}
    </Box>
  );
  
  return (
    <AppBar position="sticky" color="default" elevation={1}>
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          {/* Mobile view */}
          <Box sx={{ display: { xs: 'flex', md: 'none' }, alignItems: 'center' }}>
            <IconButton
              size="large"
              aria-label="menu"
              aria-controls="menu-appbar"
              aria-haspopup="true"
              onClick={handleDrawerToggle}
              color="inherit"
              edge="start"
              sx={{ mr: 1 }}
            >
              <MenuIcon />
            </IconButton>
            <Box sx={{ flexGrow: 1, display: 'flex', justifyContent: 'center' }}>
              <RouterLink to="/">
                <Logo height={40} />
              </RouterLink>
            </Box>
          </Box>
          
          {/* Desktop view */}
          <Box sx={{ display: { xs: 'none', md: 'flex' }, alignItems: 'center', mr: 2 }}>
            <RouterLink to="/">
              <Logo height={40} />
            </RouterLink>
          </Box>
          
          <Box sx={{ display: { xs: 'none', md: 'flex' }, flexGrow: 1 }}>
            {pages.map((page) => (
              <Button
                key={page.name}
                component={RouterLink}
                to={page.path}
                sx={{ mx: 1 }}
              >
                {page.name}
              </Button>
            ))}
          </Box>
          
          {/* Search bar */}
          {!isSmall && (
            <Box sx={{ flexGrow: 0, mr: 2 }}>
              <SearchBar size="small" />
            </Box>
          )}
          
          {/* Theme toggle */}
          <Box sx={{ mr: 1 }}>
            <ThemeToggle size="small" />
          </Box>
          
          {/* User menu */}
          {isAuthenticated ? (
            <Box sx={{ flexGrow: 0 }}>
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <Tooltip title="Notifications">
                  <IconButton sx={{ mr: 1 }}>
                    <Badge badgeContent={3} color="error">
                      <NotificationsIcon />
                    </Badge>
                  </IconButton>
                </Tooltip>
                
                <Tooltip title="Open settings">
                  <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                    <Avatar
                      alt={user?.name || 'User'}
                      src={user?.avatar || ''}
                      sx={{ width: 32, height: 32 }}
                    />
                  </IconButton>
                </Tooltip>
              </Box>
              <Menu
                sx={{ mt: '45px' }}
                id="menu-appbar"
                anchorEl={anchorElUser}
                anchorOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
                keepMounted
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
                open={Boolean(anchorElUser)}
                onClose={handleCloseUserMenu}
              >
                {userMenuItems.map((item) => (
                  <MenuItem
                    key={item.name}
                    component={RouterLink}
                    to={item.path}
                    onClick={handleCloseUserMenu}
                  >
                    <ListItemIcon>{item.icon}</ListItemIcon>
                    <Typography textAlign="center">{item.name}</Typography>
                  </MenuItem>
                ))}
                <Divider />
                <MenuItem onClick={handleLogout}>
                  <ListItemIcon>
                    <LogoutIcon fontSize="small" />
                  </ListItemIcon>
                  <Typography textAlign="center">Logout</Typography>
                </MenuItem>
              </Menu>
            </Box>
          ) : (
            <Box sx={{ flexGrow: 0, display: 'flex' }}>
              <Button
                component={RouterLink}
                to="/login"
                variant="outlined"
                sx={{ mr: 1, display: { xs: 'none', sm: 'block' } }}
              >
                Login
              </Button>
              <Button
                component={RouterLink}
                to="/register"
                variant="contained"
              >
                Register
              </Button>
            </Box>
          )}
        </Toolbar>
      </Container>
      
      {/* Mobile drawer */}
      <Drawer
        anchor="left"
        open={drawerOpen}
        onClose={handleDrawerToggle}
        ModalProps={{
          keepMounted: true, // Better open performance on mobile
        }}
      >
        {drawer}
      </Drawer>
    </AppBar>
  );
};

export default ResponsiveHeader;

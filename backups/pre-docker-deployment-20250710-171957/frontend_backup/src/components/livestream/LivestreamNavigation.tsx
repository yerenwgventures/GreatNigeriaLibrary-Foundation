import React from 'react';
import { useSelector } from 'react-redux';
import { useNavigate, useLocation } from 'react-router-dom';
import {
  Box,
  Tabs,
  Tab,
  Button,
  Divider,
  useMediaQuery,
  useTheme,
  Menu,
  MenuItem,
  IconButton
} from '@mui/material';
import {
  Videocam as VideoIcon,
  CardGiftcard as GiftIcon,
  MonetizationOn as RevenueIcon,
  Person as ProfileIcon,
  MoreVert as MoreIcon,
  AdminPanelSettings as AdminIcon,
  Dashboard as DashboardIcon
} from '@mui/icons-material';
import { RootState } from '../../store';

const LivestreamNavigation: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));

  const { user } = useSelector((state: RootState) => state.auth);
  const isAdmin = user && user.role === 'admin';

  const [mobileMenuAnchor, setMobileMenuAnchor] = React.useState<null | HTMLElement>(null);

  // Determine active tab based on current path
  const getActiveTab = () => {
    const path = location.pathname;

    if (path === '/livestream') return 0;
    if (path.startsWith('/livestream/')) return 0;
    if (path === '/coins') return 1;
    if (path === '/creator/revenue') return 2;
    if (path === '/creator/dashboard') return 3;
    if (path.startsWith('/streamer/')) return 4;
    if (path === '/admin/livestream') return 5;

    return 0;
  };

  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    switch (newValue) {
      case 0:
        navigate('/livestream');
        break;
      case 1:
        navigate('/coins');
        break;
      case 2:
        navigate('/creator/revenue');
        break;
      case 3:
        navigate('/creator/dashboard');
        break;
      case 4:
        navigate(user ? `/streamer/${user.id}` : '/login');
        break;
      case 5:
        navigate('/admin/livestream');
        break;
      default:
        navigate('/livestream');
    }
  };

  const handleMobileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setMobileMenuAnchor(event.currentTarget);
  };

  const handleMobileMenuClose = () => {
    setMobileMenuAnchor(null);
  };

  const handleMobileMenuItemClick = (path: string) => {
    navigate(path);
    handleMobileMenuClose();
  };

  // If user is not logged in, show login button
  if (!user) {
    return (
      <Box sx={{ mb: 3, display: 'flex', justifyContent: 'center' }}>
        <Button
          variant="contained"
          color="primary"
          onClick={() => navigate('/login')}
        >
          Log in to access all features
        </Button>
      </Box>
    );
  }

  return (
    <Box sx={{ width: '100%', mb: 3 }}>
      {isMobile ? (
        // Mobile view
        <Box display="flex" justifyContent="space-between" alignItems="center">
          <Tabs
            value={getActiveTab()}
            onChange={handleTabChange}
            variant="scrollable"
            scrollButtons="auto"
            sx={{ flexGrow: 1 }}
          >
            <Tab icon={<VideoIcon />} label="Streams" />
          </Tabs>

          <IconButton
            aria-label="more"
            aria-controls="livestream-menu"
            aria-haspopup="true"
            onClick={handleMobileMenuOpen}
          >
            <MoreIcon />
          </IconButton>

          <Menu
            id="livestream-menu"
            anchorEl={mobileMenuAnchor}
            keepMounted
            open={Boolean(mobileMenuAnchor)}
            onClose={handleMobileMenuClose}
          >
            <MenuItem onClick={() => handleMobileMenuItemClick('/coins')}>
              <GiftIcon sx={{ mr: 1 }} /> Coins
            </MenuItem>
            <MenuItem onClick={() => handleMobileMenuItemClick('/creator/revenue')}>
              <RevenueIcon sx={{ mr: 1 }} /> Revenue
            </MenuItem>
            <MenuItem onClick={() => handleMobileMenuItemClick('/creator/dashboard')}>
              <DashboardIcon sx={{ mr: 1 }} /> Dashboard
            </MenuItem>
            <MenuItem onClick={() => handleMobileMenuItemClick(`/streamer/${user.id}`)}>
              <ProfileIcon sx={{ mr: 1 }} /> My Profile
            </MenuItem>
            {isAdmin && (
              <MenuItem onClick={() => handleMobileMenuItemClick('/admin/livestream')}>
                <AdminIcon sx={{ mr: 1 }} /> Admin
              </MenuItem>
            )}
          </Menu>
        </Box>
      ) : (
        // Desktop view
        <Tabs
          value={getActiveTab()}
          onChange={handleTabChange}
          variant="fullWidth"
        >
          <Tab icon={<VideoIcon />} label="Streams" />
          <Tab icon={<GiftIcon />} label="Coins" />
          <Tab icon={<RevenueIcon />} label="Revenue" />
          <Tab icon={<DashboardIcon />} label="Dashboard" />
          <Tab icon={<ProfileIcon />} label="My Profile" />
          {isAdmin && <Tab icon={<AdminIcon />} label="Admin" />}
        </Tabs>
      )}

      <Divider sx={{ mt: 1 }} />
    </Box>
  );
};

export default LivestreamNavigation;

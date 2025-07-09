import React from 'react';
import { Box, Typography, useTheme } from '@mui/material';
import { SvgIcon } from '@mui/material';

interface LogoProps {
  height?: number;
  showText?: boolean;
}

const Logo: React.FC<LogoProps> = ({ height = 40, showText = true }) => {
  const theme = useTheme();
  const primaryColor = theme.palette.primary.main;
  
  return (
    <Box sx={{ display: 'flex', alignItems: 'center' }}>
      <SvgIcon
        sx={{
          height: `${height}px`,
          width: `${height}px`,
          mr: showText ? 1 : 0,
        }}
        viewBox="0 0 24 24"
      >
        <path
          fill={primaryColor}
          d="M12,2L1,21H23L12,2M12,6L19.53,19H4.47L12,6M11,10V14H13V10H11M11,16V18H13V16H11Z"
        />
      </SvgIcon>
      
      {showText && (
        <Typography
          variant="h6"
          component="div"
          sx={{
            fontWeight: 700,
            fontSize: `${height * 0.5}px`,
            lineHeight: 1,
            color: theme.palette.text.primary,
            display: 'flex',
            flexDirection: 'column',
          }}
        >
          <Box component="span">Great Nigeria</Box>
          <Box
            component="span"
            sx={{
              fontSize: '0.7em',
              fontWeight: 500,
              color: theme.palette.text.secondary,
            }}
          >
            Library
          </Box>
        </Typography>
      )}
    </Box>
  );
};

export default Logo;

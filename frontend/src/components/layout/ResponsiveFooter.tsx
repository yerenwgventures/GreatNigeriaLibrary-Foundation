import React from 'react';
import { Link as RouterLink } from 'react-router-dom';
import {
  Box,
  Container,
  Grid,
  Typography,
  Link,
  Divider,
  IconButton,
  useTheme,
  useMediaQuery,
  Accordion,
  AccordionSummary,
  AccordionDetails,
} from '@mui/material';
import {
  Facebook as FacebookIcon,
  Twitter as TwitterIcon,
  Instagram as InstagramIcon,
  YouTube as YouTubeIcon,
  LinkedIn as LinkedInIcon,
  ExpandMore as ExpandMoreIcon,
} from '@mui/icons-material';
import Logo from '../common/Logo';

const ResponsiveFooter: React.FC = () => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  
  const footerLinks = [
    {
      title: 'About',
      links: [
        { name: 'Our Mission', path: '/about/mission' },
        { name: 'Our Team', path: '/about/team' },
        { name: 'Our Partners', path: '/about/partners' },
        { name: 'Careers', path: '/about/careers' },
        { name: 'Contact Us', path: '/contact' },
      ],
    },
    {
      title: 'Resources',
      links: [
        { name: 'Books', path: '/books' },
        { name: 'Courses', path: '/courses' },
        { name: 'Tutorials', path: '/tutorials' },
        { name: 'Forum', path: '/forum' },
        { name: 'Events', path: '/events' },
      ],
    },
    {
      title: 'Community',
      links: [
        { name: 'Marketplace', path: '/marketplace' },
        { name: 'Livestream', path: '/livestream' },
        { name: 'Local Groups', path: '/groups' },
        { name: 'Volunteer', path: '/volunteer' },
        { name: 'Donate', path: '/donate' },
      ],
    },
    {
      title: 'Legal',
      links: [
        { name: 'Terms of Service', path: '/legal/terms' },
        { name: 'Privacy Policy', path: '/legal/privacy' },
        { name: 'Cookie Policy', path: '/legal/cookies' },
        { name: 'Copyright', path: '/legal/copyright' },
        { name: 'Accessibility', path: '/legal/accessibility' },
      ],
    },
  ];
  
  const socialLinks = [
    { name: 'Facebook', icon: <FacebookIcon />, url: 'https://facebook.com' },
    { name: 'Twitter', icon: <TwitterIcon />, url: 'https://twitter.com' },
    { name: 'Instagram', icon: <InstagramIcon />, url: 'https://instagram.com' },
    { name: 'YouTube', icon: <YouTubeIcon />, url: 'https://youtube.com' },
    { name: 'LinkedIn', icon: <LinkedInIcon />, url: 'https://linkedin.com' },
  ];
  
  return (
    <Box
      component="footer"
      sx={{
        py: 6,
        px: 2,
        mt: 'auto',
        backgroundColor: theme.palette.mode === 'light' ? 'grey.100' : 'grey.900',
      }}
    >
      <Container maxWidth="lg">
        <Grid container spacing={4} justifyContent="space-between">
          <Grid item xs={12} md={4}>
            <Box display="flex" flexDirection="column" alignItems={{ xs: 'center', md: 'flex-start' }}>
              <RouterLink to="/">
                <Logo height={50} />
              </RouterLink>
              <Typography variant="body2" color="text.secondary" sx={{ mt: 2, textAlign: { xs: 'center', md: 'left' } }}>
                Empowering Nigerians through knowledge, community, and resources to build a better future.
              </Typography>
              <Box sx={{ mt: 2, display: 'flex', gap: 1 }}>
                {socialLinks.map((social) => (
                  <IconButton
                    key={social.name}
                    aria-label={social.name}
                    component="a"
                    href={social.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    size="small"
                  >
                    {social.icon}
                  </IconButton>
                ))}
              </Box>
            </Box>
          </Grid>
          
          {isMobile ? (
            // Mobile accordion view
            <Grid item xs={12}>
              {footerLinks.map((section) => (
                <Accordion key={section.title} elevation={0} sx={{ backgroundColor: 'transparent' }}>
                  <AccordionSummary
                    expandIcon={<ExpandMoreIcon />}
                    aria-controls={`${section.title}-content`}
                    id={`${section.title}-header`}
                  >
                    <Typography variant="subtitle1" fontWeight="bold">
                      {section.title}
                    </Typography>
                  </AccordionSummary>
                  <AccordionDetails>
                    <Box display="flex" flexDirection="column" gap={1}>
                      {section.links.map((link) => (
                        <Link
                          key={link.name}
                          component={RouterLink}
                          to={link.path}
                          color="text.secondary"
                          underline="hover"
                        >
                          {link.name}
                        </Link>
                      ))}
                    </Box>
                  </AccordionDetails>
                </Accordion>
              ))}
            </Grid>
          ) : (
            // Desktop columns view
            footerLinks.map((section) => (
              <Grid item xs={6} md={2} key={section.title}>
                <Typography variant="subtitle1" fontWeight="bold" gutterBottom>
                  {section.title}
                </Typography>
                <Box display="flex" flexDirection="column" gap={1}>
                  {section.links.map((link) => (
                    <Link
                      key={link.name}
                      component={RouterLink}
                      to={link.path}
                      color="text.secondary"
                      underline="hover"
                    >
                      {link.name}
                    </Link>
                  ))}
                </Box>
              </Grid>
            ))
          )}
        </Grid>
        
        <Divider sx={{ my: 4 }} />
        
        <Box
          display="flex"
          flexDirection={{ xs: 'column', sm: 'row' }}
          justifyContent="space-between"
          alignItems="center"
          gap={2}
        >
          <Typography variant="body2" color="text.secondary" align={isMobile ? 'center' : 'left'}>
            © {new Date().getFullYear()} Great Nigeria Library. All rights reserved.
          </Typography>
          <Box
            display="flex"
            flexDirection={{ xs: 'column', sm: 'row' }}
            gap={2}
            alignItems="center"
          >
            <Link
              component={RouterLink}
              to="/sitemap"
              color="text.secondary"
              underline="hover"
              variant="body2"
            >
              Sitemap
            </Link>
            <Link
              component="a"
              href="mailto:info@greatnigeria.org"
              color="text.secondary"
              underline="hover"
              variant="body2"
            >
              info@greatnigeria.org
            </Link>
          </Box>
        </Box>
      </Container>
    </Box>
  );
};

export default ResponsiveFooter;

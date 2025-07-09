import React, { useState } from 'react';
import {
  Box,
  Typography,
  Paper,
  Grid,
  Card,
  CardContent,
  Divider,
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  Button,
  IconButton,
  Tooltip,
  CircularProgress
} from '@mui/material';
import {
  BarChart as ChartIcon,
  TrendingUp as TrendingUpIcon,
  TrendingDown as TrendingDownIcon,
  PieChart as PieChartIcon,
  Timeline as TimelineIcon,
  Download as DownloadIcon,
  Refresh as RefreshIcon,
  Info as InfoIcon
} from '@mui/icons-material';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip as RechartsTooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
  LineChart,
  Line,
  AreaChart,
  Area
} from 'recharts';

interface StreamAnalyticsProps {
  streamId: number;
  creatorId: number;
}

interface TabPanelProps {
  children?: React.ReactNode;
  index: number;
  value: number;
}

function TabPanel(props: TabPanelProps) {
  const { children, value, index, ...other } = props;

  return (
    <div
      role="tabpanel"
      hidden={value !== index}
      id={`analytics-tabpanel-${index}`}
      aria-controls={`analytics-tab-${index}`}
      {...other}
    >
      {value === index && (
        <Box sx={{ p: 2 }}>
          {children}
        </Box>
      )}
    </div>
  );
}

// Sample data for demonstration
const SAMPLE_VIEWER_DATA = [
  { time: '00:00', viewers: 10 },
  { time: '05:00', viewers: 25 },
  { time: '10:00', viewers: 45 },
  { time: '15:00', viewers: 80 },
  { time: '20:00', viewers: 120 },
  { time: '25:00', viewers: 95 },
  { time: '30:00', viewers: 110 },
  { time: '35:00', viewers: 130 },
  { time: '40:00', viewers: 150 },
  { time: '45:00', viewers: 145 },
  { time: '50:00', viewers: 160 },
  { time: '55:00', viewers: 175 },
  { time: '60:00', viewers: 190 }
];

const SAMPLE_GIFT_DATA = [
  { time: '00:00', gifts: 0, value: 0 },
  { time: '05:00', gifts: 2, value: 60 },
  { time: '10:00', gifts: 5, value: 150 },
  { time: '15:00', gifts: 8, value: 400 },
  { time: '20:00', gifts: 12, value: 1200 },
  { time: '25:00', gifts: 7, value: 700 },
  { time: '30:00', gifts: 10, value: 1000 },
  { time: '35:00', gifts: 15, value: 1500 },
  { time: '40:00', gifts: 20, value: 2500 },
  { time: '45:00', gifts: 18, value: 2000 },
  { time: '50:00', gifts: 25, value: 3000 },
  { time: '55:00', gifts: 30, value: 4000 },
  { time: '60:00', gifts: 35, value: 5000 }
];

const SAMPLE_GIFT_TYPES = [
  { name: 'Heart', value: 500, count: 50 },
  { name: 'Star', value: 1250, count: 25 },
  { name: 'Trophy', value: 2000, count: 20 },
  { name: 'Diamond', value: 5000, count: 10 },
  { name: 'Crown', value: 8000, count: 8 },
  { name: 'Rocket', value: 15000, count: 3 }
];

const SAMPLE_TOP_GIFTERS = [
  { rank: 1, userId: 123, username: 'User #123', gifts: 25, value: 12500 },
  { rank: 2, userId: 456, username: 'User #456', gifts: 18, value: 9000 },
  { rank: 3, userId: 789, username: 'User #789', gifts: 15, value: 7500 },
  { rank: 4, userId: 101, username: 'User #101', gifts: 12, value: 6000 },
  { rank: 5, userId: 202, username: 'User #202', gifts: 10, value: 5000 }
];

const SAMPLE_ENGAGEMENT_DATA = [
  { time: '00:00', comments: 5, reactions: 2 },
  { time: '05:00', comments: 12, reactions: 8 },
  { time: '10:00', comments: 20, reactions: 15 },
  { time: '15:00', comments: 35, reactions: 25 },
  { time: '20:00', comments: 50, reactions: 40 },
  { time: '25:00', comments: 45, reactions: 35 },
  { time: '30:00', comments: 60, reactions: 45 },
  { time: '35:00', comments: 75, reactions: 55 },
  { time: '40:00', comments: 90, reactions: 70 },
  { time: '45:00', comments: 85, reactions: 65 },
  { time: '50:00', comments: 100, reactions: 80 },
  { time: '55:00', comments: 120, reactions: 95 },
  { time: '60:00', comments: 140, reactions: 110 }
];

const SAMPLE_RETENTION_DATA = [
  { duration: '0-1 min', viewers: 200 },
  { duration: '1-5 min', viewers: 180 },
  { duration: '5-15 min', viewers: 150 },
  { duration: '15-30 min', viewers: 120 },
  { duration: '30-60 min', viewers: 90 },
  { duration: '60+ min', viewers: 60 }
];

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884d8', '#82ca9d'];

const StreamAnalytics: React.FC<StreamAnalyticsProps> = ({ streamId, creatorId }) => {
  const [tabValue, setTabValue] = useState(0);
  const [timeRange, setTimeRange] = useState('all');
  const [isLoading, setIsLoading] = useState(false);
  
  const handleTabChange = (_event: React.SyntheticEvent, newValue: number) => {
    setTabValue(newValue);
  };
  
  const handleTimeRangeChange = (event: any) => {
    setTimeRange(event.target.value);
  };
  
  const handleRefresh = () => {
    setIsLoading(true);
    
    // Simulate API call
    setTimeout(() => {
      setIsLoading(false);
    }, 1000);
  };
  
  const handleExport = () => {
    // In a real app, this would generate and download a CSV or PDF report
    alert('Analytics data would be exported');
  };
  
  // Calculate summary metrics
  const peakViewers = Math.max(...SAMPLE_VIEWER_DATA.map(d => d.viewers));
  const totalGifts = SAMPLE_GIFT_TYPES.reduce((sum, item) => sum + item.count, 0);
  const totalGiftValue = SAMPLE_GIFT_TYPES.reduce((sum, item) => sum + item.value, 0);
  const totalComments = SAMPLE_ENGAGEMENT_DATA.reduce((sum, item) => sum + item.comments, 0);
  const totalReactions = SAMPLE_ENGAGEMENT_DATA.reduce((sum, item) => sum + item.reactions, 0);
  
  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
        <Typography variant="h5" component="h2">
          Stream Analytics
        </Typography>
        
        <Box display="flex" alignItems="center" gap={2}>
          <FormControl size="small" sx={{ minWidth: 150 }}>
            <InputLabel>Time Range</InputLabel>
            <Select
              value={timeRange}
              label="Time Range"
              onChange={handleTimeRangeChange}
            >
              <MenuItem value="all">All Time</MenuItem>
              <MenuItem value="hour">Last Hour</MenuItem>
              <MenuItem value="30min">Last 30 Minutes</MenuItem>
              <MenuItem value="15min">Last 15 Minutes</MenuItem>
            </Select>
          </FormControl>
          
          <Tooltip title="Refresh Data">
            <IconButton onClick={handleRefresh} disabled={isLoading}>
              {isLoading ? <CircularProgress size={24} /> : <RefreshIcon />}
            </IconButton>
          </Tooltip>
          
          <Tooltip title="Export Data">
            <IconButton onClick={handleExport}>
              <DownloadIcon />
            </IconButton>
          </Tooltip>
        </Box>
      </Box>
      
      {/* Summary Cards */}
      <Grid container spacing={2} sx={{ mb: 3 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                Peak Viewers
              </Typography>
              <Typography variant="h4">
                {peakViewers}
              </Typography>
              <Box display="flex" alignItems="center" mt={1}>
                <TrendingUpIcon color="success" fontSize="small" sx={{ mr: 0.5 }} />
                <Typography variant="body2" color="success.main">
                  +15% from avg
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                Total Gifts
              </Typography>
              <Typography variant="h4">
                {totalGifts}
              </Typography>
              <Box display="flex" alignItems="center" mt={1}>
                <TrendingUpIcon color="success" fontSize="small" sx={{ mr: 0.5 }} />
                <Typography variant="body2" color="success.main">
                  +23% from avg
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                Gift Revenue
              </Typography>
              <Typography variant="h4">
                ₦{totalGiftValue.toLocaleString()}
              </Typography>
              <Box display="flex" alignItems="center" mt={1}>
                <TrendingUpIcon color="success" fontSize="small" sx={{ mr: 0.5 }} />
                <Typography variant="body2" color="success.main">
                  +42% from avg
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>
        
        <Grid item xs={12} sm={6} md={3}>
          <Card>
            <CardContent>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom>
                Engagement
              </Typography>
              <Typography variant="h4">
                {totalComments + totalReactions}
              </Typography>
              <Box display="flex" alignItems="center" mt={1}>
                <Typography variant="body2" color="text.secondary">
                  {totalComments} comments, {totalReactions} reactions
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
      
      <Paper elevation={2}>
        <Tabs
          value={tabValue}
          onChange={handleTabChange}
          variant="scrollable"
          scrollButtons="auto"
          sx={{ borderBottom: 1, borderColor: 'divider' }}
        >
          <Tab icon={<ChartIcon />} label="Overview" id="analytics-tab-0" aria-controls="analytics-tabpanel-0" />
          <Tab icon={<GiftIcon />} label="Gifts" id="analytics-tab-1" aria-controls="analytics-tabpanel-1" />
          <Tab icon={<TimelineIcon />} label="Engagement" id="analytics-tab-2" aria-controls="analytics-tabpanel-2" />
          <Tab icon={<PieChartIcon />} label="Audience" id="analytics-tab-3" aria-controls="analytics-tabpanel-3" />
        </Tabs>
        
        {/* Overview Tab */}
        <TabPanel value={tabValue} index={0}>
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <Typography variant="h6" gutterBottom>
                Viewers Over Time
              </Typography>
              <Box sx={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart
                    data={SAMPLE_VIEWER_DATA}
                    margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="time" />
                    <YAxis />
                    <RechartsTooltip />
                    <Area type="monotone" dataKey="viewers" stroke="#8884d8" fill="#8884d8" fillOpacity={0.3} />
                  </AreaChart>
                </ResponsiveContainer>
              </Box>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Typography variant="h6" gutterBottom>
                Gift Value Over Time
              </Typography>
              <Box sx={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart
                    data={SAMPLE_GIFT_DATA}
                    margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="time" />
                    <YAxis />
                    <RechartsTooltip />
                    <Area type="monotone" dataKey="value" stroke="#82ca9d" fill="#82ca9d" fillOpacity={0.3} />
                  </AreaChart>
                </ResponsiveContainer>
              </Box>
            </Grid>
            
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>
                Viewer Retention
              </Typography>
              <Box sx={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={SAMPLE_RETENTION_DATA}
                    margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="duration" />
                    <YAxis />
                    <RechartsTooltip />
                    <Bar dataKey="viewers" fill="#8884d8" />
                  </BarChart>
                </ResponsiveContainer>
              </Box>
            </Grid>
          </Grid>
        </TabPanel>
        
        {/* Gifts Tab */}
        <TabPanel value={tabValue} index={1}>
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <Typography variant="h6" gutterBottom>
                Gift Distribution by Type
              </Typography>
              <Box sx={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={SAMPLE_GIFT_TYPES}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {SAMPLE_GIFT_TYPES.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <RechartsTooltip />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              </Box>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Typography variant="h6" gutterBottom>
                Gift Count vs Value
              </Typography>
              <Box sx={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={SAMPLE_GIFT_TYPES}
                    margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis yAxisId="left" orientation="left" stroke="#8884d8" />
                    <YAxis yAxisId="right" orientation="right" stroke="#82ca9d" />
                    <RechartsTooltip />
                    <Legend />
                    <Bar yAxisId="left" dataKey="count" fill="#8884d8" name="Count" />
                    <Bar yAxisId="right" dataKey="value" fill="#82ca9d" name="Value (₦)" />
                  </BarChart>
                </ResponsiveContainer>
              </Box>
            </Grid>
            
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>
                Top Gifters
              </Typography>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Rank</TableCell>
                      <TableCell>User</TableCell>
                      <TableCell align="right">Gifts Sent</TableCell>
                      <TableCell align="right">Total Value</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {SAMPLE_TOP_GIFTERS.map((gifter) => (
                      <TableRow key={gifter.userId}>
                        <TableCell>{gifter.rank}</TableCell>
                        <TableCell>{gifter.username}</TableCell>
                        <TableCell align="right">{gifter.gifts}</TableCell>
                        <TableCell align="right">₦{gifter.value.toLocaleString()}</TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              </TableContainer>
            </Grid>
          </Grid>
        </TabPanel>
        
        {/* Engagement Tab */}
        <TabPanel value={tabValue} index={2}>
          <Grid container spacing={3}>
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>
                Comments and Reactions Over Time
              </Typography>
              <Box sx={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart
                    data={SAMPLE_ENGAGEMENT_DATA}
                    margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="time" />
                    <YAxis />
                    <RechartsTooltip />
                    <Legend />
                    <Line type="monotone" dataKey="comments" stroke="#8884d8" activeDot={{ r: 8 }} name="Comments" />
                    <Line type="monotone" dataKey="reactions" stroke="#82ca9d" name="Reactions" />
                  </LineChart>
                </ResponsiveContainer>
              </Box>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Typography variant="h6" gutterBottom>
                Engagement Rate
              </Typography>
              <Box sx={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart
                    data={SAMPLE_VIEWER_DATA.map((item, index) => ({
                      time: item.time,
                      rate: index > 0 ? 
                        ((SAMPLE_ENGAGEMENT_DATA[index].comments + SAMPLE_ENGAGEMENT_DATA[index].reactions) / item.viewers * 100).toFixed(1) : 
                        0
                    }))}
                    margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="time" />
                    <YAxis unit="%" />
                    <RechartsTooltip />
                    <Line type="monotone" dataKey="rate" stroke="#ff7300" name="Engagement Rate (%)" />
                  </LineChart>
                </ResponsiveContainer>
              </Box>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Typography variant="h6" gutterBottom>
                Engagement Summary
              </Typography>
              <TableContainer>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Metric</TableCell>
                      <TableCell align="right">Value</TableCell>
                      <TableCell align="right">Per Viewer</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    <TableRow>
                      <TableCell>Total Comments</TableCell>
                      <TableCell align="right">{totalComments}</TableCell>
                      <TableCell align="right">{(totalComments / peakViewers).toFixed(2)}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>Total Reactions</TableCell>
                      <TableCell align="right">{totalReactions}</TableCell>
                      <TableCell align="right">{(totalReactions / peakViewers).toFixed(2)}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>Total Gifts</TableCell>
                      <TableCell align="right">{totalGifts}</TableCell>
                      <TableCell align="right">{(totalGifts / peakViewers).toFixed(2)}</TableCell>
                    </TableRow>
                    <TableRow>
                      <TableCell>Overall Engagement</TableCell>
                      <TableCell align="right">{totalComments + totalReactions + totalGifts}</TableCell>
                      <TableCell align="right">{((totalComments + totalReactions + totalGifts) / peakViewers).toFixed(2)}</TableCell>
                    </TableRow>
                  </TableBody>
                </Table>
              </TableContainer>
            </Grid>
          </Grid>
        </TabPanel>
        
        {/* Audience Tab */}
        <TabPanel value={tabValue} index={3}>
          <Grid container spacing={3}>
            <Grid item xs={12} md={6}>
              <Typography variant="h6" gutterBottom>
                Viewer Demographics
              </Typography>
              <Box sx={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={[
                        { name: 'Age 13-17', value: 15 },
                        { name: 'Age 18-24', value: 35 },
                        { name: 'Age 25-34', value: 30 },
                        { name: 'Age 35-44', value: 15 },
                        { name: 'Age 45+', value: 5 }
                      ]}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {COLORS.map((color, index) => (
                        <Cell key={`cell-${index}`} fill={color} />
                      ))}
                    </Pie>
                    <RechartsTooltip />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              </Box>
            </Grid>
            
            <Grid item xs={12} md={6}>
              <Typography variant="h6" gutterBottom>
                Viewer Locations
              </Typography>
              <Box sx={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={[
                      { name: 'Lagos', value: 45 },
                      { name: 'Abuja', value: 20 },
                      { name: 'Port Harcourt', value: 15 },
                      { name: 'Kano', value: 10 },
                      { name: 'Ibadan', value: 5 },
                      { name: 'Other', value: 5 }
                    ]}
                    margin={{ top: 10, right: 30, left: 0, bottom: 0 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="name" />
                    <YAxis unit="%" />
                    <RechartsTooltip />
                    <Bar dataKey="value" fill="#8884d8" name="Viewers (%)" />
                  </BarChart>
                </ResponsiveContainer>
              </Box>
            </Grid>
            
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>
                Viewer Acquisition
              </Typography>
              <Box sx={{ height: 300 }}>
                <ResponsiveContainer width="100%" height="100%">
                  <PieChart>
                    <Pie
                      data={[
                        { name: 'Homepage', value: 30 },
                        { name: 'Notifications', value: 25 },
                        { name: 'Search', value: 20 },
                        { name: 'Profile', value: 15 },
                        { name: 'Shared Link', value: 10 }
                      ]}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percent }) => `${name}: ${(percent * 100).toFixed(0)}%`}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {COLORS.map((color, index) => (
                        <Cell key={`cell-${index}`} fill={color} />
                      ))}
                    </Pie>
                    <RechartsTooltip />
                    <Legend />
                  </PieChart>
                </ResponsiveContainer>
              </Box>
            </Grid>
          </Grid>
        </TabPanel>
      </Paper>
      
      <Box display="flex" justifyContent="flex-end" mt={2}>
        <Button
          variant="contained"
          color="primary"
          startIcon={<DownloadIcon />}
          onClick={handleExport}
        >
          Export Full Report
        </Button>
      </Box>
    </Box>
  );
};

export default StreamAnalytics;

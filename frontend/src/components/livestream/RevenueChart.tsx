import React, { useState } from 'react';
import {
  Box,
  Typography,
  ToggleButtonGroup,
  ToggleButton
} from '@mui/material';
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer
} from 'recharts';
import { CreatorRevenue } from '../../api/livestreamService';

interface RevenueChartProps {
  data: CreatorRevenue[];
}

const RevenueChart: React.FC<RevenueChartProps> = ({ data }) => {
  const [chartType, setChartType] = useState('revenue');
  
  const handleChartTypeChange = (_event: React.MouseEvent<HTMLElement>, newType: string | null) => {
    if (newType !== null) {
      setChartType(newType);
    }
  };
  
  // Process data for chart
  const processData = () => {
    if (!data || data.length === 0) return [];
    
    // Sort data by period start date
    const sortedData = [...data].sort((a, b) => 
      new Date(a.periodStart).getTime() - new Date(b.periodStart).getTime()
    );
    
    // Format data for chart
    return sortedData.map(item => ({
      name: new Date(item.periodStart).toLocaleDateString(),
      revenue: item.netRevenue,
      gifts: item.totalGifts,
      coins: item.totalCoins
    }));
  };
  
  const chartData = processData();
  
  // Custom tooltip
  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <Box sx={{ bgcolor: 'background.paper', p: 2, border: '1px solid #ccc', borderRadius: 1 }}>
          <Typography variant="body2" fontWeight="bold" mb={1}>
            {label}
          </Typography>
          {payload.map((entry: any, index: number) => (
            <Typography 
              key={`item-${index}`} 
              variant="body2" 
              color={entry.color}
            >
              {entry.name}: {entry.name === 'revenue' ? '₦' : ''}{entry.value.toLocaleString()}
            </Typography>
          ))}
        </Box>
      );
    }
    
    return null;
  };
  
  return (
    <Box sx={{ height: '100%', width: '100%' }}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={2}>
        <Typography variant="h6">
          {chartType === 'revenue' ? 'Revenue Over Time' : 
           chartType === 'gifts' ? 'Gifts Received Over Time' : 
           'Coins Received Over Time'}
        </Typography>
        
        <ToggleButtonGroup
          value={chartType}
          exclusive
          onChange={handleChartTypeChange}
          size="small"
        >
          <ToggleButton value="revenue">
            Revenue
          </ToggleButton>
          <ToggleButton value="gifts">
            Gifts
          </ToggleButton>
          <ToggleButton value="coins">
            Coins
          </ToggleButton>
        </ToggleButtonGroup>
      </Box>
      
      {chartData.length > 0 ? (
        <ResponsiveContainer width="100%" height={350}>
          <BarChart
            data={chartData}
            margin={{
              top: 5,
              right: 30,
              left: 20,
              bottom: 5,
            }}
          >
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="name" />
            <YAxis />
            <Tooltip content={<CustomTooltip />} />
            <Legend />
            {chartType === 'revenue' && (
              <Bar dataKey="revenue" name="Revenue (₦)" fill="#8884d8" />
            )}
            {chartType === 'gifts' && (
              <Bar dataKey="gifts" name="Gifts" fill="#82ca9d" />
            )}
            {chartType === 'coins' && (
              <Bar dataKey="coins" name="Coins" fill="#ffc658" />
            )}
          </BarChart>
        </ResponsiveContainer>
      ) : (
        <Box display="flex" justifyContent="center" alignItems="center" height={350}>
          <Typography variant="body1" color="text.secondary">
            No data available for chart
          </Typography>
        </Box>
      )}
    </Box>
  );
};

export default RevenueChart;

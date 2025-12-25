import React from 'react';
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip, Legend } from 'recharts';
import { AlertCircle } from 'lucide-react';

const COLORS = ['#ef4444', '#f97316', '#eab308', '#06b6d4', '#8b5cf6', '#ec4899'];

const RejectionChart = ({ data }) => {
  // Transform object map to array for Recharts
  const chartData = Object.entries(data).map(([name, value]) => ({
    name: name.replace(/_/g, ' '),
    value
  }));

  const total = chartData.reduce((sum, item) => sum + item.value, 0);

  // Custom label to show percentage
  const renderCustomLabel = ({ cx, cy, midAngle, innerRadius, outerRadius, percent }) => {
    if (percent < 0.05) return null; // Don't show label if less than 5%
    
    const radius = innerRadius + (outerRadius - innerRadius) * 0.5;
    const x = cx + radius * Math.cos(-midAngle * Math.PI / 180);
    const y = cy + radius * Math.sin(-midAngle * Math.PI / 180);

    return (
      <text 
        x={x} 
        y={y} 
        fill="white" 
        textAnchor={x > cx ? 'start' : 'end'} 
        dominantBaseline="central"
        className="text-xs font-bold"
      >
        {`${(percent * 100).toFixed(0)}%`}
      </text>
    );
  };

  // Custom tooltip
  const CustomTooltip = ({ active, payload }) => {
    if (active && payload && payload.length) {
      const data = payload[0];
      const percentage = ((data.value / total) * 100).toFixed(1);
      
      return (
        <div className="bg-white px-3 py-2 rounded-lg border border-neutral-200 shadow-lg">
          <p className="text-xs font-semibold text-neutral-800 mb-1 capitalize">
            {data.name}
          </p>
          <p className="text-xs text-neutral-600">
            <span className="font-bold text-neutral-900">{data.value}</span> jobs ({percentage}%)
          </p>
        </div>
      );
    }
    return null;
  };

  // Custom legend - compact version
  const renderLegend = (props) => {
    const { payload } = props;
    return (
      <div className="flex flex-wrap gap-3 justify-center text-xs">
        {payload.map((entry, index) => {
          const percentage = ((entry.payload.value / total) * 100).toFixed(0);
          return (
            <div key={`legend-${index}`} className="flex items-center gap-1.5">
              <div 
                className="w-2.5 h-2.5 rounded-sm flex-shrink-0" 
                style={{ backgroundColor: entry.color }}
              />
              <span className="text-neutral-700 font-medium capitalize">
                {entry.value}
              </span>
              <span className="text-neutral-500 font-semibold">
                ({percentage}%)
              </span>
            </div>
          );
        })}
      </div>
    );
  };

  return (
    <div className="bg-white p-5 rounded-lg border border-neutral-200 shadow-sm h-full flex flex-col">
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <AlertCircle className="w-5 h-5 text-red-600" />
          <h3 className="text-lg font-semibold text-neutral-800">Rejection Reasons</h3>
        </div>
        <span className="text-xs text-neutral-600 bg-neutral-100 px-3 py-1 rounded-full font-medium">
          {total} Total
        </span>
      </div>
      
      {chartData.length === 0 ? (
        <div className="flex-1 flex items-center justify-center text-neutral-400">
          <div className="text-center">
            <AlertCircle className="w-12 h-12 mx-auto mb-3 opacity-50" />
            <p className="text-sm font-medium">No rejection data</p>
            <p className="text-xs mt-1">Data will appear as jobs are rejected</p>
          </div>
        </div>
      ) : (
        <div className="flex-1 min-h-0">
          <ResponsiveContainer width="100%" height="100%">
            <PieChart>
              <Pie
                data={chartData}
                cx="50%"
                cy="50%"
                innerRadius="45%"
                outerRadius="70%"
                paddingAngle={2}
                dataKey="value"
                label={renderCustomLabel}
                labelLine={false}
              >
                {chartData.map((entry, index) => (
                  <Cell 
                    key={`cell-${index}`} 
                    fill={COLORS[index % COLORS.length]}
                    stroke="white"
                    strokeWidth={2}
                  />
                ))}
              </Pie>
              <Tooltip content={<CustomTooltip />} />
              <Legend 
                content={renderLegend}
                verticalAlign="bottom"
                height={36}
              />
            </PieChart>
          </ResponsiveContainer>
        </div>
      )}
    </div>
  );
};

export default RejectionChart;
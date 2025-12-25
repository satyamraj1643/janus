import React from 'react';
import { TrendingUp, TrendingDown } from 'lucide-react';

const KPICard = ({ title, value, icon: Icon, trend, trendValue, color }) => {
  return (
    <div className="bg-white p-5 rounded-lg border border-neutral-200 shadow-sm hover:shadow-md transition-shadow">
      <div className="flex items-start justify-between mb-3">
        <div className={`p-2 rounded-md ${color}`}>
          <Icon className="w-5 h-5 text-white" strokeWidth={2.5} />
        </div>
        {trend && (
          <div className={`flex items-center gap-1 px-2 py-1 rounded-md text-xs font-semibold ${
            trend === 'up' 
              ? 'text-green-700 bg-green-50 border border-green-200' 
              : 'text-red-700 bg-red-50 border border-red-200'
          }`}>
            {trend === 'up' ? <TrendingUp className="w-3.5 h-3.5" /> : <TrendingDown className="w-3.5 h-3.5" />}
            {trendValue}
          </div>
        )}
      </div>
      <div className="space-y-1">
        <h3 className="text-2xl font-bold text-neutral-900 tracking-tight">{value}</h3>
        <p className="text-sm text-neutral-600 font-medium">{title}</p>
      </div>
    </div>
  );
};

export default KPICard;
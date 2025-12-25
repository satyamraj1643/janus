import React, { useEffect, useState } from 'react';
import { Layers, CheckCircle, XCircle, Activity, Play } from 'lucide-react';
import { mockStats } from '../data/mockData';
import KPICard from '../components/dashboard/KPICard';
import RejectionChart from '../components/dashboard/RejectionChart';
import ActivityFeed from '../components/dashboard/ActivityFeed';
import { useNavigate } from 'react-router-dom';

const Dashboard = () => {
    const navigate = useNavigate();
    const [stats, setStats] = useState(null);

    // Simulate fetching data
    useEffect(() => {
        // In real app, this would be an API call
        setStats(mockStats);
    }, []);

    if (!stats) {
        return (
            <div className="h-full flex items-center justify-center">
                <div className="text-center">
                    <div className="w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full animate-spin mx-auto mb-3"></div>
                    <p className="text-sm text-neutral-600">Loading dashboard...</p>
                </div>
            </div>
        );
    }

    const acceptanceRate = ((stats.admitted_requests / stats.total_requests) * 100).toFixed(1);

    return (
        <div className="h-full flex flex-col bg-gray-50">
            {/* Header */}
            <div className="flex-shrink-0 bg-white border-b border-neutral-200 px-6 py-4">
                <div className="flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-bold text-neutral-900">Dashboard</h1>
                        <p className="text-sm text-neutral-500 mt-1">Real-time system overview and analytics</p>
                    </div>
                    <button 
                        onClick={() => navigate('/simulate')}
                        className="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-semibold transition-colors shadow-sm"
                    >
                        <Play className="w-4 h-4" />
                        Simulate Job
                    </button>
                </div>
            </div>

            {/* Content Area */}
            <div className="flex-1 overflow-auto">
                <div className="p-6 max-w-[1920px] mx-auto">
                    <div className="flex flex-col gap-6">
                        {/* KPI Cards Row */}
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                            <KPICard 
                                title="Total Jobs" 
                                value={stats.total_requests.toLocaleString()} 
                                icon={Layers} 
                                color="bg-blue-500" 
                            />
                            <KPICard 
                                title="Admitted" 
                                value={stats.admitted_requests.toLocaleString()} 
                                icon={CheckCircle} 
                                color="bg-green-500" 
                                trend="up" 
                                trendValue="+12%" 
                            />
                            <KPICard 
                                title="Rejected" 
                                value={stats.rejected_requests.toLocaleString()} 
                                icon={XCircle} 
                                color="bg-red-500" 
                                trend="up" 
                                trendValue="+2.4%" 
                            />
                            <KPICard 
                                title="Acceptance Rate" 
                                value={`${acceptanceRate}%`} 
                                icon={Activity} 
                                color="bg-purple-500" 
                            />
                        </div>

                        {/* Charts and Activity Grid */}
                        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                            {/* Rejection Chart */}
                            <div className="lg:col-span-1 h-[400px]">
                                <RejectionChart data={stats.rejection_reasons} />
                            </div>

                            {/* Activity Feed */}
                            <div className="lg:col-span-2 h-[400px] lg:h-auto">
                                <ActivityFeed decisions={stats.recent_decisions} />
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
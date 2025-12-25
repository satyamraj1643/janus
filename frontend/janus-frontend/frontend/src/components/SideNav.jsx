import React, { useState, useEffect } from 'react';
import { LayoutDashboard, Zap, Settings as SettingsIcon, PanelLeftClose, PanelLeft } from 'lucide-react';
import { useNavigate, useLocation } from 'react-router-dom';

const SideNav = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [isCollapsed, setIsCollapsed] = useState(window.innerWidth < 1024);

  const navItems = [
    { id: '/', label: 'Dashboard', icon: LayoutDashboard },
    { id: '/simulate', label: 'Simulate Jobs', icon: Zap },
    { id: '/global', label: 'Global Job Settings', icon: SettingsIcon }
  ];

  useEffect(() => {
    const handleResize = () => {
      const width = window.innerWidth;
      
      // Auto-collapse below 1024px, auto-expand at 1024px and above
      if (width < 1024) {
        setIsCollapsed(true);
      } else {
        setIsCollapsed(false);
      }
    };

    window.addEventListener('resize', handleResize);
    
    return () => window.removeEventListener('resize', handleResize);
  }, []);

  return (
    <div 
      className={`h-full bg-neutral-50 border-r border-neutral-200 flex flex-col transition-all duration-300 ${
        isCollapsed ? 'w-16' : 'w-64'
      }`}
    >
      {/* Top Navigation Items */}
      <nav className="flex-1 px-3 py-4">
        {navItems.map((item) => {
          const Icon = item.icon;
          const isActive = location.pathname === item.id;
          
          return (
            <button
              key={item.id}
              onClick={() => navigate(item.id)}
              className={`w-full flex items-center px-3 py-2.5 mb-1 rounded transition-all group relative ${
                isActive
                  ? 'bg-blue-50 text-blue-600'
                  : 'text-neutral-700 hover:bg-neutral-100'
              } ${isCollapsed ? 'justify-center' : 'gap-3'}`}
              title={isCollapsed ? item.label : ''}
            >
              <Icon className="w-5 h-5 flex-shrink-0" strokeWidth={isActive ? 2.5 : 2} />
              {!isCollapsed && (
                <span 
                  className={`text-sm whitespace-nowrap ${
                    isActive ? 'font-semibold' : 'font-normal'
                  }`}
                >
                  {item.label}
                </span>
              )}
              
              {/* Tooltip for collapsed state */}
              {isCollapsed && (
                <div className="absolute left-full ml-2 px-2 py-1 bg-neutral-800 text-white text-xs rounded opacity-0 group-hover:opacity-100 pointer-events-none transition-opacity whitespace-nowrap z-50">
                  {item.label}
                </div>
              )}
            </button>
          );
        })}
      </nav>

      {/* Toggle Button at Bottom */}
      <div className={`px-3 py-3 border-t border-neutral-200 flex ${isCollapsed ? 'justify-center' : 'justify-end'}`}>
        <button
          onClick={() => setIsCollapsed(!isCollapsed)}
          className="p-1.5 rounded hover:bg-neutral-200 text-neutral-600 transition-colors"
          aria-label={isCollapsed ? 'Expand sidebar' : 'Collapse sidebar'}
        >
          {isCollapsed ? (
            <PanelLeft className="w-5 h-5" />
          ) : (
            <PanelLeftClose className="w-5 h-5" />
          )}
        </button>
      </div>
    </div>
  );
};

export default SideNav;
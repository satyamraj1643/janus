import React from 'react'
import { Routes, Route } from 'react-router-dom'
import SideNav from './components/SideNav'
import Dashboard from './pages/Dashboard'
import Global from './pages/Global'
import Simulate from './pages/Simulate'

const App = () => {
  return (
    <div className='flex flex-row h-screen w-full overflow-hidden'>
      
      {/* Sidebar - Responsive width handled by SideNav component */}
      <aside className='h-full flex-shrink-0'>
        <SideNav />
      </aside>

      {/* Main Content Area - Takes remaining space */}
      <main className='flex-1 h-full overflow-y-auto bg-gray-50'>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/global" element={<Global />} />
          <Route path="/simulate" element={<Simulate />} />
          <Route path="*" element={<Dashboard />} />
        </Routes>
      </main>

    </div>
  )
}

export default App
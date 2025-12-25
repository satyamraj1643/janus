import React, { useState } from 'react';
import { CheckCircle, XCircle, Clock, ChevronLeft, ChevronRight, Activity, FileJson, X, Copy, Check } from 'lucide-react';

const ITEMS_PER_PAGE = 7; // Matching the requested pagination size usually, kept consistent

const ActivityFeed = ({ decisions = [] }) => {
  const [currentPage, setCurrentPage] = useState(1);
  const [selectedJob, setSelectedJob] = useState(null); // For payload modal
  const [copied, setCopied] = useState(false);

  const totalPages = Math.ceil(decisions.length / ITEMS_PER_PAGE);
  const startIndex = (currentPage - 1) * ITEMS_PER_PAGE;
  const currentItems = decisions.slice(startIndex, startIndex + ITEMS_PER_PAGE);

  const handleCopy = () => {
    if (!selectedJob) return;
    navigator.clipboard.writeText(JSON.stringify(selectedJob.payload, null, 2));
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <>
      <div className="bg-white rounded-lg border border-neutral-200 shadow-sm h-full flex flex-col relative">
        {/* Header */}
        <div className="p-5 border-b border-neutral-200">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <Activity className="w-5 h-5 text-blue-600" />
              <h3 className="text-lg font-semibold text-neutral-800">Live Activity Feed</h3>
            </div>
            <span className="text-xs text-neutral-600 bg-neutral-100 px-3 py-1 rounded-full font-medium">
              {decisions.length} Jobs
            </span>
          </div>
        </div>
        
        {/* Table */}
        <div className="flex-1 overflow-auto">
          {decisions.length === 0 ? (
            <div className="flex flex-col items-center justify-center h-full text-neutral-400 p-8">
              <Clock className="w-12 h-12 mb-3 opacity-50" />
              <p className="text-sm font-medium">No jobs yet</p>
              <p className="text-xs mt-1">Jobs will appear here as they are processed</p>
            </div>
          ) : (
            <table className="w-full text-left border-collapse">
              <thead className="bg-neutral-50 sticky top-0 z-10">
                <tr>
                  <th className="px-5 py-3 text-xs font-semibold text-neutral-600 uppercase tracking-wider">Status</th>
                  <th className="px-5 py-3 text-xs font-semibold text-neutral-600 uppercase tracking-wider">Job ID</th>
                  <th className="px-5 py-3 text-xs font-semibold text-neutral-600 uppercase tracking-wider">Reason</th>
                  <th className="px-5 py-3 text-xs font-semibold text-neutral-600 uppercase tracking-wider text-right">Payload</th>
                  <th className="px-5 py-3 text-xs font-semibold text-neutral-600 uppercase tracking-wider text-right">Time</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-neutral-100">
                {currentItems.map((decision, idx) => (
                  <tr key={idx} className="hover:bg-neutral-50 transition-colors group">
                    <td className="px-5 py-3.5">
                      {decision.admitted ? (
                        <div className="flex items-center gap-2 text-green-600">
                          <CheckCircle className="w-4 h-4 flex-shrink-0" />
                          <span className="text-xs font-semibold">Admitted</span>
                        </div>
                      ) : (
                        <div className="flex items-center gap-2 text-red-600">
                          <XCircle className="w-4 h-4 flex-shrink-0" />
                          <span className="text-xs font-semibold">Rejected</span>
                        </div>
                      )}
                    </td>
                    <td className="px-5 py-3.5">
                      <span className="text-sm font-mono text-neutral-700 bg-neutral-100 px-2 py-0.5 rounded group-hover:bg-neutral-200 transition-colors">
                        {decision.job_id}
                      </span>
                    </td>
                    <td className="px-5 py-3.5">
                      <span className={`px-2.5 py-1 rounded-md text-xs font-medium inline-block ${
                        decision.admitted 
                          ? 'bg-green-100 text-green-800 border border-green-200' 
                          : 'bg-red-100 text-red-800 border border-red-200'
                      }`}>
                        {decision.reason.replace(/_/g, ' ')}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 text-right">
                       <button 
                         onClick={() => setSelectedJob(decision)}
                         className="inline-flex items-center gap-1.5 px-2 py-1 rounded text-xs font-medium text-blue-600 hover:text-blue-700 hover:bg-blue-50 transition-colors"
                       >
                         <FileJson className="w-3.5 h-3.5" />
                         View
                       </button>
                    </td>
                    <td className="px-5 py-3.5 text-right">
                      <div className="flex items-center justify-end gap-1.5 text-neutral-500">
                        <Clock className="w-3.5 h-3.5" />
                        <span className="text-xs font-medium">
                          {new Date(decision.timestamp).toLocaleTimeString()}
                        </span>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>

        {/* Pagination Controls */}
        {decisions.length > ITEMS_PER_PAGE && (
          <div className="p-4 border-t border-neutral-200 flex items-center justify-between bg-neutral-50 flex-shrink-0">
            <div className="text-xs text-neutral-600 font-medium">
              Showing {startIndex + 1}-{Math.min(startIndex + ITEMS_PER_PAGE, decisions.length)} of {decisions.length}
            </div>
            <div className="flex items-center gap-1">
              <button
                onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
                disabled={currentPage === 1}
                className="p-1.5 rounded-md hover:bg-white border border-transparent hover:border-neutral-300 disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:bg-transparent disabled:hover:border-transparent transition-all"
                aria-label="Previous page"
              >
                <ChevronLeft className="w-4 h-4 text-neutral-600" />
              </button>
              
              <div className="px-3 py-1 text-xs font-semibold text-neutral-700 bg-white border border-neutral-300 rounded-md min-w-[60px] text-center">
                {currentPage} / {totalPages}
              </div>
              
              <button
                onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
                disabled={currentPage === totalPages}
                className="p-1.5 rounded-md hover:bg-white border border-transparent hover:border-neutral-300 disabled:opacity-40 disabled:cursor-not-allowed disabled:hover:bg-transparent disabled:hover:border-transparent transition-all"
                aria-label="Next page"
              >
                <ChevronRight className="w-4 h-4 text-neutral-600" />
              </button>
            </div>
          </div>
        )}
      </div>

      {/* Payload Modal Overlay */}
      {selectedJob && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 backdrop-blur-sm p-4 animate-in fade-in duration-200">
          <div className="bg-white rounded-xl shadow-2xl w-full max-w-2xl overflow-hidden flex flex-col max-h-[90vh] animate-in zoom-in-95 duration-200">
            {/* Modal Header */}
            <div className="px-5 py-4 border-b border-neutral-200 flex items-center justify-between bg-neutral-50">
               <div className="flex items-center gap-3">
                 <div className="p-2 bg-blue-100 rounded-lg">
                    <FileJson className="w-5 h-5 text-blue-600" />
                 </div>
                 <div>
                    <h3 className="text-lg font-bold text-neutral-900">Job Payload</h3>
                    <div className="text-xs text-neutral-500 font-mono mt-0.5">{selectedJob.job_id}</div>
                 </div>
               </div>
               <button 
                 onClick={() => setSelectedJob(null)}
                 className="p-2 text-neutral-400 hover:text-neutral-600 hover:bg-neutral-200 rounded-full transition-colors"
               >
                 <X className="w-5 h-5" />
               </button>
            </div>

            {/* Modal Body */}
            <div className="p-0 flex-1 overflow-hidden flex flex-col bg-neutral-900 relative group">
                <div className="absolute top-3 right-3 z-10 opacity-0 group-hover:opacity-100 transition-opacity">
                    <button 
                        onClick={handleCopy}
                        className="flex items-center gap-1.5 px-3 py-1.5 bg-white/10 hover:bg-white/20 text-white rounded-md text-xs font-medium backdrop-blur-md border border-white/10 transition-colors"
                    >
                        {copied ? <Check className="w-3.5 h-3.5 text-green-400" /> : <Copy className="w-3.5 h-3.5" />}
                        {copied ? "Copied!" : "Copy JSON"}
                    </button>
                </div>
                <div className="overflow-auto p-5 custom-scrollbar dark-scrollbar">
                    <pre className="text-sm font-mono text-green-400 leading-relaxed whitespace-pre-wrap">
                        {JSON.stringify(selectedJob.payload || {}, null, 2)}
                    </pre>
                </div>
            </div>
            
            {/* Footer Status info */}
            <div className="px-5 py-3 bg-neutral-50 border-t border-neutral-200 text-xs text-neutral-500 flex justify-between items-center">
                <span>Timestamp: {new Date(selectedJob.timestamp).toLocaleString()}</span>
                {selectedJob.admitted ? (
                    <span className="flex items-center gap-1.5 text-green-600 font-medium px-2 py-0.5 bg-green-500/10 rounded">
                        <CheckCircle className="w-3.5 h-3.5" /> Admitted
                    </span>
                ) : (
                    <span className="flex items-center gap-1.5 text-red-600 font-medium px-2 py-0.5 bg-red-500/10 rounded">
                        <XCircle className="w-3.5 h-3.5" /> Rejected
                    </span>
                )}
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default ActivityFeed;
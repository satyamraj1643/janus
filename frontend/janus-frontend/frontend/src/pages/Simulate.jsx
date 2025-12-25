import React, { useState, useEffect } from 'react';
import { Play, Code, Layout, Plus, Trash2, CheckCircle, AlertCircle, FileJson, Copy, Check, Layers, RotateCcw } from 'lucide-react';

const generateJobId = () => "job_" + Math.floor(Math.random() * 10000);

const INITIAL_DATA = {
  job_id: generateJobId(),
  tenant_id: "tenant_A",
  priority: 10,
  dependencies: {
    "db_shard_1": 5,
    "redis_cache": 10
  },
  payload: {
    "task_type": "data_processing"
  }
};

const Simulate = () => {
  const [mode, setMode] = useState('visual'); // 'visual' | 'json'
  
  // State
  const [visualData, setVisualData] = useState({
    ...INITIAL_DATA,
    dependencies: Object.entries(INITIAL_DATA.dependencies).map(([key, value]) => ({ key, value })),
    payloadString: JSON.stringify(INITIAL_DATA.payload, null, 2)
  });
  const [jsonInput, setJsonInput] = useState(JSON.stringify(INITIAL_DATA, null, 2));

  // Batch State
  const [batchName, setBatchName] = useState("experiment_run");
  const [batch, setBatch] = useState([]);
  
  const [results, setResults] = useState(null); // Array of results or null
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);

  // Sync Logic
  useEffect(() => {
    if (mode === 'visual') {
      let payloadObj = {};
      try {
        payloadObj = JSON.parse(visualData.payloadString || "{}");
      } catch (e) { }

      const exportData = {
        ...visualData,
        dependencies: visualData.dependencies.reduce((acc, dep) => {
          if (dep.key) acc[dep.key] = dep.value;
          return acc;
        }, {}),
        payload: payloadObj
      };
      delete exportData.payloadString;
    }
  }, [visualData, mode]);

  const handleModeSwitch = (newMode) => {
    if (newMode === 'json') {
      let payloadObj = {};
      try {
        payloadObj = JSON.parse(visualData.payloadString || "{}");
      } catch (e) {
        alert("Invalid Payload JSON. Please fix before switching.");
        return;
      }

      const exportData = {
        ...visualData,
        dependencies: visualData.dependencies.reduce((acc, dep) => {
          if (dep.key) acc[dep.key] = dep.value;
          return acc;
        }, {}),
        payload: payloadObj
      };
      delete exportData.payloadString;
      setJsonInput(JSON.stringify(exportData, null, 2));
    } else {
      try {
        const parsed = JSON.parse(jsonInput);
        setVisualData({
          ...parsed,
          dependencies: parsed.dependencies 
            ? Object.entries(parsed.dependencies).map(([key, value]) => ({ key, value }))
            : [],
          payloadString: JSON.stringify(parsed.payload || {}, null, 2)
        });
      } catch (e) {
        alert("Invalid JSON. Please fix errors before switching to Visual mode.");
        return;
      }
    }
    setMode(newMode);
  };

  const getJobFromEditor = () => {
     if (mode === 'json') {
       try { return JSON.parse(jsonInput); }
       catch(e) { return null; }
     } else {
        let payloadObj = {};
        try { payloadObj = JSON.parse(visualData.payloadString || "{}"); } catch (e) {}
        
        const job = {
          ...visualData,
          dependencies: visualData.dependencies.reduce((acc, dep) => {
            if (dep.key) acc[dep.key] = dep.value;
            return acc;
          }, {}),
          payload: payloadObj
        };
        delete job.payloadString;
        return job;
     }
  };

  const addToBatch = () => {
    const job = getJobFromEditor();
    if (!job) return;

    setBatch([...batch, job]);
    
    // Reset editor for new job
    const newId = generateJobId();
    setVisualData(prev => ({
      ...prev,
      job_id: newId
    }));
    const newJson = { ...job, job_id: newId };
    setJsonInput(JSON.stringify(newJson, null, 2));
  };

  const clearBatch = () => {
    setBatch([]);
    setResults(null);
  };

  const handleSimulate = async () => {
    setLoading(true);
    setResults(null);

    const currentDraft = getJobFromEditor();
    if (!currentDraft && batch.length === 0) {
        setLoading(false);
        return;
    }

    const jobsToRun = [...batch];
    if (currentDraft && batch.length === 0) jobsToRun.push(currentDraft);

    const simulatedResults = [];

    // Simulate sequentially
    for (const job of jobsToRun) {
        // Mock API call
        await new Promise(r => setTimeout(r, 400)); 
        const isSuccess = Math.random() > 0.3;
        
        simulatedResults.push({
            success: isSuccess,
            jobId: job.job_id || 'unknown',
            message: isSuccess ? "Admitted" : "Rejected",
            reason: isSuccess 
                ? "Resources secured successfully" 
                : "Policy/Quota violation detected",
            latencies: [
                { step: "Auth Check", time: "1ms" },
                { step: "Quota Limit", time: "2ms" },
                { step: "Dependency Check", time: "1ms" }
            ]
        });
    }

    setResults(simulatedResults);
    setLoading(false);
  };

  const getLiveJSON = () => {
    const current = getJobFromEditor();
    const fullList = [...batch];
    if (current && batch.length === 0) fullList.push(current);
    
    // Naming Logic
    let currentBatchName = batchName;
    if (batch.length === 0) {
        currentBatchName = `standalone_job_${Date.now()}`;
    }

    return JSON.stringify({
        batch_name: currentBatchName,
        jobs: fullList
    }, null, 2);
  };

  const handleCopy = () => {
    const content = getLiveJSON();
    navigator.clipboard.writeText(content);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const updateVisualField = (field, value) => {
    setVisualData(prev => ({ ...prev, [field]: value }));
  };

  const addDependency = () => {
    setVisualData(prev => ({
      ...prev,
      dependencies: [...prev.dependencies, { key: "", value: 1 }]
    }));
  };

  const updateDependency = (idx, field, value) => {
    const newDeps = [...visualData.dependencies];
    newDeps[idx][field] = value;
    setVisualData(prev => ({ ...prev, dependencies: newDeps }));
  };

  const removeDependency = (idx) => {
    const newDeps = visualData.dependencies.filter((_, i) => i !== idx);
    setVisualData(prev => ({ ...prev, dependencies: newDeps }));
  };

  return (
    <div className="h-full flex flex-col bg-gray-50 relative overflow-hidden">
      {/* Header */}
      <div className="flex-shrink-0 bg-white border-b border-neutral-200 px-6 py-4">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-xl font-bold text-neutral-900">Job Simulator</h1>
            <p className="text-sm text-neutral-500 mt-1">Test admission policies with custom job payloads</p>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="flex-1 overflow-auto p-6">
        <div className="max-w-5xl mx-auto grid grid-cols-1 lg:grid-cols-3 gap-6">
          
          {/* Main Editor */}
          <div className="lg:col-span-2 flex flex-col gap-4">
            {/* Mode Toggle Tabs */}
            <div className="bg-white rounded-lg border border-neutral-200 p-1 flex shadow-sm w-fit">
              <button
                onClick={() => handleModeSwitch('visual')}
                className={`flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium transition-all ${
                  mode === 'visual' 
                    ? 'bg-neutral-100 text-neutral-900 shadow-sm' 
                    : 'text-neutral-500 hover:text-neutral-900 hover:bg-neutral-50'
                }`}
              >
                <Layout className="w-4 h-4" />
                Visual Editor
              </button>
              <button
                onClick={() => handleModeSwitch('json')}
                className={`flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium transition-all ${
                  mode === 'json' 
                    ? 'bg-neutral-100 text-neutral-900 shadow-sm' 
                    : 'text-neutral-500 hover:text-neutral-900 hover:bg-neutral-50'
                }`}
              >
                <Code className="w-4 h-4" />
                JSON Payload
              </button>
            </div>

            <div className="bg-white rounded-xl border border-neutral-200 shadow-sm overflow-hidden flex flex-col">
              <div className="p-6 flex-1">
                {mode === 'visual' ? (
                  <div className="space-y-6">
                    <div className="grid grid-cols-2 gap-6">
                      <div>
                        <label className="block text-xs font-semibold text-neutral-500 uppercase tracking-wider mb-2">Job Identity</label>
                        <input
                          type="text"
                          value={visualData.job_id}
                          onChange={(e) => updateVisualField('job_id', e.target.value)}
                          className="w-full px-3 py-2 bg-white border border-neutral-300 rounded-md text-neutral-900 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all text-sm font-mono"
                          placeholder="e.g. job_123"
                        />
                      </div>
                      <div>
                        <label className="block text-xs font-semibold text-neutral-500 uppercase tracking-wider mb-2">Tenant Scope</label>
                        <input
                          type="text"
                          value={visualData.tenant_id}
                          onChange={(e) => updateVisualField('tenant_id', e.target.value)}
                          className="w-full px-3 py-2 bg-white border border-neutral-300 rounded-md text-neutral-900 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all text-sm font-mono"
                          placeholder="e.g. tenant_A"
                        />
                      </div>
                    </div>

                    <div>
                      <div className="flex justify-between items-center mb-2">
                         <label className="block text-xs font-semibold text-neutral-500 uppercase tracking-wider">Priority</label>
                      </div>
                      <div className="flex items-center gap-4">
                        <input
                            type="range"
                            min="1"
                            max="100"
                            value={visualData.priority}
                            onChange={(e) => updateVisualField('priority', parseInt(e.target.value))}
                            className="flex-1 h-2 bg-neutral-200 rounded-lg appearance-none cursor-pointer accent-blue-600"
                        />
                        <input
                            type="number"
                            min="1"
                            max="100"
                            value={visualData.priority}
                            onChange={(e) => updateVisualField('priority', parseInt(e.target.value))}
                            className="w-16 px-2 py-1 bg-white border border-neutral-300 rounded text-center text-sm font-bold text-neutral-700 focus:ring-1 focus:ring-blue-500 outline-none"
                        />
                      </div>
                    </div>

                    <div>
                      <div className="flex items-center justify-between mb-3">
                        <label className="block text-xs font-semibold text-neutral-500 uppercase tracking-wider">Dependencies</label>
                        <button 
                          onClick={addDependency}
                          className="text-xs flex items-center gap-1 text-blue-600 hover:text-blue-700 font-medium px-2 py-1 transition-colors"
                        >
                          <Plus className="w-3 h-3" /> Add
                        </button>
                      </div>
                      <div className="space-y-2 bg-neutral-50 p-4 rounded-lg border border-neutral-100 min-h-[100px]">
                        {visualData.dependencies.length === 0 && (
                          <div className="text-sm text-neutral-400 text-center py-2">
                            No resource dependencies defined
                          </div>
                        )}
                        {visualData.dependencies.map((dep, idx) => (
                          <div key={idx} className="flex gap-3 items-center group">
                            <div className="flex-1">
                              <input
                                placeholder="Key (e.g. db_shard)"
                                value={dep.key}
                                onChange={(e) => updateDependency(idx, 'key', e.target.value)}
                                className="w-full px-3 py-1.5 text-sm border border-neutral-300 rounded focus:ring-1 focus:ring-blue-500 bg-white outline-none"
                              />
                            </div>
                            <div className="w-24">
                              <input
                                type="number"
                                placeholder="Cost"
                                value={dep.value}
                                onChange={(e) => updateDependency(idx, 'value', parseInt(e.target.value))}
                                className="w-full px-3 py-1.5 text-sm border border-neutral-300 rounded focus:ring-1 focus:ring-blue-500 bg-white outline-none"
                              />
                            </div>
                            <button 
                              onClick={() => removeDependency(idx)}
                              className="p-1.5 text-neutral-400 hover:text-red-600 transition-colors"
                            >
                              <Trash2 className="w-4 h-4" />
                            </button>
                          </div>
                        ))}
                      </div>
                    </div>

                    <div>
                      <label className="block text-xs font-semibold text-neutral-500 uppercase tracking-wider mb-2">Optional Payload (JSON)</label>
                      <textarea
                        value={visualData.payloadString}
                        onChange={(e) => updateVisualField('payloadString', e.target.value)}
                        className="w-full h-24 px-3 py-2 bg-white border border-neutral-300 rounded-md text-neutral-900 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all text-sm font-mono resize-none"
                        placeholder="{}"
                        spellCheck="false"
                      />
                    </div>
                  </div>
                ) : (
                  <div className="h-[400px]">
                    <textarea
                      value={jsonInput}
                      onChange={(e) => setJsonInput(e.target.value)}
                      className="w-full h-full font-mono text-sm p-4 bg-neutral-50 text-neutral-800 rounded-lg resize-none focus:outline-none border border-neutral-200"
                      spellCheck="false"
                    />
                  </div>
                )}
              </div>

              {/* Editor Footer Actions */}
              <div className="px-6 py-4 bg-neutral-50 border-t border-neutral-200 flex items-center justify-between round-b-xl">
                 <div className="flex items-center gap-2">
                    <button
                        onClick={addToBatch}
                        className="flex items-center gap-2 px-4 py-2 rounded-lg font-medium text-neutral-600 bg-white border border-neutral-200 hover:bg-neutral-100 hover:border-neutral-300 transition-all text-sm shadow-sm"
                    >
                        <Layers className="w-4 h-4" />
                        Add to Batch
                    </button>
                    {batch.length > 0 && (
                        <button
                            onClick={clearBatch}
                            className="flex items-center gap-2 px-4 py-2 rounded-lg font-medium text-red-600 hover:bg-red-50 transition-colors text-sm"
                        >
                            <Trash2 className="w-4 h-4" />
                            Clear ({batch.length})
                        </button>
                    )}
                </div>
                
                <button
                onClick={handleSimulate}
                disabled={loading}
                className={`flex items-center gap-2 px-6 py-2.5 rounded-lg font-semibold text-white shadow-sm transition-all ${
                    loading 
                    ? 'bg-neutral-400 cursor-not-allowed' 
                    : 'bg-blue-600 hover:bg-blue-700 shadow-blue-200'
                }`}
                >
                {loading ? (
                    <>
                    <div className="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin" />
                    Processing...
                    </>
                ) : (
                    <>
                    <Play className="w-4 h-4 fill-current" />
                    {batch.length > 0 ? `Run Batch (${batch.length})` : 'Run Job'}
                    </>
                )}
                </button>
              </div>
            </div>
          </div>

          {/* Right Column: Result & Payload Preview */}
          <div className="lg:col-span-1 flex flex-col gap-6">
            
            {/* 0. Batch Name Input */}
             {batch.length > 0 && (
                <div>
                <label className="block text-xs font-semibold text-neutral-500 uppercase tracking-wider mb-2">Batch/Experiment Name</label>
                <input
                    type="text"
                    value={batchName}
                    onChange={(e) => setBatchName(e.target.value)}
                    className="w-full px-3 py-2 bg-white border border-neutral-300 rounded-md text-neutral-900 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all text-sm"
                    placeholder="e.g. Experiment 1"
                />
                </div>
             )}

            {/* 1. Simulation Result */}
            <div>
              <h3 className="text-xs font-semibold text-neutral-500 uppercase tracking-wider mb-4">Simulation Result</h3>
              {results ? (
                <div className="space-y-3 max-h-[400px] overflow-auto custom-scrollbar pr-1">
                  {results.length > 1 && (
                      <div className="text-xs font-bold text-neutral-500 mb-2">
                          Processing {results.length} Jobs
                      </div>
                  )}
                  {results.map((res, idx) => (
                    <div key={idx} className={`bg-white rounded-xl border shadow-sm ${
                        res.success ? 'border-green-200' : 'border-red-200'
                    }`}>
                        <div className={`px-4 py-2 border-b flex items-center justify-between ${
                        res.success ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'
                        }`}>
                        <div className="flex items-center gap-2">
                            {res.success ? (
                            <CheckCircle className="w-4 h-4 text-green-600" />
                            ) : (
                            <AlertCircle className="w-4 h-4 text-red-600" />
                            )}
                            <span className={`font-bold text-xs ${
                            res.success ? 'text-green-800' : 'text-red-800'
                            }`}>
                            {res.jobId} - {res.message}
                            </span>
                        </div>
                        </div>
                        
                        <div className="p-3 space-y-2">
                           <div className="text-xs text-neutral-600">{res.reason}</div>
                        </div>
                    </div>
                  ))}
                </div>
              ) : (
                <div className="border-2 border-dashed border-neutral-200 rounded-xl p-8 text-center bg-neutral-50/50">
                  <div className="text-sm font-medium text-neutral-400">No Simulation Run</div>
                  <p className="text-xs text-neutral-400 mt-1">
                    Results will appear here after you run a simulation.
                  </p>
                </div>
              )}
            </div>

            {/* 2. Live Payload Preview */}
            <div className="flex flex-col">
              <h3 className="text-xs font-semibold text-neutral-500 uppercase tracking-wider mb-4">
                  Live Payload {batch.length > 0 && `(Batch: ${batch.length})`}
              </h3>
              <div className="bg-neutral-900 rounded-xl shadow-inner border border-neutral-800 overflow-hidden flex flex-col h-[350px]">
                <div className="flex items-center justify-between px-4 py-2 bg-neutral-800 border-b border-neutral-700 flex-shrink-0">
                  <div className="flex items-center gap-2">
                     <FileJson className="w-3 h-3 text-neutral-500" />
                     <span className="text-[10px] font-mono text-neutral-400 uppercase tracking-wider">
                        {batch.length > 0 ? "batch_export.json" : "single_job.json"}
                     </span>
                  </div>
                  <button 
                    onClick={handleCopy}
                    className="flex items-center gap-1.5 text-[10px] font-medium text-neutral-400 hover:text-white transition-colors bg-neutral-700/50 hover:bg-neutral-700 px-2 py-1 rounded"
                  >
                    {copied ? <Check className="w-3 h-3 text-green-400" /> : <Copy className="w-3 h-3" />}
                    {copied ? 'Copied' : 'Copy'}
                  </button>
                </div>
                <div className="p-4 overflow-auto custom-scrollbar dark-scrollbar flex-1">
                  <pre className="text-xs font-mono text-green-400 leading-relaxed whitespace-pre-wrap">
                    {getLiveJSON()}
                  </pre>
                </div>
              </div>
            </div>

          </div>

        </div>
      </div>
    </div>
  );
};

export default Simulate;
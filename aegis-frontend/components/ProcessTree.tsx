"use client";

import { Activity, ShieldCheck, AlertTriangle, ShieldAlert } from "lucide-react";

// pindahin data dummy-nya ke sini biar page.tsx lebih bersih
const dummyProcessData = [
  { id: 1, name: "explorer.exe", pid: 3424, status: "safe", parentId: null },
  { id: 2, name: "cmd.exe", pid: 5122, status: "suspicious", parentId: 1 },
  { id: 3, name: "powershell.exe", pid: 8990, status: "malicious", parentId: 2 },
  { id: 4, name: "svchost.exe", pid: 1124, status: "safe", parentId: null },
];

export default function ProcessTree({ onTargetPid }: { onTargetPid: (cmd: string) => void }) {
  return (
    <div className="mt-6 p-6 bg-card border border-line rounded-xl">
      <h2 className="text-xl font-semibold mb-4 text-text-p border-b border-line pb-2 flex items-center gap-2">
        <Activity className="w-5 h-5 text-blue-400" />
        Active Process Tree
      </h2>
      
      <div className="bg-main p-4 rounded-lg border border-line font-mono text-sm overflow-x-auto">
        {dummyProcessData.map((proc) => (
          <div 
            key={proc.id} 
            className={`flex items-center gap-3 py-2 border-b border-line/50 last:border-0
              ${proc.parentId ? "ml-8" : ""} 
              ${proc.parentId === 2 ? "ml-16" : ""} 
            `}
          >
            {proc.parentId && (
              <div className="w-4 h-px bg-slate-600"></div>
            )}
            
            <div className={`p-1.5 rounded-md ${
              proc.status === 'safe' ? 'bg-green-500/10 text-green-500' :
              proc.status === 'suspicious' ? 'bg-yellow-500/10 text-yellow-500' :
              'bg-red-500/10 text-red-500'
            }`}>
              {proc.status === 'safe' ? <ShieldCheck className="w-4 h-4" /> : 
               proc.status === 'suspicious' ? <AlertTriangle className="w-4 h-4" /> : 
               <ShieldAlert className="w-4 h-4" />}
            </div>

            <div className="flex-1">
              <span className={`font-bold ${
                proc.status === 'safe' ? 'text-text-p' : 
                proc.status === 'suspicious' ? 'text-yellow-500' : 
                'text-red-500'
              }`}>
                {proc.name}
              </span>
              <span className="text-text-s ml-3 text-xs">PID: {proc.pid}</span>
            </div>

            {proc.status !== 'safe' && (
              <button 
                onClick={() => onTargetPid(`kill ${proc.pid}`)}
                className="px-3 py-1 bg-main/50 border border-line hover:bg-main text-text-p rounded text-xs transition-colors"
              >
                Target PID
              </button>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
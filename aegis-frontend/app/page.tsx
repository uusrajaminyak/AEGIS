"use client";

import ProcessTree from "../components/ProcessTree";
import TerminalEmulator from "../components/TerminalEmulator";


import {
  ShieldAlert,
  Activity,
  Server,
  ShieldCheck,
  Terminal,
  AlertTriangle,
  Moon,
  Sun,
  Cloud,
} from "lucide-react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
  LabelList,
} from "recharts";
import axios from "axios";
import { useState, useEffect, useRef } from "react";

export default function Home() {
  const [alerts, setAlerts] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [isTerminalOpen, setIsTerminalOpen] = useState(false);

  // STATE BARU: Buat ngoper perintah "kill PID" dari ProcessTree ke Terminal
  const [targetCommand, setTargetCommand] = useState("");

  const [theme, setTheme] = useState("blue");

  // Efek buat nempel tema ke tag HTML
  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
  }, [theme]);

  useEffect(() => {
    const fetchAlerts = async () => {
      try {
        const response = await axios.get("http://localhost:8888/api/alerts");
        if (response.data && response.data.data) {
          setAlerts(response.data.data);
        }
      } catch (error) {
        console.error("Error fetching alerts:", error);
      } finally {
        setLoading(false);
      }
    };

    fetchAlerts();

    const interval = setInterval(fetchAlerts, 5000);
    return () => clearInterval(interval);
  }, []);

  const getChartData = () => {
    const count = { low: 0, medium: 0, high: 0 };
    alerts.forEach((a) => {
      const sev = parseInt(a.severity);
      if (sev >= 5) count.high++;
      else if (sev >= 3) count.medium++;
      else count.low++;
    });
    return [
      { name: "Low", count: count.low, fill: "#3b82f6" },
      { name: "Medium", count: count.medium, fill: "#f59e0b" },
      { name: "High", count: count.high, fill: "#ef4444" },
    ];
  };

  const chartData = getChartData();

  return (
    <div className="min-h-screen p-8 font-sans">
      <header className="flex items-center justify-between mb-10 border-b border-line pb-6">
        <div className="flex items-center gap-4">
          <ShieldAlert className="w-10 h-10 text-blue-500" />
          <div>
            <h1 className="text-3xl font-bold tracking-wider text-text-primary">
              AEGIS
            </h1>
            <p className="text-text-secondary">
              Cloud-Native Endpoint Detection & Response
            </p>
          </div>
        </div>
        <div className="flex items-center gap-4">

          {/* TOMBOL THEME SWITCHER */}
          <div className="flex bg-card border border-slate-700 rounded-lg p-1 gap-1">
            <button onClick={() => setTheme("blue")} className={`p-1.5 rounded-md transition-all ${theme === 'blue' ? 'bg-blue-500 text-text-primary' : 'text-text-secondary hover:text-text-primary'}`}>
              <Cloud className="w-4 h-4" />
            </button>
            <button onClick={() => setTheme("dark")} className={`p-1.5 rounded-md transition-all ${theme === 'dark' ? 'bg-gray-700 text-text-primary' : 'text-text-secondary hover:text-text-primary'}`}>
              <Moon className="w-4 h-4" />
            </button>
            <button onClick={() => setTheme("light")} className={`p-1.5 rounded-md transition-all ${theme === 'light' ? 'bg-slate-200 text-slate-900' : 'text-text-secondary hover:text-text-primary'}`}>
              <Sun className="w-4 h-4" />
            </button>
          </div>

          {/* TOMBOL ISOLATE HOST (BARU) */}
          <button 
            onClick={() => alert("Mengirim sinyal darurat ke agen... Host diisolasi!")}
            className="flex items-center gap-2 px-4 py-2 bg-red-950/40 hover:bg-red-900/80 border border-red-800 rounded-lg text-sm font-bold text-red-500 hover:text-red-400 transition-all shadow-[0_0_15px_rgba(239,68,68,0.2)]"
          >
            <AlertTriangle className="w-4 h-4" />
            ISOLATE HOST
          </button>

          <span className="flex items-center gap-2 px-4 py-2 bg-card border border-slate-700 rounded-lg text-sm font-medium">
            <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
            System Online
          </span>
        </div>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-10">
        <div className="p-6 bg-card border border-line rounded-xl flex items-center gap-4">
          <div className="p-3 bg-blue-500/10 rounded-lg">
            <Server className="w-6 h-6 text-blue-500" />
          </div>
          <div>
            <p className="text-text-secondary text-sm">Active Agents</p>
            <p className="text-2xl font-bold">1</p>
          </div>
        </div>

        <div className="p-6 bg-card border border-line rounded-xl flex items-center gap-4">
          <div className="p-3 bg-red-500/10 rounded-lg">
            <Activity className="w-6 h-6 text-red-500" />
          </div>
          <div>
            <p className="text-text-secondary text-sm">Total Alerts</p>
            <p className="text-2xl font-bold">{alerts.length}</p>
          </div>
        </div>

        <div className="p-6 bg-card border border-line rounded-xl flex items-center gap-4">
          <div className="p-3 bg-yellow-500/10 rounded-lg">
            <ShieldCheck className="w-6 h-6 text-yellow-500" />
          </div>
          <div>
            <p className="text-text-secondary text-sm">Threats Blocked</p>
            <p className="text-2xl font-bold">0</p>
          </div>
        </div>

        <div
          onClick={() => setIsTerminalOpen(!isTerminalOpen)}
          className={`p-6 rounded-xl flex items-center gap-4 cursor-pointer transition-all duration-300 ${
            isTerminalOpen
              ? "bg-main border border-green-500/50 shadow-[0_0_20px_rgba(34,197,94,0.2)]"
              : "bg-card border border-line hover:bg-main/80 hover:border-green-500/30"
          }`}
        >
          <div className="p-3 bg-green-500/20 rounded-lg">
            <Terminal
              className={`w-6 h-6 text-green-400 ${isTerminalOpen ? "animate-pulse" : ""}`}
            />
          </div>
          <div>
            <p className="text-text-secondary text-sm">Web Terminal</p>
            <p className="text-lg font-bold text-green-400">
              {isTerminalOpen ? "Active" : "Standby"}
            </p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 p-6 bg-card border border-line rounded-xl h-[400px] flex flex-col">
          <h2 className="text-xl font-semibold mb-4 text-text-primary border-b border-line pb-2">
            Threat Severity Distribution
          </h2>
          <div className="flex-1 w-full mt-4">
            {loading ? (
              <div className="h-full flex items-center justify-center text-slate-500 animate-pulse">
                Calibrating Radar...
              </div>
            ) : (
              <ResponsiveContainer
                width="100%"
                height="100%"
                className="font-sans"
              >
                <BarChart data={chartData} style={{ fontFamily: "inherit" }}>
                  <CartesianGrid
                    strokeDasharray="3 3"
                    stroke="var(--color-gridline)"
                    vertical={false}
                    opacity={0.6}
                  />
                  <XAxis
                    dataKey="name"
                    stroke="var(--text-muted)"
                    fontSize={12}
                    tickLine={false}
                    axisLine={false}
                    allowDecimals={false}
                    fontFamily="inherit"
                  />
                  <YAxis
                    stroke="var(--text-muted)"
                    fontSize={12}
                    tickLine={false}
                    axisLine={false}
                    allowDecimals={false}
                    fontFamily="inherit"
                  />
                  <Tooltip
                  cursor={{ fill: 'var(--color-text-s)', opacity: 0.1}} 
                  contentStyle={{
                    backgroundColor: "var(--bg-card)",
                    border: "1px solid var(--border-line)",
                    borderRadius: "8px",
                    padding: "8px 12px",
                    boxShadow: "0 4px 12px rgba(0,0,0,0.1)",
                  }}
                  labelStyle={{ 
                    color: "var(--text-muted)", 
                    fontWeight: "bold", 
                    fontSize: '12px', 
                    marginBottom: '4px' 
                  }}
                  itemStyle={{ 
                    color: "var(--text-main)", 
                    fontSize: '14px', 
                    fontWeight: 'bold', 
                    padding: 0 
                  }}
                  />
                  {chartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.fill} />
                  ))}
                  <Bar dataKey="count" radius={[4, 4, 0, 0]}>
                    <LabelList
                      dataKey="count"
                      position="top"
                      fill="#94a3b8"
                      fontSize={12}
                      fontWeight="bold"
                    />
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            )}
          </div>
        </div>

        <div className="p-6 bg-card border border-line rounded-xl h-[400px] flex flex-col">
          <h2 className="text-xl font-semibold mb-4 text-text-primary border-b border-line pb-2">
            Recent Alerts
          </h2>
          <div className="flex-1 overflow-y-auto pr-2 space-y-3 custom-scrollbar">
            {loading ? (
              <p className="text-center text-text-secondary mt-10 animate-pulse">
                Scanning frequencies...
              </p>
            ) : alerts.length === 0 ? (
              <p className="text-center text-text-secondary mt-10">
                No alerts detected.
              </p>
            ) : (
              alerts.map((alert, index) => (
                <div
                  key={index}
                  className="p-3 bg-main border border-line rounded-lg flex items-start gap-3 shrink-0"
                >
                  <AlertTriangle
                    className={`w-5 h-5 shrink-0 mt-0.5 ${
                      alert.severity === "5"
                        ? "text-red-500"
                        : alert.severity === "4"
                          ? "text-orange-500"
                          : "text-yellow-500"
                    }`}
                  />
                  <div className="overflow-hidden w-full">
                    <div className="flex justify-between items-center mb-1">
                      <span className="text-xs font-bold text-slate-300">
                        {alert.event_type}
                      </span>
                      <span className="text-[10px] text-slate-500">
                        {new Date(alert.created_at).toLocaleTimeString()}
                      </span>
                    </div>
                    <p
                      className="text-xs text-text-secondary truncate"
                      title={alert.description}
                    >
                      {alert.description}
                    </p>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>

      {/* Panggil ProcessTree dan oper fungsinya  */}
      <ProcessTree onTargetPid={(cmd) => { setTargetCommand(cmd); setIsTerminalOpen(true); }} />

      {/* Panggil TerminalEmulator yang udah dirapihin */}
      <TerminalEmulator isOpen={isTerminalOpen} externalCommand={targetCommand} />
    </div>
  );
}
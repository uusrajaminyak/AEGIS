"use client";

import {
  ShieldAlert,
  Activity,
  Server,
  ShieldCheck,
  Terminal,
  AlertTriangle,
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
  const [commandInput, setCommandInput] = useState("");
  const [terminalLogs, setTerminalLogs] = useState<string[]>([
    "AEGIS Tactical Terminal [Version 1.0.0]",
  ]);
  const terminalEndRef = useRef<HTMLDivElement>(null);
  const [isTerminalOpen, setIsTerminalOpen] = useState(false);
  const terminalContainerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (isTerminalOpen && terminalContainerRef.current) {
      setTimeout(() => {
        terminalContainerRef.current?.scrollIntoView({
          behavior: "smooth",
          block: "end",
        });
      }, 150);
    }
  }, [isTerminalOpen]);

  useEffect(() => {
    terminalEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [terminalLogs]);

  const handleCommandSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!commandInput.trim()) return;

    const currentCmd = commandInput;
    setTerminalLogs((prev) => [...prev, `> root@aegis:~# ${currentCmd}`]);
    setCommandInput("");

    try {
      const response = await axios.post("http://localhost:8888/api/command", {
        agent_id: "agent-123",
        command: currentCmd,
        target_process: "cmd.exe",
      });
      setTerminalLogs((prev) => [
        ...prev,
        `[Success] ${response.data.message}`,
      ]);
    } catch (error: any) {
      setTerminalLogs((prev) => [
        ...prev,
        `[Error] Failed to execute command: ${error.message}`,
      ]);
    }
  };

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
      <header className="flex items-center justify-between mb-10 border-b border-slate-800 pb-6">
        <div className="flex items-center gap-4">
          <ShieldAlert className="w-10 h-10 text-blue-500" />
          <div>
            <h1 className="text-3xl font-bold tracking-wider text-white">
              AEGIS
            </h1>
            <p className="text-slate-400">
              Cloud-Native Endpoint Detection & Response
            </p>
          </div>
        </div>
        <div className="flex items-center gap-4">
          <span className="flex items-center gap-2 px-4 py-2 bg-slate-900 border border-slate-700 rounded-lg text-sm font-medium">
            <span className="w-2 h-2 rounded-full bg-green-500 animate-pulse"></span>
            System Online
          </span>
        </div>
      </header>

      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-10">
        <div className="p-6 bg-slate-900 border border-slate-800 rounded-xl flex items-center gap-4">
          <div className="p-3 bg-blue-500/10 rounded-lg">
            <Server className="w-6 h-6 text-blue-500" />
          </div>
          <div>
            <p className="text-slate-400 text-sm">Active Agents</p>
            <p className="text-2xl font-bold">1</p>
          </div>
        </div>

        <div className="p-6 bg-slate-900 border border-slate-800 rounded-xl flex items-center gap-4">
          <div className="p-3 bg-red-500/10 rounded-lg">
            <Activity className="w-6 h-6 text-red-500" />
          </div>
          <div>
            <p className="text-slate-400 text-sm">Total Alerts</p>
            <p className="text-2xl font-bold">{alerts.length}</p>
          </div>
        </div>

        <div className="p-6 bg-slate-900 border border-slate-800 rounded-xl flex items-center gap-4">
          <div className="p-3 bg-yellow-500/10 rounded-lg">
            <ShieldCheck className="w-6 h-6 text-yellow-500" />
          </div>
          <div>
            <p className="text-slate-400 text-sm">Threats Blocked</p>
            <p className="text-2xl font-bold">0</p>
          </div>
        </div>

        <div
          onClick={() => setIsTerminalOpen(!isTerminalOpen)}
          className={`p-6 rounded-xl flex items-center gap-4 cursor-pointer transition-all duration-300 ${
            isTerminalOpen
              ? "bg-slate-800 border border-green-500/50 shadow-[0_0_20px_rgba(34,197,94,0.2)]"
              : "bg-slate-900 border border-slate-800 hover:bg-slate-800 hover:border-green-500/30"
          }`}
        >
          <div className="p-3 bg-green-500/20 rounded-lg">
            <Terminal
              className={`w-6 h-6 text-green-400 ${isTerminalOpen ? "animate-pulse" : ""}`}
            />
          </div>
          <div>
            <p className="text-slate-400 text-sm">Web Terminal</p>
            <p className="text-lg font-bold text-green-400">
              {isTerminalOpen ? "Active" : "Standby"}
            </p>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 p-6 bg-slate-900 border border-slate-800 rounded-xl h-[400px] flex flex-col">
          <h2 className="text-xl font-semibold mb-4 text-slate-200 border-b border-slate-800 pb-2">
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
                    stroke="#1e293b"
                    vertical={false}
                  />
                  <XAxis
                    dataKey="name"
                    stroke="#64748b"
                    fontSize={12}
                    tickLine={false}
                    axisLine={false}
                    allowDecimals={false}
                    fontFamily="inherit"
                  />
                  <YAxis
                    stroke="#64748b"
                    fontSize={12}
                    tickLine={false}
                    axisLine={false}
                    allowDecimals={false}
                    fontFamily="inherit"
                  />
                  <Tooltip
                    cursor={{ fill: "#1e293b" }}
                    contentStyle={{
                      backgroundColor: "#0f172a",
                      border: "1px solid #1e293b",
                      borderRadius: "8px",
                      color: "#f8fafc",
                    }}
                    itemStyle={{ color: "#cbd5e1" }}
                    labelStyle={{ color: "#94a3b8", fontWeight: "bold" }}
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

        <div className="p-6 bg-slate-900 border border-slate-800 rounded-xl h-[400px] flex flex-col">
          <h2 className="text-xl font-semibold mb-4 text-slate-200 border-b border-slate-800 pb-2">
            Recent Alerts
          </h2>
          <div className="flex-1 overflow-y-auto pr-2 space-y-3 custom-scrollbar">
            {loading ? (
              <p className="text-center text-slate-500 mt-10 animate-pulse">
                Scanning frequencies...
              </p>
            ) : alerts.length === 0 ? (
              <p className="text-center text-slate-500 mt-10">
                No alerts detected.
              </p>
            ) : (
              alerts.map((alert, index) => (
                <div
                  key={index}
                  className="p-3 bg-slate-950 border border-slate-800 rounded-lg flex items-start gap-3 shrink-0"
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
                      className="text-xs text-slate-400 truncate"
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

      <div
        ref={terminalContainerRef}
        className={`transition-all duration-500 ease-in-out overflow-hidden ${
          isTerminalOpen
            ? "max-h-[400px] opacity-100 mt-6"
            : "max-h-0 opacity-0 mt-0"
        }`}
      >
        <div className="p-6 bg-[#0a0a0a] border border-slate-800 rounded-xl h-[250px] flex flex-col font-mono shadow-inner relative">
          <div className="absolute inset-0 pointer-events-none bg-[linear-gradient(transparent_50%,rgba(0,0,0,0.25)_50%)] bg-[length:100%_4px] z-10 opacity-20"></div>

          <div className="flex items-center gap-2 mb-2 border-b border-slate-800 pb-2 shrink-0 z-20">
            <Terminal className="w-4 h-4 text-green-500" />
            <h2 className="text-sm font-semibold text-slate-400">
              root@aegis:~#
            </h2>
          </div>

          <div className="flex-1 overflow-y-auto text-sm text-green-400 space-y-1 mb-2 custom-scrollbar z-20">
            {terminalLogs.map((log, index) => (
              <div
                key={index}
                className={`
                ${log.includes("[Error]") ? "text-red-400" : ""} 
                ${log.includes("[Success]") ? "text-blue-400" : ""}
              `}
              >
                {log}
              </div>
            ))}
            <div ref={terminalEndRef} />
          </div>

          <form
            onSubmit={handleCommandSubmit}
            className="flex gap-2 z-20 shrink-0 mt-2"
          >
            <span className="text-green-500 font-bold">{">"}</span>
            <input
              type="text"
              value={commandInput}
              onChange={(e) => setCommandInput(e.target.value)}
              className="flex-1 bg-transparent outline-none text-green-400 placeholder-green-800 font-mono"
              placeholder="Type Command (e.g: kill 4042)"
              autoComplete="off"
              spellCheck="false"
            />
          </form>
        </div>
      </div>
    </div>
  );
}
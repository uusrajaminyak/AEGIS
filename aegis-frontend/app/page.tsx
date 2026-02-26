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
import { useState, useEffect } from "react";

export default function Home() {
  const [alerts, setAlerts] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

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

        <div className="p-6 bg-slate-900 border border-slate-800 rounded-xl flex items-center gap-4 cursor-pointer hover:bg-slate-800 transition-colors">
          <div className="p-3 bg-green-500/10 rounded-lg">
            <Terminal className="w-6 h-6 text-green-500" />
          </div>
          <div>
            <p className="text-slate-400 text-sm">Web Terminal</p>
            <p className="text-lg font-bold text-green-400">Offline</p>
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
          <div className="flex-1 overflow-y-auto pr-2 space-y-3">
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
                  className="p-3 bg-slate-950 border border-slate-800 rounded-lg flex items-start gap-3"
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
    </div>
  );
}

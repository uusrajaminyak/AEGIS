"use client";

import { useState, useEffect, useRef } from "react";
import { Terminal } from "lucide-react";
import axios from "axios";

interface TerminalProps {
  isOpen: boolean;
  externalCommand: string;
}

export default function TerminalEmulator({ isOpen, externalCommand }: TerminalProps) {
  const [commandInput, setCommandInput] = useState("");
  const [terminalLogs, setTerminalLogs] = useState<string[]>([
    "AEGIS Tactical Terminal [Version 1.0.0]",
  ]);
  const terminalEndRef = useRef<HTMLDivElement>(null);
  const terminalContainerRef = useRef<HTMLDivElement>(null);

  // Auto-scroll saat terminal dibuka
  useEffect(() => {
    if (isOpen && terminalContainerRef.current) {
      setTimeout(() => {
        terminalContainerRef.current?.scrollIntoView({
          behavior: "smooth",
          block: "end",
        });
      }, 150);
    }
  }, [isOpen]);

  // Auto-scroll saat ada log baru
  useEffect(() => {
    terminalEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [terminalLogs]);

  // Menangkap perintah dari Process Tree
  useEffect(() => {
    if (externalCommand) {
      setCommandInput(externalCommand);
    }
  }, [externalCommand]);

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
      setTerminalLogs((prev) => [...prev, `[+] ${response.data.message}`]);
    } catch (error: any) {
      setTerminalLogs((prev) => [...prev, `[-] Failed to execute command: ${error.message}`]);
    }
  };

  return (
    <div
      ref={terminalContainerRef}
      className={`transition-all duration-500 ease-in-out overflow-hidden ${
        isOpen ? "max-h-[400px] opacity-100 mt-6" : "max-h-0 opacity-0 mt-0"
      }`}
    >
      <div className="p-6 bg-[#0a0a0a] border border-line rounded-xl h-[250px] flex flex-col font-mono shadow-inner relative">
        <div className="absolute inset-0 pointer-events-none bg-[linear-gradient(transparent_50%,rgba(0,0,0,0.25)_50%)] bg-[length:100%_4px] z-10 opacity-20"></div>

        <div className="flex items-center gap-2 mb-2 border-b border-line pb-2 shrink-0 z-20">
          <Terminal className="w-4 h-4 text-green-500" />
        </div>

        <div className="flex-1 overflow-y-auto text-sm text-green-400 space-y-1 mb-2 custom-scrollbar z-20">
          {terminalLogs.map((log, index) => (
            <div
              key={index}
              className={`
              ${log.includes("[-] Failed") ? "text-red-400" : ""} 
              ${log.includes("[+]") ? "text-blue-400" : ""}
            `}
            >
              {log}
            </div>
          ))}
          <div ref={terminalEndRef} />
        </div>

        <form onSubmit={handleCommandSubmit} className="flex gap-2 z-20 shrink-0 mt-2">
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
  );
}
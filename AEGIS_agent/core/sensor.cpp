#include <iostream>
#include <windows.h>
#include <vector>

extern "C" {
    __declspec(dllexport) void InitSensor() {
        std::cout << "[*] C++ Sensor Initialized" << std::endl;

        std::vector<int> dummy_pids = {1024, 2048, 4096};
        std::cout << "[*] System scan initialized " << dummy_pids.size() << " PIDs found" << std::endl;

        for(int pid : dummy_pids) {
            std::cout << "[*] Scanning PID: " << pid << std::endl;
        }
        
    }
}
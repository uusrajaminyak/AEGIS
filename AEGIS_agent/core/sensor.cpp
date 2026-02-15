#include <iostream>
#include <windows.h>
#include "MinHook.h"

std::string WideToUTF8(const std::wstring& wstr) {
    if (wstr.empty()) return std::string();
    int size_needed = WideCharToMultiByte(CP_UTF8, 0, &wstr[0], (int)wstr.size(), NULL, 0, NULL, NULL);
    std::string strTo(size_needed, 0);
    WideCharToMultiByte(CP_UTF8, 0, &wstr[0], (int)wstr.size(), &strTo[0], size_needed, NULL, NULL);
    return strTo;
}

typedef uintptr_t(*AlertCallback)(uintptr_t severity, uintptr_t messagePtr);
AlertCallback g_AlertCallback = nullptr;

typedef BOOL(WINAPI *CREATEPROCESSW)(LPCWSTR, LPWSTR, LPSECURITY_ATTRIBUTES, LPSECURITY_ATTRIBUTES, BOOL, DWORD, LPVOID, LPCWSTR, LPSTARTUPINFOW, LPPROCESS_INFORMATION);
CREATEPROCESSW fpCreateProcessW = NULL;

BOOL WINAPI DetourCreateProcessW(LPCWSTR lpApplicationName, LPWSTR lpCommandLine, LPSECURITY_ATTRIBUTES lpProcessAttributes, LPSECURITY_ATTRIBUTES lpThreadAttributes, BOOL bInheritHandles, DWORD dwCreationFlags, LPVOID lpEnvironment, LPCWSTR lpCurrentDirectory, LPSTARTUPINFOW lpStartupInfo, LPPROCESS_INFORMATION lpProcessInformation) {
    std::wstring wsAppName(lpApplicationName ? lpApplicationName : L"");
    std::wstring wsCmdLine(lpCommandLine ? lpCommandLine : L"");
    std::wstring fullTarget = wsAppName;
    if (fullTarget.empty()) {
        fullTarget = wsCmdLine;
    } else if (!wsCmdLine.empty()) {
        fullTarget += L" Cmd: " + wsCmdLine;
    }
    if (g_AlertCallback && !fullTarget.empty()) {
        std::string utf8Target = WideToUTF8(fullTarget);
        std::string alertMessage  = "Process Creation: " + utf8Target;
        g_AlertCallback(2, (uintptr_t)alertMessage.c_str());
    }
    return fpCreateProcessW(lpApplicationName, lpCommandLine, lpProcessAttributes, lpThreadAttributes, bInheritHandles, dwCreationFlags, lpEnvironment, lpCurrentDirectory, lpStartupInfo, lpProcessInformation);
}

extern "C" {
    __declspec(dllexport) void SetAlertCallback(AlertCallback callback) {
        g_AlertCallback = callback;
    }

    __declspec(dllexport) void InitSensor() {
        if (MH_Initialize () != MH_OK) {
            std::cerr << "Failed to initialize MinHook." << std::endl;
            return;
        }

        if (MH_CreateHookApi(L"kernel32", "CreateProcessW", (LPVOID)&DetourCreateProcessW, (LPVOID*)&fpCreateProcessW) != MH_OK) {
            std::cerr << "Failed to create hook for CreateProcessW." << std::endl;
            return;
        }

        if (MH_EnableHook(MH_ALL_HOOKS) != MH_OK) {
            std::cerr << "Failed to enable hooks." << std::endl;
            return;
        }

        std::cout << "Sensor initialized and CreateProcessW hooked successfully." << std::endl;
    }
}
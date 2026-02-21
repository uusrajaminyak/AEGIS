#include <winsock2.h>
#include <ws2tcpip.h>
#include <iostream>
#include <string>
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

typedef int (WSAAPI *CONNECT_FUNC)(SOCKET s, const struct sockaddr*, int);
CONNECT_FUNC fpConnect = NULL;

typedef HANDLE(WINAPI *CREATEFILEW)(LPCWSTR, DWORD, DWORD, LPSECURITY_ATTRIBUTES, DWORD, DWORD, HANDLE);
CREATEFILEW fpCreateFileW = NULL;

int WSAAPI DetourConnect(SOCKET s, const struct sockaddr* name, int namelen) {
    if (name->sa_family == AF_INET) {
        struct sockaddr_in* addr = (struct sockaddr_in*)name;
        char ipStr[INET_ADDRSTRLEN];

        inet_ntop(AF_INET, &(addr->sin_addr), ipStr, INET_ADDRSTRLEN);
        int port = ntohs(addr->sin_port);
        if (g_AlertCallback) {
            std::string alertMessage = "Network Connection: " + std::string(ipStr) + ":" + std::to_string(port);
            g_AlertCallback(3, (uintptr_t)alertMessage.c_str());
        }
    }
    return fpConnect(s, name, namelen);
}

HANDLE WINAPI DetourCreateFileW(LPCWSTR lpFileName, DWORD dwDesiredAccess, DWORD dwShareMode, LPSECURITY_ATTRIBUTES lpSecurityAttributes, DWORD dwCreationDisposition, DWORD dwFlagsAndAttributes, HANDLE hTemplateFile) {
    if (dwDesiredAccess & GENERIC_WRITE) {
        if(g_AlertCallback && lpFileName) {
            std::wstring wsFileName(lpFileName);
            if (wsFileName.find(L"C:\\Windows") == std::wstring::npos) {
                std::string utf8FileName = WideToUTF8(wsFileName);
                std::string alertMessage = "File Write: " + utf8FileName;
                g_AlertCallback(4, (uintptr_t)alertMessage.c_str());
            }
        }
    }
    return fpCreateFileW(lpFileName, dwDesiredAccess, dwShareMode, lpSecurityAttributes, dwCreationDisposition, dwFlagsAndAttributes, hTemplateFile);
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

        if (MH_CreateHookApi(L"ws2_32", "connect", (LPVOID)&DetourConnect, (LPVOID*)&fpConnect) != MH_OK) {
            std::cerr << "Failed to create hook for connect." << std::endl;
            return;
        }

        if (MH_CreateHookApi(L"kernel32", "CreateFileW", (LPVOID)&DetourCreateFileW, (LPVOID*)&fpCreateFileW) != MH_OK) {
            std::cerr << "Failed to create hook for CreateFileW." << std::endl;
            return;
        }

        if (MH_EnableHook(MH_ALL_HOOKS) != MH_OK) {
            std::cerr << "Failed to enable hooks." << std::endl;
            return;
        }

        std::cout << "Sensor initialized and all hooks installed successfully." << std::endl;
    }

    __declspec(dllexport) void TestNetworkHook() {
        WSADATA wsaData;
        WSAStartup(MAKEWORD(2, 2), &wsaData);
        SOCKET sock = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP);

        if (sock != INVALID_SOCKET) {
            sockaddr_in clientService;
            clientService.sin_family = AF_INET;
            clientService.sin_addr.s_addr = inet_addr("93.184.215.14");
            clientService.sin_port = htons(80);

            connect(sock, (SOCKADDR*)&clientService, sizeof(clientService));
            closesocket(sock);
        }
        WSACleanup();
    }

    __declspec(dllexport) void TestFileHook() {
        HANDLE hFile = CreateFileW(L"C:\\testfile.txt", GENERIC_WRITE, 0, NULL, CREATE_ALWAYS, FILE_ATTRIBUTE_NORMAL, NULL);
        if (hFile != INVALID_HANDLE_VALUE) {
            CloseHandle(hFile);
        }
    }
}
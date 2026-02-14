#include <iostream>
#include <windows.h>
#include "MinHook.h"

typedef void (__stdcall *ALERT_CALLBACK)(int severity, const char* message);
ALERT_CALLBACK goAlertChannel  = nullptr;

typedef int (WINAPI* MESSAGEBOXW)(HWND, LPCWSTR, LPCWSTR, UINT);
MESSAGEBOXW fpMessageBoxW = NULL;

int WINAPI DetourMessageBoxW(HWND hWnd, LPCWSTR lpText, LPCWSTR lpCaption, UINT uType) {
    std::cout << "Detoured MessageBoxW called!" << std::endl;
    std::wcout << L"Text: " << lpText << std::endl;
    std::wcout << L"Caption: " << lpCaption << std::endl;

    if (goAlertChannel != nullptr) {
        goAlertChannel(1, "Alert from DetourMessageBoxW!");
    }

    std::cout << "Hijacking..." << std::endl;
    return fpMessageBoxW(hWnd, L"This is a test",L"Hacked by AEGIS", uType);
}

extern "C" {
    __declspec(dllexport) void SetAlertCallback(ALERT_CALLBACK callback) {
        goAlertChannel = callback;
        std::cout << "Alert callback set." << std::endl;
    }

    __declspec(dllexport) void InitSensor() {
        std::cout << "Initializing Sensor..." << std::endl;
        if (MH_Initialize() != MH_OK) {
            std::cout << "Failed to initialize MinHook." << std::endl;
            return;
        }

        if (MH_CreateHook((LPVOID)&MessageBoxW, (LPVOID)&DetourMessageBoxW, reinterpret_cast<LPVOID*>(&fpMessageBoxW)) != MH_OK) {
            std::cout << "Failed to create hook for MessageBoxW." << std::endl;
            return;
        }

        if (MH_EnableHook(MH_ALL_HOOKS) != MH_OK) {
            std::cout << "Failed to enable hook for MessageBoxW." << std::endl;
            return;
        }

        std::cout << "Sensor initialized and hook applied." << std::endl;

        MessageBoxW(NULL, L"Original Message", L"Test", MB_OK);
    }
}
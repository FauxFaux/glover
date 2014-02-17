#include <windows.h>
#include <stdio.h>

#pragma comment(lib, "user32.lib")

static HWND msgwnd;

LRESULT CALLBACK keyboardProc(int nCode, WPARAM wParam, LPARAM lParam) {
  if (nCode == HC_ACTION) {
    switch (wParam) {
    case WM_KEYUP:
    case WM_KEYDOWN:
    case WM_SYSKEYDOWN:
    case WM_SYSKEYUP:
      PKBDLLHOOKSTRUCT p = (PKBDLLHOOKSTRUCT)lParam;
      if (p->vkCode == VK_TAB) {
        PostQuitMessage(0);
        return 1;
      }
    }
  }

  return CallNextHookEx(NULL, nCode, wParam, lParam);
}

int main() {
  HINSTANCE hInst = GetModuleHandle(NULL);
  HHOOK hook = SetWindowsHookEx(WH_KEYBOARD_LL, keyboardProc, hInst, 0);
  msgwnd = CreateWindow(TEXT("STATIC"), TEXT("Glover window"), 0, 0, 0, 0, 0,
                             HWND_MESSAGE, 0, hInst, 0);

  MSG msg;
  BOOL ret;
  while (0 != (ret = GetMessage(&msg, msgwnd, 0, 0))) {
    TranslateMessage(&msg);
    DispatchMessage(&msg);
  }

  UnhookWindowsHookEx(hook);
  return msg.wParam;
}

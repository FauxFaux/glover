#include <windows.h>
#include <stdio.h>

#pragma comment(lib, "user32.lib")

LRESULT CALLBACK keyboardProc(int nCode, WPARAM wParam, LPARAM lParam) {
  if (nCode == HC_ACTION) {
    switch (wParam) {
    case WM_KEYUP:
    case WM_KEYDOWN:
    case WM_SYSKEYDOWN:
    case WM_SYSKEYUP:
      PKBDLLHOOKSTRUCT p = (PKBDLLHOOKSTRUCT)lParam;
      if (p->vkCode == VK_TAB) {
        return 1;
      }
    }
  }

  return CallNextHookEx(NULL, nCode, wParam, lParam);
}

int main() {
  HINSTANCE hInst = GetModuleHandle(NULL);
  HHOOK hook = SetWindowsHookEx(WH_KEYBOARD_LL, keyboardProc, hInst, 0);
  MessageBox(NULL, "blocked", "blocker", MB_OK);
  UnhookWindowsHookEx(hook);
  return 0;
}

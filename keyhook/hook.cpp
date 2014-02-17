#include <windows.h>
#include <stdio.h>
#include <limits.h>

#pragma comment(lib, "user32.lib")

// http://c-faq.com/misc/bitsets.html
#define BITMASK(b) (1 << ((b) % CHAR_BIT))
#define BITSLOT(b) ((b) / CHAR_BIT)
#define BITSET(a, b) ((a)[BITSLOT(b)] |= BITMASK(b))
#define BITCLEAR(a, b) ((a)[BITSLOT(b)] &= ~BITMASK(b))
#define BITTEST(a, b) ((a)[BITSLOT(b)] & BITMASK(b))
#define BITNSLOTS(nb) ((nb + CHAR_BIT - 1) / CHAR_BIT)

// vkeys 1-254 inclusive
// http://msdn.microsoft.com/en-us/library/windows/desktop/ms644967.aspx
static const DWORD MAX_KEY = 254;
static char currently_down[BITNSLOTS(MAX_KEY)] = {};
typedef unsigned char uchar;

static HWND msgwnd;

static unsigned short total_set(char key_bitset[]) {
  unsigned short ret = 0;
  for (uchar i = 1; i < MAX_KEY; ++i) {
    ret += BITTEST(key_bitset, i);
  }
  return ret;
}

static void print(char key_bitset[]) {
  for (uchar i = 1; i < MAX_KEY; ++i) {
    if (!BITTEST(key_bitset, i)) {
      continue;
    }
    printf("%d ", i);
  }
  printf("\n");
}

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

      if (p->vkCode < 1 || p->vkCode >= MAX_KEY) {
        break;
      }

      uchar vk = p->vkCode;

      if (WM_KEYUP == wParam || WM_SYSKEYUP == wParam) {
        BITCLEAR(currently_down, vk);
      } else {
        BITSET(currently_down, vk);
      }

      if (total_set(currently_down)) {
        print(currently_down);
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

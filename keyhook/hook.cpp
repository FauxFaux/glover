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
static const size_t array_size = BITNSLOTS(MAX_KEY);
static char currently_down[array_size] = {};
static char now_released[array_size] = {};
static char denied_keys[array_size] = {};
typedef unsigned char uchar;

static HWND msgwnd;

static void print(char key_bitset[]) {
  for (uchar i = 1; i < MAX_KEY; ++i) {
    if (BITTEST(key_bitset, i)) {
      printf("%d ", i);
    }
  }
  printf("\n");
}

static BOOL valid_key(DWORD key) { return key >= 1 && key <= MAX_KEY; }

LRESULT CALLBACK keyboardProc(int nCode, WPARAM wParam, LPARAM lParam) {
  if (nCode == HC_ACTION) {
    switch (wParam) {
    case WM_KEYUP:
    case WM_KEYDOWN:
    case WM_SYSKEYDOWN:
    case WM_SYSKEYUP:

      PKBDLLHOOKSTRUCT p = (PKBDLLHOOKSTRUCT)lParam;

      if (!valid_key(p->vkCode)) {
        break;
      }

      uchar vk = p->vkCode;

      if (WM_KEYUP == wParam || WM_SYSKEYUP == wParam) {
        BITSET(now_released, vk);

        for (size_t i = 0; i < array_size; ++i) {
          now_released[i] &= currently_down[i];
        }

        if (0 == memcmp(currently_down, now_released, array_size)) {
          print(currently_down);
          memset(currently_down, 0, array_size);
          memset(now_released, 0, array_size);
        }
      } else {
        BITSET(currently_down, vk);
      }

      if (BITTEST(denied_keys, vk)) {
        return 1;
      }
    }
  }

  return CallNextHookEx(NULL, nCode, wParam, lParam);
}

int main(int argc, char *argv[]) {
  for (int i = 1; i < argc; ++i) {
    int val = atoi(argv[i]);
    if (!valid_key(val)) {
      fprintf(stderr, "argument %d isn't a valid vkey code: %s\n", i, argv[i]);
      return 2;
    }
    BITSET(denied_keys, val);
  }
  HINSTANCE hInst = GetModuleHandle(NULL);
  HHOOK hook = SetWindowsHookEx(WH_KEYBOARD_LL, keyboardProc, hInst, 0);
  if (NULL == hook) {
    fprintf(stderr, "hooking failed\n");
    return 3;
  }
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

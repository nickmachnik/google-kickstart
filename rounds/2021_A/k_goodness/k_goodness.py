#!/usr/bin/env python

import sys

def main():
    for i, (n, k, s) in enumerate(read_tests()):
        print("case #{}: {}".format(i + 1, abs(k - score(s, n))))
    
def score(s, n):
    res = 0
    for i in range(n // 2):
        if s[i] != s[n - i - 1]:
            res += 1
    return res

def read_tests():
    lines = sys.stdin.readlines()
    n_tests = int(lines[0].strip())
    for i in range(1, n_tests * 2, 2):
        n, k = [int(e) for e in lines[i].strip().split()]
        s = lines[i + 1].strip()
        yield n, k, s
    
if __name__ == "__main__":
    main()
    
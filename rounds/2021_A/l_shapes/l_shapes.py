#!/usr/bin/env python

import sys
from pprint import pprint

def main():
    for i, grid in enumerate(read_tests()):
        print("Case #{}: {}".format(i + 1, count_ls_in_grid(grid)))

def find_segments(grid):
    r, c = len(grid), len(grid[0])
    col_segs = {}
    row_segs = {}
    for row in range(r):
        for col in range(c):
            if grid[row][col] == 1:
                # update vertical segments
                if col not in col_segs:
                    col_segs[col] = [[row, row]]
                elif col_segs[col][-1][1] == row - 1:
                    col_segs[col][-1][1] = row
                elif col_segs[col][-1][0] == col_segs[col][-1][1]:
                    col_segs[col][-1] = [row, row]
                else:
                    col_segs[col].append([row, row])
                # update horizontal segments
                if row not in row_segs:
                    row_segs[row] = [[col, col]]
                elif row_segs[row][-1][1] == col - 1:
                    row_segs[row][-1][1] = col
                elif row_segs[row][-1][0] == row_segs[row][-1][1]:
                    row_segs[row][-1] = [col, col]
                else:
                    row_segs[row].append([col, col])
    
    # clean up too short fragments at ends
    clean_col_segs = {}
    for col, seg in col_segs.items():
        if seg[-1][0] == seg[-1][1]:
            seg.pop()
        if len(seg) > 0:
            clean_col_segs[col] = seg
        
    clean_row_segs = {}
    for row, seg in row_segs.items():
        if seg[-1][0] == seg[-1][1]:
            seg.pop()
        if len(seg) > 0:
            clean_row_segs[row] = seg

    return clean_row_segs, clean_col_segs

def bin_find_seg(pos, segs, lo, hi):
    """Binary search for segment"""
    if hi < lo:
        return None
    mid = (hi + lo) // 2
    if segs[mid][1] < pos:
        return bin_find_seg(pos, segs, mid + 1, hi)
    elif segs[mid][0] > pos:
        return bin_find_seg(pos, segs, lo, mid - 1)
    else:
        return segs[mid]

def count_ls_in_grid(grid) -> int:
    row_segs, col_segs = find_segments(grid)
    l_count = 0
    r, c = len(grid), len(grid[0])
    for row in range(r):
        for col in range(c):
            if grid[row][col] == 0 or row not in row_segs or col not in col_segs:
                continue
            row_seg = bin_find_seg(col, row_segs[row], 0, len(row_segs[row]) - 1)
            if row_seg is None:
                continue
            col_seg = bin_find_seg(row, col_segs[col], 0, len(col_segs[col]) - 1)
            if col_seg is None:
                continue
            
            left_up = overlapping_ls(col - row_seg[0] + 1, row - col_seg[0] + 1)
            up_right = overlapping_ls(row - col_seg[0] + 1, row_seg[1] - col + 1)
            right_down = overlapping_ls(row_seg[1] - col + 1, col_seg[1] - row + 1)
            down_left = overlapping_ls(col_seg[1] - row + 1, col - row_seg[0] + 1)

            l_count += left_up + up_right + right_down + down_left
    return l_count
    
def overlapping_ls(long, short):    
    if short <= 1 or long <= 1:
        return 0

    if long < short:
        long, short = short, long
    
    # -2 because we are always counting the 1, 2 L, which we shouldn't
    return (min(long // 2, short) + short // 2) - 2

def read_tests():
    lines = sys.stdin.readlines()
    n_tests = int(lines[0].strip())
    returned = 0
    curr_line = 1
    
    while returned != n_tests:
        r, c = [int(e) for e in lines[curr_line].strip().split()]
        rows = []
        curr_line += 1
        for i in range(r):
            rows.append([int(e) for e in lines[curr_line + i].strip().split()])
    
        assert len(rows) == r
        assert len(rows[0]) == c
    
        curr_line += i + 1
        returned += 1
        yield rows
    
if __name__ == "__main__":
    main()
    
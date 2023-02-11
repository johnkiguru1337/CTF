#!/usr/bin/python3
# Author : trustie_rity


data = [ "0e","c9","9d","b8","26","83","26","41","74","e9","26","a5","83","94","0e","63","37","37","37","00"]

for j in data:
	cmp = int(j,16)
	print(j)
	for i in range(1000):
		num = ( i * 17 ) % 253
		if num == cmp:
			print(f"j = { j } : num = { num } : i = { i }")
# *i = (char)((long)(*i * 17) % 253);

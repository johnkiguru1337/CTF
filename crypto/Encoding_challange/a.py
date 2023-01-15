#!/usr/bin/env python3
# Author : @trustie_rity

from pwn import * # pip install pwntools
import json , subprocess , base64
from Crypto.Util.number import *

r = remote('socket.cryptohack.org', 13377, level = 'info')

def json_recv():
    line = r.recvline()
    return json.loads(line.decode())

def json_send(hsh):
    request = json.dumps(hsh).encode()
    r.sendline(request)

while True:
	try:
		received = json_recv()
        
		print("Received type: ")
		b = received["type"]
		log.info(f"{b}")
		print("Received encoded value: ")
		c = received["encoded"]
		log.info(f"{c}")
        
		if "utf-8" in b:
			#decoded = "".join([chr(i) for i in c])
			decoded = "".join(map(chr,c))
		elif "hex" in b:
			decoded = bytes.fromhex(c).decode()
		elif "rot13" in b:
			decoded = subprocess.check_output(f"echo {c} | rot13",shell=True).strip().decode()
		elif "bigint" in b:
			#decoded = bytes.fromhex(c.strip("0x")).decode()
			decoded = long_to_bytes(int(c.strip("0x"),16)).decode()
		elif "base64" in b:
			decoded = base64.b64decode(c).decode()
		print(decoded)
		to_send = {
			"decoded": f"{decoded}"
		}
		json_send(to_send)
        
		#print(json_recv())
	except:
		print(received)
		break	

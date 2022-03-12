import zipfile
import pyfiglet
import colorama
from colorama import Fore, Style ,Back

print(Fore.BLUE + "Author : trustie_rity")
ascii_banner = pyfiglet.figlet_format("ZIP Tool")
print(Fore.WHITE + ascii_banner)
message = "rockyou is a sweet wordlist to use :) \n Path : /usr/share/wordlists/rockyou.txt \n Thank me later :) \n"
print(Fore.YELLOW + message)
#print(Style.RESET_ALL)
def crack_password(password_list, obj):
	idx = 0
	with open(password_list, 'rb') as file:
		for line in file:
			for word in line.split():
				try:
					idx += 1
					obj.extractall(pwd=word)
					print("Password found at line", idx)
					print(Fore.GREEN + Back.BLACK + "Password is", word.decode())
					return True
				except:
					continue
	return False


password_list = str(input("Enter the full path of passwordlist: "))

zip_file = str(input("Enter full path of zip file to bruteforce: "))


try:
    obj = zipfile.ZipFile(zip_file)
    cnt = len(list(open(password_list, "rb")))
    print("There are total", cnt,
	"number of passwords to test")
    if crack_password(password_list,obj) == False:
        print("Password not found in this wordlist :( ")
except:
    print(Fore.RED + "[x] File not found...check if you mispelled :(")


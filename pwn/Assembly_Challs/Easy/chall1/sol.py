#! /usr/bin/python3 

encoded = "fhwu`foqg}s`cwgm"

def main():
    print("EncodeStringLength: %s " % len(encoded))

    decoded =  [(i-1) for i in [(ord(ii) ^ 19) for ii in encoded]]
    print("".join(map(chr, decoded)))

if __name__ == "__main__":
    main()

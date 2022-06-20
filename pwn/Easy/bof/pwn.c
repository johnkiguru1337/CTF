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

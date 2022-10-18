#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<unistd.h>
#include<stdbool.h>
#include<sys/ptrace.h>

// author @ trustie_rity
struct case_study{
	int priority;
	char *name;
};
int kid();
void nerd();
char naming[40];
int j = -1;

int main(int argc,char argv[]){
	int input;
	puts("Welcome here,They told you its cool...hmm!they Lied :(\nFirst road will lead you to the first flag ^_^,but \nWhy not take the second one and grab both flags!\nSelect\n> 1 \n or \n> 2\nOh man!I hope you take the second one: ");
	fflush(stdout);
	scanf("%d",&input);
	
	if(input == 2){
		nerd();
	}else if(input == 1){
		kid();
	}else{
		puts("Oh man!Cant follow simple instructions!");
		exit(1);
	}
	return 0;
}

int kid(){
	char input[0x30];
	char password[0xa];
	char reg[0x50];
	char *pass = "S3cur3p4$$";
	int i;
	while( true ){
		while( true ){
			memset(password,0,0x10);
			puts(" _\t_\n|_\t_|\n|_\t_|\n|_\t_|\tOkay there kid!Welcome to this awesome system!!\n\nWhat would you like to do?\n1. Register \n2. Login \n> ");
			fflush(stdout);
			scanf("%d",&i);
			if(i != 1) break;
			printf("Enter your names: ");
			fflush(stdout);
			read(0,reg,0x50);
			printf("Welcome %s\n",reg);
		}
		if(i != 2) break;
		puts("Enter the password: ");
		fflush(stdout);
		scanf("%s",password);
		if(0 == strcmp(password , pass)){
			printf("Logged in,Enter commands to navigate! with passowrd : %s",pass);                		    read(0,input,0x70);
		}else{
			puts("Wrong password !");
			exit(0);
		}		
	}
	puts("Bye :)");
	return 0;
}

void nerd(){
	system("clear");
	struct case_study *case1,*case2,*case3;
	int i;
	int long var2;

	puts(" _\t_\n|_\t_|\n|_\t_|\n|_\t_|\tBright student!Welcome to this case study");
	case1 = malloc(sizeof(struct case_study));
	case1->priority = 1;
	case1->name = malloc(8);

	case2 = malloc(sizeof(struct case_study));
	case2->priority = 2;
	case2->name = malloc(8);

	for(i = 0;i < 2; i++){
		memset(naming,0,0xc);
		printf("Enter the name of your case %d: ",(i+1));
		fflush(stdout);
		scanf("%s",naming);
		var2 = strlen(naming);
		if(i == 0)
			strcpy(case1->name,naming);
		if(i == 1)
			strcpy(case2->name,naming);
		else{
			if(false || j == 0){
				goto weh;
				weh:
				system("/bin/sh");
			}
		}
	}
	i = ptrace(PTRACE_TRACEME,0,1,0);
	if(i == -1){
		puts("Must you do dynamic analysis");
		exit(0);
	}
	else{
		puts("\nAwesome study");
	}
}

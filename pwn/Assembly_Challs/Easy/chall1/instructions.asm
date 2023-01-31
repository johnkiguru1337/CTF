x:

        .ascii  "fhwu`foqg}s`cwgm"

.LC0:

        .string "Partially Might be right!"

.LC1:

        .string "Welcome"

asm_trick:

        push    rbp

        mov     rbp, rsp

        sub     rsp, 32

        mov     QWORD PTR [rbp-24], rdi

        mov     DWORD PTR [rbp-4], 0

        jmp     .L2

.L5:

        mov     eax, DWORD PTR [rbp-4]

        movsx   rdx, eax

        mov     rax, QWORD PTR [rbp-24]

        add     rax, rdx

        movzx   eax, BYTE PTR [rax]

        xor     eax, 19

        movsx   eax, al

        lea     edx, [rax+1]

        mov     eax, DWORD PTR [rbp-4]

        cdqe

        movzx   eax, BYTE PTR x[rax]

        movzx   eax, al

        cmp     edx, eax

        je      .L3

        mov     edi, OFFSET FLAT:.LC0

        call    puts

        mov     eax, 1

        jmp     .L4

.L3:

        add     DWORD PTR [rbp-4], 1

.L2:

        cmp     DWORD PTR [rbp-4], 32

        jle     .L5

        mov     edi, OFFSET FLAT:.LC1

        call    puts

        mov     eax, 1

.L4:

        leave

        ret

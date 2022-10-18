# sudo docker build -t libcorheap .
# sudo docker run -d -p 1025:1025 --rm -it libcorheap

FROM ubuntu:19.10


RUN useradd -d /home/ctf/ -m -p ctf -s /bin/bash ctf
RUN echo "ctf:ctf" | chpasswd

WORKDIR /home/ctf

COPY chall .
COPY flag .
COPY ynetd .

RUN chown -R root:root /home/ctf

USER ctf
EXPOSE 1025
CMD ./ynetd -p 1025 ./chall

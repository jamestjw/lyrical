FROM golang:1.14.1-buster
WORKDIR /app
SHELL ["/bin/bash", "-c"] 

# Install manually all the missing libraries
RUN apt-get update
RUN apt-get install software-properties-common -y
RUN apt-get update

RUN mkdir $HOME/go
RUN echo 'export GOPATH=$HOME/go' >> ~/.bashrc
RUN echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
RUN source ~/.bashrc
RUN apt-get install vim -y
RUN apt-get install ffmpeg -y
RUN apt-get install sqlite3 -y
RUN go get github.com/jonas747/dca/cmd/dca

CMD ["/bin/bash"]

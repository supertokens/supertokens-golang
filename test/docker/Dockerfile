FROM rishabhpoddar/supertokens_website_sdk_testing
RUN cd /tmp
RUN wget https://dl.google.com/go/go1.17.linux-amd64.tar.gz
RUN tar -xvf go1.17.linux-amd64.tar.gz
RUN mv go /usr/local
env GOROOT /usr/local/go
env GOPATH $HOME/go
env PATH $GOPATH/bin:$GOROOT/bin:$PATH
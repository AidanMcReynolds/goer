cd /tmp
wget https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz
tar xzf go1.13.1.linux-amd64.tar.gz
export PATH=$PATH:/tmp/go/bin/go
/tmp/go/bin/go version
cd $HOME
goer

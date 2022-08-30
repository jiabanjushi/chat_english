wget https://studygolang.com/dl/golang/go1.17.5.linux-amd64.tar.gz
tar -C /usr/local -xvf go1.17.5.linux-amd64.tar.gz
mv go1.17.5.linux-amd64.tar.gz /tmp
echo "PATH=\$PATH:/usr/local/go/bin" >> /etc/profile
echo "PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
source /etc/profile
go version
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
FROM golang:1.21-buster

RUN curl -s https://api.github.com/repos/qiniu/goc/releases/latest | grep "browser_download_url.*-linux-amd64.tar.gz" | cut -d : -f 2,3 | tr -d \" | xargs -n 1 curl -L | tar -zx && chmod +x goc && mv goc /usr/local/bin

CMD ["goc", "server"]
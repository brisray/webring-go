# Webring-Go

An attempt to build a [webring] server in [Go]

![Screenshot of example webring](https://github.com/alifeee/webring-go/assets/13833017/997b912d-4705-415d-9f87-b5d468f15b10)

[webring]: https://indieweb.org/webring
[Go]: https://go.dev/

## Commands

Requirements: Go 1.21.4

### Build

```bash
go build ./serve.go
```

### Run

```bash
.\serve.exe # Windows
./serve # Linux
```

## Endpoints

This will run on a server, and any website that wants to participate can add themselves to the webring via a Pull Request, and put a webring at the bottom of their webpage.

### `/webring.js`

this should be included like `<script src="[root]/webring.js"></script>` where you want the webring to go. See [`http://server.alifeee.co.uk:8080/home`](http://server.alifeee.co.uk:8080/home) for an example.

### `/home`

This is the homepage for the webring. For example: [`http://server.alifeee.co.uk:8080/home`](http://server.alifeee.co.uk:8080/home)

### `/next`

Given the header of the requesting site, gets the next site in the ring

### `/previous`

Given the header of the requesting site, gets the previous site in the ring

### Data

There needs to be information about which sites are participating in the webring, and some metadata. This is in the form of a configuration file. For example, see [`arthena.json`].

[`arthena.json`]: https://github.com/mldangelo/open-webring/blob/main/public/ring/arthena.json

- `webring.toml` or `webring.json`

## Deploy on remote server

### Initial deployment

```bash
ssh $USER@$SERVER
cd ~/go
git clone https://github.com/alifeee/webring-go.git
cd webring-go
# install go
wget https://go.dev/dl/go1.21.4.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz
rm go1.21.4.linux-amd64.tar.gz
# edit config
cp webring.example.toml webring.toml
nano webring.toml
# set up tmux
tmux new -s webring
cd ~/go/webring
# build and execute
go build serve.go
./serve
# Ctrl+B, D to detach from tmux
```

### Update deployment

```bash
ssh $USER@$SERVER
tmux ls
tmux attach -t webring
# send ctrl+C
git pull
go build serve.go
./serve
# Ctrl+B, D to detach from tmux
```

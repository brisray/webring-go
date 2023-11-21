# Webring-Go

A [webring] server built with [Go]

[webring]: https://indieweb.org/webring
[Go]: https://go.dev/

To add yourself to the webring, add your site to [`webring.toml`](./webring.toml) and create a Pull Request.

## To include the webring on your site

Firstly, you must be part of the webring, or the next/previous links will fail.

### Simple

Add `<script src="[root]/webring.js"></script>` where you want the webring to appear in the DOM. This will add HTML which will look like the following:

![Screenshot of example webring](images/webring.png)

You can style it via the `.webring` class (for the root element) and the `.previous`, `.name`, `description`, and `next` classes for the child elements.

### More customisable

The script above just adds the following HTML to the DOM, with templates replaced by the items in the [config](./webring.toml):

```html
<section class="webring">
  <a class="previous" href="{{ Root }}/previous">Previous</a>
  <a class="name" href="{{ Root }}/">{{ Name }}</a>
  <p class="description">{{ Description }}</p>
  <a class="next" href="{{ Root }}/next">Next</a>
</section>
```

So long as you include links to:

- homepage `/`
- next `/next`
- previous `/previous`

...you can write the HTML and style it as you want.

## Endpoints

### `/webring.js`

This should be included like `<script src="[root]/webring.js"></script>` where you want the webring to go. See <http://webring.alifeee.co.uk> for an example.

### `/`

This is the homepage for the webring. For example: <http://webring.alifeee.co.uk>

### `/next`

Given the header of the requesting site, returns a redirect to the next site in the ring.

### `/previous`

Given the header of the requesting site, returns a redirect to the previous site in the ring.

## Development

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

### Test

```bash
go test
```

#### Generate and view coverage report (HTML)

```bash
go test -coverprofile="c.out"; go tool cover -html="c.out"
```

### Deploy on remote server

#### Initial deployment

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
nano webring.toml
# set up tmux
tmux new -s webring
cd ~/go/webring
# build and execute
go build serve.go
./serve
# Ctrl+B, D to detach from tmux
```

#### Update deployment

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

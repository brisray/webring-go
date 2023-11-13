# Webring-Go

An attempt to build a [webring] server in [Go]

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

## Theory

This will run on a server, and any website that wants to participate can add themselves to the webring via a Pull Request, and put a webring at the bottom of their webpage.

### HTTP endpoints

- The webring component. This is either an HTML endpoint or a JavaScript endpoint
  - /webring.html - returns HTML for webring, with information and next/previous buttons. It can be fetched into the DOM with JavaScript, or, for example, [`htmx`](https://htmx.org/)
  - /webring.js - this would be included like `<script src="...webring.js"></script>`, and like how [`giscus`](https://giscus.app/) works, the script would replace itself with the HTML in the DOM
- /next - given the header of the requesting site, gets the next site in the ring
- /previous - given the header of the requesting site, gets the previous site in the ring

### Data

There needs to be information about which sites are participating in the webring, and some metadata. This is in the form of a configuration file. For example, see [`arthena.json`].

[`arthena.json`]: https://github.com/mldangelo/open-webring/blob/main/public/ring/arthena.json

- `webring.toml` or `webring.json`

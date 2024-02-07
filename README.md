# Webring-Go

A [webring] server built with [Go]

[webring]: https://indieweb.org/webring
[Go]: https://go.dev/

## What is a webring?

See [`webring.alifeee.co.uk`] for more information and links, but briefly: A webring is a bit of HTML that a group of people put on their websites, which you can click through to navigate between the sites. It is a way of discovering content on the web without search engines, advertising, SEO, and all the things which have made the Internet a... *more corporate* place.

[`webring.alifeee.co.uk`]: https://webring.alifeee.co.uk/

## How to join the webring

1. Add your site to [`webring.toml`](./webring.toml) and create a Pull Request.

    ```toml
    [[Websites]]
    Name = "your website here!"
    Url = "https://your.website"
    Image = "https://your.website/image.png"
    Description = "your description!"
    ```

2. Add the webring to your site!

    Put `<script src="https://webring.alifeee.co.uk/webring.js"></script>` where you want the webring to appear on your site. The webring HTML will be added here!

    ![Screenshot of example webring](images/webring.png)

3. Style the webring how you want!

    You can style it via the `.webring` class (for the root element) and the `.previous`, `.name`, `description`, and `next` classes for the child elements. Inspect the webring HTML for more information, or see the [webring.html template](./templates/webring.html.template).

## More customisation

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

# set up service
cp webring.service /etc/systemd/system/webring.service
sudo systemctl enable webring.service
sudo systemctl start webring.service
sudo systemctl status webring.service
```

#### Update deployment

```bash
ssh $USER@$SERVER
git pull
go build serve.go
cp webring.service /etc/systemd/system/webring.service
sudo systemctl enable webring.service
sudo systemctl start webring.service
sudo systemctl status webring.service
```

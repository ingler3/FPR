# PDX, City of ????

Who doesn't love food? No-one I know.

Here's [the easiest way to find 'em](https://ingler3.github.io/fpr).


## Install

Install dependencies as usual...

```ShellSession
$ go get ./...
```

Then just add API keys!

```ShellSession
# https://developers.google.com/places/web-service/get-api-key
$ export GOOGLE_API_KEY='<your api key here>'

# https://www.mapbox.com/help/how-access-tokens-work/
$ export MAPBOX_ACCESS_TOKEN='<your access token here>'

# Build your own donut map
$ go run main.go \
    -keyword donut \
    -type 'restaurant|bakery' \
    -location 'Portland, OR'
```

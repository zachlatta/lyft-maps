# lyft-maps

## Getting Started

Prerequisites:

* [Docker](https://docker.io)
* [fig](http://orchardup.github.io/fig/)

Set the variables in `config/config.yml`. You'll have to get both your Lyft
Facebook access token and user ID by inspecting requests to their server. Some
good resources on that: [Charles and Android](http://jaanus.com/blog/2012/02/12/debugging-http-on-an-android-phone-or-tablet-with-charles-proxy-for-fun-and-profit/), [Charles and iPhone](http://www.charlesproxy.com/documentation/faqs/using-charles-from-an-iphone/).
Remember to properly setup everything for SSL proxying to work correctly, as
Lyft's API uses HTTPS.

After you have the configuration set up, run `fig up` and you should be good.

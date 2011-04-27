# Lannister

Lannister is a static website generator.

## Features:

* Uses Go's standard library for templating. Yay!
* Multiplexes output. What? It means you can define multiple output templates for a piece of content.
* Multiplexes input. What? Means you can mix in other content together. See the tutorial section.
* Settings are stored in an SQLite database. Don't worry about this, it's just convenient and means you can do more with them.

## Installation.

Sorry. This is Go code, so it's not simple or straightforward to distrubute code. At least, I've not heard of any easy way.

* Install Go, if you haven't already. It involves setting some stupid environment variables unlike any other language, which is lame, but what can you do? -> http://golang.org/doc/install.html
* Then run `goinstall github.com/wjlroe/lannister`
* Or (probably more sensible), clone this repository and type `make`
* You now have a binary, called lannister, that you can copy or link somewhere in your $PATH. I'd recommend ~/bin but that's just me.
* Once again, I am truly sorry about that whole installation thing. One day we will look back and laugh about it. Probably.

## Tutorial

This tutorial is aimed at exactly what I wrote this tool for. I wanted a static site, with content loaded using the awesome PJAX jQuery plugin (https://github.com/defunkt/jquery-pjax). PJAX requires, like any AJAX-based tech, that the content requested is returned without the entire page. This means, if you want to serve up /about.html both as a static file (fully templated) and as a page fragment (for PJAX), you need different layouts/formats/whatever. This, in a nutshell, is the problem Lannister has arrived to solve.

So let's get on with creating a website. It's pretty exciting doing this, every single time and this time is no exception.

 > ./lannister createsite ~/Dropbox/Sites/newsite

See what I did there? Dropbox is cool, you have to admit. Anyway, whether or not you're cool with syncing your new website on every computer you own (seriously, who wouldn't want that?), let's examine what we have got.

In the new directory that was created, we have the following:

 /newsite/
   -
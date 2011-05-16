# Lannister

Lannister is a static website generator. It is a work in progress.

## Features:

* Uses Go's standard library for templating. Yay!
* Multiplexes output. What? It means you can define multiple output templates for a piece of content.
* Multiplexes input. What? Means you can mix in other content together. See the tutorial section.
* Settings are stored in an SQLite database. Don't worry about this, it's just convenient and means you can do more with them.

## Installation.

Sorry. This is Go code, so it's not simple or straightforward to distrubute code. At least, I've not heard of any easy way.

* Install Go, if you haven't already. It involves setting some stupid environment variables unlike any other language, which is lame, but what can you do? -> [http://golang.org/doc/install.html](http://golang.org/doc/install.html)
* Then run `goinstall github.com/wjlroe/lannister`
* Or (probably more sensible), clone this repository and type `make`
* You now have a binary, called lannister, that you can copy or link somewhere in your $PATH. I'd recommend ~/bin but that's just me.
* Once again, I am truly sorry about that whole installation thing. One day we will look back and laugh about it. Probably.

## Tutorial

This tutorial is aimed at exactly what I wrote this tool for. I wanted a static site, with content loaded using the awesome PJAX jQuery plugin [https://github.com/defunkt/jquery-pjax](https://github.com/defunkt/jquery-pjax). PJAX requires, like any AJAX-based tech, that the content requested is returned without the entire page. This means, if you want to serve up /about.html both as a static file (fully templated) and as a page fragment (for PJAX), you need different layouts/formats/whatever. This, in a nutshell, is the problem Lannister has arrived to solve.

So let's get on with creating a website. It's pretty exciting doing this, every single time and this time is no exception.

   ./lannister createsite ~/Dropbox/Sites/newsite

See what I did there? Dropbox is cool, you have to admit. Anyway, whether or not you're cool with syncing your new website on every computer you own (seriously, who wouldn't want that?), let's examine what we have got.

In the new directory that was created, we have the following:

    /newsite/
    |-> layouts/
        |-> default.html
        |-> default-pjax.html
    |-> pages/
        |-> about.md
        |-> index.md
    |-> images/
        |-> loading.gif
    |-> stylesheets/
        |-> style.css
    |-> javascript/
        |-> jquery.min.js
        |-> jquery-ui.min.js
        |-> jquery.pjax.js
        |-> app.js
    |-> site/

OK, that seems reasonable. You will notice from the above that the conventions for using Lannister are for using Markdown for the content.

So now let's generate the site:

   ./lannister generate

That will spit the entire site in the site/ sub-directory.

So what does it do? Well, this bootstrapped site provides a really simple site that shows how you can use PJAX to make cool websites with static content. Put it behind Varnish to cache things and speed up your site even more (saving you some hosting costs when you write that super amazing article - or rant - and get slashdotted).

Now if you want to test the site in your browser, run:

   ./lannister serve

This will serve the content of your site on [localhost:6565](http://localhost:6565) - check it out!

## Alternatives

Yeah, alternatives. Makes it sound like this is the default. Yeah so you *could* use one of those other wacky systems, but you don't want to do that. Still, there *are* other systems out there. No doubt you've heard of at least two of these.

* [Jekyl](https://github.com/mojombo/jekyll). github use this for pages and blogs. I use it for my blog. Yeah. I use it. What an admission.
* [Hyde](http://ringce.com/hyde). As soon as somebody created a project called 'Jekyl', somebody just *had* to create one called 'Hyde'. Full marks there!
* [Bloxsom](http://www.blosxom.com/). Oh yeah! You remember that thing, the one in Perl, right? Yeah. Still exists. Doesn't really do much useful to me. Very blog-focused. Baroque.
* [Hakyll](http://jaspervdj.be/hakyll/index.html) in Haskell
* [utternson](https://github.com/pepijndevos/utterson) in Clojure
* [blatter](http://bitbucket.org/jek/blatter/)
* [lanyon](http://bitbucket.org/arthurk/lanyon/)
* [Webby](http://webby.rubyforge.org/)
* [poole](https://bitbucket.org/obensonne/poole/src)
* [toto](http://cloudhead.io/toto)
* [chisel](https://github.com/dz/chisel)
* [tahchee](http://ivy.fr/tahchee/)
* [middleman](https://github.com/tdreyno/middleman)
* [frank](https://github.com/blahed/frank)
* [bricolage](http://bricolagecms.org/)
* [nanoc](http://nanoc.stoneship.org/)
* [webgen](http://webgen.rubyforge.org/)
* [StaticMatic](http://staticmatic.rubyforge.org/)
* [static](http://static.newqdev.com/)
* [Ensemble](https://launchpad.net/ensemble) I have no idea what this is

## License

<a rel="license" href="http://creativecommons.org/licenses/by-sa/3.0/"><img alt="Creative Commons License" style="border-width:0" src="http://i.creativecommons.org/l/by-sa/3.0/88x31.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by-sa/3.0/">Creative Commons Attribution-ShareAlike 3.0 Unported License</a>.


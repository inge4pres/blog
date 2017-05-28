---
title: My first Golang web project is online
author: inge4pres
layout: post
date: 2015-06-11T12:15:34+00:00
categories:
  - social
  - tech
tags:
  - GO
  - golang
  - shortener
  - software
  - URL

---
It is true: I fell in love with <a href="http://golang.org" target="_blank">Go</a>, not because I love Google and his products, but because it really fits my ideology of simplicity and power in a programming language. I started experimenting with the language and thank to his web-oriented approach I quickly came up with one of the simplest single task web application I could write: a URL shortener.

What is a URL shortener? It&#8217;s a service that will give you a short link for a long URL.

Why should I use it? A short URL is easier to remember and to copy and paste, it lets you write more on social media where characters are limited (Twitter).

Why writing one when there are plenty of them already? Other shortener have an expiration date on short links, <a href="http://4pr.es" target="_blank">4pr.es</a> doesn&#8217;t!

Take a look at the <a href="http://4pr.es/rwfiwK" target="_blank">code</a> and you&#8217;ll see it is very simple: it takes a URL as text input filed of a form

<pre><span class="pl-smi">short</span>, <span class="pl-smi">err</span> <span class="pl-k">:=</span> <span class="pl-c1">createUrl</span>(req.<span class="pl-c1">FormValue</span>(<span class="pl-s"><span class="pl-pds">"</span>url<span class="pl-pds">"</span></span>))
</pre>

generates a random string of 6 charachters, checking that the string is not in the database

<pre>...
for urlPresent(coder.Shrt) {
        coder.Shrt = shorten(coder.Length)
}
...
func shorten(c uint) string {
    rand.Seed(time.Now().UnixNano())
    b := make([]rune, c)
    for i := range b {
        b[i] = letters[rand.Intn(len(letters))]
    }
    return string(b)
}</pre>

and keep the URL &#8211; string association in the database for further redirection; once the random short URL is visited the client gets redirected to the origianl long URL!

<pre>...
err := getUrl(params["short"], w, req)
...
func getUrl(short string, w http.ResponseWriter, req *http.Request) error {
    var redir string
    err := db.QueryRow("SELECT url FROM urls WHERE short = ?", short).Scan(&redir)
    if err != nil {
        return err
    }
    http.Redirect(w, req, redir, 301)
    return nil
}</pre>

So please try it and if you have any suggestion to increase performance or security please let me know!

Cheers

&nbsp;

&nbsp;
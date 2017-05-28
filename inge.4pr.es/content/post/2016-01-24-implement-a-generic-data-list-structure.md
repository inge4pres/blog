---
title: Implement a generic data list structure
author: inge4pres
layout: post
date: 2016-01-24T16:40:37+00:00
categories:
  - social
  - tech
tags:
  - aws
  - challenge
  - coding
  - design
  - golang
  - software

---
As a coding challenge I was asked to provide a generic list implementation using a language of my choice and using only primitive types, avoiding the use of high level built-ins. I chose <a href="https://golang.org/" target="_blank">Go</a> because I want to learn it and I know it can be useful to create an abstract, generic implementation.

The challenge request to implement at least 4 methods on the generic type:

  * Filter() &#8211;  returns a subset of the List satisfying an operation
  * Map() &#8211; returns the List objects&#8217; map
  * Reverse() &#8211; reverse the ordering of the List objects
  * FoldLeft() &#8211; join the objects from left to right using a join character

As a bonus question I was asked to code unit tests for the aforementioned methods and give an explanation on how the implementation guarantees concurrent access on resources.

So here is my implementation: the type List has only one attribute, an array of type `interface{}`

<pre class="theme:sublime-text font:consolas lang:go decode:true" title="List">type List struct {
    data []interface{}
}</pre>

Every type will be convertible to the `interface{}` type, but as Golang has strong types the conversion is not implicit and must be declared.
  
Reverse() will create a new array of `interface{}` to hold the reversed list

<pre class="theme:sublime-text font:consolas lang:go decode:true" title="Reverse">func (m *List) Reverse() *List {
	var ret []interface{}
	for i := 1; i &lt;= len(m.data); i++ {
		ret = append(ret, m.data[len(m.data)-i])
	}
	return &List{
		data: ret,
	}
}</pre>

Map() returns the List elements&#8217; array, so it can be accessed as a whole

<pre class="theme:sublime-text font:consolas lang:go decode:true " title="Map">func (m *List) Map() []interface{} {
	return m.data
}
</pre>

This two function types will help define a custom operation to be used in Filter() and FoldLeft(): <a href="http://jordanorelli.com/post/42369331748/function-types-in-go-golang" target="_blank">functions are types</a> in Go and this enable a great level of abstraction.

<pre class="theme:sublime-text font:consolas lang:go decode:true " title="">type filterFn func(interface{}) interface{}
type foldFn func([]interface{}) interface{}
</pre>

Filter() will use a filter function, without the need to define it (!), and return a portion of the List data array.

<pre class="theme:sublime-text font:consolas lang:go decode:true " title="Filter">func (m *List) Filter(filter filterFn) *List {
	var ret []interface{}
	for d := range m.data {
		if data := filter(m.data[d]); data != nil {
			ret = append(ret, data)
		}
	}
	return &List{
		data: ret,
	}
}</pre>

FoldLeft() will use a fold function, again not yet defined, the return a single element made of the entire list.

<pre class="theme:sublime-text font:consolas lang:go decode:true " title="FoldLeft">func (m *List) FoldLeft(fold foldFn) *List {
	var ret []interface{}
	ret = append(ret, fold(m.data))
	return &List{
		data: ret,
	}
}</pre>

You can find all the code <a href="https://github.com/inge4pres/blog/tree/master/implement-a-generic-data-list-structure" target="_blank">here</a>, any comment is welcome on how to improve the abstraction or efficiency of the implementation.
  
The opportunity to dig into the language ability to abstract is a very helpful way to better understand the language itself, so this coding challenge was a great opportunity to learn a little bit more Go!
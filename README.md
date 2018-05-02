# yuuri

[![Build Status](https://travis-ci.org/johynpapin/yuuri.svg?branch=master)](https://travis-ci.org/johynpapin/yuuri)
[![Coverage Status](https://coveralls.io/repos/github/johynpapin/yuuri/badge.svg?branch=master)](https://coveralls.io/github/johynpapin/yuuri?branch=master)
[![GoDoc](https://godoc.org/github.com/johynpapin/yuuri?status.svg)](https://godoc.org/github.com/johynpapin/yuuri)

Repository used for my internship at the NII.

## The project

I am doing this project as part of my internship at the NII.

The basic idea is to enrich a Linked Data database. This database contains recipes and ingredients. The aim here is to add nutritional and price information to this database.

## Packages

In order to clarify the structure of this project, but also to facilitate the reuse of what can be, some modules are accessible as packages. You are free to use them if you have the use.

These packages are accessible in the [pkg](pkg) folder, as recommended by the [Standard Go Project Layout](https://github.com/golang-standards/project-layout).

### agrovoc

Package [agrovoc](pkg/agrovoc) offers a simple way to search for a concept on agrovoc.

### rdfgraph

Package [rdfgraph](pkg/rdfgraph) allows to represent RDF graphs, to decode them and to encode them.

## FAQ

### Why "yuuri"?

Yuuri is the name of a character who loves to eat in the manga [Shoujo Shuumatsu Ryoukou](http://girls-last-tour.com/). However, this project concerns cooking recipes...
package main

var srcMod = `foo=bar


hello=LET h = "Hello" IN k0 h "disc\"arded"


world =
  LET d _ = "ditched"
      w = "World"
  IN k1 (d '?') w



/*
helloOrWorld h0w1 =
  LET foo = h0w1
      h = hello
      w = world
  IN CASE foo OF
    H -> h
    W -> w
*/

moo = "bar"
`

package main

var srcMod = `foo="bar"

moo=321

foo = "bar"

world =
  (k1 "ditched") "World" (yo 123)



/*
helloOrWorld h0w1 =
  let foo = h0w1
      h = hello
      w = world
  in case foo of
    H -> h
    W -> w
*/

hello =
  k0 "Hello" // "discarded"


`

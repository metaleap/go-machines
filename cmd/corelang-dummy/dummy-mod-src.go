package main

var srcMod = `

hello =
  (k0 "Hello") // "(discarded)"

foo = "bar"

world =
  (k1 "(ditched)") // "World"


/*
helloOrWorld h0w1 =
  let foo = h0w1
      h = hello
      w = world
  in case foo of
    H -> h
    W -> w
*/
`

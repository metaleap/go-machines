package main

var srcMod = `

hello =
  k0 "Hello" "(discarded)"

foo = "bar"

world =
  k1 "(ditched)" "World"



helloOrWorld h0w1 =
  let foo = case h0w1 of
    H -> 0
    W -> 1
  in case foo of
    0 -> hello
    1 -> world
`

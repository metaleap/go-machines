package main

var srcMod = `foo=bar


hello=let h = "Hello" in k0 h "disc\"arded"


world =
  let d _ = "ditched"
      w = "World"
  in k1 (d '?') w



/*
helloOrWorld h0w1 =
  let foo = h0w1
      h = hello
      w = world
  in case foo of
    H -> h
    W -> w
*/

moo = "bar"
`

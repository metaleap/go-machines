package main

const srcMod = `


pair x y f = f x y

fst p = p k0

snd p = p k1

f x y =
    LET REC a = pair x b
            b = pair y a
    IN fst (snd (snd (snd a)))

main = f 3 4

main2 = LET REC f = f x IN f


// random noisy rubbish..

foo=bar


hello = LET
        h = "Hello"
        IN k0 h "disc\"arded"


world =
  LET d _ = "ditched"
      w = "World"
  IN k1 (d '?') w



helloOrWorld h0w1 =
  LET foo = h0w1
      h = hello
      w = world
  IN CASE foo OF
    0 -> h
    1 -> w


moo = "bar"
`

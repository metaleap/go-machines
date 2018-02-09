package main

const srcMod = `


pair l r f = f l r

fst p = p k0

snd p = p k1

cons a b cc cn = cc a b
nil cc cn = cn
hd list = list k0 abort
tl list = list k1 abort
abort = abort
infinite n = cons n (infinite n)
listish = hd (tl (infinite 4))


freakish x y =
    LET REC a = pair x b
            b = pair y a
    IN fst (snd (snd (snd a)))

main0 = freakish 3 4

main1 k = LET
        p = pair 123
        pp = p 456
        fun = k
    IN pp fun

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

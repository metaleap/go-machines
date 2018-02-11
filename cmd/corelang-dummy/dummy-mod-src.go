package main

// https://www.youtube.com/watch?v=hrBq8R_kxI0 8m
// https://www.youtube.com/watch?v=GhERMBT7u4w 21m

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


fac n = if (n==0) 1 (n * (fac (n - 1))) // 'when' instead of 'if' will work (executes approx. ~25-40% more steps)

checkifopish = 3 × (4 ÷ 5)


when cond then else =
    CASE cond OF
    0 -> 42
    1 -> else
    2 -> then


test ctor =
    CASE ctor OF    0 -> 42
                    1 n -> neg n
                    2 x y -> x + y
                    3 x y -> x * y


do = test ((3 2) 5 3) // call to test with ctor of (3 2) returns the result of 5*3


Ycomb f = LET REC x = f x IN x


main1 k = LET REC
        pa = pair
        pp = LET n = 123 IN pa n
        fun = k0 k k
    IN (pp 567) fun

main2 = LET REC f = f x IN f


// random noisy rubbish..

foo=bar


hello = LET
        h = "Hello"
        IN k0 h "disc\"arded"


world =
  LET d = "ditched"
      w = "World"
  IN k1 (k0 d "ditched") w



helloOrWorld h0w1 =
  LET foo = h0w1
      h = hello
      w = world
  IN CASE foo OF
    0 -> h
    1 -> w


moo = "bar"
`

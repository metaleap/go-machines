package main

// https://www.youtube.com/watch?v=hrBq8R_kxI0 8m
// https://www.youtube.com/watch?v=GhERMBT7u4w 21m

const srcMod = `

page136 x =
    LET foo = CASE x OF 1 -> 111
                        2 -> 222
    IN (2 2) (123 0) foo

p136 = page136 (2 0)


page137_1 x = (99 3) 123 x

p137_1 = (page137_1 456) 789


page137_2 incompletector = incompletector 654

p137_2 = (page137_2 ((99 3) 321)) 987



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


fac n = when (n==0) 1 (n * (fac (n - 1))) // 'when' instead of 'if' executes approx. ~20-30% more steps & appls

checkIfLexedOpish = 3 × (4 ÷ 5)

when cond then else =
    CASE cond OF
    1 -> else
    2 -> then



test ctor =
    CASE ctor OF    1 n -> neg n
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

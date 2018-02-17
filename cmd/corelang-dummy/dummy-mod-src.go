package main

// https://www.youtube.com/watch?v=hrBq8R_kxI0 8m
// https://www.youtube.com/watch?v=GhERMBT7u4w 21m

const srcMod = `

page136 x =
    LET foo = CASE x OF 12 -> 111
                        34 -> 222
    IN (2 2) foo ((234 1) 77)

z x = CASE x OF 12 -> 444
                34 -> 555
zz= (1234567890 1) (z (12 0))

p136 = page136 (34 0)


page137_1 x = (99 5) 123 x 333

p137_1 =
    LET REC
        p = page137_1
        oo = p 456 666
    IN oo 789


page137_2 incompletector = incompletector 654

p137_2 = (page137_2 ((88 3) 321)) 987



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

when cond then else =
    CASE cond OF
    1 -> else
    2 -> then


fac n =                         // using 'when' instead of 'if' here works equivalently: but executes ~20-30% more steps and ~20-30% more calls; plus tends to take ~2x as long
    if (n==0)
    /*then*/ 1
    /*else*/ (n * (fac (n - 1)))



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


checkIfLexedOpish = 3 × (4 ÷ 5)

moo = "bar"
`

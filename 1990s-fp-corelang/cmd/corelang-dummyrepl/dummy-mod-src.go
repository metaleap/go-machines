package main

// https://www.youtube.com/watch?v=hrBq8R_kxI0 8m
// https://www.youtube.com/watch?v=GhERMBT7u4w 21m

const srcMod = `
map fn lst =
    CASE lst OF
        Nil -> (Nil 0)
        Cons x xs -> (Cons 2) (fn x) (map fn xs)

pow n = n * n

times n foo =
    when (n == 0)
    /*then*/ (N 0)
    /*else*/ ((C 2) foo (times (n-1) foo))

t4 = times 4 123

powlst =
    LET nums = ((Cons 2) 12 ((Cons 2) 34 ((Cons 2) 56 ((Cons 2) 78 (Nil 0)))))
    IN map pow nums


page136 x =
    LET foo = CASE x OF Foo -> 111
                        Bar -> 222
                        _ -> 42
    IN (P136 2) foo ((Wot 1) 77)

p136 = page136 (Nope 0)

z x = CASE x OF One -> 444
                Two -> 555
zz= (Zz 1) (z (One 0))

ctpar foo = (Partial 4) 22 foo

ctp = ctpar 33 44 55

page137_1 x = (P137_1 5) 123 x 333

p137_1 =
    LET REC
        p = page137_1
        oo = p 456 666
    IN oo 789


page137_2 incompletector = incompletector 654

p137_2 = (page137_2 ((P137_2 3) 321)) 987



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
    False -> else
    True -> then


fac n =                         // using 'when' instead of builtin 'if' here executes ~15-25% more steps and ~15-25% more calls; plus tends to take ~2x as long
    when (n==0)
    /*then*/ 1
    /*else*/ (n * (fac (n - 1)))



test ctor =
    CASE ctor OF    Neg n -> neg n
                    Add x y -> x + y
                    Mul x y -> x * y

do = test ((Mul 2) 5 3) // call to test with ctor of (3 2) returns the result of 5*3





Ycomb f = LET REC x = f x IN x


main1 k = LET REC
        pa = pair
        pp = LET n = 123 IN pa n
        fun = k0 k k
    IN (pp 567) fun


// random noisy rubbish..

foo=bar



checkIfLexedOpish = 3 × (4 ÷ 5)

`

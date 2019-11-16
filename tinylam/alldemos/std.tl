
Enum      := Foo | Bar | Baz
Bool      := True | False
Maybe     := Ok value | None
Dual      := Neat neat | Weak weak
Tuple     := Pair fst snd
List      := ListEnd | ListLink h (List t)
ToString  := StrNum n | StrOther o

// Enum toString it := it "Foo" "Bar" "Baz"
// Bool toString it := it "True" "False"
// Tuple toString it := it (pair -> pair (fst snd -> ** concat ["(Pair " fst.toString " " snd.toString ")"]))

same it   := it   // I combinator aka 'id' / identity
ever      := True // K combinator aka 'konst'

and b1 b2 := b1 b2 False
or b1 b2  := b1 True b2
xor b1 b2 := b1 b2.not b2

not       := _ False True
not:      := <· not

if        := same
else      := same

nEQ v1 v2 := not (v1.EQ v2)
gEQ v1 v2 := or (v1.GT v2) (v1.EQ v2)
lEQ v1 v2 := or (v1.LT v2) (v1.EQ v2)





/> it fn    := fn it

>> it fns   := fns it (item rest -> >> (item it) rest)

// flip args for `fn`
<> fn       := a1 a2 -> fn a2 a1

// compose right to left
<· f2 f1 a  := f2 (f1 a)

// compose left to right
·> f1 f2 a  := f2 (f1 a)

// see std.num.sum or std.num.product for the point of this =)
** f2 args  := args.>-.forth args.-< f2

__msgPanic  := "forced crash via `()`"

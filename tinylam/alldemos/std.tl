
Enum          := Foo | Bar | Baz
Bool          := True | False
May t         := None | That t
Tried yay nay := Yay yay | Nay nay
Tuple one two := Pair one two
List t        := ListNil | ListLink t (List t)

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

>> it fns   := fns it (hd tl -> >> (hd it) tl)

// flip args for `fn`
<> fn       := a1 a2 -> fn a2 a1

// compose right to left
<· f2 f1 a  := f2 (f1 a)

// compose left to right
·> f1 f2 a  := f2 (f1 a)

// like a fold/reduce without an initial zero val, hence assumes a non-nil `l`
** binOp l  := l.iter
  iter args := args.>- args.-< (_ _ -> binOp args.-< ( iter args.>- ))












//begin interpreter-required internal defs used for syntax sugars

__tlListIsNil list    := list True (_ _ -> False)
__tlListIsntNil list  := list False (_ _ -> True)

//end

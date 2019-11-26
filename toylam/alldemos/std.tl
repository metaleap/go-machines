
// pseudo sum types: are not really "type-checked" but help signify intents
// and data flows, and create additional representative constructor funcs
// Scott-encoding these, plus allow for case-scrutinizing syntax-sugar with
// constructor exhaustiveness checked.

Bool            := True | False
Could           := Have _ | None
Result          := Yay _ | Nay _
Tuple           := Pair fst snd
List            := ListEnd | ListLink _ _



// I combinator aka 'id' / identity
same it         := it

// K combinator aka 'konst'
ever            := True

and b1 b2       := if (b1) b2 (else False)
or b1 b2        := if (b1) True (else b2)
xor b1 b2       := if (b1) b2.not (else b2)

not             := if _ False (else True)
not·            := ·· not

if              := same
then            := same
else            := same

// `when` is like `if` with its first 2 args flipped, allowing for handy dot-notation
// where sensible: `dis.when foo (else dat)` equiv. to `if foo dis (else dat)`
when            := <> if



nEQ v1 v2       := not (v1.EQ v2)
gEQ v1 v2       := or (v1.GT v2) (v1.EQ v2)
lEQ v1 v2       := or (v1.LT v2) (v1.EQ v2)



cases chks it   := chks ? [] => None | .. => {chk ifso} rest -> it.chk (Have it.ifso) (cases rest it)



// helper for REPL sessions: given a Could or Result, force-crashes for the
// failure case or extracts and returns the contained result of the success case
?!              := _ same ()



/> it fn        := fn it

>> it fns       := fns it (item rest -> >> (item it) rest)

// flip args for `fn`
<> fn a1 a2     := fn a2 a1

// compose right to left
·· f2 f1 a      := f2 (f1 a)

// compose right to left with 2 input args
··· f2 f1 a1 a2 := f2 (f1 a1 a2)

// see std.num.sum or std.num.product for the point of this =)
** f args       := args.>-.forth args.-< f

incr            := ADD 1
decr            := ADD -1
sum             := ** ADD
product         := ** MUL
neg             := SUB 0
abs i           := if (i.LT 0) i.neg (else i)
min x y         := x.when (x.LT y) (else y)
max x y         := x.when (x.GT y) (else y)
pow i exp       := if (exp.GT -1) exp.mults (else ERR "std.num.pow called with negative exp:" exp)
  mults e       := if (e.EQ 0) 1 (else i.MUL (mults e.decr))
isEven i        := EQ 0 (i.MOD 2)
isOdd i         := i.sign.EQ (i.MOD 2)
sign i          := if (i.LT 0) -1 (else 1)
factorial i     := if (i.EQ 0) 1 (else i.MUL (factorial (i.SUB 1)))
fibonacci i     := if (i.LT 2) i (else ADD (fibonacci (i.SUB 1)) (fibonacci (i.SUB 2)))

toString numval := if (numval.GT -1) numval.tostr.reversed (else +> '-' numval.abs.tostr.reversed)
  tostr i       := if (i.LT 10) [i.tochr] (else +> (tochr (i.MOD 10)) (tostr (i.DIV 10)) )
  tochr d       := d.ADD 48 // d is always >= 0 and < 10, char-code for '0' is 48..

parseSigned str := str ? [] => Nay "empty input" | .. => unsign
  unsign c rest := parsed (if (c.EQ '-') rest (else str))
  parsed s      := s.parseUnsigned Result? Nay => Nay | Yay => ok
  ok numval     := Yay (numval.neg.when (str.-<.EQ '-') (else numval))

parseUnsigned s := s ? [] => Nay "empty input" | .. => _ _ -> checkresult
  checkresult   := badchar Could? None => Yay finalresult | Have => bad -> Nay (bad.neg.+> " isn't between 0-9")
  badchar       := digits.findWhere (LT _ 0) // bad chars are marked thusly by `chartodigit`
  finalresult   := digits.more s.len.decr
  digits        := s.each chartodigit
  more l n      := ADD (l.-<.MUL (pow 10 n)) (l.>- 0 (_ _ -> more l.>- n.decr))
  chartodigit c := (and (c.GT 47) (c.LT 58)) (c.SUB 48) (else c.neg) // will be caught by `badchar`

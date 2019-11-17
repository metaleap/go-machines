incr            := ADD 1
decr            := ADD -1
sum             := ** ADD
product         := ** MUL
neg             := SUB 0
abs i           := if (i.LT 0) i.neg (else i)
min x y         := x.when (x.LT y) (else y)
max x y         := x.when (x.GT y) (else y)
pow i exp       := if (exp.LT 0) (ERR exp) (else mults exp)
  mults e       := if (e.EQ 0) 1 (else i.MUL (mults e.decr))
isEven i        := EQ 0 (i.MOD 2)
isOdd i         := i.sign.EQ (i.MOD 2)
sign i          := if (i.LT 0) -1 (else 1)
factorial i     := if (i.EQ 0) 1 (else i.MUL (factorial (i.SUB 1)))
fibonacci i     := if (i.LT 2) 1 (else ADD (fibonacci (i.SUB 1)) (fibonacci (i.SUB 2)))

parseSigned str := str ? [] => ERR "empty input" | => _ _ -> signdropped
  signdropped   := if (EQ '-' str.-<) (neg (parseSigned str.>-)) (else unsigned)
  unsigned      := str.parseUnsigned

parseUnsigned s := s ? [] => ERR "empty input" | => _ _ -> digits.more s.len.decr
  digits        := s.each chartodigit
  more l n      := ADD (l.-<.MUL (pow 10 n)) (l.>- 0 (_ _ -> more l.>- n.decr))
  chartodigit c := (and (c.GT 47) (c.LT 58)) (c.SUB 48) (ERR (c.+> " isn't between 0-9") )

toString n      := if (n.LT 0) (+> '-' n.abs.tostring) (else n.tostring)
  tostring i    := i.digitchar.when (i.LT 10) (else concat (tostring (i.DIV 10)) (digitchar (i.MOD 10)))
  digitchar d   := [d.ADD 48]

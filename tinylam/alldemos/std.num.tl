incr            := ADD 1
decr            := ADD -1
neg             := SUB 0
abs n           := if (n.LT 0) n.neg (else n)
pow exp n       := if (exp.LT 0) (ERR exp) (else n_pow exp)
  n_pow e       := if (EQ 0 e) 1 (else (n.MUL (n_pow e.decr)))
sum             := ** ADD
product         := ** MUL
min x y         := if (x.LT y) x (else y)
max x y         := if (x.GT y) x (else y)
isEven n        := EQ 0 (n.MOD 2)
isOdd n         := EQ (if n.LT 0 -1 1) (n.MOD 2)

parseUnsigned s := if !s (ERR "empty input") (else digits.iter digits.len.decr)
  digits        := s.each chartodigit
  iter l n      := ADD (l.-<.MUL (pow n 10)) (if !l.>- 0 (else iter l.>- n.decr))
  chartodigit c := (and (c.GT 47) (c.LT 58)) (c.SUB 48) (ERR ({c} " isn't between 0-9") )

parseSigned str := if !str (ERR "empty input") (else sanssign)
  sanssign      := if (EQ '-' str.-<) (neg (parseSigned str.>-)) (else unsigned)
  unsigned      := str.parseUnsigned

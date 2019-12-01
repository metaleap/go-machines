main args env := m1 args env

m7 args env   := pow 10 args.len

m4 args env   := concat (+> \t s1) "!"
  \t          := 9
  s1          := s0 l
  s0          := num.toString
  l           := env.len

m0 args _     := num.factorial args.len

m5 args _     := if (n.EQ 0) (Nay "no args") (else Yay (concat n.toString " arg(s)")) Result? Yay => yo | Nay => no
  n           := args.len
  yo          := concat "neat: "
  no          := concat "weak: "

m1 args _     := (num.parseSigned (args.at 0)) Result? Yay => n -> concat "Num: " (num.toString n) | Nay => same

m2 args env   := if (args.len.EQ 0) env.len (else max (len (env.at 0)) (len (args.at 0)))

m3 args _     := (cases [{ LT _ 8 , ever 111 }, { GT _ 23 , ever 777 }, { ever True, ever 444 }] 17) same 999

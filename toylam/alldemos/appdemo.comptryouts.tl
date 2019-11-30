main args env := m4 args env

m4 args env   := concat (+> \t s1) "!"
  \t          := 9
  s1          := s0 l
  s0          := num.toString
  l           := env.len

m0 _ _ := num.factorial 5

m1 _ _      := num.parseUnsigned "7" Result? Yay => num -> "yo" | Nay => msg -> "no"

m2 args env := if (args.len.EQ 0) env.len (else max (len (env.at 0)) (len (args.at 0)))

m3 args _   := (cases [{ LT _ 8 , ever 111 }, { GT _ 23 , ever 777 }, { ever True, ever 444 }] 7) same 999

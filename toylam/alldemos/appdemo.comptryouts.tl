main args env := m13 args env


m10 args _    := args 777 (arg _ -> 555)

m8 args _     := args ? [] => "arg missing" | .. => arg1 _ -> (arg1.EQ "foo") "got foo" "not foo"

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

m3 args _     := (cases [{ LT _ 8 , ever 111 }, { GT _ 23 , ever 777 }, { ever True, ever 444 }] args.len) same 999

m6 args _     := [on, 0, -0, ""]
  on _ jsrc   := [same , (json.parsed jsrc) Result? Yay => json.stringified | Nay => msg -> (concat "BAH: " msg)]

m9 args _     := (json.parsed (args.at 0)) Result? Yay => _ -> "neat" | Nay => msg -> (concat "!?" msg)


m11 args _    := [on, 0, -0, ""]
  on _ jsrc   := [same , num.toString (len ( jsrc)) ]

m12 args _    := num.parseUnsigned (args.at 0) Result? Nay => _ -> 123 | Yay => n -> n

m13 args _    := +> '!' (args ? [] => "?" | .. => hd tl -> hd)

main args env := num.parseUnsigned "1"

// m1 _ _      := num.parseUnsigned "7" Result? Yay => "same" | Nay => msg -> "msg"

// m2 args env := if (args.len.EQ 0) env.len (else max (len (env.at 0)) (len (args.at 0)))

// m3 args _   := (cases [{ LT _ 8 , ever 111 }, { GT _ 23 , ever 777 }, { ever True, ever 444 }] 7) same 999

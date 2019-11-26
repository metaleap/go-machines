main args env := if (args.len.EQ 0) env.len (else max (len (env.at 0)) (len (args.at 0)))

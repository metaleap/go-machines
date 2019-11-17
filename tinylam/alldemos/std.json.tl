
Token             := Ident | NumLit | StrLit | SepComma | SepColon | CurlyOpen | CurlyClose | SquareOpen | SquareClose
Value             := Null | Boolean _ | Int _ | Str _ | Arr (List _) | Obj (List (Pair k v))

tokenized src     := src.more
  more s          := s ? => cur rest -> if (cur.iswhitespace) (more rest) (else None)
  iswhitespace c  := ** or [c.EQ 32 , and (c.GT 7) (c.LT 12) , c.EQ 13]

parse src         := src ? [] => Nay "empty input" | => _ _ -> src.tryparse
  tryparse s      := s.attempt [pnull, pbool, pnum]
  attempt s trys  := trys ? [] => Nay (errmsg s) | => cur rest -> ((cur s) Could? Have => Yay | None => attempt s rest)
  pnull toknull   := if (toknull.EQ "null") (Have Null) (else None)
  pbool tokb      := (tokb.EQ "true") (Have (Boolean True)) ((tokb.EQ "false") (Have (Boolean False)) None)
  pnum toknum     := (num.parseSigned toknum) Result? Yay => ·· Have Int | Nay => None
  errmsg tokstr   := concat.** ["unrecognized token: `" , tokstr ,"`"]

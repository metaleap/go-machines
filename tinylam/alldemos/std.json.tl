
Token             := TIdent _ | TNumLit _ | TStrLit _ | TSep commaBool | TBracket openBool curlyBool
Value             := VNull | VBoolean _ | VInt _ | VStr _ | VArr (List _) | VObj (List (Pair k v))

tokenized src     := src.more
  more s          := s ? [] => None | => cur rest -> if (cur.iswhitespace) (more rest) (else None)
  iswhitespace c  := or (and (c.LT 12) (c.GT 7)) (or (c.EQ 32) (c.EQ 13))

parse src         := src ? [] => Nay "empty input" | => _ _ -> src.tryparse
  tryparse s      := s.attempt [pnull, pbool, pnum]
  attempt s trys  := trys ? [] => Nay (errmsg s) | => cur rest -> ((cur s) Could? Have => Yay | None => attempt s rest)
  pnull toknull   := if (toknull.EQ "null") (Have VNull) (else None)
  pbool tokb      := (tokb.EQ "true") (Have (VBoolean True)) ((tokb.EQ "false") (Have (VBoolean False)) None)
  pnum toknum     := (num.parseSigned toknum) Result? Yay => ·· Have VInt | Nay => None
  errmsg tokstr   := concat.** ["unrecognized token: `" , tokstr ,"`"]

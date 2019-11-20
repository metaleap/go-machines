
Token := TSepColon | TSepComma | TBracketSq _ | TBracketCy _ | TNumLit _ | TStrLit _ | TIdent _ | __TTmp _ _
Value := VNull | VBoolean _ | VInt _ | VStr _ | VArr (List _) | VObj (List (Pair k v))



tokenized src     := src ? .. => more
  more chr rest   := (simple.cases chr) Could? Have => +> _ rest.tokenized | None => others chr rest
  simple          := [{EQ ',', TSepComma.ever}, {EQ '[', True.TBracketSq.ever}, {EQ ']', False.TBracketSq.ever}, {EQ '{', True.TBracketCy.ever}, {EQ '}', False.TBracketCy.ever}, {EQ ':', TSepColon.ever}]
  others c r      := c.checknum (casenum c r) (c.checkname (casename c r) (c.checkstr (casestr c r) (tokenized r)))
  casenum c rest  := (rest.brokenUpOn (_ -> checknum.not·)) (fst snd -> +> (TNumLit (c.+> fst)) (tokenized snd))
  casename c rest := (rest.brokenUpOn (_ -> checkname.not·)) (fst snd -> +> (TIdent (c.+> fst)) (tokenized snd))
  casestr c rest  := (rest.brokenUpOn (a c´ -> and c´.checkstr (a.last2 ? [] => True | .. => c1 c2 -> or (c2.-<.nEQ 92) (c1.EQ 92)))) (fst snd -> +> (TStrLit (concat (c.+> fst) [34])) (snd ? .. => _ r -> tokenized r))
  checknum chr    := or (chr.EQ '-') (and (chr.GT 47) (chr.LT 58))
  checkname chr   := and (chr.GT 96) (chr.LT 123)
  checkstr chr    := chr.EQ 34


parsed src        := src.tokenized ? [] => Nay "empty input" | .. => ··· handle parse
  treeified toks  := toks.rewritten (t -> t Token? TBracketSq => not | => False) (t -> t Token? TBracketSq => same | => False) (inner -> [inner.where (t -> t Token? TSepComma => False | => True )])
  handle          := {tail ret} -> tail ? [] => ret | .. => _ _ -> Nay "trailing tokens remaining after completed expression"
  parse tok rest  := tok Token? TNumLit => pnum rest | TStrLit => pstr rest | TIdent => pident rest | => {[], Nay "unexpected token"}

  pident r t      := { r, (t.EQ "null") (Yay VNull) ((t.EQ "true") (Yay True.VBoolean) ((t.EQ "false") (Yay False.VBoolean) (Nay t))) }
  pnum r t        := { r, t.parseSigned Result? Nay => Nay | Yay => ·· Yay VInt }
  pstr r t        := { r, t ? [] => Nay t | .. => quotemark rest -> Yay rest.sansLast.VStr } // caution, not handling escapes for now..
  parr p r t      := if (t.not) {[], Nay "]"} (else arr p r) // r is all that follows opening `[`
    arr p ts      := ts ? [] => {[], Nay "array non-terminated"} | .. => arr´ p
    arr´ p t r    := t Token? => {[], Nay "todo"}

  pobj p r t      := { r, Nay "notyet" }



toString json     := json Value? VNull => casenull | VBoolean => casebool | VInt => caseint | VStr => casestr | VArr => casearr | VObj => caseobj
  casenull        := "null"
  casebool        := _ "true" "false"
  caseint         := num.toString
  casestr         := quoted
  casearr l       := concat (+> '[' (joinedVia (l.each toString) ", ")) "]"
  caseobj l       := concat (+> '{' (joinedVia (l.each ({k v} -> joinedVia [k.quoted, toString v] ": ")) ", ")) "}"
  quoted s        := concat (+> 34 s) [34]

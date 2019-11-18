
Token := TSepColon | TSepComma | TBracketSq _ | TBracketCy _ | TNumLit _ | TStrLit _ | TIdent _
Value := VNull | VBoolean _ | VInt _ | VStr _ | VArr (List _) | VObj (List (Pair k v))


parse src         := src ? [] => Nay "empty input" | => _ _ -> src.tokenized.tryparse
  tryparse toks   := toks


tokenized src     := src ? => more
  more chr rest   := (simple.cases chr) Could? Have => +>.<> (tokenized rest) | None => others chr rest
  simple          := [{EQ ',', TSepComma.ever}, {EQ '[', True.TBracketSq.ever}, {EQ ']', False.TBracketSq.ever}, {EQ '{', True.TBracketCy.ever}, {EQ '}', False.TBracketCy.ever}, {EQ ':', TSepColon.ever}]
  others c r      := c.checknum (casenum c r) (c.checkname (casename c r) (c.checkstr (casestr c r) (tokenized r)))
  casenum c rest  := (rest.brokenUpOn (_ -> checknum.not·)) (fst snd -> +> (TNumLit (c.+> fst)) (tokenized snd))
  casename c rest := (rest.brokenUpOn (_ -> checkname.not·)) (fst snd -> +> (TIdent (c.+> fst)) (tokenized snd))
  casestr c rest  := (rest.brokenUpOn (a c´ -> and c´.checkstr (a.last2 ? [] => True | => c1 c2 -> or (c2.-<.nEQ 92) (c1.EQ 92)))) (fst snd -> +> (TStrLit (concat (c.+> fst) [34])) (snd ? => _ r -> tokenized r))
  checknum chr    := or (chr.EQ '-') (and (chr.GT 47) (chr.LT 58))
  checkname chr   := and (chr.GT 96) (chr.LT 123)
  checkstr chr    := chr.EQ 34

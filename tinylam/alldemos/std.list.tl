
-< list                 := list (ERR "forgot nil check before -<") True
>- list                 := list (ERR "forgot nil check before >-") False
+>                      := ListLink

len list                := list 0 (_ tail -> incr (len tail))

at list i               := if (EQ 0 i) list.-< (else at list.>- i.decr)

each list fn            := list.iter
  iter l                := l [] (head tail -> head.fn.+> (iter tail))

slice list start length := take (list.drop start) length

drop list n             := if (EQ 0 n) list (else drop list.>- n.decr)

take list n             := if (EQ 0 n) [] (else list.-<.+> (take list.>- n.decr))

where list isMatch      := list.iter
  iter l                := l [] (h t -> (if h.isMatch (h.+>) (else same)) (iter t))

whereNot list isMatch   := list.where (not: isMatch)

without list item       := list.where (nEQ item)

findWhere list isMatch  := list.iter
  iter l                := l same (h t -> if (h.isMatch) h.ever (else iter t))

countWhere list isMatch := list.iter
  iter l                := l 0 (head tail -> ADD (head.isMatch 1 0) (iter tail))

has list item           := list.hasSome (EQ item)

hasSome list isMatch    := list.iter
  iter l                := l False (head tail -> or (head.isMatch) (iter tail))

hasNone list isMatch    := list.hasOnly (not: isMatch)

hasOnly list isMatch    := list False list.iter
  iter l                := if (not l.-<.isMatch) False (else or (!l.>-) (iter l.>-))

sortedBy list isLess    := list.iter
  iter l                := l [] (hd tl -> concat (iter (filt same hd tl)) (hd.+> (iter (filt not hd tl))))
  filt how hd tl        := tl.where (item -> how (isLess item hd))

one tup                 := tup (l r -> l)
two tup                 := tup (l r -> r)

dividedBy list isMatch  := { (items where) (items whereNot) }
  items how             := each (tagged.how two) one
  tagged                := list.iter
  iter l                := l [] (hd tl -> +> {hd hd.isMatch} (iter tl) )

deduped list            := list.iter
  iter l                := l [] (hd tl -> (if (tl.has hd) same (else hd.+>)) (iter tl))

last list               := list.>- list.-< (_ _ -> last list.>-)

indexed list            := pairUp list.indices list

indices list            := list.iter 0
  iter l i              := l [] (hd tl -> i.+> (iter tl i.incr))

pairUp l1 l2            := l1 [] (hd1 tl1 -> l2 [] (result hd1 tl1))
  result h1 t1 h2 t2    := +> {h1 h2} (pairUp t1 t2)

concat l1 l2            := l2 l1 (_ _ -> l1.prepend)
  prepend l             := l l2 (hd tl -> hd.+> (prepend tl))

forth list z binOp      := list.acc z
  acc l cur             := l cur (hd tl -> acc tl (binOp cur hd))

from list z binOp       := list.acc
  acc l                 := l z (hd tl -> binOp hd (acc tl))

reverse list            := acc [] list
  acc cur l             := l cur (hd tl -> acc (hd.+> cur) tl)

mGet m key              := m.iter
  iter l                := l same (hd tl -> hd (k v ->  if (k.EQ key) v.ever (else iter tl)))

mDrop m key             := m.iter
  iter l                := l [] (hd tl -> hd (k _ -> (if (k.EQ key) same (else hd.+>))) (iter tl))

mWith m key value       := +> {key value} (m.mDrop key)

mHasKeyWhere m isMatch  := m.hasSome (Tuple (k v -> k.isMatch))
mHasValWhere m isMatch  := m.hasSome (Tuple (k v -> v.isMatch))

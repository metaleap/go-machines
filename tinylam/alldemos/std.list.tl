
-< list                 := list (ERR "forgot nil check before -<") True
>- list                 := list (ERR "forgot nil check before >-") False
+>                      := ListLink
-| list                 := list []

len list                := list 0 (_ rest -> incr (len rest))

at list i               := if (EQ 0 i) list.-< (else at list.>- i.decr)

each list fn            := list.more
  more l                := l.-| (item rest -> item.fn.+> (more rest))

slice list start length := take (list.drop start) length

drop list n             := if (EQ 0 n) list (else drop list.>- n.decr)

take list n             := if (EQ 0 n) [] (else list.-<.+> (take list.>- n.decr))

where list isMatch      := list.more
  more l                := l.-| (item rest -> (if item.isMatch (item.+>) (else same)) (more rest))

whereNot list isMatch   := list.where (not: isMatch)

without list item       := list.where (nEQ item)

findWhere list isMatch  := list.more
  more l                := l same (item rest -> if (item.isMatch) item.ever (else more rest))

countWhere list isMatch := list.more
  more l                := l 0 (item rest -> ADD (item.isMatch 1 0) (more rest))

has list item           := list.hasSome (EQ item)

hasSome list isMatch    := list.more
  more l                := l False (item rest -> or (item.isMatch) (more rest))

hasNone list isMatch    := list.hasOnly (not: isMatch)

hasOnly list isMatch    := list False list.more
  more l                := if (not l.-<.isMatch) False (else or (!l.>-) (more l.>-))

sortedBy list isLess    := list.more
  more l                := l.-| (i r -> (lts gts -> concat lts (i.+> gts)) (more (filt same i r)) (more (filt not i r)))
  filt how item rest    := rest.where (comparee -> how (isLess comparee item))

bisectedBy list isMatch := { (extract where tagged) (extract whereNot tagged) }
  extract filt tups     := each (tups.filt sndOfTuplePair) fstOfTuplePair
  tagged                := list.more
  more l                := l.-| (item rest -> +> {item item.isMatch} (more rest) )

deduped list            := list.more
  more l                := l.-| (item rest -> item.+> (more (rest.without item)))

last list               := list.>- list.-< (_ _ -> last list.>-)

indexed list            := list.more 0
  more l i              := l.-| (item rest -> +> {i item} (more rest i.incr))

indices list            := list.more 0
  more l i              := l.-| (_ rest -> i.+> (more rest i.incr))

pairedUpWith list list2 := list.-| (item1 rest1 -> list2.-| (result item1 rest1))
  result i1 r1 i2 r2    := +> {i1 i2} (pairedUpWith r1 r2)

concat l1 l2            := l2 l1 (_ _ -> l1.prepend)
  prepend l             := l l2 (item rest -> item.+> (prepend rest))

forth list z f2         := list.more z
  more l accval         := l accval (item rest -> more rest (f2 accval item))

from list z f2          := list.more
  more l                := l z (item rest -> f2 item (more rest))

reverse list            := more [] list
  more cur l            := l cur (item rest -> more (item.+> cur) rest)

mGet m key              := m.more
  more l                := l same (item rest -> item (k v ->  if (k.EQ key) v.ever (else more rest)))

mDrop m key             := m.more
  more l                := l.-| (item rest -> item (k _ -> (if (k.EQ key) same (else item.+>))) (more rest))

mWith m key value       := +> {key value} (m.mDrop key)

mHasKeyWhere m isMatch  := m.hasSome (Tuple (k v -> k.isMatch))
mHasValWhere m isMatch  := m.hasSome (Tuple (k v -> v.isMatch))

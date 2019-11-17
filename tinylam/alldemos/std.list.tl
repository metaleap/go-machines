
-< list                 := list (ERR "forgot nil check before -<") True
>- list                 := list (ERR "forgot nil check before >-") False
+>                      := ListLink
-| list                 := list []
! list                  := list True (_ _ -> False)
// ? list                  := list False (_ _ -> True)

len list                := list 0 (_ rest -> incr (len rest))

at list i               := if (EQ 0 i) list.-< (else at list.>- i.decr)

each list fn            := list.more
  more l                := l.-| (item rest -> item.fn.+> (more rest))

slice list start length := take (list.drop start) length

drop list n             := if (EQ 0 n) list (else drop list.>- n.decr)

dropWhile list isMatch  := list.more
  more l                := l.-| (item rest -> if item.isMatch (more rest) (else l))

dropUntil list isMatch  := list.dropWhile isMatch.not·

take list n             := if (EQ 0 n) [] (else list.-<.+> (take list.>- n.decr))

takeUntil list isMatch  := list.more
  more l                := l.-| (item rest -> item.isMatch [] (item.+> (more rest)))

takeWhile list isMatch  := list.takeUntil isMatch.not·

where list isMatch      := list.more
  more l                := l.-| (item rest -> (item.isMatch (item.+>) same) (more rest))

whereNot list isMatch   := list.where (not· isMatch)

without list item       := list.where (nEQ item)

findWhere list isMatch  := list.more
  more l                := l None (item rest -> item.isMatch item.Ok (more rest))

countWhere list isMatch := list.more
  more l                := l 0 (item rest -> item.isMatch (ADD 1 (more rest)) (more rest))

has list item           := list.hasSome (EQ item)

hasSome list isMatch    := list.more
  more l                := l False (item rest -> or (item.isMatch) (more rest))

hasNone list isMatch    := list.hasOnly (not· isMatch)

hasOnly list isMatch    := list False more
  more item rest        := and item.isMatch (rest True more)

sortedBy list isLess    := list.sort
  sort l                := l.-| (item rest -> item.joined (sort (rest.below item)) (sort (rest.above item)))
  below                 := filt same
  above                 := filt not
  filt how rest item    := rest.where (comparee -> how (isLess comparee item))
  joined cur pref suff  := concat pref (cur.+> suff)

bisectedBy list isMatch := {extract where tagged, extract whereNot tagged}
  tagged                := list.more
    more l              := l.-| (item rest -> +> {item, item.isMatch} (more rest) )
  extract filt tups     := each (tups.filt sndOfTuplePair) fstOfTuplePair

deduped list            := list.more
  more l                := l.-| (item rest -> item.+> (more (rest.without item)))

last list               := list.>- list.-< (_ _ -> last list.>-)

indexed list            := list.more 0
  more l i              := l.-| (item rest -> +> {i, item} (more rest i.incr))

indices list            := list.more 0
  more l i              := l.-| (_ rest -> i.+> (more rest i.incr))

brokenUpOn list isMatch := list.more []
  more l acc            := l None (item rest -> item.isMatch (Ok {acc, rest}) (more rest (concat acc [item])))

pairedUpWith list list2 := list.-| (item1 rest1 -> list2.-| (result item1 rest1))
  result i1 r1 i2 r2    := +> {i1, i2} (pairedUpWith r1 r2)

intersperse list value  := list.-| (more False)
  more +? item rest    := value.+? (item.+> (rest [] (more +>)))

splitUpOn list isMatch  := list.-| (_ _ -> list.more)
  more l                := (l.brokenUpOn isMatch) (caseOk more) (else l.+> [])
  caseOk split          := _ (left right -> left.+> (split right))

joinedVia list sep      := concat.** (list.intersperse sep)

trimStart list items    := list.dropWhile items.has

trimEnd list items      := reverse (list.reverse.dropWhile items.has)

trim list items         := trimEnd (list.trimStart items) items

concat l1 l2            := l2 l1 (_ _ -> l1.prepend)
  prepend l             := l l2 (item rest -> item.+> (prepend rest))

from list z f2          := list.more
  more l                := l z (item rest -> f2 item (more rest))

forth list z f2         := list.more z
  more l accval         := l accval (item rest -> more rest (f2 accval item))

reverse list            := more [] list
  more acc l            := l acc (item rest -> more (item.+> acc) rest)

mGet m key              := m.more
  more l                := l None (item rest -> item (k v ->  if (k.EQ key) v.Ok (else more rest)))

mDrop m key             := m.more
  more l                := l.-| (item rest -> item (k v -> (if (k.EQ key) same (else item.+>))) (more rest))

mWith m key value       := +> {key, value} (m.mDrop key)

mHasKeyWhere m isMatch  := m.hasSome (Tuple (k v -> k.isMatch))
mHasValWhere m isMatch  := m.hasSome (Tuple (k v -> v.isMatch))

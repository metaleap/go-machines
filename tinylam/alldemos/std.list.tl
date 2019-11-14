
-<                      := _ [] True

>-                      := _ [] False

len list                := if (!list) 0 (else incr (len list.>-))

at list i               := if (EQ 0 i) list.-< (else at list.>- i.decr)

each list fn            := list.iter
  iter l                := if (!l) [] (else {l.-<.fn (iter l.>-)})

slice list start length := take (list.drop start) length

drop list n             := if (EQ 0 n) list (else drop list.>- n.decr)

take list n             := if (EQ 0 n) [] (else {list.-< (take list.>- n.decr)})

where list isMatch      := list.iter
  iter l                := if (!l) [] (else (if l.-<.isMatch {l.-<} (else self)) (iter l.>-))

whereNot list isMatch   := list.where (not: isMatch)

without list item       := list.where (nEQ _ item)

findWhere list isMatch  := list.iter
  iter l                := if (!l) self (else if l.-<.isMatch l.-<.ever (else iter l.>-))

countWhere list isMatch := list.iter
  iter l                := if (!l) 0 (else (ADD (l.-<.isMatch 1 0) (iter l.>-)))

has list item           := list.hasSome (EQ _ item)

hasSome list isMatch    := list.iter
  iter l                := if (!l) False (else or (l.-<.isMatch) (iter l.>-))

hasNone list isMatch    := list.hasOnly (not: isMatch)

hasOnly list isMatch    := and ?list list.iter
  iter l                := if (not l.-<.isMatch) False (else or (!l.>-) (iter l.>-))

sortedBy list isLess    := list.iter
  iter l                := if (!l) [] (else concat (iter (filt self l)) {l.-< (iter (filt not l))} )
  filt how l            := l.>-.where (item -> how (isLess item l.-<))

partitionedBy list test := { (each (elemsmarked.where >-) -<) (each (elemsmarked.whereNot >-) -<) }
  elemsmarked           := list.iter
  iter l                := if (!l) [] (else {{l.-< l.-<.test} (iter l.>-)})

deduped list            := list.iter
  iter l                := if (!l) [] (else (if (l.>-.has l.-<) self (else {l.-<})) (iter l.>-))

last list               := if (!list.>-) list.-< (else last list.>-)

indexed list            := pairUp list.indices list

indices list            := list.iter 0
  iter l i              := if (!l) [] (else {i (iter l.>- i.incr)})

pairUp l1 l2            := ifboth {{l1.-< l2.-<} (pairUp l1.>- l2.>-)}
  ifboth                := if (or !l1 !l2) []

concat l1 l2            := if (!l2) l1 (else l1.prepend)
  prepend l             := if (!l) l2 ( else {l.-< (prepend l.>-)} )

forth list z binOp      := list.acc z
  acc l cur             := if (!l) cur (else acc l.>- (binOp cur l.-<))

from list z binOp       := list.acc
  acc l                 := if (!l) z (else binOp l.-< (acc l.>-))

reverse list            := list.acc []
  acc l cur             := if (!l) cur (else acc l.>- {l.-< cur})

mGet m key              := m.iter
  iter l                := if (!l) self (else if (l.-<.-<.EQ key) l.-<.>-.ever (else iter l.>-))

mDrop m key             := m.iter
  iter l                := if (!l) [] (else ((l.-<.-<.EQ key) self {l.-<}) (iter l.>-))

mWith m key value       := {{key value} (mDrop m key)}

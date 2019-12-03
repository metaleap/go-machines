
-< list                   := list ? [] => ERR "forgot nil check before -<" same | .. => True
>- list                   := list ? [] => ERR "forgot nil check before >-" same | .. => False
+>                        := ListLink
! list                    := list ? [] => True | .. => _ _ -> False

len list                  := list ? [] => 0 | .. => _ rest -> incr (len rest)

at list i                 := if (i.EQ 0) list.-< (else at list.>- i.decr)

each list fn              := list.more
  more l                  := l ? .. => item rest -> item.fn.+> (more rest)

eachTry list fn           := list.more
  more l                  := l ? [] => Yay [] | .. => (item rest -> (r -> r Result? Nay => Nay | Yay => v -> (mr -> mr Result? Nay => Nay | Yay => vs -> Yay (v.+> vs)) (more rest)) item.fn)

slice list start length   := take (list.drop start) length

drop list n               := if (n.EQ 0) list (else drop list.>- n.decr)

dropWhile list isMatch    := list.more
  more l                  := l ? .. => item rest -> if item.isMatch (more rest) (else l)

dropUntil list isMatch    := list.dropWhile isMatch.not·

take list n               := if (n.EQ 0) [] (else list.-<.+> (take list.>- n.decr))

takeUntil list isMatch    := list.more
  more l                  := l ? .. => item rest -> item.isMatch [] (item.+> (more rest))

takeWhile list isMatch    := list.takeUntil isMatch.not·

where list isMatch        := list.more
  more l                  := l ? .. => item rest -> (item.+>.when item.isMatch (else same)) (more rest)

whereNot list isMatch     := list.where (not· isMatch)

without list item         := list.where (nEQ item)

findWhere list isMatch    := list.more
  more l                  := l ? [] => None | .. => item rest -> item.Have.when item.isMatch (else more rest)

countWhere list isMatch   := list.more
  more l                  := l ? [] => 0 | .. => item rest -> (check item.isMatch) (more rest)
  check ismatch           := if ismatch (ADD 1) (else same)

has list item             := list.hasSome (EQ item)

hasSome list isMatch      := list.more
  more l                  := l ? [] => False | .. => item rest -> or (item.isMatch) (more rest)

hasNone list isMatch      := list.hasOnly (not· isMatch)

hasOnly list isMatch      := list False more
  more item rest          := and item.isMatch (rest ? [] => True | .. => more)

sortedBy list isLess      := list.sort
  sort l                  := l ? .. => item rest -> concat (sort (rest.lt item)) (item.+> (sort (rest.geq item)))
  lt                      := filt same
  geq                     := filt not
  filt how rest item      := rest.where (comparee -> how (comparee.isLess item))

bisectedBy list isMatch   := {extract where tagged, extract whereNot tagged}
  tagged                  := list.more
    more l                := l ? .. => item rest -> +> {item, item.isMatch} (more rest)
  extract filt tups       := each (tups.filt sndOfTuplePair) fstOfTuplePair

deduped list              := list.more
  more l                  := l ? .. => item rest -> item.+> (more (rest.without item))

last list                 := list.>- ? [] => list.-< | .. => _ _ -> last list.>-

last2 list                := list ? .. => item rest -> (rest ? .. => item2 rest2 -> (rest2 ? [] => [item, item2] | .. =>  _ _ -> last2 rest))

sansLast list             := list ? .. => item rest -> (rest ? .. => _ _ -> item.+> (sansLast rest))

indexed list              := list.more 0
  more l i                := l ? .. => item rest -> +> {i, item} (more rest i.incr)

indices list              := list.more 0
  more l i                := l ? .. => _ rest -> i.+> (more rest i.incr)

brokenUpOn list isMatch   := list.more [] // caution with isMatch, takes 2 args of which the 2nd is the running item (the 1st being the accumulated predecessors)
  more l accum            := l ? [] => {accum, []} | .. => ok more accum
  ok cont acc item rest   := if (isMatch acc item) {acc, item.+> rest} (else rest.cont (concat acc [item]))

pairedUpWith list list2   := list ? .. => item1 rest1 -> list2 ? .. => result item1 rest1
  result i1 r1 i2 r2      := +> {i1, i2} (pairedUpWith r1 r2)

interspersed list value   := list ? .. => more False
  more useval item rest   := value.useval (item.+> (rest ? .. => more +>))

breakOnLast list isMatch  := list.lastwhere []
  lastwhere l acc         := l ? [] => None | .. => item rest -> pick item rest acc (lastwhere rest (concat acc [item]))
  pick item rest acc ret  := ret Could? Have => Have | None => if (item.isMatch) (Have {item, acc, rest}) (else None)

breakOnFirst list isMatch := list.firstwhere []
  firstwhere l acc        := l ? [] => None | .. => item rest -> (item.maybe acc rest) (else firstwhere rest (concat acc [item]))
  maybe item acc rest     := if (item.isMatch) (Have {item, acc, rest})

// eg. call:              rewritten   "foo[z,z,z]bar[x[z,z]y]baz"   (EQ ']')   (EQ '[')   (inner -> concat (+> '(' (inner.whereNot (EQ ','))) ")")
// result:                "foo(zzz)bar(x(zz)y)baz"
rewritten list d1 d2 re   := list.more
  more l                  := l.breakOnFirst d1 Could? None => l | Have => handleclosing more l
  handleclosing m l       := {_ pref suff} -> pref.breakOnLast d2 Could? None => l | Have => handleopening m suff
  handleopening m rest    := {_ pref suff} -> m (concat (concat pref suff.re) rest)

splitUpOn list isMatch    := list ? .. => _ _ -> list.more
  more l                  := (l.brokenUpOn (_ -> isMatch)) (andthen more)
  andthen cont fst snd    := fst.+> (snd ? .. => _ rest -> rest.cont)

joinedVia list sep        := list ? .. => _ _ -> concat.** (list.interspersed sep)

concat l1 l2              := l2 ? [] => l1 | .. => _ _ -> l1.prepend
  prepend l               := l ? [] => l2 | .. => item rest -> item.+> (prepend rest)

from list z f2            := list.more
  more l                  := l ? [] => z | .. => item rest -> f2 item (more rest)

forth list z f2           := list.more z
  more l accval           := l ? [] => accval | .. => item rest -> more rest (f2 accval item)

reversed list             := more [] list
  more acc l              := l ? [] => acc | .. => item rest -> more (item.+> acc) rest

toPairs list              := list ? .. => next
  next item rest          := rest ? .. => valitem more -> +> {item, valitem} (toPairs more)

mGet m key                := m.more
  more l                  := l ? [] => None | .. => {k v} rest -> v.Have.when (k.EQ key) (else more rest)

mDrop m key               := m.more
  more l                  := l ? .. => {k v} rest -> (unless (k.EQ key) (+> {k, v})) (more rest)
  unless                  := _ same

mWith m key value         := +> {key, value} (m.mDrop key)


step s start max f        := countup start
  countup i               := if (i.cmp max) [] (else +> i.f (countup (i.ADD s)))
  cmp                     := if (s.LT 0) lEQ (else gEQ)

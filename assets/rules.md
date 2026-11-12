# Diplomacy Orders and Their Resolution

This document is a detailed account of the rules for writing and resolving orders in
*Diplomacy*, based solely on the supplied rulebook. It is organized as an adjudication
specification: first determine whether each order is legal, then resolve all legal orders
simultaneously, then resolve retreats, and—after Fall—ownership and adjustments.

## 1. Fundamental principles

- Every Army and Fleet has a base strength of one. No unit is intrinsically stronger than
another.
- At most one unit may occupy a province at the end of a movement or retreat phase. There
are no exceptions.
- Movement orders are simultaneous. Do not resolve them sequentially or let an earlier
apparent result change a later order.
- An order may be legal yet fail because of other orders. Illegality and failure are
different:
  - an illegal, missing, or unusably ambiguous movement-phase order is treated as a Hold;
  - a legal order participates in adjudication even if it ultimately fails;
  - some failed attacks still cut support or create a standoff.
- A badly written order that is nevertheless legal and has only one possible meaning is
carried out. An ambiguous order that does not identify a single legal meaning is not
carried out.
- Unit type may be omitted when the unit's province identifies exactly one unit, since
only one unit can occupy a province.
- Orders are revealed simultaneously. Promises and agreements have no rules force.

Order sets are dated by season and year, beginning with Spring 1901, then Fall 1901,
Spring 1902, and so on. Each Spring has movement and, if needed, retreat phases. Each
Fall has movement and retreat phases followed by supply-center counting and Winter
adjustments. The usual notation identifies the unit and its current province, followed by
Hold, a destination, Support, or Convoy; notation is secondary to a single clear legal
meaning.

## 2. Board geography and legal reach

### 2.1 Province types

- Armies may occupy and move among land and coastal provinces, but never sea provinces.
- Fleets may occupy and move among sea and coastal provinces, but never inland provinces.
- An ordinary move goes only to an adjacent province that the moving unit could legally
occupy.
- No unit may occupy an unnamed space, Switzerland, or an island other than the named
provinces of the British Isles.

### 2.2 Fleet movement along coasts

A Fleet in a coastal province may move to another coastal province only when the two are
connected by a coastline. Land adjacency alone is insufficient. Thus fleet adjacency must
be evaluated by coast, not merely by shared provincial boundary.

### 2.3 Split-coast provinces

Bulgaria, Spain, and St. Petersburg each have two separately identified coasts.

- A Fleet entering one of these provinces enters a particular coast.
- If the Fleet could legally enter either coast, its move order must identify the
destination coast; otherwise the move fails as an invalid/ambiguous order.
- Once present, the Fleet occupies the entire province for occupancy, defense, attack,
and supply-center purposes, but its subsequent movement options are limited by the coast
it occupies.
- A Fleet may support an action into the province whenever it could move to at least one
of that province's coasts. A support order into the province does not need to name a
coast.
- An Army ignores coast distinctions.

### 2.4 Special waterways

- Kiel and Constantinople each count as a single-coast province. A Fleet may enter from
one side and leave by the other on a later turn. Armies also cross normally. A unit may
not skip over either province.
- An Army or Fleet may move directly between Denmark and Sweden. The common boundary does
not divide Sweden's coast. Denmark is not adjacent to Berlin.
- A Fleet moving between the Baltic Sea and Skagerrak must pass through Denmark or Sweden
on separate turns; those sea provinces are not directly adjacent.

## 3. Movement-phase order types

Each unit may receive exactly one of four kinds of order:

1. Hold.
2. Move.
3. Support.
4. Convoy (Fleet only).

A unit with no order, an illegal order, or an order that cannot be uniquely interpreted
Holds. A failed support or convoy order likewise leaves the unit in place, but the
distinction matters: it may still defend with strength one and may itself be attacked or
supported to hold.

## 4. Hold orders

A Hold orders a unit to remain in its province. A unit also effectively Holds when it
receives no order or an illegal order.

A unit that is Holding, Supporting, Convoying, or unordered may receive hold support. A
unit ordered to Move cannot receive hold support, even if its move fails.

## 5. Move orders and basic conflicts

### 5.1 Unopposed moves

A legal move succeeds if the destination can be occupied after all simultaneous conflicts
are resolved. A unit vacating a province does not block another unit from entering it,
provided the first unit's move succeeds and no other rule prevents the entry.

### 5.2 Equal-strength attacks: standoffs

If two or more units of equal effective strength attempt to enter the same province, none
enters. This is a standoff. All unsuccessful attackers remain in their origins unless
separately dislodged there.

A standoff does not dislodge a unit already in the contested province. If an occupant
remains there, it stays unless one attacker has strictly greater valid strength than
every relevant opposition and is allowed to dislodge it.

A province left empty because equal attacks bounced is still a standoff province. This
status matters during retreats.

### 5.3 A stationary or failed-to-move unit can block many moves

If an occupant does not successfully leave, it continues to occupy and defend its
province. This can cause an incoming move to fail, which can in turn cause another move
into the incoming unit's origin to fail. Resolution must therefore reach a consistent
simultaneous result rather than assuming that every ordered move vacates its origin.

### 5.4 Head-to-head exchanges

Two units ordered directly into one another's provinces may not exchange places by
ordinary land/sea moves. They engage head-to-head; compare their valid strengths against
one another. If neither has enough strength to dislodge the other, both remain. A convoy
creates the exception described below.

### 5.5 Rotations

Three or more units may rotate provinces in one turn, so long as no two are directly
exchanging places. A complete cycle such as A→B, B→C, C→A succeeds if no outside conflict
or other rule defeats one of its moves.

## 6. Support orders

### 6.1 What may support what

An Army or Fleet may support an Army or Fleet of any nation. Permission is unnecessary,
and support cannot be refused merely because it is unwanted.

The supporting unit must be able to make a legal move to the province in which it is
giving support. It need not be adjacent to the supported unit; it must be adjacent, with
the correct movement capability, to the supported unit's destination or holding province.

Consequences:

- an Army cannot support an action in a sea province;
- a Fleet cannot support an action in an inland province;
- coastal adjacency for Fleets is judged by navigable coastline;
- support into a split-coast province is support to the whole province, not a particular
coast.

### 6.2 Support to hold

An order naming only the supported unit's province supports that unit in place. It is
valid only if the supported unit is not ordered to Move. It may support a unit that
Holds, Supports, Convoys, has no order, or has a non-move order that fails.

If the supported unit was ordered to Move, hold support is invalid even when that move
fails. It does not convert into support for the attempted move.

### 6.3 Support to move

Support for a move must exactly match the supported unit's actual ordered origin and
destination. If the supported unit was ordered elsewhere, or was not ordered to Move, the
support does nothing. It does not become hold support.

### 6.4 Strength calculation

Each valid, uncut support adds one to the supported action's strength. A unit with two
valid supports therefore acts with strength three. To enter and dislodge an occupant, an
attack must have strictly greater applicable strength; equality produces a bounce.

Only support that is legal, matches the relevant order, is not cut, and is not
invalidated by dislodgement counts.

## 7. Cutting and invalidating support

### 7.1 Ordinary support cutting

An attack on a supporting unit cuts its support even if the attack does not succeed and
even if the attacker has no support. The key exception is the province against which
support is being given:

- an attack from any province other than the province where the support is directed cuts
support;
- an attack from the province where the support is directed does not cut support merely
by attacking; it must dislodge the supporter.

For support to hold, the province where support is directed is the supported unit's
province. For support to move, it is the supported move's destination.

### 7.2 Dislodgement of a supporter

A supporting unit that is dislodged provides no support. This includes dislodgement by an
attack from the province into which support was directed.

### 7.3 A dislodged attacker may still cut support elsewhere

A unit can cut support in a province other than the one from which it is dislodged. Being
dislodged does not retroactively erase its attack for support-cutting purposes. The
exception remains: a dislodged unit has no effect on the province of the attacker that
dislodged it.

### 7.4 Attacks by one's own nation do not cut support

A nation's attack on one of its own units does not cut that unit's support. A power
cannot use a nominal self-attack to cancel its own support.

### 7.5 Convoy-related precedence

Rare convoy interactions modify ordinary support cutting:

- A convoyed Army does **not** cut the support of a unit that is supporting an attack on
a Fleet required for that Army's convoy route. This prevents circular reasoning in which
the Army arrives only if it cuts the support, while it can cut the support only if it
arrives. The convoy is therefore broken if that supported attack dislodges the necessary
Fleet.
- If the Army has at least one other complete, successful convoy route, it **does** cut
support given by a unit in the Army's destination when that support attacks a Fleet on an
alternative convoy route. The independently successful route establishes the Army's
arrival and attack on the destination.

## 8. Dislodgement and the effect of dislodged units

### 8.1 Requirements for dislodgement

An occupant is dislodged only by an attack permitted to dislodge it whose valid strength
is strictly greater than the occupant's applicable defense and all competing attacks.

A unit ordered to Move does not receive hold support. When resolving a direct
head-to-head battle, its attack strength is relevant against the opposing move. In other
cases, a unit whose move fails remains physically present, but supports must match the
order actually given rather than being reassigned after failure.

### 8.2 A dislodged unit has no effect on its dislodger's origin

A dislodged unit—even if supported—has no effect on the province from which the
successful dislodging attack came. It cannot bounce the dislodger there, cut support
directed there merely through the dislodged attack, or prevent another unit from entering
that origin.

This creates an important edge case: if units from A and C both attack B, while B attacks
A, and the attack A→B dislodges B, B's attack on A has no effect. Another eligible move
into A can therefore succeed rather than bouncing with B.

### 8.3 A dislodged unit can still affect a different province

Dislodgement does not erase every aspect of the dislodged unit's order. It may:

- participate in a standoff in a province other than its dislodger's origin; and
- cut support in another province, subject to the ordinary support-cut rules.

These two rules must be applied together with the no-effect-on-the-dislodger rule.

## 9. No self-dislodgement

A nation may not dislodge one of its own units and may not provide support that would
enable another nation to dislodge one of its units. This remains true when the
self-dislodgement was unintended or caused by foreign support.

### 9.1 Own attack against own occupant

An attack cannot dislodge a friendly occupant. Friendly and foreign support cannot make
that self-dislodgement legal. The friendly occupant remains unless dislodged by a
separately sufficient foreign attack.

### 9.2 Supporting a foreign attack against one's own unit

Support supplied by a nation for an attack on that nation's own unit is denied for
purposes of dislodging that unit. If the foreign attacker has enough other valid
strength, it may still dislodge the unit.

### 9.3 Self-standoffs are permitted

A nation may order equal attacks by its own units into the same empty province to arrange
a self-standoff. Those attacks can also oppose foreign attacks. However, unexpected
support may make one friendly attack stronger and cause it to enter the empty province.
If the province is occupied by a friendly unit, the prohibition on self-dislodgement
still prevents entry through that unit.

An attempted self-dislodging attack can still be relevant as opposition to another
nation's attack; it is not simply erased from all comparisons.

## 10. Convoys

### 10.1 Basic requirements

Only an Army may be convoyed. Only a Fleet occupying a sea province may convoy; a Fleet
in any coastal province—including Kiel, Constantinople, or Denmark—may not convoy.

For a convoyed move to be available:

- the Army must be in a coastal origin province and be ordered to move to a coastal
destination province;
- every Fleet used in a route must be in a sea province;
- the Fleets in a route must form an unbroken chain of adjacent sea provinces from one
adjacent to the Army's origin to one adjacent to its destination;
- every Fleet in that route must have a matching Convoy order naming the Army's exact
origin and destination;
- a Fleet may convoy no more than one Army in a turn.

The convoy can span any number of contiguous sea provinces and may use Fleets belonging
to different nations. Convoying provides movement, not strength; support must be ordered
separately.

### 10.2 Mismatched convoy orders

If a Convoy order names a different destination from the Army's Move order, that Fleet
does not create the stated route. The Army remains at its origin unless some other valid
method or complete matching convoy route carries it to its ordered destination.

### 10.3 Disruption

An attack on a convoying Fleet has no effect unless it dislodges the Fleet. If any Fleet
required by a particular convoy route is dislodged, that route fails.

If every possible matching route fails, the Army remains in its origin. It may
nevertheless be dislodged there by a successful attack.

### 10.4 Multiple convoy routes

Matching orders may establish several routes for the same Army. The Army reaches its
destination if at least one complete route remains intact. Dislodgement of a Fleet on one
route does not stop a different intact route.

### 10.5 Bounce at the destination

The convoyed Army attacks its destination normally and may receive matching move support
there. If it ties another applicable force and cannot enter, it remains at its original
coastal province; it does not stop at sea or at an intermediate point. It can still be
dislodged from its origin by a separate successful attack.

### 10.6 Exchanging places by convoy

Two units may exchange provinces when one or both moves are by convoy. This is the
exception to the ordinary prohibition on two-unit exchanges.

### 10.7 Land route and convoy route both available

An Army's destination may be adjacent by land while matching convoy orders also exist.

- If at least one convoying Fleet belongs to the Army's own nation, the convoy route must
be used and the land route is discarded.
- If no convoying Fleet belongs to the Army's nation, the land route is used by default.
This prevents foreign powers from kidnapping the Army through an unwanted convoy.
- In the latter case, the Army's owner may explicitly order the move “via convoy,” in
which case the convoy is used.

These distinctions affect whether a head-to-head exchange is permitted and whether
disruption of the convoy causes the move to fail.

## 11. A practical simultaneous adjudication order

The rulebook resolves everything simultaneously, but the following dependency order helps
produce the correct fixed result:

1. **Normalize orders.** Identify units and provinces; convert missing, illegal, and
irreducibly ambiguous orders to Holds.
2. **Validate geometry and matching.** Check move adjacency, unit terrain, fleet coasts,
support reach, exact support matching, convoying Fleet location, and exact convoy
matching.
3. **Enumerate convoy routes.** Record every complete matching route for each potentially
convoyed Army and whether a land route is selected under the land-versus-convoy rules.
4. **Identify attacks on supporters and convoying Fleets.** Apply the exception for
attacks from the province where support is directed, the self-attack rule, and the convoy
paradox precedence rules.
5. **Determine which supporters and convoying Fleets are dislodged.** This may remove
support or break routes; re-evaluate dependent moves until the result is internally
consistent under the stated exceptions.
6. **Compute effective strengths.** Count only legal, matching, uncut support, excluding
support forbidden by self-dislodgement.
7. **Resolve destination conflicts, head-to-head battles, rotations, and occupancy
dependencies.** Require strict superiority; ties bounce. A unit whose departure fails
continues to occupy its origin.
8. **Mark dislodgements and standoff provinces.** Preserve the special effects of
dislodged units in provinces other than their dislodger's origin.
9. **Apply all successful moves simultaneously.** No province may end with more than one
unit.
10. **Resolve retreats simultaneously.** Use the post-movement board plus the recorded
attackers and standoff provinces.

This is a dependency order, not a sequence in which units move one at a time. Any
implementation must arrive at the single result required by all simultaneous orders and
the convoy precedence rules.

## 12. Retreat and disband orders

### 12.1 When retreat orders are written

After movement adjudication, every dislodged unit must retreat or disband. Retreat orders
are written without negotiation and revealed simultaneously.

### 12.2 Legal retreat destinations

A retreat must be to an adjacent province the unit could normally enter. A retreat may go
deeper into enemy territory, but it may not go:

- to any occupied province, using the board position after successful movement;
- to the province from which the attack that dislodged the unit came; or
- to an empty province in which a standoff occurred during the just-completed movement
phase.

A retreat cannot be supported and cannot use a convoy. Fleet coast restrictions still
apply; a Fleet must retreat to a legally reachable coast.

The prohibition is specifically the dislodger's origin, not the origin of every unit that
supported the attack. A province vacated for reasons other than a standoff is not
forbidden merely because it became empty.

### 12.3 Retreat conflicts

If two or more dislodged units retreat to the same province, all such units are
disbanded. Retreat strength and nationality do not matter; there is no support or
tie-breaker.

### 12.4 Mandatory and voluntary disbandment

A dislodged unit is disbanded when:

- it has no legal retreat destination;
- it receives no retreat order;
- its retreat order is illegal or conflicts with another retreat into the same province;
or
- its owner voluntarily orders it disbanded instead of retreating.

## 13. Fall supply-center ownership

Supply-center control is evaluated only after the Fall turn has been completed, including
Fall retreats.

- A nation gains control of a supply center when one of its units occupies that province
at the end of Fall.
- Once gained, control persists while the center is empty or occupied by that nation's
unit at the end of subsequent Falls.
- Merely passing through a center in Spring, or occupying it in Spring and leaving it in
Fall, does not change control.
- An occupation at the end of Fall changes control even if the center previously belonged
to another nation.

## 14. Winter adjustments

After Fall retreats and the ownership count, each nation's unit count must be reconciled
with its number of controlled supply centers. Adjustment orders are written and revealed
simultaneously, without negotiation.

### 14.1 Forced removals

If a nation has more units than controlled centers, its owner chooses enough units to
disband to make the counts equal.

For a nation in civil disorder, forced removals are automatic: remove units farthest from
the home nation first; among equally distant units, remove Fleets before Armies, then use
alphabetical order of their provinces.

### 14.2 Builds

If a nation has fewer units than controlled centers, it may build no more than the
difference, subject to all of these limits:

- a build must be in an unoccupied home supply center belonging to that nation at the
start of the game;
- the nation must currently control that home center;
- at most one unit may be built in a province;
- an Army may be built in an inland or coastal home center;
- a Fleet may be built only in a coastal home center;
- a Fleet build in St. Petersburg must specify north or south coast;
- no build may occur in a captured non-home center;
- an occupied home center cannot be used, regardless of who occupies it;
- a nation may waive any or all builds.

A nation with units but no controlled home center may continue playing, but cannot build
until it again controls an unoccupied home center at the relevant Fall adjustment.

Confusing or invalid build/removal orders are ignored. A waived or unusable build is lost
for that adjustment phase; it is not banked.

## 15. Civil disorder and neutral units in reduced-player variants

If a player supplies no Spring or Fall orders because the nation has collapsed into civil
disorder:

- all its units Hold;
- they do not support one another;
- they may be supported to hold by other nations;
- if dislodged, they are disbanded rather than retreated;
- the nation receives no new builds;
- forced adjustment removals use the automatic distance rule above.

In the specified reduced-player variants, removed nations' starting units similarly
remain in place and defend without supporting one another; other players may support them
to hold, and any such unit that is dislodged is disbanded.

## 16. Edge-case checklist

Before finalizing an adjudication, verify all of the following:

- Did every Fleet move, support, retreat, or build respect coastline connectivity and
split-coast identity?
- Was every support geometrically possible for the supporter, and did move support
exactly match the supported order?
- Was hold support excluded from every unit ordered to Move, even when that move failed?
- Was support cut by any attack from outside the supported province, even an unsuccessful
or subsequently dislodged attacker?
- Was support preserved against a mere attack from the province into which support was
directed, but removed if that attacker dislodged the supporter?
- Were attacks by a power on its own units prevented from cutting its own support?
- Was support that would enable dislodgement of the supporter's own unit excluded?
- Did every dislodgement require strict superiority rather than equality?
- Was a dislodged unit given no effect on the province from which its dislodger attacked,
while retaining permitted effects elsewhere?
- Were direct two-unit swaps blocked unless at least one unit used a convoy?
- Were longer rotations allowed?
- Did failed departures continue to block their origins and possibly create cascading
failures?
- Were all matching convoy routes considered, rather than treating one broken route as
breaking every route?
- Was a convoy broken only by dislodgement of a required Fleet, not by a mere attack?
- If both a land path and convoy existed, was the correct route selected based on
convoying Fleet ownership or an explicit “via convoy” instruction?
- Were the two convoy/support precedence rules applied to avoid circular convoy
adjudication?
- Were retreat destinations tested against final occupancy, the dislodger's origin, and
every movement-phase standoff province?
- Were colliding retreats all disbanded?
- Was supply-center control changed only after Fall, after retreats?
- Were builds restricted to controlled, vacant home supply centers with legal unit and
coast types?


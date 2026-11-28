# Diplomacy order resolution test plan

Each case describes one resolution of a complete turn. Its numeric ID is permanent and doubles as its fixture basename. The expected results come from [the adjudication rules](rules.md) (§ numbers below) and [the printed rulebook](diplomacy-instructions.pdf) (printed page numbers). The following house rules take precedence.

- **H1.** Submit conditional retreats with the other orders. Resolve movement first. Then process retreats for the units dislodged during movement, using the resulting board. To voluntarily disband a dislodged unit, omit its retreat order.
- **H2.** Process movement commands in this order: convoy, hold, move, support. Keep each unit's first valid command in that processing order. These four order types share one choice per unit. During retreats, independently keep each dislodged unit's first valid retreat. In each stage, skip invalid orders and ignore later orders for that unit once a valid one has been accepted.
- **H3.** Ignore adjustment orders in Spring. In Fall, apply each power's valid builds or removals in submission order until its quota is filled. If a power still has more units than controlled supply centers, randomly disband its remaining units until the counts match. This also applies to countries in civil disorder.

- **H4.** Orders must explicitly specify the unit type (army or fleet).
- **H5.** A fleet moving into a province with multiple coasts must explicitly specify its destination coast, even when only one coast is reachable.

### Coverage matrix

Use the links to open the complete case. Section 11 supplies the simultaneous resolution procedure and §16 cross-checks the interactions below.

| Rule | Behavior | Cases |
|---|---|---|
| H2 · movement order selection | First valid order in convoy, hold, move, support processing order, including an invalid first instruction and a valid move that bounces | [106](#case-106), [107](#case-107), [108](#case-108), [109](#case-109), [129](#case-129), [130](#case-130), [131](#case-131), [431](#case-431), [432](#case-432) |
| H1–H2 · conditional orders | Independent movement and retreat choices in either submission order, invalid first refuge, duplicate retreats and collision after acceptance | [501](#case-501), [502](#case-502), [503](#case-503), [504](#case-504), [505](#case-505), [506](#case-506), [523](#case-523), [524](#case-524), [525](#case-525), [526](#case-526) |
| H3 · adjustments by season | Identical surplus or deficit in Spring and Fall, movement plus an adjustment for the same unit | [608](#case-608), [609](#case-609), [610](#case-610), [611](#case-611), [636](#case-636) |
| H4, §§1, 3–4 · interpret orders | Explicit/implicit hold, required unit type, ambiguity, invalid issuer, nonexistent province or unit, illegal order versus legal failure | [101](#case-101), [102](#case-102), [103](#case-103), [104](#case-104), [105](#case-105), [110](#case-110), [114](#case-114), [131](#case-131), [132](#case-132), [303](#case-303), [304](#case-304), [305](#case-305), [328](#case-328) |
| §§2.1–2.2 · access | Army/fleet terrain, direct distance, coastline, and conflicting map edges | [103](#case-103), [110](#case-110), [111](#case-111), [112](#case-112), [113](#case-113), [126](#case-126), [132](#case-132) |
| H5, §2.3 · split coasts | Require an explicit entry coast, restrict exit, contest and support a whole province, ignore coasts for armies, respect coast on retreat and build | [114](#case-114), [115](#case-115), [116](#case-116), [117](#case-117), [118](#case-118), [119](#case-119), [120](#case-120), [121](#case-121), [311](#case-311), [514](#case-514), [515](#case-515), [522](#case-522), [617](#case-617), [618](#case-618), [633](#case-633) |
| §2.4 · waterways | Canal entry and exit, army passage, Denmark–Sweden, prohibited shortcuts | [122](#case-122), [123](#case-123), [124](#case-124), [125](#case-125), [127](#case-127), [128](#case-128) |
| §§1, 5.1–5.3, 11 · simultaneous occupancy | Equal forces and stationary defenders, multiple attackers, successful and cascading failed departures, input permutation | [201](#case-201), [202](#case-202), [203](#case-203), [204](#case-204), [205](#case-205), [209](#case-209), [210](#case-210), [312](#case-312), [313](#case-313) |
| §§5.4–5.5, 8 · exchanges, cycles and dislodgement | Direct swap, stronger opponent, rotation and broken rotation, suppressed attack at dislodger's origin, standoff elsewhere | [206](#case-206), [207](#case-207), [208](#case-208), [209](#case-209), [210](#case-210), [211](#case-211), [212](#case-212), [419](#case-419), [420](#case-420), [421](#case-421) |
| §§4, 6.1–6.2 · support reach and defense | Foreign support, destination reach by terrain and coast, hold support for holding, unordered, supporting and convoying units, failed moves and failed non-moves | [301](#case-301), [302](#case-302), [303](#case-303), [304](#case-304), [308](#case-308), [309](#case-309), [310](#case-310), [311](#case-311), [320](#case-320), [327](#case-327), [328](#case-328), [417](#case-417) |
| §§6.3–6.4, 8.1 · offensive strength | Match actual move and origin, no conversion to hold support, equal and superior supported attacks | [301](#case-301), [302](#case-302), [303](#case-303), [304](#case-304), [305](#case-305), [306](#case-306), [307](#case-307), [312](#case-312), [313](#case-313), [411](#case-411) |
| §§7.1–7.4, 8.3 · support cuts | Outside attack, own attack, destination exception and dislodgement, dislodged cutter, hold support cut, illegal attack | [314](#case-314), [315](#case-315), [316](#case-316), [317](#case-317), [318](#case-318), [319](#case-319), [320](#case-320), [330](#case-330) |
| §9 · self-dislodgement | Own occupant, denied foreign assistance, sufficient independent strength, self-standoff, involuntary support, opposition to a rival, friendly departure | [321](#case-321), [322](#case-322), [323](#case-323), [324](#case-324), [325](#case-325), [326](#case-326), [329](#case-329) |
| §§6.1, 10.1–10.2 · available convoys | Fleet and passenger type, coast endpoints, complete matching orders, multisea chain and missing link, one army per fleet, support cannot be carried | [401](#case-401), [402](#case-402), [403](#case-403), [404](#case-404), [405](#case-405), [406](#case-406), [407](#case-407), [408](#case-408), [409](#case-409), [418](#case-418), [427](#case-427), [428](#case-428), [429](#case-429) |
| §§10.3–10.5 · convoy strength and disruption | Landing bounce, separate landing support, attacked versus dislodged carrier, alternate surviving route or all routes cut, attack on passenger at origin | [410](#case-410), [411](#case-411), [412](#case-412), [413](#case-413), [414](#case-414), [415](#case-415), [416](#case-416), [417](#case-417), [430](#case-430) |
| §§10.6–10.7 · route selection | Swap by one or two convoys, own versus foreign fleet, explicit convoy and disruption with parallel land route | [419](#case-419), [420](#case-420), [421](#case-421), [422](#case-422), [423](#case-423) |
| §7.5 · convoy/support precedence | Necessary fleet under supported attack, independent route enables cut, landing with added support | [424](#case-424), [425](#case-425), [426](#case-426) |
| §§12.1–12.2, 12.4 · legal retreats and removals | Conditional trigger, disband by omitting a retreat, no legal destination, origin, terrain/coast, no convoy or support, end-of-movement occupancy/standoff | [501](#case-501), [502](#case-502), [503](#case-503), [504](#case-504), [505](#case-505), [506](#case-506), [507](#case-507), [508](#case-508), [509](#case-509), [510](#case-510), [511](#case-511), [512](#case-512), [513](#case-513), [514](#case-514), [515](#case-515), [516](#case-516), [517](#case-517), [521](#case-521), [524](#case-524) |
| §12.3 · simultaneous retreats | Different refuges, two-way and three-way collisions, support cannot break a tie, different coasts of one province | [518](#case-518), [519](#case-519), [520](#case-520), [521](#case-521), [522](#case-522) |
| §13 · center control | Spring visit, Fall capture/recapture, persistence in vacancy, neutral center, ownership after retreats | [601](#case-601), [602](#case-602), [603](#case-603), [604](#case-604), [605](#case-605), [606](#case-606), [607](#case-607), [622](#case-622), [624](#case-624) |
| H3, §§13–14 · count and remove | Post-retreat census, ordered removals up to quota, invalid requests, random removals after no choice or a partial choice | [606](#case-606), [607](#case-607), [610](#case-610), [611](#case-611), [612](#case-612), [613](#case-613), [630](#case-630), [631](#case-631), [632](#case-632), [637](#case-637), [638](#case-638), [639](#case-639) |
| §14.2 · build eligibility | Home and current control, home recovered, site vacant after movement, legal unit and coast, one unit per province, over-quota and ambiguous orders | [609](#case-609), [614](#case-614), [615](#case-615), [616](#case-616), [617](#case-617), [618](#case-618), [619](#case-619), [620](#case-620), [621](#case-621), [622](#case-622), [623](#case-623), [624](#case-624), [625](#case-625), [626](#case-626), [633](#case-633), [634](#case-634), [635](#case-635) |
| §14.2 · waiver and carryover | Waive one or all, no stored builds in a later year | [627](#case-627), [628](#case-628), [629](#case-629) |
| H3, §§14.1, 15 · civil disorder | Hold without mutual support, outside support, disband on defeat, no builds, random removal of excess units | [701](#case-701), [702](#case-702), [703](#case-703), [704](#case-704), [705](#case-705), [706](#case-706), [707](#case-707) |
| §15, PDF p. 19 · variants | Eliminated powers defend, outside support, national adjustment counts for a shared player, Italian neutrality and alliance timing | [708](#case-708), [709](#case-709), [710](#case-710), [711](#case-711), [712](#case-712), [713](#case-713) |

### Building a fixture

The future runner lives in `/tests` at the repository root. Give each case one file in `/tests/maps/<test-number>.json`, e.g. `/tests/maps/115.json`. Use exactly this envelope.

```json
{
  "name": "",
  "description": "",
  "expectedOutcome": "",
  "inputBoard": [],
  "outputBoard": [],
  "orders": []
}
```

Use the case title for `name`. Put the purpose, rule reference and relevant setup in `description`. Expand the result into plain language in `expectedOutcome`. Both boards contain **all 75 provinces** from [the original map](original-diplomacy.json), in source order, including geometry, coordinates, coasts and supply-center home designations.

1. Deep-copy the map's province list. Clear all units and set every existing center's `controlledBy` to null. Install exactly the units and control listed in the case. This gives sparse positions with explicit counts for each power. Keep non-centers as non-centers. A center may start under different control from its occupant.
2. In the tables, `FR A par pic / GE F hol` puts two French armies in Paris and Picardy and a German fleet in Holland. The power and type apply to each following location until changed. Power codes AU, EN, FR, GE, IT, RU, TR stand for Austria-Hungary, England, France, Germany, Italy, Russia and Turkey. Store the map's full power names. Province codes are map IDs. A suffix such as `spa-nc` gives the fleet's coast. Use a null unit coast elsewhere. Set each unit's `location.id` to its province and give it a distinct stable ID, e.g. `u-par` for the unit starting in Paris. `SC FR par bre` means France controls exactly those listed centers in addition to any others explicitly assigned in the row.
3. Read each orders cell **left to right**. Commas separate successive commands. Each origin identifies the issuing unit, type, power and starting coast from that row. Encode the ordered commands as `models.Order` entries with `playerName` and `value` (use deterministic IDs and increasing creation times). Consecutive commands from the same power may share one comma-delimited value. Preserve the original sequence across powers and command types. Use the raw text indicated for input-format tests. An empty list is written `[]`. An unlisted unit has no movement order and defends in place.

   | Plan notation | Example command value |
   |---|---|
   | `par H`, `par>bur` | `A par - H`, `A par - bur` |
   | `pic S par>bur`, `pic S par H` | `A pic S A par - bur`, `A pic S A par - H` |
   | `eng C lon>bel` | `F eng C A lon - bel` |
   | `bur R gas` | `A R bur - gas` |
   | `FR +F bre`, `FR -par` | French build `F bre`, French adjustment removal `D A par` |

4. Deep-copy the input to construct the expected board. In results, `par>bur` places the **original Paris unit** in Burgundy and updates its location. `×bur` means the original Burgundy unit was dislodged and disbanded. `bur R gas` means it was dislodged and survived in Gascony. Builds and adjustment removals are identified separately. Apply these changes simultaneously. **Every unmentioned unit stays at its listed origin, with its original ID, power, type and coast. All unmentioned control persists.** “Same” means the entire final board equals the input. These conventions specify both boards even in rows giving only the differences.
5. Create a fresh `models.Game` for each fixture. Load its board and players, then set `turn.TurnNumber` for the case's season when assigning `turn.Orders` from the fixture. Pass that turn directly to `board.ApplyTurn(game.Board, turn)`, the same path used when the turn handler closes a turn. Use the returned board. Fail on error or panic. Compare province and unit data to the independent expected board by province ID, checking for missing, duplicate or extra provinces before comparison. Use the count assertion below for units eligible for random removal. Require unique unit IDs, correct nested locations and at most one unit per province. Keep the input for comparison independent of the mutable game.

### Runner setup and limits

The six fixture fields have no season, player-state or per-command result field. Use the case number to select the following runner setup and assertions. Codes flag particularly important limitations. Wherever a row specifies a command disposition or intermediate event, check it using D or T in addition to the board comparison.

- **Season.** Use the season and year in the table heading unless the row specifies its own. Odd turn numbers are Spring and even turn numbers are Fall. For example, the Spring build in 608 uses `TurnNumber = 1` and the Fall build in 609 uses `TurnNumber = 2`. Spring and Fall 1902 use 3 and 4. Case 629 also exercises a prior waiver, so its setup resolves the earlier turns described in that row.
- **Players.** Assume normal seven-power play and participating powers submitting orders (an intentionally empty order sheet is allowed, PDF p. 7). A missing order for one unit alone does not establish civil disorder. Explicit civil disorder and variants use **V** below.
- **D — order disposition.** Record the original command ordinal, issuer, text, and the validation/adjudication reason. Distinguish accepted but unsuccessful orders, invalid orders, duplicates, and inapplicable conditional or seasonal orders where specified. For example, a first hold followed by a move can leave exactly the same board as two rejected orders. Inspect the returned validation results and, for unrecognized text, parser results. Where the path drops a command, groups by kind, or cannot identify identical occurrences, a per-occurrence outcome/trace is needed. Never infer a reason from the final board alone.
- **T — transient outcomes.** Assert the indicated movement-stage event as well as the final board. Standoff provinces, dislodger origins, cuts, route viability and retreat collisions are transient. Where no public result exposes one, the future runner needs an adjudication observation or a board-distinguishing companion case. `×` always implies a movement dislodgement followed by a disband under §12.4, including when legal refuges exist but no retreat was submitted.
- **B — generated build IDs.** Compare power, type, province and coast exactly. Bind an expected placeholder such as `new-bre` to the fresh ID actually produced for that build, checking non-emptiness, uniqueness and freshness against all input IDs. Preserve exact IDs for existing units. Alternatively inject a deterministic ID source.
- **R — random adjustment removals (631, 639, 704–707).** Each of these fixtures has one surviving unit for the affected power. In `outputBoard`, use the first eligible unit listed in the starting position as the example survivor. The runner accepts any eligible survivor and checks that the power's final unit count equals its number of controlled centers. At each province eligible for the draw, accept either an empty province or its original unit with its original data. Compare ordered removals, other units, province data and supply-center control exactly.
- **F — order language.** Feed each indicated literal to the production parser. The current grammar requires a unit type and lacks the explicit “via convoy” modifier used in 422.
- **V — game context.** Seed civil disorder, a reduced-player assignment, neutrality, or history as stated. The board and turn alone cannot carry all of this through the current `ApplyTurn` signature. A game-level entry point or explicit resolution context and, for history, a prior-turn assertion are required for those claims.

**When to hold or disband.** A unit holds by default once its movement orders have been checked and none is valid (102, 108). A legal move that bounces still uses the unit's movement order (131). For retreats, check each destination against post-movement occupancy, the dislodger's origin and movement standoffs (§12.2). Once each dislodged unit's first legal retreat has been selected, resolve those retreats simultaneously. Disband a unit with no valid retreat or one whose selected retreat collides with another (§§12.3–12.4, PDF p. 18).

**Board evidence.** The setup tables use reciprocal routes and check fleet moves against the relevant coasts, with deliberately illegal orders and case 126 identified explicitly. The map includes asymmetric data. For example, `par.routes` lists `mar`, while `mar.routes` omits `par` and the PDF p. 7 enumerates Paris's four actual neighbors. The PDF also lists English Channel–North Sea, absent on both JSON sides. Warsaw references a nonexistent `svn` where Livonia lists `war`. Preserve the fixture geometry as supplied. Use rulebook outcomes where explicit and identify any map correction as a separate dependency.

### Orders, terrain and coasts — Spring 1901

| ID and purpose | Rule | Starting units | Orders in submission order | Expected result |
|---|---|---|---|---|
| <a id="case-101"></a>101 Explicit hold [D] | §§1,4 | FR A par | par H | Same. Hold accepted. |
| <a id="case-102"></a>102 Missing instruction | §§1,4 | FR A par | [] | Same. Paris holds. |
| <a id="case-103"></a>103 Unopposed army | §§2.1,5.1 | FR A par | par>bur | par>bur. |
| <a id="case-104"></a>104 Missing unit type [F,D] | H4 | FR A par | FR literal `par - bur` | Same. Reject the order for omitting the unit type. The army holds in Paris. |
| <a id="case-105"></a>105 Ambiguous destination [F,D] | §1 | FR A par | FR literal `A par - pic or bur` | Same. Two legal readings, neither executes. |
| <a id="case-106"></a>106 Hold before move [D] | H2 | FR A par | par H, par>bur | Same. Hold accepted, move ignored. |
| <a id="case-107"></a>107 Move before hold [D] | H2 | FR A par | par>bur, par H | par H. Hold processed before move; move ignored. |
| <a id="case-108"></a>108 Invalid first, valid second [D] | H2, §§1,3 | FR A par | par>lon, par>bur | par>bur. Reject the illegal Paris–London move, then accept Paris–Burgundy as the army's first valid order. |
| <a id="case-109"></a>109 First of two moves [D] | H2 | FR A par | par>pic, par>bur | par>pic only. |
| <a id="case-110"></a>110 Unknown space and foreign issuer [D] | §§1,2.1, PDF pp. 6–8 | FR A bur / GE A mun | FR literal `A bur - swi`, FR literal `A mun - ruh` | Same. Switzerland has no occupiable ID. France cannot order Germany's unit. |
| <a id="case-111"></a>111 Unit terrain and distance [D] | §§2.1,3 | EN A lon / FR F bre / GE A mun | lon>nth, bre>par, mun>bel | Same. Army to sea, fleet inland, and nonadjacent inland move each fail. |
| <a id="case-112"></a>112 Shared land boundary | §2.2, PDF p. 8 | IT F rom | rom>apu | Same. No continuous coast. |
| <a id="case-113"></a>113 Navigable coastline | §2.2 | IT F rom | rom>nap | rom>nap. |
| <a id="case-114"></a>114 Two reachable coasts, none specified [D] | §2.3 | TR F con | con>bul | Same. Ambiguous east/south coast. |
| <a id="case-115"></a>115 Chosen east coast | §2.3 | TR F con | con>bul-ec | con>bul-ec. |
| <a id="case-116"></a>116 Missing coast with one reachable coast [D] | H5 | FR F gas | gas>spa | Same. Reject the order for omitting the destination coast. The fleet holds in Gascony even though only Spain’s north coast is reachable. |
| <a id="case-117"></a>117 Wrong destination coast | §2.3 | FR F gas | gas>spa-sc | Same. |
| <a id="case-118"></a>118 Restricted exit | §2.3 | FR F spa-nc | spa-nc>wes | Same. |
| <a id="case-119"></a>119 Reachable exit | §2.3 | FR F spa-sc | spa-sc>wes | spa-sc>wes. Clear the coast at sea. |
| <a id="case-120"></a>120 One province across coasts [T] | §§2.3,5.2 | FR F gas / IT F gol | gas>spa-nc, gol>spa-sc | Same. Spain is one empty standoff province. |
| <a id="case-121"></a>121 Army ignores coasts | §2.3 | AU A ser | ser>bul | ser>bul, null coast. |
| <a id="case-122"></a>122 Kiel and Constantinople entrances | §2.4 | GE F hol / TR F bla | hol>kie, bla>con | hol>kie, bla>con. Both destinations have a single coast. |
| <a id="case-123"></a>123 Canal exits on later turn | §2.4 | GE F kie / TR F con | kie>ber, con>aeg | kie>ber, con>aeg. No remembered entry side needed. |
| <a id="case-124"></a>124 No shortcuts through waterways | §2.4 | GE F hol / TR F bla / RU F bal / EN A den | hol>ber, bla>aeg, bal>ska, den>ber | Same. Each attempted direct edge is illegal. |
| <a id="case-125"></a>125 Denmark–Sweden, either unit type | §2.4 | GE A den / RU F bal | den>swe, bal>den | den>swe, bal>den. The fleet follows into vacated Denmark. |
| <a id="case-126"></a>126 Rulebook/map regression | §§2.1–2.2, PDF p. 7 | EN F eng / FR A par | eng>nth, par>mar | eng>nth, Paris holds. The rulebook enumerates both origins' neighbors. JSON alone gives an incorrect oracle. |
| <a id="case-127"></a>127 Single Swedish coastline | §2.4 | RU F den / GE F bal | den>swe, bal>den | den>swe, bal>den. |
| <a id="case-128"></a>128 Armies cross the canal provinces | §2.4 | TR A con / GE A kie | con>smy, kie>ber | con>smy, kie>ber. |
| <a id="case-129"></a>129 Move takes precedence over earlier support [D] | H2, §§5.1,6.4 | FR A pic bur / GE A bel | bur S pic>bel, pic>bel, bur>gas | bur>gas. Moves are processed before supports, so Burgundy's earlier support is ignored. Picardy bounces, Belgium stays. |
| <a id="case-130"></a>130 Move submitted before support [D] | H2, §§5.1,6.4 | FR A pic bur / GE A bel | bur>gas, pic>bel, bur S pic>bel | bur>gas. Picardy bounces, Belgium stays. The later support is ignored. |
| <a id="case-131"></a>131 Legal but unsuccessful first move [D] | H2, §5.2 | FR A pic / GE A bel | pic>bel, pic>par | Same. The first instruction attacks and bounces. The later legal move cannot replace it after resolution. |
| <a id="case-132"></a>132 Unnamed island and unknown origin [F,D] | §§1,2.1, PDF p. 8 | IT F tyn | IT literal `F tyn - sic`, IT literal `A xyz - rom` | Same. Sicily has no province entry. The second command identifies no unit and cannot create one. |

### Contested movement — Spring 1901

| ID and purpose | Rule | Starting units | Orders in submission order | Expected result |
|---|---|---|---|---|
| <a id="case-201"></a>201 Equal attacks on vacancy [T] | §§1,5.2 | FR A pic / GE A ruh | pic>bel, ruh>bel | Same. Belgium empty by standoff. |
| <a id="case-202"></a>202 Equal attacks on occupant | §5.2 | FR A pic / GE A ruh / EN F bel | pic>bel, ruh>bel | Same. England keeps its fleet. |
| <a id="case-203"></a>203 Three-way tie [T] | §5.2 | FR A pic / GE A ruh / EN F nth | pic>bel, ruh>bel, nth>bel | Same. Belgium empty. Army and fleet have equal base strength. |
| <a id="case-204"></a>204 Follow a departure | §§5.1,11 | FR A par bur / GE A mun | par>bur, bur>mun, mun>boh | par>bur, bur>mun, mun>boh. |
| <a id="case-205"></a>205 Block a chain at its end | §§5.3,11 | FR A par bur / GE A mun boh | par>bur, bur>mun, mun>boh | Same. Bohemia holds, so neither Munich nor Burgundy vacates. |
| <a id="case-206"></a>206 Ordinary exchange | §5.4 | FR A par / GE A bur | par>bur, bur>par | Same. Both moves fail. |
| <a id="case-207"></a>207 Stronger head-to-head, third party follows | §§5.4,8.2 | FR A par pic / GE A bur / IT A gas | par>bur, bur>par, pic S par>bur, gas>par | par>bur, ×bur, gas>par. The dislodged army cannot contest Paris. |
| <a id="case-208"></a>208 Three-unit rotation | §5.5 | FR A par / GE A bur / IT A pic | par>bur, bur>pic, pic>par | par>bur, bur>pic, pic>par, preserving individual identities. |
| <a id="case-209"></a>209 Break one leg of rotation | §§5.2–5.5 | FR A par / GE A bur / IT A pic / EN A bre | par>bur, bur>pic, pic>par, bre>pic | Same. The Picardy tie blocks every departure in the cycle. |
| <a id="case-210"></a>210 Reorder simultaneous moves | §§1,11 | FR A par / GE A bur / IT A pic | pic>par, bur>pic, par>bur | par>bur, bur>pic, pic>par, exactly as in 208. |
| <a id="case-211"></a>211 Dislodged attacker still contests elsewhere [T] | §8.3 | FR A par pic / GE A bur / IT A tyr | par>bur, pic S par>bur, bur>mun, tyr>mun | par>bur, ×bur. Tyrolia stays. Munich empty by standoff. |
| <a id="case-212"></a>212 Even supported dislodged attack has no effect at dislodger's origin | §8.2, PDF p. 12 | TR A bul F bla / RU A rum gre ser sev | bul>rum, bla S bul>rum, rum>bul, gre S rum>bul, ser S rum>bul, sev>rum | rum>bul, ×bul, sev>rum. Strength three beats two. |

### Support, cuts and self-dislodgement — Spring 1901

| ID and purpose | Rule | Starting units | Orders in submission order | Expected result |
|---|---|---|---|---|
| <a id="case-301"></a>301 Foreign support reaches destination | §§6.1,6.3–6.4 | FR A pic / EN F nth / GE A bel | pic>bel, nth S pic>bel | pic>bel, ×bel. North Sea need not touch Picardy. |
| <a id="case-302"></a>302 Hold support ties supported attack | §§4,6.2,6.4 | FR A pic / EN F nth / GE A bel ruh | pic>bel, nth S pic>bel, bel H, ruh S bel H | Same, two against two. |
| <a id="case-303"></a>303 Unordered defender can receive support | §§4,6.2 | FR A pic / EN F nth / GE A bel ruh | pic>bel, nth S pic>bel, ruh S bel H | Same, including the unordered army in Belgium. |
| <a id="case-304"></a>304 Moving defender loses hold support | §§4,6.2 | FR A pic / EN F nth / GE A bel ruh hol | pic>bel, nth S pic>bel, bel>hol, ruh S bel H | pic>bel, ×bel. The move into friendly occupied Holland fails. |
| <a id="case-305"></a>305 Support mismatches destination | §6.3 | FR A pic / EN F eng / GE A bel | pic>bel, eng S pic>bre | Same. Both proposed destinations are reachable. Picardy actually attacks Belgium at one. |
| <a id="case-306"></a>306 Support mismatches origin | §6.3 | FR A pic bur / EN F nth / GE A bel | pic>bel, nth S bur>bel | Same. Burgundy never attacks. |
| <a id="case-307"></a>307 Move support cannot defend a holder | §§6.2–6.3 | FR A pic bur / GE A bel ruh | pic>bel, bur S pic>bel, ruh S bel>hol | pic>bel, ×bel. Ruhr does not increase Belgium's defense. |
| <a id="case-308"></a>308 Fleet cannot support inland | §6.1 | GE F ber A pru / RU A sil | pru>sil, ber S pru>sil | Same. Berlin's land border with Silesia gives the fleet no support reach. |
| <a id="case-309"></a>309 Army cannot support at sea | §6.1 | EN A lon F yor / GE F nth | yor>nth, lon S yor>nth | Same. |
| <a id="case-310"></a>310 Coastline limits support | §§2.2,6.1 | IT F rom A ven / AU A apu | ven>apu, rom S ven>apu | Same. Rome's fleet cannot sail to Apulia. |
| <a id="case-311"></a>311 Opposite-coast support counts | §§2.3,6.1 | FR F gas / IT F gol / GE A spa | gas>spa-nc, gol S gas>spa | gas>spa-nc, ×spa. Gulf of Lyon reaches only the south coast and supports the whole province. |
| <a id="case-312"></a>312 Equal supported attacks on vacancy [T] | §§5.2,6.4 | FR A pic bur / GE A ruh hol | pic>bel, bur S pic>bel, ruh>bel, hol S ruh>bel | Same. Belgium empty, two against two. |
| <a id="case-313"></a>313 Extra support breaks tie | §§5.2,6.4 | FR A pic bur / GE A ruh hol / EN F nth | pic>bel, bur S pic>bel, ruh>bel, hol S ruh>bel, nth S pic>bel | pic>bel. Three beats two. |
| <a id="case-314"></a>314 Outside attack cuts, even on bounce [T] | §7.1 | GE A pru sil / RU A war boh | pru>war, sil S pru>war, boh>sil | Same. Bohemia cuts Silesia while failing to enter. |
| <a id="case-315"></a>315 Same-power attack preserves support | §7.4 | GE A pru sil boh / RU A war | pru>war, sil S pru>war, boh>sil | pru>war, ×war. Germany's attack on Silesia leaves its support intact. |
| <a id="case-316"></a>316 Attack from supported destination | §§7.1–7.2 | GE A pru sil / RU A war | pru>war, sil S pru>war, war>sil | pru>war, ×war. Silesia survives the unsupported attack from Warsaw. |
| <a id="case-317"></a>317 Dislodgement from supported destination [T] | §§7.1–7.2, PDF p. 12 | GE F ber A sil / RU A pru war F bal | ber>pru, sil S ber>pru, pru>sil, war S pru>sil, bal>pru | pru>sil, ×sil. Berlin and Baltic Sea stay. Prussia ends empty by standoff after Silesia loses its support. |
| <a id="case-318"></a>318 Dislodged attacker cuts elsewhere [T] | §§7.3,8.3 | GE A pru sil / RU A war boh / AU A tyr mun | pru>war, sil S pru>war, boh>sil, tyr>boh, mun S tyr>boh | tyr>boh, ×boh. Silesia's support is cut, so Prussia and Warsaw stay. |
| <a id="case-319"></a>319 Cut hold support | §§6.2,7.1 | FR A pic bur / GE A bel ruh / IT A mun | pic>bel, bur S pic>bel, ruh S bel H, mun>ruh | pic>bel, ×bel. Munich bounces in Ruhr and cuts its defensive support. |
| <a id="case-320"></a>320 Support a supporter to hold | §§4,6.2,7.1 | FR A pic bur par / GE A mun ruh bel | pic>bel, bur S pic>bel, par S bur H, mun>bur, ruh S mun>bur | Same. Paris saves Burgundy from dislodgement. The attack still cuts Burgundy's support, saving Belgium. |
| <a id="case-321"></a>321 Own occupant resists foreign-backed attacker | §9.1 | FR A pic bel / EN F nth | pic>bel, nth S pic>bel | Same. Foreign support cannot force French self-dislodgement. |
| <a id="case-322"></a>322 Deny defender's support to foreign attack | §9.2 | FR A pic / GE A bel ruh | pic>bel, ruh S pic>bel | Same. Ruhr cannot help evict Germany's own army. |
| <a id="case-323"></a>323 Independently sufficient foreign force | §9.2 | FR A pic bur / GE A bel ruh | pic>bel, ruh S pic>bel, bur S pic>bel | pic>bel, ×bel. French strength is two without Ruhr. |
| <a id="case-324"></a>324 Own attacks oppose foreign attack | §9.3 | FR A pic bur / GE A ruh | pic>bel, bur>bel, ruh>bel | Same. All three bounce, including the two French armies. |
| <a id="case-325"></a>325 Unrequested support breaks self-standoff | §§6.1,9.3 | FR A pic bur / GE A ruh / EN F nth | pic>bel, bur>bel, ruh>bel, nth S pic>bel | pic>bel. Picardy must accept England's support. |
| <a id="case-326"></a>326 Attempted self-dislodgement opposes a rival | §§9.1,9.3, PDF p. 15 | FR A pic bel bur / GE A hol ruh | pic>bel, bur S pic>bel, hol>bel, ruh S hol>bel | Same. The French attack ties the German attack at two while the French occupant stays. |
| <a id="case-327"></a>327 Hold support for failed non-move | §§3,6.2 | FR A pic bur / GE A bel ruh par | pic>bel, bur S pic>bel, bel S par>pic, ruh S bel H | Same. The support issued from Belgium has no matching move. Belgium can still receive Ruhr's hold support. |
| <a id="case-328"></a>328 Illegal move becomes a hold | §§1,4,6.2 | FR A pic bur / GE A bel ruh | pic>bel, bur S pic>bel, bel>mos, ruh S bel H | Same. The impossible move is normalized to hold before comparing strength. Compare 304's legal but failed move. |
| <a id="case-329"></a>329 Defender's support when defender vacates | §§5.1,9.2 | FR A pic / GE A bel ruh / EN F nth | pic>bel, ruh S pic>bel, bel>hol, nth>bel | bel>hol, pic>bel. Ruhr can help Picardy beat North Sea since no German unit is being dislodged. |
| <a id="case-330"></a>330 Illegal attack cannot cut support | §§1,6.4,7.1 | GE A pru sil / RU A war F ber | pru>war, sil S pru>war, ber>sil | pru>war, ×war. A fleet's illegal order into inland Silesia leaves its support intact. |

### Convoys — Spring 1901

| ID and purpose | Rule | Starting units | Orders in submission order | Expected result |
|---|---|---|---|---|
| <a id="case-401"></a>401 Single-fleet crossing | §10.1 | EN A lon F eng | lon>bel, eng C lon>bel | lon>bel. Fleet stays at sea. |
| <a id="case-402"></a>402 Missing transport | §§2.1,10.1 | EN A lon F eng | lon>bel | Same. Fleet presence alone supplies no convoy. |
| <a id="case-403"></a>403 Wrong destination on convoy order | §10.2 | EN A lon F eng | lon>bel, eng C lon>pic | Same. |
| <a id="case-404"></a>404 Wrong army on convoy order | §§10.1–10.2 | EN A lon wal F eng | lon>bel, eng C wal>bel | Same. |
| <a id="case-405"></a>405 Multinational chain | §10.1, PDF p. 14 | EN A lon F eng / FR F mid / IT F wes | lon>tun, eng C lon>tun, mid C lon>tun, wes C lon>tun | lon>tun. All three fleets stay. |
| <a id="case-406"></a>406 Missing middle link | §10.1 | EN A lon F eng / FR F mid / IT F wes | lon>tun, eng C lon>tun, wes C lon>tun | Same. The endpoints cannot bridge the gap. |
| <a id="case-407"></a>407 Coastal fleet cannot bridge two seas | §10.1 | TR A ank F bla con aeg | ank>gre, bla C ank>gre, con C ank>gre, aeg C ank>gre | Same. Constantinople cannot act as a convoying fleet between Black Sea and Aegean Sea. |
| <a id="case-408"></a>408 Fleet as passenger [F,D] | §10.1 | EN F lon eng | lon>bel, EN literal `F eng C F lon - bel` | Same. Only an army can be carried. |
| <a id="case-409"></a>409 Inland convoy endpoint | §10.1 | EN A lon F eng | lon>par, eng C lon>par | Same. No convoy can end in Paris. |
| <a id="case-410"></a>410 Convoy adds no strength | §§10.1,10.5 | EN A lon F eng / GE A bel | lon>bel, eng C lon>bel | Same. Army bounces back to London. |
| <a id="case-411"></a>411 Supported landing | §§6.3,10.5 | EN A lon F eng / GE A bel / FR A pic | lon>bel, eng C lon>bel, pic S lon>bel | lon>bel, ×bel. |
| <a id="case-412"></a>412 Mere attack on transport | §10.3 | EN A lon F eng / GE A bel / FR A pic / IT F bre | lon>bel, eng C lon>bel, pic S lon>bel, bre>eng | lon>bel, ×bel. Brest bounces, transport survives. |
| <a id="case-413"></a>413 Transport dislodged, supported landing fails | §§10.3,10.5 | EN A lon F eng / GE A bel / FR A pic / IT F bre iri | lon>bel, eng C lon>bel, pic S lon>bel, bre>eng, iri S bre>eng | bre>eng, ×eng. London and Belgium stay. Picardy cannot deliver an army without a route. |
| <a id="case-414"></a>414 Alternate route survives | §10.4 | EN A lon F eng nth / GE A bel / FR A pic / IT F bre iri | lon>bel, eng C lon>bel, nth C lon>bel, pic S lon>bel, bre>eng, iri S bre>eng | bre>eng, ×eng, lon>bel, ×bel. North Sea carries the army. |
| <a id="case-415"></a>415 Every alternate route broken | §§10.3–10.4 | EN A lon F eng nth / IT F bre iri / GE F hel den | lon>bel, eng C lon>bel, nth C lon>bel, bre>eng, iri S bre>eng, hel>nth, den S hel>nth | bre>eng, ×eng, hel>nth, ×nth. London stays. Belgium empty. |
| <a id="case-416"></a>416 Bounced passenger dislodged at home | §10.5 | EN A lon F eng / GE A bel / FR A yor wal | lon>bel, eng C lon>bel, yor>lon, wal S yor>lon | yor>lon, ×lon. Belgium's defender stays, no intermediate army at sea. |
| <a id="case-417"></a>417 Fleet receives hold support while convoying | §§4,6.2,10.3 | EN A lon F eng mid / FR F bre iri | lon>bel, eng C lon>bel, mid S eng H, bre>eng, iri S bre>eng | lon>bel. Fleet in Channel survives two against two. |
| <a id="case-418"></a>418 One fleet ordered to ferry two armies [D] | H2, §10.1 | EN A lon wal F eng | lon>bel, wal>pic, eng C lon>bel, eng C wal>pic | lon>bel only. Wales stays. Second convoy instruction ignored. |
| <a id="case-419"></a>419 Both sides exchange by convoy | §10.6, PDF p. 16 | EN A lon F nth / FR A bel F eng | lon>bel, bel>lon, nth C lon>bel, eng C bel>lon | lon>bel, bel>lon. Fleets stay. |
| <a id="case-420"></a>420 Adjacent exchange, own transport | §§10.6–10.7 | EN A lon F eng / FR A wal | lon>wal, wal>lon, eng C lon>wal | lon>wal, wal>lon. England's move must use the convoy. |
| <a id="case-421"></a>421 Adjacent exchange, foreign transport | §§5.4,10.7 | EN A lon / IT F eng / FR A wal | lon>wal, wal>lon, eng C lon>wal | Same. The English move defaults to land, creating a head-to-head. |
| <a id="case-422"></a>422 Explicit foreign convoy [F] | §10.7, PDF p. 16 | EN A lon / IT F eng / FR A wal | EN literal `A lon - wal via convoy`, wal>lon, eng C lon>wal | lon>wal, wal>lon. Requires parsing explicit convoy intent. |
| <a id="case-423"></a>423 Disruption cannot restore discarded land route | §10.7 | EN A lon F eng / FR F bre iri | lon>wal, eng C lon>wal, bre>eng, iri S bre>eng | bre>eng, ×eng. London stays even though Wales is adjacent and empty. |
| <a id="case-424"></a>424 Necessary route under attack from destination [T] | §7.5, PDF p. 17 | FR A tun F tyn / IT F rom nap | tun>nap, tyn C tun>nap, rom>tyn, nap S rom>tyn | rom>tyn, ×tyn. Naples keeps its support, so Tunis stays. |
| <a id="case-425"></a>425 Independent route changes the cut [T] | §§7.5,10.4, PDF p. 17 | FR A tun F tyn ion / IT F rom nap | tun>nap, tyn C tun>nap, ion C tun>nap, rom>tyn, nap S rom>tyn | Same. Ionian route lets the army attack Naples and cut support. Rome bounces and the landing also bounces. |
| <a id="case-426"></a>426 Independent route with landing support | §7.5, PDF p. 17 | FR A tun apu F tyn ion / IT F rom nap | tun>nap, tyn C tun>nap, ion C tun>nap, rom>tyn, nap S rom>tyn, apu S tun>nap | tun>nap, ×nap. Rome bounces, both convoy fleets stay. |
| <a id="case-427"></a>427 Inland origin | §10.1 | FR A par F eng | par>lon, eng C par>lon | Same. The convoy must embark at the coast. |
| <a id="case-428"></a>428 Army cannot transport another army [F,D] | §10.1 | EN A lon wal | lon>bel, EN literal `A wal C A lon - bel` | Same. Wales supplies no convoy. |
| <a id="case-429"></a>429 Support cannot be transported [F,D] | §6.1, PDF p. 13 | EN A pic lon F eng / FR F bre | pic>bre, lon S pic>bre, EN literal `F eng C A lon S A pic - bre` | Same. London lacks direct support reach. The fleet cannot ferry its support. |
| <a id="case-430"></a>430 Disrupted passenger can be dislodged | §§8.1,10.3 | EN A lon F eng / FR F bre iri / GE A yor wal | lon>bel, eng C lon>bel, bre>eng, iri S bre>eng, yor>lon, wal S yor>lon | bre>eng, ×eng, yor>lon, ×lon. Belgium stays empty. |
| <a id="case-431"></a>431 Convoy before fleet move [D] | H2, §10.1 | EN A lon F eng | lon>bel, eng C lon>bel, eng>bre | lon>bel. The Channel fleet convoys and stays at sea. Its later move is ignored. |
| <a id="case-432"></a>432 Fleet move before convoy [D] | H2, §10.1 | EN A lon F eng | lon>bel, eng>bre, eng C lon>bel | eng>bre. Ignore the later convoy. The army stays in London because it has no transport. |

### Conditional retreats — Spring 1901

All instructions in these rows arrive before resolution. A unit with only conditional retreat orders defends its origin normally during movement. Read its submitted retreats in order if it is dislodged.

| ID and purpose | Rule | Starting units | Orders in submission order | Expected result |
|---|---|---|---|---|
| <a id="case-501"></a>501 Condition false [D] | H1, §12.1 | FR A pic / GE A bel | pic>bel, bel R hol | Same. Belgium defends at one. The conditional retreat never activates. |
| <a id="case-502"></a>502 Condition true | H1, §§8.1,12.2 | FR A pic bur / GE A bel | pic>bel, bur S pic>bel, bel R hol | pic>bel, bel R hol. A retreat order gives no movement-phase advantage or penalty. |
| <a id="case-503"></a>503 No retreat submitted | H1, §12.4, PDF p. 18 | FR A pic bur / GE A bel | pic>bel, bur S pic>bel | pic>bel, ×bel, even with Holland available. |
| <a id="case-504"></a>504 Movement submitted before retreat [D] | H1–H2, §12.2 | FR A pic bur / GE A bel ruh | bel>ruh, bel R hol, pic>bel, bur S pic>bel | pic>bel, bel R hol. Belgium's accepted move fails against the friendly occupant in Ruhr. Dislodgement activates the separate retreat to Holland. |
| <a id="case-505"></a>505 Retreat submitted before movement [D] | H1–H2, §12.2 | FR A pic bur / GE A bel ruh | bel R hol, bel>ruh, pic>bel, bur S pic>bel | pic>bel, bel R hol. The move to Ruhr is accepted and fails during movement. The earlier-submitted retreat then succeeds. |
| <a id="case-506"></a>506 Retreat submitted first and unit escapes by moving [D] | H1–H2 | FR A pic bur / GE A bel | bel R hol, bel>ruh, pic>bel, bur S pic>bel | bel>ruh, pic>bel. The German army reaches vacant Ruhr during movement. Its conditional retreat stays dormant and Holland stays empty. |
| <a id="case-507"></a>507 Occupied refuge | §12.2, §12.4 | FR A pic bur / GE A bel hol | pic>bel, bur S pic>bel, bel R hol | pic>bel, ×bel. Holland's occupant stays. |
| <a id="case-508"></a>508 Newly vacated refuge | §12.2 | FR A pic bur / GE A bel hol | pic>bel, bur S pic>bel, hol>kie, bel R hol | pic>bel, hol>kie, bel R hol. Use post-movement occupancy. |
| <a id="case-509"></a>509 Newly occupied refuge | §12.2 | FR A pic bur / GE A bel / EN F hel | pic>bel, bur S pic>bel, hel>hol, bel R hol | pic>bel, hel>hol, ×bel. |
| <a id="case-510"></a>510 Empty dislodger origin forbidden | §§8.2,12.2 | FR A pic bur / GE A bel | pic>bel, bur S pic>bel, bel R pic | pic>bel, ×bel. Picardy stays empty. |
| <a id="case-511"></a>511 Empty standoff province forbidden [T] | §§5.2,12.2 | FR A pic bur / GE A bel / EN F nth / IT A kie | pic>bel, bur S pic>bel, nth>hol, kie>hol, bel R hol | pic>bel, ×bel. Holland remains empty after the bounce and the failed retreat. |
| <a id="case-512"></a>512 Nonadjacent refuge | §§2.1,12.2 | FR A pic bur / GE A bel | pic>bel, bur S pic>bel, bel R ber | pic>bel, ×bel. |
| <a id="case-513"></a>513 Retreat cannot use a convoy | §§10.1,12.2 | FR A pic bur / GE A bel F nth | pic>bel, bur S pic>bel, nth C bel>lon, bel R lon | pic>bel, ×bel. London empty and fleet stays in North Sea. |
| <a id="case-514"></a>514 Fleet retreats to reachable coast | §§2.3,12.2 | FR F gas / EN F mid A par | mid>gas, par S mid>gas, gas R spa-nc | mid>gas, gas R spa-nc. |
| <a id="case-515"></a>515 Fleet retreats to unreachable coast | §§2.3,12.2,12.4 | FR F gas / EN F mid A par | mid>gas, par S mid>gas, gas R spa-sc | mid>gas, ×gas. |
| <a id="case-516"></a>516 No legal escape | §§12.2,12.4 | FR F por / EN F mid spa-nc | mid>por, spa-nc S mid>por | mid>por, ×por. Every reachable province is the dislodger's origin or occupied. |
| <a id="case-517"></a>517 Omit retreat after ordering a hold | H1, §12.4, PDF p. 18 | FR A pic bur / GE A bel | bel H, pic>bel, bur S pic>bel | pic>bel, ×bel. Germany ordered its army to hold and submitted no retreat. Disband it on dislodgement even though Holland is open. |
| <a id="case-518"></a>518 Two retreats collide | §12.3 | EN A par / GE A bel / FR A bre gas / IT A hol ruh | bre>par, gas S bre>par, hol>bel, ruh S hol>bel, par R pic, bel R pic | bre>par, hol>bel, ×par, ×bel. Picardy vacant. |
| <a id="case-519"></a>519 Distinct simultaneous retreats | §12.3 | EN A par / GE A bel / FR A bre gas / IT A hol ruh | bre>par, gas S bre>par, hol>bel, ruh S hol>bel, par R pic, bel R bur | bre>par, hol>bel, par R pic, bel R bur. |
| <a id="case-520"></a>520 Three-way retreat collision | §12.3 | EN A par bel bur / FR A bre gas / IT A hol ruh / GE A mun mar | bre>par, gas S bre>par, hol>bel, ruh S hol>bel, mun>bur, mar S mun>bur, par R pic, bel R pic, bur R pic | bre>par, hol>bel, mun>bur, ×par, ×bel, ×bur. Third arrival cannot claim Picardy, even with three units of one power. |
| <a id="case-521"></a>521 Retreats cannot receive support | §§6.3,12.2–12.3 | EN A par / GE A bel F eng / FR A bre gas / IT A hol ruh | bre>par, gas S bre>par, hol>bel, ruh S hol>bel, eng S bel>pic, par R pic, bel R pic | bre>par, hol>bel, ×par, ×bel. No movement order matches the support. It cannot break the retreat tie. |
| <a id="case-522"></a>522 Retreats to different coasts collide | §§2.3,12.3 | FR F gas / IT F mar / GE A par bre / AU A bur pie | par>gas, bre S par>gas, bur>mar, pie S bur>mar, gas R spa-nc, mar R spa-sc | par>gas, bur>mar, ×gas, ×mar. All of Spain stays vacant. |
| <a id="case-523"></a>523 Two usable conditional retreats [D] | H1–H2 | FR A pic bur / GE A bel | bel R hol, bel R ruh, pic>bel, bur S pic>bel | pic>bel, bel R hol. Ignore second retreat even though Ruhr is open. |
| <a id="case-524"></a>524 First refuge occupied after movement [D] | H1–H2, §12.2 | FR A pic bur / GE A bel / EN F nth | bel R hol, bel R ruh, pic>bel, bur S pic>bel, nth>hol | pic>bel, nth>hol, bel R ruh. England occupies Holland before retreats. Skip that invalid retreat and accept Ruhr. |
| <a id="case-525"></a>525 Nonadjacent first refuge, legal second refuge [D] | H1–H2, §12.2 | FR A pic bur / GE A bel | bel R ber, bel R hol, pic>bel, bur S pic>bel | pic>bel, bel R hol. Skip Berlin because it is nonadjacent. Holland is the army's first legal retreat. |
| <a id="case-526"></a>526 Accepted retreat collides despite later alternative [D,T] | H1–H2, §12.3 | EN A par / GE A bel / FR A bre gas / IT A hol ruh | bre>par, gas S bre>par, hol>bel, ruh S hol>bel, par R pic, bel R pic, bel R bur | bre>par, hol>bel, ×par, ×bel. Both first retreats legally target empty Picardy. Ignore Belgium's later retreat to Burgundy. The collision disbands both armies and leaves Picardy and Burgundy empty. |

### Ownership and adjustments — Fall 1901 unless stated

Control lists describe the start of movement. Count centers and surviving units **after Fall retreats** and before executing any builds or required removals. A normal power with spare builds and no build instruction declines them. Newly built units never move that turn. Observe H3 in the explicitly marked Spring controls.

| ID and purpose | Rule | Starting units and control | Orders in submission order | Expected result |
|---|---|---|---|---|
| <a id="case-601"></a>601 Spring occupation | §13 | **Spring**. FR A pic. SC GE bel | pic>bel | pic>bel. Belgium remains German. |
| <a id="case-602"></a>602 Fall occupation transfers control | §13 | FR A pic. SC GE bel | pic>bel | pic>bel. Belgium becomes French, supplying the one army. |
| <a id="case-603"></a>603 Spring visitor leaves by Fall | §13, PDF p. 19 | FR A bel. SC GE bel, FR par | bel>pic | bel>pic. Empty Belgium remains German. This is the Fall snapshot after a Spring visit. |
| <a id="case-604"></a>604 Capture neutral by staying | §13 | EN F bel. No controlled centers | bel H | Same units, Belgium becomes English. |
| <a id="case-605"></a>605 Vacant center retains control | §13 | GE A bel. SC GE bel | bel>ruh | bel>ruh. Belgium remains German and supplies Ruhr. |
| <a id="case-606"></a>606 Retreat captures center and funds build | H1, §§12.2,13,14.2 | FR A pic bur / GE A bel. SC FR par bre, GE bel mun | pic>bel, bur S pic>bel, bel R hol, GE +A mun | pic>bel, bel R hol. Belgium becomes French, Holland German. Germany has two centers and one survivor, so builds a German army in Munich [B]. France waives its spare build. |
| <a id="case-607"></a>607 Disband before census | H1, §§12.4,13,14 | FR A pic bur / GE A bel. SC FR par bre, GE bel mun | pic>bel, bur S pic>bel, GE +A mun | pic>bel, ×bel. Belgium becomes French. Germany has one center and zero survivors, so builds a German army in Munich [B]. |
| <a id="case-608"></a>608 Spring build ignored [D] | H3, §14.2 | **Spring**. No units. SC FR bre | FR +F bre | Same. Surplus and valid site, adjustment ignored for season. |
| <a id="case-609"></a>609 Fall coastal fleet build [B] | §14.2 | No units. SC FR bre | FR +F bre | Build French fleet at Brest. |
| <a id="case-610"></a>610 Spring removal ignored [D] | H3, §14.1 | **Spring**. FR A pic bur. SC FR par | FR -pic | Same, despite a one-unit deficit. |
| <a id="case-611"></a>611 Fall owner selects removal | §14.1 | FR A pic bur. SC FR par | FR -pic | Remove French army in Picardy. Burgundy stays. |
| <a id="case-612"></a>612 Multiple required removals | §14.1 | FR A pic bur F eng. SC FR par | FR -bur, FR -eng | Remove Burgundy and Channel units. Picardy stays. |
| <a id="case-613"></a>613 Balanced count | §14 | FR A pic. SC FR par | FR +A par, FR -pic | Same. Neither an extra build nor a discretionary adjustment removal is due. |
| <a id="case-614"></a>614 Fleet build inland | §14.2 | No units. SC FR par | FR +F par | Same. |
| <a id="case-615"></a>615 Army build inland [B] | §14.2 | No units. SC FR par | FR +A par | Build French army at Paris. |
| <a id="case-616"></a>616 Army build on coast [B] | §14.2 | No units. SC FR bre | FR +A bre | Build French army at Brest, null coast. |
| <a id="case-617"></a>617 Fleet build at named split coast [B] | §14.2, PDF p. 19 | No units. SC RU stp | RU +F stp-sc | Build Russian fleet at St. Petersburg south coast. |
| <a id="case-618"></a>618 Fleet build missing coast [D] | §14.2 | No units. SC RU stp | RU +F stp | Same. Russia must name a coast even when the center is empty. |
| <a id="case-619"></a>619 Capture cannot create a home build site | §14.2 | FR A pic. SC FR bel hol | FR +A bel | Same. France retains its army and two foreign centers but controls no home site. |
| <a id="case-620"></a>620 Home site under foreign control | §14.2 | No units. SC FR bre, GE par | FR +A par | Same. Paris is French home territory under German control. |
| <a id="case-621"></a>621 Home site occupied by builder | §14.2 | FR A bre. SC FR bre par | FR +F bre | Same. French surplus cannot stack a second unit at Brest. |
| <a id="case-622"></a>622 Home site occupied by enemy | §§13,14.2 | GE A bre. SC FR bre par | FR +F bre | Units stay. Brest becomes German. France still has one build but cannot use Brest. |
| <a id="case-623"></a>623 Vacate home site and rebuild there [B] | §§5.1,14.2 | FR A bre. SC FR bre par | bre>pic, FR +F bre | bre>pic and build a new French fleet at Brest. The build names a site, not the departed army. |
| <a id="case-624"></a>624 Recaptured home is still occupied | §§13,14.2 | FR A pic. SC FR bel, GE par | pic>par, FR +A par | pic>par, Paris becomes French. France has one spare build but cannot place it on the occupying army. |
| <a id="case-625"></a>625 Two builds, one province | §14.2 | No units. SC FR bre par | FR +A bre, FR +A bre | Exactly one new French army at Brest [B]. Paris stays empty. Duplicate commands have the same visible effect [D]. |
| <a id="case-626"></a>626 Excess build requests for distinct sites [B,D] | H3, §14.2 | FR A pic. SC FR par bre | FR +A par, FR +F bre | Build a French army at Paris. This fills France's one available build. Ignore the later fleet request and leave Brest empty. |
| <a id="case-627"></a>627 Partial waiver [B] | §14.2 | No units. SC FR par bre | FR +A par | Build French army at Paris only. Brest stays empty, one build declined. |
| <a id="case-628"></a>628 Entire surplus waived | §14.2 | No units. SC FR par bre | [] | Same. No automatic builds. |
| <a id="case-629"></a>629 Waiver does not bank a build [V] | §14.2 | **Fall 1902**. FR A pic / GE F gas. SC FR par bre, GE kie. Same position in Fall 1901 and Spring 1902. | gas>bre, FR +A par | gas>bre, Brest becomes German. Paris remains vacant. France has one unit and one center. In runner setup resolve the preceding Fall and Spring with `pic H, gas H` each time, including no French build in the preceding Fall. Assert the intermediate boards equal this input. |
| <a id="case-630"></a>630 Invalid removal names empty site [D] | §§14.1–14.2 | FR A pic bur. SC FR par | FR literal `D A gas`, FR -bur | Remove Burgundy only. Invalid Gascony request cannot satisfy the deficit. |
| <a id="case-631"></a>631 Active player submits no required removal [R] | H3, §14.1 | FR A pic bur. SC FR par. France is active. | pic H | Randomly remove one French army. France must finish with exactly one unit for its one center. Either original army may survive in place. |
| <a id="case-632"></a>632 Surplus removals requested [D] | H3, §14.1 | FR A pic bur. SC FR par | FR -pic, FR -bur | Remove Picardy, the first valid choice. France now has one unit for one center. Ignore the later removal and retain the army in Burgundy. |
| <a id="case-633"></a>633 Split-coast home army [B] | §§2.3,14.2 | No units. SC RU stp | RU +A stp | Build Russian army at St. Petersburg with null coast. |
| <a id="case-634"></a>634 Unusable coastal build without type [F,D] | §14.2, PDF p. 19 | No units. SC FR bre | FR literal `bre` | Same. Coastal build must distinguish army from fleet. |
| <a id="case-635"></a>635 Build after releasing recovered center [B] | §§13,14.2 | **Fall 1902**. FR A par. SC FR par bel. Position after 624 and a quiet Spring. | par>pic, FR +A par | par>pic, build a new French army at Paris. Control of both centers persists. |
| <a id="case-636"></a>636 Capture eliminates need for submitted removal | §§13,14.1 | FR A pic bur. SC FR par | pic>bre, FR -pic | pic>bre, Brest becomes French. The Fall count gives France two centers for its two armies. Both survive. |
| <a id="case-637"></a>637 Capture funds build in the same turn [B] | §§13,14.2 | FR A pic. SC FR par | pic>bel, FR +A par | pic>bel, Belgium becomes French. Build a French army at Paris. Before movement France had zero spare builds. |
| <a id="case-638"></a>638 Enemy retreat cancels apparent build [D] | H1, §§12.2,13,14.2 | FR A pic bur / GE A bel. SC FR par hol, GE bel | pic>bel, bur S pic>bel, bel R hol, FR +A par | pic>bel, bel R hol. Belgium becomes French and Holland German. With Paris and Belgium supplying its two armies, France has no build. Counting before Germany's retreat would incorrectly add an army at Paris. |
| <a id="case-639"></a>639 Active player submits one of two required removals [R] | H3, §14.1 | FR A pic bur F eng. SC FR par. France is active. | FR -eng | Remove the Channel fleet as ordered, then randomly remove one of the two French armies. The Channel must be empty. Exactly one French army survives in place for France's one center. |

### Civil disorder and reduced-player games

Every row in this table requires **V** setup outside the province list. In 701–707 the named power is in civil disorder and submits no orders. Under H3, every surviving unit of that power is eligible for random removal when the Fall count shows a deficit.

| ID and purpose | Rule | Season, units, control and context | Orders in submission order | Expected result |
|---|---|---|---|---|
| <a id="case-701"></a>701 Civil units hold independently [T] | §15, PDF p. 19 | Spring 1901. IT A rom ven, Italy in disorder / FR A tus apu | tus>rom, apu S tus>rom | tus>rom, ×rom. Venice gives no automatic support and the civil unit disbands despite available retreats. |
| <a id="case-702"></a>702 Other powers may defend civil units | §15 | Spring 1901. IT A rom ven, Italy in disorder / FR A tus apu / AU F tyn | tus>rom, apu S tus>rom, tyn S rom H | Same. Foreign hold support saves Rome. |
| <a id="case-703"></a>703 No civil builds [D] | §15 | Fall 1901. IT A ven, Italy in disorder. SC IT ven rom | [] | Same. Verify Italy receives no build despite its surplus and vacant Rome. An active player waiving a build would give the same board. |
| <a id="case-704"></a>704 Random removal at home or abroad [R] | H3, §§14.1,15 | Fall 1901. FR A gal F bre, France in disorder. SC FR bre | [] | Randomly remove one French unit. Exactly one survives in place, either the army in Galicia or the fleet in Brest. Brest remains French even if the fleet is removed. |
| <a id="case-705"></a>705 Random removal from a mixed force [R] | H3, §§14.1,15 | Fall 1901. FR F mid A pie, France in disorder. SC FR par | [] | France retains one unit for Paris. Either the fleet in the Mid Atlantic or the army in Piedmont survives in place. |
| <a id="case-706"></a>706 Random removal from two armies [R] | H3, §§14.1,15 | Fall 1901. FR A pie spa, France in disorder. SC FR spa | [] | France retains one army, in either Piedmont or Spain. Spain remains French whichever army is removed. |
| <a id="case-707"></a>707 English armies abroad during civil disorder [R] | H3, §§14.1,15 | Fall 1901. EN A tyr pic, England in disorder. SC EN lon | [] | Randomly remove one English army. London supplies exactly one survivor, which may be the original army in either Tyrolia or Picardy. |
| <a id="case-708"></a>708 Six-player neutral dislodged | §15, PDF p. 19 | Spring 1901. Six players, Italy eliminated. IT A rom ven F nap / AU A tyr F adr | tyr>ven, adr S tyr>ven | tyr>ven, ×ven. The other neutral Italian units stay. |
| <a id="case-709"></a>709 Neutral receives outside hold support | §15, PDF p. 19 | Spring 1901. Six players, Italy eliminated. IT A rom ven F nap / AU A tyr F adr / FR A pie | tyr>ven, adr S tyr>ven, pie S ven H | Same. Venice survives two against two. |
| <a id="case-710"></a>710 Five-player second neutral power | §15, PDF p. 19 | Spring 1901. Five players, Germany and Italy eliminated. GE A mun ber F kie / IT A rom ven F nap / FR A bur ruh | bur>mun, ruh S bur>mun | bur>mun, ×mun. Kiel and Berlin do not automatically defend Munich. The other neutral units stay. |
| <a id="case-711"></a>711 Shared player, separate national counts [B] | PDF p. 19, §14 | Fall 1901. Four players, same person plays AU and FR. FR A pic bur. SC FR par, AU vie bud | FR -pic, AU +A vie, AU +A bud | Remove Picardy. Build Austrian armies in Vienna and Budapest. France keeps Burgundy. Compute each country's adjustment independently. |
| <a id="case-712"></a>712 Italian neutrality bars entry | PDF p. 19 | **Spring 1914**, two-player WWI variant, Italy neutral. FR A mar | mar>pie | Same. Piedmont is Italian territory even when empty. Seed the variant calendar and pre-alliance state. |
| <a id="case-713"></a>713 Italian neutrality ends | PDF p. 19 | **Spring 1915**, same variant, Italy has joined England/France/Russia after the fixed Fall 1914 coin result. FR A mar | mar>pie | mar>pie. Fix the alliance outcome in game setup. No random draw during the test. |

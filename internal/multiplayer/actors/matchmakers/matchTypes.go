package matchmakers

type MatchType int

const (
	RandomMatch MatchType = iota
	RankedMatch
	GuildMatch
	CustomMatch
)

var matchTypeNames = map[MatchType]string{
	RandomMatch: "Random",
	RankedMatch: "Ranked",
	GuildMatch:  "Guild",
	CustomMatch: "Custom",
}

func (t MatchType) String() string {
	return matchTypeNames[t]
}

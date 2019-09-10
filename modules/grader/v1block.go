package grader

import "sort"

// V1GradedBlock is an opr set that has been graded. The set should be read only through it's interface
// implementation.
type V1GradedBlock struct {
	baseGradedBlock
}

var _ GradedBlock = (*V1GradedBlock)(nil)

func (g *V1GradedBlock) Version() uint8 {
	return 1
}

func (g *V1GradedBlock) Winners() []*GradingOPR {
	if len(g.oprs) < 10 {
		return nil
	}

	return g.oprs[:10]
}

func (g *V1GradedBlock) WinnersShortHashes() []string {
	return g.winnersShortHashes(10)
}

// WinnerAmount is the number of OPRs that receive a payout
func (g *V1GradedBlock) WinnerAmount() int {
	return 10
}

func (g *V1GradedBlock) grade() {
	if len(g.oprs) < 10 {
		return
	}

	if g.cutoff > len(g.oprs) {
		g.cutoff = len(g.oprs)
	}

	for i := g.cutoff; i >= 10; i-- {
		avg := averageV1(g.oprs[:i])
		for j := 0; j < i; j++ {
			gradeV1(avg, g.oprs[j])
		}
		// Because this process can scramble the sorted fields, we have to resort with each pass.
		sort.SliceStable(g.oprs[:i], func(i, j int) bool { return g.oprs[i].SelfReportedDifficulty > g.oprs[j].SelfReportedDifficulty })
		sort.SliceStable(g.oprs[:i], func(i, j int) bool { return g.oprs[i].Grade < g.oprs[j].Grade })
	}
}

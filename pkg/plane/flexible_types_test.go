package plane

import "testing"

func TestStateOutputFromIssuePrefersExpandedStateName(t *testing.T) {
	issue := Issue{
		State: FlexibleState{
			ID:   "state-1",
			Name: "Done",
		},
		StateName: "Fallback",
	}

	got := StateOutputFromIssue(issue)
	if got.ID != "state-1" {
		t.Fatalf("expected state id state-1, got %q", got.ID)
	}
	if got.Name != "Done" {
		t.Fatalf("expected state name Done, got %q", got.Name)
	}
}

func TestStateOutputFromIssueFallsBackToStateName(t *testing.T) {
	issue := Issue{
		State: FlexibleState{
			ID:   "state-2",
			Name: "",
		},
		StateName: "Todo",
	}

	got := StateOutputFromIssue(issue)
	if got.ID != "state-2" {
		t.Fatalf("expected state id state-2, got %q", got.ID)
	}
	if got.Name != "Todo" {
		t.Fatalf("expected state name Todo, got %q", got.Name)
	}
}

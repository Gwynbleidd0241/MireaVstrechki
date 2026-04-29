package service

import "testing"

func FuzzIsValidTaskStatus(f *testing.F) {
	f.Add("todo")
	f.Add("in_progress")
	f.Add("done")
	f.Add("")
	f.Add("TODO")
	f.Add("\x00")
	f.Add("todo\n")
	f.Add(" todo ")

	f.Fuzz(func(t *testing.T, s string) {
		valid := isValidTaskStatus(s)

		if !valid {
			return
		}

		if s != "todo" && s != "in_progress" && s != "done" {
			t.Errorf("isValidTaskStatus(%q) = true, but %q is not in the allowed set", s, s)
		}
	})
}

func FuzzIsValidParticipantRole(f *testing.F) {
	f.Add("participant")
	f.Add("responsible")
	f.Add("")
	f.Add("admin")
	f.Add("PARTICIPANT")
	f.Add(" participant")

	f.Fuzz(func(t *testing.T, s string) {
		valid := isValidParticipantRole(s)

		if !valid {
			return
		}

		if s != "participant" && s != "responsible" {
			t.Errorf("isValidParticipantRole(%q) = true, but %q is not in the allowed set", s, s)
		}
	})
}

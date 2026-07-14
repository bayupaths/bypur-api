package model

import "testing"

func TestBeforeCreateHooksAssignIDs(t *testing.T) {
	cases := []struct {
		name string
		hook func() (string, error)
	}{
		{"profile", func() (string, error) {
			m := &Profile{}
			return m.ID, m.BeforeCreate(nil)
		}},
		{"social link", func() (string, error) {
			m := &SocialLink{}
			err := m.BeforeCreate(nil)
			return m.ID, err
		}},
		{"offering", func() (string, error) {
			m := &Offering{}
			err := m.BeforeCreate(nil)
			return m.ID, err
		}},
		{"skill", func() (string, error) {
			m := &Skill{}
			err := m.BeforeCreate(nil)
			return m.ID, err
		}},
		{"experience", func() (string, error) {
			m := &Experience{}
			err := m.BeforeCreate(nil)
			return m.ID, err
		}},
		{"project", func() (string, error) {
			m := &Project{}
			err := m.BeforeCreate(nil)
			return m.ID, err
		}},
		{"contact message", func() (string, error) {
			m := &ContactMessage{}
			err := m.BeforeCreate(nil)
			return m.ID, err
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := tc.hook()
			if err != nil {
				t.Fatalf("BeforeCreate returned error: %v", err)
			}
			if id == "" {
				t.Fatal("expected hook to assign ID")
			}
		})
	}
}

func TestBeforeCreateHooksKeepExistingIDs(t *testing.T) {
	profile := &Profile{ID: "existing"}
	if err := profile.BeforeCreate(nil); err != nil {
		t.Fatalf("BeforeCreate returned error: %v", err)
	}
	if profile.ID != "existing" {
		t.Fatalf("expected existing ID to be preserved, got %s", profile.ID)
	}
}

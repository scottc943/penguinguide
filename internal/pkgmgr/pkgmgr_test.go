package pkgmgr

import (
    "reflect"
    "testing"

    "penguinguide/internal/distro"
)

// This test checks that New returns the expected concrete manager type
// for each known distro family and that unknown families get the noop manager.
func TestNewReturnsExpectedManager(t *testing.T) {
    tests := []struct {
        name string
        d    *distro.Distro
        want interface{}
    }{
        {
            name: "debian family gets aptManager",
            d: &distro.Distro{
                Family: distro.FamilyDebian,
                ID:     "debian",
            },
            want: &aptManager{},
        },
        {
            name: "rhel family gets dnfManager",
            d: &distro.Distro{
                Family: distro.FamilyRHEL,
                ID:     "rhel",
            },
            want: &dnfManager{},
        },
        {
            name: "arch family gets pacmanManager",
            d: &distro.Distro{
                Family: distro.FamilyArch,
                ID:     "arch",
            },
            want: &pacmanManager{},
        },
        {
            name: "alpine family gets apkManager",
            d: &distro.Distro{
                Family: distro.FamilyAlpine,
                ID:     "alpine",
            },
            want: &apkManager{},
        },
        {
            name: "unknown family gets noopManager",
            d: &distro.Distro{
                Family: "custom",
                ID:     "mydistro",
            },
            want: &noopManager{},
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            mgr := New(tc.d)
            if mgr == nil {
                t.Fatalf("New returned nil for distro %+v", tc.d)
            }

            gotType := reflect.TypeOf(mgr)
            wantType := reflect.TypeOf(tc.want)

            if gotType != wantType {
                t.Fatalf("New(%+v) type = %v, want %v", tc.d, gotType, wantType)
            }
        })
    }
}


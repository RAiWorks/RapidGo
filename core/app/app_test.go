package app

import (
	"testing"

	"github.com/RAiWorks/RapidGo/core/container"
)

// mockProvider tracks Register/Boot calls for testing.
type mockProvider struct {
	registerCalled bool
	bootCalled     bool
	bootOrder      *[]string
	name           string
	bindName       string
	bindValue      string
}

func (p *mockProvider) Register(c *container.Container) {
	p.registerCalled = true
	if p.bindName != "" {
		c.Instance(p.bindName, p.bindValue)
	}
}

func (p *mockProvider) Boot(c *container.Container) {
	p.bootCalled = true
	if p.bootOrder != nil {
		*p.bootOrder = append(*p.bootOrder, p.name)
	}
}

// TC-13: App.Register calls provider.Register immediately
func TestApp_RegisterCallsProviderRegister(t *testing.T) {
	a := New()
	p := &mockProvider{bindName: "greeting", bindValue: "hello"}

	a.Register(p)

	if !p.registerCalled {
		t.Fatal("Register should be called immediately")
	}

	result := a.Make("greeting").(string)
	if result != "hello" {
		t.Errorf("expected 'hello', got '%s'", result)
	}
}

// TC-14: App.Boot calls all providers' Boot in registration order
func TestApp_BootCallsInOrder(t *testing.T) {
	a := New()
	var order []string

	pA := &mockProvider{name: "A", bootOrder: &order}
	pB := &mockProvider{name: "B", bootOrder: &order}

	a.Register(pA)
	a.Register(pB)
	a.Boot()

	if len(order) != 2 {
		t.Fatalf("expected 2 boot calls, got %d", len(order))
	}
	if order[0] != "A" || order[1] != "B" {
		t.Errorf("expected boot order [A, B], got %v", order)
	}

	if !pA.bootCalled || !pB.bootCalled {
		t.Error("both providers should have Boot called")
	}
}

// TC-15: App.Make resolves service from container
func TestApp_MakeResolvesService(t *testing.T) {
	a := New()
	p := &mockProvider{bindName: "greeting", bindValue: "world"}

	a.Register(p)

	result := a.Make("greeting").(string)
	if result != "world" {
		t.Errorf("expected 'world', got '%s'", result)
	}
}

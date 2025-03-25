package trbac

import "testing"

func TestLiteralContext(t *testing.T) {
    action := "read"
    resourceType := "document"
    roles := []string{"reader", "editor"}

    ctx := NewLiteralContext(action, resourceType, roles)

    if ctx.Action() != action {
        t.Errorf("Expected action %s, got %s", action, ctx.Action())
    }

    if ctx.ResourceType() != resourceType {
        t.Errorf("Expected resource type %s, got %s", resourceType, ctx.ResourceType())
    }

    if len(ctx.Roles()) != len(roles) {
        t.Errorf("Expected roles length %d, got %d", len(roles), len(ctx.Roles()))
    }

    for i, role := range roles {
        if ctx.Roles()[i] != role {
            t.Errorf("Expected role %s, got %s", role, ctx.Roles()[i])
        }
    }
}

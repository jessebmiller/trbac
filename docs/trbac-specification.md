# TRBAC: Typed Role-Based Access Control with Constraints
## Formal Specification v1.0

## 1. Introduction

TRBAC (Typed Role-Based Access Control) is an authorization framework designed to provide a simple, extensible, and language-agnostic approach to role-based access control with powerful constraint capabilities. This specification defines the core concepts, behaviors, and properties that any TRBAC implementation must support, regardless of programming language or execution environment.

### 1.1 Design Goals

1. **Simplicity**: Authorization logic should be clear and understandable
2. **Type Safety**: Resource types should be explicitly defined and checked
3. **Extensibility**: The system should be easily extended via constraints
4. **Language Agnosticism**: Core concepts should be implementable in any language
5. **Minimal Dependencies**: Core functionality should have minimal external dependencies
6. **Declarative Configuration**: Authorization rules should be defined declaratively

## 2. Core Concepts

### 2.1 Entities

The TRBAC model consists of the following core entities:

#### 2.1.1 Resource Types

A **ResourceType** is a named category of resources that can be protected by the system. Resource types define the "what" in authorization decisions.

Properties:
- Resource types are represented as identifiers
- Resource types should follow a consistent naming convention within an implementation

#### 2.1.2 Actions

An **Action** represents an operation that can be performed on a resource. Actions define the "how" in authorization decisions.

Properties:
- Actions are represented as identifiers
- Actions should follow a consistent naming convention within an implementation
- Common actions include but are not limited to: "read", "write", "create", "update", "delete", "list"

#### 2.1.3 Roles

A **Role** represents a named set of permissions that can be assigned to an actor. Roles define the "who" in authorization decisions.

Properties:
- Roles are represented as identifiers
- A single actor may have multiple roles
- Implementations may support role hierarchies (role inheritance)

#### 2.1.4 Constraints

A **Constraint** is a rule that can restrict when a permission applies based on the context of an authorization request.

Properties:
- Constraints are identified by a name
- Constraints evaluate to a boolean value (pass/fail)
- Constraints have access to the full authorization context
- Constraints may access external systems to make decisions

#### 2.1.5 Permissions

A **Permission** associates actions on resource types with optional constraints, defining what operations are allowed.

Properties:
- A collection of allowed actions
- A collection of applicable resource types
- An optional collection of constraints that must all pass for the permission to be granted

#### 2.1.6 Privileges

**Privileges** define the assignment of permissions to roles, establishing which roles have which permissions.

Properties:
- Map roles to their granted permissions
- Define the complete authorization model of the system

### 2.2 Authorization Context

An **Authorization Context** represents all relevant information about an access attempt that is needed to make an authorization decision.

Properties:
- The action being performed
- The resource type being accessed
- The roles assigned to the actor
- May include additional contextual information for constraint evaluation

## 3. Core Behaviors

### 3.1 Authorization Decision

The core authorization behavior is determining whether a given context is authorized. This is the primary function of a TRBAC system.

Authorization logic:
1. Retrieve all permissions granted to the roles in the context
2. For each permission:
   a. Check if the permission is relevant (matches the action and resource type)
   b. If relevant, evaluate all constraints for the permission
   c. If all constraints pass, the action is permitted
3. If no permission grants the action, it is denied

### 3.2 Constraint Evaluation

Constraints are evaluated in the context of an authorization attempt:

1. Each constraint is evaluated with access to the authorization context
2. A constraint returns a boolean result (pass/fail)
3. All constraints associated with a relevant permission must pass for the permission to be granted

### 3.3 Permission Composition

When multiple roles are assigned to an actor, the effective permissions are the union of all permissions granted to any of those roles.

## 4. Required Properties

### 4.1 Type Safety

Resource types and actions must be explicitly checked during authorization decisions to provide type safety:

1. The resource type of the context must match a resource type of the permission
2. The action of the context must match an action of the permission

### 4.2 Extensibility

The constraint mechanism must be extensible to accommodate custom authorization logic:

1. Implementations must support adding custom constraints
2. Constraints must have access to the full authorization context
3. The constraint evaluation mechanism must be pluggable

### 4.3 Determinism

Authorization decisions must be deterministic:

1. Given the same context and privileges, the authorization decision must always be the same
2. Authorization decisions should not depend on global state unless that state is explicitly part of the context

## 5. Abstract Data Model

### 5.1 Representation

While implementations may vary, conceptually the TRBAC model can be represented as:

```
Context := {
    action: Identifier,
    resourceType: Identifier,
    roles: Set<Identifier>
}

Permission := {
    actions: Set<Identifier>,
    resourceTypes: Set<Identifier>,
    constraints: Set<Identifier>
}

Privileges := Map<Role, Set<Permission>>

AuthZ(context: Context, privileges: Privileges) -> Boolean
```

### 5.2 Logical Representation

The authorization decision can be expressed in predicate logic as:

```
AuthZ(context, privileges) := 
  ∃ role ∈ context.roles:
    ∃ permission ∈ privileges[role]:
      context.action ∈ permission.actions ∧
      context.resourceType ∈ permission.resourceTypes ∧
      ∀ constraint ∈ permission.constraints:
        EvaluateConstraint(constraint, context)
```

## 6. Configuration

### 6.1 Declarative Configuration

TRBAC privileges should be configurable in a declarative format. The abstract structure is:

```
roles:
  <role_name>:
    permissions:
      - actions: [<action1>, <action2>, ...]
        resource_types: [<type1>, <type2>, ...]
        constraints: [<constraint1>, <constraint2>, ...]
```

Implementations may use formats like TOML, YAML, JSON, or others as appropriate.

## 7. Extension Points

### 7.1 Hierarchical Roles

Implementations may support role hierarchies where roles can inherit permissions from other roles.

Properties:
- Role inheritance relationships define a directed acyclic graph
- When retrieving permissions for a role, permissions from all ancestor roles are included
- Cycle detection must be implemented to prevent infinite recursion

### 7.2 Resource-Aware Constraints

Implementations may support resource-aware constraints that have access to the specific resource instance being accessed.

Properties:
- The Authorization Context may be extended to provide access to the resource
- Constraints can examine resource properties to make decisions
- This enables attribute-based access control (ABAC) functionality

### 7.3 Dynamic Privileges

Implementations may support dynamic privileges that can change at runtime.

Properties:
- Privileges may be updated during system operation
- Changes to privileges are propagated to all authorization decision points
- Consistency guarantees must be documented

## 8. Security Considerations

Implementations should:

1. Provide protection against constraint evaluation failures
2. Document the behavior when constraints cannot be evaluated
3. Include logging/auditing capabilities for authorization decisions
4. Validate all inputs to prevent injection attacks
5. Consider performance implications for frequent authorization checks

## Appendix A: Examples

### A.1 Example Authorization Scenario

**Privileges definition:**
```
roles:
  reader:
    permissions:
      - actions: [read, list]
        resource_types: [document, folder]
        constraints: []
  
  writer:
    permissions:
      - actions: [read, list, create, update]
        resource_types: [document]
        constraints: [business_hours_only]
  
  admin:
    permissions:
      - actions: [read, list, create, update, delete]
        resource_types: [document, folder, user]
        constraints: [audit_logged]
```

**Authorization scenarios:**
1. User with role "reader" attempts to read a document
   - Outcome: Permitted (matches permission with no constraints)
   
2. User with role "writer" attempts to update a document outside business hours
   - Outcome: Denied (constraint "business_hours_only" fails)
   
3. User with roles "reader" and "writer" attempts to update a document during business hours
   - Outcome: Permitted (writer role grants permission with passing constraint)
   
4. User with no roles attempts to read a document
   - Outcome: Denied (no applicable permissions)

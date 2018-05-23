# trbac
What is typed role-based access control with constraints?

# Goal

* Simple authorization mechanism for microservices
* Declarative RBAC definitions as configuration
* Language agnostic contstraints through shell scripts

## Libraries for declarative Typed RBAC

* Users of the library provide a Context implementation that provides an
  `Action` on a `ResourceType` and the `Roles` of the requesting agent as well
  as any metadata your constraints need.

* Users of the library also provide a `Privileges` definition in TOML.

* auth.May(c Context) bool

# Model

drawing: https://docs.google.com/drawings/d/1MFizaLM7BWfPZSmduzvzDIYvQer8vfR2nGHbrjzVtbM/edit?usp=sharing

* Resource Types: Anything managed by the service being protected
* Resources: a particular instance of a resource type
* Action: Something the service being protected can do with resources (Read,
  Write, List, etc.)
* Permissions: The general right to take an action on resource of a type under
  Constraints
* Roles: Classes of actors on protected resources
* Context: The particular properties of a request to take an action on a
  resource
  * For a proxy, this could be the `*http.Request` object
* Constraint: Arbitrary rule that may deny permission depending on context
* [Future work] Dependant constraint: Constraint that depends on particular
  resources that the request applies to
* Privileges: The assignment of a Permission to a role
  * [Future work] declarative constraints on Privileges, and Privileges state
    changes

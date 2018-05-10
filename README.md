# trbac
What is role-based access control with constraints?

# Goal

## A Reverse Proxy for authorization middleware

* Simple authorization mechanism for microservices
* Declarative RBAC definitions as configuration
* Language agnostic contstraints (or at least Go, Python, and JavaScript)
  * Simple client constraint libraries for ease of language portability
  * Shell scripts as constraints?
    * Language agnostic with the #!
    * Can be managed as configuration

## Libraries for non url based authorization

* may(roles, action, resource, context)

# Model

drawing: https://docs.google.com/drawings/d/1MFizaLM7BWfPZSmduzvzDIYvQer8vfR2nGHbrjzVtbM/edit?usp=sharing

* Resource Types: Anything managed by the service being protected
* Resources: a particular instance of a resource type
* Action: Something the service being protected can do with resources (Read, Write, List, etc.)
* Permissions: The general right to take an action on resource of a type
* Roles: Classes of clients
* Context: The particular properties of a request to take an action on a resource
  * For the proxy, this is the request object
* Constraint: Arbitrary rule that may revoke permission depending on context
* Dependant constraint: Constraint that depends on particular resources that the request applies to
* Privileges: The assignment of a Permission to a role under Constraints

## Open Questions

* How do we inform trbac what roles are relevant?
  * mapping in a configuration file? :\
    * start with this, but make it pluggable so we can write a "Roles" interface
      for whatever storage mechanism we want when we need it
  * Templated URI to user service?
    * would need templated parsing of result


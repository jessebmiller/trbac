# trbac
What are role based access controls with constraints?

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

* Resource Types: Anything managed by the service being protected
* Resources: a particular instance of a resource type
* Action: Something the service being protected can do with resources (Read, Write, List, etc.)
* Permissions: The general right to take an action on resource of a type
* Roles: Classes of clients
* Context: The particular properties of a request to take an action on a resource
* Constraint: Arbitrary rule that may revoke permission depending on context
* Dependant constraint: Constraint that depends on particular resources that the request applies to
* Privileges: The assignment of a Permission to a role under Constraints


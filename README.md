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
* constraint(context, resource)

# Model

* Resource Types: A class of resource relevant to the service being protected
* Resources: a particular instance of a resource type
* Action: Something users can do with resources (Read, Write, List, etc.)
* Permissions: The general right to take an action on a type of resource
* Roles: Classes of users
* Context: The particular properties of a request
* Constraint: Arbitrary rule that may revoke permission depending on context
* Dependant constraint: Constraint that depends on the particular resource
* Privileges: The assignment of a Permission to a role under Constraints


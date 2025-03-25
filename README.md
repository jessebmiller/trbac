# TRBAC

What is Typed Role Based Access Control with constraints?

# Design goals

- **Simplicity**: Authorization logic should be clear and understandable
- **Type Safety**: Resource types should be explicitly defined and checked
- **Extensibility**: The system should be easily extended via constraints
- **Language Agnosticism**: Core concepts should be implementable in any language
- **Minimal Dependencies**: Core functionality should have minimal external dependencies
- **Declarative Configuration**: Authorization rules should be defined declaratively

## Target Use Cases

### Edge Computing & IoT Environments

Edge devices often have limited resources and intermittent connectivity. TRBAC's lightweight nature and minimal dependencies make it ideal for:

- Embedded systems with constrained environments
- Devices that need to make authorization decisions offline
- IoT gateways that manage access to multiple connected devices

Key features: Lightweight persistence layer optimized for embedded systems; offline caching with synchronization mechanisms.

### DevOps & Infrastructure Tooling

The shell script constraint runner is unique and particularly valuable for infrastructure-focused tools where:

- Authorization needs to integrate with existing shell scripts and system commands
- Permissions depend on infrastructure state that's best queried through CLI tools
- Operators are already comfortable with shell scripting

Key features: Built-in integration with common DevOps tools; template library for infrastructure-specific constraints.

### Multi-Language Microservice Architectures

Your design is conceptually simple and the shell script constraint runner provides language-agnostic extensibility:

- Services written in different languages can share the same authorization model
- Shell script constraints can be written in any language and called consistently
- Authorization logic can be centralized while implementation remains distributed

Key features: Standardized constraint interfaces for multiple languages; gRPC/REST services for centralized policy management.

### Legacy System Integration

Many organizations struggle with access control when integrating with legacy systems:

- Shell script constraints could easily wrap existing CLI-based authentication systems
- The simple model makes it easier to map between modern and legacy paradigms
- Lightweight implementation means it can be deployed alongside legacy components

Key features: Adapters for common legacy authorization systems; migration tooling.

### Educational & Teaching Environments

TRBAC's clarity makes it excellent for teaching authorization concepts:

- Clean separation between model components
- Clear implementation that follows the conceptual model
- Minimal "magic" compared to larger frameworks

Key features: Interactive examples; visualization of authorization decisions; configurable complexity levels.

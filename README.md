# Diplomacy API

```mermaid
classDiagram
    class User {
        +String name
        +String email
        +login()
    }
    class Admin {
        +manageUsers()
    }
    User <|-- Admin
```

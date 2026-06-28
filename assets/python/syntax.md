---
languages: 
    - python
file_types:
    - *.py
priority: required
related_resources:
    - python/logging
    - python/architecture
---

# Python Syntax Standards
Standards for how python code should look. Naming conventions, language structure preferences


- ALWAYS favor legibility over performance or speed. ALWAYS choose common, popular and well-documented implementation strategies. NEVER write novel or unique implementations.


## Conditional Operations

- NEVER use ternary operators - ex:
```python
status = "Adult" if age > 18 else "Minor"
```

## Class Naming Conventions

- ALWAYS Postfix classes with appropriate pattern names that match the role of the class - Factory, Singleton, Service, Repository. 
    - Example: If there should only ever be one "Settings" class instance - it should be called "SettingsSingleton"
    - Example: If a class handles business logic for registering users - it should be called "RegistrationService"

- ALWAYS annotate DTO classes with the @dataclass annotation
- IF there is a clear inbound and outbound contract, PREFER DTO classes postfixed with "Request" and "Response", not "DTO".
Example:
```python
@dataclass
class CreateUserRequest()
    username: str
    password: str

@dataclass
class CreateUserResponse()
    username: str
    account_active: bool
```

## Interfaces

- ALWAYS postfix interface classes with "Interface"
- ALWAYS extend abc.ABC in interfaces
- ALWAYS use @abstractmethod annotation for abstract methods
- NEVER write implementation in an interface. Interfaces ONLY contain method signatures

Example:
```python
from abc import ABC, abstractmethod

class CarInterface(ABC)

    @abstractmethod
    def start_car():
        ...

    @abstractmethod
    def pop_hood():
        ...
```

# Abstract Classes

- ALWAYS Prefix class name with "Abstract"
- ALWAYS implement abc.ABC for abstract classes
- ALWAYS use @abstractmethod annotation for abstract methods within the class

Example:
```python
from abc import ABC, abstractmethod

class AbstractTruck(ABC):

    def __init__(self):
        self._logger = logging.logger(__class__.__name__)

    @abstractmethod
    def start_truck():
        ...
```
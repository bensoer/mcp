---
uri: standards://python/syntax
name: Python Syntax Standards
description: "Standards for how python code should look. Naming conventions, language structure preferences"
languages: 
    - python
file_types:
    - "*.py"
priority: required
related_resources:
    - standards://python/logging
    - standards://python/architecture
---

# Python Syntax Standards
Standards for how python code should look. Naming conventions, language structure preferences


- ALWAYS favor legibility over performance or speed. ALWAYS choose common, popular and well-documented implementation strategies. NEVER write novel or unique implementations.


## Conditional Operations

- NEVER use ternary operators. ALWAYS use full "if", "elif", "else" syntax

Example:
```python
# BAD:
status = "Adult" if age > 18 else "Minor"

# GOOD:
if age > 18:
    status = "Adult"
else:
    status = "Minor
```

- ALWAYS explicitly check for None. NEVER check for None via "truthiness"

Example:
```python:
results: Optional[ResultResponse] = None

# BAD
if results:
    raise Exception("Results Not Found")

# BAD
if not results:
    return results.status

# GOOD
if results is None:
    raise Exception("Results Not Found")

# GOOD
if results is not None:
    return results.status
```

## Typing

- PREFER explicite typing as much as possible
- ALWAYS specify types in method parameters
- ALWAYS use typing module for Optional, Union, Any types. NEVER use pipe (|) operator syntax
- ALWAYS specify return type for methods INCLUDING None for void return types

Example:
```python
from typing import Optional, Any

# GOOD
username: str = "foxtrot94"
metadata: Optional[dict[string,Any]] = None

# BAD
username = "foxtrot94" # Implicite type - UNDESIRABLE
metadata: dict | None # Pipe syntax - BAD . Implcite dict types - UNDESIRABLE

# GOOD - username explicite type, metadata explitely typed. Returned type, even for void
def create_user(username: str, metadata: Optional[dict[string,Any]]) -> None:
    ...

# BAD - username missing type, metadata using pipe syntax. Missing return type
def create_user(username, metadata: dict | None)
    ...

```

## Parameter Passing & Method Signatures

- PREFER non-optional types

Examples:
```python

# GOOD
def create_user(username: str, password: str, email: str) -> None:
    ...
# NOT IDEAL - default variable arguments used
def create_user(username: str, password: str, email = "default@gmail.com") -> None:
    ...

```

- Use DTO objects when parameter number exceeds 5 parameters
- NEVER return tuples in methods. Return DTOs when multiple variables are required in response

Examples:
```python
# GOOD - Passes DTO, and returns DTO when parameter is more then 5 and response is a complex object
def create_user(request: CreateUserRequest) -> CreateUserResponse:
    ...
# BAD - Too many parameters & returned Tuple
def create_user(username: str, password: str, email: str, address: str, phone_number: str, active: bool) -> (str, str, str,):
    ...


```

- NEVER use variable arguments - `*args` or `*kwargs` method signatures


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
- ALWAYS extend abc.ABC in interfaces. NEVER use Protocol
- ALWAYS use @abstractmethod annotation for abstract methods
- NEVER write implementation in an interface. Interfaces ONLY contain method signatures

Example:
```python
from abc import ABC, abstractmethod

class CarInterface(ABC)

    @abstractmethod
    def start_car() -> None:
        ...

    @abstractmethod
    def pop_hood() -> None:
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
    def start_truck() -> None:
        ...
```
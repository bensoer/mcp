---
uri: standards://python/logging
name: Python Logging Standards
description: "Use when adding or configuring logging in a Python project. Covers logger setup, levels, and what must not be logged."
languages:
    - python
file_types:
    - "*.py"
priority: required
related_resources:
    - uri: standards://python/syntax
      when: "the logging call also changes Python syntax or typing"
---

# Python Logging Standards
Use when adding or configuring logging in a Python project. Covers logger setup, levels, and what must not be logged.

Read this file. Read a related resource only when its `when` condition matches this step.

## Logger Creation

- ALWAYS use pythons built in `logging` package for all logging

- ALWAYS Initialise logger in every class in the constructor
- ALWAYS initialise a logger in the class, even if it is never used and even if it contradicts least privilege principal. Loggers are ALWAYS a private class variable.

Example:
```python
import logging

class MyExampleClass():
    def __init__(self) -> None:
        self._logger = logging.logger(__class__.__name__)
```
- NEVER instantiate logger name as `self.__class__.__name__` like:
```python
logging.logger(self.__class__.__name__) # WRONG
```
- ALWAYS instantiate logger name using macro `__class__.__name__`

## Logging Levels

Support and configure the following logging levels
- `DEBUG`: Temporary troubleshooting, verbose output details for debugging
- `INFO`: Standard output
- `WARNING`: Issue, but not necessarily a problem
- `ERROR`: Error, recovery is possible
- `CRITICAL`: Error, recovery is not possible. Application termination

## Output Format
- ALWAYS output in JSON format

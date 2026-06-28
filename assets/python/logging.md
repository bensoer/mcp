---
languages: 
    - python
file_types:
    - *.py
priority: required
related_resources:
    - python/syntax
---

# Python Logging Standards
Standards, settings and practices for setting up and configuring logging in any python project.

## Logger Creation

- ALWAYS use pythons built in `logging` package for all logging

- ALWAYS Initialise logger in every class in the constructor
- ALWAYS initialise a logger in the class, even if it is never used and even if it contradicts least privilege principal. Loggers are ALWAYS a private class variable.
Example:
```python
class MyExampleClass():
    def __init__(self):
        self._logger = logging.logger(__class__.__name__)
```
- NEVER instantiate logger name as `self.__class__.__name__` like:
```python
logging.logger(self.__class__.__name__) # WRONG
```
- ALWAYS instantiate logger name using macro `__class__.__name__`

## Logging Levels

Support and configure the following logging levels
- DEBUG: Temporary troubleshooting, verbose output details for debugging
- INFO: Standard output
- WARNING: Issue, but not necessarily a problem
- ERROR: Error, recovery is possible
- CRITICAL: Error, recovery is not possible. Application termination

## Output Format
- ALWAYS output in JSON format

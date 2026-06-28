---
languages: 
    - python
file_types:
    - *.py
priority: required
related_resources:
    - python/syntax
---

# Python Coding General Practices
Best practices for code organisation and application architecture

- PREFER Composition over Inheritance
- IF choosing Interhitence, implementation MUST pass ALL S.O.L.I.D. principals. Failure of ANY SOLID principal means Inheritance is the wrong architecture. Use Composition instead.
- Keep code D.R.Y.
- NEVER create a method unless code within it is being reused in 2 or more places. ALWAYS place code used only once inline.
- NEVER create package level methods. ALWAYS place methods closest to the code that use it. Create private methods first within the class using it. If that does not work, Create a new class with static methods to access.
- NEVER create a `utils.py` or like package and place methods within it
- NEVER "poach" methods from other classes. Poaching breaks encapsulation and single responsibility principals. Poaching also signals bad architecture design and needs refactoring.
- AVOID novel architectures or solutions. PREFER mainstream, common or "textbook" implementations as closely as possible
- PREFER Dependency Injection and Inversion-of-Control patterns
- PREFER Low Coupling and High Cohesion
- Classes ALWAYS follow Encapsulation and Single Responsibility Principals
- NEVER access class attributes directly UNLESS it is a DTO or @dataclass. 
- ALWAYS implement "setters" and "getters" for variables if they are need to be accessed outside of the class.

- ALWAYS implement code following Least Privelege Principal - if it can be private, make it private
    - NEVER create "setters" and "getters" for variables that are not accessed by 1 or more classes or code outside the class
    - NEVER create class variables if they are not used outside of a single method within the class
    - NEVER create public variables in a class - create "setters" and "getters" methods.
    - NOTE: Loggers are an EXCEPTION. Loggers are always class variables and should always exist even if not used
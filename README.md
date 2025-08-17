## GOALS
- Provide a robust configuration loading system for Go applications.
- Support loading configuration from various sources: files (YAML/JSON) and environment variables.
- Allow easy mapping of configuration fields to struct fields using tags.
- Facilitate the dumping of configuration data in a user-friendly format.

## NON-GOALS
- Provide a complete solution for all possible configuration scenarios. Instead this should be a robust starting point that can be extended as needed and just adds a little bit of structure and organization on top of YAML/JSON/.. file loading.
- Replace existing configuration management tools or libraries. This library is not intended to be a drop-in replacement for existing tools, but rather a complementary solution specifically with a focus on simplicity and ease of use.
- Verification of configuration values. There are other frameworks that specialize in struct validation f.e. https://github.com/go-playground/validator
- Performance and Memory optimization. A configuration is typically loaded once at startup and then accessed in-memory with a small set of keys/values, so the performance impact is minimal.

# typed: strict

# Sorbet does not support recursive type aliases, so nested values stay untyped.
JsonHash = T.type_alias { T::Hash[String, T.untyped] }

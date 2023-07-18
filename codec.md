- Serialzer/Deserializer:
	- analogs of serde::{Serialize,Deserialize}
	- Serialize accepts a Serializer
	- Deserialize accepts a Deserializer
- Encoder/Decoder:
	- analogs of serde::{Serializer,Deserializer}
	- Decoder uses a Visitor to drive deserialization
- Visitor:
	- analog of serde::de::Visitor

- Open questions:
    - how to allow formats to special-case based on destination type?
        - formats only see the visitor, not the actual destination
        - could add a Type() method or something on visitor s.t. the format can see the type
    - how to allow formats to special-case based on struct fields?
        - this seems clearer: first-class struct decoding/encoding?

- Serialize(enc Encoder)
    - receiver drives the encoder
    - primitives call the encode methods directly
    - composites call the composite methods
    - pointers are an exception to this rule at the moment

---



---

- Type _must_ implement Visitor / Deserializer?
    - idiomatic Go puts serialization behavior in the hands of the destination type
    - wrapping types to provide alternate serialization is problematic b/c it requires recursively wrapping everything
    - no associated types...
    - could use some sort of registration scheme? format.Register, etc.
    - as long as there is always a Decoder / Encoder for composites, this could just work?
        - Decoder / Encoder can check for format-specific bits

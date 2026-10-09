# MRA adapter and source

`data.json` contains the complete class/property/schema objects (57 class records), definitions and metadata from official **MRA English 1.3.0, Release R rev.2**, published 2024-07-26. It combines the original separate files into one embedded JSON document; meanings and releases are preserved. `COPYRIGHT.txt` is the original MIT license/notice. The adapter itself is under the repository license.

- Official reference-data page: https://echonet.jp/spec_mra_rr2_en/
- Original ZIP: https://echonet.jp/wp/wp-content/uploads/pdf/General/Standard/MRA/MRA_en_v1.3.0.zip
- ZIP SHA256: `db8ebf5fe33027255cba4f9da7f7f747ade4a1d65783a9e7e1cb4a95f7d65780`

The vendored records and definitions were compared with that official ZIP. This is a pinned reference dataset, not a claim of current specification compliance or certification. It matches the existing core database generation's MRA version; the core public database does not expose the data schemas needed by an editor. No dependency/core API change was made.

Actual validated property maps govern availability. EPC 82's fresh Get supplies the release when present. Without a release, a definition is usable only if its schema/access is identical and covers releases A through R. A future release or differing historical definition is raw/read only. A device value outside the schema is raw/unknown; no enum/range/unit is inferred from the bytes.

Supported editor schemas: exact `state` EDT choices; bounded uint8/int8/uint16/int16/uint32/int32 numbers with optional fixed `multiple`; length-constrained raw; simple oneOf unions of those. Coefficient-dependent values, numeric enumerations, multipleOf constraints, level/bitmap/array/object/time/date/numericValue and unsupported references are displayed as raw and cannot be edited. Raw values have only size validation, not semantic safety guarantees. Only a writable MRA definition in both the device's Get and Set maps can be written, because a fresh Get is required for verification.

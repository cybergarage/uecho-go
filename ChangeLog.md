# ChangeLog

## v1.3.1 (2026-02-02)

### Fixed
- Prevent panic when closing UDP sockets concurrently ([c652da7](https://github.com/cybergarage/uecho-go/commit/c652da7))

### Improved
- Add `Message.EHD()` to return the Echonet Header (EHD) ([06d94bd](https://github.com/cybergarage/uecho-go/commit/06d94bd))
- Add `ESV::String()` and `ObjectCode::String()` ([0ce2337](https://github.com/cybergarage/uecho-go/commit/0ce2337), [6ea9112](https://github.com/cybergarage/uecho-go/commit/6ea9112))
- Reduced allocations by preallocating slices/buffers in core ECHONET paths ([01c41b7](https://github.com/cybergarage/uecho-go/commit/01c41b7))

## v1.3.0 (2025-10-12)

### Improved
- New abstraction interfaces: `Config`, `Manufacture`, `StandardDatabase`, `DeviceOption` ([9fec903](https://github.com/cybergarage/uecho-go/commit/9fec903edc3457fdd599751295a0c8b046256946), [a8f5a57](https://github.com/cybergarage/uecho-go/commit/a8f5a579700e7e8e80304d44e824f38439e98339), [d02fbb07](https://github.com/cybergarage/uecho-go/commit/d02fbb07ef57b0e9954890462de90c31c036c2cd), [81937072](https://github.com/cybergarage/uecho-go/commit/819370726e9d9fea54a31c2e77b5a7523edfd71d))
- `ObjectRequestHandler` interface ([f67a303](https://github.com/cybergarage/uecho-go/commit/f67a3037048e054ff9042fa8430d66006380fbd2))
- Socket option handling robustness ([46f851f](https://github.com/cybergarage/uecho-go/commit/46f851fb9afd89406c9fc87726d8a06e28e75039))

### Breaking Changes
- Listener methods unified: `NodeMessageReceived` / `PropertyRequestReceived` → `OnMessage` / `OnPropertyRequest`; controller events simplified ([05b01d42](https://github.com/cybergarage/uecho-go/commit/05b01d421138fc32016a597a7483271c6a0a4778), [f86fe82](https://github.com/cybergarage/uecho-go/commit/f86fe82a1138a42cff5003fd2751df37bb3d24c8))
- Property helpers renamed: `ByteData` / `StringData` / `IntegerData` → `AsByte` / `AsString` / `AsInteger` ([f86fe82](https://github.com/cybergarage/uecho-go/commit/f86fe82a1138a42cff5003fd2751df37bb3d24c8))

### Data Updates
- MRA objects refreshed to latest spec ([87475f2](https://github.com/cybergarage/uecho-go/commit/87475f2f4ad86adfe4ab625f390c3d87c856963c))
- Manufacturer codes updated ([08532d09](https://github.com/cybergarage/uecho-go/commit/08532d09f23167fbbc6ab29b28e66bf7fd8792f7))

## v1.2.2 (2024-08-05)
- Updated the standard object database based on the latest MRA (Machine Readable Appendix) version 1.3.0 from the ECHONET Consortium.
- Updated the standard manufacturer code database based on the latest MCA (Manufacturer Code List) from the ECHONET Consortium.

## v1.2.1 (2024-01-16)
- Fixed Update transport servers to exit normally when the socket encounters any errors
- Enabled ubuntu-latest testing on Github actions

## v1.2.0 (2023-07-06)
- Added support for Windows platforms (Thanks to David González Filoso).
- Updated the codebase to be compatible with golang 1.20.

## v1.1.3 (2023-05-07)
- Updated the standard manufacturer code database based on the latest MCA (Manufacturer Code List) from the ECHONET Consortium.
- Updated the standard object database based on the latest MRA (Machine Readable Appendix) version 1.2.0 from the ECHONET Consortium.

## v1.1.2 (2023-01-05)
- Updated to enable self messages.
- Updated StandardDatabase to store property types and capacities more accurately.
- Added uechobench for benchmarking.

## v1.1.1 (2022-11-27)
- Updated Controller to add the standard node profile object as the default.
- Added experimental manufacturer codes to the standard database.

## v1.1.0 (2022-08-26)
- Fixed the issue with setting correct description format 2 property maps.
- Updated the standard object database based on the latest MRA (Machine Readable Appendix) version 1.1.1 from the ECHONET Consortium.
- Removed the 'Get' prefix from all getter methods.
- Extended Property::*Attribute() to set the standard properties more accurately.

## v1.0.2 (2022-08-20)
- Updated the transport layer to bind all available interfaces by default.

## v1.0.1 (2022-08-05)
- Updated the codebase to be compatible with golang 1.18.
- Fixed golangci-lint warnings.

## v1.0.0 (2019-02-04)
- The first public release.

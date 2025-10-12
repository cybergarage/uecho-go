## [v1.3.0] - 2025-10-12

Changes since v1.2.2 (based on the latest ~30 fetched commits). Because of an API pagination/quantity limit this list may be incomplete. Before the final release, please verify with `git log v1.2.2..` or the GitHub compare view.

### 💥 Breaking Changes

Backward-incompatible interface / method name changes and public surface adjustments that may break builds:

1. StandardDatabase  
   - `LookupObjectByCode()` → `LookupObject()` ([b821e66](https://github.com/cybergarage/uecho-go/commit/b821e66ae14b1be6381d3e82aa84bb50e1a5570d))

2. Listener interface naming unification ([f86fe82](https://github.com/cybergarage/uecho-go/commit/f86fe82a1138a42cff5003fd2751df37bb3d24c8), [05b01d42](https://github.com/cybergarage/uecho-go/commit/05b01d421138fc32016a597a7483271c6a0a4778))  
   - `NodeListener.NodeMessageReceived` → `NodeListener.OnMessage`  
   - `ObjectListener.PropertyRequestReceived` → `ObjectListener.OnPropertyRequest`  
   - `ControllerListener.ControllerMessageReceived` → `ControllerListener.OnMessage`  
   - `ControllerListener.ControllerNewNodeFound` → `ControllerListener.OnNodeFound`

3. Property helper methods renamed ([f86fe82](https://github.com/cybergarage/uecho-go/commit/f86fe82a1138a42cff5003fd2751df37bb3d24c8))  
   - `ByteData()` → `AsByte()`  
   - `StringData()` → `AsString()`  
   - `IntegerData()` → `AsInteger()`  

4. Error model reworked ([f86fe82](https://github.com/cybergarage/uecho-go/commit/f86fe82a1138a42cff5003fd2751df37bb3d24c8), [44e52a9](https://github.com/cybergarage/uecho-go/commit/44e52a92b07cc04301cb362a8d6ee784b59a07e7))  
   - Introduced exported sentinel errors (`ErrInvalid`, `ErrNoData`, `ErrNotFound`, etc.) and changed formatted message patterns.  
   - String-based error comparisons may now fail; migrate to `errors.Is` / `errors.As`.

5. Object / handler interface name simplification ([262b907](https://github.com/cybergarage/uecho-go/commit/262b907c9ec4f7fb41b1b899d32c673be640de34), [2cf52977](https://github.com/cybergarage/uecho-go/commit/2cf529770ef896ceb390c6b955a8a2633fb8f9aa), others)  
   - Simplified “object handlerr interface method names” (see commit [262b907](https://github.com/cybergarage/uecho-go/commit/262b907c9ec4f7fb41b1b899d32c673be640de34)).  
   - `PropertyMap.ParentObject()` → `PropertyMap.Object()` ([2cf52977](https://github.com/cybergarage/uecho-go/commit/2cf529770ef896ceb390c6b955a8a2633fb8f9aa))

6. Function privatization ([4a9d6ed](https://github.com/cybergarage/uecho-go/commit/4a9d6ed40fb237188513317c54f3700e28ea2297))  
   - `NewObjectWithCodeBytes()` made private. Replace with existing public factory APIs.

### ✨ Features

- Added `Message::TID()` accessor ([c15e2f5](https://github.com/cybergarage/uecho-go/commit/c15e2f5cb9b98cd3556edb7c2a44f2d446349c0c))  
- Added `ObjectRequestHandler` interface for property request handling ([f67a303](https://github.com/cybergarage/uecho-go/commit/f67a3037048e054ff9042fa8430d66006380fbd2))

### 🛠 Refactoring / Internal Improvements

- Introduced `Config` interface ([9fec903](https://github.com/cybergarage/uecho-go/commit/9fec903edc3457fdd599751295a0c8b046256946))  
- Introduced `Manufacture` interface ([a8f5a57](https://github.com/cybergarage/uecho-go/commit/a8f5a579700e7e8e80304d44e824f38439e98339))  
- Introduced `StandardDatabase` interface ([d02fbb07](https://github.com/cybergarage/uecho-go/commit/d02fbb07ef57b0e9954890462de90c31c036c2cd))  
- Introduced `DeviceOption` interface ([81937072](https://github.com/cybergarage/uecho-go/commit/819370726e9d9fea54a31c2e77b5a7523edfd71d))  
- Socket option handling refactored: use `syscall.RawConn` instead of `os.File` (improved reuse/portability) ([46f851f](https://github.com/cybergarage/uecho-go/commit/46f851fb9afd89406c9fc87726d8a06e28e75039))  
- Parser & protocol layer error normalization and formatting ([f86fe82](https://github.com/cybergarage/uecho-go/commit/f86fe82a1138a42cff5003fd2751df37bb3d24c8), [75fba86](https://github.com/cybergarage/uecho-go/commit/75fba86c433f7cd113860ee8069ece78513f05a2), [618fc52](https://github.com/cybergarage/uecho-go/commit/618fc529c312ce2e1f4b8cde79dd7ff2a2ffcbc4), [741cbc7](https://github.com/cybergarage/uecho-go/commit/741cbc78a3ef5e0d4c8b195423b9c752a8b13662), [dcaf06a3](https://github.com/cybergarage/uecho-go/commit/dcaf06a3bf45e3ec26f8836466186e3b1bba0cfb), [60bbfa0](https://github.com/cybergarage/uecho-go/commit/60bbfa063be1723a8fc48f26bc014efc15c297ac))  
- Unified property / protocol error handling ([f86fe82](https://github.com/cybergarage/uecho-go/commit/f86fe82a1138a42cff5003fd2751df37bb3d24c8), [44e52a9](https://github.com/cybergarage/uecho-go/commit/44e52a92b07cc04301cb362a8d6ee784b59a07e7))

### 🐛 Bug Fixes

- Fixed `TestLocalNode()` ([e5a75ae](https://github.com/cybergarage/uecho-go/commit/e5a75ae276a670e123c87ddfee7c8dcf0877237a))  
- `NewLocalNode()` now forwards `LocalNodeOption` parameters correctly ([e5d6b87](https://github.com/cybergarage/uecho-go/commit/e5d6b877956e8d2d3a22d72ce6af4090a8b44c48))

### 📚 Documentation

- Expanded examples and documentation updates ([683f1f9](https://github.com/cybergarage/uecho-go/commit/683f1f983629c2957b497c0f1d11db04672e9552), [c8748a7](https://github.com/cybergarage/uecho-go/commit/c8748a7ca151e0e54debf8aecf987e6160a8d437), [31572b2](https://github.com/cybergarage/uecho-go/commit/31572b2affdfe38c5149fcc5949d77ba037a6e54), [f6abaf43](https://github.com/cybergarage/uecho-go/commit/f6abaf43654578112dc8f5f22ac94af45550e217), [5eb5d28](https://github.com/cybergarage/uecho-go/commit/5eb5d28a93944da074b9dd263d7b2557f167f8ef))  
- Updated docs for listener renames ([f6abaf43](https://github.com/cybergarage/uecho-go/commit/f6abaf43654578112dc8f5f22ac94af45550e217), others)

### 🧪 Tests

- Updated `TestController()` and `TestLocalNode()` ([b3e7f22](https://github.com/cybergarage/uecho-go/commit/b3e7f229bb566dcdd2e97caefbeca87ad3252c8b))  
- Added/extended example-based tests ([683f1f9](https://github.com/cybergarage/uecho-go/commit/683f1f983629c2957b497c0f1d11db04672e9552))

### 🧩 Data / Database Updates

- Updated MRA objects to latest version (MRA_en_v1.3.0) ([87475f2](https://github.com/cybergarage/uecho-go/commit/87475f2f4ad86adfe4ab625f390c3d87c856963c))  
- Updated manufacturer codes to the latest list ([08532d09](https://github.com/cybergarage/uecho-go/commit/08532d09f23167fbbc6ab29b28e66bf7fd8792f7))  
- Large object definition diff (1,155 additions / 1,155 deletions) ([08532d09](https://github.com/cybergarage/uecho-go/commit/08532d09f23167fbbc6ab29b28e66bf7fd8792f7))

### 🧹 Chore / Maintenance / CI

- Updated `.golangci.yaml` ([1672ab87](https://github.com/cybergarage/uecho-go/commit/1672ab8793fa738c4d34935fb7f7f8c15bd769b9), [13481d7](https://github.com/cybergarage/uecho-go/commit/13481d7b807d503abcd7cd2de33b4c913385b9ea))  
- GitHub Actions / make workflow tweak ([6f38926](https://github.com/cybergarage/uecho-go/commit/6f38926cdbed5f53cea5ab9f59e7852cb090c77a))  
- README update ([5eb5d28](https://github.com/cybergarage/uecho-go/commit/5eb5d28a93944da074b9dd263d7b2557f167f8ef))  
- Function comment / code comment grooming ([5e7ca9ca](https://github.com/cybergarage/uecho-go/commit/5e7ca9ca59e1fcbcb665d283f7059fce1701aa2c))  
- Script / generator adjustments for objects ([08532d09](https://github.com/cybergarage/uecho-go/commit/08532d09f23167fbbc6ab29b28e66bf7fd8792f7))  
- Ignore Docker interfaces in `GetAvailableInterfaces()` ([d86f0ad](https://github.com/cybergarage/uecho-go/commit/d86f0ad058db42c37583ec2de89a261bb591340b))

### 🔧 Migration Guide (Summary of Required Actions)

| Change | Old | New | Action |
|--------|-----|-----|--------|
| Standard DB lookup | `LookupObjectByCode()` | `LookupObject()` | Rename usages ([b821e66](https://github.com/cybergarage/uecho-go/commit/b821e66ae14b1be6381d3e82aa84bb50e1a5570d)) |
| Property helpers | `ByteData()/StringData()/IntegerData()` | `AsByte()/AsString()/AsInteger()` | Rename + add error handling (`ErrNoData`, etc.) ([f86fe82](https://github.com/cybergarage/uecho-go/commit/f86fe82a1138a42cff5003fd2751df37bb3d24c8)) |
| Listener methods | `NodeMessageReceived`, etc. | `OnMessage`, `OnPropertyRequest`, `OnNodeFound` | Rename implementations ([05b01d42](https://github.com/cybergarage/uecho-go/commit/05b01d421138fc32016a597a7483271c6a0a4778)) |
| PropertyMap parent | `ParentObject()` | `Object()` | Rename calls ([2cf52977](https://github.com/cybergarage/uecho-go/commit/2cf529770ef896ceb390c6b955a8a2633fb8f9aa)) |
| Privatized constructor | `NewObjectWithCodeBytes()` | (private) | Use public object factory APIs ([4a9d6ed](https://github.com/cybergarage/uecho-go/commit/4a9d6ed40fb237188513317c54f3700e28ea2297)) |
| Error handling | String compare | `errors.Is` / `errors.As` | Refactor tests & logic ([f86fe82](https://github.com/cybergarage/uecho-go/commit/f86fe82a1138a42cff5003fd2751df37bb3d24c8)) |

### 🔗 Key Commits Referenced

[c15e2f5](https://github.com/cybergarage/uecho-go/commit/c15e2f5cb9b98cd3556edb7c2a44f2d446349c0c), [f67a303](https://github.com/cybergarage/uecho-go/commit/f67a3037048e054ff9042fa8430d66006380fbd2), [b821e66](https://github.com/cybergarage/uecho-go/commit/b821e66ae14b1be6381d3e82aa84bb50e1a5570d), [f86fe82](https://github.com/cybergarage/uecho-go/commit/f86fe82a1138a42cff5003fd2751df37bb3d24c8), [05b01d42](https://github.com/cybergarage/uecho-go/commit/05b01d421138fc32016a597a7483271c6a0a4778), [262b907](https://github.com/cybergarage/uecho-go/commit/262b907c9ec4f7fb41b1b899d32c673be640de34), [2cf52977](https://github.com/cybergarage/uecho-go/commit/2cf529770ef896ceb390c6b955a8a2633fb8f9aa), [46f851f](https://github.com/cybergarage/uecho-go/commit/46f851fb9afd89406c9fc87726d8a06e28e75039), [44e52a9](https://github.com/cybergarage/uecho-go/commit/44e52a92b07cc04301cb362a8d6ee784b59a07e7), [e5d6b87](https://github.com/cybergarage/uecho-go/commit/e5d6b877956e8d2d3a22d72ce6af4090a8b44c48), [e5a75ae](https://github.com/cybergarage/uecho-go/commit/e5a75ae276a670e123c87ddfee7c8dcf0877237a), [08532d09](https://github.com/cybergarage/uecho-go/commit/08532d09f23167fbbc6ab29b28e66bf7fd8792f7), [87475f2](https://github.com/cybergarage/uecho-go/commit/87475f2f4ad86adfe4ab625f390c3d87c856963c)

### ✅ Recommended Upgrade Steps

1. Update dependency: `go get github.com/cybergarage/uecho-go@v1.3.0`  
2. Resolve build errors from renamed methods (listeners, property helpers, lookup functions).  
3. Migrate error comparisons to `errors.Is` / `errors.As`.  
4. Review generated / object definition dependent code for changes (MRA updates).  
5. Adopt new interfaces (`ObjectRequestHandler`, `Message::TID()`) where beneficial.  

### 📌 Note on Completeness

Because only a limited window of commits was fetched through the API, verify full history with:
```
git log --oneline v1.2.2..HEAD
```
or the GitHub compare view.

---

[Compare v1.2.2...v1.3.0](https://github.com/cybergarage/uecho-go/compare/v1.2.2...v1.3.0)

# 2024-08-05 v1.2.2
* Updated the standard object database based on the latest MRA (Machine Readable Appendix) version 1.3.0 from the ECHONET Consortium.
* Updated the standard manufacturer code database based on the latest MCA (Manufacturer Code List) from the ECHONET Consortium.

# 2024-01-16 v1.2.1
* Fixed Update transport servers to exit normally when the socket encounters any errors
* Enabled ubuntu-latest testing on Github actions

# 2023-07-06 v1.2.0
* Added support for Windows platforms (Thanks to David González Filoso).
* Updated the codebase to be compatible with golang 1.20.

# 2023-05-07 v1.1.3
* Updated the standard manufacturer code database based on the latest MCA (Manufacturer Code List) from the ECHONET Consortium.
* Updated the standard object database based on the latest MRA (Machine Readable Appendix) version 1.2.0 from the ECHONET Consortium.

# 2023-01-05 v1.1.2
* Updated to enable self messages.
* Updated StandardDatabase to store property types and capacities more accurately.
* Added uechobench for benchmarking.

# 2022-11-27 v1.1.1
* Updated Controller to add the standard node profile object as the default.
* Added experimental manufacturer codes to the standard database.

# 2022-08-26 v1.1.0
* Fixed the issue with setting correct description format 2 property maps.
* Updated the standard object database based on the latest MRA (Machine Readable Appendix) version 1.1.1 from the ECHONET Consortium.
* Removed the 'Get' prefix from all getter methods.
* Extended Property::*Attribute() to set the standard properties more accurately.

# 2022-08-20 v1.0.2
* Updated the transport layer to bind all available interfaces by default.

# 2022-08-05 v1.0.1
* Updated the codebase to be compatible with golang 1.18.
* Fixed golangci-lint warnings.

# 2019-02-04 v1.0.0
* The first public release.

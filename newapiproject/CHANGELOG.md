# Changelog

## [1.0.0] - 2024-05-27
### Added
- Added Changelog.md file

### Changed
- Switched to Linux environment
- Ubuntu environment set up
- Docker and Go installed on Ubuntu
- Project set up on Ubuntu

### Fixed
- Go staticcheck error fixed


## [1.0.1] - 2024-05-28
### Added
- Config package for handling API configuration and request limits
- Singleton pattern for loading and accessing configuration
- Request limit middleware to enforce daily request limits

### Changed
- Integrated config settings in the main application
- Config file loading added to the application initialization process


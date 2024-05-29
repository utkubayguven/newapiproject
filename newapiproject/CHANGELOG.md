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



## [1.0.1] - 2024-05-29
### Added
- Configurations for Docker Compose to set up both the application and the PostgreSQL database.
- Added Dockerfile for app-container

### Changed
- Refactored Dockerfile to copy `config.json` correctly.
- Adjusted `GetConfig` function to print debug information about the configuration loading process.

### Fixed
- Resolved issues with Docker container not finding `GLIBC_2.34` and `GLIBC_2.32`.
- Corrected file paths for `.env` and `config.json` in the Docker setup.
- Addressed issue with database connection refusal in Docker setup.



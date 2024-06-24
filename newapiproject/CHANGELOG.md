bridge# Changelog

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



## [1.0.0] - 2024-05-28
### Added
- Config package for handling API configuration and request limits
- Singleton pattern for loading and accessing configuration
- Request limit middleware to enforce daily request limits

### Changed
- Integrated config settings in the main application
- Config file loading added to the application initialization process



## [1.0.0] - 2024-05-29
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


## [1.0.0] - 2024-05-31
### Added
- Added infra0 infra1 infra2 etcd node 

### Fixed
- Some error with creating node2 and node3



## [1.0.0] - 2024-06-04

### Changed
- Updated the etcd cluster setup based on the latest conversation and findings
#### Notes from conversation:
- I initially faced issues with local etcd cluster setup, referring to production documentation.
- Oğuzhan clarified the requirement for single device clustering.
- I successfully resolved the issue but needed to redo the setup after restarting.
- The setup was based on the clustering guide from the etcd documentation.

### Fixed
- Addressed issue with setting up etcd clusters by adding nodes one by one instead of all at once
- Fixed the need to redo the setup after restarting by updating the configuration process based on the etcd clustering guide



## [1.0.0] - 2024-06-05

### Added
- Configurations for a local etcd cluster setup
- Added amd64 etcd to local directory and configured it to run with ./
- Configured to run with goreman and used nano to write addresses in the Procfile

### Changed
- Updated the etcd cluster setup based on the latest findings
- Adjusted setup to add nodes sequentially to avoid mismatch errors between nodes

### Fixed
- Addressed issues with setting up etcd clusters by adding nodes one by one instead of all at once
- Resolved node mismatch issues by sequentially adding nodes to the cluster



## [1.0.0] - 2024-06-06

### Added
- Added Docker Compose setup for etcd with three containers: node1, node2, and node3

### Changed
- Set up etcd cluster within Docker Compose with separate containers for each node

### Fixed
- Tested etcd cluster setup by performing the following tests:
- Verified that data written to etcd1 can be read from etcd2
- Checked that the cluster continues to function when etcd3 is stopped
- Confirmed that data added to the cluster while one node is down is accessible when the node is brought back up



## [2.0.0] - 2024-06-24

### Added
- Added database package with functions for initializing and managing the ETCD client.
- Introduced database.InitEtcd(endpoints []string) to start the ETCD client.
- Implemented database.GetClient(endpoints []string) to return the existing ETCD client or create a new one.
- Created database.TestPutGet(client *clientv3.Client) to perform read and write tests on ETCD.

### Changed
- Updated Docker network configuration to ensure all ETCD containers are connected to the same network using the bridge driver.
- Manually connected ETCD containers to the etcd-net network to resolve connection issues.

### Fixed
- Resolved network configuration issues by ensuring all ETCD containers are connected to the same Docker network (etcd-net).
- Improved logging and error handling within handlers.Register to provide detailed error messages and user data logging.
- Ensured proper initialization and connection handling of the ETCD client in the main application.










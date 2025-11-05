# Roadmap to Version 1.0

## Current State

CSV Query is a working proof-of-concept that allows users to query CSV files using SQL through both a CLI tool and an HTTP server. The core functionality works, but the project is marked as "WORK IN PROGRESS" and needs additional polish, testing, and features to be production-ready.

**What Currently Works:**
- Query CSV files with SQL using an in-memory SQLite database
- Join multiple CSV files
- CLI tool (`csvquery`) for local file querying
- HTTP server (`csvquery-server`) for querying remote CSV files via URLs
- Basic test coverage
- Go module structure

## Critical Path to v1.0

### 1. Testing & Quality Assurance

**Priority: Critical**

- **Expand test coverage**
  - Add tests for CLI argument parsing (`cmd/csvquery/main.go:44-82`)
  - Add tests for HTTP server endpoints (`cmd/csvquery-server/main.go:26-88`)
  - Add edge case tests: empty CSV files, malformed CSV, missing headers
  - Add tests for error conditions in `Insert()` and `Query()` methods
  - Add tests for different CSV formats (different delimiters, quoted fields, multiline)
  - Target: >80% code coverage

- **Add integration tests**
  - End-to-end CLI tests with real CSV files
  - HTTP server integration tests with mock HTTP requests
  - Test multi-table joins with realistic data

- **Add benchmarks**
  - Performance tests for large CSV files (1MB, 10MB, 100MB)
  - Memory usage profiling
  - Concurrent query handling for server version

### 2. Error Handling & Validation

**Priority: Critical**

- **Improve input validation**
  - Validate table names to prevent SQL injection
  - Validate CSV structure (ensure consistent column counts)
  - Better error messages for invalid SQL queries
  - Validate file paths and URLs before processing

- **Enhanced error reporting**
  - Add context to errors (which file, which table, which row)
  - User-friendly error messages vs debug information
  - Proper HTTP status codes for server errors

- **Graceful failure handling**
  - Handle missing files gracefully
  - Handle network timeouts for server CSV fetching
  - Handle out-of-memory situations for large files

### 3. Documentation

**Priority: High**

- **User documentation**
  - Expand README with more examples
  - Document SQL dialect limitations (SQLite-specific features)
  - Add troubleshooting section
  - Document CSV format requirements
  - Add performance guidelines (file size limits, best practices)

- **API documentation**
  - Document HTTP server API endpoints and parameters
  - Add example requests/responses for server version
  - Document CORS behavior

- **Developer documentation**
  - Add godoc comments to all exported functions and types
  - Add architecture overview (how CSV data flows through the system)
  - Contributing guidelines
  - Code examples for using as a library

- **Examples directory**
  - Add more example CSV files with realistic scenarios
  - Add shell scripts demonstrating common use cases
  - Add example server queries

### 4. Security

**Priority: High** (especially for server version)

- **Server security hardening**
  - Implement rate limiting to prevent abuse
  - Add timeout for CSV downloads from remote URLs
  - Implement maximum file size limits
  - Add allowlist/blocklist for CSV URLs
  - Consider adding authentication/API keys
  - Improve CORS configuration (currently allows any origin)
  - Add request logging for security auditing

- **SQL injection prevention**
  - Audit current parameterization in `csvquery.go:106-112`
  - Add validation for table and column names
  - Document safe usage patterns

- **Resource limits**
  - Add memory limits for in-memory database
  - Add row count limits
  - Add query timeout limits

### 5. Features & Usability

**Priority: Medium**

- **CLI improvements**
  - Add `--version` flag
  - Add `--output` flag for different output formats (JSON, TSV, formatted table)
  - Support reading from stdin
  - Support multiple queries in one invocation
  - Better help text with examples
  - Add `--delimiter` flag for non-comma delimiters (TSV support)
  - Add `--no-header` flag for headerless CSV files

- **Server improvements**
  - Add health check endpoint
  - Add metrics endpoint (query count, latency, etc.)
  - Support POST requests for large queries
  - Add JSON output format option
  - Make listen address configurable (currently hardcoded to `localhost:8080`)
  - Add TLS/HTTPS support

- **Query enhancements**
  - Support for CSV files without headers
  - Auto-detect column types (currently all strings)
  - Support for different CSV dialects
  - Better handling of NULL values

### 6. Performance

**Priority: Medium**

- **Optimize for large files**
  - Streaming CSV parsing for large files
  - Consider pagination for large result sets
  - Add progress indicators for long operations

- **Caching (server version)**
  - Cache downloaded CSV files
  - Add ETags support for HTTP caching
  - Add TTL configuration for cache

- **Memory optimization**
  - Profile memory usage with large datasets
  - Consider using disk-backed SQLite instead of `:memory:` for large files
  - Stream results instead of loading all into memory

### 7. Build & Release Process

**Priority: Medium**

- **Cross-platform builds**
  - Automate builds for Linux, macOS, Windows
  - Test on all platforms
  - Document build requirements (CGO for SQLite)
  - Consider providing static binaries

- **Release automation**
  - Set up GitHub Actions for CI/CD
  - Automated testing on pull requests
  - Automated releases with version tagging
  - Generate release notes from commits

- **Distribution**
  - Publish to package managers (Homebrew, apt, etc.)
  - Provide Docker images for server version
  - Document installation methods

### 8. Configuration

**Priority: Low**

- **Configuration file support**
  - Support for `.csvqueryrc` configuration file
  - Environment variable support
  - Default table mappings
  - Server configuration (port, timeout, limits)

- **Server configuration**
  - Command-line flags for all server options
  - Configuration file for server settings
  - Environment-based configuration for cloud deployment

### 9. API Stability

**Priority: Critical**

- **Stabilize public API**
  - Review all exported functions and types
  - Ensure naming conventions are consistent
  - Mark internal functions as unexported
  - Version the API (semantic versioning)
  - Document breaking changes policy

- **Server API versioning**
  - Design URL structure for API versioning
  - Document API stability guarantees

### 10. Additional Features (Nice to Have)

**Priority: Low**

- **Advanced SQL features**
  - Support for CREATE VIEW
  - Support for temporary tables
  - Support for user-defined functions

- **Data type inference**
  - Auto-detect integers, floats, dates
  - Configurable type hints

- **Output formatting**
  - Pretty-printed tables
  - Colored output
  - Excel output format

- **Monitoring**
  - Structured logging
  - Integration with observability tools (Prometheus, etc.)

## Minimum Viable v1.0

To ship a minimal but solid v1.0, the following are **ESSENTIAL**:

1. **Testing**: Test coverage >70% with integration tests
2. **Documentation**: Comprehensive README and API docs
3. **Security**: Basic security hardening for server (rate limits, size limits, timeouts)
4. **Error Handling**: All error cases handled with clear messages
5. **CLI**: Add `--version`, `--help`, and `--output` flags
6. **Server**: Make listen address configurable
7. **Build**: Automated builds for major platforms
8. **API Stability**: Finalize public API and commit to semantic versioning

## Success Metrics for v1.0

- Test coverage >70%
- All exported functions have godoc comments
- Zero known critical bugs
- Works reliably with CSV files up to 100MB
- Server handles at least 10 concurrent requests
- Clear upgrade path documented
- Installation instructions for Linux, macOS, and Windows

## Timeline Estimate

**Assuming 1 developer working part-time:**

- Testing & Quality: 2-3 weeks
- Error Handling: 1 week
- Documentation: 1-2 weeks
- Security: 1-2 weeks
- Features & Usability: 2-3 weeks
- Performance: 1-2 weeks
- Build & Release: 1 week
- API Stability: 1 week

**Total: 10-17 weeks (2.5-4 months)**

## Post-v1.0 Considerations

After v1.0, consider:

- Plugin system for custom functions
- Support for other data formats (JSON, Parquet, etc.)
- Query optimization and explain plans
- Distributed querying for very large datasets
- Web UI for the server version
- Saved queries / query library
- Data transformation pipelines

package runner

// RunnerInfo merges Docker container info with optional GitHub runner status.
type RunnerInfo struct {
	Num          int
	Name         string
	ContainerID  string
	DockerState  string // running, exited, etc.
	DockerStatus string // human-readable Docker status

	// GitHub-sourced fields (populated when --github flag is used)
	GitHubID     int64
	GitHubStatus string // online, offline
	Busy         bool
}

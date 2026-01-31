package runner

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/lamtuanvu/gh-runner-ctl/internal/config"
	"github.com/lamtuanvu/gh-runner-ctl/internal/docker"
)

// Manager orchestrates runner lifecycle operations.
type Manager struct {
	Docker *docker.Client
	Config *config.Config
}

// NewManager creates a new runner manager.
func NewManager(cfg *config.Config, dc *docker.Client) *Manager {
	return &Manager{Config: cfg, Docker: dc}
}

// Up creates and starts `count` new runners, filling the lowest available numbers.
func (m *Manager) Up(ctx context.Context, count int) ([]string, error) {
	token, err := config.ResolveToken(m.Config.Token)
	if err != nil {
		return nil, err
	}

	existing, err := m.Docker.ListManagedContainers(ctx)
	if err != nil {
		return nil, err
	}

	var nums []int
	for _, c := range existing {
		nums = append(nums, c.Num)
	}

	newNums := NextNumbers(nums, count)
	var created []string
	for _, num := range newNums {
		name := fmt.Sprintf("%s-runner-%d", m.Config.Runners.NamePrefix, num)
		fmt.Printf("Creating %s...\n", name)
		id, err := m.Docker.CreateRunner(ctx, m.Config, num, token)
		if err != nil {
			return created, fmt.Errorf("creating runner %d: %w", num, err)
		}
		created = append(created, id)
		fmt.Printf("  Started %s (%s)\n", name, id[:12])
	}
	return created, nil
}

// Down stops and removes `count` runners, starting from the highest-numbered.
// If all is true, removes all managed runners.
func (m *Manager) Down(ctx context.Context, count int, all bool) error {
	existing, err := m.Docker.ListManagedContainers(ctx)
	if err != nil {
		return err
	}
	if len(existing) == 0 {
		fmt.Println("No managed runners found.")
		return nil
	}

	var targets []docker.RunnerContainer
	if all {
		targets = existing
	} else {
		var nums []int
		numToContainer := make(map[int]docker.RunnerContainer)
		for _, c := range existing {
			nums = append(nums, c.Num)
			numToContainer[c.Num] = c
		}
		toRemove := HighestNumbers(nums, count)
		for _, n := range toRemove {
			targets = append(targets, numToContainer[n])
		}
	}

	for _, c := range targets {
		fmt.Printf("Removing %s...\n", c.Name)
		if err := m.Docker.RemoveRunner(ctx, c.Name); err != nil {
			fmt.Printf("  Warning: %v\n", err)
			continue
		}
		fmt.Printf("  Removed %s\n", c.Name)
	}
	return nil
}

// Scale adjusts to exactly `target` runners.
func (m *Manager) Scale(ctx context.Context, target int) error {
	existing, err := m.Docker.ListManagedContainers(ctx)
	if err != nil {
		return err
	}

	current := len(existing)
	if current == target {
		fmt.Printf("Already at %d runners.\n", target)
		return nil
	}

	if current < target {
		diff := target - current
		fmt.Printf("Scaling up: %d -> %d (adding %d)\n", current, target, diff)
		_, err := m.Up(ctx, diff)
		return err
	}

	diff := current - target
	fmt.Printf("Scaling down: %d -> %d (removing %d)\n", current, target, diff)
	return m.Down(ctx, diff, false)
}

// List returns info about all managed runners.
func (m *Manager) List(ctx context.Context) ([]RunnerInfo, error) {
	containers, err := m.Docker.ListManagedContainers(ctx)
	if err != nil {
		return nil, err
	}

	var infos []RunnerInfo
	for _, c := range containers {
		infos = append(infos, RunnerInfo{
			Num:          c.Num,
			Name:         c.Name,
			ContainerID:  c.ID,
			DockerState:  c.State,
			DockerStatus: c.Status,
		})
	}
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Num < infos[j].Num
	})
	return infos, nil
}

// Stop stops runners. If nameOrNum is empty and all is true, stops all.
func (m *Manager) Stop(ctx context.Context, nameOrNum string, all bool) error {
	if all {
		return m.forEachManaged(ctx, "Stopping", func(ctx context.Context, c docker.RunnerContainer) error {
			return m.Docker.StopRunner(ctx, c.Name)
		})
	}
	c, err := m.findRunner(ctx, nameOrNum)
	if err != nil {
		return err
	}
	fmt.Printf("Stopping %s...\n", c.Name)
	return m.Docker.StopRunner(ctx, c.Name)
}

// Start starts stopped runners. If nameOrNum is empty and all is true, starts all.
func (m *Manager) Start(ctx context.Context, nameOrNum string, all bool) error {
	if all {
		return m.forEachManaged(ctx, "Starting", func(ctx context.Context, c docker.RunnerContainer) error {
			return m.Docker.StartRunner(ctx, c.Name)
		})
	}
	c, err := m.findRunner(ctx, nameOrNum)
	if err != nil {
		return err
	}
	fmt.Printf("Starting %s...\n", c.Name)
	return m.Docker.StartRunner(ctx, c.Name)
}

// Remove removes stopped runners. If nameOrNum is empty and all is true, removes all.
func (m *Manager) Remove(ctx context.Context, nameOrNum string, all bool) error {
	if all {
		return m.forEachManaged(ctx, "Removing", func(ctx context.Context, c docker.RunnerContainer) error {
			return m.Docker.RemoveRunner(ctx, c.Name)
		})
	}
	c, err := m.findRunner(ctx, nameOrNum)
	if err != nil {
		return err
	}
	fmt.Printf("Removing %s...\n", c.Name)
	return m.Docker.RemoveRunner(ctx, c.Name)
}

func (m *Manager) findRunner(ctx context.Context, nameOrNum string) (docker.RunnerContainer, error) {
	existing, err := m.Docker.ListManagedContainers(ctx)
	if err != nil {
		return docker.RunnerContainer{}, err
	}

	prefix := m.Config.Runners.NamePrefix
	for _, c := range existing {
		if c.Name == nameOrNum ||
			c.Name == prefix+"-runner-"+nameOrNum ||
			fmt.Sprintf("%d", c.Num) == nameOrNum ||
			strings.HasPrefix(c.ID, nameOrNum) {
			return c, nil
		}
	}
	return docker.RunnerContainer{}, fmt.Errorf("runner %q not found", nameOrNum)
}

func (m *Manager) forEachManaged(ctx context.Context, action string, fn func(context.Context, docker.RunnerContainer) error) error {
	existing, err := m.Docker.ListManagedContainers(ctx)
	if err != nil {
		return err
	}
	if len(existing) == 0 {
		fmt.Println("No managed runners found.")
		return nil
	}
	for _, c := range existing {
		fmt.Printf("%s %s...\n", action, c.Name)
		if err := fn(ctx, c); err != nil {
			fmt.Printf("  Warning: %v\n", err)
		}
	}
	return nil
}

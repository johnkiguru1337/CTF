package libcontainer

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/moby/sys/mountinfo"
	"github.com/opencontainers/runc/libcontainer/configs"
	"github.com/opencontainers/runc/libcontainer/utils"
	"github.com/opencontainers/runtime-spec/specs-go"

	"golang.org/x/sys/unix"
)

func TestFactoryNew(t *testing.T) {
	root := t.TempDir()
	factory, err := New(root)
	if err != nil {
		t.Fatal(err)
	}
	if factory == nil {
		t.Fatal("factory should not be nil")
	}
	lfactory, ok := factory.(*LinuxFactory)
	if !ok {
		t.Fatal("expected linux factory returned on linux based systems")
	}
	if lfactory.Root != root {
		t.Fatalf("expected factory root to be %q but received %q", root, lfactory.Root)
	}

	if factory.Type() != "libcontainer" {
		t.Fatalf("unexpected factory type: %q, expected %q", factory.Type(), "libcontainer")
	}
}

func TestFactoryNewIntelRdt(t *testing.T) {
	root := t.TempDir()
	factory, err := New(root, IntelRdtFs)
	if err != nil {
		t.Fatal(err)
	}
	if factory == nil {
		t.Fatal("factory should not be nil")
	}
	lfactory, ok := factory.(*LinuxFactory)
	if !ok {
		t.Fatal("expected linux factory returned on linux based systems")
	}
	if lfactory.Root != root {
		t.Fatalf("expected factory root to be %q but received %q", root, lfactory.Root)
	}

	if factory.Type() != "libcontainer" {
		t.Fatalf("unexpected factory type: %q, expected %q", factory.Type(), "libcontainer")
	}
}

func TestFactoryNewTmpfs(t *testing.T) {
	root := t.TempDir()
	factory, err := New(root, TmpfsRoot)
	if err != nil {
		t.Fatal(err)
	}
	if factory == nil {
		t.Fatal("factory should not be nil")
	}
	lfactory, ok := factory.(*LinuxFactory)
	if !ok {
		t.Fatal("expected linux factory returned on linux based systems")
	}
	if lfactory.Root != root {
		t.Fatalf("expected factory root to be %q but received %q", root, lfactory.Root)
	}

	if factory.Type() != "libcontainer" {
		t.Fatalf("unexpected factory type: %q, expected %q", factory.Type(), "libcontainer")
	}
	mounted, err := mountinfo.Mounted(lfactory.Root)
	if err != nil {
		t.Fatal(err)
	}
	if !mounted {
		t.Fatalf("Factory Root is not mounted")
	}
	mounts, err := mountinfo.GetMounts(mountinfo.SingleEntryFilter(lfactory.Root))
	if err != nil {
		t.Fatal(err)
	}
	if len(mounts) != 1 {
		t.Fatalf("Factory Root is not listed in mounts list")
	}
	m := mounts[0]
	if m.FSType != "tmpfs" {
		t.Fatalf("FSType of root: %s, expected %s", m.FSType, "tmpfs")
	}
	if m.Source != "tmpfs" {
		t.Fatalf("Source of root: %s, expected %s", m.Source, "tmpfs")
	}
	err = unix.Unmount(root, unix.MNT_DETACH)
	if err != nil {
		t.Error("failed to unmount root:", err)
	}
}

func TestFactoryLoadNotExists(t *testing.T) {
	factory, err := New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	_, err = factory.Load("nocontainer")
	if err == nil {
		t.Fatal("expected nil error loading non-existing container")
	}
	if !errors.Is(err, ErrNotExist) {
		t.Fatalf("expected ErrNotExist, got %v", err)
	}
}

func TestFactoryLoadContainer(t *testing.T) {
	root := t.TempDir()
	// setup default container config and state for mocking
	var (
		id            = "1"
		expectedHooks = configs.Hooks{
			configs.Prestart: configs.HookList{
				configs.CommandHook{Command: configs.Command{Path: "prestart-hook"}},
			},
			configs.Poststart: configs.HookList{
				configs.CommandHook{Command: configs.Command{Path: "poststart-hook"}},
			},
			configs.Poststop: configs.HookList{
				unserializableHook{},
				configs.CommandHook{Command: configs.Command{Path: "poststop-hook"}},
			},
		}
		expectedConfig = &configs.Config{
			Rootfs: "/mycontainer/root",
			Hooks:  expectedHooks,
			Cgroups: &configs.Cgroup{
				Resources: &configs.Resources{},
			},
		}
		expectedState = &State{
			BaseState: BaseState{
				InitProcessPid: 1024,
				Config:         *expectedConfig,
			},
		}
	)
	if err := os.Mkdir(filepath.Join(root, id), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := marshal(filepath.Join(root, id, stateFilename), expectedState); err != nil {
		t.Fatal(err)
	}
	factory, err := New(root, IntelRdtFs)
	if err != nil {
		t.Fatal(err)
	}
	container, err := factory.Load(id)
	if err != nil {
		t.Fatal(err)
	}
	if container.ID() != id {
		t.Fatalf("expected container id %q but received %q", id, container.ID())
	}
	config := container.Config()
	if config.Rootfs != expectedConfig.Rootfs {
		t.Fatalf("expected rootfs %q but received %q", expectedConfig.Rootfs, config.Rootfs)
	}
	expectedHooks[configs.Poststop] = expectedHooks[configs.Poststop][1:] // expect unserializable hook to be skipped
	if !reflect.DeepEqual(config.Hooks, expectedHooks) {
		t.Fatalf("expects hooks %q but received %q", expectedHooks, config.Hooks)
	}
	lcontainer, ok := container.(*linuxContainer)
	if !ok {
		t.Fatal("expected linux container on linux based systems")
	}
	if lcontainer.initProcess.pid() != expectedState.InitProcessPid {
		t.Fatalf("expected init pid %d but received %d", expectedState.InitProcessPid, lcontainer.initProcess.pid())
	}
}

func marshal(path string, v interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close() //nolint: errcheck
	return utils.WriteJSON(f, v)
}

type unserializableHook struct{}

func (unserializableHook) Run(*specs.State) error {
	return nil
}
